package main

import (
	"github.com/riza/linx/internal/banner"
	"github.com/riza/linx/internal/options"
	"github.com/riza/linx/linx"
	"github.com/riza/linx/pkg/logger"
)

const Version = "v0.0.1"

func main() {
	banner.Show(Version)

	opts, err := options.Get().Parse()
	if err != nil {
		logger.Get().Fatal(err)
	}

	err = linx.Run(opts)
	if err != nil {
		logger.Get().Error(err)
	}
}
