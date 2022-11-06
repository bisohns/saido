package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"github.com/bisohns/saido/config"
	"github.com/bisohns/saido/driver"
	"github.com/bisohns/saido/inspector"
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	}}

type HostsController struct {
	Info *config.DashboardInfo
	// Connections : hostname mapped to connection instances to reuse
	// across metrics
	mu      sync.Mutex
	Drivers map[string]*driver.Driver
	// ReadOnlyHosts : restrict pinging every other server except these
	ReadOnlyHosts []string
	// ClientConnected : shows that a client is connected
	ClientConnected bool
	Client          chan *Client
	Received        chan *ReceiveMessage
	StopPolling     chan bool
}

func (hosts *HostsController) getDriver(address string) *driver.Driver {
	hosts.mu.Lock()
	defer hosts.mu.Unlock()
	return hosts.Drivers[address]
}

func (hosts *HostsController) resetDriver(host config.Host) {
	hosts.mu.Lock()
	defer hosts.mu.Unlock()
	hostDriver := host.Connection.ToDriver()
	hosts.Drivers[host.Address] = &hostDriver
}

func (hosts *HostsController) clientConnected() bool {
	hosts.mu.Lock()
	defer hosts.mu.Unlock()
	return hosts.ClientConnected
}

func (hosts *HostsController) setClientConnected(connected bool) {
	hosts.mu.Lock()
	defer hosts.mu.Unlock()
	hosts.ClientConnected = connected
}

func (hosts *HostsController) setReadOnlyHost(hostlist config.HostList) {
	hosts.mu.Lock()
	defer hosts.mu.Unlock()
	hosts.ReadOnlyHosts = hostlist
}

func (hosts *HostsController) sendMetric(host config.Host, client *Client) {
	if hosts.getDriver(host.Address) == nil {
		hosts.resetDriver(host)
	}
	for _, metric := range hosts.Info.Metrics {
		driver := hosts.getDriver(host.Address)
		initializedMetric, err := inspector.Init(metric, driver)
		data, err := initializedMetric.Execute()
		if err == nil {
			var unmarsh interface{}
			json.Unmarshal(data, &unmarsh)
			message := &SendMessage{
				Message: Message{
					Host:     host.Address,
					Platform: (*driver).GetDetails().Name,
					Name:     metric,
					Data:     unmarsh,
				},
				Error: false,
			}
			if config.Contains(hosts.ReadOnlyHosts, host) {
				client.Send <- message
			}
		} else {
			// check for error 127 which means command was not found
			var errorContent string
			if !strings.Contains(fmt.Sprintf("%s", err), "127") {
				errorContent = fmt.Sprintf("Could not retrieve metric %s from driver %s with error %s, resetting connection...", metric, host.Address, err)
			} else {
				errorContent = fmt.Sprintf("Command %s not found on driver %s", metric, host.Address)
			}
			log.Error(errorContent)
			hosts.resetDriver(host)
			message := &SendMessage{
				Message: ErrorMessage{
					Error: errorContent,
					Host:  host.Address,
				},
				Error: true,
			}
			client.Send <- message
		}
	}
}

func (hosts *HostsController) Poll(client *Client) {
	for {
		for _, host := range hosts.Info.Hosts {
			if !hosts.clientConnected() {
				return
			}
			if config.Contains(hosts.ReadOnlyHosts, host) {
				go hosts.sendMetric(host, client)
			}
		}
		log.Infof("Delaying for %d seconds", hosts.Info.PollInterval)
		time.Sleep(time.Duration(hosts.Info.PollInterval) * time.Second)
	}
}

func (hosts *HostsController) Run() {
	for {
		select {
		case client := <-hosts.Client:
			go hosts.Poll(client)
		case received := <-hosts.Received:
			if received.FilterBy == "" {
				hosts.setReadOnlyHost(hosts.Info.GetAllHostAddresses())
			} else {
				hosts.setReadOnlyHost([]string{received.FilterBy})
			}
		case poll := <-hosts.StopPolling:
			hosts.setClientConnected(!poll)
		}
	}

}

func (hosts *HostsController) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	client := &Client{
		Socket:          socket,
		Send:            make(chan *SendMessage, messageBufferSize),
		Received:        hosts.Received,
		StopHostPolling: hosts.StopPolling,
	}
	hosts.Client <- client
	hosts.StopPolling <- false
	go client.Write()
	client.Read()
}

// NewHostsController : initialze host controller with config file
func NewHostsController(cfg *config.Config) *HostsController {
	dashboardInfo := config.GetDashboardInfoConfig(cfg)
	hosts := &HostsController{
		Info:            dashboardInfo,
		Drivers:         make(map[string]*driver.Driver),
		ReadOnlyHosts:   dashboardInfo.GetAllHostAddresses(),
		Client:          make(chan *Client),
		Received:        make(chan *ReceiveMessage),
		StopPolling:     make(chan bool),
		ClientConnected: true,
	}
	return hosts
}
