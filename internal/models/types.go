package models

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"time"
)

type JSONTime time.Time

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

type JSONByte []byte

func (uuid JSONByte) MarshalJSON() ([]byte, error) {
	if len(uuid) != 16 {
		return nil, errors.New("WRONG GUID LENGTH")
	}
	uuidStr := (fmt.Sprintf("%08x-%04x-%04x-%04x-%12x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16]))

	return json.Marshal(uuidStr)
}

func (uuid *JSONByte) UnmarshalJSON(data []byte) error {
	var uuidStr string
	if err := json.Unmarshal(data, &uuidStr); err != nil {
		return err
	}

	re := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$")
	if !re.MatchString(uuidStr) {
		return errors.ErrUnsupported
	}

	noDashes := uuidStr[0:8] + uuidStr[9:13] + uuidStr[14:18] + uuidStr[19:23] + uuidStr[24:]
	parsedUUID, err := hex.DecodeString(noDashes)
	if err != nil {
		return err
	}
	*uuid = JSONByte(parsedUUID)

	return nil
}
