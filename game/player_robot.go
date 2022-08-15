package game

import (
	"github.com/CuteReimu/FengSheng/game/interfaces"
	"github.com/CuteReimu/FengSheng/protos"
	"strconv"
	"time"
)

var AIMainPhase = make(map[protos.CardType]func(player interfaces.IPlayer, card interfaces.ICard) bool)

type RobotPlayer struct {
	interfaces.BasePlayer
}

func (r *RobotPlayer) String() string {
	return strconv.Itoa(r.Location()) + "号[机器人]"
}

func (r *RobotPlayer) NotifyMainPhase(_ uint32) {
	if r.Location() != r.GetGame().GetWhoseTurn() {
		return
	}
	cards := r.GetCards()
	for cardId := range cards {
		card := cards[cardId]
		ai := AIMainPhase[card.GetType()]
		if ai != nil && ai(r, card) {
			return
		}
	}
	time.AfterFunc(time.Second, func() {
		r.GetGame().Post(r.GetGame().SendPhase)
	})
}
