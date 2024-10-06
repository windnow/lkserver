package sqlite

import (
	"fmt"
	m "lkserver/internal/models"
)

func InitCoordinators(repo *reportsRepo) error {

	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %[1]s (
			ref			BLOB PRIMARY KEY,
			report		BLOB,
			coordinator	BLOB,
			author		BLOB,
			when_added	INTEGER,

			FOREIGN KEY (report) REFERENCES %[2]s(ref),
			FOREIGN KEY (coordinator) REFERENCES %[3]s(ref),
			FOREIGN KEY (author) REFERENCES %[3]s(ref)
		);
		
		CREATE INDEX IF NOT EXISTS idx_%[1]s_report ON    %[1]s(report);
		`, tabCoordinators, tabReport, tabUsers)

	err := repo.source.Exec(query)
	if err != nil {
		return m.HandleError(err, "sqliteRepo.initEducation")
	}
	return nil

}
