package main

import "net/http"

func listImages(cl *http.Client) error {
	return debugRequest(cl, "/images/json")
}

func listContainers(cl *http.Client) error {
	return debugRequest(cl, "/containers/json?all=true")
}
