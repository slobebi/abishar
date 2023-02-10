package server

import (
	"abishar/internal/config"
	"abishar/internal/controller"
	"abishar/internal/server/routes"
	"log"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tylerb/graceful"
)

type HTTPServerItf interface {
	ListenAndServe() error
}

type httpServer struct {
	server         *graceful.Server
	config         config.Config
	echo           *echo.Echo
	ctrl           *controller.Controllers
	preMiddlewares []echo.MiddlewareFunc
	allMiddlewares []echo.MiddlewareFunc
}

func NewHTTPServer(
	cfg config.Config,
	ctrl *controller.Controllers,
	pm []echo.MiddlewareFunc,
	am []echo.MiddlewareFunc,
) HTTPServerItf {
	server := &httpServer{
		echo:           echo.New(),
		config:         cfg,
		ctrl:           ctrl,
		preMiddlewares: pm,
		allMiddlewares: am,
	}

	server.connectCoreWithEcho(server.config.JWT)
	server.initGracefulServer()

	return server
}

func (h *httpServer) ListenAndServe() error {
	config := h.config.HTTPServer
	echo := h.echo
	setServerObj(h.echo, config)

	log.Println("server started")
	log.Println("time_human", time.Now().Format(time.RFC3339))
	log.Println("address", echo.Server.Addr)
	log.Println("graceful_timeout", config.GracefulTimeout)

	return h.server.ListenAndServe()
}

func (h *httpServer) connectCoreWithEcho(jwt config.JWT) {
	e := h.echo
	controller := h.ctrl

	e.Pre(h.preMiddlewares...)
	e.Use(h.allMiddlewares...)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	routes.Register(e, controller, jwt)

	// Set custom error handler
	setServerObj(e, h.config.HTTPServer)
}

func (h *httpServer) initGracefulServer() {
	config := h.config.HTTPServer
	echo := h.echo

	h.server = &graceful.Server{
		Server:  echo.Server,
		Timeout: config.GracefulTimeout,
		Logger:  graceful.DefaultLogger(),
	}
}

func setServerObj(e *echo.Echo, server config.HTTPServer) {
	e.Server.Addr = server.ListenAddress + ":" + strconv.Itoa(server.Port)
	if server.ReadTimeout > 0 {
		e.Server.ReadTimeout = server.ReadTimeout
	}
	if server.WriteTimeout > 0 {
		e.Server.WriteTimeout = server.WriteTimeout
	}
	if server.IdleTimeout > 0 {
		e.Server.IdleTimeout = server.IdleTimeout
	}
}
