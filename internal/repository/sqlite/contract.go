package sqlite

type contractRepo struct {
	db *src
}

func (c *contractRepo) Close() {
	//nop
}

func (r *sqliteRepo) initContractRepo() error {
	r.contract = &contractRepo{
		db: r.db,
	}
	return nil
}
