package card

import (
	"github.com/CuteReimu/FengSheng/game"
	"github.com/CuteReimu/FengSheng/game/interfaces"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
)

var logger = utils.GetLogger("card")

func init() {
	for color := range protos.Color_name {
		game.DefaultDeck = append(game.DefaultDeck,
			&ShiTan{BaseCard: interfaces.BaseCard{Id: 1 + uint32(color)*6, Direction: protos.Direction_Right,
				Color: []protos.Color{protos.Color(color)}}, WhoDrawCard: []protos.Color{protos.Color_Black}},
			&ShiTan{BaseCard: interfaces.BaseCard{Id: 2 + uint32(color)*6, Direction: protos.Direction_Right,
				Color: []protos.Color{protos.Color(color)}}, WhoDrawCard: []protos.Color{protos.Color_Blue}},
			&ShiTan{BaseCard: interfaces.BaseCard{Id: 3 + uint32(color)*6, Direction: protos.Direction_Right,
				Color: []protos.Color{protos.Color(color)}}, WhoDrawCard: []protos.Color{protos.Color_Red, protos.Color_Black}},
			&ShiTan{BaseCard: interfaces.BaseCard{Id: 4 + uint32(color)*6, Direction: protos.Direction_Left,
				Color: []protos.Color{protos.Color(color)}}, WhoDrawCard: []protos.Color{protos.Color_Red, protos.Color_Blue}},
			&ShiTan{BaseCard: interfaces.BaseCard{Id: 5 + uint32(color)*6, Direction: protos.Direction_Left,
				Color: []protos.Color{protos.Color(color)}}, WhoDrawCard: []protos.Color{protos.Color_Blue, protos.Color_Black}},
			&ShiTan{BaseCard: interfaces.BaseCard{Id: 6 + uint32(color)*6, Direction: protos.Direction_Left,
				Color: []protos.Color{protos.Color(color)}}, WhoDrawCard: []protos.Color{protos.Color_Red}},
		)
	}
	game.DefaultDeck = append(game.DefaultDeck,
		&PingHeng{BaseCard: interfaces.BaseCard{Id: 19, Color: []protos.Color{protos.Color_Black}, Direction: protos.Direction_Left, Lockable: true}},
		&PingHeng{BaseCard: interfaces.BaseCard{Id: 20, Color: []protos.Color{protos.Color_Black}, Direction: protos.Direction_Right, Lockable: true}},
		&PingHeng{BaseCard: interfaces.BaseCard{Id: 21, Color: []protos.Color{protos.Color_Blue}, Direction: protos.Direction_Left, Lockable: true}},
		&PingHeng{BaseCard: interfaces.BaseCard{Id: 22, Color: []protos.Color{protos.Color_Red}, Direction: protos.Direction_Right, Lockable: true}},
		&PingHeng{BaseCard: interfaces.BaseCard{Id: 23, Color: []protos.Color{protos.Color_Red, protos.Color_Black}, Direction: protos.Direction_Up}},
		&PingHeng{BaseCard: interfaces.BaseCard{Id: 24, Color: []protos.Color{protos.Color_Blue, protos.Color_Black}, Direction: protos.Direction_Up}},
		&PingHeng{BaseCard: interfaces.BaseCard{Id: 25, Color: []protos.Color{protos.Color_Red, protos.Color_Black}, Direction: protos.Direction_Left}},
		&PingHeng{BaseCard: interfaces.BaseCard{Id: 26, Color: []protos.Color{protos.Color_Blue, protos.Color_Black}, Direction: protos.Direction_Right}},
	)
	game.DefaultDeck = append(game.DefaultDeck,
		&WeiBi{BaseCard: interfaces.BaseCard{Id: 27, Color: []protos.Color{protos.Color_Red}, Direction: protos.Direction_Left}},
		&WeiBi{BaseCard: interfaces.BaseCard{Id: 28, Color: []protos.Color{protos.Color_Red}, Direction: protos.Direction_Left}},
		&WeiBi{BaseCard: interfaces.BaseCard{Id: 29, Color: []protos.Color{protos.Color_Red}, Direction: protos.Direction_Left}},
		&WeiBi{BaseCard: interfaces.BaseCard{Id: 30, Color: []protos.Color{protos.Color_Red}, Direction: protos.Direction_Right}},
		&WeiBi{BaseCard: interfaces.BaseCard{Id: 31, Color: []protos.Color{protos.Color_Blue}, Direction: protos.Direction_Left}},
		&WeiBi{BaseCard: interfaces.BaseCard{Id: 32, Color: []protos.Color{protos.Color_Blue}, Direction: protos.Direction_Right}},
		&WeiBi{BaseCard: interfaces.BaseCard{Id: 33, Color: []protos.Color{protos.Color_Blue}, Direction: protos.Direction_Right}},
		&WeiBi{BaseCard: interfaces.BaseCard{Id: 34, Color: []protos.Color{protos.Color_Blue}, Direction: protos.Direction_Right}},
		&WeiBi{BaseCard: interfaces.BaseCard{Id: 35, Color: []protos.Color{protos.Color_Black}, Direction: protos.Direction_Left}},
		&WeiBi{BaseCard: interfaces.BaseCard{Id: 36, Color: []protos.Color{protos.Color_Black}, Direction: protos.Direction_Left}},
		&WeiBi{BaseCard: interfaces.BaseCard{Id: 37, Color: []protos.Color{protos.Color_Black}, Direction: protos.Direction_Right}},
		&WeiBi{BaseCard: interfaces.BaseCard{Id: 38, Color: []protos.Color{protos.Color_Black}, Direction: protos.Direction_Right}},
		&WeiBi{BaseCard: interfaces.BaseCard{Id: 39, Color: []protos.Color{protos.Color_Blue, protos.Color_Black}, Direction: protos.Direction_Left}},
		&WeiBi{BaseCard: interfaces.BaseCard{Id: 40, Color: []protos.Color{protos.Color_Red, protos.Color_Black}, Direction: protos.Direction_Right}},
	)
	for i := uint32(0); i < 3; i++ {
		game.DefaultDeck = append(game.DefaultDeck,
			&LiYou{BaseCard: interfaces.BaseCard{Id: 41 + i*2, Color: []protos.Color{protos.Color_Black}, Direction: protos.Direction_Left, Lockable: true}},
			&LiYou{BaseCard: interfaces.BaseCard{Id: 42 + i*2, Color: []protos.Color{protos.Color_Black}, Direction: protos.Direction_Right, Lockable: true}},
		)
	}
	game.DefaultDeck = append(game.DefaultDeck,
		&LiYou{BaseCard: interfaces.BaseCard{Id: 47, Color: []protos.Color{protos.Color_Blue}, Direction: protos.Direction_Left, Lockable: true}},
		&LiYou{BaseCard: interfaces.BaseCard{Id: 48, Color: []protos.Color{protos.Color_Red}, Direction: protos.Direction_Right, Lockable: true}},
	)
	for _, color := range []protos.Color{protos.Color_Red, protos.Color_Blue} {
		game.DefaultDeck = append(game.DefaultDeck,
			&ChengQing{BaseCard: interfaces.BaseCard{Id: 45 + uint32(color)*4, Color: []protos.Color{color}, Direction: protos.Direction_Up, Lockable: true}},
			&ChengQing{BaseCard: interfaces.BaseCard{Id: 46 + uint32(color)*4, Color: []protos.Color{color}, Direction: protos.Direction_Up, Lockable: true}},
			&ChengQing{BaseCard: interfaces.BaseCard{Id: 47 + uint32(color)*4, Color: []protos.Color{protos.Color_Black}, Direction: protos.Direction_Up, Lockable: true}},
			&ChengQing{BaseCard: interfaces.BaseCard{Id: 48 + uint32(color)*4, Color: []protos.Color{protos.Color_Black}, Direction: protos.Direction_Up, Lockable: true}},
		)
	}
	for _, direction := range []protos.Direction{protos.Direction_Left, protos.Direction_Right} {
		game.DefaultDeck = append(game.DefaultDeck,
			&PoYi{BaseCard: interfaces.BaseCard{Id: 52 + uint32(direction)*5, Color: []protos.Color{protos.Color_Red, protos.Color_Black}, Direction: direction, Lockable: true}},
			&PoYi{BaseCard: interfaces.BaseCard{Id: 53 + uint32(direction)*5, Color: []protos.Color{protos.Color_Blue, protos.Color_Black}, Direction: direction, Lockable: true}},
			&PoYi{BaseCard: interfaces.BaseCard{Id: 54 + uint32(direction)*5, Color: []protos.Color{protos.Color_Red}, Direction: direction, Lockable: true}},
			&PoYi{BaseCard: interfaces.BaseCard{Id: 55 + uint32(direction)*5, Color: []protos.Color{protos.Color_Blue}, Direction: direction, Lockable: true}},
			&PoYi{BaseCard: interfaces.BaseCard{Id: 56 + uint32(direction)*5, Color: []protos.Color{protos.Color_Black}, Direction: direction, Lockable: true}},
		)
	}
	game.DefaultDeck = append(game.DefaultDeck,
		&DiaoBao{BaseCard: interfaces.BaseCard{Id: 67, Color: []protos.Color{protos.Color_Red}, Direction: protos.Direction_Up}},
		&DiaoBao{BaseCard: interfaces.BaseCard{Id: 68, Color: []protos.Color{protos.Color_Red}, Direction: protos.Direction_Left}},
		&DiaoBao{BaseCard: interfaces.BaseCard{Id: 69, Color: []protos.Color{protos.Color_Red}, Direction: protos.Direction_Right}},
		&DiaoBao{BaseCard: interfaces.BaseCard{Id: 70, Color: []protos.Color{protos.Color_Blue}, Direction: protos.Direction_Up}},
		&DiaoBao{BaseCard: interfaces.BaseCard{Id: 71, Color: []protos.Color{protos.Color_Blue}, Direction: protos.Direction_Left}},
		&DiaoBao{BaseCard: interfaces.BaseCard{Id: 72, Color: []protos.Color{protos.Color_Blue}, Direction: protos.Direction_Right}},
		&DiaoBao{BaseCard: interfaces.BaseCard{Id: 73, Color: []protos.Color{protos.Color_Black}, Direction: protos.Direction_Left}},
		&DiaoBao{BaseCard: interfaces.BaseCard{Id: 74, Color: []protos.Color{protos.Color_Black}, Direction: protos.Direction_Right}},
		&DiaoBao{BaseCard: interfaces.BaseCard{Id: 75, Color: []protos.Color{protos.Color_Red, protos.Color_Black}, Direction: protos.Direction_Up}},
		&DiaoBao{BaseCard: interfaces.BaseCard{Id: 76, Color: []protos.Color{protos.Color_Red, protos.Color_Black}, Direction: protos.Direction_Right}},
		&DiaoBao{BaseCard: interfaces.BaseCard{Id: 77, Color: []protos.Color{protos.Color_Blue, protos.Color_Black}, Direction: protos.Direction_Up}},
		&DiaoBao{BaseCard: interfaces.BaseCard{Id: 78, Color: []protos.Color{protos.Color_Blue, protos.Color_Black}, Direction: protos.Direction_Left}},
	)
	game.DefaultDeck = append(game.DefaultDeck,
		&JieHuo{BaseCard: interfaces.BaseCard{Id: 79, Color: []protos.Color{protos.Color_Red}, Direction: protos.Direction_Up}},
		&JieHuo{BaseCard: interfaces.BaseCard{Id: 80, Color: []protos.Color{protos.Color_Red}, Direction: protos.Direction_Up}},
		&JieHuo{BaseCard: interfaces.BaseCard{Id: 81, Color: []protos.Color{protos.Color_Red}, Direction: protos.Direction_Up, Lockable: true}},
		&JieHuo{BaseCard: interfaces.BaseCard{Id: 82, Color: []protos.Color{protos.Color_Blue}, Direction: protos.Direction_Up}},
		&JieHuo{BaseCard: interfaces.BaseCard{Id: 83, Color: []protos.Color{protos.Color_Blue}, Direction: protos.Direction_Up}},
		&JieHuo{BaseCard: interfaces.BaseCard{Id: 84, Color: []protos.Color{protos.Color_Blue}, Direction: protos.Direction_Up, Lockable: true}},
		&JieHuo{BaseCard: interfaces.BaseCard{Id: 85, Color: []protos.Color{protos.Color_Black}, Direction: protos.Direction_Up}},
		&JieHuo{BaseCard: interfaces.BaseCard{Id: 86, Color: []protos.Color{protos.Color_Black}, Direction: protos.Direction_Up}},
		&JieHuo{BaseCard: interfaces.BaseCard{Id: 87, Color: []protos.Color{protos.Color_Black}, Direction: protos.Direction_Up, Lockable: true}},
		&JieHuo{BaseCard: interfaces.BaseCard{Id: 88, Color: []protos.Color{protos.Color_Black}, Direction: protos.Direction_Up, Lockable: true}},
		&JieHuo{BaseCard: interfaces.BaseCard{Id: 89, Color: []protos.Color{protos.Color_Red, protos.Color_Black}, Direction: protos.Direction_Up}},
		&JieHuo{BaseCard: interfaces.BaseCard{Id: 90, Color: []protos.Color{protos.Color_Blue, protos.Color_Black}, Direction: protos.Direction_Up}},
	)
	game.DefaultDeck = append(game.DefaultDeck,
		&WuDao{BaseCard: interfaces.BaseCard{Id: 91, Color: []protos.Color{protos.Color_Red}, Direction: protos.Direction_Up}},
		&WuDao{BaseCard: interfaces.BaseCard{Id: 92, Color: []protos.Color{protos.Color_Red}, Direction: protos.Direction_Left}},
		&WuDao{BaseCard: interfaces.BaseCard{Id: 93, Color: []protos.Color{protos.Color_Red}, Direction: protos.Direction_Right}},
		&WuDao{BaseCard: interfaces.BaseCard{Id: 94, Color: []protos.Color{protos.Color_Blue}, Direction: protos.Direction_Up}},
		&WuDao{BaseCard: interfaces.BaseCard{Id: 95, Color: []protos.Color{protos.Color_Blue}, Direction: protos.Direction_Left}},
		&WuDao{BaseCard: interfaces.BaseCard{Id: 96, Color: []protos.Color{protos.Color_Blue}, Direction: protos.Direction_Right}},
		&WuDao{BaseCard: interfaces.BaseCard{Id: 97, Color: []protos.Color{protos.Color_Black}, Direction: protos.Direction_Left}},
		&WuDao{BaseCard: interfaces.BaseCard{Id: 98, Color: []protos.Color{protos.Color_Black}, Direction: protos.Direction_Right}},
		&WuDao{BaseCard: interfaces.BaseCard{Id: 99, Color: []protos.Color{protos.Color_Blue, protos.Color_Black}, Direction: protos.Direction_Left}},
		&WuDao{BaseCard: interfaces.BaseCard{Id: 100, Color: []protos.Color{protos.Color_Red, protos.Color_Black}, Direction: protos.Direction_Right}},
	)
}
