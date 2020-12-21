package goo

import (
	"fmt"
	"github.com/facebookgo/grace/gracenet"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type GRPCGraceful struct {
	nett string
	addr string
	s    *grpc.Server
	net  *gracenet.Net
}

func NewGRPCGraceful(nett, addr string, s *grpc.Server) *GRPCGraceful {
	return &GRPCGraceful{
		nett: nett,
		addr: addr,
		s:    s,
		net:  &gracenet.Net{},
	}
}

func (g *GRPCGraceful) Serve() error {
	lis, err := g.net.Listen(g.nett, g.addr)
	if err != nil {
		return err
	}

	errs := make(chan error)

	AsyncFunc(func() {
		errs <- g.s.Serve(lis)
	})
	AsyncFunc(g.killPPID)
	AsyncFunc(g.storePID)

	quit := g.handleSignal(errs)

	select {
	case err := <-errs:
		return err
	case <-quit:
		return nil
	}
}

func (g *GRPCGraceful) handleSignal(errs chan error) <-chan struct{} {
	quit := make(chan struct{})

	AsyncFunc(func() {
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)

		for sig := range ch {
			switch sig {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				signal.Stop(ch)
				g.s.GracefulStop()
				close(quit)
				return

			case syscall.SIGUSR1, syscall.SIGUSR2:
				if _, err := g.net.StartProcess(); err != nil {
					errs <- err
				}
			}
		}
	})

	return quit
}

func (g *GRPCGraceful) storePID() {
	pid := fmt.Sprintf("%d", os.Getpid())
	ioutil.WriteFile(".pid", []byte(pid), 0644)
	log.Println(fmt.Sprintf("server is running, address=%s, pid=%s", g.addr, pid))
}

func (g *GRPCGraceful) killPPID() {
	inherit := os.Getenv("LISTEN_FDS") != ""
	if !inherit {
		return
	}
	ppid := os.Getppid()
	if ppid == 1 {
		return
	}
	syscall.Kill(ppid, syscall.SIGTERM)
}
