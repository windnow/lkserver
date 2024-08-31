package models

type Individuals struct {
	Code             string   `json:"code"`
	FirstName        string   `json:"first_name"`
	LastName         string   `json:"last_name"`
	Patronymic       string   `json:"patronymic"`
	IndividualNumber string   `json:"iin"`
	Image            string   `json:"image"`
	BirthDate        JSONTime `json:"birth_date"`
	BirthPlace       string   `json:"birth_place"`
	PersonalNumber   string   `json:"personal_number"`
}
