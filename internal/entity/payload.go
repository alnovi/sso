package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mileusna/useragent"
)

type Payload map[string]string

func (e *Payload) Scan(val interface{}) error {
	return json.Unmarshal(val.([]byte), e)
}

func (e *Payload) Value() (driver.Value, error) {
	return json.Marshal(e)
}

func (e *Payload) IP() string {
	return (*e)["ip"]
}

func (e *Payload) Agent() string {
	return (*e)["agent"]
}

func (e *Payload) Device() string {
	ua := useragent.Parse(e.Agent())
	device := fmt.Sprintf("%s %s %s", ua.OS, ua.OSVersion, ua.Name)
	device = strings.Replace(device, "  ", " ", -1)
	return strings.Trim(device, " ")
}
