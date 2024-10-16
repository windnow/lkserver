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
	"report_ref":        types.Report,
	"supervisor":        types.Users,
	"acting_supervisor": types.Users,
	"basis":             types.String,
	"transport_type":    types.String,
}
