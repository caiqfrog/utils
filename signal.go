package utils

import (
	"os"
	"os/signal"
	"syscall"
)

// InitSignal register signals handler.
func Signal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL:
			return
		case syscall.SIGHUP:
			reload()
		default:
			return
		}
	}
}

func reload() {
	// TODO reload
}
