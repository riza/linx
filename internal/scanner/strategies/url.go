package strategies

import (
	"fmt"
	"io/ioutil"
	"linx/pkg/logger"
	"net/http"
)

type URLStrategy struct {
	Target string
}

func (us URLStrategy) GetContent() ([]byte, error) {
	logger.Get().Debugf("selected url strategy target=%s", us.Target)
	return us.getFileContent()
}

func (us URLStrategy) getFileContent() ([]byte, error) {
	logger.Get().Debugf("getting file from %s", us.Target)
	resp, err := http.Get(us.Target)
	if err != nil {
		return nil, err
	}

	logger.Get().Debugf("response: status code=%d", resp.StatusCode)
	if !(resp.StatusCode >= 200 && resp.StatusCode <= 299) {
		return nil, fmt.Errorf("getting url content fail. status code is not success code=%d", resp.StatusCode)
	}

	logger.Get().Debugf("response: content length=%d", resp.ContentLength)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
