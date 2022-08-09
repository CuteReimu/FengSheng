package game

import (
	"github.com/CuteReimu/FengSheng/protos"
	"math/rand"
	"time"
)

type Deck struct {
	cards       []*protos.Card
	discardPile []*protos.Card
	random      *rand.Rand
}

func NewDeck() *Deck {
	d := new(Deck)
	d.random = rand.New(rand.NewSource(time.Now().Unix()))
	id := uint32(1)
	for i := uint32(1); i < 108; i++ {
		card := &protos.Card{
			CardId:   id,
			CardDir:  protos.Direction(d.random.Intn(len(protos.Direction_name))),
			CardType: protos.CardType(d.random.Intn(len(protos.CardType_name))),
		}
		color := d.random.Intn(8)
		switch color {
		case 6:
			card.CardColor = []protos.Color{protos.Color_Black, protos.Color_Red}
		case 7:
			card.CardColor = []protos.Color{protos.Color_Black, protos.Color_Blue}
		default:
			card.CardColor = []protos.Color{protos.Color(color / 2)}
		}
		d.cards = append(d.cards, card)
		id++
	}
	d.Shuffle()
	return d
}

func (d *Deck) Shuffle() {
	d.cards = append(d.cards, d.discardPile...)
	d.discardPile = nil
	d.random.Shuffle(len(d.cards), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
}

func (d *Deck) Draw(n int) []*protos.Card {
	if n > len(d.cards) {
		d.Shuffle()
	}
	if n > len(d.cards) {
		n = len(d.cards)
	}
	result := d.cards[:n]
	d.cards = d.cards[n:]
	return result
}

func (d *Deck) Discard(cards ...*protos.Card) {
	d.discardPile = append(d.discardPile, cards...)
}
