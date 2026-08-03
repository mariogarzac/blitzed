package player

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	repo *Repository
}

func NewController(repo *Repository) *Controller {
	return &Controller{repo: repo}
}

func (c *Controller) List(ctx echo.Context) error {
	players, err := c.repo.List(ctx.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if players == nil {
		players = []Player{}
	}
	return ctx.JSON(http.StatusOK, players)
}

func (c *Controller) Create(ctx echo.Context) error {
	var body struct {
		Name string `json:"name"`
	}
	if err := ctx.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if body.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "name is required")
	}

	p, err := c.repo.Create(ctx.Request().Context(), body.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusCreated, p)
}

func (c *Controller) Delete(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid player id")
	}

	if err := c.repo.Delete(ctx.Request().Context(), id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "player not found")
		}
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return ctx.NoContent(http.StatusNoContent)
}
