package render

import (
	"fmt"
	"html/template"
	"io"
	"strings"

	"github.com/labstack/echo/v4"
)

type Render struct {
	templates *template.Template
}

func New(path string) *Render {
	path = strings.Trim(path, "/")
	path = fmt.Sprintf("%s/*", path)

	return &Render{templates: template.Must(template.ParseGlob(path))}
}

func (t *Render) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}
	return t.templates.ExecuteTemplate(w, name, data)
}
