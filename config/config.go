package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

var globalConfig *viper.Viper

const (
	listenAddress    = "listen_address"
	totalPlayerCount = "player.total_count"
	robotPlayerCount = "player.robot_count"
	tcpDebugLogOpen  = "log.tcp_debug_log"
)

func init() {
	globalConfig = viper.New()
	globalConfig.SetConfigName("application")
	globalConfig.SetConfigType("yaml")
	globalConfig.AddConfigPath(".")
	_ = globalConfig.ReadInConfig()
	var newConfig []string
	if !globalConfig.InConfig(listenAddress) {
		globalConfig.Set(listenAddress, "127.0.0.1:9091")
		newConfig = append(newConfig, listenAddress)
	}
	if !globalConfig.InConfig(totalPlayerCount) {
		globalConfig.Set(totalPlayerCount, 5)
		newConfig = append(newConfig, totalPlayerCount)
	}
	if !globalConfig.InConfig(robotPlayerCount) {
		globalConfig.Set(robotPlayerCount, 4)
		newConfig = append(newConfig, robotPlayerCount)
	}
	if !globalConfig.InConfig(tcpDebugLogOpen) {
		globalConfig.Set(tcpDebugLogOpen, true)
		newConfig = append(newConfig, tcpDebugLogOpen)
	}
	if len(newConfig) > 0 {
		if err := globalConfig.WriteConfigAs("application.yaml"); err != nil {
			panic(fmt.Sprintf("写入配置失败: %+v", err))
		}
		log.Println("在application.yaml有新增的配置：" + strings.Join(newConfig, ", ") + "，请检查后重新运行程序")
		b := make([]byte, 1)
		_, _ = os.Stdin.Read(b)
		os.Exit(0)
	}
}

// GetListenAddress 服务端监听端口号
func GetListenAddress() string {
	return globalConfig.GetString(listenAddress)
}

// GetTotalPlayerCount 总玩家人数
func GetTotalPlayerCount() int {
	return globalConfig.GetInt(totalPlayerCount)
}

// GetRobotPlayerCount 机器人人数
func GetRobotPlayerCount() int {
	return globalConfig.GetInt(robotPlayerCount)
}

// IsTcpDebugLogOpen 是否开启tcp调试日志
func IsTcpDebugLogOpen() bool {
	return globalConfig.GetBool(tcpDebugLogOpen)
}
