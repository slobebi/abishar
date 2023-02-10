package vote

import (
	ctrls "abishar/internal/controller"

	"abishar/internal/server/middleware"

	"github.com/labstack/echo/v4"
)

func Register(e *echo.Echo, controllers *ctrls.Controllers, jwt echo.MiddlewareFunc) {

	vote := e.Group("/transaction", jwt, middleware.MidParseSession)

	vote.POST("/vote", controllers.Vote.Vote)
	vote.GET("/unvote", controllers.Vote.UnVote)
	vote.GET("/get", controllers.Vote.GetVotedMovieByUser)
	vote.GET("/get-all", controllers.Vote.GetMostVotedMovie)

}
