package tun

import (
	"fmt"
	"vpngui/internal/app/command"
)

func setMacOSTun() error {
	err := clearMacOSTun()
	if err != nil {
		return err
	}

	commands := [][]string{
		{"sudo", "ifconfig", "utun100", "198.18.0.1", "198.18.0.1", "up"},
		{"sudo", "route", "add", "-net", "1.0.0.0/8", "198.18.0.1"},
		{"sudo", "route", "add", "-net", "2.0.0.0/7", "198.18.0.1"},
		{"sudo", "route", "add", "-net", "4.0.0.0/6", "198.18.0.1"},
		{"sudo", "route", "add", "-net", "8.0.0.0/5", "198.18.0.1"},
		{"sudo", "route", "add", "-net", "16.0.0.0/4", "198.18.0.1"},
		{"sudo", "route", "add", "-net", "32.0.0.0/3", "198.18.0.1"},
		{"sudo", "route", "add", "-net", "64.0.0.0/2", "198.18.0.1"},
		{"sudo", "route", "add", "-net", "128.0.0.0/1", "198.18.0.1"},
		{"sudo", "route", "add", "-net", "198.18.0.0/15", "198.18.0.1"},
	}

	return command.RunCommands(commands, false)
}

func clearMacOSTun() error {
	commands := [][]string{
		{"sudo", "ifconfig", "utun100", "198.18.0.1", "198.18.0.1", "down"},
		{"sudo", "route", "delete", "default"},
		{"sudo", "route", "delete", "1.0.0.0/8"},
		{"sudo", "route", "delete", "2.0.0.0/7"},
		{"sudo", "route", "delete", "4.0.0.0/6"},
		{"sudo", "route", "delete", "8.0.0.0/5"},
		{"sudo", "route", "delete", "16.0.0.0/4"},
		{"sudo", "route", "delete", "32.0.0.0/3"},
		{"sudo", "route", "delete", "64.0.0.0/2"},
		{"sudo", "route", "delete", "128.0.0.0/1"},
		{"sudo", "route", "delete", "198.18.0.0/15"},
		{"sudo", "route", "add", "default", DefaultGateway},
		{"sudo", "route", "add", "-net", fmt.Sprintf("%s/32", DefaultIP), "-interface", DefaultInterface},
	}

	return command.RunCommands(commands, true)
}

func setLinuxTun() error {
	commands := [][]string{
		{"ip", "route", "del", "default"},
		{"ip", "route", "add", "default", "via", "198.18.0.1", "dev", "tun0", "metric", "1"},
		{"ip", "route", "add", "default", "via", "172.17.0.1", "dev", DefaultInterface, "metric", "10"},
	}

	return command.RunCommands(commands, false)
}

func clearLinuxTun() error {
	commands := [][]string{
		{"ip", "route", "del", "default"},
	}

	return command.RunCommands(commands, false)
}

func setWindowsTun() error {
	err := clearWindowsTun()
	if err != nil {
		return err
	}

	commands := [][]string{
		{"netsh", "interface", "ipv4", "set", "address", "name=\"wintun\"", "source=static", "addr=192.168.123.1", "mask=255.255.255.0"},
		{"netsh", "interface", "ipv4", "set", "dnsservers", "name=\"wintun\"", "static", "address=8.8.8.8", "register=none", "validate=no"},
		{"netsh", "interface", "ipv4", "add", "route", "0.0.0.0/0", "\"wintun\"", "192.168.123.1", "metric=1"},
	}

	return command.RunCommands(commands, false)
}

func clearWindowsTun() error {
	commands := [][]string{
		{"netsh", "interface", "ipv4", "delete", "route", "0.0.0.0/0", "\"wintun\""},
	}

	return command.RunCommands(commands, true)
}
