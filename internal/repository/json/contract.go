package json

type ContractRepo struct{}

func NewContractRepo(folder string) (*ContractRepo, error) {
	return &ContractRepo{}, nil
}

func (c *ContractRepo) Close() {

}
