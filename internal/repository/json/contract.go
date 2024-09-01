package json

type contractRepo struct{}

func (r *repo) initContractRepo() error {
	r.contract = &contractRepo{}
	return nil
}

func (c *contractRepo) Close() {

}
