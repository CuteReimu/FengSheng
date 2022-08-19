package interfaces

import (
	"github.com/CuteReimu/FengSheng/protos"
	"math/rand"
)

type IGame interface {
	Start(totalCount, robotCount int)
	GetPlayers() []IPlayer
	GetDeck() IDeck
	GetRandom() *rand.Rand
	GetWhoDie() int
	GetWhoseTurn() int
	GetCurrentCard() *CurrentCard
	SetCurrentCard(currentCard *CurrentCard)
	GetCurrentPhase() protos.Phase
	GetWhoseSendTurn() int
	SetWhoseSendTurn(whoseSendTurn int)
	GetWhoseFightTurn() int
	GetMessageCardDirection() protos.Direction
	GetCurrentMessageCard() ICard
	SetCurrentMessageCard(currentMessageCard ICard)
	IsMessageCardFaceUp() bool
	SetMessageCardFaceUp(messageCardFaceUp bool)
	GetLockPlayers() []int
	IsIdleTimePoint() bool
	Post(callback func())
	DrawPhase()
	MainPhase()
	SendPhaseStart()
	OnSendCard(card ICard, dir protos.Direction, targetLocation int, lockLocations []int)
	MessageMoveNext()
	OnChooseReceiveCard()
	FightPhaseNext()
	ReceivePhase()
	AskForChengQing()
	AskNextForChengQing()
	AfterChengQing()
	NextTurn()
	PlayerDiscardCard(player IPlayer, cards ...ICard)
	GetDieState() DieState
}

type DieState int32

const (
	DieStateNone DieState = iota
	DieStateWaitForChengQing
	DieStateDying
)
