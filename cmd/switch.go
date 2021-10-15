package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xvrzhao/gvm/internal"
)

var cmdSwitch = &cobra.Command{
	Use:     "switch SEMANTIC_VERSION",
	Aliases: []string{"s"},
	Short:   "Switch to the specified Go version",
	Long:    internal.CmdDescriptionSwitch,

	PreRun:  checkPermission,
	RunE:    runCmdSwitch,
	PostRun: printDone,
}

func runCmdSwitch(cmd *cobra.Command, args []string) error {
	if len(args) <= 0 {
		return internal.ErrNoVersionSpecified
	}

	inCn, _ := cmd.Flags().GetBool("cn")
	v, err := internal.NewVersion(args[0], inCn)
	if err != nil {
		return fmt.Errorf("failed to NewVersion: %w", err)
	}

	wantToInstall, _ := cmd.Flags().GetBool("install")
	if !v.IsInstalled() && wantToInstall {
		if err = v.Install(false); err != nil {
			return fmt.Errorf("failed to Install: %w", err)
		}
	}

	if err := v.Switch(); err != nil {
		return fmt.Errorf("failed to switch: %w", err)
	}

	return nil
}

func init() {
	cmdSwitch.Flags().Bool("cn", false, "use https://golang.google.cn to download")
	cmdSwitch.Flags().BoolP("install", "i", false, "install if the version is not installed")

	App.AddCommand(cmdSwitch)
}
