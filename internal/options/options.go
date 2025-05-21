package options

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/riza/linx/pkg/logger"
)

type Options struct {
	Target   string
	Output   string
	Debug    bool
	Parallel bool
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

	flag.BoolVar(&o.Debug, "debug", false, "do you want to know what's inside the engine?")
	flag.StringVar(&o.Output, "output", "", "output file name (supports html and json formats)")
	flag.BoolVar(&o.Parallel, "parallel", false, "scan multiple targets in parallel (only works with comma separated targets)")

	// Parse flags, but the first non-flag argument will be our target
	flag.Parse()

	if o.Debug {
		logger.Get().SetLevelDebug()
	}

	// Get positional arguments
	args := flag.Args()
	if len(args) == 0 {
		printDefaults()
		return nil, fmt.Errorf("target is required")
	}

	// First positional argument is the target
	o.Target = strings.Join(args, ",")

	isValid := validateTarget(o.Target)
	if !isValid {
		printDefaults()
		return nil, fmt.Errorf("target is invalid: %s", o.Target)
	}

	return o, nil
}

func validateTarget(target string) (isValid bool) {
	// Support for multiple targets (comma separated)
	targets := strings.Split(target, ",")

	for _, t := range targets {
		t = strings.TrimSpace(t)
		isCurrentValid := false

		// Check if it's a URL to a JS file
		if (strings.Contains(t, "http://") || strings.Contains(t, "https://")) && strings.Contains(t, ".js") {
			isCurrentValid = true
		}

		// Check if it's a local JS file
		if strings.Contains(t, ".js") {
			isCurrentValid = true
		}

		if !isCurrentValid {
			return false
		}
	}

	return true
}

func printDefaults() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options] <url or path>\n\nOptions:\n", os.Args[0])
	flag.PrintDefaults()
}
