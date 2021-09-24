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
	Long:    `List all Go versions installed by GVM, the current version marked by '>'.`,
	RunE:    runCmdList,
}

func runCmdList(cmd *cobra.Command, args []string) error {
	versions, err := internal.GetInstalledGoVersionStrings()
	if err != nil {
		return fmt.Errorf("failed to GetInstalledGoVersionStrings: %w", err)
	}
	if len(versions) <= 0 {
		fmt.Println("no Go version managed by GVM yet")
		return nil
	}

	curVersion, noVersionErr := internal.GetCurrentVersionStr()

	for _, version := range versions {
		if noVersionErr == nil && curVersion == version {
			fmt.Printf("> \033[36m%s\033[0m\n", curVersion)
		} else {
			fmt.Printf("  %s\n", version)
		}
	}

	return nil
}

func init() {
	app.AddCommand(cmdList)
}
