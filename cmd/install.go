package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xvrzhao/gvm/funcs"
	e "github.com/xvrzhao/utils/errors"
)

var installCmd = &cobra.Command{
	Use:     "install SEMANTIC_VERSION",
	Aliases: []string{"i", "add", "a"},
	Short:   "Install a specific version of Go",
	Long:    "Install a specific version of Go, such as: `sudo gvm install 1.14.6`. If you are in China, add the flag `--cn`.",
	PreRun:  isRootUser,
	RunE:    runInstall,
}

func runInstall(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("need a version of Go to install")
	}

	inCn, _ := cmd.Flags().GetBool("cn")
	v, err := funcs.NewVersion(args[0], inCn)
	if err != nil {
		return e.Wrapper(err, "new version error")
	}

	force, _ := cmd.Flags().GetBool("force")

	fmt.Print("downloading ... ")

	if err = v.Download(force); err != nil {
		return e.Wrapper(err, "download error")
	}

	fmt.Print("done\ndecompressing ... ")

	if err = v.Decompress(force); err != nil {
		return e.Wrapper(err, "decompress error")
	}

	fmt.Println("done")

	wantToSwitch, _ := cmd.Flags().GetBool("switch")
	if wantToSwitch {
		rootCmd.SetArgs([]string{"switch", v.String()})
		if err = rootCmd.Execute(); err != nil {
			return e.Wrapper(err, "switch command executing error")
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().Bool("cn", false,
		"Use https://golang.google.cn to download")
	installCmd.Flags().BoolP("force", "f", false,
		"Ignore the installation, download and install again")
	installCmd.Flags().BoolP("switch", "s", false,
		"Switch to the version after its installation")
}
