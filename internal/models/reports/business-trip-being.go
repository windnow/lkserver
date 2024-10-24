package reports

import (
	"database/sql"
	m "lkserver/internal/models"
	"lkserver/internal/models/types"
)

type BeingOnBusinesTrip struct {
	Acting        m.JSONByte                       `json:"acting"`
	OrderSource   string                           `json:"order_source"`
	Basis         string                           `json:"basis"`
	TransportType string                           `json:"transport_type"`
	TripGoal      string                           `json:"trip_goal"`
	TripFrom      m.JSONTime                       `json:"trip_from"`
	TripTo        m.JSONTime                       `json:"trip_to"`
	TripDuration  int                              `json:"trip_duration"`
	Destinations  []BeingOnBusinessTripDestination `json:"destinations"`
}

var BeingOnBusinesTripMETA = m.META{
	"acting":         m.Desc(types.Users, map[string]string{"ru": "Ф.И.О. временного И.О."}, 1),
	"order_source":   m.Desc(types.String, map[string]string{"ru": "Чьи распоряжения"}, 2),
	"basis":          m.Desc(types.String, map[string]string{"ru": "Основание"}, 3),
	"transport_type": m.Desc(types.String, map[string]string{"ru": "Вид транспорта"}, 4),
	"trip_goal":      m.Desc(types.String, map[string]string{"ru": "Цель поездки"}, 5),
	"trip_from":      m.Desc(types.Date, map[string]string{"ru": "Дата начала командировки"}, 6),
	"trip_to":        m.Desc(types.Date, map[string]string{"ru": "Дата окончания командировки"}, 7),
	"trip_duration":  m.Desc(types.Number, map[string]string{"ru": "Количество суток"}, 12),
}

func (bbtd *BeingOnBusinesTrip) Scan(rows *sql.Rows) error {

	return rows.Scan(
		&bbtd.Acting,
		&bbtd.OrderSource,
		&bbtd.Basis,
		&bbtd.TransportType,
		&bbtd.TripGoal,
		&bbtd.TripFrom,
		&bbtd.TripTo,
		&bbtd.TripDuration,
	)

}

type BeingOnBusinessTripDestination struct {
	Destination  m.JSONByte `json:"destination"`
	Organization string     `json:"organization"`
	From         m.JSONTime `json:"date_from"`
	To           m.JSONTime `json:"date_to"`
	Duration     int        `json:"duration"`
}

func NewBeingBusinesTripDestination() m.Scanable {
	return &BeingOnBusinessTripDestination{}
}

func (bbtd *BeingOnBusinessTripDestination) Scan(rows *sql.Rows) error {
	return rows.Scan(
		&bbtd.Destination,
		&bbtd.Organization,
		&bbtd.From,
		&bbtd.To,
		&bbtd.Duration,
	)
}

var BeingBusinessTripDestinationMETA = m.META{
	"report_ref":   m.Desc(types.Report, map[string]string{"ru": "Идентификатор рапорта"}, 0),
	"destination":  m.Desc(types.Cato, map[string]string{"ru": "Место назначения"}, 1),
	"organization": m.Desc(types.String, map[string]string{"ru": "Организация"}, 2),
	"from":         m.Desc(types.Date, map[string]string{"ru": "Период проживания с"}, 3),
	"to":           m.Desc(types.Date, map[string]string{"ru": "Период проживания по"}, 4),
	"duration":     m.Desc(types.Number, map[string]string{"ru": "Количество суток проживания"}, 5),
}
