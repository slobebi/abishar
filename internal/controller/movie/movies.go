package movie

import (
	enMovie "abishar/internal/entity/movies"
	enUser "abishar/internal/entity/user"
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type (
	userUsecase interface {
		GetUserSession(sess enUser.Session) *enUser.SessionData
	}

	movieUsecase interface {
		InsertMovie(ctx context.Context, form enMovie.MovieRequest) (int64, error)
		UpdateMovie(ctx context.Context, form enMovie.Movies) error
		GetMovies(ctx context.Context, page, limit int) ([]enMovie.Movies, error)
		SearchMovies(ctx context.Context, page, limit int, title, description, artists, genres string) ([]enMovie.Movies, error)
	}
)

type Controller struct {
	userUsc  userUsecase
	movieUsc movieUsecase
}

func NewController(
	userUsc userUsecase,
	movieUsc movieUsecase,
) *Controller {
	return &Controller{
		userUsc:  userUsc,
		movieUsc: movieUsc,
	}
}

func (c *Controller) InsertMovies(ctx echo.Context) error {

	session := ctx.Get(enUser.SessionContextKey).(enUser.Session)

	sessionData := c.userUsc.GetUserSession(session)
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

	form := enMovie.MovieRequest{}

	if err := ctx.Bind(&form); err != nil {
		return ctx.JSON(http.StatusBadRequest,
			map[string]interface{}{
				"Error": "Bad Request",
			},
		)
	}

	movieId, err := c.movieUsc.InsertMovie(ctx.Request().Context(), form)
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
			"ID":     movieId,
		},
	)
}

func (c *Controller) UpdateMovie(ctx echo.Context) error {

	session := ctx.Get(enUser.SessionContextKey).(enUser.Session)

	sessionData := c.userUsc.GetUserSession(session)
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

	form := enMovie.Movies{}

	if err := ctx.Bind(&form); err != nil {
		return ctx.JSON(http.StatusBadRequest,
			map[string]interface{}{
				"Error": "Bad Request",
			},
		)
	}

	err := c.movieUsc.UpdateMovie(ctx.Request().Context(), form)
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

func (c *Controller) GetMovies(ctx echo.Context) error {
	page := ctx.QueryParam("page")
	limit := ctx.QueryParam("limit")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest,
			map[string]interface{}{
				"Error": "Bad Request",
			},
		)
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest,
			map[string]interface{}{
				"Error": "Bad Request",
			},
		)
	}

	if pageInt == 0 {
		pageInt = 1
	}

	if limitInt == 0 {
		limitInt = 10
	}

	movies, err := c.movieUsc.GetMovies(ctx.Request().Context(), pageInt, limitInt)
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
			"Data":   movies,
		},
	)
}

func (c *Controller) SearchMovie(ctx echo.Context) error {
	page := ctx.QueryParam("page")
	limit := ctx.QueryParam("limit")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest,
			map[string]interface{}{
				"Error": "Bad Request",
			},
		)
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest,
			map[string]interface{}{
				"Error": "Bad Request",
			},
		)
	}

	title := ctx.QueryParam("title")
	description := ctx.QueryParam("description")
	artists := ctx.QueryParam("artists")
	genres := ctx.QueryParam("genres")

	if pageInt == 0 {
		pageInt = 1
	}

	if limitInt == 0 {
		limitInt = 10
	}

	movies, err := c.movieUsc.SearchMovies(ctx.Request().Context(), pageInt, limitInt, title, description, artists, genres)
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
			"Data":   movies,
		},
	)
}
