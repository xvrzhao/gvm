/*
Copyright © 2020 Xavier Zhao <xvrzhao@gmail.com>
Licensed under the MIT License. See LICENSE file in the project root for license information.
*/

package funcs

const (
	prefixOfDownloadURL   = "https://golang.org/dl/"
	prefixOfDownloadURLCn = "https://golang.google.cn/dl/"

	tmpPath  = "/tmp/gvm"
	usrLocal = "/usr/local"
	gvmRoot  = usrLocal + "/gvm"
	goRoot   = usrLocal + "/go"
)

type myBool int

const (
	unknown myBool = iota
	yes
	no
)
