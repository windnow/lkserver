package models

import "time"

type User struct {
	Name      string
	Iin       string
	Pin       string
	BirthDate time.Time
}
