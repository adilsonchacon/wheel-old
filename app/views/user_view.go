package views

import (
	"time"
	"wheel.smart26.com/app/models"
)

type UserPaginationJson struct {
	Pagination MainPagination `json:"pagination"`
	Users      []UserJson     `json:"users"`
}

type UserSuccessfullySavedJson struct {
	SystemMessage SystemMessage `json:"system_message"`
	User          UserJson      `json:"user"`
}

type UserJson struct {
	ID        uint       `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Admin     bool       `json:"admin"`
	Locale    string     `json:"locale"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func SetUserJson(user models.User) UserJson {
	return UserJson{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Admin:     user.Admin,
		Locale:    user.Locale,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}
}
