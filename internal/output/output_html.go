package output

import (
	"html/template"
	"os"

	"github.com/riza/linx/pkg/logger"
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
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.8.3/font/bootstrap-icons.css">
    <style>
        .highlight {
            background-color: #ffff99;
            padding: 2px;
            border-radius: 3px;
        }
        .url-container {
            word-break: break-all;
        }
        pre {
            max-height: 150px;
            overflow-y: auto;
            background-color: #f8f9fa;
            padding: 10px;
            border-radius: 5px;
            font-size: 0.85rem;
        }
        .copy-btn {
            cursor: pointer;
        }
        .filter-container {
            margin-bottom: 20px;
        }
        #stats {
            font-size: 0.9rem;
            margin-bottom: 15px;
        }
        .type-badge {
            font-size: 0.7rem;
            margin-left: 5px;
        }
        th {
            position: sticky;
            top: 0;
            background-color: #fff;
            z-index: 1;
        }
        .table-container {
            max-height: 80vh;
            overflow-y: auto;
        }
    </style>
</head>
<body>

<div class="container-fluid py-4">
    <div class="row mb-3">
        <div class="col-12">
            <h4>
                <i class="bi bi-link-45deg"></i> Linx Report 
                <small class="text-muted">{{ .Target }}</small>
            </h4>
        </div>
    </div>
    
    <div class="filter-container">
        <div class="row g-3">
            <div class="col-md-6">
                <div class="input-group">
                    <span class="input-group-text"><i class="bi bi-search"></i></span>
                    <input type="text" class="form-control" id="searchInput" placeholder="Search URLs...">
                    <button class="btn btn-outline-secondary" type="button" id="clearSearch"><i class="bi bi-x-lg"></i></button>
                </div>
            </div>
            <div class="col-md-6">
                <div class="input-group">
                    <span class="input-group-text">Filter</span>
                    <select class="form-select" id="typeFilter">
                        <option value="all">All Types</option>
                        <option value="api">API Endpoints</option>
                        <option value="static">Static Resources</option>
                        <option value="external">External URLs</option>
                        <option value="relative">Relative Paths</option>
                    </select>
                </div>
            </div>
        </div>
    </div>
    
    <div id="stats" class="alert alert-info">
        <strong>Total URLs found:</strong> <span id="totalCount">{{ len .Results }}</span> | 
        <strong>Displayed:</strong> <span id="displayedCount">{{ len .Results }}</span>
    </div>

    <div class="table-container">
        <table class="table table-striped table-hover">
            <thead class="table-light">
            <tr>
                <th scope="col" style="width: 40%">URL</th>
                <th scope="col" style="width: 60%">Context</th>
            </tr>
            </thead>
            <tbody id="resultsTable">
            {{range $index, $result := .Results}}
            <tr class="result-row" data-url="{{ .URL }}">
                <td class="url-container">
                    <div class="d-flex justify-content-between">
                        <div>
                            <span class="url-text">{{ .URL }}</span>
                            <span class="type-badge badge bg-secondary" data-type="unknown">analyzing...</span>
                        </div>
                        <div>
                            <i class="bi bi-clipboard copy-btn" data-clipboard-text="{{ .URL }}" title="Copy URL"></i>
                        </div>
                    </div>
                </td>
                <td>
                    <pre><code>{{ .Location }}</code></pre>
                </td>
            </tr>
            {{end}}
            </tbody>
        </table>
    </div>

    <div class="mt-4 pt-3 text-muted border-top">
        <div class="d-flex justify-content-between">
            <div>Generated with <a href="https://github.com/riza/linx" target="_blank">linx</a></div>
            <div><button class="btn btn-sm btn-primary" id="exportBtn"><i class="bi bi-download"></i> Export JSON</button></div>
        </div>
    </div>
</div>

<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.2.0-beta1/dist/js/bootstrap.bundle.min.js"
        integrity="sha384-pprn3073KE6tl6bjs2QrFaJGz5/SUsLqktiwsUTF55Jfv3qYSDhgCecCxMW52nD2"
        crossorigin="anonymous"></script>
