package game

import (
	"github.com/CuteReimu/FengSheng/game/card"
	"github.com/CuteReimu/FengSheng/game/interfaces"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/davyxu/cellnet"
)

type Deck struct {
	game        *Game
	cards       []interfaces.ICard
	discardPile []interfaces.ICard
}

func NewDeck(game *Game) *Deck {
	d := &Deck{game: game}
	for color := range protos.Color_name {
		d.cards = append(d.cards,
			&card.ShiTan{BaseCard: interfaces.BaseCard{Id: 1 + uint32(color)*6, Direction: protos.Direction_Right,
				Color: []protos.Color{protos.Color(color)}}, WhoDrawCard: []protos.Color{protos.Color_Black}},
			&card.ShiTan{BaseCard: interfaces.BaseCard{Id: 2 + uint32(color)*6, Direction: protos.Direction_Right,
				Color: []protos.Color{protos.Color(color)}}, WhoDrawCard: []protos.Color{protos.Color_Blue}},
			&card.ShiTan{BaseCard: interfaces.BaseCard{Id: 3 + uint32(color)*6, Direction: protos.Direction_Right,
				Color: []protos.Color{protos.Color(color)}}, WhoDrawCard: []protos.Color{protos.Color_Red, protos.Color_Black}},
			&card.ShiTan{BaseCard: interfaces.BaseCard{Id: 4 + uint32(color)*6, Direction: protos.Direction_Left,
				Color: []protos.Color{protos.Color(color)}}, WhoDrawCard: []protos.Color{protos.Color_Red, protos.Color_Blue}},
			&card.ShiTan{BaseCard: interfaces.BaseCard{Id: 5 + uint32(color)*6, Direction: protos.Direction_Left,
				Color: []protos.Color{protos.Color(color)}}, WhoDrawCard: []protos.Color{protos.Color_Blue, protos.Color_Black}},
			&card.ShiTan{BaseCard: interfaces.BaseCard{Id: 6 + uint32(color)*6, Direction: protos.Direction_Left,
				Color: []protos.Color{protos.Color(color)}}, WhoDrawCard: []protos.Color{protos.Color_Red}},
		)
	}
	d.cards = append(d.cards,
		&card.PingHeng{BaseCard: interfaces.BaseCard{Id: 19, Color: []protos.Color{protos.Color_Black}, Direction: protos.Direction_Left, Lockable: true}},
		&card.PingHeng{BaseCard: interfaces.BaseCard{Id: 20, Color: []protos.Color{protos.Color_Black}, Direction: protos.Direction_Right, Lockable: true}},
		&card.PingHeng{BaseCard: interfaces.BaseCard{Id: 21, Color: []protos.Color{protos.Color_Blue}, Direction: protos.Direction_Left, Lockable: true}},
		&card.PingHeng{BaseCard: interfaces.BaseCard{Id: 22, Color: []protos.Color{protos.Color_Red}, Direction: protos.Direction_Right, Lockable: true}},
		&card.PingHeng{BaseCard: interfaces.BaseCard{Id: 23, Color: []protos.Color{protos.Color_Red, protos.Color_Black}, Direction: protos.Direction_Up}},
		&card.PingHeng{BaseCard: interfaces.BaseCard{Id: 24, Color: []protos.Color{protos.Color_Blue, protos.Color_Black}, Direction: protos.Direction_Up}},
		&card.PingHeng{BaseCard: interfaces.BaseCard{Id: 25, Color: []protos.Color{protos.Color_Red, protos.Color_Black}, Direction: protos.Direction_Left}},
		&card.PingHeng{BaseCard: interfaces.BaseCard{Id: 26, Color: []protos.Color{protos.Color_Blue, protos.Color_Black}, Direction: protos.Direction_Right}},
	)
	d.cards = append(d.cards,
		&card.WeiBi{BaseCard: interfaces.BaseCard{Id: 27, Color: []protos.Color{protos.Color_Red}, Direction: protos.Direction_Left}},
		&card.WeiBi{BaseCard: interfaces.BaseCard{Id: 28, Color: []protos.Color{protos.Color_Red}, Direction: protos.Direction_Left}},
		&card.WeiBi{BaseCard: interfaces.BaseCard{Id: 29, Color: []protos.Color{protos.Color_Red}, Direction: protos.Direction_Left}},
		&card.WeiBi{BaseCard: interfaces.BaseCard{Id: 30, Color: []protos.Color{protos.Color_Red}, Direction: protos.Direction_Right}},
		&card.WeiBi{BaseCard: interfaces.BaseCard{Id: 31, Color: []protos.Color{protos.Color_Blue}, Direction: protos.Direction_Left}},
		&card.WeiBi{BaseCard: interfaces.BaseCard{Id: 32, Color: []protos.Color{protos.Color_Blue}, Direction: protos.Direction_Right}},
		&card.WeiBi{BaseCard: interfaces.BaseCard{Id: 33, Color: []protos.Color{protos.Color_Blue}, Direction: protos.Direction_Right}},
		&card.WeiBi{BaseCard: interfaces.BaseCard{Id: 34, Color: []protos.Color{protos.Color_Blue}, Direction: protos.Direction_Right}},
		&card.WeiBi{BaseCard: interfaces.BaseCard{Id: 35, Color: []protos.Color{protos.Color_Black}, Direction: protos.Direction_Left}},
		&card.WeiBi{BaseCard: interfaces.BaseCard{Id: 36, Color: []protos.Color{protos.Color_Black}, Direction: protos.Direction_Left}},
		&card.WeiBi{BaseCard: interfaces.BaseCard{Id: 37, Color: []protos.Color{protos.Color_Black}, Direction: protos.Direction_Right}},
		&card.WeiBi{BaseCard: interfaces.BaseCard{Id: 38, Color: []protos.Color{protos.Color_Black}, Direction: protos.Direction_Right}},
		&card.WeiBi{BaseCard: interfaces.BaseCard{Id: 39, Color: []protos.Color{protos.Color_Blue, protos.Color_Black}, Direction: protos.Direction_Left}},
		&card.WeiBi{BaseCard: interfaces.BaseCard{Id: 40, Color: []protos.Color{protos.Color_Red, protos.Color_Black}, Direction: protos.Direction_Right}},
	)
	for i := uint32(0); i < 3; i++ {
		d.cards = append(d.cards,
			&card.LiYou{BaseCard: interfaces.BaseCard{Id: 41 + i*2, Color: []protos.Color{protos.Color_Black}, Direction: protos.Direction_Left, Lockable: true}},
			&card.LiYou{BaseCard: interfaces.BaseCard{Id: 42 + i*2, Color: []protos.Color{protos.Color_Black}, Direction: protos.Direction_Right, Lockable: true}},
		)
	}
	d.cards = append(d.cards,
		&card.LiYou{BaseCard: interfaces.BaseCard{Id: 47, Color: []protos.Color{protos.Color_Blue}, Direction: protos.Direction_Left, Lockable: true}},
		&card.LiYou{BaseCard: interfaces.BaseCard{Id: 48, Color: []protos.Color{protos.Color_Red}, Direction: protos.Direction_Right, Lockable: true}},
	)
	for _, color := range []protos.Color{protos.Color_Red, protos.Color_Blue} {
		d.cards = append(d.cards,
			&card.ChengQing{BaseCard: interfaces.BaseCard{Id: 45 + uint32(color)*4, Color: []protos.Color{color}, Direction: protos.Direction_Up, Lockable: true}},
			&card.ChengQing{BaseCard: interfaces.BaseCard{Id: 46 + uint32(color)*4, Color: []protos.Color{color}, Direction: protos.Direction_Up, Lockable: true}},
			&card.ChengQing{BaseCard: interfaces.BaseCard{Id: 47 + uint32(color)*4, Color: []protos.Color{protos.Color_Black}, Direction: protos.Direction_Up, Lockable: true}},
			&card.ChengQing{BaseCard: interfaces.BaseCard{Id: 48 + uint32(color)*4, Color: []protos.Color{protos.Color_Black}, Direction: protos.Direction_Up, Lockable: true}},
		)
	}
	for _, direction := range []protos.Direction{protos.Direction_Left, protos.Direction_Right} {
		d.cards = append(d.cards,
			&card.PoYi{BaseCard: interfaces.BaseCard{Id: 52 + uint32(direction)*5, Color: []protos.Color{protos.Color_Red, protos.Color_Black}, Direction: direction, Lockable: true}},
			&card.PoYi{BaseCard: interfaces.BaseCard{Id: 53 + uint32(direction)*5, Color: []protos.Color{protos.Color_Blue, protos.Color_Black}, Direction: direction, Lockable: true}},
			&card.PoYi{BaseCard: interfaces.BaseCard{Id: 54 + uint32(direction)*5, Color: []protos.Color{protos.Color_Red}, Direction: direction, Lockable: true}},
			&card.PoYi{BaseCard: interfaces.BaseCard{Id: 55 + uint32(direction)*5, Color: []protos.Color{protos.Color_Blue}, Direction: direction, Lockable: true}},
			&card.PoYi{BaseCard: interfaces.BaseCard{Id: 56 + uint32(direction)*5, Color: []protos.Color{protos.Color_Black}, Direction: direction, Lockable: true}},
		)
	}
	d.cards = append(d.cards,
		&card.DiaoBao{BaseCard: interfaces.BaseCard{Id: 67, Color: []protos.Color{protos.Color_Red}, Direction: protos.Direction_Up}},
		&card.DiaoBao{BaseCard: interfaces.BaseCard{Id: 68, Color: []protos.Color{protos.Color_Red}, Direction: protos.Direction_Left}},
		&card.DiaoBao{BaseCard: interfaces.BaseCard{Id: 69, Color: []protos.Color{protos.Color_Red}, Direction: protos.Direction_Right}},
		&card.DiaoBao{BaseCard: interfaces.BaseCard{Id: 70, Color: []protos.Color{protos.Color_Blue}, Direction: protos.Direction_Up}},
		&card.DiaoBao{BaseCard: interfaces.BaseCard{Id: 71, Color: []protos.Color{protos.Color_Blue}, Direction: protos.Direction_Left}},
		&card.DiaoBao{BaseCard: interfaces.BaseCard{Id: 72, Color: []protos.Color{protos.Color_Blue}, Direction: protos.Direction_Right}},
		&card.DiaoBao{BaseCard: interfaces.BaseCard{Id: 73, Color: []protos.Color{protos.Color_Black}, Direction: protos.Direction_Left}},
		&card.DiaoBao{BaseCard: interfaces.BaseCard{Id: 74, Color: []protos.Color{protos.Color_Black}, Direction: protos.Direction_Right}},
		&card.DiaoBao{BaseCard: interfaces.BaseCard{Id: 75, Color: []protos.Color{protos.Color_Red, protos.Color_Black}, Direction: protos.Direction_Up}},
		&card.DiaoBao{BaseCard: interfaces.BaseCard{Id: 76, Color: []protos.Color{protos.Color_Red, protos.Color_Black}, Direction: protos.Direction_Right}},
		&card.DiaoBao{BaseCard: interfaces.BaseCard{Id: 77, Color: []protos.Color{protos.Color_Blue, protos.Color_Black}, Direction: protos.Direction_Up}},
		&card.DiaoBao{BaseCard: interfaces.BaseCard{Id: 78, Color: []protos.Color{protos.Color_Blue, protos.Color_Black}, Direction: protos.Direction_Left}},
	)
	d.cards = append(d.cards,
		&card.JieHuo{BaseCard: interfaces.BaseCard{Id: 79, Color: []protos.Color{protos.Color_Red}, Direction: protos.Direction_Up}},
		&card.JieHuo{BaseCard: interfaces.BaseCard{Id: 80, Color: []protos.Color{protos.Color_Red}, Direction: protos.Direction_Up}},
		&card.JieHuo{BaseCard: interfaces.BaseCard{Id: 81, Color: []protos.Color{protos.Color_Red}, Direction: protos.Direction_Up, Lockable: true}},
		&card.JieHuo{BaseCard: interfaces.BaseCard{Id: 82, Color: []protos.Color{protos.Color_Blue}, Direction: protos.Direction_Up}},
		&card.JieHuo{BaseCard: interfaces.BaseCard{Id: 83, Color: []protos.Color{protos.Color_Blue}, Direction: protos.Direction_Up}},
		&card.JieHuo{BaseCard: interfaces.BaseCard{Id: 84, Color: []protos.Color{protos.Color_Blue}, Direction: protos.Direction_Up, Lockable: true}},
		&card.JieHuo{BaseCard: interfaces.BaseCard{Id: 85, Color: []protos.Color{protos.Color_Black}, Direction: protos.Direction_Up}},
		&card.JieHuo{BaseCard: interfaces.BaseCard{Id: 86, Color: []protos.Color{protos.Color_Black}, Direction: protos.Direction_Up}},
		&card.JieHuo{BaseCard: interfaces.BaseCard{Id: 87, Color: []protos.Color{protos.Color_Black}, Direction: protos.Direction_Up, Lockable: true}},
		&card.JieHuo{BaseCard: interfaces.BaseCard{Id: 88, Color: []protos.Color{protos.Color_Black}, Direction: protos.Direction_Up, Lockable: true}},
		&card.JieHuo{BaseCard: interfaces.BaseCard{Id: 89, Color: []protos.Color{protos.Color_Red, protos.Color_Black}, Direction: protos.Direction_Up}},
		&card.JieHuo{BaseCard: interfaces.BaseCard{Id: 90, Color: []protos.Color{protos.Color_Blue, protos.Color_Black}, Direction: protos.Direction_Up}},
	)
	d.cards = append(d.cards,
		&card.WuDao{BaseCard: interfaces.BaseCard{Id: 91, Color: []protos.Color{protos.Color_Red}, Direction: protos.Direction_Up}},
		&card.WuDao{BaseCard: interfaces.BaseCard{Id: 92, Color: []protos.Color{protos.Color_Red}, Direction: protos.Direction_Left}},
		&card.WuDao{BaseCard: interfaces.BaseCard{Id: 93, Color: []protos.Color{protos.Color_Red}, Direction: protos.Direction_Right}},
		&card.WuDao{BaseCard: interfaces.BaseCard{Id: 94, Color: []protos.Color{protos.Color_Blue}, Direction: protos.Direction_Up}},
		&card.WuDao{BaseCard: interfaces.BaseCard{Id: 95, Color: []protos.Color{protos.Color_Blue}, Direction: protos.Direction_Left}},
		&card.WuDao{BaseCard: interfaces.BaseCard{Id: 96, Color: []protos.Color{protos.Color_Blue}, Direction: protos.Direction_Right}},
		&card.WuDao{BaseCard: interfaces.BaseCard{Id: 97, Color: []protos.Color{protos.Color_Black}, Direction: protos.Direction_Left}},
		&card.WuDao{BaseCard: interfaces.BaseCard{Id: 98, Color: []protos.Color{protos.Color_Black}, Direction: protos.Direction_Right}},
		&card.WuDao{BaseCard: interfaces.BaseCard{Id: 99, Color: []protos.Color{protos.Color_Blue, protos.Color_Black}, Direction: protos.Direction_Left}},
		&card.WuDao{BaseCard: interfaces.BaseCard{Id: 100, Color: []protos.Color{protos.Color_Red, protos.Color_Black}, Direction: protos.Direction_Right}},
	)
	d.Shuffle()
	return d
}

func (d *Deck) Shuffle() {
	d.cards = append(d.cards, d.discardPile...)
	d.discardPile = nil
	d.game.Random.Shuffle(len(d.cards), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
	d.notifyDeckCount(true)
}

func (d *Deck) Draw(n int) []interfaces.ICard {
	if n > len(d.cards) {
		d.Shuffle()
	}
	if n > len(d.cards) {
		n = len(d.cards)
	}
	result := d.cards[:n]
	d.cards = d.cards[n:]
	d.notifyDeckCount(false)
	return result
}

func (d *Deck) Discard(cards ...interfaces.ICard) {
	d.discardPile = append(d.discardPile, cards...)
}

func (d *Deck) GetDeckCount() int {
	return len(d.cards)
}

func (d *Deck) notifyDeckCount(shuffled bool) {
	for _, player := range d.game.Players {
		if s, ok := player.(cellnet.Session); ok {
			s.Send(&protos.SyncDeckNumToc{
				Num:      uint32(len(d.cards)),
				Shuffled: shuffled,
			})
		}
	}
}
