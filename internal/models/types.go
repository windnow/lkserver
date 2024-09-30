package models

import (
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"time"
)

type JSONTime time.Time

var (
	ErrNotFound           = errors.New("NOT FOUND")
	ErrWrongLength        = errors.New("WRONG GUID LENGTH")
	ErrRefIntegrity       = errors.New("REFERENCE INTEGRITY IS VIOLATED")
	ErrInvalidCredentials = errors.New("INVALID CREDENTIALS")
)

func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("02.01.2006"))
	return []byte(stamp), nil
}

func (t *JSONTime) UnmarshalJSON(b []byte) error {

	var dateStr string
	if err := json.Unmarshal(b, &dateStr); err != nil {
		return err
	}

	result, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return err
	}

	*t = JSONTime(result)

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
		result, err := time.Parse("2006-01-02", v)
		if err != nil {
			return HandleError(err, "JSONTime.Scan")
		}
		*t = JSONTime(result)
	default:
		return HandleError(fmt.Errorf("Unsupported type for JSONTime: %T", v))
	}
	return nil
}

func (t JSONTime) Value() (driver.Value, error) {
	return time.Time(t), nil
}

type JSONByte [16]byte

func (uuid JSONByte) MarshalJSON() ([]byte, error) {
	if len(uuid) != 16 {
		return nil, HandleError(ErrWrongLength, "JSONByte.MarshalJSON")
	}
	uuidStr := uuid.String()

	return json.Marshal(uuidStr)
}

func (uuid JSONByte) String() string {
	return (fmt.Sprintf("%08x-%04x-%04x-%04x-%12x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16]))
}

func (uuid JSONByte) Value() (driver.Value, error) {
	if len(uuid) != 16 {
		return nil, HandleError(ErrWrongLength, "JSONByte.Value")
	}
	return uuid[:], nil
}

func (uuid *JSONByte) Scan(src any) error {
	b, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan nonbyte value into JSONByte value")
	}
	if len(b) != 16 {
		return fmt.Errorf("wrong GUID length: expected 16 bytes, got %d", len(b))
	}
	copy(uuid[:], b)
	return nil
}
func (uuid *JSONByte) UnmarshalJSON(data []byte) error {
	var uuidStr string
	if err := json.Unmarshal(data, &uuidStr); err != nil {
		return HandleError(err, "JSONByte.UnmarsharJSON")
	}

	re := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$")
	if !re.MatchString(uuidStr) {
		return HandleError(errors.ErrUnsupported, "JSONByte.UnmarsharJSON")

	}

	noDashes := uuidStr[0:8] + uuidStr[9:13] + uuidStr[14:18] + uuidStr[19:23] + uuidStr[24:]
	parsedUUID, err := hex.DecodeString(noDashes)
	if err != nil {
		return HandleError(err, "JSONByte.UnmarsharJSON")
	}
	*uuid = JSONByte(parsedUUID)

	return nil
}

func (left JSONByte) Equal(righ JSONByte) bool {
	return left == righ
}
