package entity

import (
	"database/sql/driver"
	"encoding/json"
)

const (
	PayloadQuery = "query"
	PayloadIP    = "ip"
	PayloadAgent = "agent"
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

func (p *Payload) Query() string {
	return (*p)[PayloadQuery]
}

func (p *Payload) IP() string {
	return (*p)[PayloadIP]
}

func (p *Payload) Agent() string {
	return (*p)[PayloadAgent]
}
