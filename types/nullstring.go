package types

import (
	"database/sql"
	"encoding/json"
	"time"
)

type JSONNullString struct {
	sql.NullString
}

func (j *JSONNullString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		j.Valid = false
		j.String = ""
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	j.String = s
	j.Valid = true
	return nil
}

func (j *JSONNullString) MarshalJSON() ([]byte, error) {
	json, err := json.Marshal(j.String)
	return json, err
}

type JSONNullInt64 struct {
	sql.NullInt64
}

func (j *JSONNullInt64) UnmarshalJSON(data []byte) error {
	var i int64
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	j.Int64 = i
	j.Valid = true
	return nil
}

func (j *JSONNullInt64) MarshalJSON() ([]byte, error) {
	json, err := json.Marshal(j.Int64)
	return json, err
}

type JSONNullTime struct {
	sql.NullTime
}

func (j *JSONNullTime) UnmarshalJSON(data []byte) error {
	var i time.Time
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	j.Time = i
	j.Valid = true
	return nil
}

func (j *JSONNullTime) MarshalJSON() ([]byte, error) {
	json, err := json.Marshal(j.Time)
	return json, err
}
