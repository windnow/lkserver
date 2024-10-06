package sqlite

import (
	"context"
	"encoding/json"
	"fmt"
	m "lkserver/internal/models"
	"time"
)

type rankHistoryRepo struct {
	source *src
	repo   *sqliteRepo
}

type draftRH struct {
	D  m.JSONTime `json:"date"`
	Rg m.JSONByte `json:"rank"`
	Ig m.JSONByte `json:"individual"`
}

func (draft *draftRH) update(r *sqliteRepo) (*m.RankHistory, error) {
	ret := func(err error) (*m.RankHistory, error) {
		return nil, m.HandleError(err, "draftRH.update")
	}

	rank, err := r.ranks.Get(draft.Rg)
	if err != nil {
		return ret(err)
	}
	individ, err := r.individuals.Get(draft.Ig)
	if err != nil {
		return ret(err)
	}

	return &m.RankHistory{
		Date:       draft.D,
		Rank:       rank,
		Individual: individ,
	}, nil
}

func (s *sqliteRepo) initRankHistory() error {
	rh := &rankHistoryRepo{
		source: s.db,
		repo:   s,
	}

	err := rh.source.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %[1]s (
			date INTEGER,
			rank BLOB,
			individual BLOB,

			FOREIGN KEY (rank) REFERENCES %[2]s(ref),
			FOREIGN KEY (individual) REFERENCES %[3]s(ref)
		);

		CREATE INDEX IF NOT EXISTS idx_%[1]s_date ON %[1]s(date);
		CREATE INDEX IF NOT EXISTS idx_%[1]s_individ ON %[1]s(individual);
	`, tabRankHistory, tabRanks, tabIndividuals))
	if err != nil {
		return m.HandleError(err, "sqliteRepo.initRankHistory")
	}
	err = rh.source.Exec(`
	`)
	if err != nil {
		return m.HandleError(err, "sqliteRepo.initRankHistory")
	}

	var count int64
	rh.source.db.QueryRow(fmt.Sprintf(`select count(*) from %[1]s`, tabRankHistory)).Scan(&count)
	if count == 0 {
		var result *[]draftRH
		if err := json.Unmarshal([]byte(mockData), &result); err != nil {
			return m.HandleError(err, "sqliteRepo.initRankHistory")
		}
		for _, r := range *result {
			row, err := r.update(s)
			if err != nil {
				return m.HandleError(err, "sqliteRepo.initRankHistory")
			}
			if err = rh.Save(context.Background(), row); err != nil {
				return m.HandleError(err, "sqliteRepo.initRankHistory")
			}
			fmt.Printf("%v\n", r)
		}
	}

	s.rankHistory = rh

	return nil

	// return m.HandleError(errors.ErrUnsupported, "sqliteRepo.initRankHistory")

}

func (r *rankHistoryRepo) Save(ctx context.Context, rh *m.RankHistory) error {
	return m.HandleError(r.source.ExecContextInTransaction(ctx, fmt.Sprintf(`INSERT OR REPLACE INTO %[1]s (date, rank, individual) VALUES (?, ?, ?)`, tabRankHistory),
		time.Time(rh.Date).Unix(), rh.Rank.Key, rh.Individual.Key,
	), "rankHistoryRepo.Save")
}

func (r *rankHistoryRepo) Close() {}

func (r *rankHistoryRepo) GetHistoryByIin(indivIin string) ([]*m.RankHistory, error) {
	individ, err := r.repo.individuals.GetByIin(indivIin)
	if err != nil {
		return nil, m.HandleError(err)
	}
	result, err := r.GetHistory(context.Background(), individ)
	if err != nil {
		return nil, m.HandleError(err)
	}
	if len(result) == 0 {
		return nil, m.HandleError(m.ErrNotFound, "rankHistoryRepo.GetLastByIin")
	}

	return result, nil

}

func (r *rankHistoryRepo) GetLastByIin(individIin string) (*m.RankHistory, error) {

	var last *m.RankHistory
	history, err := r.GetHistoryByIin(individIin)
	if err != nil {
		return nil, m.HandleError(err, "rankHistoryRepo.GetLastByIin")
	}
	if len(history) == 0 {
		return nil, m.HandleError(m.ErrNotFound, "rankHistoryRepo.GetLastByIin")
	}
	for _, record := range history {
		if last == nil || time.Time(record.Date).Unix() > time.Time(last.Date).Unix() {
			last = record
		}
	}

	return last, nil

}

func (r *rankHistoryRepo) GetHistory(ctx context.Context, individ *m.Individuals) ([]*m.RankHistory, error) {

	rows, err := r.source.db.QueryContext(ctx, fmt.Sprintf("SELECT date, rank FROM %[1]s WHERE individual = ?", tabRankHistory), individ.Key)
	if err != nil {
		return nil, m.HandleError(err, "rankHistoryRepo.getHistory")
	}
	defer rows.Close()
	select {
	case <-ctx.Done():
		return nil, m.HandleError(ctx.Err(), "rankHistoryRepo.getHistory")
	default:
	}

	var records []*draftRH

	for rows.Next() {
		record := &draftRH{}
		if err := rows.Scan(&record.D, &record.Rg); err != nil {
			return nil, m.HandleError(err, "rankHistoryRepo.getHistory")
		}
		records = append(records, record)
	}

	var result []*m.RankHistory
	for _, record := range records {

		rank, err := r.repo.ranks.Get(record.Rg)
		if err != nil {
			return nil, m.HandleError(err, "rankHistoryRepo.getHistory")
		}
		result = append(result, &m.RankHistory{
			Date:       record.D,
			Individual: individ,
			Rank:       rank,
		})

	}

	return result, nil
}

var mockData string = `
[
    {
        "date": "2023.10.23",
        "rank": "86bf503e-9327-46d4-8d6c-35dd19b88cfa",
        "individual": "27f74b66-cba7-486d-a263-81b6cb9a3e57"
    },
    {
        "date": "2022.07.02",
        "rank": "758ebb53-eea6-4fde-84fa-1153527a3883",
        "individual": "19db2753-68f9-4b5d-998a-727e347a958a"
    },
    {
        "date": "2023.07.02",
        "rank": "f5e5f01c-6a27-4ae2-a3a3-5d714f9b871f",
        "individual": "52efc72d-ba0d-4f87-ae73-e902936395fe"
    }
]`
