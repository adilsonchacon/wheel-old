package user

import (
	"time"
	"wheel.smart26.com/app/entity"
	"wheel.smart26.com/commons/view"
)

type PaginationJson struct {
	Pagination view.MainPagination `json:"pagination"`
	Users      []Json              `json:"users"`
}

type SuccessfullySavedJson struct {
	SystemMessage view.SystemMessage `json:"system_message"`
	User          Json               `json:"user"`
}

type Json struct {
	ID        uint       `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Admin     bool       `json:"admin"`
	Locale    string     `json:"locale"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func SetJson(user entity.User) Json {
	return Json{
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
