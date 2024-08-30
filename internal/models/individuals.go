package models

type Individuals struct {
	FirstName        string   `json:"first_name"`
	LastName         string   `json:"last_name"`
	Patronymic       string   `json:"patronymic"`
	IndividualNumber string   `json:"iin"`
	Image            string   `json:"image"`
	BirthDate        JSONTime `json:"birth_date"`
}
