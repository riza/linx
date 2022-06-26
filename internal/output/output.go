package output

type Output interface {
	RenderAndSave(data *OutputData) error
}

type OutputData struct {
	Target   string
	Filename string
	Results  []Result
}

type Result struct {
	URL      string
	Location string
}
