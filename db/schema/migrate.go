package schema

import (
	"wheel.smart26.com/app/user"
	"wheel.smart26.com/commons/app/model"
	"wheel.smart26.com/commons/crypto"
	"wheel.smart26.com/db/entity"
)

func Migrate() {
	model.Db.AutoMigrate(&entity.User{})

	_, err := user.FindByEmail("wheel@smart26.com")
	if err != nil {
		model.Db.Create(&entity.User{Name: "Wheel Smart26", Email: "wheel@smart26.com", Password: crypto.SetPassword("secret123"), Locale: "en", Admin: true})
	}

	model.Db.AutoMigrate(&entity.Session{})
}