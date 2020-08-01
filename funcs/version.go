package funcs

import (
	"errors"
	"fmt"
	e "github.com/xvrzhao/utils/errors"
	"log"
	"os"
	"path/filepath"
	"runtime"
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

type version struct {
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

func NewVersion(semver string, inCn bool) (v *version, err error) {
	sem, err := checkSemver(semver)
	if err != nil {
		err = e.Wrapper(err, "checkSemver error")
		return
	}

	v = &version{
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

	if downloaded, downloadedTarGzFile, err := v.checkDownloading(); err != nil {
		err = e.Wrapper(err, "checkDownloading error")
		return
	} else if downloaded {
		v.isDownloaded, v.downloadedTarGzFile = yes, downloadedTarGzFile
	} else {
		v.isDownloaded = no
	}

	if installed, dir, err := v.checkInstallation(); err != nil {
		err = e.Wrapper(err, "checkInstallation error")
		return
	} else if installed {
		v.isDecompressed, v.dir = yes, dir
	} else {
		v.isDecompressed = no
	}

	return
}

func (v *version) buildTarGzFile() string {
	v.tarGzFile = fmt.Sprintf("go%v.%s-%s.tar.gz", v.semantics, v.os, v.arch)
	return v.tarGzFile
}

func (v *version) buildDownloadURL(inCn bool) string {
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

func (v *version) Download(force bool) (err error) {
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

func (v *version) Decompress(force bool) error {
	if v.isDecompressed == yes && !force {
		return nil
	}

	if v.isDownloaded != yes {
		return errors.New("version is not downloaded")
	}

	goDir := filepath.Join(gvmRoot, "go")
	vgoDir := filepath.Join(gvmRoot, fmt.Sprintf("go%v", v.semantics))

	if err := os.RemoveAll(goDir); err != nil {
		return e.Wrapper(err, "RemoveAll %s error", goDir)
	}

	if err := os.RemoveAll(vgoDir); err != nil {
		return e.Wrapper(err, "RemoveAll %s error", vgoDir)
	}

	if err := decompress(v.downloadedTarGzFile, gvmRoot); err != nil {
		return e.Wrapper(err, "decompress error")
	}

	if err := os.Rename(goDir, vgoDir); err != nil {
		return e.Wrapper(err, "error when renaming %s to %s", goDir, vgoDir)
	}

	v.isDecompressed, v.dir = yes, vgoDir
	return nil
}

func (v *version) checkDownloading() (downloaded bool, downloadedTarGzFile string, err error) {
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

	downloaded = true
	return
}

func (v *version) checkInstallation() (installed bool, versionDir string, err error) {
	thisVersionStr := fmt.Sprintf("go%v", v.semantics)
	versionStrings, err := GetInstalledGoVersionStrings()
	if err != nil {
		err = e.Wrapper(err, "GetInstalledGoVersionStrings error")
		return
	}

	for _, versionStr := range versionStrings {
		if versionStr == thisVersionStr {
			installed, versionDir = true, filepath.Join(gvmRoot, thisVersionStr)
			return
		}
	}

	installed = false
	return
}

func (v *version) IsInstalled() bool {
	if v.isDecompressed == no {
		return true
	}

	return false
}
