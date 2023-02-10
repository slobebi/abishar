package user

import (
	ctrls "abishar/internal/controller"

	"abishar/internal/server/middleware"

	"github.com/labstack/echo/v4"
)

func Register(e *echo.Echo, controllers *ctrls.Controllers, jwt echo.MiddlewareFunc) {

	user := e.Group("/user")

	// without authentication
	user.POST("/login", controllers.User.Login)
	user.POST("/register", controllers.User.Register)

	// with authentication
	userAuth := e.Group("/user", jwt, middleware.MidParseSession)
	userAuth.POST("/logout", controllers.User.Logout)
	userAuth.GET("/wallet", controllers.User.GetUserWallet)
	userAuth.POST("/wallet", controllers.User.AddWallet)
}
