package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

var logger = log.New(os.Stderr, "", 0)

func debugRequest(cl *http.Client, s string) error {
	rsp, err := cl.Get(dh + s)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	if rsp.Status[0] != '2' {
		return fmt.Errorf("expected status 2xx, got %s", rsp.Status)
	}
	return debugJSON(rsp.Body)
}

func debugJSON(r io.Reader) error {
	c := exec.Command("jq", ".")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = r
	return c.Run()
}

func handleBuildResponse(r io.Reader) error {
	dec := json.NewDecoder(r)
	for {
		var raw json.RawMessage
		err := dec.Decode(&raw)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		m := &Message{}
		if err := json.Unmarshal(raw, &m); err != nil {
			return err
		}
		if m.Stream == "" {
			logger.Printf("%s", raw)
		} else {
			logger.Printf("%s", m.Stream)
		}
	}
	return nil
}
