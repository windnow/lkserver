package json

import (
	"fmt"
	"lkserver/internal/models"
	"lkserver/internal/repository"
)

type RankRepo struct {
	dataDir string
	ranks   []*models.Rank
}

func (r *repo) initRankRepo() error {
	repo := &RankRepo{dataDir: r.dataDir}
	if err := repo.init(); err != nil {
		return err
	}

	r.ranks = repo
	return nil
}

func (r *RankRepo) Get(id int) (*models.Rank, error) {
	for _, rank := range r.ranks {
		if rank.Id == id {
			return rank, nil
		}
	}

	return nil, repository.ErrNotFound
}

func (r *RankRepo) Close() {}

func (r *RankRepo) init() error {
	return initFile(
		fmt.Sprintf("%s/ranks.json", r.dataDir),
		&r.ranks,
	)
}
