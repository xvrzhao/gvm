package funcs

import "testing"

func TestNewVersion(t *testing.T) {
	v := NewVersion(1, 14, 0, true)
	t.Log(v.filename, v.fullURL)
}

func TestVersion_Download(t *testing.T) {
	v := NewVersion(1, 14, 0, false)
	v.Download()
}
