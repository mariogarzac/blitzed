package match

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mariogarzac/blitzed/internal/player"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Pool() *pgxpool.Pool {
	return r.pool
}

// --- Match CRUD ---

func (r *Repository) Create(ctx context.Context, winningScore int, playerIDs []int) (Match, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return Match{}, err
	}
	defer tx.Rollback(ctx)

	isStandard := winningScore == 75

	var m Match
	err = tx.QueryRow(ctx,
		`INSERT INTO matches (winning_score, is_standard) VALUES ($1, $2)
		 RETURNING id, winning_score, is_standard, status, started_at, finished_at`,
		winningScore, isStandard,
	).Scan(&m.ID, &m.WinningScore, &m.IsStandard, &m.Status, &m.StartedAt, &m.FinishedAt)
	if err != nil {
		return Match{}, err
	}

	for _, pid := range playerIDs {
		if _, err = tx.Exec(ctx,
			`INSERT INTO match_players (match_id, player_id) VALUES ($1, $2)`,
			m.ID, pid,
		); err != nil {
			return Match{}, err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return Match{}, err
	}

	return r.GetDetail(ctx, m.ID)
}

func (r *Repository) List(ctx context.Context, status string) ([]MatchSummary, error) {
	query := `
		SELECT m.id, m.winning_score, m.is_standard, m.status, m.started_at, m.finished_at,
		       COUNT(DISTINCT mp.player_id) as player_count
		FROM matches m
		LEFT JOIN match_players mp ON mp.match_id = m.id`

	args := []any{}
	if status != "" {
		query += ` WHERE m.status = $1`
		args = append(args, status)
	}
	query += ` GROUP BY m.id ORDER BY m.started_at DESC`

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []MatchSummary
	for rows.Next() {
		var s MatchSummary
		if err := rows.Scan(&s.ID, &s.WinningScore, &s.IsStandard, &s.Status, &s.StartedAt, &s.FinishedAt, &s.PlayerCount); err != nil {
			return nil, err
		}
		summaries = append(summaries, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i, s := range summaries {
		winners, err := r.getWinners(ctx, s.ID)
		if err != nil {
			return nil, err
		}
		summaries[i].Winners = winners
	}

	if summaries == nil {
		summaries = []MatchSummary{}
	}
	return summaries, nil
}

func (r *Repository) GetDetail(ctx context.Context, id int) (Match, error) {
	var m Match
	err := r.pool.QueryRow(ctx,
		`SELECT id, winning_score, is_standard, status, started_at, finished_at FROM matches WHERE id = $1`, id,
	).Scan(&m.ID, &m.WinningScore, &m.IsStandard, &m.Status, &m.StartedAt, &m.FinishedAt)
	if err != nil {
		return Match{}, err
	}

	players, err := r.getPlayers(ctx, id)
	if err != nil {
		return Match{}, err
	}
	m.Players = players

	winners, err := r.getWinners(ctx, id)
	if err != nil {
		return Match{}, err
	}
	m.Winners = winners

	rounds, err := r.getRounds(ctx, id)
	if err != nil {
		return Match{}, err
	}
	m.Rounds = rounds

	return m, nil
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM matches WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("match not found")
	}
	return nil
}

func (r *Repository) UpdateFinished(ctx context.Context, tx pgx.Tx, matchID int, status string, finishedAt time.Time, winnerIDs []int) error {
	if _, err := tx.Exec(ctx,
		`UPDATE matches SET status = $1, finished_at = $2 WHERE id = $3`,
		status, finishedAt, matchID,
	); err != nil {
		return err
	}

	for _, wid := range winnerIDs {
		if _, err := tx.Exec(ctx,
			`INSERT INTO match_winners (match_id, player_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
			matchID, wid,
		); err != nil {
			return err
		}
	}
	return nil
}

// --- Round operations ---

type ScoreInput struct {
	PlayerID int
	Delta    int
}

func (r *Repository) InsertRound(ctx context.Context, matchID int, inputs []ScoreInput) (int, map[int]int, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, nil, err
	}
	defer tx.Rollback(ctx)

	var status string
	if err = tx.QueryRow(ctx, `SELECT status FROM matches WHERE id = $1`, matchID).Scan(&status); err != nil {
		return 0, nil, err
	}
	if status != "in_progress" {
		return 0, nil, errors.New("match is already finished")
	}

	var maxRound *int
	_ = tx.QueryRow(ctx, `SELECT MAX(round_num) FROM rounds WHERE match_id = $1`, matchID).Scan(&maxRound)
	roundNum := 1
	if maxRound != nil {
		roundNum = *maxRound + 1
	}

	var roundID int
	if err = tx.QueryRow(ctx,
		`INSERT INTO rounds (match_id, round_num) VALUES ($1, $2) RETURNING id`,
		matchID, roundNum,
	).Scan(&roundID); err != nil {
		return 0, nil, err
	}

	cumulatives := make(map[int]int, len(inputs))
	for _, inp := range inputs {
		var lastCum *int
		_ = tx.QueryRow(ctx, `
			SELECT rs.cumulative FROM round_scores rs
			JOIN rounds ro ON ro.id = rs.round_id
			WHERE ro.match_id = $1 AND rs.player_id = $2
			ORDER BY ro.round_num DESC LIMIT 1`,
			matchID, inp.PlayerID,
		).Scan(&lastCum)

		prev := 0
		if lastCum != nil {
			prev = *lastCum
		}
		cumulatives[inp.PlayerID] = prev + inp.Delta
	}

	for _, inp := range inputs {
		if _, err = tx.Exec(ctx,
			`INSERT INTO round_scores (round_id, player_id, delta, cumulative) VALUES ($1, $2, $3, $4)`,
			roundID, inp.PlayerID, inp.Delta, cumulatives[inp.PlayerID],
		); err != nil {
			return 0, nil, err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, nil, err
	}

	return roundNum, cumulatives, nil
}

func (r *Repository) UndoLastRound(ctx context.Context, matchID int) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var status string
	if err = tx.QueryRow(ctx, `SELECT status FROM matches WHERE id = $1`, matchID).Scan(&status); err != nil {
		return err
	}
	if status != "in_progress" {
		return errors.New("can only undo rounds in an in_progress match")
	}

	var roundID *int
	_ = tx.QueryRow(ctx,
		`SELECT id FROM rounds WHERE match_id = $1 ORDER BY round_num DESC LIMIT 1`, matchID,
	).Scan(&roundID)

	if roundID == nil {
		return errors.New("no rounds to undo")
	}

	if _, err = tx.Exec(ctx, `DELETE FROM rounds WHERE id = $1`, *roundID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *Repository) GetCurrentCumulatives(ctx context.Context, matchID int, playerIDs []int) (map[int]int, error) {
	result := make(map[int]int, len(playerIDs))
	for _, pid := range playerIDs {
		result[pid] = 0
	}

	rows, err := r.pool.Query(ctx, `
		SELECT DISTINCT ON (rs.player_id) rs.player_id, rs.cumulative
		FROM round_scores rs
		JOIN rounds ro ON ro.id = rs.round_id
		WHERE ro.match_id = $1
		ORDER BY rs.player_id, ro.round_num DESC`,
		matchID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var pid, cum int
		if err := rows.Scan(&pid, &cum); err != nil {
			return nil, err
		}
		result[pid] = cum
	}
	return result, rows.Err()
}

func (r *Repository) GetPlayerIDs(ctx context.Context, matchID int) ([]int, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT player_id FROM match_players WHERE match_id = $1`, matchID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// --- Private helpers ---

func (r *Repository) getPlayers(ctx context.Context, matchID int) ([]player.Player, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT p.id, p.name, p.created_at FROM players p
		 JOIN match_players mp ON mp.player_id = p.id
		 WHERE mp.match_id = $1 ORDER BY p.name`, matchID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []player.Player
	for rows.Next() {
		var p player.Player
		if err := rows.Scan(&p.ID, &p.Name, &p.CreatedAt); err != nil {
			return nil, err
		}
		players = append(players, p)
	}
	if players == nil {
		players = []player.Player{}
	}
	return players, rows.Err()
}

func (r *Repository) getWinners(ctx context.Context, matchID int) ([]player.Player, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT p.id, p.name, p.created_at FROM players p
		 JOIN match_winners mw ON mw.player_id = p.id
		 WHERE mw.match_id = $1`, matchID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var winners []player.Player
	for rows.Next() {
		var p player.Player
		if err := rows.Scan(&p.ID, &p.Name, &p.CreatedAt); err != nil {
			return nil, err
		}
		winners = append(winners, p)
	}
	if winners == nil {
		winners = []player.Player{}
	}
	return winners, rows.Err()
}

func (r *Repository) getRounds(ctx context.Context, matchID int) ([]Round, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, match_id, round_num, created_at FROM rounds WHERE match_id = $1 ORDER BY round_num`, matchID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rounds []Round
	for rows.Next() {
		var round Round
		if err := rows.Scan(&round.ID, &round.MatchID, &round.RoundNum, &round.CreatedAt); err != nil {
			return nil, err
		}
		rounds = append(rounds, round)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i, round := range rounds {
		scoreRows, err := r.pool.Query(ctx,
			`SELECT player_id, delta, cumulative FROM round_scores WHERE round_id = $1`, round.ID,
		)
		if err != nil {
			return nil, err
		}
		var scores []RoundScore
		for scoreRows.Next() {
			var s RoundScore
			if err := scoreRows.Scan(&s.PlayerID, &s.Delta, &s.Cumulative); err != nil {
				scoreRows.Close()
				return nil, err
			}
			scores = append(scores, s)
		}
		scoreRows.Close()
		if scores == nil {
			scores = []RoundScore{}
		}
		rounds[i].Scores = scores
	}

	if rounds == nil {
		rounds = []Round{}
	}
	return rounds, nil
}
