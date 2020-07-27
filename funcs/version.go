package funcs

import (
	"fmt"
	e "github.com/xvrzhao/utils/errors"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type Semantics struct {
	major, minor, patch uint8
}

func (s Semantics) String() string {
	v := fmt.Sprintf("%d.%d", s.major, s.minor)
	if s.patch != 0 {
		v = fmt.Sprintf("%s.%d", v, s.patch)
	}
	return v
}

type version struct {
	// env
	Semantics
	os       string
	arch     string
	filename string
	fullURL  string

	// download
	dling      bool
	finishDl   finishState
	dlFilename string

	// unCompress
	ucing     bool
	finishUc  finishState
	ucDirname string
}

func NewVersion(ver string, inCn bool) (v *version, err error) {
	sem, err := checkSemVer(ver)
	if err != nil {
		err = e.Wrapper(err, "checkSemVer error")
		return
	}

	v = &version{
		Semantics: sem,
		os:        runtime.GOOS,
		arch:      runtime.GOARCH,
		filename:  "",
		fullURL:   "",

		dling:      false,
		finishDl:   fsUnknown,
		dlFilename: "",

		ucing:     false,
		finishUc:  fsUnknown,
		ucDirname: "",
	}

	v.buildPackFilename()
	v.buildFullURL(inCn)

	if installed, ucDir, err := v.checkInstallation(); err != nil {
		err = e.Wrapper(err, "check is installed error")
		return
	} else if installed {
		v.finishUc, v.ucDirname = fsFinished, ucDir
	} else {
		v.finishUc = fsUnFinished
	}

	return
}

func (v *version) buildPackFilename() string {
	v.filename = fmt.Sprintf("go%v.%s-%s.tar.gz", v.Semantics, v.os, v.arch)
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
	fmt.Printf("> Downloading the go archive from %s ... ", v.fullURL)

	df, err := download(v)
	if err != nil {
		fmt.Println()
		e.Log(e.Wrapper(err, "version download error"))
		os.Exit(1)
	}

	v.dling, v.finishDl, v.dlFilename = false, fsUnFinished, df
	fmt.Println("Done.")
}

func (v *version) UnCompress() {
	v.ucing = true
	fmt.Printf("> Extract to %s from downloaded archive ... ", gvmRoot)

	goDir := filepath.Join(gvmRoot, "go")
	vgoDir := filepath.Join(gvmRoot, fmt.Sprintf("go%v", v.Semantics))

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

	v.ucing, v.finishUc, v.ucDirname = false, fsUnFinished, vgoDir
	fmt.Println("Done.")
}

func (v *version) checkInstallation() (installed bool, ucDir string, err error) {
	vStr := "go" + v.Semantics.String()
	versions, err := GetInstalledGoVersions()
	if err != nil {
		err = e.Wrapper(err, "GetInstalledGoVersions error")
		return
	}
	for _, ver := range versions {
		if ver == vStr {
			installed, ucDir = true, filepath.Join(gvmRoot, vStr)
			return
		}
	}
	installed = false
	return
}

func (v *version) IsInstalled() bool {
	if v.finishUc == fsUnFinished {
		return true
	}
	return false
}
