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
	"key":     Desc(types.Users, map[string]string{"ru": "Идентификатор"}, 0),
	"iin":     Desc(types.String, map[string]string{"ru": "ИИН"}, 2),
	"name":    Desc(types.String, map[string]string{"ru": "Имя"}, 1),
	"pin":     Desc(types.String, map[string]string{"ru": "ПИН-код"}, 3),
	"individ": Desc(types.Individuals, map[string]string{"ru": "Физическое лицо"}, 4),
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
