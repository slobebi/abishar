package app

import (
	"net/http"

	"abishar/internal/config"

	"abishar/internal/server"

	// Repositories
	movieRepo "abishar/internal/repository/movies"
	userRepo "abishar/internal/repository/user"
	voteRepo "abishar/internal/repository/votes"

	// Usecases
	movieUsc "abishar/internal/usecase/movies"
	userUsc "abishar/internal/usecase/user"
	voteUsc "abishar/internal/usecase/votes"

	// Controllers
	ctrls "abishar/internal/controller"
	movieCtrl "abishar/internal/controller/movie"
	userCtrl "abishar/internal/controller/user"
	voteCtrl "abishar/internal/controller/vote"

	"github.com/labstack/echo/v4"
	echoMid "github.com/labstack/echo/v4/middleware"
)

func InitHTTPServer(cfg config.Config) server.HTTPServerItf {
	// Drivers
	db := connectDatabase(cfg.Database)
	redis := connectRedis(cfg.Redis)

	// Initialize Repositories
	userRepository := userRepo.NewRepository(db, redis, cfg.JWT)
	movieRepository := movieRepo.NewRepository(db, redis)
	voteRepository := voteRepo.NewRepository(db, redis)

	// Initialize Usecases
	userUsecase := userUsc.NewUsecase(userRepository)
	movieUsecase := movieUsc.NewUsecase(movieRepository)
	voteUsecase := voteUsc.NewUsecase(voteRepository, movieRepository)

	// Initialize Controllers
	userController := userCtrl.NewController(userUsecase)
	movieController := movieCtrl.NewController(userUsecase, movieUsecase)
	voteController := voteCtrl.NewController(userUsecase, voteUsecase)

	controllers := ctrls.NewControllers(
		userController,
		movieController,
		voteController,
	)

	preMiddlewares := []echo.MiddlewareFunc{
		echoMid.RemoveTrailingSlashWithConfig(echoMid.TrailingSlashConfig{
			RedirectCode: http.StatusMovedPermanently,
		}),
	}

	// Middlewares that runs after the router
	allMiddlewares := []echo.MiddlewareFunc{
		echoMid.Gzip(),
		echoMid.CORSWithConfig(echoMid.CORSConfig{
			Skipper:      echoMid.DefaultSkipper,
			AllowOrigins: []string{"*"},
			AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		}),
		echoMid.Recover(),
		echoMid.RequestID(),
	}

	httpServerItf := server.NewHTTPServer(cfg, controllers, preMiddlewares, allMiddlewares)
	return httpServerItf
}
