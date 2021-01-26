package goo

import (
	"fmt"
	"github.com/facebookgo/grace/gracenet"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type TCPGraceful struct {
	addr    string
	net     *gracenet.Net
	handler func(net.Conn) func()
}

func NewTCPGraceful(addr string, handler func(net.Conn) func()) *TCPGraceful {
	return &TCPGraceful{addr: addr, handler: handler, net: &gracenet.Net{}}
}

func (g *TCPGraceful) Serve() error {
	l, err := g.net.Listen("tcp", g.addr)
	if err != nil {
		return err
	}

	AsyncFunc(g.killPPID)
	AsyncFunc(g.storePID)

	errs := make(chan error)
	AsyncFunc(func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				errs <- err
				break
			}
			AsyncFunc(g.handler(conn))
		}
	})

	quit := g.handleSignal(l, errs)
	select {
	case err := <-errs:
		return err
	case <-quit:
		return nil
	}
	return nil
}

func (g *TCPGraceful) handleSignal(l net.Listener, errs chan error) <-chan struct{} {
	quit := make(chan struct{})
	AsyncFunc(func() {
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
		for sig := range ch {
			switch sig {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				signal.Stop(ch)
				l.Close()
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

func (g *TCPGraceful) storePID() {
	pid := fmt.Sprintf("%d", os.Getpid())
	ioutil.WriteFile(".pid", []byte(pid), 0644)
	log.Println(fmt.Sprintf("server is running, address=%s, pid=%s", g.addr, pid))
}

func (g *TCPGraceful) killPPID() {
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
