package interfaces

import (
	_ "github.com/CuteReimu/FengSheng/core"
	"github.com/CuteReimu/FengSheng/protos"
	_ "github.com/davyxu/cellnet/peer/tcp"
	_ "github.com/davyxu/cellnet/proc/tcp"
)

type IGame interface {
	Start(totalCount, robotCount int)
	GetPlayers() []IPlayer
	GetDeck() IDeck
	GetWhoseTurn() int
	GetCurrentCard() ICard
	SetCurrentCard(card ICard)
	GetCurrentPhase() protos.Phase
	Post(callback func())
	DrawPhase()
	MainPhase()
	SendPhase()
	FightPhase()
	ReceivePhase()
	NextTurn()
}
