package cmd

import (
	"fmt"
	"log"
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
		log.Fatal("Permission denied, please execute this command as the root user.")
	}
}

func checkOS(cmd *cobra.Command, args []string) {
	if runtime.GOOS == "windows" {
		log.Fatal("Sorry, GVM does not support Windows platform at the moment.")
	}
}

func printDone(cmd *cobra.Command, args []string) {
	fmt.Println("\033[2KDone!")
}
