package config

import (
	"fmt"
	"io/ioutil"

	// "github.com/bisohns/saido/driver"

	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
)

type HostList = []string
type Metrics = map[string]string

type DashboardInfo struct {
	Hosts        []Host
	Metrics      Metrics
	Title        string
	PollInterval int
}

func Contains(hostList HostList, host Host) bool {
	for _, compare := range hostList {
		if host.Address == compare {
			return true
		}
	}
	return false
}

func MergeMetrics(a, b Metrics) (metrics Metrics) {
	metrics = Metrics{}
	inputs := [2]Metrics{a, b}
	for _, metric := range inputs {
		for k, v := range metric {
			metrics[k] = v
		}
	}
	return
}

// GetAllHostAddresses : returns list of all hosts in the dashboard
func (dashboardInfo *DashboardInfo) GetAllHostAddresses() (addresses HostList) {
	addresses = []string{}
	for _, host := range dashboardInfo.Hosts {
		addresses = append(addresses, host.Address)
	}
	return
}

type Connection struct {
	Type                 string `mapstructure:"type"`
	Username             string `mapstructure:"username"`
	Password             string `mapstructure:"password"`
	PrivateKeyPath       string `mapstructure:"private_key_path"`
	PrivateKeyPassPhrase string `mapstructure:"private_key_passphrase"`
	Port                 int32  `mapstructure:"port"`
	Host                 string
}

type Host struct {
	Address    string
	Alias      string
	Connection *Connection
	// Metrics : extend global metrics with single metrics
	Metrics Metrics
}

type Config struct {
	Hosts        map[interface{}]interface{} `yaml:"hosts"`
	Metrics      map[interface{}]interface{} `yaml:"metrics"`
	Title        string                      `yaml:"title"`
	PollInterval int                         `yaml:"poll-interval"`
}

func LoadConfig(configPath string) *Config {
	var config = &Config{}
	confYaml, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Errorf("yamlFile.Get err   %v ", err)
	}
	err = yaml.Unmarshal([]byte(confYaml), &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return config
}

func coerceMetrics(rawMetrics map[interface{}]interface{}) map[string]string {
	metrics := make(map[string]string)
	for metric, customCommand := range rawMetrics {
		metric := fmt.Sprintf("%v", metric)
		metrics[metric] = fmt.Sprintf("%v", customCommand)
	}
	return metrics
}

func GetDashboardInfoConfig(config *Config) *DashboardInfo {
	dashboardInfo := &DashboardInfo{
		Title: "Saido",
	}
	if config.Title != "" {
		dashboardInfo.Title = config.Title
	}

	dashboardInfo.Hosts = parseConfig("root", "", config.Hosts, &Connection{})
	dashboardInfo.Metrics = coerceMetrics(config.Metrics)
	for _, host := range dashboardInfo.Hosts {
		log.Debugf("%s: %v", host.Address, host.Connection)
	}
	if config.PollInterval < 5 {
		log.Fatal("Cannot set poll interval below 5 seconds")
	}
	dashboardInfo.PollInterval = config.PollInterval
	return dashboardInfo
}

func parseConnection(conn map[interface{}]interface{}) *Connection {
	var c Connection
	mapstructure.Decode(conn, &c)
	if c.Type == "ssh" && c.Port == 0 {
		c.Port = 22
	}
	if c.Password != "" && c.PrivateKeyPath != "" {
		log.Fatal("Cannot specify both password login and private key login on same connection")
	}
	return &c
}

func parseConfig(name string, host string, group map[interface{}]interface{}, currentConnection *Connection) []Host {
	currentConn := currentConnection
	allHosts := []Host{}
	log.Debugf("Loading config for %s and host: %s with Connection: %+v", name, host, currentConn)
	isParent := false // Set to true for groups that contain just children data i.e children
	if conn, ok := group["connection"]; ok {
		v, ok := conn.(map[interface{}]interface{})
		if !ok {
			log.Fatalf("Failed to parse connection for %s", name)
		}

		currentConn = parseConnection(v)
	}

	if children, ok := group["children"]; ok {
		isParent = true
		parsedChildren, ok := children.(map[interface{}]interface{})
		if !ok {
			log.Fatalf("Failed to parse children of %s", name)
			return nil
		}

		for k, v := range parsedChildren {
			host := make(map[interface{}]interface{})
			host, ok := v.(map[interface{}]interface{})
			if !ok && v != nil { // some leaf nodes do not contain extra data under
				log.Errorf("Faled to parse children of %s", name)
			}
			allHosts = append(allHosts, parseConfig(fmt.Sprintf("%s:%s", name, k), fmt.Sprintf("%s", k), host, currentConn)...)
		}
	}

	if !isParent {
		currentConn.Host = host

		newHost := Host{
			Address:    host,
			Connection: currentConn,
		}
		if alias, ok := group["alias"]; ok {
			newHost.Alias = alias.(string)
		}
		if metrics, ok := group["metrics"]; ok {
			rawMetrics, ok := metrics.(map[interface{}]interface{})
			if !ok {
				log.Fatalf("Failed to parse metrics for %s", name)
			}
			individualMetrics := coerceMetrics(rawMetrics)
			newHost.Metrics = individualMetrics
		}

		allHosts = append(allHosts, newHost)
	}
	return allHosts
}
