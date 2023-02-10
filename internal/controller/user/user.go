package user

import (
	"context"
	"net/http"

	enUser "abishar/internal/entity/user"

	"github.com/labstack/echo/v4"
)

type userUsecase interface {
	RegisterUser(ctx context.Context, form enUser.RegisterForm) (*enUser.RegisterResponse, error)
	Login(ctx context.Context, form enUser.LoginRequest) (*enUser.RegisterResponse, error)
	Logout(sess enUser.Session) error
	GetUserSession(sess enUser.Session) *enUser.SessionData
	GetUserWallet(ctx context.Context, username string) (*enUser.UserWallet, error)
	AddWallet(ctx context.Context, amount int64, userID int64) error
}

type Controller struct {
	user userUsecase
}

func NewController(user userUsecase) *Controller {
	return &Controller{
		user: user,
	}
}

func (c *Controller) Register(ctx echo.Context) error {

	form := enUser.RegisterForm{}

	if err := ctx.Bind(&form); err != nil {
		return ctx.JSON(http.StatusBadRequest,
			map[string]interface{}{
				"Error": "Bad Request",
			},
		)
	}

	response, err := c.user.RegisterUser(ctx.Request().Context(), form)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError,
			map[string]interface{}{
				"Error": err.Error(),
			},
		)
	}

	return ctx.JSON(http.StatusOK,
		map[string]interface{}{
			"Data": response,
		},
	)

}

func (c *Controller) Login(ctx echo.Context) error {

	form := enUser.LoginRequest{}

	if err := ctx.Bind(&form); err != nil {
		return ctx.JSON(http.StatusBadRequest,
			map[string]interface{}{
				"Error": "Bad Request",
			},
		)
	}

	response, err := c.user.Login(ctx.Request().Context(), form)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError,
			map[string]interface{}{
				"Error": err.Error(),
			},
		)
	}

	return ctx.JSON(http.StatusOK,
		map[string]interface{}{
			"Data": response,
		},
	)
}

func (c *Controller) Logout(ctx echo.Context) error {

	session := ctx.Get(enUser.SessionContextKey).(enUser.Session)

	sessionData := c.user.GetUserSession(session)
	if sessionData == nil {
		return ctx.JSON(http.StatusUnauthorized,
			map[string]interface{}{
				"Error": "Unauthorized",
			},
		)
	}

	err := c.user.Logout(session)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError,
			map[string]interface{}{
				"Error": err.Error(),
			},
		)
	}

	return ctx.JSON(http.StatusOK, "Success Log out")
}

func (c *Controller) GetUserWallet(ctx echo.Context) error {
	session := ctx.Get(enUser.SessionContextKey).(enUser.Session)

	sessionData := c.user.GetUserSession(session)
	if sessionData == nil {
		return ctx.JSON(http.StatusUnauthorized,
			map[string]interface{}{
				"Error": "Unauthorized",
			},
		)
	}

	response, err := c.user.GetUserWallet(ctx.Request().Context(), session.Username)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError,
			map[string]interface{}{
				"Error": err.Error(),
			},
		)
	}

	return ctx.JSON(http.StatusOK,
		map[string]interface{}{
			"Data": response,
		},
	)
}

func (c *Controller) AddWallet(ctx echo.Context) error {
	session := ctx.Get(enUser.SessionContextKey).(enUser.Session)

	sessionData := c.user.GetUserSession(session)
	if sessionData == nil {
		return ctx.JSON(http.StatusUnauthorized,
			map[string]interface{}{
				"Error": "Unauthorized",
			},
		)
	}

	form := enUser.WalletRequest{}

	if err := ctx.Bind(&form); err != nil {
		return ctx.JSON(http.StatusBadRequest,
			map[string]interface{}{
				"Error": "Bad Request",
			},
		)
	}

	err := c.user.AddWallet(ctx.Request().Context(), form.Amount, session.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError,
			map[string]interface{}{
				"Error": err.Error(),
			},
		)
	}

	return ctx.JSON(http.StatusOK, "Success")

}
