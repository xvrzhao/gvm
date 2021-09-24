package internal

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	e "github.com/xvrzhao/utils/errors"
)

type semantics struct {
	major, minor, patch uint8
}

func (s semantics) String() string {
	v := fmt.Sprintf("%d.%d", s.major, s.minor)
	if s.patch != 0 {
		v = fmt.Sprintf("%s.%d", v, s.patch)
	}
	return v
}

type Version struct {
	// env
	semantics
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

func NewVersion(semver string, inCn bool) (v *Version, err error) {
	sem, err := checkSemver(semver)
	if err != nil {
		err = e.Wrapper(err, "checkSemver error")
		return
	}

	v = &Version{
		semantics:   sem,
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
	if downloaded, downloadedTarGzFile, err := v.checkDownloading(); err != nil {
		return e.Wrapper(err, "checkDownloading error")
	} else if downloaded {
		v.isDownloaded, v.downloadedTarGzFile = yes, downloadedTarGzFile
	} else {
		v.isDownloaded = no
	}

	if installed, dir, err := v.checkInstallation(); err != nil {
		return e.Wrapper(err, "checkInstallation error")
	} else if installed {
		v.isDecompressed, v.dir = yes, dir
	} else {
		v.isDecompressed = no
	}

	return nil
}

func (v *Version) buildTarGzFile() string {
	v.tarGzFile = fmt.Sprintf("go%v.%s-%s.tar.gz", v.semantics, v.os, v.arch)
	return v.tarGzFile
}

func (v *Version) buildDownloadURL(inCn bool) string {
	if v.tarGzFile == "" {
		log.Fatal("*version.buildDownloadURL: *version.tarGzFile is empty")
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
		err = nil
		return
	}

	file, err := download(v)
	if err != nil {
		err = e.Wrapper(err, "download error")
		return
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

	dir, err := decompress(v.semantics.String(), v.downloadedTarGzFile)
	if err != nil {
		return fmt.Errorf("failed to decompress: %w", err)
	}

	v.isDecompressed, v.dir = yes, dir
	return nil
}

func (v *Version) checkDownloading() (downloaded bool, downloadedTarGzFile string, err error) {
	if v.tarGzFile == "" {
		log.Fatal("*version.checkDownloading: *version.tarGzFile is empty")
	}

	downloadedTarGzFile = filepath.Join(tmpPath, v.tarGzFile)
	_, err = os.Stat(downloadedTarGzFile)
	if os.IsNotExist(err) {
		err = nil
		return
	}
	if err != nil {
		err = e.Wrapper(err, "error when getting fileInfo of downloaded tarGzFile")
		return
	}

	if !isArchiveValid(downloadedTarGzFile) {
		return
	}

	downloaded = true
	return
}

func (v *Version) checkInstallation() (installed bool, versionDir string, err error) {
	thisVersionStr := fmt.Sprintf("go%v", v.semantics)
	versionStrings, err := GetInstalledGoVersionStrings()
	if err != nil {
		err = e.Wrapper(err, "GetInstalledGoVersionStrings error")
		return
	}

	for _, versionStr := range versionStrings {
		if "go"+versionStr == thisVersionStr {
			installed, versionDir = true, filepath.Join(gvmRoot, thisVersionStr)
			return
		}
	}

	installed = false
	return
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
