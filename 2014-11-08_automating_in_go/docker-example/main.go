package main

import (
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
			return nil
		}
	}
	return nil
}
