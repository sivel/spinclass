package common

import (
	openstack "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
)

type Config struct {
	Server struct {
		Port string
		Cert string
		Key  string
	}
	OpenStack struct {
		IdentityEndpoint string `yaml:"identity"`
		Username         string
		APIKey           string
		Password         string
		Region           string
		ImageRef         string `yaml:"image"`
		FlavorRef        string `yaml:"flavor"`
	}
}

type SpinUp struct {
	Prefix    string   `json:"prefix"`
	ServerIDs []string `json:"server_ids"`
}

type Odometer struct {
	Instances map[string]*openstack.Server `json:"instances"`
}
