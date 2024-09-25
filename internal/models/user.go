package models

type User struct {
	Key JSONByte `json:"key"`
	Iin string   `json:"iin"`
	Pin string   `json:"pin,omitempty"`
}

func (u *User) Sanitize() {
	u.Pin = ""
}
