package options

import (
	"flag"
	"fmt"
	"linx/pkg/logger"
	"os"
	"strings"
)

type Options struct {
	Target string
	Debug  bool
}

var (
	o *Options
)

func init() {
	o = &Options{}
}

func Get() *Options {
	return o
}

func (o *Options) Parse() (*Options, error) {
	flag.StringVar(&o.Target, "target", "", "can be *.js file path or url")
	flag.BoolVar(&o.Debug, "debug", false, "do you want to know what's inside the engine?")
	flag.Parse()

	if o.Debug {
		logger.Get().SetLevelDebug()
	}

	if o.Target == "" {
		printDefaults()
		return nil, fmt.Errorf(errTargetIsRequired, o.Target)
	}

	isValid := validateTarget(o.Target)
	if !isValid {
		printDefaults()
		return nil, fmt.Errorf(errTargetIsInvalid, o.Target)
	}

	return o, nil
}

func validateTarget(target string) (isValid bool) {
	if (strings.Contains(target, "http://") || strings.Contains(target, "https://")) && strings.Contains(target, ".js") {
		return true
	}
	if strings.Contains(target, ".js") {
		return true
	}
	return false
}

func printDefaults() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}
