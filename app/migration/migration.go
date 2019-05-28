package migration

import (
	"wheel.smart26.com/app/entity"
	"wheel.smart26.com/app/user"
	"wheel.smart26.com/commons/crypto"
	"wheel.smart26.com/commons/db"
)

func Run() {
	db.Conn.AutoMigrate(&entity.User{})

	_, err := user.FindByEmail("wheel@smart26.com")
	if err != nil {
		db.Conn.Create(&entity.User{Name: "Wheel Smart26", Email: "wheel@smart26.com", Password: crypto.SetPassword("secret123"), Locale: "en", Admin: true})
	}

	db.Conn.AutoMigrate(&entity.Session{})
}
