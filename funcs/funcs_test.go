package funcs

import "testing"

func TestUnCompress(t *testing.T) {
	if err := unCompress("/tmp/gvm/go1.14.6.darwin-amd64.tar.gz", "/tmp"); err != nil {
		t.Errorf("failed to uncompress: %v", err)
		return
	}
}
