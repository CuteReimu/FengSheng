package game

import "github.com/CuteReimu/FengSheng/protos"

type CurrentCard struct {
	Card         ICard
	Player       int
	TargetPlayer int
}

type ICard interface {
	GetId() uint32
	GetType() protos.CardType
	GetColors() []protos.Color
	GetDirection() protos.Direction
	CanLock() bool
	CanUse(g *Game, user IPlayer, args ...interface{}) bool
	Execute(g *Game, user IPlayer, args ...interface{})
	ToPbCard() *protos.Card
	String() string
}

type BaseCard struct {
	Id        uint32
	Color     []protos.Color
	Direction protos.Direction
	Lockable  bool
}

func (card *BaseCard) GetId() uint32 {
	return card.Id
}

func (card *BaseCard) GetColors() []protos.Color {
	return card.Color
}

func (card *BaseCard) GetDirection() protos.Direction {
	return card.Direction
}

func (card *BaseCard) CanLock() bool {
	return card.Lockable
}

func (card *BaseCard) ToPbCard() *protos.Card {
	pb := &protos.Card{
		CardId:  card.Id,
		CardDir: card.Direction,
		CanLock: card.Lockable,
	}
	pb.CardColor = append(pb.CardColor, card.Color...)
	return pb
}
