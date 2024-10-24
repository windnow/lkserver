package sqlite

import (
	m "lkserver/internal/models"
	"lkserver/internal/models/reports"
)

type BeingOnBusinesTrip struct {
	source *src
}

func (r *BeingOnBusinesTrip) GetStructure() interface{} {
	details := &reports.BeingOnBusinesTrip{}
	details.Destinations = []reports.BeingOnBusinessTripDestination{}
	return &reports.ReportData{
		Head:         &m.Report{},
		Coordinators: []*reports.Coordinators{},
		Details:      details,
	}
}

func (r *BeingOnBusinesTrip) META() map[string]m.META {
	return map[string]m.META{
		"details":              reports.BeingOnBusinesTripMETA,
		"details.destinations": reports.BeingBusinessTripDestinationMETA,
	}
}

/*
func (r *BeingOnBusinesTrip) Get(ctx context.Context, ref m.JSONByte, txs ...*sql.Tx) (any, map[string]m.META, error) {

}
*/
