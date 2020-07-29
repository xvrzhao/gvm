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
