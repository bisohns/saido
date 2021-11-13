package config

import (
	"io/ioutil"
	"strings"

	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
)

var config = &Config{}

type Connection struct {
	Type           string
	Username       string
	Password       string
	PrivateKeyPath string
}

// Group represents group specified under children
type Group struct {
	Name       string
	Hosts      map[string]*Host
	Children   map[string]*Group
	Connection *Connection
}

type Host struct {
	Alias        string
	Port         int32
	Connection   *Connection
	Hosts        map[string]*Host
	directGroups map[string]*Group
}

type Config struct {
	Hosts   map[string]interface{} `yaml:"hosts"`
	Metrics []string               `yaml:"metrics"`
}

func LoadConfig(configPath string) *Config {
	confYaml, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Errorf("yamlFile.Get err   %v ", err)
	}
	err = yaml.Unmarshal([]byte(confYaml), &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	for k, v := range config.Hosts {
		// Parser Children
		if strings.ToLower(k) == "children" {

		} else {

		}
		log.Info(k, v)
	}

	return config
}

func GetConfig() *Config {
	return config
}

func parseHosts(h interface{}) []*Host {
	return []*Host{}
}
