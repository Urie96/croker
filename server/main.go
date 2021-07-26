package server

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/urie96/croker/server/consts"
	chttp "github.com/urie96/croker/server/http"
)

func Run() {
	os.Remove(consts.SockPath)
	listener, err := net.Listen("unix", consts.SockPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	http.Serve(listener, chttp.Handler())
}
