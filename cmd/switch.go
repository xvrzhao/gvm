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
	"github.com/xvrzhao/gvm/funcs"
	e "github.com/xvrzhao/utils/errors"
)

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:     "switch SEMANTIC_VERSION",
	Aliases: []string{"s"},
	Short:   "A brief description of your command",
	Long:    ``,
	PreRun:  isRootUser,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) <= 0 {
			return errors.New("need a version of Go to switch")
		}

		inCn, _ := cmd.Flags().GetBool("cn")
		v, err := funcs.NewVersion(args[0], inCn)
		if err != nil {
			return e.Wrapper(err, "new version error")
		}

		wantToInstall, _ := cmd.Flags().GetBool("install")
		if !v.IsInstalled() && wantToInstall {
			args := []string{"install", v.String()}
			if inCn {
				args = append(args, "--cn")
			}

			rootCmd.SetArgs(args)
			if err = rootCmd.Execute(); err != nil {
				return e.Wrapper(err, "install command executing error")
			}
		}

		fmt.Print("switching version ... ")

		if err := funcs.SwitchVersion(v); err != nil {
			return e.Wrapper(err, "switch version error")
		}

		fmt.Println("done")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// switchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// switchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	switchCmd.Flags().Bool("cn", false, "Use https://golang.google.cn to download.")
	switchCmd.Flags().BoolP("install", "i", false, "Install if the version is not installed.")
}
