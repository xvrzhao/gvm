package funcs

import (
	"bytes"
	"errors"
	"fmt"
	e "github.com/xvrzhao/utils/errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func mkGvmRoot() error {
	fileInfo, err := os.Stat(gvmRoot)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(gvmRoot, os.ModePerm); err != nil {
			return e.Wrapper(err, "gvmRoot mkdir error")
		}
		return nil
	}
	if err != nil {
		return e.Wrapper(err, "gvmRoot get fileInfo error")
	}
	if !fileInfo.IsDir() {
		return errors.New(fmt.Sprintf("%s exists but isn't a dir", gvmRoot))
	}
	return nil
}

func unCompress(tarGzFile, dstPath string) (err error) {
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
