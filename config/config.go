package config

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/sivel/spinclass/common"
	"gopkg.in/yaml.v2"
)

func parseConfig() common.Config {
	var config common.Config
	configPath, err := filepath.Abs(".")
	configFile := filepath.Join(configPath, "spinclass.yaml")
	text, err := ioutil.ReadFile(configFile)
	if err == nil {
		yaml.Unmarshal(text, &config)
	}
	if config.Server.Port == "" {
		config.Server.Port = ":3000"
	}
	if config.OpenStack.IdentityEndpoint == "" {
		config.OpenStack.IdentityEndpoint = "https://identity.api.rackspacecloud.com/v2.0"
	}
	return config
}

func Config() common.Config {
	config := parseConfig()

	Flags(&config)

	if config.OpenStack.Username == "" || (config.OpenStack.Password == "" && config.OpenStack.APIKey == "") {
		log.Fatal("OpenStack Username and Password or APIKey are required")
	}
	if config.OpenStack.Region == "" || config.OpenStack.ImageRef == "" || config.OpenStack.FlavorRef == "" {
		log.Fatal("OpenStack Region, Image and Flavor are required")
	}

	return config
}
