package template

import (
	"net/http"
	"github.com/labstack/echo"
	"fmt"
	"html/template"
	"io"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}





func MainPage() echo.HandlerFunc {
	return func(c echo.Context) error { //c をいじって Request, Responseを色々する

		fmt.Println(c)
		return c.String(http.StatusOK, "Hello World")
	}
}
