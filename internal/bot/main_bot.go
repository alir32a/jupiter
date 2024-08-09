package bot

import (
	"context"
	"fmt"
	"github.com/alir32a/jupiter/config"
	"github.com/alir32a/jupiter/internal/model"
	"github.com/alir32a/jupiter/pkg/tg"
	"github.com/alir32a/jupiter/pkg/util"
	"github.com/charmbracelet/log"
	"strconv"
	"time"
)

const (
	TimeFormat = "2006-01-02T15:04:05 -07:00"
)

type UserService interface {
	CreateUser(ctx context.Context, req model.CreateUserRequest) (model.CreateUserResponse, error)
	ChangePassword(ctx context.Context, username string) (string, error)
}

type ConnectionService interface {
	GetActiveConnections(ctx context.Context, req model.GetActiveConnectionsRequest) (model.GetActiveConnectionsResponse, error)
	GetUserActiveConnections(ctx context.Context, username string) ([]model.ConnectionEntity, error)
}

type PackageService interface {
	GetUserActivePackages(ctx context.Context, username string) (model.GetUserPackages, error)
}

type MainBot struct {
	userSvc        UserService
	connectionSvc  ConnectionService
	packageSvc     PackageService
	bot            *tg.Bot
	cfg            *config.MainBotConfig
	logger         *log.Logger
	queryCommander *QueryCommander
}

func NewMainBot(cfg *config.MainBotConfig, logger *log.Logger, userSvc UserService, connectionSvc ConnectionService,
	packageSvc PackageService) *MainBot {
	return &MainBot{
		userSvc:        userSvc,
		connectionSvc:  connectionSvc,
		packageSvc:     packageSvc,
		bot:            tg.NewBot(cfg.Token),
		cfg:            cfg,
		logger:         logger,
		queryCommander: NewQueryCommander(),
	}
}

func (m MainBot) Run() error {
	if err := m.CheckCommands(); err != nil {
		return err
	}

	m.logger.Info("starting main bot ...")

	m.bot.Run(func(updates []tg.Update) error {
		for _, update := range updates {
			if update.CallbackQuery != nil {
				if err := m.queryCommander.Handle(*update.CallbackQuery); err != nil {
					m.logger.Error(err.Error())
				}

				continue
			}

			if err := m.parseCommand(update.Message); err != nil {
				m.logger.Error(err.Error())

				return err
			}
		}

		return nil
	})

	return nil
}

