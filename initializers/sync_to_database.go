package initializers

import "github.com/Waris-Shaik/todo-backend/models"

func SyncDatabase() {

	DB.AutoMigrate(&models.User{}, &models.Todo{})

}
