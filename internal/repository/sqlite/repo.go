package sqlite

type sqliteRepo struct {
	db          *src
	userRepo    *UserRepository
	contract    *contractRepo
	individuals *individualsRepo
	ranks       *rankRepo
}
