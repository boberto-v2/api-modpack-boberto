package entities_user

import "time"

type User struct {
	ID       string
	Email    string
	Password string
	Username string
	CreateAt time.Time
	UpdateAt time.Time
}
