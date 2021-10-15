package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xvrzhao/gvm/internal"
)

var cmdInstall = &cobra.Command{
	Use:     "install VERSION",
	Aliases: []string{"i", "add", "get"},
	Short:   "Install the specified Go version",
	Long:    internal.CmdDescriptionInstall,
	Example: "gvm i 1.17.1 -s -c",

	PreRun:  checkPermission,
	RunE:    runCmdInstall,
	PostRun: printDone,
}

func runCmdInstall(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return internal.ErrNoVersionSpecified
	}

	inCn, _ := cmd.Flags().GetBool("cn")
	v, err := internal.NewVersion(args[0], inCn)
	if err != nil {
		return fmt.Errorf("failed to NewVersion: %w", err)
	}

	force, _ := cmd.Flags().GetBool("force")
	if err := v.Install(force); err != nil {
		return fmt.Errorf("failed to Install: %w", err)
	}

	wantToSwitch, _ := cmd.Flags().GetBool("switch")
	if !wantToSwitch {
		return nil
	}

	if err = v.Switch(); err != nil {
		return fmt.Errorf("failed to switch: %w", err)
	}

	return nil
}

func init() {
	cmdInstall.Flags().BoolP("cn", "c", false, "use https://golang.google.cn to download")
	cmdInstall.Flags().BoolP("force", "f", false, "ignore the installation and download, install again")
	cmdInstall.Flags().BoolP("switch", "s", false, "switch to the version after its installation")

	App.AddCommand(cmdInstall)
}
