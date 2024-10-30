package app

import (
	"context"
	"embed"
	"fmt"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"vpngui/config"
	"vpngui/internal/app/log"
	"vpngui/internal/app/repository"
	xray_api "vpngui/internal/app/xray-api"
	"vpngui/pkg/db"
	embedded "vpngui/pkg/embed"
	"vpngui/pkg/logger"
)

func New(assets embed.FS) error {
	err := config.LoadConfig()
	if err != nil {
		return err
	}

	capLog := log.New()
	go capLog.CaptureStdout()

	logger.Init("info")

	err = embedded.Init()
	if err != nil {
		return err
	}

	err = db.Init("db/vpngui.db")
	if err != nil {
		return err
	}

	app := NewApp()
	cfg := &config.Config{}

	configRepository := repository.NewConfig()
	routesRepository := repository.NewRoutes()
	runXrayApi := xray_api.NewRun(configRepository)
	routesXrayApi := xray_api.NewRoutes(runXrayApi, routesRepository)

	err = wails.Run(&options.App{
		Title:             "VPN-GUI",
		Width:             750,
		Height:            500,
		DisableResize:     true,
		HideWindowOnClose: true,
		Frameless:         false,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
			configRepository,
			routesRepository,
			runXrayApi,
			routesXrayApi,
			cfg,
			capLog,
		},
		Mac: &mac.Options{
			TitleBar: mac.TitleBarHiddenInset(),
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}

	return nil
}

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
