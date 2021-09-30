package internal

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	e "github.com/xvrzhao/utils/errors"
)

func isArchiveValid(tarGzFile string) bool {
	cmd := exec.Command("tar", "-tzf", tarGzFile)
	cmdStdErrBuf := new(bytes.Buffer)
	cmd.Stderr = cmdStdErrBuf

	if err := cmd.Run(); err != nil {
		return false
	}

	return true
}

func GetInstalledGoVersionStrings() (versions []string, err error) {
	versions = make([]string, 0)

	fis, err := ioutil.ReadDir(gvmRoot)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadDir failed: %w", err)
	}

	for _, fi := range fis {
		if fi.IsDir() && fi.Name()[0] != '.' && fi.Name()[:2] == "go" {
			versions = append(versions, fi.Name()[2:])
		}
	}

	return
}

func SwitchVersion(v *Version) error {
	resetGlobalProgressBar(10, "Switching...")
	globalProgressBar.Add(2)
	defer func() {
		globalProgressBar.Finish()
		time.Sleep(time.Second)
		globalProgressBar.Clear()
	}()

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
		return fmt.Errorf("failed to GetCurrentVersionStr: %w", err)
	}

	if v.String() == currentVersion {
		return errors.New("can not remove current version")
	}

	if err = os.RemoveAll(v.GetInstallationDir()); err != nil {
		return fmt.Errorf("failed to remove the version directory: %w", err)
	}

	return nil
}
