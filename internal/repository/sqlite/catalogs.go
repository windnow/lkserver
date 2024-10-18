package sqlite

import m "lkserver/internal/models"

func (repo *sqliteRepo) initCatalogs() error {
	if err := repo.initCato(); err != nil {
		return m.HandleError(err, "sqliteRepo.initCatalogs")
	}

	if err := repo.initVus(); err != nil {
		return m.HandleError(err, "sqliteRepo.initCatalogs")
	}

	if err := repo.initOrganization(); err != nil {
		return m.HandleError(err, "sqliteRepo.initCatalogs")
	}

	if err := repo.initDevision(); err != nil {
		return m.HandleError(err, "sqliteRepo.initCatalogs")
	}

	return nil
}
