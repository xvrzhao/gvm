package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xvrzhao/gvm/internal"
)

var cmdInstall = &cobra.Command{
	Use:     "install SEMANTIC_VERSION",
	Aliases: []string{"i", "add", "get"},
	Short:   "Install the specified Go version",
	Long:    internal.CmdDescriptionInstall,
	PreRun:  isRootUser,
	RunE:    runCmdInstall,
	PostRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("\033[2KDone!")
	},
}

func runCmdInstall(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return internal.ErrNoVersionSpecified
	}

	inCn, _ := cmd.Flags().GetBool("cn")
	v, err := internal.NewVersion(args[0], inCn)
	if err != nil {
		return fmt.Errorf("internal.NewVersion failed: %w", err)
	}

	force, _ := cmd.Flags().GetBool("force")
	if err = v.Download(force); err != nil {
		return fmt.Errorf("version.Download failed: %w", err)
	}

	if err = v.Decompress(force); err != nil {
		return fmt.Errorf("version.Decompress failed: %w", err)
	}

	wantToSwitch, _ := cmd.Flags().GetBool("switch")
	if !wantToSwitch {
		return nil
	}

	if err = internal.SwitchVersion(v); err != nil {
		return fmt.Errorf("SwitchVersion failed: %w", err)
	}

	return nil
}

func init() {
	cmdInstall.Flags().Bool("cn", false, "use https://golang.google.cn to download")
	cmdInstall.Flags().BoolP("force", "f", false, "ignore the installation, download and install again")
	cmdInstall.Flags().BoolP("switch", "s", false, "switch to the version after its installation")

	App.AddCommand(cmdInstall)
}
