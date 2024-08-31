package json

import (
	"fmt"
	"lkserver/internal/models"
	"lkserver/internal/repository"
)

type IndividualsRepo struct {
	dataDir     string
	individuals []*models.Individuals
}

func (r *IndividualsRepo) init() error {
	return initFile(
		fmt.Sprintf("%s/individuals.json", r.dataDir),
		&r.individuals,
	)
}

func (r *repo) initIndividualsRepo() error {
	repo := &IndividualsRepo{dataDir: r.dataDir}
	if err := repo.init(); err != nil {
		return err
	}
	r.individuals = repo
	return nil
}

func (r *IndividualsRepo) Get(iin string) (*models.Individuals, error) {
	for _, individual := range r.individuals {
		if individual.IndividualNumber == iin {
			return individual, nil
		}
	}
	return nil, repository.ErrNotFound
}
