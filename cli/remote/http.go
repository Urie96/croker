package remote

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/fatih/structs"
	"github.com/urie96/croker/server/consts"
)

var client = &http.Client{
	Transport: &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "unix", consts.SockPath)
		},
	},
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
	resp, err := client.Do(req)
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
