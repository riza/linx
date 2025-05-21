package output

import (
	"encoding/json"
	"os"

	"github.com/riza/linx/pkg/logger"
)

type OutputJSON struct {
}

func (oj OutputJSON) RenderAndSave(data *OutputData) error {
	f, err := os.Create(data.Filename)
	if err != nil {
		return err
	}
	defer f.Close()

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	_, err = f.Write(jsonData)
	if err != nil {
		return err
	}

	logger.Get().Infof("results saved in JSON format: %s", data.Filename)
	return nil
}
