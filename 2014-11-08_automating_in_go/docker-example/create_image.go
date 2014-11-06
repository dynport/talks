package main

import (
	"net/http"
	"net/url"
)

func createImage(c *http.Client) error {
	files := []*File{
		{"Dockerfile", dockerfile},
		{"config.ru", `run lambda { |env| [200, {}, ["OK\n"]] }`},
	}
	b, err := createArchive(files)
	if err != nil {
		return err
	}
	v := url.Values{"t": {"hello"}}
	req, err := http.NewRequest("POST", dh+"/build?"+v.Encode(), b)
	if err != nil {
		return err
	}
	rsp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	if err := checkResponse(rsp); err != nil {
		return err
	}
	return handleBuildResponse(rsp.Body)
}
