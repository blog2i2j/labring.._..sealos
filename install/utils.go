package install

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

//username
var (
	User   string
	Passwd string
)

//Cmd is
func Cmd(host string, cmd string) []byte {
	session, err := Connect(User, Passwd, host)
	if err != nil {
		fmt.Println("	Error create ssh session failed", err)
		panic(1)
		return []byte{}
	}
	defer session.Close()

	b, err := session.CombinedOutput(cmd)
	fmt.Println("\n\n exec command\n")
	fmt.Println(host, "    ", cmd)
	fmt.Printf("%s\n\n", b)
	if err != nil {
		fmt.Println("	Error exec command failed", err)
		panic(1)
		return []byte{}
	}
	return b
}

//Connect is
func Connect(user, passwd, host string) (*ssh.Session, error) {
	auth := []ssh.AuthMethod{ssh.Password(passwd)}
	config := ssh.Config{
		Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"},
	}

	clientConfig := &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: time.Duration(5) * time.Minute,
		Config:  config,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	addr := fmt.Sprintf("%s:22", host)
	client, err := ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		return nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return nil, err
	}

	return session, nil
}
