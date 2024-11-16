//go:build !windows

package proxy

func notifySettingsChange() error {
	return nil
}
