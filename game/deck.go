package game

import (
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
)

var DefaultDeck []ICard

type Deck struct {
	nextId      uint32
	game        *Game
	cards       []ICard
	discardPile []ICard
}

func NewDeck(game *Game) *Deck {
	d := &Deck{game: game, cards: append(make([]ICard, 0), DefaultDeck...)}
	d.nextId = uint32(len(d.cards)) - 1
	d.Shuffle()
	return d
}

func (d *Deck) Shuffle() {
	d.cards = append(d.cards, d.discardPile...)
	d.discardPile = nil
	utils.Random.Shuffle(len(d.cards), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
	d.notifyDeckCount(true)
}

func (d *Deck) Draw(n int) []ICard {
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

func (d *Deck) Discard(cards ...ICard) {
	d.discardPile = append(d.discardPile, cards...)
}

func (d *Deck) GetDeckCount() int {
	return len(d.cards)
}

func (d *Deck) notifyDeckCount(shuffled bool) {
	for _, player := range d.game.GetPlayers() {
		if s, ok := player.(*HumanPlayer); ok {
			s.Send(&protos.SyncDeckNumToc{
				Num:      uint32(len(d.cards)),
				Shuffled: shuffled,
			})
		}
	}
}

func (d *Deck) GetNextId() uint32 {
	d.nextId++
	return d.nextId
}
