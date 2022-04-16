package config

import "github.com/syth0le/microservice-load-balancer/sources/structures"

type Config struct {
	ProxyPort string              `json:"proxy_port"`
	Servers   []structures.Server `json:"servers"`
}
