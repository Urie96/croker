package remote

import (
	"github.com/urie96/croker/model"
)

const host = "http://sock"

type Server struct{}

func (Server) IsRunning() bool {
	var resp model.GetServerDetailResponse
	err := httpDo("GET", "/detail", nil, &resp)
	if err != nil {
		return false
	}
	return resp.Status == "running"
}

func (Server) Version() string {
	var resp model.GetServerDetailResponse
	httpDo("GET", "/detail", nil, &resp)
	return resp.Version
}
