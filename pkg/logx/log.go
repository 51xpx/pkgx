package logx

// Mod log mod
type Mod string

const (
	// ModConsole console
	ModConsole = "console"
	// ModFileSizeRotate log file with size rotate
	ModFileSizeRotate = "fileSizeRotate"
	// ModK8s log file with internal k8s
	ModK8s = "k8s"
)
