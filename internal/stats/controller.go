package stats

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	repo *Repository
}

func NewController(repo *Repository) *Controller {
	return &Controller{repo: repo}
}

func (c *Controller) HeadToHead(ctx echo.Context) error {
	playerA, err := strconv.Atoi(ctx.QueryParam("player_a"))
	if err != nil || playerA == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "player_a is required")
	}
	playerB, err := strconv.Atoi(ctx.QueryParam("player_b"))
	if err != nil || playerB == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "player_b is required")
	}
	if playerA == playerB {
		return echo.NewHTTPError(http.StatusBadRequest, "player_a and player_b must be different")
	}

	result, err := c.repo.HeadToHead(ctx.Request().Context(), playerA, playerB)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, result)
}
