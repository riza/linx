package strategies

type ScanStrategy interface {
	GetContent() ([]byte, error)
	GetFileName() string
}
