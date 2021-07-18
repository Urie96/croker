package job

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/fatih/structs"
	"github.com/urie96/croker/model"
)

const host = "http://localhost:1996"

func Infos() []*model.JobInfo {
	var resp model.GetJobResponse
	logFatal(httpDo("GET", "/jobs", nil, &resp))
	logFatal(resp.Error())
	return resp.JobInfos
}

func Log(id string, follow bool) {
	url := fmt.Sprintf("%s/job/%s/log?follow=%v", host, id, follow)
	resp, err := http.Get(url)
	logFatal(err)
	_, err = io.Copy(os.Stdout, resp.Body)
	logFatal(err)
}

func Add(schedule, command string) string {
	var resp model.AddJobResponse
	req := model.AddJobRequest{
		Command:  command,
		CronSpec: schedule,
	}
	logFatal(httpDo("POST", "/job", req, &resp))
	logFatal(resp.Error())
	return resp.JobID
}

func Stop(id string) string {
	var resp model.StopJobResponse
	path := fmt.Sprintf("/job/%s/stop", id)
	logFatal(httpDo("PUT", path, nil, &resp))
	logFatal(resp.Error())
	return resp.JobID
}

func Start(id string) {
	var resp model.BaseResponse
	path := fmt.Sprintf("/job/%s/start", id)
	logFatal(httpDo("PUT", path, nil, &resp))
	logFatal(resp.Error())
}

func Remove(id string) {
	var resp model.BaseResponse
	path := "/job/" + id
	logFatal(httpDo("DELETE", path, nil, &resp))
	logFatal(resp.Error())
}

func ServerVersion() string {
	var resp model.ServerVersionResponse
	logFatal(httpDo("GET", "/version", nil, &resp))
	logFatal(resp.Error())
	return resp.Version
}

func httpDo(method, path string, reqBody, respBody interface{}) error {
	url := host + path
	var body io.Reader
	if reqBody != nil {
		if method == "" || method == "GET" || method == "DELETE" {
			query := "?"
			for k, v := range structs.Map(reqBody) {
				query += fmt.Sprintf("%s=%v", k, v)
			}
			method += query
		} else {
			data, err := json.Marshal(reqBody)
			if err != nil {
				return err
			}
			body = bytes.NewReader(data)
		}
	}
	req, err := http.NewRequest(method, url, body)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, respBody)
}

func logFatal(err error) {
	if err != nil {
		log.New(os.Stderr, "", 0).Fatalln(err.Error())
	}
}
