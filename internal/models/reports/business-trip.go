package reports

import (
	m "lkserver/internal/models"
	"lkserver/internal/models/types"
)

type BussinesTripDetails struct {
	ReportRef        m.JSONByte `json:"report_ref"`
	Supervisor       m.JSONByte `json:"supervisor"`
	ActingSupervisor m.JSONByte `json:"acting_supervisor"`
	Basis            string     `json:"basis"`
	TransportType    string     `json:"transport_type"`
}

var BussinesTripDetailsMeta m.META = m.META{
	"report_ref":        m.Description{Type: types.Report, Labels: map[string]string{"ru": "Идентификатор рапорта"}},
	"supervisor":        m.Description{Type: types.Users, Labels: map[string]string{"ru": "Руководитель"}},
	"acting_supervisor": m.Description{Type: types.Users, Labels: map[string]string{"ru": "ИО руководителя"}},
	"basis":             m.Description{Type: types.String, Labels: map[string]string{"ru": "Основание"}},
	"transport_type":    m.Description{Type: types.String, Labels: map[string]string{"ru": "Вид транспорта"}},
}
