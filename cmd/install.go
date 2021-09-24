package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/xvrzhao/gvm/internal"
	e "github.com/xvrzhao/utils/errors"
)

var cmdInstall = &cobra.Command{
	Use:     "install SEMANTIC_VERSION",
	Aliases: []string{"i", "add", "a"},
	Short:   "Install the specified Go version",
	Long:    internal.CmdDescriptionInstall,
	PreRun:  isRootUser,
	RunE:    runCmdInstall,
}

func runCmdInstall(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("need a version of Go to install")
	}

	inCn, _ := cmd.Flags().GetBool("cn")
	v, err := internal.NewVersion(args[0], inCn)
	if err != nil {
		return e.Wrapper(err, "new version error")
	}

	force, _ := cmd.Flags().GetBool("force")

	if err = v.Download(force); err != nil {
		return e.Wrapper(err, "download error")
	}

	if err = v.Decompress(force); err != nil {
		return e.Wrapper(err, "decompress error")
	}

	wantToSwitch, _ := cmd.Flags().GetBool("switch")
	if wantToSwitch {
		app.SetArgs([]string{"switch", v.String()})
		if err = app.Execute(); err != nil {
			return e.Wrapper(err, "switch command executing error")
		}
	}

	return nil
}

func init() {
	cmdInstall.Flags().Bool("cn", false, "use https://golang.google.cn to download")
	cmdInstall.Flags().BoolP("force", "f", false, "ignore the installation, download and install again")
	cmdInstall.Flags().BoolP("switch", "s", false, "switch to the version after its installation")

	app.AddCommand(cmdInstall)
}
