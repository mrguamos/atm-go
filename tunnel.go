package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/spf13/viper"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
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

	log.Printf("SSH tunnel established: %s -> %s\n", localAddr, targetAddr)
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

func publicKeyFile(file string) ssh.AuthMethod {
	key, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("Failed to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKeyWithPassphrase(key, []byte(viper.GetString("SSH_PASSPHRASE")))
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	return ssh.PublicKeys(signer)
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
		return publicKeyFile(sshKey), nil
	}

	sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		return nil, err
	}

	agentClient := agent.NewClient(sshAgent)
	signers, err := agentClient.Signers()
	if err != nil {
		return nil, err
	}

	return ssh.PublicKeys(signers...), nil
}
