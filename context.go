package goo

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

var (
	sig         = make(chan os.Signal)
	ctx, cancel = context.WithCancel(context.Background())
	Context     = ctx
)

func init() {
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)

	go func() {
		for si := range sig {
			switch si {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				cancel()
			case syscall.SIGUSR1:
			case syscall.SIGUSR2:
			default:
			}
		}
	}()
}
