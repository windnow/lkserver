package json

import (
	"encoding/json"
	"fmt"
	"io"
	"lkserver/internal/models"
	"log"
	"os"
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

	return nil, nil

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
