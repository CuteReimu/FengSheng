package interfaces

import (
	_ "github.com/CuteReimu/FengSheng/core"
	"github.com/CuteReimu/FengSheng/protos"
	_ "github.com/davyxu/cellnet/peer/tcp"
	_ "github.com/davyxu/cellnet/proc/tcp"
	"math/rand"
)

type IGame interface {
	Start(totalCount, robotCount int)
	GetPlayers() []IPlayer
	GetDeck() IDeck
	GetRandom() *rand.Rand
	GetWhoseTurn() int
	GetCurrentCard() *CurrentCard
	SetCurrentCard(currentCard *CurrentCard)
	GetCurrentPhase() protos.Phase
	GetWhoseSendTurn() int
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
	NextTurn()
	PlayerDiscardCard(player IPlayer, cards ...ICard)
}
