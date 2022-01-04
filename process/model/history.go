package process

import (
	"encoding/json"

	"gorm.io/gorm"
)

type UserHistory struct {
	gorm.Model
	Username   string `gorm:"user_name;size:100;not null"`
	IsAccepted bool   `gorm:"is_accepted;default:true"`
	IsLogin    bool   `gorm:"is_login";default:true`
}

func (c *UserHistory) ToString() (string, error) {
	res, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func NewUserHistoryDB(userName string, is_accepted bool, is_login bool) *UserHistory {
	return &UserHistory{Username: userName, IsAccepted: is_accepted, IsLogin: is_login}
}
