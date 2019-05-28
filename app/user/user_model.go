package user

import (
	"errors"
	"regexp"
	"strings"
	"time"
	"wheel.smart26.com/app/entity"
	"wheel.smart26.com/commons/app/model"
	"wheel.smart26.com/commons/app/model/pagination"
	"wheel.smart26.com/commons/crypto"
	"wheel.smart26.com/config"
)

const NotFound = "user was not found"

var Current entity.User

func Find(id interface{}) (entity.User, error) {
	var user entity.User
	var err error

	model.Db.First(&user, id, "deleted_at IS NULL")
	if model.Db.NewRecord(user) {
		err = errors.New(NotFound)
	}

	return user, err
}

func FindAll() []entity.User {
	var users []entity.User

	model.Db.Order("name").Find(&users, "deleted_at IS NULL")

	return users
}

func IsValid(user *entity.User) (bool, []error) {
	var count int
	var errs []error
	var validEmail = regexp.MustCompile(`\A[^@]+@([^@\.]+\.)+[^@\.]+\z`)

	if len(user.Name) == 0 {
		errs = append(errs, errors.New("name can't be blank"))
	} else if len(user.Name) > 255 {
		errs = append(errs, errors.New("name is too long"))
	}

	if len(user.Email) == 0 {
		errs = append(errs, errors.New("email can't be blank"))
	} else if len(user.Email) > 255 {
		errs = append(errs, errors.New("email is too long"))
	} else if !validEmail.MatchString(user.Email) {
		errs = append(errs, errors.New("email is invalid"))
	} else if model.Db.Model(&entity.User{}).Where("id <> ? AND email = ? AND deleted_at IS NULL", user.ID, user.Email).Count(&count); count > 0 {
		errs = append(errs, errors.New("email has already been taken"))
	}

	if len(user.Password) < 8 {
		errs = append(errs, errors.New("password is too short minimum is 8 characters"))
	} else if len(user.Password) > 255 {
		errs = append(errs, errors.New("password is too long"))
	}

	if !isLocaleValid(user.Locale) {
		errs = append(errs, errors.New("locale is invalid"))
	}

	return (len(errs) == 0), errs
}

func Update(user *entity.User) (bool, []error) {
	var newValue, currentValue interface{}
	var valid bool
	var errs []error

	mapUpdate := make(map[string]interface{})

	currentUser, findErr := Find(user.ID)
	if findErr != nil {
		return false, []error{findErr}
	}

	valid, errs = IsValid(user)

	if valid {
		columns := model.ColumnsFromTable(user, false)
		for _, column := range columns {
			newValue, _ = model.GetColumnValue(user, column)
			currentValue, _ = model.GetColumnValue(currentUser, column)

			if newValue != currentValue {
				mapUpdate[column] = newValue

				if column == "password" {
					mapUpdate[column] = crypto.SetPassword(mapUpdate[column].(string))
				}
			}
		}

		if len(mapUpdate) > 0 {
			model.Db.Model(&user).Updates(mapUpdate)
		}

	}

	return valid, errs
}

func Create(user *entity.User) (bool, []error) {
	valid, errs := IsValid(user)
	if valid && model.Db.NewRecord(user) {
		user.Password = crypto.SetPassword(user.Password)

		model.Db.Create(&user)

		if model.Db.NewRecord(user) {
			errs = append(errs, errors.New("database error"))
			return false, errs
		}
	}

	return valid, errs
}

func Save(user *entity.User) (bool, []error) {
	if model.Db.NewRecord(user) {
		return Create(user)
	} else {
		return Update(user)
	}
}

func Destroy(user *entity.User) bool {
	if model.Db.NewRecord(user) {
		return false
	} else {
		model.Db.Delete(&user)
		return true
	}
}

func FindByEmail(email string) (entity.User, error) {
	var user entity.User
	var err error

	model.Db.Where("email = ? AND deleted_at IS NULL", email).First(&user)
	if model.Db.NewRecord(user) {
		user = entity.User{}
		err = errors.New(NotFound)
	}

	return user, err
}

func FindByResetPasswordToken(token string) (entity.User, error) {
	var user entity.User
	var err error

	enconded_token := crypto.EncryptText(token, config.SecretKey())
	two_days_ago := time.Now().Add(time.Second * time.Duration(config.ResetPasswordExpirationSeconds()) * (-1))

	model.Db.Where("reset_password_token = ? AND reset_password_sent_at >= ? AND deleted_at IS NULL", enconded_token, two_days_ago).First(&user)
	if model.Db.NewRecord(user) {
		user = entity.User{}
		err = errors.New(NotFound)
	}

	return user, err
}

func Paginate(criteria map[string]string, page interface{}, perPage interface{}) ([]entity.User, int, int, int) {
	var users []entity.User
	var user entity.User

	search, currentPage, totalPages, totalEntries := pagination.Query(&user, criteria, page, perPage)

	search.Order("name").Find(&users, "deleted_at IS NULL")

	return users, currentPage, totalPages, totalEntries
}

func Authenticate(email string, password string) (entity.User, error) {
	user, err := FindByEmail(email)

	if model.Db.NewRecord(user) || !crypto.CheckPassword(password, user.Password) {
		user = entity.User{}
		err = errors.New("invalid credentials")
	}

	return user, err
}

func IsNil(user *entity.User) bool {
	return model.Db.NewRecord(user)
}

func Exists(user *entity.User) bool {
	return !IsNil(user)
}

func SetCurrent(id interface{}) error {
	var err error
	Current, err = Find(id)

	return err
}

func IdExists(id interface{}) bool {
	_, err := Find(id)

	return (err == nil)
}

func SetRecovery(user *entity.User) (string, []error) {
	token := crypto.RandString(20)

	if model.Db.NewRecord(user) {
		return "", []error{errors.New(NotFound)}
	} else {
		t := time.Now()
		user.ResetPasswordSentAt = &t
		user.ResetPasswordToken = crypto.EncryptText(token, config.SecretKey())

		valid, errs := Save(user)

		if valid {
			return token, errs
		} else {
			return "", errs
		}
	}
}

func ClearRecovery(user *entity.User) (bool, []error) {
	if model.Db.NewRecord(user) {
		return false, []error{errors.New(NotFound)}
	} else {
		user.ResetPasswordToken = ""
		user.ResetPasswordSentAt = nil
		valid, errs := Save(user)

		return valid, errs
	}
}

func FirstName(user *entity.User) string {
	return strings.Split(user.Name, " ")[0]
}

// local methods

func isLocaleValid(locale string) bool {
	locales := config.Locales()

	for _, a := range locales {
		if a == locale {
			return true
		}
	}

	return false
}
