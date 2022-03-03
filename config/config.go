package config

import (
	"fmt"
	"io/ioutil"

	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
)

type DashboardInfo struct {
	Hosts   []Host
	Metrics []string
	Title   string
}

type Connection struct {
	Type           string `mapstructure:"type"`
	Username       string `mapstructure:"type"`
	Password       string `mapstructure:"type"`
	PrivateKeyPath string `mapstructure:"private_key_path"`
	Port           int32  `mapstructure:"port"`
}

type Host struct {
	Address    string
	Alias      string
	Connection *Connection
}

type Config struct {
	Hosts   map[interface{}]interface{} `yaml:"hosts"`
	Metrics []string                    `yaml:"metrics"`
	Title   string                      `yaml:"title"`
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

func GetDashboardInfoConfig(config *Config) *DashboardInfo {
	dashboardInfo := &DashboardInfo{
		Title: "Saido",
	}
	if config.Title != "" {
		dashboardInfo.Title = config.Title
	}

	dashboardInfo.Hosts = parseConfig("root", "", config.Hosts, &Connection{})
	dashboardInfo.Metrics = config.Metrics
	return dashboardInfo
}

func parseConnection(conn map[interface{}]interface{}) *Connection {
	var c Connection
	mapstructure.Decode(conn, &c)
	return &c
}

func parseConfig(name string, host string, group map[interface{}]interface{}, currentConnection *Connection) []Host {
	currentConn := currentConnection
	allHosts := []Host{}
	log.Infof("Loading config for %s and host: %s with Connection: %+v", name, host, currentConn)
	isParent := false // Set to true for groups that contain just children data i.e children
	if conn, ok := group["connection"]; ok {
		v, ok := conn.(map[interface{}]interface{})
		if !ok {
			log.Errorf("Failed to parse connection for %s", name)
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
				log.Errorf("Failed to parse children of %s", name)
			}
			allHosts = append(allHosts, parseConfig(fmt.Sprintf("%s:%s", name, k), fmt.Sprintf("%s", k), host, currentConn)...)
		}
	}

	if !isParent {
		newHost := Host{
			Address:    host,
			Connection: currentConn,
		}

		allHosts = append(allHosts, newHost)
	}
	return allHosts
}
