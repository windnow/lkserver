package json

import (
	"lkserver/internal/repository"
)

type repo struct {
	dataDir         string
	user            *userRepo
	individuals     *individualsRepo
	contract        *contractRepo
	ranks           *rankRepo
	ranksHistory    *rankHistoryRepo
	eduInstitutions *educationInstitutionRepo
	specialties     *specialtiesRepo
	education       *educationRepo
}

func (r *repo) init() error {
	if err := r.initUserRepo(); err != nil {
		return err
	}
	if err := r.initContractRepo(); err != nil {
		return err
	}
	if err := r.initIndividualsRepo(); err != nil {
		return err
	}
	if err := r.initRankRepo(); err != nil {
		return err
	}
	if err := r.initRankHistoryRepo(); err != nil {
		return err
	}

	if err := r.initEducationInstitutions(); err != nil {
		return err
	}

	if err := r.initSpcialties(); err != nil {
		return err
	}

	if err := r.initEducation(); err != nil {
		return err
	}

	return nil
}

func NewJSONProvider(dataDir string) (*repository.Repo, error) {

	r := repo{dataDir: dataDir}
	if err := r.init(); err != nil {
		return nil, err
	}

	return &repository.Repo{
		User:                 r.user,
		Individuals:          r.individuals,
		Contract:             r.contract,
		Ranks:                r.ranks,
		RanksHistory:         r.ranksHistory,
		EducationInstitution: r.eduInstitutions,
		Specialties:          r.specialties,
		Education:            r.education,
	}, nil
}
