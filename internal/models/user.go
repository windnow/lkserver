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
	"key":     types.Users,
	"iin":     types.String,
	"name":    types.String,
	"pin":     types.String,
	"individ": types.Individuals,
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
