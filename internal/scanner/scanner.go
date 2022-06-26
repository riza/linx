package scanner

import (
	"fmt"
	"github.com/riza/linx/internal/output"
	"github.com/riza/linx/internal/scanner/strategies"
	"github.com/riza/linx/pkg/logger"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"unsafe"
)

// rule from LinkFinder https://github.com/GerbenJavado/LinkFinder/blob/master/linkfinder.py#L29 ty @GerbenJavado
const rule = `(?:"|')(((?:[a-zA-Z]{1,10}://|//)[^"'/]{1,}\.[a-zA-Z]{2,}[^"']{0,})|((?:/|\.\./|\./)[^"'><,;| *()(%%$^/\\\[\]][^"'><,;|()]{1,})|([a-zA-Z0-9_\-/]{1,}/[a-zA-Z0-9_\-/]{1,}\.(?:[a-zA-Z]{1,4}|action)(?:[\?|#][^"|']{0,}|))|([a-zA-Z0-9_\-/]{1,}/[a-zA-Z0-9_\-/]{3,}(?:[\?|#][^"|']{0,}|))|([a-zA-Z0-9_\-]{1,}\.(?:php|asp|aspx|jsp|json|action|html|js|txt|xml)(?:[\?|#][^"|']{0,}|)))(?:"|')`

type task struct {
	target   string
	strategy strategies.ScanStrategy
}

type scanner struct {
	task task
}

func NewScanner(target string) scanner {
	return scanner{
		task{
			target:   target,
			strategy: defineStrategyForTarget(target),
		},
	}
}

func (s scanner) Run() error {
	r, _ := regexp.Compile(rule)

	strategy := s.task.strategy
	content, err := strategy.GetContent()
	if err != nil {
		return fmt.Errorf(errGetFileContent, err)
	}

	out := output.OutputData{
		Target:   s.task.target,
		Filename: strategy.GetFileName() + "_result.html",
		Results:  []output.Result{},
	}

	for _, s := range r.FindAllStringSubmatchIndex(*(*string)(unsafe.Pointer(&content)), -1) {
		url := content[s[0]:s[1]]
		closeLines := content[s[0]-100 : s[1]+100]

		out.Results = append(out.Results, output.Result{
			URL:      string(url),
			Location: string(closeLines),
		})

		logger.Get().Infof("found possible url: %s", string(url))
	}

	logger.Get().Infof("%d possible url found", len(out.Results))

	err = output.NewOutputHTML(out).RenderAndSave()
	if err != nil {
		return fmt.Errorf(errOutputFailed, err)
	}

	err = openResults(out.Filename)
	if err != nil {
		return fmt.Errorf(errOutputOpenFailed, err)
	}

	return nil
}

func defineStrategyForTarget(target string) strategies.ScanStrategy {
	if strings.Contains(target, "http://") || strings.Contains(target, "https://") {
		return strategies.URLStrategy{target}
	}
	return strategies.FileStrategy{target}
}

func openResults(fileName string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", fileName).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", fileName).Start()
	case "darwin":
		err = exec.Command("open", fileName).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		return err
	}

	return nil
}
