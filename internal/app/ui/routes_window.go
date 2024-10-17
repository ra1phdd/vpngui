package ui

//
//import (
//	"fyne.io/fyne/v2"
//	"fyne.io/fyne/v2/container"
//	"fyne.io/fyne/v2/widget"
//	"vpngui/config"
//	"vpngui/internal/app/xray-api"
//)
//
//func ApplyStartRoutes(radioEnableBlackList, radioEnableWhiteList *widget.Check) {
//	if config.JSON.EnableBlackList {
//		radioEnableBlackList.SetChecked(true)
//	} else {
//		radioEnableWhiteList.SetChecked(true)
//	}
//}
//
//func SetupRoutesWindow(w fyne.Window) {
//	var entryDomainBlackList, entryIPBlackList, entryPortBlackList, entryDomainWhiteList, entryIPWhiteList, entryPortWhiteList *fyne.Container
//	var radioEnableBlackList, radioEnableWhiteList *widget.Check
//
//	radioEnableBlackList = widget.NewCheck("Чёрные списки", func(value bool) {
//		if value {
//			radioEnableWhiteList.SetChecked(false)
//
//			config.JSON.EnableBlackList = true
//			config.JSON.EnableWhiteList = false
//
//			err := xray_api.SwapOutbounds(&config.Xray.Outbounds, "direct", "proxy")
//			if err != nil {
//				return
//			}
//		} else {
//			radioEnableWhiteList.SetChecked(true)
//
//			config.JSON.EnableBlackList = false
//			config.JSON.EnableWhiteList = true
//
//			err := xray_api.SwapOutbounds(&config.Xray.Outbounds, "proxy", "direct")
//			if err != nil {
//				return
//			}
//		}
//		err := config.SaveConfig()
//		if err != nil {
//			return
//		}
//	})
//	radioEnableWhiteList = widget.NewCheck("Белые списки", func(value bool) {
//		if value {
//			radioEnableBlackList.SetChecked(false)
//
//			config.JSON.EnableBlackList = false
//			config.JSON.EnableWhiteList = true
//
//			err := xray_api.SwapOutbounds(&config.Xray.Outbounds, "proxy", "direct")
//			if err != nil {
//				return
//			}
//		} else {
//			radioEnableBlackList.SetChecked(true)
//
//			config.JSON.EnableBlackList = true
//			config.JSON.EnableWhiteList = false
//
//			err := xray_api.SwapOutbounds(&config.Xray.Outbounds, "direct", "proxy")
//			if err != nil {
//				return
//			}
//		}
//		err := config.SaveConfig()
//		if err != nil {
//			return
//		}
//	})
//
//	listDomainBlackList := createEntry("Список доменов", xray_api.GetDomain("proxy"), true)
//	listIPBlackList := createEntry("Список IP-адресов", xray_api.GetIP("proxy"), true)
//	listPortBlackList := createEntry("Список портов", xray_api.GetPort("proxy"), false)
//	listDomainWhiteList := createEntry("Список доменов", xray_api.GetDomain("direct"), true)
//	listIPWhiteList := createEntry("Список IP-адресов", xray_api.GetIP("direct"), true)
//	listPortWhiteList := createEntry("Список портов", xray_api.GetPort("direct"), false)
//
//	entryDomainBlackList = createEntryWithButton("Введите домен...",
//		func(text string) {
//			xray_api.AddDomain("proxy", text)
//			listDomainBlackList.SetText(xray_api.GetDomain("proxy"))
//		},
//		func(text string) {
//			xray_api.DelDomain("proxy", text)
//			listDomainBlackList.SetText(xray_api.GetDomain("proxy"))
//		},
//	)
//	entryIPBlackList = createEntryWithButton("Введите IP-адрес...",
//		func(text string) {
//			xray_api.AddPort("proxy", text)
//			listIPBlackList.SetText(xray_api.GetIP("proxy"))
//		},
//		func(text string) {
//			xray_api.DelPort("proxy", text)
//			listIPBlackList.SetText(xray_api.GetIP("proxy"))
//		},
//	)
//	entryPortBlackList = createEntryWithButton("Введите порт...",
//		func(text string) {
//			xray_api.AddPort("proxy", text)
//			listPortBlackList.SetText(xray_api.GetPort("proxy"))
//		},
//		func(text string) {
//			xray_api.DelPort("proxy", text)
//			listPortBlackList.SetText(xray_api.GetPort("proxy"))
//		},
//	)
//	entryDomainWhiteList = createEntryWithButton("Введите домен...",
//		func(text string) {
//			xray_api.AddDomain("direct", text)
//			listDomainWhiteList.SetText(xray_api.GetDomain("direct"))
//		},
//		func(text string) {
//			xray_api.DelDomain("direct", text)
//			listDomainWhiteList.SetText(xray_api.GetDomain("direct"))
//		},
//	)
//	entryIPWhiteList = createEntryWithButton("Введите IP-адрес...",
//		func(text string) {
//			xray_api.AddIP("direct", text)
//			listIPWhiteList.SetText(xray_api.GetIP("direct"))
//		},
//		func(text string) {
//			xray_api.DelIP("direct", text)
//			listIPWhiteList.SetText(xray_api.GetIP("direct"))
//		},
//	)
//	entryPortWhiteList = createEntryWithButton("Введите порт...",
//		func(text string) {
//			xray_api.AddPort("direct", text)
//			listPortWhiteList.SetText(xray_api.GetPort("direct"))
//		},
//		func(text string) {
//			xray_api.DelPort("direct", text)
//			listPortWhiteList.SetText(xray_api.GetPort("direct"))
//		},
//	)
//
//	firstContainer := container.NewVBox(
//		container.NewCenter(radioEnableBlackList),
//		container.NewVBox(listDomainBlackList),
//		container.NewGridWithRows(1, entryDomainBlackList),
//		widget.NewLabel(""),
//		container.NewVBox(listIPBlackList),
//		container.NewGridWithRows(1, entryIPBlackList),
//		widget.NewLabel(""),
//		container.NewVBox(listPortBlackList),
//		container.NewGridWithRows(1, entryPortBlackList),
//	)
//
//	secondContainer := container.NewVBox(
//		container.NewCenter(radioEnableWhiteList),
//		container.NewVBox(listDomainWhiteList),
//		container.NewGridWithRows(1, entryDomainWhiteList),
//		widget.NewLabel(""),
//		container.NewVBox(listIPWhiteList),
//		container.NewGridWithRows(1, entryIPWhiteList),
//		widget.NewLabel(""),
//		container.NewVBox(listPortWhiteList),
//		container.NewGridWithRows(1, entryPortWhiteList),
//	)
//
//	ApplyStartRoutes(radioEnableBlackList, radioEnableWhiteList)
//	horizontalContainer := container.NewGridWithRows(1,
//		firstContainer,
//		secondContainer,
//	)
//
//	w.SetContent(horizontalContainer)
//	w.Resize(fyne.NewSize(600, 0))
//	w.Show()
//}
//
//func createEntry(placeholder, text string, multiline bool) *widget.Entry {
//	entry := widget.NewEntry()
//	entry.MultiLine = multiline
//	//entry.Disable()
//	entry.SetPlaceHolder(placeholder)
//	entry.SetText(text)
//
//	return entry
//}
//
//func createEntryWithButton(placeholder string, onButtonAddClick func(text string), onButtonDelClick func(text string)) *fyne.Container {
//	entry := widget.NewEntry()
//	entry.SetPlaceHolder(placeholder)
//
//	buttonAdd := widget.NewButton(" + ", func() {
//		text := entry.Text
//		onButtonAddClick(text)
//		entry.SetText("")
//	})
//
//	buttonDel := widget.NewButton(" - ", func() {
//		text := entry.Text
//		onButtonDelClick(text)
//		entry.SetText("")
//	})
//
//	buttonsContainer := container.NewHBox(buttonAdd, buttonDel)
//
//	return container.NewBorder(nil, nil, nil, buttonsContainer, entry)
//}
