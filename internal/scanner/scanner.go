package scanner

import (
	"fmt"
	"github.com/riza/linx/internal/options"
	"github.com/riza/linx/internal/output"
	"github.com/riza/linx/internal/scanner/strategies"
	"github.com/riza/linx/pkg/logger"
	"path/filepath"
	"regexp"
	"strings"
	"unsafe"
)

// rule from LinkFinder https://github.com/GerbenJavado/LinkFinder/blob/master/linkfinder.py#L29 ty @GerbenJavado
const (
	rule                = `(?:"|')(((?:[a-zA-Z]{1,10}://|//)[^"'/]{1,}\.[a-zA-Z]{2,}[^"']{0,})|((?:/|\.\./|\./)[^"'><,;| *()(%%$^/\\\[\]][^"'><,;|()]{1,})|([a-zA-Z0-9_\-/]{1,}/[a-zA-Z0-9_\-/]{1,}\.(?:[a-zA-Z]{1,4}|action)(?:[\?|#][^"|']{0,}|))|([a-zA-Z0-9_\-/]{1,}/[a-zA-Z0-9_\-/]{3,}(?:[\?|#][^"|']{0,}|))|([a-zA-Z0-9_\-]{1,}\.(?:php|asp|aspx|jsp|json|action|html|js|txt|xml)(?:[\?|#][^"|']{0,}|)))(?:"|')`
	excludeFileTypeRule = `.css|.jpg|.jpeg|.png|.svg|.img|.gif|.mp4|.flv|.ogv|.webm|.webp|.mov|.mp3|.m4a|.m4p|.scss|.tif|.tiff|.ttf|.otf|.woff|.woff2|.bmp|.ico|.eot|.htc|.rtf|.swf|.image|w3.org|doubleclick.net|youtube.com|.vue|jquery|bootstrap|font|jsdelivr.net|vimeo.com|pinterest.com|facebook|linkedin|twitter|instagram|google|mozilla.org|jibe.com|schema.org|schemas.microsoft.com|wordpress.org|w.org|wix.com|parastorage.com|whatwg.org|polyfill.io|typekit.net|schemas.openxmlformats.org|openweathermap.org|openoffice.org|reactjs.org|angularjs.org|java.com|purl.org|/image|/img|/css|/wp-json|/wp-content|/wp-includes|/theme|/audio|/captcha|/font|robots.txt|node_modules|.wav|.gltf|.js`
	excludeMimeTypeRule = `text/css|image/jpeg|image/jpg|image/png|image/svg+xml|image/gif|image/tiff|image/webp|image/bmp|image/x-icon|image/vnd.microsoft.icon|font/ttf|font/woff|font/woff2|font/x-woff2|font/x-woff|font/otf|audio/mpeg|audio/wav|audio/webm|audio/aac|audio/ogg|audio/wav|audio/webm|video/mp4|video/mpeg|video/webm|video/ogg|video/mp2t|video/webm|video/x-msvideo|application/font-woff|application/font-woff2|application/vnd.android.package-archive|binary/octet-stream|application/octet-stream|application/pdf|application/x-font-ttf|application/x-font-otf|application/json|text/javascript|text/plain|text/x-yaml|text/html|text/babel|text/markdown|text/tsx|application/typescript|application/javascript|text/x-handlebars-template|application/x-typescript|text/x-gfm|text/jsx`
)

var (
	outputEngines = map[string]output.Output{
		"":      output.OutputNoop{},
		".html": output.OutputHTML{},
	}
)

type task struct {
	target   string
	output   string
	strategy strategies.ScanStrategy
}

type scanner struct {
	task task
}

func NewScanner(opts *options.Options) scanner {
	return scanner{
		task{
			target:   opts.Target,
			output:   opts.Output,
			strategy: defineStrategyForTarget(opts.Target),
		},
	}
}

func (s scanner) Run() error {
	r, _ := regexp.Compile(rule)
	rFt, _ := regexp.Compile(excludeFileTypeRule)
	rMt, _ := regexp.Compile(excludeMimeTypeRule)

	strategy := s.task.strategy
	content, err := strategy.GetContent()
	if err != nil {
		return fmt.Errorf(errGetFileContent, err)
	}

	out := &output.OutputData{
		Target:   s.task.target,
		Filename: s.task.output,
		Results:  []output.Result{},
	}

	for _, f := range r.FindAllStringSubmatchIndex(*(*string)(unsafe.Pointer(&content)), -1) {
		url := content[f[0]:f[1]]
		if rFt.MatchString(string(url)) || rMt.MatchString(string(url)) {
			continue
		}

		if s.isAlreadyExists(string(url), out) {
			continue
		}

		closeLines := content[f[0]-100 : f[1]+100]
		out.Results = append(out.Results, output.Result{
			URL:      string(url),
			Location: string(closeLines),
		})

		logger.Get().Infof("found possible url: %s", string(url))
	}

	logger.Get().Infof("%d possible url found", len(out.Results))
	oE, ok := outputEngines[s.getOutputEngineKey()]
	if !ok {
		return fmt.Errorf(errOutputEngineNotFound, s.getOutputEngineKey())
	}

	err = oE.RenderAndSave(out)
	if err != nil {
		return fmt.Errorf(errOutputFailed, err)
	}

	return nil
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
		return strategies.URLStrategy{target}
	}
	return strategies.FileStrategy{target}
}
