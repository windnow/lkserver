package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	m "lkserver/internal/models"
	"lkserver/internal/models/reports"
	"lkserver/internal/models/types"
)

var (
	ErrWrongType = fmt.Errorf("INVALID REPORT DATA TYPE")
)

type DepartureOnBusinessTrip struct {
	source *src
}

func (r *DepartureOnBusinessTrip) GetStructure() interface{} {
	details := &reports.BussinesTripDetails{}
	details.Destinations = []reports.BusinessTripDestination{}
	return &reports.ReportData{Details: details}
}

func (r *DepartureOnBusinessTrip) Get(ctx context.Context, ref m.JSONByte, txs ...*sql.Tx) (any, map[string]m.META, error) {
	query := fmt.Sprintf(`
	SELECT acting, unscheduled, devision, article_number, order_source, order_number, order_date, other, transport_type, trip_goal, trip_from, trip_to, trip_duration
	FROM %[1]s WHERE report = ?
	`, types.BusinessTrip)

	var tx *sql.Tx
	if len(txs) > 0 {
		tx = txs[0]
	}
	result, err := m.Query[*reports.BussinesTripDetails](r.source.db, ctx, reports.NewBusinesTripDetails, tx, query, ref)
	if err != nil {
		return nil, nil, m.HandleError(err, "DepartureOnBusinessTrip.Get")
	}
	if len(result) != 1 {
		return nil, nil, m.HandleError(m.ErrNotFound, "DepartureOnBusinessTrip.Get")
	}
	details := result[0]
	{
		query := fmt.Sprintf(`
        SELECT
            destination, organization, date_from, date_to, duration
        FROM %[1]s WHERE report_ref = ?`, types.BusinessTripDest)
		result, err := m.Query[*reports.BusinessTripDestination](r.source.db, ctx, reports.NewBusinesTripDestination, tx, query, ref)
		if err != nil {
			return nil, nil, m.HandleError(m.ErrNotFound, "DepartureOnBusinessTrip.Get")
		}
		for _, dest := range result {
			details.Destinations = append(details.Destinations, *dest)
		}
	}

	return details, map[string]m.META{"details": reports.BussinesTripDetailsMeta}, nil
}
func (r *DepartureOnBusinessTrip) Save(tx *sql.Tx, ctx context.Context, report m.JSONByte, data any) error {
	details, ok := data.(*reports.BussinesTripDetails)
	if !ok {
		m.HandleError(ErrWrongType, "DepartureOnBusinessTrip.Save")
	}
	query := fmt.Sprintf(`
	INSERT OR REPLACE INTO %[1]s (
        report_ref, acting, unscheduled, devision, article_number, order_source, order_number, order_date, other, transport_type, trip_goal, trip_from, trip_to, trip_duration
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, types.BusinessTrip)

	_, err := tx.ExecContext(ctx, query,
		report,
		details.Acting,
		details.Unscheduled,
		details.Devision,
		details.ArticleNumber,
		details.OrderSource,
		details.OrderNumber,
		details.OrderDate.Unix(),
		details.Other,
		details.TransportType,
		details.TripGoal,
		details.TripFrom.Unix(),
		details.TripTo.Unix(),
		details.TripDuration,
	)
	if err != nil {
		return m.HandleError(err, "DepartureOnBusinessTrip.Save")
	}

	query = fmt.Sprintf(`
    INSERT INTO %[1]s (
        report_ref, destination, organization, date_from, date_to, duration
	) VALUES (?, ?, ?, ?, ?, ?)
	`, types.BusinessTripDest)
	sttmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return m.HandleError(err, "DepartureOnBusinessTrip.Save")
	}
	for _, destination := range details.Destinations {

		sttmt.ExecContext(ctx,
			report,
			destination.Destination,
			destination.Organization,
			destination.From.Unix(),
			destination.To.Unix(),
			destination.Destination,
		)

	}

	return nil
}
func (r *DepartureOnBusinessTrip) Init() error {

	query := fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %[1]s (
            report_ref      BLOB,
            acting          BLOB,
            unscheduled     INTEGER,
            devision        BLOB,
            article_number  INTEGER,
            order_source    BLOB,
            order_number    INTEGER,
            order_date      INTEGER,
            other           TEXT,
            transport_type  TEXT,
            trip_goal       TEXT,
            trip_from       INTEGER,
            trip_to         INTEGER,
            trip_duration   INTEGER,

            foreign key (report_ref)        references %[2]s(ref),
            foreign key (acting)            references %[3]s(ref),
            foreign key (devision)          references %[4]s(ref),
            foreign key (order_source)      references %[5]s(ref)
        );
	
        CREATE INDEX IF NOT EXISTS idx_%[1]s_report_ref ON    %[1]s(report_ref);

        CREATE TABLE IF NOT EXISTS %[6]s(
            report_ref   BLOB,
            destination  BLOB,
            organization TEXT,
            date_from    INTEGER,
            date_to      INTEGER,
            duration     INTEGER,

            FOREIGN KEY (report_ref)  REFERENCES %[2]s(ref),
            FOREIGN KEY (destination) REFERENCES %[7]s(ref)
        );
        CREATE INDEX IF NOT EXISTS idx_%[6]s_report_ref ON    %[6]s(report_ref);
		`,
		types.BusinessTrip,     // 1
		types.Report,           // 2
		types.Users,            // 3
		types.Devision,         // 4
		types.OrderSource,      // 5
		types.BusinessTripDest, // 6
		types.Cato,             // 7
	)

	err := r.source.Exec(query)
	if err != nil {
		return m.HandleError(err, "DepartureOnBusinessTrip.Init")
	}
	return nil
}
