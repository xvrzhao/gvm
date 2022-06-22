package internal

const (
	prefixOfDownloadURL   = "https://golang.org/dl/"
	prefixOfDownloadURLCn = "https://golang.google.cn/dl/"

	tmpPath  = "/tmp/gvm"
	usrLocal = "/usr/local"
	gvmRoot  = usrLocal + "/gvm"
	goRoot   = usrLocal + "/go"
)

const (
	CmdDescriptionRoot = `GVM is a go version manager. You can use commands of install, list, switch 
and remove to manage local installation of multiple Go versions.

GVM is just support for Unix-like system yet, and the working mechanism of it 
is very simple. GVM will create a gvm directory in /usr/local to host multiple 
versions of GOROOT, and create a symbol link named go in /usr/local referring 
to the specific version in gvm directory. So, you just need to add /usr/local/go/bin 
to PATH environment variable to run go command, and use gvm to switch the 
reference of the symbol link.

Multiple versions of Go installed by GVM can share the same GOPATH compatibly, 
and this is also advocated by GVM.`

	CmdDescriptionInstall = `Install a specific Go version, such as 'sudo gvm install 1.18.3' or 
'sudo gvm install 1.15 -s', if you are in China, do not forget to add the 
flag '--cn'.`

	CmdDescriptionSwitch = `Switch to the specified Go version. You can add the flag '-i' to install 
the version if it's not installed yet, do not forget to add the flag 
'--cn' if you are in China and add '-i'.`
)
