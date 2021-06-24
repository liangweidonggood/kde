package main

import (
	"fmt"
	"kde/config"
	"kde/server"
)

func init() {
	config.GVA_VP = config.Viper()          // 初始化Viper
	config.GVA_REDIS = config.RedisClient() //初始化redis
	config.GVA_MQTT = config.MqttClient()   //初始化mqtt
}
func main() {
	go server.HttpServer()
	server.TCPServer(fmt.Sprintf(":%d", config.GVA_CONFIG.System.TcpPort))
}
