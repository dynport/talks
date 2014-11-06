package main

import (
	"net"
	"os"

	"code.google.com/p/go.crypto/ssh"
	"code.google.com/p/go.crypto/ssh/agent"
)

func openSSH() (*ssh.Client, error) {
	auths := []ssh.AuthMethod{}
	if a, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		auths = append(auths, ssh.PublicKeysCallback(agent.NewClient(a).Signers))
	}
	config := &ssh.ClientConfig{User: "root", Auth: auths}
	return ssh.Dial("tcp", "128.199.54.55:22", config)
}
