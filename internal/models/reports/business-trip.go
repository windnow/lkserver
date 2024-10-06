package reports

import m "lkserver/internal/models"

type BussinesTripDetails struct {
	ReportRef        m.JSONByte `json:"report_ref"`
	Supervisor       m.JSONByte `json:"supervisor"`
	ActingSupervisor m.JSONByte `json:"acting_supervisor"`
	Basis            string     `json:"basis"`
	TransportType    string     `json:"transport_type"`
}
