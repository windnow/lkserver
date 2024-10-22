package reports

import (
	"database/sql"
	m "lkserver/internal/models"
	"lkserver/internal/models/types"
)

type BussinesTripDetails struct {
	Acting        m.JSONByte                `json:"acting"`
	Unscheduled   bool                      `json:"unscheduled"`
	Devision      m.JSONByte                `json:"devision"`
	ArticleNumber int                       `json:"article_number"`
	OrderSource   m.JSONByte                `json:"order_source"`
	OrderNumber   int                       `json:"order_number"`
	OrderDate     m.JSONTime                `json:"order_date"`
	Other         string                    `json:"other"`
	TransportType string                    `json:"transport_type"`
	TripGoal      string                    `json:"trip_goal"`
	TripFrom      m.JSONTime                `json:"trip_from"`
	TripTo        m.JSONTime                `json:"trip_to"`
	TripDuration  int                       `json:"trip_duration"`
	Destinations  []BusinessTripDestination `json:"destinations"`
}

func NewBusinesTripDetails() m.Scanable {
	return &BussinesTripDetails{}
}

func (btd *BussinesTripDetails) Scan(rows *sql.Rows) error {

	return rows.Scan(
		&btd.Acting,
		&btd.Unscheduled,
		&btd.Devision,
		&btd.ArticleNumber,
		&btd.OrderSource,
		&btd.OrderNumber,
		&btd.OrderDate,
		&btd.Other,
		&btd.TransportType,
		&btd.TripGoal,
		&btd.TripFrom,
		&btd.TripTo,
		&btd.TripDuration,
	)

}

var BussinesTripDetailsMeta m.META = m.META{
	"report_ref":     m.Desc(types.Report, map[string]string{"ru": "Идентификатор рапорта"}, 0),
	"acting":         m.Desc(types.Users, map[string]string{"ru": "Временно исполняющий обязанности"}, 2),
	"unscheduled":    m.Desc(types.Bool, map[string]string{"ru": "Внеплановое"}, 1),
	"devision":       m.Desc(types.Devision, map[string]string{"ru": "Подразделение"}, 3),
	"article_number": m.Desc(types.Number, map[string]string{"ru": "Номер пункта"}, 4),
	"order_source":   m.Desc(types.OrderSource, map[string]string{"ru": "Чей приказ"}, 5),
	"order_number":   m.Desc(types.String, map[string]string{"ru": "Номер приказа"}, 6),
	"order_date":     m.Desc(types.Date, map[string]string{"ru": "Дата приказа"}, 7),
	"transport_type": m.Desc(types.String, map[string]string{"ru": "Вид транспорта"}, 8),
	"trip_goal":      m.Desc(types.String, map[string]string{"ru": "Цель поездки"}, 9),
	"trip_from":      m.Desc(types.Date, map[string]string{"ru": "Дата начала командировки"}, 10),
	"trip_to":        m.Desc(types.Date, map[string]string{"ru": "Дата окончания командировки"}, 11),
	"trip_duration":  m.Desc(types.Number, map[string]string{"ru": "Количество суток"}, 12),
	"destinations":   m.Desc(types.BusinessTripDest, map[string]string{"ru": "Затраты на проживание"}, 13),
}

type BusinessTripDestination struct {
	Destination  m.JSONByte `json:"destination"`
	Organization string     `json:"organization"`
	From         m.JSONTime `json:"date_from"`
	To           m.JSONTime `json:"date_to"`
	Duration     int        `json:"duration"`
}

func NewBusinesTripDestination() m.Scanable {
	return &BusinessTripDestination{}
}

func (btd *BusinessTripDestination) Scan(rows *sql.Rows) error {
	return rows.Scan(
		&btd.Destination,
		&btd.Organization,
		&btd.From,
		&btd.To,
		&btd.Duration,
	)
}

var BusinessTripDestinationMETA = m.META{
	"report_ref":   m.Desc(types.Report, map[string]string{"ru": "Идентификатор рапорта"}, 0),
	"destination":  m.Desc(types.Cato, map[string]string{"ru": "Место назначения"}, 1),
	"organization": m.Desc(types.String, map[string]string{"ru": "Организация"}, 2),
	"from":         m.Desc(types.Date, map[string]string{"ru": "Период проживания с"}, 3),
	"to":           m.Desc(types.Date, map[string]string{"ru": "Период проживания по"}, 4),
	"duration":     m.Desc(types.Number, map[string]string{"ru": "Количество суток проживания"}, 5),
}
