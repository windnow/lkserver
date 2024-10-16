package sqlite

import (
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
	repo.catalogs = &repository.Catalogs{}

	if err := repo.initCatalogs(); err != nil {
		return nil, m.HandleError(err, "NewSQLiteProvider")
	}

	if err := repo.initIndividualsRepo(); err != nil {
		return nil, m.HandleError(err, "NewSQLiteProvider")
	}
	if err := repo.initUserRepo(); err != nil {
		return nil, m.HandleError(err, "NewSQLiteProvider")
	}
	if err := repo.initContractRepo(); err != nil {
		return nil, m.HandleError(err, "NewSQLiteProvider")
	}
	if err := repo.initRankRepo(); err != nil {
		return nil, m.HandleError(err, "NewSQLiteProvider")
	}
	if err := repo.initRankHistory(); err != nil {
		return nil, m.HandleError(err, "NewSQLiteProvider")
	}

	if err := repo.initEducationInstitutions(); err != nil {
		return nil, m.HandleError(err, "NewSQLiteProvider")
	}
	if err := repo.initSpecialties(); err != nil {
		return nil, m.HandleError(err, "NewSQLiteProvider")
	}
	if err := repo.initEducation(); err != nil {
		return nil, m.HandleError(err, "NewSQLiteProvider")
	}
	if err := repo.initReports(); err != nil {
		return nil, m.HandleError(err, "NewSQLiteProvider")
	}

	return &repository.Repo{
		User:                 repo.userRepo,
		Contract:             repo.contract,
		Individuals:          repo.individuals,
		Ranks:                repo.ranks,
		RanksHistory:         repo.rankHistory,
		EducationInstitution: repo.institutions,
		Specialties:          repo.specialties,
		Education:            repo.education,
		Reports:              repo.reports,
		Catalogs:             repo.catalogs,
	}, nil
}
