package models

import (
	"github.com/jinzhu/gorm"
	"regexp"
	"strings"
	"time"
	"wheel.smart26.com/commons/crypto"
	"wheel.smart26.com/commons/log"
	"wheel.smart26.com/config"
)

type User struct {
	gorm.Model
	Name                string     `gorm:"type:varchar(255)"`
	Email               string     `gorm:"type:varchar(255);unique_index"`
	Admin               bool       `gorm:"default:false"`
	Password            string     `gorm:"type:varchar(255)"`
	ResetPasswordToken  string     `gorm:"type:varchar(255)"`
	ResetPasswordSentAt *time.Time `gorm:"default:null"`
	Locale              string     `gorm:"type:varchar(255);default:'en'"`
	Sessions            []Session
}

var UserCurrent User

func userValidate(user *User) bool {
	var count int
	var users []User
	var validEmail = regexp.MustCompile(`\A[^@]+@([^@\.]+\.)+[^@\.]+\z`)

	Errors = nil

	if len(user.Name) == 0 {
		Errors = append(Errors, "name can't be blank")
	} else if len(user.Name) > 255 {
		Errors = append(Errors, "name is too long")
	}

	if len(user.Email) == 0 {
		Errors = append(Errors, "email can't be blank")
	} else if len(user.Email) > 255 {
		Errors = append(Errors, "email is too long")
	} else if !validEmail.MatchString(user.Email) {
		Errors = append(Errors, "email is invalid")
	} else if db.Where("id <> ? AND email = ? AND deleted_at IS NULL", user.ID, user.Email).Find(&users).Count(&count); count > 0 {
		Errors = append(Errors, "email has already been taken")
	}

	if len(user.Password) < 8 {
		Errors = append(Errors, "password is too short minimum is 8 characters")
	} else if len(user.Password) > 255 {
		Errors = append(Errors, "password is too long")
	}

	if !userLocaleValid(user.Locale) {
		Errors = append(Errors, "locale is invalid")
	}

	return (len(Errors) == 0)
}

func userLocaleValid(locale string) bool {
	s := config.Locales()

	for _, a := range s {
		if a == locale {
			return true
		}
	}

	return false
}

func UserUpdate(user *User) bool {
	var valueNew interface{}
	var valueCurrent interface{}

	mapUpdate := make(map[string]interface{})
	currentUser := UserFind(user.ID)

	if userValidate(user) && !db.NewRecord(user) {
		columns := ColumnsFromTable(user, false)
		for _, column := range columns {
			valueNew, _ = GetColumnValue(user, column)
			valueCurrent, _ = GetColumnValue(currentUser, column)

			if valueNew != valueCurrent {
				mapUpdate[column] = valueNew

				if column == "password" {
					log.Info.Println(mapUpdate[column])
					log.Info.Println(mapUpdate[column].(string))
					mapUpdate[column] = crypto.SetPassword(mapUpdate[column].(string))
				}
			}
		}

		if len(mapUpdate) > 0 {
			db.Model(&user).Updates(mapUpdate)
		}

		return true
	} else if db.NewRecord(user) {
		Errors = nil
		Errors = append(Errors, "user was not found")
		return false
	} else {
		return false
	}
}

func UserCreate(user *User) bool {
	log.Info.Println("models: SaveUser")

	if userValidate(user) && db.NewRecord(user) {
		user.Password = crypto.SetPassword(user.Password)

		db.Create(&user)

		if db.NewRecord(user) {
			Errors = append(Errors, "database error")
			return false
		}
		return true
	} else {
		return false
	}
}

func UserSave(user *User) bool {
	if db.NewRecord(user) {
		return UserCreate(user)
	} else {
		return UserUpdate(user)
	}
}

func UserFind(id interface{}) User {
	var user User

	db.First(&user, id, "deleted_at IS NULL")

	return user
}

func UserFindByEmail(email string) User {
	var user User

	db.Where("email = ? AND deleted_at IS NULL", email).First(&user)
	if db.NewRecord(user) {
		user = User{}
	}

	return user
}

func UserFindByResetPasswordToken(token string) User {
	var user User

	enconded_token := crypto.EncryptText(token, config.SecretKey())
	two_days_ago := time.Now().Add(time.Second * time.Duration(config.ResetPasswordExpirationSeconds()) * (-1))

	db.Where("reset_password_token = ? AND reset_password_sent_at >= ? AND deleted_at IS NULL", enconded_token, two_days_ago).First(&user)
	if db.NewRecord(user) {
		user = User{}
	}

	return user
}

func UserDestroy(user *User) bool {
	if db.NewRecord(user) {
		return false
	} else {
		db.Delete(&user)
		return true
	}
}

func UserList() []User {
	var users []User

	db.Order("name").Find(&users, "deleted_at IS NULL")

	return users
}

func UserPaginate(criteria map[string]string, page interface{}, perPage interface{}) ([]User, int, int, int) {
	var users []User
	var user User

	search, currentPage, totalPages, totalEntries := GenericPaginateQuery(&user, criteria, page, perPage)

	search.Order("name").Find(&users, "deleted_at IS NULL")

	return users, currentPage, totalPages, totalEntries
}

func UserAuthenticate(email string, password string) User {
	var user User
	user = UserFindByEmail(email)
	if !db.NewRecord(user) && !crypto.CheckPassword(password, user.Password) {
		user = User{}
	}

	return user
}

func UserNil(user User) bool {
	return db.NewRecord(user)
}

func UserExists(user User) bool {
	return !UserNil(user)
}

func UserSetCurrent(id uint) User {
	UserCurrent = UserFind(id)

	return UserCurrent
}

func UserIdExists(id uint) bool {
	return UserExists(UserCurrent)
}

func UserSetRecovery(user *User) string {
	t := time.Now()
	token := crypto.RandString(20)

	if db.NewRecord(user) {
		return ""
	} else {
		user.ResetPasswordToken = crypto.EncryptText(token, config.SecretKey())
		user.ResetPasswordSentAt = &t
		UserSave(user)

		return token
	}
}

func UserClearRecovery(user *User) bool {
	if db.NewRecord(user) {
		return false
	} else {
		user.ResetPasswordToken = ""
		user.ResetPasswordSentAt = nil
		UserSave(user)

		return true
	}
}

func UserFirstName(user *User) string {
	return strings.Split(user.Name, " ")[0]
}
