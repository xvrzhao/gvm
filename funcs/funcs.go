package funcs

import (
	"bytes"
	"errors"
	"fmt"
	e "github.com/xvrzhao/utils/errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func unCompress(tarGzFile, dstPath string) (err error) {
	oldUmask := syscall.Umask(0)
	defer syscall.Umask(oldUmask)

	if err = os.MkdirAll(dstPath, os.ModePerm); err != nil {
		err = e.Wrapper(err, "dstPath mkdir error")
		return
	}

	cmd := exec.Command("tar", "-C", dstPath, "-xzf", tarGzFile)
	cmdStdErrBuf := new(bytes.Buffer)
	cmd.Stderr = cmdStdErrBuf

	if err = cmd.Run(); err != nil {
		return e.Wrapper(fmt.Errorf("Cmd.Run error: %w, stderr output: %q", err, cmdStdErrBuf.String()),
			"command run error")
	}

	return nil
}

func download(v *version) (dlFilename string, err error) {
	res, err := http.Get(v.fullURL)
	if err != nil {
		err = e.Wrapper(err, "http get go package file error")
		return
	}
	defer res.Body.Close()

	err = os.MkdirAll(tmpPath, os.ModePerm)
	if err != nil {
		err = e.Wrapper(err, "tmp path mkdir error")
		return
	}

	df := filepath.Join(tmpPath, v.filename)
	file, err := os.Create(df)
	if err != nil {
		err = e.Wrapper(err, "create local tmp file error")
		return
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		err = e.Wrapper(err, "downloading error")
		return
	}

	dlFilename = df
	return
}

var semVerError = errors.New("invalid semantic version")

func checkSemVer(semVer string) (v Semantics, err error) {
	s := strings.Split(semVer, ".")
	if len(s) < 2 || len(s) > 3 {
		err = semVerError
		return
	}
	for idx, semverItem := range s {
		var num int
		num, err = strconv.Atoi(semverItem)
		if err != nil {
			err = semVerError
			return
		}
		switch idx {
		case 0:
			v.major = uint8(num)
		case 1:
			v.minor = uint8(num)
		case 2:
			v.patch = uint8(num)
		}
	}
	return
}

func GetInstalledGoVersions() (versions []string, err error) {
	versions = make([]string, 0)

	fis, err := ioutil.ReadDir(gvmRoot)
	if os.IsNotExist(err) {
		err = nil
		return
	}
	if err != nil {
		err = e.Wrapper(err, "gvmRoot readdir error")
		return
	}
	for _, fi := range fis {
		if fi.IsDir() && fi.Name()[0] != '.' && fi.Name()[:2] == "go" {
			versions = append(versions, fi.Name()[2:])
		}
	}
	return
}
