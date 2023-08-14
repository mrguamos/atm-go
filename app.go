package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"github.com/spf13/viper"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/crypto/ssh"
)

type App struct {
	ctx            context.Context
	messageService *messageService
	db             *sqlx.DB
	configService  *configService
	sshClient      *ssh.Client
	listener       net.Listener
}

func NewApp() *App {
	return &App{}
}

type atmSwitch interface {
	pack(message Message) ([]byte, error)
	unpack(r io.Reader) (AtmResponse, error)
	build(message *Message, reversal bool)
}

func (a *App) shutdown(ctx context.Context) {
	if a.db != nil {
		a.db.Close()
	}
	if a.listener != nil {
		a.listener.Close()
	}
	if a.sshClient != nil {
		a.sshClient.Close()
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	db, err := sqlx.Connect("sqlite3", dirname+"/atm.db")
	if err != nil {
		log.Fatal(err)
	}
	a.db = db
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		log.Fatal(err)
	}

	if err := goose.Up(db.DB, "migrations"); err != nil {
		log.Fatal(err)
	}

	configService := &configService{db: db}
	a.configService = configService
	configs, err := configService.loadConfigs()
	if err != nil {
		log.Fatal("unable to load configs:", err)
	}
	for _, v := range configs {
		viper.SetDefault(v.Key, v.Value)
	}

	messageService := &messageService{db: db}
	a.messageService = messageService

}

func (a *App) SendMessage(message Message) (AtmResponse, error) {
	atmSwitch, err := getAtmSwitch(message)
	if err != nil {
		return AtmResponse{}, err
	}
	atmSwitch.build(&message, false)
	return a.messageService.sendTcpMessage(atmSwitch, message)
}

func (a *App) SendReversalMessage(id int) (AtmResponse, error) {
	message, err := a.messageService.getMessage(id)
	if err != nil {
		return AtmResponse{}, err
	}
	atmSwitch, err := getAtmSwitch(message)
	if err != nil {
		return AtmResponse{}, err
	}
	atmSwitch.build(&message, true)
	return a.messageService.sendTcpMessage(atmSwitch, message)
}

func getAtmSwitch(message Message) (atmSwitch, error) {
	var atmSwitch atmSwitch
	switch message.Switch {
	case CORTEX:
		atmSwitch = cortex
	default:
		return nil, fmt.Errorf("atm switch not supported")
	}
	return atmSwitch, nil
}

func (a *App) GetMessages(page int) ([]Message, error) {
	messages, err := a.messageService.getMessages(page)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (a *App) GetConfigs() []Config {
	configs, err := a.configService.getConfigs()
	if err != nil {
		log.Println(err)
	}
	return configs
}

func (a *App) UpdateConfigs(configs []Config) error {
	err := a.configService.updateConfigs(configs)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) UseTunnel() error {
	sshClient, listener, err := tunnel(a.ctx)
	if err != nil {
		log.Println("unable to load ssh tunnel:", err)
		return err
	}
	a.sshClient = sshClient
	a.listener = listener
	log.Println("tunnel successful")
	return nil
}

func (a *App) CloseTunnel() error {
	if a.listener != nil {
		err := a.listener.Close()
		if err != nil {
			return err
		}

	}
	if a.sshClient != nil {
		err := a.sshClient.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) PingTunnel() {
	if a.sshClient != nil {
		_, _, err := a.sshClient.SendRequest("keepalive@openssh.com", true, nil)
		if err != nil {
			runtime.EventsEmit(a.ctx, "tunnel", false, false)
			return
		}
		runtime.EventsEmit(a.ctx, "tunnel", true, false)
		return
	}
	runtime.EventsEmit(a.ctx, "tunnel", false, false)
}
