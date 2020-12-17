package goo

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

var (
	sigs            = make(chan os.Signal)
	Context, cancel = context.WithCancel(context.Background())
)

func init() {
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)

	AsyncFunc(func() {
		for sig := range sigs {
			switch sig {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				cancel()
				return
			case syscall.SIGUSR1:
			case syscall.SIGUSR2:
			}
		}
	})

	AsyncFunc(func() {
		for {
			select {
			case <-Context.Done():
				os.Exit(0)
			}
		}
	})
}
