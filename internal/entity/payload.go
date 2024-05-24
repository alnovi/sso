package entity

import (
	"database/sql/driver"
	"encoding/json"
)

type Payload map[string]string

func (e *Payload) Scan(val interface{}) error {
	return json.Unmarshal(val.([]byte), e)
}

func (e *Payload) Value() (driver.Value, error) {
	return json.Marshal(e)
}
