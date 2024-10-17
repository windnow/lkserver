package models

import (
	"lkserver/internal/models/types"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Key          JSONByte `json:"key"`
	Iin          string   `json:"iin"`
	Name         string   `json:"name"`
	Pin          string   `json:"pin,omitempty"`
	Individual   JSONByte `json:"individ"`
	PasswordHash []byte   `json:"-"`
}

var UserMETA META = META{
	"key":     Description{Type: types.Users, Labels: map[string]string{"ru": "Идентификатор"}},
	"iin":     Description{Type: types.String, Labels: map[string]string{"ru": "ИИН"}},
	"name":    Description{Type: types.String, Labels: map[string]string{"ru": "Имя"}},
	"pin":     Description{Type: types.String, Labels: map[string]string{"ru": "ПИН-код"}},
	"individ": Description{Type: types.Individuals, Labels: map[string]string{"ru": "Физическое лицо"}},
}

func (u *User) Sanitize() {
	u.Pin = ""
}

func (u *User) GeneratePasswordHash() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Pin), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = bytes
	return nil
}

func (u *User) CheckPassword() bool {
	err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(u.Pin))
	return err == nil
}
