package config

import (
	"flag"

	"github.com/sivel/spinclass/common"
)

func Flags(config *common.Config) {
	flag.StringVar(&config.Server.Port, "port", config.Server.Port, "HOST:PORT to listen on, HOST not required to listen on all addresses")
	flag.StringVar(&config.Server.Cert, "cert", config.Server.Cert, "SSL cert file path. This option with 'key' enables SSL communication")
	flag.StringVar(&config.Server.Key, "key", config.Server.Key, "SSL key file path. This option with 'cert' enables SSL communication")
	flag.StringVar(&config.OpenStack.IdentityEndpoint, "identity", config.OpenStack.IdentityEndpoint, "OpenStack Identity V2 Endpoint URL")
	flag.StringVar(&config.OpenStack.Username, "username", config.OpenStack.Username, "OpenStack username")
	flag.StringVar(&config.OpenStack.Password, "password", config.OpenStack.Password, "OpenStack password")
	flag.StringVar(&config.OpenStack.APIKey, "apikey", config.OpenStack.APIKey, "OpenStack API Key")
	flag.StringVar(&config.OpenStack.Region, "region", config.OpenStack.Region, "OpenStack Region")
	flag.StringVar(&config.OpenStack.ImageRef, "image", config.OpenStack.ImageRef, "OpenStack Image ID")
	flag.StringVar(&config.OpenStack.FlavorRef, "flavor", config.OpenStack.FlavorRef, "OpenStack Flavor ID")
	flag.Parse()
}
