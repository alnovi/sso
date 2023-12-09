package template

import (
	"fmt"
	"html/template"
	"io"
	"strings"

	"github.com/labstack/echo/v4"
)

type HtmlRenderer struct {
	templates *template.Template
}

func NewHtmlRenderer(path string) *HtmlRenderer {
	path = strings.Trim(path, "/")
	path = fmt.Sprintf("%s/*", path)

	return &HtmlRenderer{templates: template.Must(template.ParseGlob(path))}
}

func (t *HtmlRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}
