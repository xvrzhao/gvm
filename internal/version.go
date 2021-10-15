package internal

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"time"
)

type Version struct {
	// env
	Semantics
	os          string
	arch        string
	tarGzFile   string
	downloadURL string

	// download
	isDownloaded        bool
	downloadedTarGzFile string

	// decompress
	isDecompressed bool
	dir            string
}

func NewVersion(versionName string, inCn bool) (*Version, error) {
	sem, err := NewSemantics(versionName)
	if err != nil {
		return nil, fmt.Errorf("failed to NewSemantics: %w", err)
	}

	v := &Version{
		Semantics:   sem,
		os:          runtime.GOOS,
		arch:        runtime.GOARCH,
		tarGzFile:   "",
		downloadURL: "",

		isDownloaded:        false,
		downloadedTarGzFile: "",

		isDecompressed: false,
		dir:            "",
	}

	v.buildTarGzFile()
	if _, err = v.buildDownloadURL(inCn); err != nil {
		return nil, fmt.Errorf("failed to buildDownloadURL: %w", err)
	}

	if err = v.Reload(); err != nil {
		return nil, fmt.Errorf("failed to Reload: %w", err)
	}

	return v, nil
}

func (v *Version) Reload() error {
	if filePath, isDownloaded, err := v.checkDownload(); err != nil {
		return fmt.Errorf("failed to checkDownload: %w", err)
	} else if isDownloaded {
		v.isDownloaded, v.downloadedTarGzFile = true, filePath
	} else {
		v.isDownloaded, v.downloadedTarGzFile = false, ""
	}

	if dir, isInstalled, err := v.checkInstallation(); err != nil {
		return fmt.Errorf("failed to checkInstallation: %w", err)
	} else if isInstalled {
		v.isDecompressed, v.dir = true, dir
	} else {
		v.isDecompressed, v.dir = false, ""
	}

	return nil
}

func (v *Version) Install(force bool) error {
	if err := v.Download(force); err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}

	if err := v.Decompress(force); err != nil {
		return fmt.Errorf("failed to decompress: %w", err)
	}

	return nil
}

func (v *Version) Download(force bool) error {
	if v.isDownloaded == true && !force {
		return nil
	}

	file, err := v.download()
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}

	v.isDownloaded, v.downloadedTarGzFile = true, file
	return nil
}

func (v *Version) download() (downloadedTarGzFile string, err error) {
	res, err := http.Get(v.downloadURL)
	if err != nil {
		return "", fmt.Errorf("failed to GET %s: %w", v.downloadURL, err)
	}
	defer res.Body.Close()

	oldUmask := syscall.Umask(0)
	defer syscall.Umask(oldUmask)
	if err = os.MkdirAll(tmpPath, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to mkdir %s: %w", tmpPath, err)
	}

	dstFile := filepath.Join(tmpPath, v.tarGzFile)
	file, err := os.Create(dstFile)
	if err != nil {
		return "", fmt.Errorf("failed to create dstFile(%s): %w", dstFile, err)
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to copy from res.Body to file: %w", err)
	}

	return dstFile, nil
}

func (v *Version) Decompress(force bool) error {
	if v.isDecompressed == true && !force {
		return nil
	}

	if v.isDownloaded != true {
		return errors.New("version is not downloaded")
	}

	dir, err := v.decompress()
	if err != nil {
		return fmt.Errorf("failed to decompress: %w", err)
	}

	v.isDecompressed, v.dir = true, dir
	return nil
}

