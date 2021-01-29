package goo

import (
	"fmt"
	"github.com/facebookgo/grace/gracenet"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

type TCPGraceful struct {
	addr    string
	net     *gracenet.Net
	handler func(net.Conn)
	wg      sync.WaitGroup
}

func NewTCPGraceful(addr string, handler func(net.Conn)) *TCPGraceful {
	return &TCPGraceful{addr: addr, handler: handler, net: &gracenet.Net{}}
}

func (g *TCPGraceful) Serve() {
	addr, err := net.ResolveTCPAddr("tcp", g.addr)
	if err != nil {
		log.Fatal(err.Error())
	}

	l, err := g.net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err.Error())
	}

	quit := make(chan struct{})

	AsyncFunc(g.killPPID)
	AsyncFunc(g.storePID)
	AsyncFunc(g.handleSignal(l, quit))

	for {
		conn, err := l.Accept()
		if err != nil {
			Log.Error(err.Error())
			if strings.Contains(err.Error(), "use of closed network connection") {
				break
			}
			continue
		}
		g.wg.Add(1)
		AsyncFunc(func() {
			defer g.wg.Done()
			defer conn.Close()
			go g.handler(conn)
			<-quit
		})
	}

	g.wg.Wait()
}

func (g *TCPGraceful) handleSignal(l *net.TCPListener, quit chan struct{}) func() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)

	return func() {
		for sig := range ch {
			switch sig {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				signal.Stop(ch)
				l.Close()
				close(quit)
				return
			case syscall.SIGUSR1, syscall.SIGUSR2:
				if _, err := g.net.StartProcess(); err != nil {
					Log.Error(err.Error())
				}
			}
		}
	}
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
