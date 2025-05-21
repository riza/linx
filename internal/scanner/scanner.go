package scanner

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"unsafe"

	"github.com/riza/linx/internal/options"
	"github.com/riza/linx/internal/output"
	"github.com/riza/linx/internal/scanner/strategies"
	"github.com/riza/linx/pkg/logger"
)

// rule from LinkFinder https://github.com/GerbenJavado/LinkFinder/blob/master/linkfinder.py#L29 ty @GerbenJavado
// Enhanced to catch more edge cases and modern URL patterns
const (
	// Main rule enhanced to catch more patterns
	rule = `(?:"|'|(?:url\()|(?:URL\()|(?:href=)|(?:src=)|(?:data-url=)|(?:location(?:\s*=|\.href\s*=))|(?:\.open\(\s*))((\b(?:[a-zA-Z]{1,20}://|//)[^"'/]{1,}\.[a-zA-Z]{2,}[^"']{0,})|(?:data:(?:image|application|text)/[a-zA-Z0-9;,.+]{1,20}base64[^"']{0,})|([a-zA-Z0-9_\-]{1,}\.(?:php|asp|aspx|jsp|json|action|html|js|txt|xml|do|py|rb)(?:\?[^"|']{0,}|))|(?:(?:/|\.\./|\./)[^"'><,;| *()(%%$^/\\\[\]][^"'><,;|()]{1,})|(?:[a-zA-Z0-9_\-/]{1,}/[a-zA-Z0-9_\-/]{1,}\.(?:[a-zA-Z]{1,4}|action)(?:[\?|/|#][^"|']{0,}|))|(?:[a-zA-Z0-9_\-/]{1,}/[a-zA-Z0-9_\-/]{3,}(?:[\?|/|#][^"|']{0,}|))|(?:[a-zA-Z0-9_\-]{1,}\.(?:php|asp|aspx|jsp|json|action|html|js|txt|xml|do|py|rb)(?:[\?|/|#][^"|']{0,}|)))(?:"|'|\))`

	// Additional rule for URL patterns using template literals (backticks) in modern JS
	templateLiteralRule = "`((?:(?:https?:)?//)[^`\r\n]{1,})`"

	// Additional rule for URL patterns in config objects
	configObjectRule = `["'](endpoint|url|href|src|uri|path|target)["']\s*:\s*["']((?:(?:https?:)?//|/)[^"'\r\n]{1,})["']`

	// Additional rule for API endpoints
	apiEndpointRule = `["'](/api/v?[0-9.]*?/[^"'\r\n]{1,})["']`

	// Additional rule for URLs in JavaScript variable assignments
	assignmentRule = `(?:const|let|var)\s+\w+\s*=\s*["']((?:(?:https?:)?//|/)[^"'\r\n]{1,})["']`

	// Enhanced exclude rule with more patterns
	excludeFileTypeRule = `.css|.jpg|.jpeg|.png|.svg|.img|.gif|.mp4|.flv|.ogv|.webm|.webp|.mov|.mp3|.m4a|.m4p|.scss|.tif|.tiff|.ttf|.otf|.woff|.woff2|.bmp|.ico|.eot|.htc|.rtf|.swf|.image|w3.org|doubleclick.net|youtube.com|.vue|jquery|bootstrap|font|jsdelivr.net|vimeo.com|pinterest.com|facebook|linkedin|twitter|instagram|google|mozilla.org|jibe.com|schema.org|schemas.microsoft.com|wordpress.org|w.org|wix.com|parastorage.com|whatwg.org|polyfill.io|typekit.net|schemas.openxmlformats.org|openweathermap.org|openoffice.org|reactjs.org|angularjs.org|java.com|purl.org|/image|/img|/css|/wp-json|/wp-content|/wp-includes|/theme|/audio|/captcha|/font|robots.txt|node_modules|.wav|.gltf`

	excludeMimeTypeRule = `text/css|image/jpeg|image/jpg|image/png|image/svg+xml|image/gif|image/tiff|image/webp|image/bmp|image/x-icon|image/vnd.microsoft.icon|font/ttf|font/woff|font/woff2|font/x-woff2|font/x-woff|font/otf|audio/mpeg|audio/wav|audio/webm|audio/aac|audio/ogg|audio/wav|audio/webm|video/mp4|video/mpeg|video/webm|video/ogg|video/mp2t|video/webm|video/x-msvideo|application/font-woff|application/font-woff2|application/vnd.android.package-archive|binary/octet-stream|application/octet-stream|application/pdf|application/x-font-ttf|application/x-font-otf|application/json|text/javascript|text/plain|text/x-yaml|text/html|text/babel|text/markdown|text/tsx|application/typescript|application/javascript|text/x-handlebars-template|application/x-typescript|text/x-gfm|text/jsx`
)

var (
	outputEngines = map[string]output.Output{
		"":      output.OutputNoop{},
		".html": output.OutputHTML{},
		".json": output.OutputJSON{},
	}

	// Compile all rules at init
	patterns []*regexp.Regexp
)

func init() {
	// Compile all patterns for efficient reuse
	patterns = []*regexp.Regexp{
		regexp.MustCompile(rule),
		regexp.MustCompile(templateLiteralRule),
		regexp.MustCompile(configObjectRule),
		regexp.MustCompile(apiEndpointRule),
		regexp.MustCompile(assignmentRule),
	}
}

type task struct {
	target   string
	output   string
	strategy strategies.ScanStrategy
}

type scanner struct {
	task task
	opts *options.Options
}