func (v *Version) decompress() (dir string, err error) {
	goDir := filepath.Join(gvmRoot, "go")
	vgoDir := filepath.Join(gvmRoot, fmt.Sprintf("go%v", v.Semantics))

	if err = os.RemoveAll(goDir); err != nil {
		return "", fmt.Errorf("failed to remove directory(%s): %w", goDir, err)
	}
	if err = os.RemoveAll(vgoDir); err != nil {
		return "", fmt.Errorf("failed to remove directory(%s): %w", vgoDir, err)
	}

	oldUmask := syscall.Umask(0)
	defer syscall.Umask(oldUmask)
	if err = os.MkdirAll(gvmRoot, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to mkdir %q: %w", gvmRoot, err)
	}

	if err = DecompressUsingTar(tmpPath+"/"+v.tarGzFile, gvmRoot); err != nil {
		return "", fmt.Errorf("failed to decompressUsingTar: %w", err)
	}

	if err = os.Rename(goDir, vgoDir); err != nil {
		return "", fmt.Errorf("failed to rename %s to %s: %w", goDir, vgoDir, err)
	}

	return vgoDir, nil
}

func (v *Version) Switch() error {
	if err := v.Reload(); err != nil {
		return fmt.Errorf("failed to Reload: %w", err)
	}

	if !v.IsInstalled() {
		return ErrVersionNotInstalled
	}

	fi, err := os.Lstat(goRoot)
	if !os.IsNotExist(err) {
		if err != nil {
			return fmt.Errorf("failed to Lstat goRoot: %w", err)
		}

		if fi.Mode()&os.ModeSymlink != 0 {
			if err := os.Remove(goRoot); err != nil {
				return fmt.Errorf("failed to remove goRoot: %w", err)
			}
		} else {
			if err := os.Rename(goRoot, fmt.Sprintf("%s.backup.%s", goRoot, time.Now().Format("20060102150405"))); err != nil {
				return fmt.Errorf("failed to backupOldGoRoot: %w", err)
			}
		}
	}

	if err := os.Symlink(v.GetInstallationDir(), goRoot); err != nil {
		return fmt.Errorf("failed to Symlink: %w", err)
	}

	return nil
}

func (v *Version) Remove() error {
	currentVersion, err := GetCurrentVersion()
	if err != nil {
		return fmt.Errorf("failed to GetCurrentVersion: %w", err)
	}

	if v.String() == currentVersion {
		return ErrVersionIsInUse
	}

	if err = os.RemoveAll(v.GetInstallationDir()); err != nil {
		return fmt.Errorf("failed to remove the version directory: %w", err)
	}

	return nil
}

func (v *Version) checkDownload() (filePath string, isDownloaded bool, err error) {
	if v.tarGzFile == "" {
		return "", false, errors.New("version.tarGzFile not built")
	}

	filePath = filepath.Join(tmpPath, v.tarGzFile)
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		return filePath, false, nil
	}
	if err != nil {
		return filePath, false, fmt.Errorf("failed to os.State: %w", err)
	}

	if !IsTarGzFileValid(filePath) {
		return filePath, false, nil
	}

	return filePath, true, nil
}

func (v *Version) checkInstallation() (versionDir string, isInstalled bool, err error) {
	thisVersionStr := fmt.Sprintf("go%v", v.Semantics)
	versions, err := GetAllInstalledVersions()
	if err != nil {
		return "", false, fmt.Errorf("failed to GetAllInstalledVersions: %w", err)
	}

	for _, version := range versions {
		if "go"+version == thisVersionStr {
			return filepath.Join(gvmRoot, thisVersionStr), true, nil
		}
	}

	return "", false, nil
}

func (v *Version) buildTarGzFile() string {
	v.tarGzFile = fmt.Sprintf("go%v.%s-%s.tar.gz", v.Semantics, v.os, v.arch)
	return v.tarGzFile
}

func (v *Version) buildDownloadURL(inCn bool) (string, error) {
	if v.tarGzFile == "" {
		return "", errors.New("version.tarGzFile not built")
	}

	var prefix string
	if inCn {
		prefix = prefixOfDownloadURLCn
	} else {
		prefix = prefixOfDownloadURL
	}

	v.downloadURL = prefix + v.tarGzFile
	return v.downloadURL, nil
}

func (v *Version) IsInstalled() bool {
	if v.isDecompressed == true {
		return true
	}

	return false
}

func (v *Version) GetInstallationDir() string {
	return v.dir
}
