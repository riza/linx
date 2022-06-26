package output

import (
	"html/template"
	"os"
)

const templateFile = "./internal/output/output_html_template.html"

type OutputHTML struct {
	output OutputData
}

func NewOutputHTML(output OutputData) OutputHTML {
	return OutputHTML{output: output}
}

func (oh OutputHTML) RenderAndSave() error {
	f, err := os.Create(oh.output.Filename)
	if err != nil {
		return err
	}
	defer f.Close()

	t, err := template.ParseFiles(templateFile)
	if err != nil {
		return err
	}

	err = t.Execute(f, oh.output)
	if err != nil {
		return err
	}

	return nil
}