<script>
    document.addEventListener('DOMContentLoaded', function() {
        const results = [
            {{range .Results}}
            { url: '{{.URL}}', location: '{{.Location}}' },
            {{end}}
        ];
        
        // Analyze URL types
        analyzeURLs();
        
        // Setup search
        const searchInput = document.getElementById('searchInput');
        searchInput.addEventListener('input', filterResults);
        
        // Setup clear button
        document.getElementById('clearSearch').addEventListener('click', () => {
            searchInput.value = '';
            filterResults();
        });
        
        // Setup type filter
        document.getElementById('typeFilter').addEventListener('change', filterResults);
        
        // Setup copy buttons
        document.querySelectorAll('.copy-btn').forEach(btn => {
            btn.addEventListener('click', function() {
                const text = this.dataset.clipboardText;
                navigator.clipboard.writeText(text).then(() => {
                    const originalClass = this.className;
                    this.className = 'bi bi-check-lg text-success';
                    setTimeout(() => {
                        this.className = originalClass;
                    }, 1500);
                });
            });
        });
        
        // Export JSON button
        document.getElementById('exportBtn').addEventListener('click', function() {
            const data = {
                target: '{{ .Target }}',
                results: results
            };
            
            const blob = new Blob([JSON.stringify(data, null, 2)], {type: 'application/json'});
            const url = URL.createObjectURL(blob);
            
            const a = document.createElement('a');
            a.href = url;
            a.download = '{{ .Target }}_linx_results.json';
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
            URL.revokeObjectURL(url);
        });
        
        function analyzeURLs() {
            document.querySelectorAll('.result-row').forEach(row => {
                const url = row.dataset.url;
                const badge = row.querySelector('.type-badge');
                
                // Determine URL type
                let type = 'unknown';
                let badgeClass = 'bg-secondary';
                
                if (url.match(/^(https?:)?\/\//) && url.match(/api|graphql|service|\/v[0-9]+\//i)) {
                    type = 'api';
                    badgeClass = 'bg-danger';
                } else if (url.match(/^(https?:)?\/\//)) {
                    type = 'external';
                    badgeClass = 'bg-primary';
                } else if (url.match(/\.(js|css|jpg|jpeg|png|gif|svg|woff|ttf)/i)) {
                    type = 'static';
                    badgeClass = 'bg-success';
                } else if (url.match(/^(\/|\.\/)/) || !url.match(/^[a-z]+:/i)) {
                    type = 'relative';
                    badgeClass = 'bg-warning text-dark';
                }
                
                badge.textContent = type;
                badge.dataset.type = type;
                badge.className = 'type-badge badge ' + badgeClass;
                row.dataset.type = type;
            });
        }
        
        function filterResults() {
            const searchTerm = searchInput.value.toLowerCase();
            const typeFilter = document.getElementById('typeFilter').value;
            let displayCount = 0;
            
            document.querySelectorAll('.result-row').forEach(row => {
                const url = row.dataset.url.toLowerCase();
                const type = row.dataset.type;
                
                const matchesSearch = url.includes(searchTerm);
                const matchesType = typeFilter === 'all' || type === typeFilter;
                
                if (matchesSearch && matchesType) {
                    row.style.display = '';
                    displayCount++;
                    
                    // Highlight matching text if there's a search term
                    if (searchTerm) {
                        const urlElement = row.querySelector('.url-text');
                        const originalText = row.dataset.url;
                        const regex = new RegExp('(' + escapeRegExp(searchTerm) + ')', 'gi');
                        urlElement.innerHTML = originalText.replace(regex, '<span class="highlight">$1</span>');
                    } else {
                        const urlElement = row.querySelector('.url-text');
                        urlElement.textContent = row.dataset.url;
                    }
                } else {
                    row.style.display = 'none';
                }
            });
            
            document.getElementById('displayedCount').textContent = displayCount;
        }
        
        function escapeRegExp(string) {
            return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
        }
    });
</script>
</body>
</html>`

type OutputHTML struct {
}

func (oh OutputHTML) RenderAndSave(data *OutputData) error {
	f, err := os.Create(data.Filename)
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

	logger.Get().Infof("results saved: %s", data.Filename)
	return nil
}
