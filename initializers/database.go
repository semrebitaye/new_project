package initializers

import (
	"os"

	"golang.org/x/exp/slog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error
	dsn := os.Getenv("DB_URL")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		FullSaveAssociations: true,
	})

	if err != nil {
		slog.Error("Failed to connect to database", err)
		return
	}
}
