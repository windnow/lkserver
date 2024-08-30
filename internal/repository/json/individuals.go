package json

import (
	"errors"
	"fmt"
	"lkserver/internal/models"
)

type IndividualsRepo struct {
	dataDir     string
	individuals []models.Individuals
}

func (r *IndividualsRepo) init() error {
	return initFile(
		fmt.Sprintf("%s/individuals.json", r.dataDir),
		&r.individuals,
	)
}

func NewIndividualsRepo(dataDir string) (*IndividualsRepo, error) {
	repo := &IndividualsRepo{dataDir: dataDir}
	if err := repo.init(); err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *IndividualsRepo) Get(iin string) (*models.Individuals, error) {
	for _, individual := range r.individuals {
		if individual.IndividualNumber == iin {
			return &individual, nil
		}
	}
	return nil, errors.New("NOT FOUND")
}
