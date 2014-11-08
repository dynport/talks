package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Container struct {
	Image           string `json:",omitempty"`
	ID              string `json:",omitempty"`
	NetworkSettings struct {
		IPAddress string
	}
}

func createContainer(cl *http.Client) (*Container, error) {
	c := &Container{Image: imageName}
	b := &bytes.Buffer{}
	if err := json.NewEncoder(b).Encode(c); err != nil {
		return nil, err
	}
	rsp, err := cl.Post(dh+"/containers/create", "application/json", b)
	if err != nil {
		return nil, err
	}
	if err := checkResponse(rsp); err != nil {
		return nil, err
	}
	var cn *Container
	return cn, json.NewDecoder(rsp.Body).Decode(&cn)
}
