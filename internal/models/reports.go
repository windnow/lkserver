package models

import (
	"lkserver/internal/models/types"
)

type Report struct {
	Ref       JSONByte `json:"ref"`
	Type      JSONByte `json:"type"`
	Date      JSONTime `json:"date"`
	Number    string   `json:"number"`
	RegNumber string   `json:"reg_number"`
	Author    JSONByte `json:"author"`
}

var ReportMETA = META{
	"ref":        Description{Type: types.Report, Labels: map[string]string{"ru": "Идентификатор"}},
	"type":       Description{Type: types.ReportType, Labels: map[string]string{"ru": "Вид рапорта"}},
	"date":       Description{Type: types.Date, Labels: map[string]string{"ru": "Дата"}},
	"number":     Description{Type: types.String, Labels: map[string]string{"ru": "Номер"}},
	"reg_number": Description{Type: types.String, Labels: map[string]string{"ru": "Регистрационный номер"}},
	"author":     Description{Type: types.Users, Labels: map[string]string{"ru": "Автор"}},
}
