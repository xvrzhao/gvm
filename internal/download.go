package internal

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"syscall"
)

func checkDownload(tarGzFileName string) (filePath string, isDownloaded bool, err error) {
	filePath = filepath.Join(tmpPath, tarGzFileName)
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		return filePath, false, nil
	}
	if err != nil {
		return filePath, false, fmt.Errorf("os.State failed: %w", err)
	}

	if !isArchiveValid(filePath) {
		return filePath, false, nil
	}

	return filePath, true, nil
}

func download(v *Version) (downloadedTarGzFile string, err error) {
	res, err := http.Get(v.downloadURL)
	if err != nil {
		err = fmt.Errorf("failed to GET %s: %w", v.downloadURL, err)
		return
	}
	defer res.Body.Close()

	oldUmask := syscall.Umask(0)
	defer syscall.Umask(oldUmask)

	if err = os.MkdirAll(tmpPath, os.ModePerm); err != nil {
		err = fmt.Errorf("failed to mkdir %s: %w", tmpPath, err)
		return
	}

	dstFile := filepath.Join(tmpPath, v.tarGzFile)
	file, err := os.Create(dstFile)
	if err != nil {
		err = fmt.Errorf("failed to create dstFile(%s): %w", dstFile, err)
		return
	}
	defer file.Close()

	resetGlobalProgressBar(res.ContentLength, "Downloading...")
	_, err = io.Copy(io.MultiWriter(file, globalProgressBar), res.Body)
	if err != nil {
		err = fmt.Errorf("failed to copy from res.Body to file: %w", err)
		return
	}

	downloadedTarGzFile = dstFile
	return
}
