package interfaces

import (
	"github.com/CuteReimu/FengSheng/utils"
)

var logger = utils.GetLogger("interfaces")

type IPlayer interface {
	Init(game IGame, location int)
	GetGame() IGame
	Location() int
	GetAlternativeLocation(location int) uint32
	NotifyAddHandCard(location int, unknownCount int, cards ...ICard)
	Draw(count int)
	NotifyDrawPhase(location int)
	NotifyMainPhase(location int, waitSecond uint32)
	IsAlive() bool
}

type BasePlayer struct {
	game     IGame
	location int
	cards    map[uint32]ICard
}

func (p *BasePlayer) Init(game IGame, location int) {
	p.game = game
	p.location = location
	p.cards = make(map[uint32]ICard)
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

func (p *BasePlayer) GetAlternativeLocation(location int) uint32 {
	location -= p.Location()
	totalPlayerCount := len(p.GetGame().GetPlayers())
	if location < 0 {
		location += totalPlayerCount
	}
	return uint32(location % totalPlayerCount)
}

func (p *BasePlayer) NotifyDrawPhase(int) {
}

func (p *BasePlayer) NotifyMainPhase(int, uint32) {
	p.GetGame().Post(p.GetGame().SendPhase)
}

func (p *BasePlayer) IsAlive() bool {
	return true
}
