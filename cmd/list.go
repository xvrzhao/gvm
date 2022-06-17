package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xvrzhao/gvm/internal"
)

var cmdList = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "l"},
	Short:   "List all Go versions installed by GVM",

	RunE: runCmdList,
}

func runCmdList(cmd *cobra.Command, args []string) error {
	versions, err := internal.GetAllInstalledVersions()
	if err != nil {
		return fmt.Errorf("failed to GetAllInstalledVersions: %w", err)
	}
	if len(versions) <= 0 {
		fmt.Println("Empty.")
		return nil
	}

	curVersion, err := internal.GetCurrentVersion()

	for _, version := range versions {
		if err == nil && curVersion == version {
			fmt.Printf("\033[1;36m✔  %s\033[0m\n", curVersion)
		} else {
			fmt.Printf("   %s\n", version)
		}
	}

	return nil
}

func init() {
	App.AddCommand(cmdList)
}
