package banner

import "github.com/riza/linx/pkg/logger"

const banner = `
   ___         
  / (_)__ __ __
 / / / _ \\ \ /
/_/_/_//_/_\_\  %s
`

func Show(ver string) {
	logger.Get().Printf(banner, ver)
}
