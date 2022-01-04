package process

import (
	"encoding/json"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"user_name;size:100;not null"`
	Password string `gorm:"password;not null"`
	Email    string `gorm:"size:100"`
}

func (c *User) ToString() (string, error) {
	res, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func NewUserDB(userName string, password string, email string) *User {
	return &User{Username: userName, Email: email, Password: password}
}
