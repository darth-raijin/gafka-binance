package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Migration000NewWorldOrder(log *zap.Logger) *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "000_new_world_order",
		Migrate: func(tx *gorm.DB) error {
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

			err := tx.AutoMigrate(
				&Trade{},
			)
			if err != nil {
				return err
			}

			return nil
		},

		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	}
}
