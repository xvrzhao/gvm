package funcs

import (
	"testing"
)

func TestGetInstalledGoVersionStrings(t *testing.T) {
	versions, err := GetInstalledGoVersionStrings()
	if err != nil {
		t.Errorf("failed to GetInstalledGoVersionStrings: %v", err)
		return
	}

	t.Log(versions)
}

func TestIsArchiveValid(t *testing.T) {
	isArchiveValid("/tmp/gvm/go1.14.6.darwin-amd64.tar.gz")
}
