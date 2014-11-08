package main

import (
	"net/http"
	"net/url"
)

type Message struct {
	Stream string `json:"stream"`
}

func createImage(c *http.Client) error {
	files := []*File{
		{"Dockerfile", dockerfile},
		{"config.ru", configRU},
	}
	b, err := createArchive(files)
	if err != nil {
		return err
	}
	v := url.Values{"t": {"hello"}}
	rsp, err := c.Post(dh+"/build?"+v.Encode(), "application/octet-stream", b)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	if err := checkResponse(rsp); err != nil {
		return err
	}
	return handleBuildResponse(rsp.Body)
}

const configRU = `#!/bin/env ruby

require "json"

app = lambda do |env|
  [200, {"Content-Type" => "application/json"}, [JSON.generate({status: "OK", ruby_version: RUBY_DESCRIPTION}) + "\n"]]
end

run app
`

const dockerfile = `FROM ubuntu:utopic
RUN apt-get update
RUN apt-get install -y curl build-essential libyaml-dev libxml2-dev libxslt1-dev libreadline-dev libssl-dev zlib1g-dev
RUN apt-get install -y ruby ruby-dev
RUN gem install --no-ri --no-rdoc puma

ADD /config.ru config.ru

ENTRYPOINT puma /config.ru -b tcp://0.0.0.0:8080
`
