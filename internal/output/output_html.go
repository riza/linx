package output

import (
	"github.com/riza/linx/pkg/logger"
	"html/template"
	"os"
)

const htmlExtension = ".html"
const htmlTemplate = `<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>{{ .Target }} - linx report</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.2.0-beta1/dist/css/bootstrap.min.css" rel="stylesheet"
          integrity="sha384-0evHe/X+R7YkIZDRvuzKMRqM+OrBnVFBL6DOitfPri4tjfHxaWutUpFmBp4vmVor" crossorigin="anonymous">
</head>
<body>

<div style="margin-top: 2%;margin-bottom: 5%;" class="container">
    <div class="mb-5 pb-3 fs-4 border-bottom">
        {{ .Target }}
    </div>

    <div class="table-responsive">
        <table style="table-layout: fixed" class="table table-striped table-hover">
            <thead>
            <tr>
                <th scope="col">URL</th>
                <th scope="col">Location in file</th>
            </tr>
            </thead>
            <tbody>
            {{range .Results}}
            <tr>
                <td>{{ .URL }}</td>
                <td>
                    <pre><code>{{ .Location }}</code></pre>
                </td>
            </tr>
            {{end}}
            </tbody>
        </table>
    </div>

    <div class="mt-5 pt-3 text-muted border-top">
        created with <a href="https://github.com/riza/linx">linx</a>
    </div>
</div>

<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.2.0-beta1/dist/js/bootstrap.bundle.min.js"
        integrity="sha384-pprn3073KE6tl6bjs2QrFaJGz5/SUsLqktiwsUTF55Jfv3qYSDhgCecCxMW52nD2"
        crossorigin="anonymous"></script>
</body>
</html>`

type OutputHTML struct {
}

func (oh OutputHTML) RenderAndSave(data OutputData) error {
	fileName := data.Filename + htmlExtension

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	t, err := template.New("output").Parse(htmlTemplate)
	if err != nil {
		return err
	}

	err = t.Execute(f, data)
	if err != nil {
		return err
	}

	logger.Get().Infof("results saved: %s", fileName)
	return nil
}
