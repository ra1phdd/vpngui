//go:build windows
// +build windows

package proxy

import (
	"errors"
	"go.uber.org/zap"
	"syscall"
	"vpngui/pkg/logger"
)

func notifySettingsChange() error {
	dll, err := syscall.LoadDLL("wininet.dll")
	if err != nil {
		logger.Error("Failed to load wininet.dll", zap.Error(err))
		return err
	}
	defer dll.Release()

	proc, err := dll.FindProc("InternetSetOptionW")
	if err != nil {
		logger.Error("Failed to find InternetSetOptionW", zap.Error(err))
		return err
	}

	const InternetOptionSettingsChanged = 39
	const InternetOptionRefresh = 37

	if _, _, err := proc.Call(0, InternetOptionSettingsChanged, 0, 0); err != nil && !errors.Is(err, syscall.Errno(0)) {
		logger.Error("Failed to call InternetSetOptionW (SETTINGS_CHANGED)")
		return err
	}

	if _, _, err := proc.Call(0, InternetOptionRefresh, 0, 0); err != nil && !errors.Is(err, syscall.Errno(0)) {
		logger.Error("Failed to call InternetSetOptionW (REFRESH)")
		return err
	}

	return nil
}
