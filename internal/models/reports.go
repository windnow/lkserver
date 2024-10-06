package models

type Report struct {
	Ref       JSONByte `json:"ref"`
	Date      JSONTime `json:"date"`
	Number    string   `json:"number"`
	RegNumber string   `json:"reg_number"`
	Author    JSONByte `json:"author"`
}