func (m MainBot) CheckCommands() error {
	commands, err := m.bot.GetCommands()
	if err != nil {
		return err
	}

	cmdMap := mapCommandsByName(commands)

	var missingCommands []tg.BotCommand
	for name, desc := range MainBotCommands {
		if _, ok := cmdMap[name]; !ok {
			missingCommands = append(missingCommands, tg.BotCommand{
				Command:     name,
				Description: desc,
			})
		}
	}

	if len(missingCommands) > 0 {
		err := m.bot.SetCommands(missingCommands...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b MainBot) parseCommand(msg tg.Message) error {
	if msg.Type != tg.MessageTypeCommand {
		return b.SendUnknownMessage(msg.From.ID)
	}

	switch msg.Text {
	case "/start":
		return b.Start(msg)
	case "/create":
		return b.CreateUser(msg)
	case "/status":
		return b.GetStatus(msg)
	case "/password":
		return b.ChangePassword(msg)
	case "/connections":
		return b.GetActiveConnections(msg)
	default:
		return b.SendUnknownMessage(msg.From.ID)
	}
}

func (b MainBot) Start(msg tg.Message) error {
	reply := `
	hey there 👋
	this is jupiter bot, you can manage your openconnect vpn account here.
    currently, you can't buy packages from this bot (yet).
	you can create user (you can create only one user per telegram account),
    see your active package (remaining traffic, expire time etc.),
	change your account password and see active connections
    let's f**ck sansoorchi.✊

	supported commands:
    - /start: show this message
    - /status: show your active package status
    - /create: create a user and show credentials, and activate a trial package if trial is activated by administrators
    - /password: change your account password
    - /connections: show active connections
`

	_, err := b.bot.SendMessage(tg.SendMessageRequest{ChatID: msg.From.ID, Text: reply})
	if err != nil {
		return err
	}

	return nil
}

func (b MainBot) CreateUser(msg tg.Message) error {
	username := msg.From.Username
	if username == "" {
		username = strconv.Itoa(msg.From.ID)
	}

	user, err := b.userSvc.CreateUser(context.Background(), model.CreateUserRequest{
		Username:   username,
		ExternalID: strconv.Itoa(msg.From.ID),
		UserType:   model.UserTypeTelegram,
	})
	if err != nil {
		return err
	}

	reply := fmt.Sprintf(
		`here is your username and password, this message will be deleted in an hour, 
         you can change your password anytime using /passaword
         Userame: %s
         Password: %s`, user.Username, user.Password)
	resp, err := b.bot.SendMessage(tg.SendMessageRequest{ChatID: msg.From.ID, Text: reply, ProtectContent: true})
	if err != nil {
		return err
	}

	if len(resp) > 0 {
		time.AfterFunc(1*time.Hour, func() {
			if err := b.bot.DeleteMessage(msg.From.ID, resp[0].MessageID); err != nil {
				b.logger.Error(err.Error())
			}
		})
	}

	return nil
}

func (b MainBot) GetStatus(msg tg.Message) error {
	packages, err := b.packageSvc.GetUserActivePackages(context.Background(), msg.From.Username)
	if err != nil {
		_, err := b.bot.SendMessage(tg.SendMessageRequest{ChatID: msg.From.ID, Text: err.Error()})
		if err != nil {
			return err
		}

		return nil
	}

	if packages.ActivePackage.ID != 0 && len(packages.ReservedPackages) <= 0 {
		_, err := b.bot.SendMessage(tg.SendMessageRequest{ChatID: msg.From.ID, Text: "oops, you don't have any active package"})
		if err != nil {
			return err
		}
	}

	activePack := packages.ActivePackage
	reply := fmt.Sprintf(`
	Active Package:
	Traffic Limit: %s
	Download Traffic Usage: %s
	Upload Traffic Usage: %s
	Total Traffic Usage: %s
	Max Connections: %d
	Expire At: %s
	`, util.ToHumanReadableBytes(activePack.TrafficLimit), util.ToHumanReadableBytes(activePack.DownloadTrafficUsage),
		util.ToHumanReadableBytes(activePack.UploadTrafficUsage),
		util.ToHumanReadableBytes(activePack.DownloadTrafficUsage+activePack.UploadTrafficUsage), activePack.MaxConnections,
		activePack.ExpireAt.Format(TimeFormat))

	if len(packages.ReservedPackages) > 0 {
		reply += "Reserved Packages\n"

		for i, pack := range packages.ReservedPackages {
			reply += fmt.Sprintf(`
		# %d
		Traffic Limit: %s
		Max Connections: %d
		Expiration: %d Days
		`, i+1, util.ToHumanReadableBytes(pack.TrafficLimit), pack.MaxConnections, pack.ExpirationInDays)
		}
	}

	_, err = b.bot.SendMessage(tg.SendMessageRequest{ChatID: msg.From.ID, Text: reply})
	if err != nil {
		return err
	}

	return nil
}

func (b MainBot) GetActiveConnections(msg tg.Message) error {
	conns, err := b.connectionSvc.GetUserActiveConnections(context.Background(), msg.From.Username)
	if err != nil || len(conns) <= 0 {
		_, err := b.bot.SendMessage(tg.SendMessageRequest{ChatID: msg.From.ID, Text: "there is no active connections"})
		if err != nil {
			return err
		}
	}

	var reply string
	for i, conn := range conns {
		reply += fmt.Sprintf(
			`
			# %d
			IP: %s
			Location: %s
			UserAgent: %s
			Device: %s
			Download Traffic Usage: %s
			Upload Traffic Usage: %s
			Connected At: %s
			`, i+1, conn.RemoteIP, conn.Location, conn.UserAgent, conn.Hostname, util.ToHumanReadableBytes(conn.DownloadTrafficUsage),
			util.ToHumanReadableBytes(conn.UploadTrafficUsage), conn.ConnectedAt.Format(TimeFormat))
	}

	_, err = b.bot.SendMessage(tg.SendMessageRequest{ChatID: msg.From.ID, Text: reply})
	if err != nil {
		return err
	}

	return nil
}

func (b MainBot) ChangePassword(msg tg.Message) error {
	newPassword, err := b.userSvc.ChangePassword(context.Background(), msg.From.Username)
	if err != nil {
		_, err := b.bot.SendMessage(tg.SendMessageRequest{ChatID: msg.From.ID, Text: err.Error()})
		if err != nil {
			return err
		}
	}

	reply := fmt.Sprintf(
		`here is your username and your new password, this message will be deleted in an hour,
         Userame: %s
         Password: %s`, msg.From.Username, newPassword)
	resp, err := b.bot.SendMessage(tg.SendMessageRequest{ChatID: msg.From.ID, Text: reply, ProtectContent: true})
	if err != nil {
		return err
	}

	if len(resp) > 0 {
		time.AfterFunc(1*time.Hour, func() {
			if err := b.bot.DeleteMessage(msg.From.ID, resp[0].MessageID); err != nil {
				b.logger.Error(err.Error())
			}
		})
	}

	return nil
}

func (b MainBot) SendUnknownMessage(id int) error {
	_, err := b.bot.SendMessage(tg.SendMessageRequest{
		ChatID:         id,
		Text:           "huh? use /start if you don't know how to use me",
		ProtectContent: true})
	if err != nil {
		return err
	}

	return nil
}

func mapCommandsByName(commands []tg.BotCommand) map[string]tg.BotCommand {
	result := make(map[string]tg.BotCommand)

	for _, cmd := range commands {
		result[cmd.Command] = cmd
	}

	return result
}
