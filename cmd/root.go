package cmd

import (
	"log"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/xvrzhao/gvm/internal"
)

var app = &cobra.Command{
	Use:   "gvm",
	Short: "GVM is a go version manager",
	Long:  internal.CmdDescriptionRoot,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if runtime.GOOS == "windows" {
			log.Fatal("Sorry, GVM does not support Windows platform at the moment.")
		}
	},
}

func Execute() {
	if err := app.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func isRootUser(cmd *cobra.Command, args []string) {
	if os.Getuid() != 0 {
		log.Fatal("Permission denied, please execute this command as the root user.")
	}
}
