package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"lkserver/internal/models"
	"log"
	"os"
	"time"
)

type User struct {
	Name      string `json:"name"`
	Iin       string `json:"iin"`
	Pin       string `json:"pin"`
	BirthDate string `json:"birth_date"`
}
type JsonRepo struct {
	dataDir string
	users   []User
}

func (r *JsonRepo) init() error {
	log.Printf("Init JSON storage (%s)", r.dataDir)

	bytes, err := readData(fmt.Sprintf("%s/users.json", r.dataDir))
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, &r.users)

	for _, user := range r.users {
		log.Printf("User: %s", user.Name)
	}

	return err
}

func (r *JsonRepo) GetUser(iin, pin string) (*models.User, error) {

	for _, user := range r.users {
		if user.Iin == iin && user.Pin == pin {
			BirthDate, err := time.Parse("2006-01-02", user.BirthDate)
			if err != nil {
				return nil, err
			}
			return &models.User{
				Name:      user.Name,
				Iin:       user.Iin,
				Pin:       user.Pin,
				BirthDate: models.JSONTime(BirthDate),
			}, nil
		}
	}
	return nil, errors.New("NOT FOUND")

}

func (r *JsonRepo) Close() {}

func New(config string) (*JsonRepo, error) {
	repo := &JsonRepo{dataDir: config}
	if err := repo.init(); err != nil {
		return nil, err
	}
	return repo, nil
}

func readData(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return bytes, nil

}
