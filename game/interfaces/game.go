package interfaces

import (
	_ "github.com/CuteReimu/FengSheng/core"
	_ "github.com/davyxu/cellnet/peer/tcp"
	_ "github.com/davyxu/cellnet/proc/tcp"
)

type IGame interface {
	Start(totalCount, robotCount int)
	GetPlayers() []IPlayer
	GetDeck() IDeck
	GetWhoseTurn() int
	GetCurrentCard() ICard
}
