package game

import (
	"math/rand"
	"testing"
	"time"
)

func TestCardIdUnique(t *testing.T) {
	deck := NewDeck(rand.New(rand.NewSource(time.Now().UnixMilli())))
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
