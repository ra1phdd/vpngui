package ui

//
//import (
//	"fyne.io/fyne/v2"
//	"fyne.io/fyne/v2/canvas"
//	"fyne.io/fyne/v2/container"
//	"fyne.io/fyne/v2/driver/desktop"
//	"fyne.io/fyne/v2/widget"
//	"image/color"
//	"time"
//	"vpngui/config"
//	xray_api "vpngui/internal/app/xray-api"
//)
//
//var bgColor = canvas.NewRectangle(color.RGBA{25, 25, 40, 255})
//
//func ApplyStartConfig(buttonEnableVPN, buttonDisableVPN *widget.Button, checkRoutes *widget.Check) {
//	if config.JSON.ActiveVPN {
//		buttonEnableVPN.Disable()
//		buttonDisableVPN.Enable()
//
//		go xray_api.Run()
//	} else {
//		buttonEnableVPN.Enable()
//		buttonDisableVPN.Disable()
//	}
//
//	if config.JSON.DisableRoutes {
//		checkRoutes.Checked = true
//	}
//}
//
//func SetupMainWindow(a fyne.App) {
//	var label *widget.Label
//	var buttonEnableVPN, buttonDisableVPN, buttonRestartVPN *widget.Button
//	var checkRoutes *widget.Check
//
//	w := a.NewWindow("VPN-GUI")
//	w.SetCloseIntercept(func() {
//		w.Hide()
//	})
//
//	label = widget.NewLabel("VPN выключен")
//	if config.JSON.ActiveVPN {
//		label.SetText("VPN включен")
//	}
//	centeredLabel := container.NewCenter(label)
//
//	buttonEnableVPN = widget.NewButton("Включить VPN", func() {
//		buttonEnableVPN.Disable()
//		buttonDisableVPN.Enable()
//
//		label.SetText("VPN включается...")
//		go xray_api.Run()
//		for !config.JSON.ActiveVPN {
//			time.Sleep(100 * time.Millisecond)
//		}
//		label.SetText("VPN включен")
//	})
//	buttonDisableVPN = widget.NewButton("Выключить VPN", func() {
//		buttonEnableVPN.Enable()
//		buttonDisableVPN.Disable()
//
//		label.SetText("VPN выключается...")
//		xray_api.Kill()
//		for config.JSON.ActiveVPN {
//			time.Sleep(100 * time.Millisecond)
//		}
//		label.SetText("VPN выключен")
//	})
//	buttonRestartVPN = widget.NewButton("Перезапустить VPN", func() {
//		label.SetText("VPN перезапускается...")
//		xray_api.Kill()
//		for config.JSON.ActiveVPN {
//			time.Sleep(100 * time.Millisecond)
//		}
//		go xray_api.Run()
//		for !config.JSON.ActiveVPN {
//			time.Sleep(100 * time.Millisecond)
//		}
//		label.SetText("VPN включен")
//	})
//	buttonRoutes := widget.NewButton("Маршруты", func() {
//		rw := a.NewWindow("Маршруты")
//		rw.SetOnClosed(func() {
//			xray_api.Kill()
//			for config.JSON.ActiveVPN {
//				time.Sleep(100 * time.Millisecond)
//			}
//			go xray_api.Run()
//			for !config.JSON.ActiveVPN {
//				time.Sleep(100 * time.Millisecond)
//			}
//		})
//
//		SetupRoutesWindow(rw)
//	})
//	checkRoutes = widget.NewCheck("Отключить маршруты (проксировать всё)", func(value bool) {
//		if value {
//			err := xray_api.DisableRoutes()
//			if err != nil {
//				return
//			}
//		} else {
//			err := xray_api.EnableRoutes()
//			if err != nil {
//				return
//			}
//		}
//	})
//
//	ApplyStartConfig(buttonEnableVPN, buttonDisableVPN, checkRoutes)
//
//	if desk, ok := a.(desktop.App); ok {
//		m := fyne.NewMenu("VPN-GUI",
//			fyne.NewMenuItem("Показать окно", func() {
//				w.Show()
//			}),
//			fyne.NewMenuItem("Скрыть окно", func() {
//				w.Hide()
//			}),
//			fyne.NewMenuItemSeparator(),
//			fyne.NewMenuItem("Включить VPN", func() {
//				if config.JSON.ActiveVPN {
//					return
//				}
//
//				buttonEnableVPN.Disable()
//				buttonDisableVPN.Enable()
//
//				label.SetText("VPN включается...")
//				go xray_api.Run()
//				for !config.JSON.ActiveVPN {
//					time.Sleep(100 * time.Millisecond)
//				}
//				label.SetText("VPN включен")
//			}),
//			fyne.NewMenuItem("Выключить VPN", func() {
//				if !config.JSON.ActiveVPN {
//					return
//				}
//
//				buttonEnableVPN.Enable()
//				buttonDisableVPN.Disable()
//
//				label.SetText("VPN выключается...")
//				xray_api.Kill()
//				for config.JSON.ActiveVPN {
//					time.Sleep(100 * time.Millisecond)
//				}
//				label.SetText("VPN выключен")
//			}),
//			fyne.NewMenuItem("Перезапустить VPN", func() {
//				label.SetText("VPN перезапускается...")
//				xray_api.Kill()
//				for config.JSON.ActiveVPN {
//					time.Sleep(100 * time.Millisecond)
//				}
//				go xray_api.Run()
//				for !config.JSON.ActiveVPN {
//					time.Sleep(100 * time.Millisecond)
//				}
//				label.SetText("VPN включен")
//			}),
//			fyne.NewMenuItemSeparator(),
//			fyne.NewMenuItem("Включить маршруты", func() {
//				if !config.JSON.DisableRoutes {
//					return
//				}
//
//				err := xray_api.EnableRoutes()
//				if err != nil {
//					return
//				}
//
//				checkRoutes.SetChecked(false)
//			}),
//			fyne.NewMenuItem("Отключить маршруты", func() {
//				if config.JSON.DisableRoutes {
//					return
//				}
//
//				err := xray_api.DisableRoutes()
//				if err != nil {
//					return
//				}
//
//				checkRoutes.Checked = true
//				checkRoutes.Refresh()
//			}),
//		)
//		desk.SetSystemTrayMenu(m)
//	}
//
//	buttonEnableVPN.ExtendBaseWidget(buttonEnableVPN)
//	buttonEnableVPN.Importance = widget.HighImportance
//
//	buttonDisableVPN.ExtendBaseWidget(buttonDisableVPN)
//	buttonDisableVPN.Importance = widget.HighImportance
//
//	hbox := container.NewGridWithRows(1, buttonEnableVPN, buttonDisableVPN)
//	content := container.NewVBox(
//		centeredLabel,
//		hbox,
//		buttonRestartVPN,
//		widget.NewLabel(""),
//		buttonRoutes,
//		container.NewCenter(checkRoutes),
//	)
//	w.SetContent(content)
//
//	w.Resize(fyne.NewSize(400, 550))
//	w.ShowAndRun()
//}
