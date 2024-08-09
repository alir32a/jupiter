package main

import (
	"context"
	"github.com/alir32a/jupiter/config"
	"github.com/alir32a/jupiter/database"
	"github.com/alir32a/jupiter/internal/bot"
	"github.com/alir32a/jupiter/internal/handler"
	"github.com/alir32a/jupiter/internal/model"
	"github.com/alir32a/jupiter/internal/repository"
	"github.com/alir32a/jupiter/internal/service"
	"github.com/alir32a/jupiter/pkg/jwt"
	"github.com/alir32a/jupiter/pkg/ocserv"
	clog "github.com/charmbracelet/log"
	ejwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	logger := clog.NewWithOptions(os.Stdout, clog.Options{
		ReportTimestamp: true,
		ReportCaller:    true,
	})

	if err := ocserv.CheckInstallation(context.Background()); err != nil {
		logger.Fatal(err)
	}

	cfg, err := config.GetConfig()
	if err != nil {
		logger.Fatal(err)
	}

	if err := runManager(cfg, logger); err != nil {
		logger.Fatal(err)
	}
}

func runManager(cfg *config.Config, logger *clog.Logger) error {
	ocservClient := ocserv.NewClient(cfg.OCCTL.PasswordFile)

	db := setupDB(cfg, logger)

	userRepo := repository.NewUserRepository(db)
	connectionRepo := repository.NewConnectionRepository(db)
	packageRepo := repository.NewPackageRepository(db)
	adminRepo := repository.NewAdminRepository(db)

	userSvc := service.NewUserService(cfg, logger, ocservClient, userRepo, packageRepo)
	connectionSvc := service.NewConnectionService(logger, ocservClient, connectionRepo, packageRepo, userRepo)
	packageSvc := service.NewPackageService(logger, packageRepo, userRepo)
	adminSvc := service.NewAdminService(adminRepo, logger)

	server := handler.NewHTTPServer(cfg.HTTPServerConfig, logger)

	corsCfg := middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
	}
	if cfg.HTTPServerConfig.ENV == "debug" {
		// set for running web locally without cors and cookies issues, MUST NOT SET IN PRODUCTION ENV.
		corsCfg.UnsafeWildcardOriginWithAllowCredentials = true
	}

	server.Use(middleware.CORSWithConfig(corsCfg))

	noAuth := server.Group("/api/v1")
	auth := server.Group("/api/v1")
	auth.Use(ejwt.WithConfig(ejwt.Config{
		ErrorHandler: func(ctx echo.Context, err error) error {
			ctx.JSON(http.StatusUnauthorized, "invalid token")

			return err
		},
		SigningKey:  []byte(cfg.HTTPServerConfig.AccessTokenSecret),
		TokenLookup: "cookie:adminAccess",
		ParseTokenFunc: func(ctx echo.Context, token string) (interface{}, error) {
			return jwt.ParseToken(token, cfg.HTTPServerConfig.AccessTokenSecret)
		},
	}))

	healthCtrl := handler.NewHealthCheckHandler(db)
	healthCtrl.SetRoutes(server.Group(""))

	adminCtrl := handler.NewAdminHandler(adminSvc, cfg.HTTPServerConfig, logger)
	adminCtrl.SetNoAuthRoutes(noAuth)
	adminCtrl.SetRoutes(auth)

	connectionsCtrl := handler.NewConnectionHandler(connectionSvc, logger)
	connectionsCtrl.SetRoutes(auth)

	packagesCtrl := handler.NewPackageHandler(packageSvc, logger)
	packagesCtrl.SetRoutes(auth)

	usersCtrl := handler.NewUserHandler(userSvc, logger)
	usersCtrl.SetRoutes(auth)

	go func() {
		if err := server.Run(); err != nil {
			logger.Fatal(err)
		}
	}()

	mainBot := bot.NewMainBot(cfg.MainBot, logger, userSvc, connectionSvc, packageSvc)

	go func() {
		if err := mainBot.Run(); err != nil {
			logger.Fatal(err.Error())
		}
	}()

	ticker := time.Tick(cfg.Manager.UpdateInterval * time.Second)

	var failureCount int
	for _ = range ticker {
		if failureCount == cfg.Manager.MaxFailures {
			if cfg.Manager.ShutdownOnMaxFailure {
				if err := ocservClient.ShutdownServer(context.Background()); err != nil {
					logger.Fatal(err)
				}
			}
		}

		if err := manageConnections(ocservClient, connectionSvc, cfg.Manager.UpdateTimeout); err != nil {
			logger.Error(err.Error())

			failureCount++
		}
	}

	return nil
}

func setupDB(cfg *config.Config, logger *clog.Logger) *gorm.DB {
	db, err := database.GetDatabaseConnection(cfg.DB)
	if err != nil {
		logger.Fatal(err)
	}

	sqlDb, err := db.DB()
	if err != nil {
		logger.Fatal(err)
	}

	if err := goose.Up(sqlDb, cfg.DB.MigrationPath); err != nil {
		logger.Fatal(err)
	}

	return db
}

func manageConnections(ocservClient *ocserv.Client, connectionSvc *service.ConnectionService, timeout time.Duration) error {
	ctx, cancelFn := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancelFn()

	connections, err := ocservClient.GetConnections(ctx)
	if err != nil {
		return err
	}

	modelConns, updatedAt := toModelConnectionEntities(connections)

	err = connectionSvc.UpsertConnections(ctx, model.UpsertConnectionsRequest{Connections: modelConns})
	if err != nil {
		return err
	}

	return connectionSvc.ManageActiveConnections(ctx, updatedAt)
}

func toModelConnectionEntity(req ocserv.ConnectionEntity) model.ConnectionEntity {
	rx, _ := strconv.Atoi(req.RX)
	tx, _ := strconv.Atoi(req.TX)

	connectedAt, err := time.Parse("2006-01-02 15:04", req.ConnectedAt)
	if err != nil {
		clog.Error(err)
	}

	return model.ConnectionEntity{
		ExternalID:           strconv.Itoa(req.ID),
		Username:             req.Username,
		RemoteIP:             req.RemoteIP,
		Location:             req.Location,
		UserAgent:            req.UserAgent,
		Hostname:             req.Hostname,
		DownloadTrafficUsage: tx,
		UploadTrafficUsage:   rx,
		ConnectedAt:          connectedAt,
		UpdatedAt:            time.Now(),
	}
}

func toModelConnectionEntities(req []ocserv.ConnectionEntity) ([]model.ConnectionEntity, time.Time) {
	result := make([]model.ConnectionEntity, 0, len(req))

	if len(req) <= 0 {
		return nil, time.Time{}
	}

	for _, connection := range req {
		result = append(result, toModelConnectionEntity(connection))
	}

	return result, result[0].UpdatedAt
}
