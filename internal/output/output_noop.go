package output

type OutputNoop struct {
}

func (o OutputNoop) RenderAndSave(data *OutputData) error {
	return nil
}
