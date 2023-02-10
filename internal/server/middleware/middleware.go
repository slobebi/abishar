package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"

	enUser "abishar/internal/entity/user"
)

func MidParseSession(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		token, ok := ctx.Get(enUser.SessionContextKey).(*jwt.Token)
		if token == nil || !ok {
			return ctx.JSON(http.StatusUnauthorized,
				map[string]interface{}{
					"Error": "Failed to Get Session",
				},
			)
		}

		claims, ok := token.Claims.(*enUser.TokenClaim)
		if !ok {
			return ctx.JSON(http.StatusUnauthorized,
				map[string]interface{}{
					"Error": "Failed To Get JWT Claims",
				},
			)
		}

		sess := enUser.Session{
			ID:        claims.ID,
			Username:  claims.Username,
			UniqueKey: claims.UniqueKey,
			IsAdmin:   claims.IsAdmin,
		}

		ctx.Set(enUser.SessionContextKey, sess)

		return next(ctx)
	}
}
