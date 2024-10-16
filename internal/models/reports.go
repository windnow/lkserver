package models

import "lkserver/internal/models/types"

type Report struct {
	Ref       JSONByte `json:"ref"`
	Type      JSONByte `json:"type"`
	Date      JSONTime `json:"date"`
	Number    string   `json:"number"`
	RegNumber string   `json:"reg_number"`
	Author    JSONByte `json:"author"`
}

var ReportMETA = META{
	"ref":        types.Report,
	"type":       types.ReportType,
	"date":       types.Date,
	"number":     types.String,
	"reg_number": types.String,
	"author":     types.Users,
}
