package stats

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) HeadToHead(ctx context.Context, playerAID, playerBID int) (HeadToHeadResult, error) {
	var result HeadToHeadResult

	err := r.pool.QueryRow(ctx,
		`SELECT id, name, created_at FROM players WHERE id = $1`, playerAID,
	).Scan(&result.PlayerA.ID, &result.PlayerA.Name, &result.PlayerA.CreatedAt)
	if err != nil {
		return result, err
	}
	err = r.pool.QueryRow(ctx,
		`SELECT id, name, created_at FROM players WHERE id = $1`, playerBID,
	).Scan(&result.PlayerB.ID, &result.PlayerB.Name, &result.PlayerB.CreatedAt)
	if err != nil {
		return result, err
	}

	rows, err := r.pool.Query(ctx, `
		SELECT m.id, m.status
		FROM matches m
		JOIN match_players mpa ON mpa.match_id = m.id AND mpa.player_id = $1
		JOIN match_players mpb ON mpb.match_id = m.id AND mpb.player_id = $2
		WHERE m.is_standard = true AND m.status IN ('finished', 'draw')`,
		playerAID, playerBID,
	)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	type matchRow struct {
		id     int
		status string
	}
	var matchRows []matchRow
	for rows.Next() {
		var mr matchRow
		if err := rows.Scan(&mr.id, &mr.status); err != nil {
			return result, err
		}
		matchRows = append(matchRows, mr)
	}
	rows.Close()

	result.MatchesPlayed = len(matchRows)

	for _, mr := range matchRows {
		if mr.status == "draw" {
			result.Draws++
			continue
		}
		var winnerID int
		err := r.pool.QueryRow(ctx,
			`SELECT player_id FROM match_winners WHERE match_id = $1 LIMIT 1`, mr.id,
		).Scan(&winnerID)
		if err != nil {
			continue
		}
		switch winnerID {
		case playerAID:
			result.PlayerAWins++
		case playerBID:
			result.PlayerBWins++
		}
	}

	return result, nil
}
