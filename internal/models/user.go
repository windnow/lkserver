package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	Key          JSONByte     `json:"key"`
	Iin          string       `json:"iin"`
	Pin          string       `json:"pin,omitempty"`
	Individual   *Individuals `json:"individ"`
	PasswordHash []byte       `json:"-"`
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
