package funcs

const (
	dlPrefix   = "https://golang.org/dl/"
	dlPrefixCn = "https://golang.google.cn/dl/"

	tmpPath  = "/tmp/gvm"
	usrLocal = "/usr/local"
	gvmRoot  = usrLocal + "/gvm"
	goRoot   = usrLocal + "/go"
)

type finishState int

const (
	fsUnknown finishState = iota
	fsFinished
	fsUnFinished
)
