package app

import (
	"embed"
	wails "github.com/wailsapp/wails/v3/pkg/application"
	"go.uber.org/zap"
	"runtime"
	"time"
	"vpngui/internal/app/config"
	"vpngui/internal/app/log"
	"vpngui/internal/app/network"
	"vpngui/internal/app/repository"
	"vpngui/internal/app/runner"
	"vpngui/internal/app/stats"
	"vpngui/internal/app/transport"
	"vpngui/internal/app/transport/proxy"
	"vpngui/internal/app/transport/tun"
	"vpngui/internal/app/xray-core"
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
	err := runXrayCore.Kill(false)
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
	nw             *network.Network
	settingsRepo   *repository.SettingsRepository
	configRepo     *repository.ConfigRepository
	routesRepo     *repository.RoutesRepository
	runnerCmd      *runner.Command
	runnerProcess  *runner.Process
	transportProxy *proxy.Proxy
	transportTun   *tun.Tun
	transports     *transport.Transport
	runXrayCore    *xray_core.RunXrayCore
	routesXrayCore *xray_core.RoutesXrayCore
	traffic        *stats.Traffic
	capLog         *log.Log
	cfg            *config.Config
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
	nw = network.New()
	err := nw.Init()
	if err != nil {
		logger.Fatal("Failed to init network", zap.Error(err))
	}

	settingsRepo = repository.NewSettings()
	configRepo = repository.NewConfig()
	routesRepo = repository.NewRoutes()

	runnerCmd = runner.NewCmd()
	runnerProcess = runner.NewProcess()

	transportProxy = proxy.New(runnerCmd)
	transportTun = tun.New(runnerCmd, runnerProcess)
	transports = transport.New(configRepo, transportProxy, transportTun)

	traffic = stats.NewTraffic(configRepo, transports)
	capLog = log.New()
	go capLog.CaptureStdout()

	runXrayCore = xray_core.NewRun(configRepo, runnerProcess, transports)
	routesXrayCore = xray_core.NewRoutes(runXrayCore, routesRepo)

	cfg = config.New()
	routesXrayCore.ActualizeConfig()

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
			wails.NewService(runXrayCore),
			wails.NewService(routesXrayCore),
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
			err = runXrayCore.Run()
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
			err = runXrayCore.Run()
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
			err = runXrayCore.Kill(true)
			if err != nil {
				logger.Error("Failed to kill xray-core", zap.Error(err))
				return
			}
		}
		err = runXrayCore.Run()
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
			err = routesXrayCore.EnableRoutes()
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
			err = routesXrayCore.DisableRoutes()
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
