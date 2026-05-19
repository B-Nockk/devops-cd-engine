package ssh

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"time"

	out "cd-engine/internal/ports/out"

	"golang.org/x/crypto/ssh"
)

type executor struct {
	user       string
	privateKey []byte
	signer     ssh.Signer
}

func NewExecutor(user string, privateKey []byte) (out.Executor, error) {
	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	return &executor{
		user:       user,
		privateKey: privateKey,
		signer:     signer,
	}, nil
}

func (e *executor) Execute(ctx context.Context, target string, command string) (string, error) {
	config := &ssh.ClientConfig{
		User:            e.user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(e.signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", target)
	if err != nil {
		return "", fmt.Errorf("failed to dial target: %w", err)
	}
	defer conn.Close()

	c, chans, reqs, err := ssh.NewClientConn(conn, target, config)
	if err != nil {
		return "", fmt.Errorf("failed to create SSH client: %w", err)
	}
	client := ssh.NewClient(c, chans, reqs)
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	var buf bytes.Buffer
	session.Stdout = &buf
	session.Stderr = &buf

	done := make(chan error, 1)
	go func() {
		done <- session.Run(command)
	}()

	select {
	case <-ctx.Done():
		_ = session.Signal(ssh.SIGKILL)
		return buf.String(), ctx.Err()
	case err := <-done:
		return buf.String(), err
	}
}
