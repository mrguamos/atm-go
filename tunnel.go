package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/spf13/viper"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/crypto/ssh"
)

func tunnel(ctx context.Context) (*ssh.Client, net.Listener, error) {

	authMethod, err := sshAgentAuth()
	if err != nil {
		return nil, nil, err
	}

	config := &ssh.ClientConfig{
		User: viper.GetString("SSH_USERNAME"),
		Auth: []ssh.AuthMethod{
			authMethod,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	bastionAddr := fmt.Sprintf("%s:%s", viper.GetString("BASTION_HOST"), viper.GetString("BASTION_PORT"))
	bastionConn, err := ssh.Dial("tcp", bastionAddr, config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to dial bastion host: %v", err)
	}

	localAddr := fmt.Sprintf("localhost:%s", viper.GetString("SSH_LOCAL_PORT"))
	targetAddr := fmt.Sprintf("%s:%s", viper.GetString("TARGET_HOST"), viper.GetString("SSH_REMOTE_PORT"))
	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create listener on bastion host: %v", err)
	}

	log.Printf("SSH tunnel established: %s -> %s", localAddr, targetAddr)
	runtime.EventsEmit(ctx, "tunnel", true)

	go func() {
		for {
			localConn, err := listener.Accept()
			if err != nil {
				fmt.Println("Failed to accept connection:", err)
				listener.Close()
				bastionConn.Close()
				runtime.EventsEmit(ctx, "tunnel", false)
				break
			}

			go handleLocalConnection(bastionConn, localConn, targetAddr)
		}
	}()

	return bastionConn, listener, nil
}

func publicKeyFile(file string) (ssh.AuthMethod, error) {
	key, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}

	passPhrase := viper.GetString("SSH_PASSPHRASE")
	var signer ssh.Signer

	if len(passPhrase) > 0 {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(key, []byte(viper.GetString("SSH_PASSPHRASE")))
	} else {
		signer, err = ssh.ParsePrivateKey(key)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return ssh.PublicKeys(signer), nil
}

func handleLocalConnection(client *ssh.Client, localConn net.Conn, targetAddr string) {
	targetConn, err := client.Dial("tcp", targetAddr)
	if err != nil {
		log.Printf("Failed to connect to target host: %v", err)
		return
	}

	go io.Copy(localConn, targetConn)
	go io.Copy(targetConn, localConn)
}

func sshAgentAuth() (ssh.AuthMethod, error) {
	sshKey := viper.GetString("SSH_KEY")

	if len(sshKey) > 0 {
		return publicKeyFile(sshKey)
	}

	return nil, errors.New("missing SSH_KEY, please configure it in the settings")
}
