package s2i

// Interface is the interface for s2i.
type Interface interface {
	SourceToImage(code, jobName string) (string, error)
}
