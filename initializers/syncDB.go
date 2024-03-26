package initializers

import "new_projects/models"

func SyncDB() {
	DB.AutoMigrate(&models.User{})
}
