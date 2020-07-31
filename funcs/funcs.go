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

func decompress(tarGzFile, dstPath string) error {
	oldUmask := syscall.Umask(0)
	defer syscall.Umask(oldUmask)

	if err := os.MkdirAll(dstPath, os.ModePerm); err != nil {
		return e.Wrapper(err, "mkdir %s error", dstPath)
	}

	cmd := exec.Command("tar", "-C", dstPath, "-xzf", tarGzFile)
	cmdStdErrBuf := new(bytes.Buffer)
	cmd.Stderr = cmdStdErrBuf

	if err := cmd.Run(); err != nil {
		return e.Wrapper(fmt.Errorf("command run error: %w, stderr output: %q", err, cmdStdErrBuf.String()),
			"command run error")
	}

	return nil
}

func download(v *version) (downloadedTarGzFile string, err error) {
	res, err := http.Get(v.downloadURL)
	if err != nil {
		err = e.Wrapper(err, "error when HTTP GET %s", v.downloadURL)
		return
	}
	defer res.Body.Close()

	oldUmask := syscall.Umask(0)
	defer syscall.Umask(oldUmask)

	if err = os.MkdirAll(tmpPath, os.ModePerm); err != nil {
		err = e.Wrapper(err, "error when mkdir %s", tmpPath)
		return
	}

	dstFile := filepath.Join(tmpPath, v.tarGzFile)
	file, err := os.Create(dstFile)
	if err != nil {
		err = e.Wrapper(err, "error when creating dstFile")
		return
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		err = e.Wrapper(err, "error when copying from res.Body to file")
		return
	}

	downloadedTarGzFile = dstFile
	return
}

var semVerError = errors.New("invalid semantic version")

func checkSemver(semVer string) (sem semantics, err error) {
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
			sem.major = uint8(num)
		case 1:
			sem.minor = uint8(num)
		case 2:
			sem.patch = uint8(num)
		}
	}

	return
}

func GetInstalledGoVersionStrings() (versions []string, err error) {
	versions = make([]string, 0)

	fis, err := ioutil.ReadDir(gvmRoot)
	if os.IsNotExist(err) {
		err = nil
		return
	}
	if err != nil {
		err = e.Wrapper(err, "readDir of gvmRoot error")
		return
	}

	for _, fi := range fis {
		if fi.IsDir() && fi.Name()[0] != '.' && fi.Name()[:2] == "go" {
			versions = append(versions, fi.Name()[2:])
		}
	}

	return
}
