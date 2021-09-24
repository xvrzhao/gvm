package internal

import (
	"fmt"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

var globalProgressBar *progressbar.ProgressBar

func resetGlobalProgressBar(max int64, description string) {
	if globalProgressBar != nil {
		globalProgressBar.Clear()
	}

	globalProgressBar = progressbar.NewOptions64(max,
		progressbar.OptionSetWriter(os.Stdout),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription(fmt.Sprintf("[cyan]>[reset] %s\t", description)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]·[reset]",
			SaucerHead:    "[green]·[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
}

func Test() {
	resetGlobalProgressBar(100, "Downloading...")
	globalProgressBar.Add(30)
	time.Sleep(time.Second * 1)
	globalProgressBar.Finish()

	time.Sleep(time.Second * 3)

	resetGlobalProgressBar(10, "Decompressing...")
	globalProgressBar.Add(1)

	time.Sleep(time.Hour)
}
