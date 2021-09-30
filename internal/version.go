package internal

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"

	e "github.com/xvrzhao/utils/errors"
)

type Version struct {
	// env
	Semantics
	os          string
	arch        string
	tarGzFile   string
	downloadURL string

	// download
	isDownloaded        myBool
	downloadedTarGzFile string

	// decompress
	isDecompressed myBool
	dir            string
}

func NewVersion(version string, inCn bool) (v *Version, err error) {
	sem, err := NewSemantics(version)
	if err != nil {
		return nil, fmt.Errorf("NewSemantics failed: %w", err)
	}

	v = &Version{
		Semantics:   sem,
		os:          runtime.GOOS,
		arch:        runtime.GOARCH,
		tarGzFile:   "",
		downloadURL: "",

		isDownloaded:        unknown,
		downloadedTarGzFile: "",

		isDecompressed: unknown,
		dir:            "",
	}

	v.buildTarGzFile()
	v.buildDownloadURL(inCn)

	if err = v.Reload(); err != nil {
		err = e.Wrapper(err, "version reload error")
		return
	}

	return
}

func (v *Version) Reload() error {
	if filePath, isDownloaded, err := v.checkDownload(); err != nil {
		return e.Wrapper(err, "checkDownloading error")
	} else if isDownloaded {
		v.isDownloaded, v.downloadedTarGzFile = yes, filePath
	} else {
		v.isDownloaded, v.downloadedTarGzFile = no, ""
	}

	if dir, isInstalled, err := v.checkInstallation(); err != nil {
		return e.Wrapper(err, "checkInstallation error")
	} else if isInstalled {
		v.isDecompressed, v.dir = yes, dir
	} else {
		v.isDecompressed, v.dir = no, ""
	}

	return nil
}

func (v *Version) buildTarGzFile() string {
	v.tarGzFile = fmt.Sprintf("go%v.%s-%s.tar.gz", v.Semantics, v.os, v.arch)
	return v.tarGzFile
}

func (v *Version) buildDownloadURL(inCn bool) string {
	if v.tarGzFile == "" {
		panic("version.tarGzFile not built")
	}

	pf := prefixOfDownloadURL
	if inCn {
		pf = prefixOfDownloadURLCn
	}

	v.downloadURL = pf + v.tarGzFile
	return v.downloadURL
}

func (v *Version) Download(force bool) (err error) {
	if v.isDownloaded == yes && !force {
		return nil
	}

	file, err := download(v)
	if err != nil {
		return e.Wrapper(err, "download error")
	}

	v.isDownloaded, v.downloadedTarGzFile = yes, file
	return
}

func (v *Version) Decompress(force bool) error {
	if v.isDecompressed == yes && !force {
		return nil
	}

	if v.isDownloaded != yes {
		return errors.New("version is not downloaded")
	}

	dir, err := decompress(v.Semantics.String(), v.downloadedTarGzFile)
	if err != nil {
		return fmt.Errorf("failed to decompress: %w", err)
	}

	v.isDecompressed, v.dir = yes, dir
	return nil
}

func (v *Version) checkDownload() (filePath string, isDownloaded bool, err error) {
	if v.tarGzFile == "" {
		panic("version.tarGzFile not built")
	}

	return checkDownload(v.tarGzFile)
}

func (v *Version) checkInstallation() (versionDir string, isInstalled bool, err error) {
	thisVersionStr := fmt.Sprintf("go%v", v.Semantics)
	versionStrings, err := GetInstalledGoVersionStrings()
	if err != nil {
		return "", false, fmt.Errorf("GetInstalledGoVersionStrings failed: %w", err)
	}

	for _, versionStr := range versionStrings {
		if "go"+versionStr == thisVersionStr {
			return filepath.Join(gvmRoot, thisVersionStr), true, nil
		}
	}

	return "", false, nil
}

func (v *Version) IsInstalled() bool {
	if v.isDecompressed == yes {
		return true
	}

	return false
}

func (v *Version) GetInstallationDir() string {
	return v.dir
}
