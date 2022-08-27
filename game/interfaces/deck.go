package interfaces

type IDeck interface {
	Shuffle()
	Draw(n int) []ICard
	Discard(cards ...ICard)
	GetDeckCount() int
	GetNextId() uint32
}
