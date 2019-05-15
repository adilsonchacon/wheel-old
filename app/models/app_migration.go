package models

import (
	"wheel.smart26.com/commons/crypto"
)

func Migrate() {
	db.AutoMigrate(&User{})

	user := UserFindByEmail("wheel@smart26.com")
	if db.NewRecord(user) {
		db.Create(&User{Name: "Wheel Smart26", Email: "wheel@smart26.com", Password: crypto.SetPassword("secret123"), Locale: "en", Admin: true})
	}

	db.AutoMigrate(&Session{})
}
