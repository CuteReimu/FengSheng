package game

import (
	"testing"
)

func TestCardIdUnique(t *testing.T) {
	deck := NewDeck(&Game{})
	m := make(map[uint32]int)
	for _, c := range deck.cards {
		m[c.GetId()]++
	}
	for i := uint32(1); i <= uint32(len(m)); i++ {
		if count, ok := m[i]; !ok || count != 1 {
			t.Errorf("%d count: %d", i, count)
		}
	}
}
