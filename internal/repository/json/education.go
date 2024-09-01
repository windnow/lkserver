package json

import (
	"errors"
	"lkserver/internal/models"
	"lkserver/internal/repository"
)

type educationInstitutionRepo struct {
	educationInstitution []*models.EducationInstitution
}

func (c *educationInstitutionRepo) Close() {}
func (e *educationInstitutionRepo) Get(id int) (*models.EducationInstitution, error) {

	for _, inst := range e.educationInstitution {
		if inst.Id == id {
			return inst, nil
		}
	}

	return nil, repository.ErrNotFound
}

type specialtiesRepo struct {
	specialties []*models.Specialties
}

func (c *specialtiesRepo) Close() {}
func (s *specialtiesRepo) Get(id int) (*models.Specialties, error) {

	if s == nil {
		return nil, errors.New("REPO NOT INIT")
	}

	for _, spec := range s.specialties {
		if spec.Id == id {
			return spec, nil
		}
	}

	return nil, repository.ErrNotFound
}

type educationRepo struct {
	education []*models.Education
}

func (e *educationRepo) Close() {}
func (e *educationRepo) Get(iin string) ([]*models.Education, error) {
	var result []*models.Education
	for _, str := range e.education {
		if str.Individual.IndividualNumber == iin {
			result = append(result, str)
		}
	}
	return result, nil
}
func (r *repo) initEducationInstitutions() (err error) {
	repo := &educationInstitutionRepo{}
	err = initFile(r.dataDir+"/education-institutions.json", &repo.educationInstitution)
	if err != nil {
		return err
	}
	r.eduInstitutions = repo

	return
}

func (r *repo) initSpcialties() (err error) {
	repo := &specialtiesRepo{}
	err = initFile(r.dataDir+"/specialties.json", &repo.specialties)
	if err != nil {
		return err
	}
	r.specialties = repo

	return
}

func (r *repo) initEducation() error {
	var result []*models.Education

	var data []struct {
		Iin       string `json:"individual"`
		Institut  int    `json:"institution"`
		Year      int    `json:"year"`
		Specialty int    `json:"specialty"`
		Type      string `json:"type"`
	}
	err := initFile(r.dataDir+"/education.json", &data)
	if err != nil {
		return err
	}

	for _, str := range data {
		individ, err := r.individuals.Get(str.Iin)
		if err != nil {
			return err
		}
		institut, err := r.eduInstitutions.Get(str.Institut)
		if err != nil {
			return err
		}
		spec, err := r.specialties.Get(str.Specialty)
		if err != nil {
			return err
		}

		result = append(result, &models.Education{
			Individual:           individ,
			EducationInstitution: institut,
			YearOfCompletion:     str.Year,
			Specialty:            spec,
			Type:                 str.Type,
		})

	}

	r.education = &educationRepo{education: result}
	return nil
}
