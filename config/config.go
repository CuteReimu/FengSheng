package config

import (
	"github.com/spf13/viper"
)

var GlobalConfig *viper.Viper

func init() {
	GlobalConfig = viper.New()
	GlobalConfig.SetConfigName("config")
	GlobalConfig.SetConfigType("yaml")
	GlobalConfig.AddConfigPath(".")
	err := GlobalConfig.ReadInConfig()
	if err != nil {
		GlobalConfig.Set("listen_address", "127.0.0.1:9091")
		GlobalConfig.Set("player.total_count", 5)
		GlobalConfig.Set("player.robot_count", 4)
		GlobalConfig.Set("log.tcp_debug_log", true)
	}
}
