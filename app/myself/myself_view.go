package myself

import (
	"time"
	"wheel.smart26.com/app/entity"
)

type Json struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func SetJson(userMyself entity.User) Json {
	return Json{
		ID:        userMyself.ID,
		Name:      userMyself.Name,
		Email:     userMyself.Email,
		CreatedAt: userMyself.CreatedAt,
		UpdatedAt: userMyself.UpdatedAt,
	}
}
