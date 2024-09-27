package json

import (
	"errors"
	"lkserver/internal/models"
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

func (r *individualsRepo) Get(key models.JSONByte) (*models.Individuals, error) {
	for _, individual := range r.individuals {
		if individual.Key == key {
			return individual, nil
		}
	}
	return nil, models.ErrNotFound
}

func (r *individualsRepo) GetByIin(iin string) (*models.Individuals, error) {
	return nil, errors.ErrUnsupported
}
