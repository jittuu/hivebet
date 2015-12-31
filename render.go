package hivebet

import (
	"html/template"
	"io"
)

const layout = "templates/layout.html"

// RenderTemplate renders html templates
func RenderTemplate(w io.Writer, data interface{}, filenames ...string) error {
	filenames = append(filenames, layout)
	t, err := template.New("layout.html").ParseFiles(filenames...)
	if err != nil {
		return err
	}

	if err = t.Execute(w, data); err != nil {
		return err
	}

	return nil
}
