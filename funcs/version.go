package funcs

import (
	"fmt"
	e "github.com/xvrzhao/utils/errors"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type semantics struct {
	major, minor, patch uint8
}

func (s semantics) String() string {
	v := fmt.Sprintf("%d.%d", s.major, s.minor)
	if s.patch != 0 {
		v = fmt.Sprintf("%s.%d", v, s.patch)
	}
	return v
}

type version struct {
	// env
	semantics
	os       string
	arch     string
	filename string
	fullURL  string

	// download
	dling      bool
	finishDl   bool
	dlFilename string

	// unCompress
	ucing    bool
	finishUc bool
}

func NewVersion(major, minor, patch uint8, inCn bool) *version {
	v := &version{
		semantics: semantics{major, minor, patch},
		os:        runtime.GOOS,
		arch:      runtime.GOARCH,
		filename:  "",
		fullURL:   "",
	}
	v.buildPackFilename()
	v.buildFullURL(inCn)
	return v
}

func (v *version) buildPackFilename() string {
	v.filename = fmt.Sprintf("go%v.%s-%s.tar.gz", v.semantics, v.os, v.arch)
	return v.filename
}

func (v *version) buildFullURL(inCn bool) string {
	if v.filename == "" {
		log.Fatalf("failed to build fullURL of version, cause filename is not built yet")
	}
	pf := dlPrefix
	if inCn {
		pf = dlPrefixCn
	}
	v.fullURL = pf + v.filename
	return v.fullURL
}

func (v *version) Download() {
	v.dling = true
	fmt.Printf("> downloading the go archive from %s ... ", v.fullURL)

	df, err := download(v)
	if err != nil {
		fmt.Println()
		e.Log(e.Wrapper(err, "version download error"))
		os.Exit(1)
	}

	v.dling, v.finishDl, v.dlFilename = false, true, df
	fmt.Println("done")
}

func (v *version) UnCompress() {
	v.ucing = true
	fmt.Printf("> extract to %s from downloaded archive ... ", gvmRoot)

	goDir := filepath.Join(gvmRoot, "go")
	vgoDir := filepath.Join(gvmRoot, fmt.Sprintf("go%v", v.semantics))

	if err := os.RemoveAll(goDir); err != nil {
		fmt.Println()
		e.Log(e.Wrapper(err, "remove %s error", goDir))
		os.Exit(1)
	}

	if err := os.RemoveAll(vgoDir); err != nil {
		fmt.Println()
		e.Log(e.Wrapper(err, "remove %s error", vgoDir))
		os.Exit(1)
	}

	if err := unCompress(v.dlFilename, gvmRoot); err != nil {
		fmt.Println()
		e.Log(e.Wrapper(err, "unCompress error"))
		os.Exit(1)
	}

	if err := os.Rename(goDir, vgoDir); err != nil {
		fmt.Println()
		e.Log(e.Wrapper(err, "rename go to gox.x.x error"))
		os.Exit(1)
	}

	v.ucing, v.finishUc = false, true
	fmt.Println("done")
}
