package main

import (
	"bytes"
	"log"
	"os"
	"text/template"
)

var logger = log.New(os.Stderr, "", 0)

func mustRender(t string, i interface{}) string {
	s, err := render(t, i)
	if err != nil {
		logger.Fatal(err)
	}
	return s
}

func render(t string, i interface{}) (string, error) {
	tpl, err := template.New(t).Parse(t)
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	err = tpl.Execute(buf, i)
	return buf.String(), err
}
