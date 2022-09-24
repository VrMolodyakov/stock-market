package v1

import (
	"time"

	"github.com/VrMolodyakov/stock-market/internal/domain/entity"
)

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResponse struct {
	Username string `json:"username"`
	Password string `json:"password"`
	CreateAt string `json:"create_at"`
}

type User struct {
	Username string
	CreateAt time.Time
}

func ResponseFromEntity(user entity.User) UserResponse {
	dt := user.CreateAt.Format(time.RFC3339)
	return UserResponse{Username: user.Username, Password: user.Password, CreateAt: dt}
}
