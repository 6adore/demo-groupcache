package cli

import (
	"context"
	"log"
	"net"

	"github.com/6adore/demo-groupcach/config"
	"golang.org/x/crypto/ssh"
)

type SSHDialer struct {
	client *ssh.Client
}

func NewSSHDialer(client *ssh.Client) *SSHDialer {
	return &SSHDialer{client}
}

func (s *SSHDialer) Dial(context context.Context, addr string) (net.Conn, error) {
	return s.client.Dial("tcp", addr)
}

func NewSSHClient() *ssh.Client {
	pswd := config.Conf.SSH.Pswd
	addr := config.Conf.SSH.Addr
	config := &ssh.ClientConfig{
		User:            "root",
		Auth:            []ssh.AuthMethod{ssh.Password(pswd)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Fatal("SSH Dial error: ", err)
	}
	return client
}
