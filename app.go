package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/crypto/ssh"
	"gopkg.in/natefinch/lumberjack.v2"
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
	packEchoTest() ([]byte, error)
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
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	logFile := &lumberjack.Logger{
		Filename:   dirname + "/logfile.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   true,
	}

	multiWriter := zerolog.MultiLevelWriter(logFile, zerolog.ConsoleWriter{Out: os.Stdout})
	logger := zerolog.New(multiWriter).With().Caller().Timestamp().Logger()
	log.Logger = logger

	a.ctx = ctx
	db, err := sqlx.Connect("sqlite", dirname+"/atm.db")
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	a.db = db
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite"); err != nil {
		log.Fatal().Err(err).Msg("")
	}

	if err := goose.Up(db.DB, "migrations"); err != nil {
		log.Fatal().Err(err).Msg("")
	}

	configService := &configService{db: db}
	a.configService = configService
	configs, err := configService.loadConfigs()
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	for _, v := range configs {
		viper.SetDefault(v.Key, v.Value)
	}

	messageService := &messageService{db: db}
	a.messageService = messageService

}

func sendMessage(a *App, message Message, reversal bool) (AtmResponse, error) {
	atmSwitch, err := getAtmSwitch(message)
	if err != nil {
		log.Error().Err(err).Msg("")
		return AtmResponse{}, err
	}
	atmSwitch.build(&message, reversal)
	err = a.messageService.saveMessage(message)
	if err != nil {
		log.Error().Err(err).Msg("")
		return AtmResponse{}, err
	}
	b, err := atmSwitch.pack(message)
	if err != nil {
		log.Error().Err(err).Msg("")
		return AtmResponse{}, err
	}
	serverAddr := fmt.Sprintf("%s:%s", viper.GetString("HOST"), viper.GetString("PORT"))
	conn, err := net.DialTimeout("tcp", serverAddr, 30*time.Second)
	if err != nil {
		return AtmResponse{}, err
	}
	defer conn.Close()
	err = a.messageService.sendTcpMessage(conn, b)
	if err != nil {
		return AtmResponse{}, err
	}
	return atmSwitch.unpack(conn)
}

func (a *App) SendFinancialMessage(message Message) (AtmResponse, error) {
	return sendMessage(a, message, false)
}

func (a *App) SendReversalMessage(id int) (AtmResponse, error) {
	message, err := a.messageService.getMessage(id)
	if err != nil {
		log.Error().Err(err).Msg("")
		return AtmResponse{}, err
	}
	return sendMessage(a, message, true)
}

func getAtmSwitch(message Message) (atmSwitch, error) {
	var atmSwitch atmSwitch
	switch message.Switch {
	case CORTEX:
		atmSwitch = cortex
	case POSTBRIDGE:
		atmSwitch = postbridge
	case NARADA:
		atmSwitch = narada
	default:
		return nil, fmt.Errorf("atm switch not supported")
	}
	return atmSwitch, nil
}

func (a *App) GetMessages(page int) ([]Message, error) {
	messages, err := a.messageService.getMessages(page)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}
	return messages, nil
}

func (a *App) GetConfigs() ([]Config, error) {
	configs, err := a.configService.getConfigs()
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}
	return configs, nil
}

func (a *App) UpdateConfigs(configs []Config) error {
	err := a.configService.updateConfigs(configs)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}
	return nil
}

func (a *App) UseTunnel() error {
	sshClient, listener, err := tunnel(a.ctx)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}
	a.sshClient = sshClient
	a.listener = listener
	log.Print("tunnel successful")
	return nil
}

func (a *App) CloseTunnel() error {
	if a.listener != nil {
		err := a.listener.Close()
		if err != nil {
			log.Error().Err(err).Msg("")
			return err
		}

	}
	if a.sshClient != nil {
		err := a.sshClient.Close()
		if err != nil {
			log.Error().Err(err).Msg("")
			return err
		}
	}
	return nil
}

func (a *App) PingTunnel() {
	if a.sshClient != nil {
		_, _, err := a.sshClient.SendRequest("keepalive@openssh.com", true, nil)
		if err != nil {
			log.Error().Err(err).Msg("")
			runtime.EventsEmit(a.ctx, "tunnel", false, false)
			return
		}
		runtime.EventsEmit(a.ctx, "tunnel", true, false)
		return
	}
	runtime.EventsEmit(a.ctx, "tunnel", false, false)
}

func (a *App) OpenFileDialog() (string, error) {
	path, _ := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{})
	return path, nil
}
