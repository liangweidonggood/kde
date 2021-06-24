package config

type Server struct {
	Redis  Redis  `yaml:"redis"`
	Mqtt   Mqtt   `yaml:"mqtt"`
	System System `yaml:"system"`
}
