package routes

import (
	"abishar/internal/config"
	ctrls "abishar/internal/controller"
	enUser "abishar/internal/entity/user"
	"abishar/internal/server/routes/movie"
	"abishar/internal/server/routes/user"
	"abishar/internal/server/routes/vote"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func Register(e *echo.Echo, controller *ctrls.Controllers, cfg config.JWT) {
	jwtMiddleware := echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(cfg.Secret),
		ContextKey: enUser.SessionContextKey,
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(enUser.TokenClaim)
		},
	})
	user.Register(e, controller, jwtMiddleware)
	vote.Register(e, controller, jwtMiddleware)
	movie.Register(e, controller, jwtMiddleware)
}
