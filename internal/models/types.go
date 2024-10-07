package models

import (
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"lkserver/internal/lkserver/config"
	"regexp"
	"time"
)

type CtxKey string

type JSONTime time.Time

var (
	ErrNotFound           = errors.New("NOT FOUND")
	ErrWrongLength        = errors.New("WRONG GUID LENGTH")
	ErrRefIntegrity       = errors.New("REFERENCE INTEGRITY IS VIOLATED")
	ErrInvalidCredentials = errors.New("INVALID CREDENTIALS")
	DateTimeFormat        = "2006.01.02 15:04:05"
	DateFormat            = "2006.01.02"
)

func (t JSONTime) MarshalJSON() ([]byte, error) {

	stamp := fmt.Sprintf("\"%s\"", time.Time(t).In(config.ServerTimeZone).Format(DateTimeFormat))
	return []byte(stamp), nil
}

func (t *JSONTime) UnmarshalJSON(b []byte) error {

	var dateStr string
	if err := json.Unmarshal(b, &dateStr); err != nil {
		return &Error{err, "JSONTime.UnmarshalJSON"}
	}

	result, err := ParseTime(dateStr)
	if err != nil {
		return err
	}

	*t = result

	return nil
}

func (t JSONTime) After(u JSONTime) bool {
	return time.Time(t).After(time.Time(u))
}

func (t *JSONTime) Scan(value any) error {
	switch v := value.(type) {
	case time.Time:
		*t = JSONTime(v)
	case int64:
		*t = JSONTime(time.Unix(v, 0))
	case string:
		result, err := ParseTime(v)
		if err != nil {
			return &Error{err, "JSONTime.Scan"}
		}
		*t = result
	default:
		return &Error{fmt.Errorf("UNSUPPORTED TYPE FOR JSONTime: %T", v), "JSONTime.Scan"}
	}
	return nil
}

func (t JSONTime) Value() (driver.Value, error) {
	return time.Time(t), nil
}

type JSONByte [16]byte

func (uuid JSONByte) MarshalJSON() ([]byte, error) {
	if uuid.Blank() {
		return json.Marshal("")
	}

	uuidStr := uuid.String()
	return json.Marshal(uuidStr)
}

func (uuid JSONByte) Blank() bool {
	return uuid == [16]byte{}
}

func (uuid JSONByte) String() string {
	return (fmt.Sprintf("%08x-%04x-%04x-%04x-%12x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16]))
}

func (uuid JSONByte) Value() (driver.Value, error) {
	if len(uuid) != 16 {
		return nil, &Error{ErrWrongLength, "JSONByte.Value"}
	}
	if uuid.Blank() {
		return nil, nil
	}
	return uuid[:], nil
}

func (uuid *JSONByte) Scan(src any) error {
	b, ok := src.([]byte)
	if !ok {
		return &Error{fmt.Errorf("cannot scan nonbyte value into JSONByte value"), "JSONByte.Scan"}
	}
	if len(b) != 16 {
		return &Error{fmt.Errorf("wrong GUID length: expected 16 bytes, got %d", len(b)), "JSONByte.Scan"}
	}
	copy(uuid[:], b)
	return nil
}
func (uuid *JSONByte) UnmarshalJSON(data []byte) error {
	var uuidStr string
	if err := json.Unmarshal(data, &uuidStr); err != nil {
		return &Error{err, "JSONByte.UnmarsharJSON"}
	}

	if uuidStr == "" {
		*uuid = JSONByte{}
		return nil
	}

	re := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$")
	if !re.MatchString(uuidStr) {
		return &Error{errors.ErrUnsupported, "JSONByte.UnmarsharJSON"}

	}

	noDashes := uuidStr[0:8] + uuidStr[9:13] + uuidStr[14:18] + uuidStr[19:23] + uuidStr[24:]
	parsedUUID, err := hex.DecodeString(noDashes)
	if err != nil {
		return &Error{err, "JSONByte.UnmarsharJSON"}
	}
	*uuid = JSONByte(parsedUUID)

	return nil
}

func (left JSONByte) Equal(righ JSONByte) bool {
	return left == righ
}
