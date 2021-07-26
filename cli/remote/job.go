package remote

import (
	"fmt"
	"io"
	"os"

	"github.com/urie96/croker/model"
)

type Job struct{}

func (Job) Infos() []*model.JobInfo {
	var resp model.GetJobResponse
	logFatal(httpDo("GET", "/jobs", nil, &resp))
	logFatal(resp.Error())
	return resp.JobInfos
}

func (Job) Log(id string, follow bool) {
	url := fmt.Sprintf("%s/job/%s/log?follow=%v", host, id, follow)
	resp, err := client.Get(url)
	logFatal(err)
	_, err = io.Copy(os.Stdout, resp.Body)
	logFatal(err)
}

func (Job) Add(schedule, command string) string {
	var resp model.AddJobResponse
	req := model.AddJobRequest{
		Command:  command,
		CronSpec: schedule,
	}
	logFatal(httpDo("POST", "/job", req, &resp))
	logFatal(resp.Error())
	return resp.JobID
}

func (Job) Stop(id string) string {
	var resp model.StopJobResponse
	path := fmt.Sprintf("/job/%s/stop", id)
	logFatal(httpDo("PUT", path, nil, &resp))
	logFatal(resp.Error())
	return resp.JobID
}

func (Job) Start(id string) {
	var resp model.BaseResponse
	path := fmt.Sprintf("/job/%s/start", id)
	logFatal(httpDo("PUT", path, nil, &resp))
	logFatal(resp.Error())
}

func (Job) Remove(id string) {
	var resp model.BaseResponse
	path := "/job/" + id
	logFatal(httpDo("DELETE", path, nil, &resp))
	logFatal(resp.Error())
}
