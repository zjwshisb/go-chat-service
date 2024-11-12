package file

type localAdapter struct {
}

func (s *localAdapter) Url(path string) string {
	return path
}
