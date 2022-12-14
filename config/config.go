package config

import (
	"fmt"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

var globalConfig *viper.Viper

const (
	listenAddress      = "listen_address"
	totalPlayerCount   = "player.total_count"
	tcpDebugLogOpen    = "log.tcp_debug_log"
	beginHandCardCount = "rule.hand_card_count_begin"
	turnHandCardCount  = "rule.hand_card_count_each_turn"
	gmEnable           = "gm.enable"
	gmListenAddress    = "gm.listen_address"
	gmDebugRoles       = "gm.debug_roles"
	clientVersion      = "version"
	roomCount          = "room_count"
)

func Init() {
	globalConfig = viper.New()
	globalConfig.SetConfigName("application")
	globalConfig.SetConfigType("yaml")
	globalConfig.AddConfigPath(".")
	_ = globalConfig.ReadInConfig()
	var newConfig []string
	initSingleConfig(&newConfig, listenAddress, "127.0.0.1:9091")
	initSingleConfig(&newConfig, totalPlayerCount, 5)
	initSingleConfig(&newConfig, tcpDebugLogOpen, true)
	initSingleConfig(&newConfig, beginHandCardCount, 3)
	initSingleConfig(&newConfig, turnHandCardCount, 3)
	initSingleConfig(&newConfig, gmEnable, false)
	initSingleConfig(&newConfig, gmListenAddress, "127.0.0.1:9092")
	initSingleConfig(&newConfig, clientVersion, uint32(1))
	initSingleConfig(&newConfig, roomCount, 200)
	initSingleConfig(&newConfig, gmDebugRoles, []int{int(protos.Role_duan_mu_jing), int(protos.Role_lao_bie)})
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

func initSingleConfig(newConfig *[]string, key string, value interface{}) {
	if !globalConfig.InConfig(key) {
		globalConfig.Set(key, value)
		*newConfig = append(*newConfig, key)
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

// IsTcpDebugLogOpen 是否开启tcp调试日志
func IsTcpDebugLogOpen() bool {
	return globalConfig.GetBool(tcpDebugLogOpen)
}

// GetHandCardCountBegin 游戏开始时摸牌数
func GetHandCardCountBegin() int {
	return globalConfig.GetInt(beginHandCardCount)
}

// GetHandCardCountEachTurn 每回合摸牌数
func GetHandCardCountEachTurn() int {
	return globalConfig.GetInt(turnHandCardCount)
}

// IsGmEnable 是否开启GM命令
func IsGmEnable() bool {
	return globalConfig.GetBool(gmEnable)
}

// GetGmListenAddress GM监听端口号
func GetGmListenAddress() string {
	return globalConfig.GetString(gmListenAddress)
}

// GetClientVersion 获取客户端版本号
func GetClientVersion() uint32 {
	return globalConfig.GetUint32(clientVersion)
}

// GetMaxRoomCount 获取最大房间数
func GetMaxRoomCount() int {
	return globalConfig.GetInt(roomCount)
}

// GetDebugRoles 获取调试角色
func GetDebugRoles() []int {
	return globalConfig.GetIntSlice(gmDebugRoles)
}
