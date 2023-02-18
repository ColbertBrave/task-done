package server

import (
	"cloud-disk/internal/auth"
	"cloud-disk/internal/log"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var options []Option

type Server struct {
	addr      string
	server    *http.Server
	ginEngine *gin.Engine
}

type Option func(engine *gin.Engine)

func NewServer(addr string) *Server {
	if addr == "" {
		return nil
	}

	gin.ForceConsoleColor()
	gin.SetMode(gin.ReleaseMode)

	ginEngine := gin.New()
	ginEngine.Use(gin.Recovery())
	ginEngine.Use(GetCostTimeOfRequest())
	ginEngine.Use(Authenticate())

	return &Server{
		addr: addr,
		server: &http.Server{
			Addr:    addr,
			Handler: ginEngine,
		},
		ginEngine: ginEngine,
	}
}

func GetCostTimeOfRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		costTime := time.Since(startTime)
		log.Info("%s|%s|cost time %d ms", c.Request.Method, c.Request.URL, costTime)
	}
}

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := auth.VerifyRequest(auth.Auth, c.Request)
		if err != nil {
			log.Error("verify request error:%s", err)
			return
		}
		c.Next()
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	err = s.server.Serve(listener)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Close() error {
	log.Info("the cloud disk is closed")
	return s.server.Close()
}

func (s *Server) InitRouter() {
	for _, option := range options {
		option(s.ginEngine)
	}
}
