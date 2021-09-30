package internal

import (
	"fmt"
	"os"

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
		progressbar.OptionSetPredictTime(true),
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
