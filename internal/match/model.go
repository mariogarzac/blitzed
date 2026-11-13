package match

import (
	"time"

	"github.com/mariogarzac/blitzed/internal/player"
)

type Match struct {
	ID           int             `json:"id"`
	WinningScore int             `json:"winning_score"`
	IsStandard   bool            `json:"is_standard"`
	HardMode     bool            `json:"hard_mode"`
	Status       string          `json:"status"`
	StartedAt    time.Time       `json:"started_at"`
	FinishedAt   *time.Time      `json:"finished_at"`
	Players      []player.Player `json:"players"`
	Winners      []player.Player `json:"winners"`
	Rounds       []Round         `json:"rounds"`
}

type MatchSummary struct {
	ID           int             `json:"id"`
	WinningScore int             `json:"winning_score"`
	IsStandard   bool            `json:"is_standard"`
	Status       string          `json:"status"`
	StartedAt    time.Time       `json:"started_at"`
	FinishedAt   *time.Time      `json:"finished_at"`
	Winners      []player.Player `json:"winners"`
	PlayerCount  int             `json:"player_count"`
}

type Round struct {
	ID        int          `json:"id"`
	MatchID   int          `json:"match_id"`
	RoundNum  int          `json:"round_num"`
	CreatedAt time.Time    `json:"created_at"`
	Scores    []RoundScore `json:"scores"`
}

type RoundScore struct {
	PlayerID   int `json:"player_id"`
	Delta      int `json:"delta"`
	Cumulative int `json:"cumulative"`
}
