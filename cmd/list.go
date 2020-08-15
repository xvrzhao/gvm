package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/xvrzhao/gvm/funcs"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "l"},
	Short:   "",
	Long:    ``,
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
