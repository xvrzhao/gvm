/*
Copyright © 2020 Xavier Zhao <xvrzhao@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	e "github.com/xvrzhao/utils/errors"
	"gvm/funcs"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:     "install SEMANTIC_VERSION",
	Aliases: []string{"i", "add"},
	Short:   "Install a specific version of Go",
	Long:    "Install a specific version of Go, such as: `sudo gvm install 1.14.6`. If you are in China, add the flag `--cn`.",
	PreRun:  isRootUser,
	RunE: func(cmd *cobra.Command, args []string) error {
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
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().Bool("cn", false,
		"Use https://golang.google.cn to download.")
	installCmd.Flags().BoolP("force", "f", false,
		"Ignore the already installed, download and install again.")
	installCmd.Flags().BoolP("switch", "s", false,
		"Switch to the version after installation.")
}
