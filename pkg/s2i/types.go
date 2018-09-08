package s2i

type Interface interface {
	SourceToImage(code, jobName string) (string, error)
}
