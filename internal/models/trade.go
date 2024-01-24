package models

import "gorm.io/gorm"

type Trade struct {
	gorm.Model
	Ticker           string
	AggrerateTradeID int
	Price            float64
	Quantity         float64
	BuyerOrderID     int
	SellerOrderID    int
	MarketBuyer      bool
}
