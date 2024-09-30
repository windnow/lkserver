package sqlite

import (
	"errors"
	"fmt"
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

	h, _ := repo.rankHistory.GetLastByIin("821019000888")
	fmt.Printf("%v\n", h)

	return &repository.Repo{
		User:         repo.userRepo,
		Contract:     repo.contract,
		Individuals:  repo.individuals,
		Ranks:        repo.ranks,
		RanksHistory: repo.rankHistory,
	}, errors.New("Testing")
}
