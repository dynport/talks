package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	dh        = "http://127.0.0.1:4243"
	imageName = "hello"
)

func main() {
	if err := run(); err != nil {
		logger.Fatal(err)
	}
}

func deleteContainer(cl *http.Client, id string) error {
	req, err := http.NewRequest("DELETE", dh+"/containers/"+id+"?force=true", nil)
	if err != nil {
		return err
	}
	rsp, err := cl.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	if err := checkResponse(rsp); err != nil {
		return err
	}
	logger.Printf("deleted container status=%s", rsp.Status)
	return nil
}

// http://docs.docker.com/reference/api/docker_remote_api_v1.15/
func run() error {
	cl, err := httpClient()
	if err != nil {
		return err
	}

	logger.Print("creating image")
	if err := createImage(cl); err != nil {
		return err
	}
	logger.Printf("creating container")
	c, err := createContainer(cl)
	if err != nil {
		return err
	}
	defer func() {
		if err := deleteContainer(cl, c.ID); err != nil {
			logger.Printf("ERROR deleting container: %s", err)
		}
	}()
	logger.Printf("created container %s", c.ID)
	rsp, err := cl.Post(dh+"/containers/"+c.ID+"/start", "application/json", nil)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	logger.Print("started container " + c.ID)
	req, err := http.NewRequest("GET", dh+"/containers/"+c.ID+"/json", nil)
	if err != nil {
		return err
	}
	if err := loadRequest(cl, req, &c); err != nil {
		return err
	}
	logger.Printf("got ip: %q", c.NetworkSettings.IPAddress)
	printed := false
	for i := 0; i < 100; i++ {
		err := func() error {
			rsp, err = cl.Get("http://" + c.NetworkSettings.IPAddress + ":8080")
			if err != nil {
				return err
			}
			defer rsp.Body.Close()
			if err := checkResponse(rsp); err != nil {
				return err
			}
			m := "got reponse %s"
			if printed {
				m = "\n" + m
			}
			logger.Printf(m, rsp.Status)
			io.Copy(os.Stdout, rsp.Body)
			return nil
		}()
		if err != nil {
			printed = true
			fmt.Print(".")
			time.Sleep(100 * time.Millisecond)
		} else {
			break
		}
	}
	return nil
}

type Message struct {
	Stream string `json:"stream"`
}

type Container struct {
	Image           string `json:",omitempty"`
	ID              string `json:",omitempty"`
	NetworkSettings struct {
		IPAddress string
	}
}

const dockerfile = `FROM ubuntu:utopic
RUN apt-get update
RUN apt-get install -y curl build-essential libyaml-dev libxml2-dev libxslt1-dev libreadline-dev libssl-dev zlib1g-dev
RUN apt-get install -y ruby ruby-dev
RUN gem install --no-ri --no-rdoc puma

ADD /config.ru config.ru

ENTRYPOINT puma /config.ru -b tcp://0.0.0.0:8080
`

func listImages(cl *http.Client) error {
	return debugRequest(cl, "/images/json")
}

func listContainers(cl *http.Client) error {
	return debugRequest(cl, "/containers/json?all=true")
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
