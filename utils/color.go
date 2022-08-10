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
	switch len(s) {
	case 1:
		return s + "色"
	case 2:
		return s + "双色"
	}
	panic(fmt.Sprint("unknown color: ", color))
}

func IdentityColorToString(color protos.Color) string {
	switch color {
	case protos.Color_Red:
		return "红方"
	case protos.Color_Blue:
		return "蓝方"
	case protos.Color_Black:
		return "神秘人"
	}
	panic(fmt.Sprint("unknown color: ", color))
}
