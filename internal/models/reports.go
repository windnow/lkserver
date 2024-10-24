package models

import (
	"database/sql"
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

func NewReport() Scanable {
	return &Report{}
}

func (r *Report) Scan(row *sql.Rows) error {
	return row.Scan(
		&r.Ref,
		&r.Type,
		&r.Date,
		&r.Number,
		&r.RegNumber,
		&r.Author,
	)
}

var ReportMETA = META{
	"ref":        Desc(types.Report, map[string]string{"ru": "Идентификатор"}, 0),
	"type":       Desc(types.ReportType, map[string]string{"ru": "Вид рапорта"}, 1),
	"date":       Desc(types.Date, map[string]string{"ru": "Дата"}, 2),
	"number":     Desc(types.String, map[string]string{"ru": "Номер"}, 3),
	"reg_number": Desc(types.String, map[string]string{"ru": "Регистрационный номер"}, 4),
	"author":     Desc(types.Users, map[string]string{"ru": "Автор"}, 5),
}
