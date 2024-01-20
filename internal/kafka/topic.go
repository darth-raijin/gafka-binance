package kafka

type Topic string

const (
	Orderbook Topic = "Orderbook"
	Trade     Topic = "Trade"
)
