package game

import (
	"github.com/CuteReimu/FengSheng/protos"
)

type IPlayer interface {
	Init(game *Game, location int)
	Location() int
	NotifyAddHandCard(card ...*protos.Card)
	NotifyOtherAddHandCard(location int, count int)
	Draw(count int)
}

type basePlayer struct {
	game     *Game
	location int
	cards    map[uint32]*protos.Card
}

func (p *basePlayer) Init(game *Game, location int) {
	p.game = game
	p.location = location
	p.cards = make(map[uint32]*protos.Card)
}

func (p *basePlayer) Location() int {
	return p.location
}

func (p *basePlayer) NotifyAddHandCard(...*protos.Card) {
}

func (p *basePlayer) NotifyOtherAddHandCard(int, int) {
}

func (p *basePlayer) GetNextPlayer(location int) IPlayer {
	if !p.game.Dir {
		location = -location
	}
	location += p.location
	if location < 0 {
		location += p.game.TotalPlayerCount
	}
	location %= p.game.TotalPlayerCount
	return p.game.Players[location]
}

func (p *basePlayer) NotifyWin(int) {
}

func (p *basePlayer) Draw(count int) {
	cards := p.game.Deck.Draw(count)
	for _, card := range cards {
		p.cards[card.GetCardId()] = card
	}
	logger.Infof("%d号玩家摸了%d张牌, 现在还有%d张牌", p.location, count, len(p.cards))
	for _, player := range p.game.Players {
		if player.Location() == p.Location() {
			player.NotifyAddHandCard(cards...)
		} else {
			player.NotifyOtherAddHandCard(p.location, len(cards))
		}
	}
}
