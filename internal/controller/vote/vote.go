package vote

import (
	enUser "abishar/internal/entity/user"
	enVote "abishar/internal/entity/votes"
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	userUsecase interface {
		GetUserSession(sess enUser.Session) *enUser.SessionData
	}

	voteUsecase interface {
		Vote(ctx context.Context, form enVote.VoteRequest) (int64, error)
		UnVote(ctx context.Context, form enVote.VoteRequest) error
		GetMostVotedMovie(ctx context.Context) (*enVote.VotedMovies, error)
		GetVotedMovieByUser(ctx context.Context, userID int64) ([]enVote.VotedMovies, error)
	}
)

type Controller struct {
	userUc userUsecase
	voteUc voteUsecase
}

func NewController(
	userUc userUsecase,
	voteUc voteUsecase,
) *Controller {
	return &Controller{
		userUc: userUc,
		voteUc: voteUc,
	}
}

func (c *Controller) Vote(ctx echo.Context) error {

	session := ctx.Get(enUser.SessionContextKey).(enUser.Session)

	sessionData := c.userUc.GetUserSession(session)
	if sessionData == nil {
		return ctx.JSON(http.StatusUnauthorized,
			map[string]interface{}{
				"Error": "Unauthorized",
			},
		)
	}

	form := enVote.VoteRequest{}

	if err := ctx.Bind(&form); err != nil {
		return ctx.JSON(http.StatusBadRequest,
			map[string]interface{}{
				"Error": "Bad Request",
			},
		)
	}

	id, err := c.voteUc.Vote(ctx.Request().Context(), form)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError,
			map[string]interface{}{
				"Error": err.Error(),
			},
		)
	}

	return ctx.JSON(http.StatusOK,
		map[string]interface{}{
			"Status": "Success",
			"ID":     id,
		},
	)
}

func (c *Controller) UnVote(ctx echo.Context) error {

	session := ctx.Get(enUser.SessionContextKey).(enUser.Session)

	sessionData := c.userUc.GetUserSession(session)
	if sessionData == nil {
		return ctx.JSON(http.StatusUnauthorized,
			map[string]interface{}{
				"Error": "Unauthorized",
			},
		)
	}

	form := enVote.VoteRequest{}

	if err := ctx.Bind(&form); err != nil {
		return ctx.JSON(http.StatusBadRequest,
			map[string]interface{}{
				"Error": "Bad Request",
			},
		)
	}

	err := c.voteUc.UnVote(ctx.Request().Context(), form)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError,
			map[string]interface{}{
				"Error": err.Error(),
			},
		)
	}

	return ctx.JSON(http.StatusOK,
		map[string]interface{}{
			"Status": "Success",
		},
	)
}

func (c *Controller) GetMostVotedMovie(ctx echo.Context) error {

	session := ctx.Get(enUser.SessionContextKey).(enUser.Session)

	sessionData := c.userUc.GetUserSession(session)
	if sessionData == nil {
		return ctx.JSON(http.StatusUnauthorized,
			map[string]interface{}{
				"Error": "Unauthorized",
			},
		)
	}

	if !session.IsAdmin {
		return ctx.JSON(http.StatusUnauthorized,
			map[string]interface{}{
				"Error": "Not Admin",
			},
		)
	}

	response, err := c.voteUc.GetMostVotedMovie(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError,
			map[string]interface{}{
				"Error": err.Error(),
			},
		)
	}

	return ctx.JSON(http.StatusOK,
		map[string]interface{}{
			"Status": "Success",
			"Data":   response,
		},
	)
}

func (c *Controller) GetVotedMovieByUser(ctx echo.Context) error {

	session := ctx.Get(enUser.SessionContextKey).(enUser.Session)

	sessionData := c.userUc.GetUserSession(session)
	if sessionData == nil {
		return ctx.JSON(http.StatusUnauthorized,
			map[string]interface{}{
				"Error": "Unauthorized",
			},
		)
	}

	response, err := c.voteUc.GetVotedMovieByUser(ctx.Request().Context(), session.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError,
			map[string]interface{}{
				"Error": err.Error(),
			},
		)
	}

	return ctx.JSON(http.StatusOK,
		map[string]interface{}{
			"Status": "Success",
			"Data":   response,
		},
	)
}
