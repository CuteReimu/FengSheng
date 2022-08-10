package interfaces

import (
	"github.com/CuteReimu/FengSheng/utils"
)

var logger = utils.GetLogger("interfaces")

type IPlayer interface {
	Init(game IGame, location int)
	GetGame() IGame
	Location() int
	NotifyAddHandCard(card ...ICard)
	NotifyOtherAddHandCard(location int, count int)
	Draw(count int)
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

func (p *BasePlayer) NotifyAddHandCard(...ICard) {
}

func (p *BasePlayer) NotifyOtherAddHandCard(int, int) {
}

func (p *BasePlayer) Draw(count int) {
	cards := p.game.GetDeck().Draw(count)
	for _, card := range cards {
		p.cards[card.GetId()] = card
	}
	logger.Infof("%d号玩家摸了%v, 现在有%d张手牌", p.location, cards, len(p.cards))
	for _, player := range p.game.GetPlayers() {
		if player.Location() == p.Location() {
			player.NotifyAddHandCard(cards...)
		} else {
			player.NotifyOtherAddHandCard(p.location, len(cards))
		}
	}
}
