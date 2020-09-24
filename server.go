package goo

import (
	"bytes"
	"fmt"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type iController interface {
	DoHandle(c *gin.Context) *Response
}

func Handler(controller iController) gin.HandlerFunc {
	return func(c *gin.Context) {
		rsp := controller.DoHandle(c)
		c.Set("response", rsp)
		if rsp == nil {
			return
		}
		c.JSON(200, rsp)
	}
}

type server struct {
	*gin.Engine
	noLogPaths map[string]interface{}
}

func NewServer() *server {
	s := &server{
		Engine: gin.New(),
		noLogPaths: map[string]interface{}{
			"/favicon.ico": nil,
		},
	}
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

func (s *server) SetNoLogPath(paths ...string) {
	for _, v := range paths {
		s.noLogPaths[v] = nil
	}
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

func (s *server) logger() gin.HandlerFunc {
	var requestId = 1000

	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		if _, ok := s.noLogPaths[c.Request.URL.Path]; ok {
			return
		}

		requestId++

		body := ""
		switch c.ContentType() {
		case "application/x-www-form-urlencoded", "application/json", "text/xml":
			buf, _ := ioutil.ReadAll(c.Request.Body)
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
			body = string(buf)
		}

		if body != "" {
			body = strings.ReplaceAll(body, "\\", "")
		}

		data := map[string]interface{}{
			"method":              c.Request.Method,
			"uri":                 c.Request.RequestURI,
			"body":                body,
			"authorization":       c.GetHeader("Authorization"),
			"x-request-id":        c.GetHeader("X-Request-Id"),
			"x-request-source":    c.GetHeader("X-Request-Source"),
			"x-request-timestamp": c.GetHeader("X-Request-Timestamp"),
			"x-request-sign":      c.GetHeader("X-Request-Sign"),
			"content-type":        c.ContentType(),
			"client-ip":           c.ClientIP(),
			"referer":             c.GetHeader("Referer"),
			"execution-time":      fmt.Sprintf("%dms", (time.Now().UnixNano()-start.UnixNano())/1e6),
		}
		Log.Debug(fmt.Sprintf("[api-request][%d]", requestId), data)

		rsp, ok := c.Get("response")
		if ok {
			Log.Debug(fmt.Sprintf("[api-response][%d]", requestId), rsp)
			if errMsg := rsp.(*Response).ErrMsg; len(errMsg) > 0 {
				Log.Error(fmt.Sprintf("[api-response-err][%d]", requestId), errMsg)
			}
		}
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
	noAccessPaths := map[string]interface{}{
		"/favicon.ico": nil,
	}

	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		if _, ok := noAccessPaths[c.Request.URL.Path]; ok {
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
