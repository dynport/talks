package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func httpClient() (*http.Client, error) {
	c, err := openSSH()
	if err != nil {
		return nil, err
	}
	return &http.Client{Transport: &http.Transport{Dial: c.Dial}}, nil
}

func loadRequest(cl *http.Client, req *http.Request, i interface{}) error {
	rsp, err := cl.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	if err := checkResponse(rsp); err != nil {
		return err
	}
	return json.NewDecoder(rsp.Body).Decode(&i)
}

func checkResponse(rsp *http.Response) error {
	if rsp.Status[0] != '2' {
		b, _ := ioutil.ReadAll(rsp.Body)
		return fmt.Errorf("expected status 2xx, got %s: %s (%s)", rsp.Status, b, rsp.Request.URL)
	}
	return nil
}
