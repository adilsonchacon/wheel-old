package models

import (
	"wheel.smart26.com/utils"
)

func Migrate() {
	db.AutoMigrate(&User{})

	user := UserFindByEmail("wheel@smart26.com")
	if db.NewRecord(user) {
		db.Create(&User{Name: "Wheel Smart26", Email: "wheel@smart26.com", Password: utils.SaferSetPassword("secret123"), Locale: "en", Admin: true})
	}

	db.AutoMigrate(&Session{})
}
