package server

import (
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"qa_spider/config"
	"qa_spider/pkg"
	"qa_spider/pkg/services/queryQA"
	"qa_spider/server/content"
	"time"
)

type HTTPServer struct {
	config   *config.Config
	log      *zap.Logger
	engine   *gin.Engine
	ctn      map[interface{}]*content.Content
	internal map[string]pkg.Internal
}

func (s *HTTPServer) Run() error {
	for _, v := range s.internal {
		go func() {
			if err := v.Run(); err != nil {
				s.log.Error("initializing internal services", zap.String("service name", v.GetName()))
			}
		}()
	}
	if err := s.engine.Run(s.config.Server.Port); err != nil {
		return err
	}
	return nil
}

func (s *HTTPServer) Stop() {
	os.Exit(1)
}

func InitHTTPServer(config *config.Config, logger *zap.Logger, internal ...pkg.Internal) Server {
	//set mode
	gin.SetMode(gin.DebugMode)

	s := &HTTPServer{
		config:   config,
		log:      logger,
		engine:   gin.New(),
		ctn:      make(map[interface{}]*content.Content),
		internal: make(map[string]pkg.Internal),
	}
	//init content services

	//use zap as logger
	s.engine.Use(ginzap.Ginzap(s.log, time.RFC3339, true))

	if config.Server.AllowCors {
		logger.Info("Server allow cors enabled")
		s.engine.Use(Cors())
	} else {
		logger.Info("Server allow cors disabled")
	}

	s.regInternal(internal...)
	//init internal dependencies
	s.initContent()
	//init handlers
	s.regHandlers()

	//allow cors

	return s
}

func (s *HTTPServer) initContent() {
	s.ctn["qa_spider"] = content.InitContent(s.config, s.log, s.internal["QASPIDER"])

}

//router initialize
func (s *HTTPServer) regHandlers() {
	s.engine.GET("v1/spider/find", queryQA.QueryQA(s.ctn["qa_spider"]))
}

//initial internal dependencies
func (s *HTTPServer) regInternal(internals ...pkg.Internal) {
	for _, v := range internals {
		s.internal[v.GetName()] = v
	}
}

//Cors management

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Headers,Authorization,User-Agent, Keep-Alive, Content-Type, X-Requested-With,X-CSRF-Token,AccessToken,Token")
		c.Header("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, PATCH, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.Header("Access-Control-Max-Age", "600")
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
