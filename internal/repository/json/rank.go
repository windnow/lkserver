package json

import (
	"bytes"
	"context"
	"errors"
	"lkserver/internal/models"
	"lkserver/internal/repository"
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

func (r *rankRepo) Get(key []byte) (*models.Rank, error) {
	for _, rank := range r.ranks {
		if bytes.Equal(rank.Key, key) {
			return rank, nil
		}
	}

	return nil, repository.ErrNotFound
}

func (r *rankRepo) Close() {}

func (r *rankRepo) Save(ctx context.Context, rank *models.Rank) error {
	return errors.ErrUnsupported
}
