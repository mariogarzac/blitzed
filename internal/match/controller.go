package match

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	svc  *Service
	repo *Repository
}

func NewController(svc *Service, repo *Repository) *Controller {
	return &Controller{svc: svc, repo: repo}
}

func (c *Controller) List(ctx echo.Context) error {
	status := ctx.QueryParam("status")
	matches, err := c.repo.List(ctx.Request().Context(), status)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, matches)
}

func (c *Controller) Create(ctx echo.Context) error {
	var body struct {
		PlayerIDs    []int `json:"player_ids"`
		WinningScore int   `json:"winning_score"`
	}
	if err := ctx.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if body.WinningScore <= 0 {
		body.WinningScore = 75
	}

	m, err := c.svc.CreateMatch(ctx.Request().Context(), body.WinningScore, body.PlayerIDs)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return ctx.JSON(http.StatusCreated, m)
}

func (c *Controller) GetDetail(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid match id")
	}

	m, err := c.repo.GetDetail(ctx.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "match not found")
	}
	return ctx.JSON(http.StatusOK, m)
}

func (c *Controller) SubmitRound(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid match id")
	}

	var body struct {
		Scores []struct {
			PlayerID int `json:"player_id"`
			Delta    int `json:"delta"`
		} `json:"scores"`
	}
	if err := ctx.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if len(body.Scores) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "scores are required")
	}

	inputs := make([]ScoreInput, len(body.Scores))
	for i, s := range body.Scores {
		inputs[i] = ScoreInput{PlayerID: s.PlayerID, Delta: s.Delta}
	}

	m, err := c.svc.SubmitRound(ctx.Request().Context(), id, inputs)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return ctx.JSON(http.StatusOK, m)
}

func (c *Controller) UndoLastRound(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid match id")
	}

	m, err := c.svc.UndoLastRound(ctx.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return ctx.JSON(http.StatusOK, m)
}

func (c *Controller) Delete(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid match id")
	}
	if err := c.repo.Delete(ctx.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	return ctx.NoContent(http.StatusNoContent)
}

func (c *Controller) EndMatch(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid match id")
	}

	m, err := c.svc.EndMatch(ctx.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return ctx.JSON(http.StatusOK, m)
}
