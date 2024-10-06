package reports

import m "lkserver/internal/models"

type BussinesTripDetails struct {
	ReportRef        m.JSONByte
	Supervisor       m.JSONByte
	ActingSupervisor m.JSONByte
	Basis            string
	TransportType    string
}
