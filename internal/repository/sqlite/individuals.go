package sqlite

import (
	"context"
	"encoding/json"
	"errors"
	"lkserver/internal/models"
	"time"
)

type individualsRepo struct {
	source *src
}

func (r *sqliteRepo) initIndividualsRepo() error {
	i := &individualsRepo{
		source: r.db,
	}
	r.individuals = i
	err := i.source.Exec(`
		CREATE TABLE IF NOT EXISTS individuals(
			guid BLOB PRIMARY KEY,
			iin TEXT UNIQUE,
			nationality TEXT,
			first_name TEXT,
			last_name TEXT,
			patronymic TEXT,
			image TEXT,
			birth_date INTEGER,
			birth_place TEXT,
			personal_number TEXT
		)
	`)
	if err != nil {
		return err
	}

	// Создание индексов
	err = i.source.Exec(`
		CREATE INDEX IF NOT EXISTS idx_individuals_iin ON individuals(iin);
		CREATE INDEX IF NOT EXISTS idx_individuals_first_name ON individuals(first_name);
		CREATE INDEX IF NOT EXISTS idx_individuals_last_name ON individuals(last_name);
		CREATE INDEX IF NOT EXISTS idx_individuals_birth_date ON individuals(birth_date);
	`)
	if err != nil {
		return err
	}

	var count int64
	i.source.db.QueryRow(`select count(*) from individuals`).Scan(&count)
	if count == 0 {
		var individuals []models.Individuals
		json.Unmarshal([]byte(data), &individuals)
		for _, individ := range individuals {
			if individ.Key, err = GenerateUUID(); err != nil {
				return err
			}
			if err := i.Save(context.Background(), &individ); err != nil {
				return err
			}

		}
	}

	return nil

}

func (i *individualsRepo) Get(iin string) (*models.Individuals, error) {

	return nil, errors.ErrUnsupported

}

func (i *individualsRepo) Save(ctx context.Context, individ *models.Individuals) error {

	return i.source.ExecContextInTransaction(ctx, insertIndividQuery,
		individ.Key[:],
		individ.IndividualNumber,
		individ.Nationality,
		individ.FirstName,
		individ.LastName,
		individ.Patronymic,
		individ.Image,
		time.Time(individ.BirthDate).Unix(),
		individ.BirthPlace,
		individ.PersonalNumber)

}

var insertIndividQuery = `
	INSERT INTO individuals(
		guid, iin, nationality, first_name, last_name, patronymic, image, birth_date, birth_place, personal_number
	) VALUES (
	 	?,	  ?,   ?,			?,			?,			?,			?,	   ?,			?,			?
	 )
`

var data string = `[
{
        "code": "000000015",
        "iin": "821019000888",
        "nationality": "Казах",
        "first_name": "Дархан",
        "last_name": "Усенбаев",
        "patronymic": "Жаксылыкович",
        "image": "821019000888",
        "birth_date": "1981-11-19",
        "birth_place": "с. Баканас Балхашского района Алма-Атинской области",
        "personal_number": "А-000001"
    },
    {
        "code": "000000016",
        "iin": "910702000888",
        "nationality": "Казах",
        "first_name": "Алинур",
        "last_name": "Асетов",
        "patronymic": "Дастанулы",
        "birth_date": "1991-09-20",
        "birth_place": "с. Баканас Жанааркинского района Карагандинской области",
        "personal_number": "А-000002"
    },
    {
        "code": "000000017",
        "iin": "851204000888",
        "nationality": "Казах",
        "first_name": "Кайрат",
        "last_name": "Каримов",
        "patronymic": "Ганиевич",
        "birth_date": "1985-11-04",
        "birth_place": "с. Октябрьское район М.Жумабаева Северо-Казахстанской области",
        "personal_number": "А-000003"
    }]`
