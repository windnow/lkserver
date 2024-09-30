package sqlite

import (
	"context"
	"encoding/json"
	"fmt"
	m "lkserver/internal/models"
)

type rankRepo struct {
	source *src
}

func (s *sqliteRepo) initRankRepo() error {
	r := &rankRepo{
		source: s.db,
	}
	if err := r.source.Exec(`
		CREATE TABLE IF NOT EXISTS ranks(
			guid BLOB PRIMARY KEY,
			name TEXT
		)
	`); err != nil {
		return m.HandleError(err, "initRankRepo")
	}

	var count int64
	r.source.db.QueryRow(`select count(*) from ranks`).Scan(&count)
	var ranks []m.Rank
	if err := json.Unmarshal([]byte(rankMock), &ranks); err != nil {
		return m.HandleError(err, "initRankRepo")
	}
	if count == 0 {

		for _, rank := range ranks {

			if err := r.Save(context.Background(), &rank); err != nil {
				return err
			}

		}

	}
	r.Get(ranks[1].Key)

	s.ranks = r
	return nil
}

func (r *rankRepo) Get(key m.JSONByte) (*m.Rank, error) {

	rank := &m.Rank{Key: key}
	err := r.source.db.QueryRow(`
		SELECT name from ranks WHERE guid = ?	
	`, key).Scan(&rank.Name)
	if err != nil {
		value := fmt.Sprint(key)
		return nil, m.HandleError(err, "rankRepo.Get ", value)
	}

	return rank, nil

}

func (r *rankRepo) Save(ctx context.Context, rank *m.Rank) error {

	return m.HandleError(r.source.ExecContextInTransaction(ctx, `INSERT OR REPLACE INTO ranks(guid, name) VALUES (?, ?)`,
		rank.Key,
		rank.Name,
	), "rankRepo.Save")
}

func (r *rankRepo) Close() {}

var rankMock string = `[
    {
        "key": "86bf503e-9327-46d4-8d6c-35dd19b88cfa",
        "name": "Полковник"
    },
    {
        "key": "cd2ae871-a915-4531-8f8e-afa309df67de",
        "name": "Майор"
    },
    {
        "key": "758ebb53-eea6-4fde-84fa-1153527a3883",
        "name": "Лейтенант"
    },
    {
        "key": "f5e5f01c-6a27-4ae2-a3a3-5d714f9b871f",
        "name": "Старший лейтенант"
    }
]`
