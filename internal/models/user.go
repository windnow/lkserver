package models

type User struct {
	Name      string   `json:"name"`
	Iin       string   `json:"iin"`
	Pin       string   `json:"pin,omitempty"`
	BirthDate JSONTime `json:"birth_date,omitempty"`
	Image     string   `json:"image"`
}

func (u *User) Sanitize() {
	u.Pin = ""
}
