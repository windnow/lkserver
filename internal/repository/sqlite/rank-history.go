package sqlite

import (
	"encoding/json"
	"errors"
	"lkserver/internal/models"
)

type rankHistoryRepo struct {
	source *src
}

func (s *sqliteRepo) initRankHistory() error {
	rh := &rankHistoryRepo{
		source: s.db,
	}

	err := rh.source.Exec(`
		CREATE TABLE rank_history (
			date INTEGER,
			rank BLOB,
			individual BLOB
		)
	`)
	if err != nil {
		return err
	}
	err = rh.source.Exec(`
		CREATE INDEX IF NOT EXISTS idx_rank_history_date ON rank_history(date);
		CREATE INDEX IF NOT EXISTS idx_rank_history_individ ON rank_history(individ);
	`)
	if err != nil {
		return err
	}

	var count int64
	rh.source.db.QueryRow(`select count(*) from rank_history`).Scan(&count)
	if count == 0 {
		var result *[]struct {
			D  models.JSONTime
			Rg models.JSONByte
			Ig models.JSONByte
		} = nil
		if err := json.Unmarshal([]byte(mockData), &result); err != nil {
			return err
		}
	}

	s.rankHistory = rh

	return errors.ErrUnsupported

}

var mockData string = `
[
    {
        "date": "2023-10-23",
        "rank": "86bf503e-9327-46d4-8d6c-35dd19b88cfa",
        "individual": "c9aba8d6-351a-4d85-a8b6-9427ea2f8c8e"
    },
    {
        "date": "2022-07-02",
        "rank": "758ebb53-eea6-4fde-84fa-1153527a3883",
        "individual": "f31c6a0f-b07c-4632-8949-2f24fde4fc26"
    },
    {
        "date": "2023-07-02",
        "rank": "f5e5f01c-6a27-4ae2-a3a3-5d714f9b871f",
        "individual": "8c272f7c-6c2c-4dba-bba5-4062005b2400"
    }
]`
