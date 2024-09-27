package sqlite

import (
	"errors"
	m "lkserver/internal/models"
	"lkserver/internal/repository"
)

func NewSQLiteProvider(dataFile string) (*repository.Repo, error) {
	db, err := newDB(dataFile)
	if err != nil {
		return nil, m.HandleError(err, "NewSQLiteProvider")
	}
	repo := &sqliteRepo{
		db: db,
	}

	if err := repo.initUserRepo(); err != nil {
		return nil, m.HandleError(err, "NewSQLiteProvider")
	}
	if err := repo.initContractRepo(); err != nil {
		return nil, m.HandleError(err, "NewSQLiteProvider")
	}
	if err := repo.initIndividualsRepo(); err != nil {
		return nil, m.HandleError(err, "NewSQLiteProvider")
	}
	if err := repo.initRankRepo(); err != nil {
		return nil, m.HandleError(err, "NewSQLiteProvider")
	}
	if err := repo.initRankHistory(); err != nil {
		return nil, m.HandleError(err, "NewSQLiteProvider")
	}

	return &repository.Repo{
		User:        repo.userRepo,
		Contract:    repo.contract,
		Individuals: repo.individuals,
		Ranks:       repo.ranks,
	}, errors.New("Testing")
}
