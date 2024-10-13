package sqlite

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
	cato         *cato
}
