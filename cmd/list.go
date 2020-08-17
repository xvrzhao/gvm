package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/xvrzhao/gvm/funcs"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "l"},
	Short:   "List all Go versions installed by GVM",
	Long:    `List all Go versions installed by GVM, the current version marked by '*'.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		versions, err := funcs.GetInstalledGoVersionStrings()
		if err != nil {
			return errors.Wrap(err, "GetInstalledGoVersionStrings error")
		}
		if len(versions) <= 0 {
			fmt.Println("empty")
			return nil
		}

		curVersion, noVersionErr := funcs.GetCurrentVersionStr()

		for _, version := range versions {
			if noVersionErr == nil && curVersion == version {
				fmt.Printf("* %s\n", curVersion)
			} else {
				fmt.Printf("  %s\n", version)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
