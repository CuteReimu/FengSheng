package utils

import (
	"fmt"
	"github.com/CuteReimu/FengSheng/protos"
)

func CardColorToString(color ...protos.Color) string {
	var s string
	for _, c := range color {
		switch c {
		case protos.Color_Red:
			s += "红"
		case protos.Color_Blue:
			s += "蓝"
		case protos.Color_Black:
			s += "黑"
		}
	}
	switch len(color) {
	case 1:
		return s + "色"
	case 2:
		return s + "双色"
	}
	panic(fmt.Sprint("unknown color: ", color))
}

func IdentityColorToString(color protos.Color, task ...protos.SecretTask) string {
	switch color {
	case protos.Color_Red:
		return "红方"
	case protos.Color_Blue:
		return "蓝方"
	case protos.Color_Black:
		if len(task) == 0 {
			return "神秘人"
		}
		switch task[0] {
		case protos.SecretTask_Killer:
			return "神秘人[镇压者]"
		case protos.SecretTask_Stealer:
			return "神秘人[簒夺者]"
		case protos.SecretTask_Collector:
			return "神秘人[双重间谍]"
		}
	}
	panic(fmt.Sprint("unknown color: ", color))
}

func IsColorIn(color protos.Color, colors []protos.Color) bool {
	for _, c := range colors {
		if c == color {
			return true
		}
	}
	return false
}
