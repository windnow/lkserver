package json

type ContractRepo struct{}

func (r *repo) initContractRepo() error {
	r.contract = &ContractRepo{}
	return nil
}

func (c *ContractRepo) Close() {

}
