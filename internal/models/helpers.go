package models

import (
	"crypto/rand"
	"database/sql"
	"errors"
)

func handleQueryError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

func HandleError(err error, text ...string) error {
	if err == nil {
		return err
	}
	err = handleQueryError(err)
	var resultStr string
	for _, str := range text {
		resultStr += str
	}
	return errors.New("\n" + resultStr + ": " + err.Error())
}
func GenerateUUID() (JSONByte, error) {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		return JSONByte{}, HandleError(err, "GenerageUUID")
	}

	// Устанавливаем версию (4) и вариант UUID
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Версия 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Вариант RFC 4122

	var jsonUUID JSONByte
	copy(jsonUUID[:], uuid)
	return jsonUUID, nil
}
