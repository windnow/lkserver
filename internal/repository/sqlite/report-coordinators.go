package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	m "lkserver/internal/models"
	"lkserver/internal/models/reports"
	"lkserver/internal/models/types"
	"time"
)

func InitCoordinators(repo *reportsRepo) error {

	/// CHECK FOREIGN KEY coordinator
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %[1]s (
				ref			BLOB PRIMARY KEY,
				report		BLOB NOT NULL,
				coordinator	BLOB NOT NULL,
				author		BLOB NOT NULL,
				when_added	INTEGER NOT NULL,

				FOREIGN KEY (report) REFERENCES %[2]s(ref),
				FOREIGN KEY (coordinator) REFERENCES %[3]s(ref),
				FOREIGN KEY (author) REFERENCES %[3]s(ref)
			);

			CREATE INDEX IF NOT EXISTS idx_%[1]s_report ON    %[1]s(report);
			`, types.Coordinators, types.Report, types.Users)

	err := repo.source.Exec(query)
	if err != nil {
		return m.HandleError(err, "sqliteRepo.initEducation")
	}
	return nil

}

func (repo *reportsRepo) SaveCoordinators(tx *sql.Tx, ctx context.Context, coordinators []*reports.Coordinators) error {

	if len(coordinators) == 0 {
		return nil
	}

	query := fmt.Sprintf("INSERT INTO %[1]s (ref, report, coordinator, author, when_added) VALUES ", types.Coordinators)
	values := []interface{}{}
	placeholders := []string{}

	for _, c := range coordinators {
		placeholders = append(placeholders, "(?, ?, ?, ?, ?)")
		values = append(values, c.Ref, c.ReportRef, c.CoordinatorRef, c.WhoAuthor, time.Time(c.WhenAdded).Unix())
	}
	query += placeholders[0] //fmt.Sprintf("%s", placeholders[0])
	for i := 1; i < len(placeholders); i++ {
		query += fmt.Sprintf(",%s", placeholders[i])
	}

	_, err := tx.ExecContext(ctx, query, values...)
	if err != nil {
		return err
	}

	return nil
}

func (repo *reportsRepo) GetCoordinators(ctx context.Context, report *m.Report) ([]*reports.Coordinators, error) {
	query := fmt.Sprintf("SELECT ref, report, coordinator, author, when_added FROM %[1]s WHERE report = ?", types.Coordinators)
	rows, err := repo.source.db.QueryContext(ctx, query, report.Ref)
	if err != nil {
		return nil, m.HandleError(err, "reportsRepo.GetCoordinators")
	}
	result := []*reports.Coordinators{}
	for rows.Next() {

		data := &reports.Coordinators{}
		if err := rows.Scan(
			&data.Ref,
			&data.ReportRef,
			&data.CoordinatorRef,
			&data.WhoAuthor,
			&data.WhenAdded,
		); err != nil {
			return nil, m.HandleError(err, "reportsRepo.GetCoordinators")
		}

		result = append(result, data)

	}
	return result, nil
}
