package output

type OutputData struct {
	Target   string
	Filename string
	Results  []Result
}

type Result struct {
	URL      string
	Location string
}
