package models

import (
	"service/pkg/repositories"
	"time"
)

var UserName = "users"

type User struct {
	repositories.Query

	Id          int64     `json:"id" db:"id" skipInsert:"+"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	Email       string    `json:"email" db:"email"`
	Password    string    `json:"password" db:"password"`
	FirstName   string    `json:"first_name" db:"first_name"`
	LastName    string    `json:"last_name" db:"last_name"`
	DisplayName string    `json:"display_name" db:"display_name"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	IsAdmin     bool      `json:"is_admin" db:"is_admin"`
	IsSuperuser bool      `json:"is_superuser" db:"is_superuser"`
	CreatedAt   time.Time `json:"created_at" skipUpdate:"+" db:"created_at"`
}

func (u *User) InformMeToQueryProvider() *User {
	u.SetTableName(UserName)
	u.SetRowData(u)
	return u
}

func NewUser() *User {
	user := &User{
		Query: repositories.NewQuery(UserName),
	}
	user.SetRowData(user)
	return user
}
