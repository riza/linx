package main

import (
	"linx/internal/banner"
	"linx/internal/options"
	"linx/linx"
	"linx/pkg/logger"
)

const Version = "v1.0"

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
