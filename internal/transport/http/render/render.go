package render

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"strings"

	"github.com/labstack/echo/v4"
)

type Render struct {
	templates *template.Template
}

func NewFromPath(path string) *Render {
	path = strings.Trim(path, "/")
	path = fmt.Sprintf("%s/*", path)

	return &Render{templates: template.Must(template.ParseGlob(path))}
}

func NewFromFS(fs fs.FS, dir string) *Render {
	dir = strings.Trim(dir, "/")
	pattern := fmt.Sprintf("%s/*", dir)

	return &Render{templates: template.Must(template.ParseFS(fs, pattern))}
}

func (t *Render) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}
	return t.templates.ExecuteTemplate(w, name, data)
}
