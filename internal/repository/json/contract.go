package json

type ContractRepo struct{}

func NewContracRepo(folder string) (*ContractRepo, error) {
	return &ContractRepo{}, nil
}
