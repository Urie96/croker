package consts

import "os"

var (
	ConjPath      = home() + "/.croker"
	JobsSyncPath  = ConjPath + "/sync"
	JobLogPath    = ConjPath + "/log"
	DaemonLogPath = ConjPath + "/log/daemon"
	DaemonPidPath = ConjPath + "/pid"
)

func init() {
	if err := os.MkdirAll(ConjPath, 0777); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(JobLogPath, 0777); err != nil {
		panic(err)
	}
}

func home() string {
	h, _ := os.UserHomeDir()
	return h
}
