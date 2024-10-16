package sqlite

import "lkserver/internal/repository"

type sqliteRepo struct {
	db           *src
	userRepo     *UserRepository
	contract     *contractRepo
	individuals  *individualsRepo
	ranks        *rankRepo
	rankHistory  *rankHistoryRepo
	institutions *eduInstitutions
	specialties  *specialties
	education    *education
	reports      *reportsRepo
	catalogs     *repository.Catalogs
}
