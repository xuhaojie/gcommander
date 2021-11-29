package waker

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"time"
)

type Waker struct {
	url    string
	client *ssh.Client
}

func New(url, user, password string) (*Waker, error) {
	config := &ssh.ClientConfig{
		Timeout:         5 * time.Second,
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		//HostKeyCallback: hostKeyCallBackFunc(h.Host),
	}

	config.Auth = []ssh.AuthMethod{ssh.Password(password)}

	sshClient, err := ssh.Dial("tcp", url, config)
	if err != nil {
		return nil, err
	}

	return &Waker{
		url,
		sshClient,
	}, nil
}

func (w *Waker) execute_command(cmd string) error {
	session, err := w.client.NewSession()
	if err != nil {
		fmt.Println("create session failed, ", err)
		return err
	}

	defer func() {
		err := session.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	combo, err := session.CombinedOutput(cmd)
	if err != nil {
		fmt.Println("execute command failed, ", err)
	}
	fmt.Println("cmd output: ", string(combo))
	return nil
}

func (w *Waker) Wake(target string) error {
	cmd := fmt.Sprintf("ether-wake -i br0 -b %s", target)
	return w.execute_command(cmd)
}
