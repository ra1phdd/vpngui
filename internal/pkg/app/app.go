package app

import (
	"context"
	"embed"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"go.uber.org/zap"
	"time"
	"vpngui/internal/app/config"
	"vpngui/internal/app/log"
	"vpngui/internal/app/repository"
	"vpngui/internal/app/stats"
	xray_api "vpngui/internal/app/xray-api"
	"vpngui/pkg/db"
	embedded "vpngui/pkg/embed"
	"vpngui/pkg/logger"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) shutdown(_ context.Context) {
	err := runXrayApi.KillOnClose()
	if err != nil {
		logger.Error("Failed to kill Xray API", zap.Error(err))
		return
	}
}

var (
	settingsRepo  *repository.SettingsRepository
	configRepo    *repository.ConfigRepository
	routesRepo    *repository.RoutesRepository
	runXrayApi    *xray_api.RunXrayAPI
	routesXrayApi *xray_api.RoutesXrayAPI
	traffic       *stats.Traffic
	capLog        *log.Log
	cfg           *config.Config
)

func New(assets embed.FS) error {
	logger.Init()

	if err := embedded.Init(); err != nil {
		return err
	}

	if err := db.Init("db/vpngui.db"); err != nil {
		return err
	}

	app := setupApplication()

	settings, err := settingsRepo.GetSettings()
	if err != nil {
		return err
	}

	logger.SetLogLevel(settings.LoggerLevel)
	go startTrafficCapture(settings.StatsUpdateInterval)

	return runWailsApp(assets, app)
}

func setupApplication() *App {
	app := NewApp()
	settingsRepo = repository.NewSettings()
	configRepo = repository.NewConfig()
	routesRepo = repository.NewRoutes()
	runXrayApi = xray_api.NewRun(configRepo)
	routesXrayApi = xray_api.NewRoutes(runXrayApi, routesRepo)
	traffic = stats.NewTraffic()
	capLog = log.New()

	go capLog.CaptureStdout()

	cfg = config.New()
	routesXrayApi.ActualizeConfig()

	return app
}

func startTrafficCapture(interval int) {
	logger.Info("Starting traffic capture")
	for {
		traffic.CaptureTraffic()
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func runWailsApp(assets embed.FS, app *App) error {
	err := wails.Run(&options.App{
		Title:             "VPN-GUI",
		Width:             750,
		Height:            500,
		DisableResize:     true,
		HideWindowOnClose: true,
		Frameless:         false,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:  app.startup,
		OnShutdown: app.shutdown,
		Bind: []interface{}{
			app,
			settingsRepo,
			configRepo,
			routesRepo,
			runXrayApi,
			routesXrayApi,
			cfg,
			capLog,
			traffic,
		},
		Mac: &mac.Options{
			TitleBar: mac.TitleBarHiddenInset(),
		},
		Linux: &linux.Options{
			WindowIsTranslucent: false,
			WebviewGpuPolicy:    linux.WebviewGpuPolicyAlways,
		},
	})
	if err != nil {
		logger.Fatal("Failed to start app", zap.Error(err))
	}

	return nil
}
