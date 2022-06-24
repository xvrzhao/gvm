package internal

import (
	"context"
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
		return fmt.Errorf("failed to Download: %w", err)
	}

	if err := v.Decompress(force); err != nil {
		return fmt.Errorf("failed to decompress: %w", err)
	}

	return nil
}

func (v *Version) Download(force bool) error {
	fmt.Printf("Downloading go%s archive ...    ", v.Semantics)

	if v.isDownloaded && !force {
		fmt.Println("\b\b\bcached")
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()
	file, err := v.download(ctx)
	if err != nil {
		fmt.Println("\b\b\bfailed")
		return fmt.Errorf("failed to download: %w", err)
	}

	fmt.Println("\b\b\bdone")
	v.isDownloaded, v.downloadedTarGzFile = true, file
	return nil
}

func (v *Version) download(ctx context.Context) (downloadedTarGzFile string, err error) {
	// mkdir -p /tmp/gvm
	oldUmask := syscall.Umask(0)
	defer syscall.Umask(oldUmask)
	if err = os.MkdirAll(tmpPath, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to mkdir %s: %w", tmpPath, err)
	}

	// issue request
	res, err := http.Get(v.downloadURL)
	if err != nil {
		return "", fmt.Errorf("failed to GET %s: %w", v.downloadURL, err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to GET %s: status code %d", v.downloadURL, res.StatusCode)
	}

	// create dst file
	dstFile := filepath.Join(tmpPath, v.tarGzFile)
	file, err := os.Create(dstFile)
	if err != nil {
		return "", fmt.Errorf("failed to create dstFile(%s): %w", dstFile, err)
	}
	defer file.Close()

	// copy
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	buf := make([]byte, 1<<10) // 1kb
	total := res.ContentLength
	read := 0
	for {
		select {
		case <-ticker.C:
			progress := float32(read) * 100 / float32(total)
			if progress < 100 {
				fmt.Printf("\b\b\b%2.0f%%", progress)
			}
		case <-ctx.Done():
			return "", fmt.Errorf("failed to copy from res.Body to file: %w", ctx.Err())
		default:
			n, err := res.Body.Read(buf)
			read += n
			file.Write(buf[:n])
			if err == io.EOF {
				return dstFile, nil
			}
			if err != nil {
				return "", fmt.Errorf("failed to copy from res.Body to file: %w", err)
			}
		}
	}
}

func (v *Version) Decompress(force bool) error {
	fmt.Print("Decompressing the archive ... ")

	if v.isDecompressed && !force {
		fmt.Println("cached")
		return nil
	}

	if !v.isDownloaded {
		fmt.Println("failed")
		return errors.New("version is not downloaded")
	}

	dir, err := v.decompress()
	if err != nil {
		fmt.Println("failed")
		return fmt.Errorf("failed to decompress: %w", err)
	}

	fmt.Println("done")
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

	fmt.Printf("Switching to go%s ... ", v.Semantics)

	fi, err := os.Lstat(goRoot)
	if !os.IsNotExist(err) {
		if err != nil {
			fmt.Println("failed")
			return fmt.Errorf("failed to Lstat goRoot: %w", err)
		}

		if fi.Mode()&os.ModeSymlink != 0 {
			if err := os.Remove(goRoot); err != nil {
				fmt.Println("failed")
				return fmt.Errorf("failed to remove goRoot: %w", err)
			}
		} else {
			if err := os.Rename(goRoot, fmt.Sprintf("%s.backup.%s", goRoot, time.Now().Format("20060102150405"))); err != nil {
				fmt.Println("failed")
				return fmt.Errorf("failed to backupOldGoRoot: %w", err)
			}
		}
	}

	if err := os.Symlink(v.GetInstallationDir(), goRoot); err != nil {
		fmt.Println("failed")
		return fmt.Errorf("failed to Symlink: %w", err)
	}

	fmt.Println("done")
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
	return v.isDecompressed
}

func (v *Version) GetInstallationDir() string {
	return v.dir
}
