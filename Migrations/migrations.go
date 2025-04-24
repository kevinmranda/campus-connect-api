package migrations

import (
	initializers "github.com/group4/campus-connect-api/Initializers"
	models "github.com/group4/campus-connect-api/Models"
)

func SyncDatabase() {
	initializers.DB.AutoMigrate(
		&models.User{},
		&models.Job{},
		&models.Timetable{},
		&models.Event{},
		&models.Post{},
	)
}