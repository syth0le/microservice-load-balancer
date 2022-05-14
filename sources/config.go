package sources

type Config struct {
	ProxyPort string   `json:"proxy_port"`
	Servers   []Server `json:"servers"`
}

var Cfg Config
var ServPool ServerPool
