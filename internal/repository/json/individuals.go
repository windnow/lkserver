package json

import (
	"lkserver/internal/models"
	"lkserver/internal/repository"
)

type individualsRepo struct {
	individuals []*models.Individuals
}

func (r *repo) initIndividualsRepo() error {
	repo := &individualsRepo{}
	if err := initFile(r.dataDir+"/individuals.json", &repo.individuals); err != nil {
		return err
	}
	r.individuals = repo
	return nil
}

func (r *individualsRepo) Get(iin string) (*models.Individuals, error) {
	for _, individual := range r.individuals {
		if individual.IndividualNumber == iin {
			return individual, nil
		}
	}
	return nil, repository.ErrNotFound
}
