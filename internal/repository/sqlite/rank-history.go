package sqlite

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	. "lkserver/internal/models"
	"time"
)

type rankHistoryRepo struct {
	source *src
}

type draftRH struct {
	D  JSONTime `json:"date"`
	Rg JSONByte `json:"rank"`
	Ig JSONByte `json:"individual"`
}

func (draft *draftRH) update(r *sqliteRepo) (*RankHistory, error) {
	ret := func(err error) (*RankHistory, error) {
		return nil, HandleError(err, "draftRH.update")
	}

	rank, err := r.ranks.Get(draft.Rg)
	if err != nil {
		return ret(err)
	}
	individ, err := r.individuals.Get(draft.Ig)
	if err != nil {
		return ret(err)
	}

	return &RankHistory{
		Date:       draft.D,
		Rank:       rank,
		Individual: individ,
	}, nil
}

func (s *sqliteRepo) initRankHistory() error {
	rh := &rankHistoryRepo{
		source: s.db,
	}

	err := rh.source.Exec(`
		CREATE TABLE IF NOT EXISTS rank_history (
			date INTEGER,
			rank BLOB,
			individual BLOB
		)
	`)
	if err != nil {
		return HandleError(err, "sqliteRepo.initRankHistory")
	}
	err = rh.source.Exec(`
		CREATE INDEX IF NOT EXISTS idx_rank_history_date ON rank_history(date);
		CREATE INDEX IF NOT EXISTS idx_rank_history_individ ON rank_history(individual);
	`)
	if err != nil {
		return HandleError(err, "sqliteRepo.initRankHistory")
	}

	var count int64
	rh.source.db.QueryRow(`select count(*) from rank_history`).Scan(&count)
	if count == 0 {
		var result *[]draftRH
		if err := json.Unmarshal([]byte(mockData), &result); err != nil {
			return HandleError(err, "sqliteRepo.initRankHistory")
		}
		for _, r := range *result {
			row, err := r.update(s)
			if err != nil {
				return HandleError(err)
			}
			if err = rh.Save(context.Background(), row); err != nil {
				return HandleError(err)
			}
			fmt.Printf("%v\n", r)
		}
	}

	s.rankHistory = rh

	return HandleError(errors.ErrUnsupported, "sqliteRepo.initRankHistory")

}

func (r *rankHistoryRepo) Save(ctx context.Context, rh *RankHistory) error {
	return HandleError(r.source.ExecContextInTransaction(ctx, `INSERT OR REPLACE INTO rank_history (date, rank, individual) VALUES (?, ?, ?)`,
		time.Time(rh.Date), rh.Rank.Key, rh.Individual.Key,
	), "rankHistoryRepo.Save")
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
