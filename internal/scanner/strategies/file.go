package strategies

import (
	"io/ioutil"
	"linx/pkg/logger"
	"os"
)

type FileStrategy struct {
	Target string
}

func (fs FileStrategy) GetContent() ([]byte, error) {
	logger.Get().Debugf("selected file content strategy target=%s", fs.Target)
	return fs.readFileContent()
}

func (fs FileStrategy) readFileContent() ([]byte, error) {
	_, err := os.Stat(fs.Target)
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadFile(fs.Target)
	if err != nil {
		return nil, err
	}

	return content, nil
}
