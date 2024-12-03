package app

import (
	"embed"
	wails "github.com/wailsapp/wails/v3/pkg/application"
	"go.uber.org/zap"
	"runtime"
	"time"
	"vpngui/internal/app/config"
	"vpngui/internal/app/log"
	"vpngui/internal/app/repository"
	"vpngui/internal/app/stats"
	xrayapi "vpngui/internal/app/xray-api"
	"vpngui/pkg/db"
	embedded "vpngui/pkg/embed"
	"vpngui/pkg/logger"
)

type App struct {
	wails *wails.App
}

func NewApp() *App {
	return &App{}
}

func (a *App) shutdown() {
	err := runXrayApi.KillOnClose()
	if err != nil {
		logger.Error("Failed to kill Xray API", zap.Error(err))
		return
	}
}

func (a *App) Hide() {
	a.wails.Hide()
}

func (a *App) Show() {
	a.wails.Show()
}

func (a *App) IsGOOSWindows() bool {
	return runtime.GOOS == "windows"
}

var (
	settingsRepo  *repository.SettingsRepository
	configRepo    *repository.ConfigRepository
	routesRepo    *repository.RoutesRepository
	runXrayApi    *xrayapi.RunXrayAPI
	routesXrayApi *xrayapi.RoutesXrayAPI
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
	settingsRepo = repository.NewSettings()
	configRepo = repository.NewConfig()
	routesRepo = repository.NewRoutes()
	runXrayApi = xrayapi.NewRun(configRepo)
	routesXrayApi = xrayapi.NewRoutes(runXrayApi, routesRepo)
	traffic = stats.NewTraffic(configRepo)
	capLog = log.New()

	go capLog.CaptureStdout()

	cfg = config.New()
	routesXrayApi.ActualizeConfig()

	app := NewApp()
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
	app.wails = wails.New(wails.Options{
		Name:        "VPN-GUI",
		Description: "VPN",
		Services: []wails.Service{
			wails.NewService(app),
			wails.NewService(cfg),
			wails.NewService(settingsRepo),
			wails.NewService(configRepo),
			wails.NewService(routesRepo),
			wails.NewService(runXrayApi),
			wails.NewService(routesXrayApi),
			wails.NewService(traffic),
			wails.NewService(capLog),
		},
		Assets: wails.AssetOptions{
			Handler: wails.AssetFileServerFS(assets),
		},
		OnShutdown: app.shutdown,
		Mac: wails.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
	})

	app.wails.NewWebviewWindowWithOptions(wails.WebviewWindowOptions{
		Title:         "VPN-GUI",
		Width:         750,
		Height:        500,
		DisableResize: true,
		Frameless:     runtime.GOOS == "windows",
		Mac: wails.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                wails.MacBackdropTranslucent,
			TitleBar:                wails.MacTitleBarHiddenInset,
		},
		Linux: wails.LinuxWindow{
			WindowIsTranslucent: false,
			WebviewGpuPolicy:    wails.WebviewGpuPolicyAlways,
		},
		Windows: wails.WindowsWindow{
			HiddenOnTaskbar: true,
		},
		BackgroundColour: wails.NewRGB(27, 38, 54),
		URL:              "/",
	})

	go func() {
		for {
			now := time.Now().Format(time.RFC1123)
			app.wails.EmitEvent("time", now)
			time.Sleep(time.Second)
		}
	}()
	systemTray(app)

	err := app.wails.Run()
	if err != nil {
		logger.Fatal("Failed to start application", zap.Error(err))
	}

	//HideWindowOnClose: true,
	//OnStartup:  app.startup,
	//OnShutdown: app.shutdown,

	return nil
}

func systemTray(app *App) {
	tray := app.wails.NewSystemTray()
	tray.SetLabel("VPN")

	menu := app.wails.NewMenu()

	menu.Add("Открыть окно").OnClick(func(_ *wails.Context) {
		app.wails.Show()
	})

	menu.AddSeparator()

	menu.Add("Включить VPN").OnClick(func(_ *wails.Context) {
		c, err := configRepo.GetConfig()
		if err != nil {
			logger.Error("Failed to get config", zap.Error(err))
			return
		}

		if c.ActiveVPN == false {
			err = runXrayApi.Run()
			if err != nil {
				logger.Error("Failed to run xray-core", zap.Error(err))
				return
			}
		}
	})
	menu.Add("Выключить VPN").OnClick(func(_ *wails.Context) {
		c, err := configRepo.GetConfig()
		if err != nil {
			logger.Error("Failed to get config", zap.Error(err))
			return
		}

		if c.ActiveVPN == true {
			err = runXrayApi.Run()
			if err != nil {
				logger.Error("Failed to run xray-core", zap.Error(err))
				return
			}
		}
	})
	menu.Add("Перезапустить VPN").OnClick(func(_ *wails.Context) {
		c, err := configRepo.GetConfig()
		if err != nil {
			logger.Error("Failed to get config", zap.Error(err))
			return
		}

		if c.ActiveVPN == true {
			err = runXrayApi.Kill()
			if err != nil {
				logger.Error("Failed to kill xray-core", zap.Error(err))
				return
			}
		}
		err = runXrayApi.Run()
		if err != nil {
			logger.Error("Failed to run xray-core", zap.Error(err))
			return
		}
	})

	menu.AddSeparator()

	menu.Add("Включить маршруты").OnClick(func(_ *wails.Context) {
		c, err := configRepo.GetConfig()
		if err != nil {
			logger.Error("Failed to get config", zap.Error(err))
			return
		}

		if c.DisableRoutes == true {
			err = routesXrayApi.EnableRoutes()
			if err != nil {
				logger.Error("Failed to enable routes", zap.Error(err))
				return
			}
		}
	})
	menu.Add("Отключить маршруты").OnClick(func(_ *wails.Context) {
		c, err := configRepo.GetConfig()
		if err != nil {
			logger.Error("Failed to get config", zap.Error(err))
			return
		}

		if c.DisableRoutes == false {
			err = routesXrayApi.DisableRoutes()
			if err != nil {
				logger.Error("Failed to enable routes", zap.Error(err))
				return
			}
		}
	})

	menu.AddSeparator()

	menu.Add("Режим работы").OnClick(func(_ *wails.Context) {
		app.Hide()
	})

	menu.AddSeparator()

	menu.Add("Выход").OnClick(func(_ *wails.Context) {
		app.wails.Quit()
	})

	tray.SetMenu(menu)
}
