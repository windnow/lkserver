package json

import (
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

func (r *rankRepo) Get(id int) (*models.Rank, error) {
	for _, rank := range r.ranks {
		if rank.Id == id {
			return rank, nil
		}
	}

	return nil, repository.ErrNotFound
}

func (r *rankRepo) Close() {}
