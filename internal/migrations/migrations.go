package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB, logger *zap.Logger) error {
	options := gormigrate.DefaultOptions
	options.UseTransaction = false

	migrations := gormigrate.New(db, options, []*gormigrate.Migration{

		Migration000NewWorldOrder(logger),
	})

	return migrations.Migrate()
}
