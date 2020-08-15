/*
Copyright © 2020 Xavier Zhao <xvrzhao@gmail.com>
Licensed under the MIT License. See LICENSE file in the project root for license information.
*/

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
	"time"
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

func isArchiveValid(tarGzFile string) bool {
	cmd := exec.Command("tar", "-tzf", tarGzFile)
	cmdStdErrBuf := new(bytes.Buffer)
	cmd.Stderr = cmdStdErrBuf

	if err := cmd.Run(); err != nil {
		return false
	}

	return true
}

func download(v *Version) (downloadedTarGzFile string, err error) {
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

func SwitchVersion(v *Version) error {
	if err := v.Reload(); err != nil {
		return e.Wrapper(err, "version reload error")
	}

	if !v.IsInstalled() {
		return errors.New("the version is not installed")
	}

	fi, err := os.Lstat(goRoot)
	if !os.IsNotExist(err) {
		if err != nil {
			return e.Wrapper(err, "Lstat goRoot error")
		}

		if fi.Mode()&os.ModeSymlink != 0 {
			if err := os.Remove(goRoot); err != nil {
				return e.Wrapper(err, "remove goRoot error")
			}
		} else {
			if err := backupOldGoRoot(); err != nil {
				return e.Wrapper(err, "backupOldGoRoot error")
			}
		}
	}

	if err := os.Symlink(v.GetInstallationDir(), goRoot); err != nil {
		return e.Wrapper(err, "error when creating symlink")
	}

	return nil
}

func backupOldGoRoot() error {
	cmd := exec.Command("mv", goRoot, fmt.Sprintf("%s.old.%s", goRoot,
		time.Now().Format("20060102150405")))
	cmdStdErr := new(bytes.Buffer)
	cmd.Stderr = cmdStdErr

	if err := cmd.Run(); err != nil {
		return e.Wrapper(fmt.Errorf("command run error: %w, stderr output: %q", err, cmdStdErr.String()),
			"command run error")
	}

	return nil
}

func GetCurrentVersionStr() (currentVersionStr string, err error) {
	noVersionErr := errors.New("no current version")
	gvmPath, err := os.Readlink(goRoot)
	if err != nil {
		err = noVersionErr
		return
	}

	if n, er := fmt.Sscanf(gvmPath, gvmRoot+"/go%s", &currentVersionStr); n != 1 || er != nil {
		err = noVersionErr
		return
	}

	return
}

func RmVersion(v *Version) error {
	currentVersion, err := GetCurrentVersionStr()
	if err != nil {
		return e.Wrapper(err, "GetCurrentVersionStr error")
	}

	if v.String() == currentVersion {
		return errors.New("can not remove current version")
	}

	if err = os.RemoveAll(v.GetInstallationDir()); err != nil {
		return e.Wrapper(err, "error when removing version dir")
	}

	return nil
}
