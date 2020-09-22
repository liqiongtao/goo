package goo

import (
	"fmt"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
	"strings"
)

type iController interface {
	DoHandle(c *gin.Context) *Response
}

func Handler(controller iController) gin.HandlerFunc {
	return func(c *gin.Context) {
		rsp := controller.DoHandle(c)
		if rsp == nil {
			return
		}
		c.JSON(200, rsp)
	}
}

type server struct {
	*gin.Engine
}

func NewServer() *server {
	s := new(server)
	s.Engine = gin.New()
	s.Use(s.cors(), s.noAccess(), s.logger(), s.recovery())
	s.NoRoute(s.noRoute())
	return s
}

func (s *server) Run(addr string) {
	pid := fmt.Sprintf("%d", os.Getpid())
	if err := ioutil.WriteFile(".pid", []byte(pid), 0755); err != nil {
		panic(err.Error())
	}
	endless.NewServer(addr, s.Engine).ListenAndServe()
}

func (*server) cors() gin.HandlerFunc {
	allowHeaders := []string{
		"Content-Type", "Content-Length",
		"Accept", "Referer", "User-Agent, ",
		"Authorization",
		"X-Requested-Id", "X-Request-Timestamp", "X-Request-Sign",
		"X-Request-AppId", "X-Request-Source", "X-Request-Token",
	}

	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", strings.Join(allowHeaders, ","))
		c.Next()
	}
}

func (*server) logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func (s *server) recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.AbortWithStatusJSON(200, Error(500, fmt.Sprint(err)))
			}
		}()
		c.Next()
	}
}

func (*server) noAccess() gin.HandlerFunc {
	noAccessPaths := []string{
		"/favicon.ico",
	}

	noAccessPathsMap := map[string]struct{}{}
	for _, i := range noAccessPaths {
		noAccessPathsMap[i] = struct{}{}
	}

	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		if _, ok := noAccessPathsMap[c.Request.URL.Path]; ok {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}

func (*server) noRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatusJSON(200, Error(404, "Page Not Found"))
	}
}
