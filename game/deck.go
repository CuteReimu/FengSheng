package game

import (
	"github.com/CuteReimu/FengSheng/game/interfaces"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/davyxu/cellnet"
)

var DefaultDeck []interfaces.ICard

type Deck struct {
	game        interfaces.IGame
	cards       []interfaces.ICard
	discardPile []interfaces.ICard
}

func NewDeck(game interfaces.IGame) *Deck {
	d := &Deck{game: game, cards: append(make([]interfaces.ICard, 0), DefaultDeck...)}
	d.Shuffle()
	return d
}

func (d *Deck) Shuffle() {
	d.cards = append(d.cards, d.discardPile...)
	d.discardPile = nil
	d.game.GetRandom().Shuffle(len(d.cards), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
	d.notifyDeckCount(true)
}

func (d *Deck) Draw(n int) []interfaces.ICard {
	if n > len(d.cards) {
		d.Shuffle()
	}
	if n > len(d.cards) {
		n = len(d.cards)
	}
	result := d.cards[:n]
	d.cards = d.cards[n:]
	d.notifyDeckCount(false)
	return result
}

func (d *Deck) Discard(cards ...interfaces.ICard) {
	d.discardPile = append(d.discardPile, cards...)
}

func (d *Deck) GetDeckCount() int {
	return len(d.cards)
}

func (d *Deck) notifyDeckCount(shuffled bool) {
	for _, player := range d.game.GetPlayers() {
		if s, ok := player.(cellnet.Session); ok {
			s.Send(&protos.SyncDeckNumToc{
				Num:      uint32(len(d.cards)),
				Shuffled: shuffled,
			})
		}
	}
}
