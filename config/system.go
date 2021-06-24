package config

type System struct {
	TcpPort int64 `yaml:"tcpPort"`
	WebPort int64 `yaml:"webPort"`
}
