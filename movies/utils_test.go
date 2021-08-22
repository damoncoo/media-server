package movies

import "testing"

func TestFileName(t *testing.T) {
	path := "/path/Abc.mp4"
	name := FileName(path)
	t.Log(name)
}
