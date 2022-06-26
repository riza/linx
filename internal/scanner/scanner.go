package scanner

import (
	"fmt"
	"github.com/riza/linx/internal/scanner/strategies"
	"github.com/riza/linx/pkg/logger"
	"regexp"
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
			strategy: defineStrategyForTarget(target),
		},
	}
}

func (s scanner) Run() error {
	r, _ := regexp.Compile(rule)

	content, err := s.task.strategy.GetContent()
	if err != nil {
		return fmt.Errorf(errGetFileContent, err)
	}

	urls := r.FindAllString(*(*string)(unsafe.Pointer(&content)), -1)
	for _, s := range urls {
		logger.Get().Infof("found possible url: %s", s)
	}

	logger.Get().Infof("%d possible url found", len(urls))
	return nil
}

func defineStrategyForTarget(target string) strategies.ScanStrategy {
	if strings.Contains(target, "http://") || strings.Contains(target, "https://") {
		return strategies.URLStrategy{target}
	}
	return strategies.FileStrategy{target}
}
