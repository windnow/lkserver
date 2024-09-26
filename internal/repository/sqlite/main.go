package sqlite

import (
	"errors"
	"lkserver/internal/repository"
)

func NewSQLiteProvider(dataFile string) (*repository.Repo, error) {
	db, err := newDB(dataFile)
	if err != nil {
		return nil, err
	}
	repo := &sqliteRepo{
		db: db,
	}

	if err := repo.initUserRepo(); err != nil {
		return nil, err
	}
	if err := repo.initContractRepo(); err != nil {
		return nil, err
	}
	if err := repo.initIndividualsRepo(); err != nil {
		return nil, err
	}
	if err := repo.initRankRepo(); err != nil {
		return nil, err
	}

	return &repository.Repo{
		User:        repo.userRepo,
		Contract:    repo.contract,
		Individuals: repo.individuals,
		Ranks:       repo.ranks,
	}, errors.New("Testing")
}
