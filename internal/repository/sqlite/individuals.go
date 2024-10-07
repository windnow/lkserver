package sqlite

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	m "lkserver/internal/models"
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
	err := i.source.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %[1]s(
			ref BLOB PRIMARY KEY,
			iin TEXT UNIQUE,
			code TEXT DEFAULT "" NOT NULL,
			nationality TEXT,
			first_name TEXT,
			last_name TEXT,
			patronymic TEXT,
			image TEXT,
			birth_date INTEGER,
			birth_place TEXT,
			personal_number TEXT
		);
		CREATE INDEX IF NOT EXISTS idx_%[1]s_iin ON %[1]s(iin);
		CREATE INDEX IF NOT EXISTS idx_%[1]s_first_name ON %[1]s(first_name);
		CREATE INDEX IF NOT EXISTS idx_%[1]s_last_name ON %[1]s(last_name);
		CREATE INDEX IF NOT EXISTS idx_%[1]s_birth_date ON %[1]s(birth_date);
	`, tabIndividuals))
	if err != nil {
		return m.HandleError(err, "sqliteRepo.initIndividualsRepo")
	}

	var count int64
	i.source.db.QueryRow(fmt.Sprintf(`select count(*) from %[1]s`, tabIndividuals)).Scan(&count)
	var individuals []m.Individuals
	if err := json.Unmarshal([]byte(data), &individuals); err != nil {
		return m.HandleError(err, "sqliteRepo.initIndividualsRepo")
	}
	if count == 0 {
		for _, individ := range individuals {
			if err := i.Save(context.Background(), &individ); err != nil {
				return m.HandleError(err, "sqliteRepo.initIndividualsRepo")
			}

		}
	}

	return nil

}

func (i *individualsRepo) Get(key m.JSONByte) (*m.Individuals, error) {

	individ := &m.Individuals{Key: key}
	err := i.source.db.QueryRow(fmt.Sprintf(`
		SELECT iin,  
			code,
			nationality,
			first_name,
			last_name,
			patronymic,
			image,
			birth_date,
			birth_place,
			personal_number
		FROM %[1]s
		WHERE ref = ? 
	`, tabIndividuals), individ.Key).Scan(
		&individ.IndividualNumber,
		&individ.Code,
		&individ.Nationality,
		&individ.FirstName,
		&individ.LastName,
		&individ.Patronymic,
		&individ.Image,
		&individ.BirthDate,
		&individ.BirthPlace,
		&individ.PersonalNumber,
	)

	if err != nil {
		return nil, m.HandleError(err, "individualsRepo.Get")
	}

	return individ, nil

}

func (i *individualsRepo) GetByIin(iin string) (*m.Individuals, error) {
	individ := &m.Individuals{IndividualNumber: iin}
	err := i.source.db.QueryRow(fmt.Sprintf(`
		SELECT ref,  
			code,
			nationality,
			first_name,
			last_name,
			patronymic,
			image,
			birth_date,
			birth_place,
			personal_number
		FROM %[1]s
		WHERE iin = ? 
	`, tabIndividuals), individ.IndividualNumber).Scan(
		&individ.Key,
		&individ.Code,
		&individ.Nationality,
		&individ.FirstName,
		&individ.LastName,
		&individ.Patronymic,
		&individ.Image,
		&individ.BirthDate,
		&individ.BirthPlace,
		&individ.PersonalNumber,
	)

	if err != nil {
		return nil, m.HandleError(errors.ErrUnsupported, "individualsRepo.GetByIin")
	}
	return individ, nil

}

func (i *individualsRepo) Save(ctx context.Context, individ *m.Individuals) error {

	return i.source.ExecContextInTransaction(ctx, insertIndividQuery, nil,
		individ.Key,
		individ.IndividualNumber,
		individ.Code,
		individ.Nationality,
		individ.FirstName,
		individ.LastName,
		individ.Patronymic,
		individ.Image,
		time.Time(individ.BirthDate).Unix(),
		individ.BirthPlace,
		individ.PersonalNumber)

}

var insertIndividQuery = fmt.Sprintf(`
	INSERT INTO %s(
		ref, iin, code, nationality, first_name, last_name, patronymic, image, birth_date, birth_place, personal_number
	) VALUES (
	 	?,	  ?,   ?,	?,      		?,			?,			?,			?,	   ?,			?,			?
	 )
`, tabIndividuals)

var data string = `[
{
		"key":"27f74b66-cba7-486d-a263-81b6cb9a3e57",
        "code": "000000015",
        "iin": "821019000888",
        "nationality": "Казах",
        "first_name": "Дархан",
        "last_name": "Усенбаев",
        "patronymic": "Жаксылыкович",
        "image": "821019000888",
        "birth_date": "1981.11.19",
        "birth_place": "с. Баканас Балхашского района Алма-Атинской области",
        "personal_number": "А-000001"
    },
    {
		"key":"52efc72d-ba0d-4f87-ae73-e902936395fe",
        "code": "000000016",
        "iin": "910702000888",
        "nationality": "Казах",
        "first_name": "Алинур",
        "last_name": "Асетов",
        "patronymic": "Дастанулы",
        "birth_date": "1991.09.20",
        "birth_place": "с. Баканас Жанааркинского района Карагандинской области",
        "personal_number": "А-000002"
    },
    {
		"key":"19db2753-68f9-4b5d-998a-727e347a958a",
        "code": "000000017",
        "iin": "851204000888",
        "nationality": "Казах",
        "first_name": "Кайрат",
        "last_name": "Каримов",
        "patronymic": "Ганиевич",
        "birth_date": "1985.11.04",
        "birth_place": "с. Октябрьское район М.Жумабаева Северо-Казахстанской области",
        "personal_number": "А-000003"
    }]`
