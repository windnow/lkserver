package sqlite

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"lkserver/internal/repository"
)

func GenerateUUID() ([]byte, error) {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		return nil, err
	}

	// Устанавливаем версию (4) и вариант UUID
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Версия 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Вариант RFC 4122

	return uuid, nil
}

func handleQueryError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return repository.ErrNotFound
	}
	return err
}
