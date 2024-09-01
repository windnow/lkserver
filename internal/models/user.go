package models

type User struct {
	Iin string `json:"iin"`
	Pin string `json:"pin,omitempty"`
}

func (u *User) Sanitize() {
	u.Pin = ""
}
