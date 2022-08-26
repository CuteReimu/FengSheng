package interfaces

import (
	"github.com/CuteReimu/FengSheng/protos"
)

type Fsm interface {
	Resolve() (next Fsm, continueResolve bool)
}

type IGame interface {
	GetPlayers() []IPlayer
	GetDeck() IDeck
	GetWhoDie() int
	GetWhoseTurn() int
	GetCurrentCard() *CurrentCard
	SetCurrentCard(currentCard *CurrentCard)
	GetCurrentPhase() protos.Phase
	GetWhoseSendTurn() int
	SetWhoseSendTurn(whoseSendTurn int)
	GetWhoseFightTurn() int
	SetWhoseFightTurn(location int)
	GetMessageCardDirection() protos.Direction
	GetCurrentMessageCard() ICard
	SetCurrentMessageCard(currentMessageCard ICard)
	IsMessageCardFaceUp() bool
	SetMessageCardFaceUp(messageCardFaceUp bool)
	SetLockPlayers(locations []int)
	GetLockPlayers() []int
	IsIdleTimePoint() bool
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
	ContinueResolve()
}

type DieState int32

const (
	DieStateNone DieState = iota
	DieStateWaitForChengQing
	DieStateDying
)
