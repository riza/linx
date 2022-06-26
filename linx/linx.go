package linx

import (
	"github.com/riza/linx/internal/options"
	"github.com/riza/linx/internal/scanner"
)

func Run(options *options.Options) error {
	scanner := scanner.NewScanner(options.Target)
	err := scanner.Run()
	if err != nil {
		return err
	}
	return nil
}