func NewScanner(opts *options.Options) scanner {
	return scanner{
		task: task{
			target:   opts.Target,
			output:   opts.Output,
			strategy: defineStrategyForTarget(opts.Target),
		},
		opts: opts,
	}
}

func (s scanner) Run() error {
	rFt, _ := regexp.Compile(excludeFileTypeRule)
	rMt, _ := regexp.Compile(excludeMimeTypeRule)

	// Multiple targets support
	targets := strings.Split(s.task.target, ",")

	if len(targets) > 1 && s.opts.Parallel {
		return s.runParallel(targets, rFt, rMt)
	}

	if len(targets) > 1 {
		for _, target := range targets {
			s.task.target = strings.TrimSpace(target)
			s.task.strategy = defineStrategyForTarget(s.task.target)
			if err := s.processTarget(rFt, rMt); err != nil {
				logger.Get().Errorf("Error processing target %s: %v", target, err)
			}
		}
		return nil
	}

	return s.processTarget(rFt, rMt)
}

func (s scanner) runParallel(targets []string, rFt, rMt *regexp.Regexp) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(targets))

	for _, target := range targets {
		wg.Add(1)
		go func(target string) {
			defer wg.Done()

			scannerCopy := scanner{
				task: task{
					target:   strings.TrimSpace(target),
					output:   s.task.output + "." + filepath.Base(strings.TrimSpace(target)),
					strategy: defineStrategyForTarget(strings.TrimSpace(target)),
				},
				opts: s.opts,
			}

			if err := scannerCopy.processTarget(rFt, rMt); err != nil {
				errChan <- fmt.Errorf("Error processing target %s: %v", target, err)
			}
		}(target)
	}

	wg.Wait()
	close(errChan)

	// Report any errors
	for err := range errChan {
		logger.Get().Error(err)
	}

	return nil
}

func (s scanner) processTarget(rFt, rMt *regexp.Regexp) error {
	strategy := s.task.strategy
	content, err := strategy.GetContent()
	if err != nil {
		return fmt.Errorf("error getting file content: %v", err)
	}

	out := &output.OutputData{
		Target:   s.task.target,
		Filename: s.task.output,
		Results:  []output.Result{},
	}

	contentStr := *(*string)(unsafe.Pointer(&content))
	processedUrls := make(map[string]bool)

	// Apply each pattern to find URLs
	for _, pattern := range patterns {
		for _, match := range pattern.FindAllStringSubmatchIndex(contentStr, -1) {
			fullMatch := contentStr[match[0]:match[1]]

			// Extract URL from the full match - different patterns may have different group indices
			var url string
			if strings.HasPrefix(fullMatch, "`") && strings.HasSuffix(fullMatch, "`") {
				// Template literal pattern
				url = fullMatch
			} else if strings.Contains(fullMatch, "url") || strings.Contains(fullMatch, "endpoint") ||
				strings.Contains(fullMatch, "href") || strings.Contains(fullMatch, "src") {
				// Config object pattern
				parts := strings.Split(fullMatch, ":")
				if len(parts) > 1 {
					urlPart := strings.TrimSpace(parts[1])
					url = strings.Trim(urlPart, `'"`)
				} else {
					url = fullMatch
				}
			} else {
				// Other patterns - extract from quotes or parentheses
				url = strings.Trim(fullMatch, `'"()`)
			}

			// Clean URL if needed
			url = cleanUrl(url)

			if url == "" || len(url) < 4 {
				continue
			}

			// Apply exclusion rules
			if rFt.MatchString(url) || rMt.MatchString(url) {
				continue
			}

			// Skip if already processed
			if _, exists := processedUrls[url]; exists {
				continue
			}
			processedUrls[url] = true

			// Limit the context to avoid huge outputs
			startIdx := match[0] - 100
			if startIdx < 0 {
				startIdx = 0
			}

			endIdx := match[1] + 100
			if endIdx > len(content) {
				endIdx = len(content)
			}

			closeLines := content[startIdx:endIdx]
			out.Results = append(out.Results, output.Result{
				URL:      url,
				Location: string(closeLines),
			})

			logger.Get().Infof("found possible url: %s", url)
		}
	}

	logger.Get().Infof("%d possible url found", len(out.Results))
	oE, ok := outputEngines[s.getOutputEngineKey()]
	if !ok {
		return fmt.Errorf("output engine not found: %s", s.getOutputEngineKey())
	}

	err = oE.RenderAndSave(out)
	if err != nil {
		return fmt.Errorf("output failed: %v", err)
	}

	return nil
}

// cleanUrl removes common noise in URLs
func cleanUrl(url string) string {
	// Remove common JavaScript noise
	url = strings.TrimPrefix(url, "url(")
	url = strings.TrimPrefix(url, "URL(")
	url = strings.TrimSuffix(url, ")")
	url = strings.Trim(url, `'"`)

	// Clean template literals
	url = strings.Trim(url, "`")

	// Handle escaped characters
	url = strings.ReplaceAll(url, "\\\"", "\"")
	url = strings.ReplaceAll(url, "\\'", "'")

	return url
}

func (s scanner) getOutputEngineKey() string {
	return filepath.Ext(s.task.output)
}

func (s scanner) isAlreadyExists(url string, out *output.OutputData) bool {
	for _, a := range out.Results {
		if a.URL == url {
			return true
		}
	}
	return false
}

func defineStrategyForTarget(target string) strategies.ScanStrategy {
	if strings.Contains(target, "http://") || strings.Contains(target, "https://") {
		return strategies.URLStrategy{Target: target}
	}
	return strategies.FileStrategy{Target: target}
}
