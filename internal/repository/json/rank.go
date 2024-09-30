package json

import (
	"context"
	"errors"
	"lkserver/internal/models"
)

type rankRepo struct {
	ranks []*models.Rank
}

func (r *repo) initRankRepo() error {
	repo := &rankRepo{}
	if err := initFile(r.dataDir+"/ranks.json", &repo.ranks); err != nil {
		return err
	}

	r.ranks = repo
	return nil
}

func (r *rankRepo) Get(key models.JSONByte) (*models.Rank, error) {
	for _, rank := range r.ranks {
		if key.Equal(rank.Key) {
			return rank, nil
		}
	}

	return nil, models.ErrNotFound
}

func (r *rankRepo) Close() {}

func (r *rankRepo) Save(ctx context.Context, rank *models.Rank) error {
	return errors.ErrUnsupported
}
