package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/xvrzhao/gvm/internal"
)

var App = &cobra.Command{
	Use:              "gvm",
	Short:            "GVM is a go version manager",
	Long:             internal.CmdDescriptionRoot,
	PersistentPreRun: checkOS,
	SilenceErrors:    true,
	SilenceUsage:     true,
}

func checkPermission(cmd *cobra.Command, args []string) {
	if os.Getuid() != 0 {
		fmt.Fprintln(os.Stderr, "Permission denied, please execute this command as the root user.")
		os.Exit(1)
	}
}

func checkOS(cmd *cobra.Command, args []string) {
	if runtime.GOOS == "windows" {
		fmt.Fprintln(os.Stderr, "Sorry, GVM doesn't support Windows platform yet.")
		os.Exit(1)
	}
}
