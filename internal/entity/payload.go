package entity

import (
	"database/sql/driver"
	"encoding/json"
)

type Payload map[string]string

func (p *Payload) Scan(value interface{}) error {
	if str, ok := value.(string); ok {
		return json.Unmarshal([]byte(str), p)
	}
	return nil
}

func (p *Payload) Value() (driver.Value, error) {
	return json.Marshal(p)
}
