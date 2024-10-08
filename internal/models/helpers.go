package models

import (
	"crypto/rand"
	"fmt"
	"lkserver/internal/lkserver/config"
	"strings"
	"time"
)

type Error struct {
	Err         error
	Description string
}

func (e *Error) Error() string {
	if e.Description == "" {
		return e.Err.Error()
	}

	return fmt.Sprintf("%s: %v", e.Description, e.Err)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func HandleError(err error, description ...string) error {
	if err == nil {
		return nil
	}
	var desc string
	if len(description) > 0 {
		desc = strings.Join(description, ", ")
	}

	return &Error{
		Err:         fmt.Errorf("%w", err),
		Description: desc,
	}
}

func GenerateUUID() (JSONByte, error) {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		return JSONByte{}, &Error{err, "GenerageUUID"}
	}

	// Устанавливаем версию (4) и вариант UUID
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Версия 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Вариант RFC 4122

	var jsonUUID JSONByte
	copy(jsonUUID[:], uuid)
	return jsonUUID, nil
}

func ParseTime(str string) (JSONTime, error) {
	result, err := time.Parse(DateTimeFormat, str)
	if err != nil {
		result, err = time.ParseInLocation(DateFormat, str, config.ServerTimeZone)
		if err != nil {
			return JSONTime{}, err
		}
	}
	return JSONTime(result.In(time.UTC)), nil
}
