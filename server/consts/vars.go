package consts

import "os"

var (
	ConjDir      = home() + "/.croker"
	JobsSyncPath = ConjDir + "/sync"
	JobLogDir    = ConjDir + "/log"
	SockPath     = ConjDir + "/croker.sock"
)

func init() {
	if err := os.MkdirAll(ConjDir, 0777); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(JobLogDir, 0777); err != nil {
		panic(err)
	}
}

func home() string {
	h, _ := os.UserHomeDir()
	return h
}
