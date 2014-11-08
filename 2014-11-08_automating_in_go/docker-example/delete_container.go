package main

import "net/http"

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
