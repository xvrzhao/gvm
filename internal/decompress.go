package internal

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

func decompress(semanticVersion, tarGzFile string) (dir string, err error) {
	goDir := filepath.Join(gvmRoot, "go")
	vgoDir := filepath.Join(gvmRoot, fmt.Sprintf("go%v", semanticVersion))

	if err = os.RemoveAll(goDir); err != nil {
		err = fmt.Errorf("failed to remove directory(%s): %w", goDir, err)
		return
	}
	if err = os.RemoveAll(vgoDir); err != nil {
		err = fmt.Errorf("failed to remove directory(%s): %w", vgoDir, err)
		return
	}

	oldUmask := syscall.Umask(0)
	defer syscall.Umask(oldUmask)

	if err = os.MkdirAll(gvmRoot, os.ModePerm); err != nil {
		err = fmt.Errorf("failed to mkdir %q: %w", gvmRoot, err)
		return
	}

	if err = decompressUsingTar(tarGzFile, gvmRoot); err != nil {
		err = fmt.Errorf("failed to decompressUsingTar: %w", err)
		return
	}

	if err = os.Rename(goDir, vgoDir); err != nil {
		err = fmt.Errorf("failed to rename %s to %s: %w", goDir, vgoDir, err)
		return
	}

	return vgoDir, nil
}

func decompressUsingTar(tarGzFile, dstPath string) error {
	finishEvent := make(chan struct{})
	wg := new(sync.WaitGroup)

	wg.Add(1)
	go func(ch <-chan struct{}, wg *sync.WaitGroup) {
		defer wg.Done()

		resetGlobalProgressBar(100, "Decompressing...")
		defer globalProgressBar.Clear()

		ticker := time.NewTicker(time.Millisecond * 50)
		defer ticker.Stop()

		for {
			select {
			case <-ch:
				globalProgressBar.Finish()
				time.Sleep(time.Second)
				return
			case <-ticker.C:
				if !globalProgressBar.IsFinished() {
					globalProgressBar.Add(1)
				}
			}
		}
	}(finishEvent, wg)

	cmd := exec.Command("tar", "-C", dstPath, "-xzf", tarGzFile)
	cmdStdErrBuf := new(bytes.Buffer)
	cmd.Stderr = cmdStdErrBuf

	err := cmd.Run()
	finishEvent <- struct{}{}
	close(finishEvent)
	wg.Wait()

	if err != nil {
		err = fmt.Errorf("command run error(%q), stderr output(%q)", err.Error(), cmdStdErrBuf.String())
		return fmt.Errorf("failed to run command: %w", err)
	}

	return nil
}
