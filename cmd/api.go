/*
Copyright © 2021 Bisohns

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bisohns/saido/config"
	"github.com/bisohns/saido/driver"
	"github.com/bisohns/saido/inspector"
	"github.com/gorilla/handlers"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	socketBufferSize  = 1042
	messageBufferSize = 256
)

var (
	port     string
	server   = http.NewServeMux()
	upgrader = &websocket.Upgrader{
		ReadBufferSize:  socketBufferSize,
		WriteBufferSize: socketBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}
)

type FullMessage struct {
	Error   bool
	Message interface{}
}

type Message struct {
	Host     string
	Name     string
	Platform string
	Data     interface{}
}

type Client struct {
	Socket *websocket.Conn
	Send   chan *FullMessage
}

// Write to websocket
func (client *Client) Write() {
	defer client.Socket.Close()
	var err error
	for msg := range client.Send {
		err = client.Socket.WriteJSON(msg)
		if err != nil {
			log.Error("Error inside client write ", err)
		}
	}
}

type Hosts struct {
	Config *config.Config
	// Connections : hostname mapped to connection instances to reuse
	// across metrics
	mu      sync.Mutex
	Drivers map[string]*driver.Driver
	Client  chan *Client
	Start   chan bool
}

func (hosts *Hosts) getDriver(address string) *driver.Driver {
	hosts.mu.Lock()
	defer hosts.mu.Unlock()
	return hosts.Drivers[address]
}

func (hosts *Hosts) resetDriver(host config.Host) {
	hosts.mu.Lock()
	defer hosts.mu.Unlock()
	hostDriver := host.Connection.ToDriver()
	hosts.Drivers[host.Address] = &hostDriver
}

func (hosts *Hosts) sendMetric(host config.Host, client *Client) {
	if hosts.getDriver(host.Address) == nil {
		hosts.resetDriver(host)
	}
	for _, metric := range config.GetDashboardInfoConfig(hosts.Config).Metrics {
		driver := hosts.getDriver(host.Address)
		initializedMetric, err := inspector.Init(metric, driver)
		data, err := initializedMetric.Execute()
		if err == nil {
			var unmarsh interface{}
			json.Unmarshal(data, &unmarsh)
			message := &FullMessage{
				Message: Message{
					Host:     host.Address,
					Platform: (*driver).GetDetails().Name,
					Name:     metric,
					Data:     unmarsh,
				},
				Error: false,
			}
			client.Send <- message
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
			message := &FullMessage{
				Message: errorContent,
				Error:   true,
			}
			client.Send <- message
		}
	}
}

func (hosts *Hosts) Run() {
	dashboardInfo := config.GetDashboardInfoConfig(hosts.Config)
	log.Debug("In Running")
	for {
		select {
		case client := <-hosts.Client:
			for {
				for _, host := range dashboardInfo.Hosts {
					go hosts.sendMetric(host, client)
				}
				log.Infof("Delaying for %d seconds", dashboardInfo.PollInterval)
				time.Sleep(time.Duration(dashboardInfo.PollInterval) * time.Second)
			}
		}
	}

}

func (hosts *Hosts) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	client := &Client{
		Socket: socket,
		Send:   make(chan *FullMessage, messageBufferSize),
	}
	hosts.Client <- client
	client.Write()
}

func newHosts(cfg *config.Config) *Hosts {
	hosts := &Hosts{
		Config:  cfg,
		Drivers: make(map[string]*driver.Driver),
		Client:  make(chan *Client),
	}
	return hosts
}

func setHostHandler(w http.ResponseWriter, r *http.Request) {
	b, _ := json.Marshal("Hello World")
	w.Write(b)
}

var apiCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Run saido dashboard on a PORT env variable, fallback to set argument",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		//    server.HandleFunc("/set-hosts", SetHostHandler)
		// TODO: set up cfg using set-hosts endpoint
		hosts := newHosts(cfg)
		server.HandleFunc("/set-hosts", setHostHandler)
		server.Handle("/metrics", hosts)
		log.Info("listening on :", port)
		_, err := strconv.Atoi(port)
		if err != nil {
			log.Fatal(err)
		}
		go hosts.Run()
		loggedRouters := handlers.LoggingHandler(os.Stdout, server)
		if err := http.ListenAndServe(":"+port, loggedRouters); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	apiCmd.Flags().StringVarP(&port, "port", "p", "3000", "Port to run application server on")
	rootCmd.AddCommand(apiCmd)
}
