package entity

import "time"

type User struct {
	Id       int
	Username string
	Password string
	CreateAt time.Time
}
