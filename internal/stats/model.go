package stats

import "github.com/mariogarzac/blitzed/internal/player"

type HeadToHeadResult struct {
	PlayerA       player.Player `json:"player_a"`
	PlayerB       player.Player `json:"player_b"`
	MatchesPlayed int           `json:"matches_played"`
	PlayerAWins   int           `json:"player_a_wins"`
	PlayerBWins   int           `json:"player_b_wins"`
	Draws         int           `json:"draws"`
}
