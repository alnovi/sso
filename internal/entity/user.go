package entity

import (
	"fmt"
	"strings"
	"time"
)

type User struct {
	ID        string
	Image     *string
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (e *User) MaskEmail() string {
	if e.Email == "" {
		return ""
	}

	atPos := strings.LastIndex(e.Email, "@")
	domain := e.Email[atPos+1:]
	name := e.Email[:atPos]

	start, end := "", ""
	lenName := len(name)
	if lenName > 5 { // nolint:gomnd
		start = name[:2]
		end = name[lenName-2:]
	}

	return fmt.Sprintf("%s***%s@%s", start, end, domain)
}
