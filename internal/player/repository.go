package player

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) List(ctx context.Context) ([]Player, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name, created_at FROM players ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []Player
	for rows.Next() {
		var p Player
		if err := rows.Scan(&p.ID, &p.Name, &p.CreatedAt); err != nil {
			return nil, err
		}
		players = append(players, p)
	}
	return players, rows.Err()
}

func (r *Repository) Create(ctx context.Context, name string) (Player, error) {
	var p Player
	err := r.pool.QueryRow(ctx,
		`INSERT INTO players (name) VALUES ($1) RETURNING id, name, created_at`,
		name,
	).Scan(&p.ID, &p.Name, &p.CreatedAt)
	return p, err
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM match_players WHERE player_id = $1`, id,
	).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("player is part of existing matches and cannot be deleted")
	}

	tag, err := r.pool.Exec(ctx, `DELETE FROM players WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
