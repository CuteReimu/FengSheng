package interfaces

import (
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
)

var logger = utils.GetLogger("interfaces")

type IPlayer interface {
	Init(game IGame, location int, identity protos.Color, secretTask protos.SecretTask)
	GetGame() IGame
	Location() int
	GetAbstractLocation(location int) int
	GetAlternativeLocation(location int) uint32
	NotifyAddHandCard(location int, unknownCount int, cards ...ICard)
	Draw(count int)
	GetCards() map[uint32]ICard
	FindCard(cardId uint32) ICard
	DeleteCard(card uint32)
	NotifyDrawPhase()
	NotifyMainPhase(waitSecond uint32)
	IsAlive() bool
	SetIdentity(identity protos.Color, secretTask protos.SecretTask)
	GetIdentity() (protos.Color, protos.SecretTask)
	String() string
}

type BasePlayer struct {
	game       IGame
	location   int
	cards      map[uint32]ICard
	identity   protos.Color
	secretTask protos.SecretTask
}

func (p *BasePlayer) Init(game IGame, location int, identity protos.Color, secretTask protos.SecretTask) {
	logger.Info(p, "身份是", utils.IdentityColorToString(identity, secretTask))
	p.game = game
	p.location = location
	p.cards = make(map[uint32]ICard)
	p.SetIdentity(identity, secretTask)
}

func (p *BasePlayer) GetGame() IGame {
	return p.game
}

func (p *BasePlayer) Location() int {
	return p.location
}

func (p *BasePlayer) NotifyAddHandCard(int, int, ...ICard) {
}

func (p *BasePlayer) Draw(count int) {
	cards := p.game.GetDeck().Draw(count)
	for _, card := range cards {
		p.cards[card.GetId()] = card
	}
	logger.Infof("%d号玩家摸了%v, 现在有%d张手牌", p.location, cards, len(p.cards))
	for _, player := range p.game.GetPlayers() {
		if player.Location() == p.Location() {
			player.NotifyAddHandCard(p.Location(), 0, cards...)
		} else {
			player.NotifyAddHandCard(p.Location(), len(cards))
		}
	}
}

func (p *BasePlayer) GetAbstractLocation(location int) int {
	return (location + p.Location()) % len(p.GetGame().GetPlayers())
}

func (p *BasePlayer) GetAlternativeLocation(location int) uint32 {
	location -= p.Location()
	totalPlayerCount := len(p.GetGame().GetPlayers())
	if location < 0 {
		location += totalPlayerCount
	}
	return uint32(location % totalPlayerCount)
}

func (p *BasePlayer) NotifyDrawPhase() {
}

func (p *BasePlayer) IsAlive() bool {
	return true
}

func (p *BasePlayer) GetCards() map[uint32]ICard {
	return p.cards
}

func (p *BasePlayer) FindCard(cardId uint32) ICard {
	if card, ok := p.cards[cardId]; ok {
		return card
	}
	return nil
}

func (p *BasePlayer) DeleteCard(cardId uint32) {
	delete(p.cards, cardId)
}

func (p *BasePlayer) SetIdentity(identity protos.Color, secretTask protos.SecretTask) {
	p.identity = identity
	p.secretTask = secretTask
}

func (p *BasePlayer) GetIdentity() (protos.Color, protos.SecretTask) {
	return p.identity, p.secretTask
}
