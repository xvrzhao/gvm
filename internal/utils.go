package internal

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func IsTarGzFileValid(tarGzFile string) bool {
	if err := exec.Command("tar", "-tzf", tarGzFile).Run(); err != nil {
		return false
	}

	return true
}

func DecompressUsingTar(tarGzFile, dstPath string) error {
	cmd := exec.Command("tar", "-C", dstPath, "-xzf", tarGzFile)
	cmdStdErrBuf := new(bytes.Buffer)
	cmd.Stderr = cmdStdErrBuf

	if err := cmd.Run(); err != nil {
		err = fmt.Errorf("command run error(%q), stderr output(%q)", err.Error(), cmdStdErrBuf.String())
		return fmt.Errorf("failed to run command: %w", err)
	}

	return nil
}

// GetAllInstalledVersions reads go dirs under gvmRoot to get
// all installed versions.
func GetAllInstalledVersions() (versions []string, err error) {
	versions = make([]string, 0)

	fis, err := ioutil.ReadDir(gvmRoot)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to ReadDir: %w", err)
	}

	for _, fi := range fis {
		if fi.IsDir() && fi.Name()[0] != '.' && fi.Name()[:2] == "go" {
			versions = append(versions, fi.Name()[2:])
		}
	}

	return
}

// GetCurrentVersion reads the dst of symbolic link goRoot to
// get the current version.
func GetCurrentVersion() (currentVersion string, err error) {
	noVersionErr := errors.New("no current version")

	gvmPath, err := os.Readlink(goRoot)
	if err != nil {
		return "", noVersionErr
	}

	if n, err := fmt.Sscanf(gvmPath, gvmRoot+"/go%s", &currentVersion); n != 1 || err != nil {
		return "", noVersionErr
	}

	return currentVersion, nil
}
