package strategies

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/riza/linx/pkg/logger"
)

type URLStrategy struct {
	Target string
}

func (us URLStrategy) GetContent() ([]byte, error) {
	logger.Get().Debugf("selected url strategy target=%s", us.Target)
	return us.getFileContent()
}

func (us URLStrategy) GetFileName() string {
	_, file := path.Split(us.Target)
	return file
}

func (us URLStrategy) getFileContent() ([]byte, error) {
	logger.Get().Debugf("getting file from %s", us.Target)

	client := &http.Client{}
	req, err := http.NewRequest("GET", us.Target, nil)
	if err != nil {
		return nil, err
	}

	// Set default headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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
