package config

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

/**
全局常量
*/
const (
	BeginStr   = "P*"   //协议开头
	EndStr     = "*B*K" //协议结尾
	MinLength  = 30     //数据量最小个数
	MaxLength  = 100    //数据量最大个数
	MIN_TIME   = 3000   //最小时间
	MAX_TIME   = 1000 * 60 * 5
	ConfigFile = "config.yaml"
	ConfigEnv  = "GVA_CONFIG"
)

/**
全局变量
*/
var (
	GVA_VP *viper.Viper //配置

	GVA_REDIS *redis.Client

	GVA_MQTT mqtt.Client

	GVA_CONFIG Server
)
