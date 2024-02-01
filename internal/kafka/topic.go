package kafka

type Topic string

const (
	Orderbook Topic = "Orderbook"
	Trade     Topic = "Trade"
)

// Get string functions

func (t Topic) String() string {
	return string(t)
}
