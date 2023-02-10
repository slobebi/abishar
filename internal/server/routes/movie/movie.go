package movie

import (
	ctrls "abishar/internal/controller"

	"abishar/internal/server/middleware"

	"github.com/labstack/echo/v4"
)

func Register(e *echo.Echo, controllers *ctrls.Controllers, jwt echo.MiddlewareFunc) {

	// public
	product := e.Group("/movie")
	product.GET("/all", controllers.Movie.GetMovies)
	product.GET("/search", controllers.Movie.SearchMovie)

	// need admin auth
	productAdmin := e.Group("/product", jwt, middleware.MidParseSession)
	productAdmin.POST("/insert", controllers.Movie.InsertMovies)
	productAdmin.PUT("/update", controllers.Movie.UpdateMovie)
}
