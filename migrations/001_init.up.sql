CREATE TABLE players (
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(64) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE matches (
    id            SERIAL PRIMARY KEY,
    winning_score INTEGER     NOT NULL DEFAULT 75,
    is_standard   BOOLEAN     NOT NULL DEFAULT true,
    status        VARCHAR(16) NOT NULL DEFAULT 'in_progress'
                  CHECK (status IN ('in_progress', 'finished', 'draw')),
    started_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    finished_at   TIMESTAMPTZ
);

CREATE TABLE match_players (
    match_id  INTEGER NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    player_id INTEGER NOT NULL REFERENCES players(id) ON DELETE RESTRICT,
    PRIMARY KEY (match_id, player_id)
);

CREATE TABLE rounds (
    id         SERIAL PRIMARY KEY,
    match_id   INTEGER NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    round_num  INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (match_id, round_num)
);

CREATE TABLE round_scores (
    round_id   INTEGER NOT NULL REFERENCES rounds(id) ON DELETE CASCADE,
    player_id  INTEGER NOT NULL REFERENCES players(id) ON DELETE RESTRICT,
    delta      INTEGER NOT NULL,
    cumulative INTEGER NOT NULL,
    PRIMARY KEY (round_id, player_id)
);

CREATE TABLE match_winners (
    match_id  INTEGER NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    player_id INTEGER NOT NULL REFERENCES players(id) ON DELETE RESTRICT,
    PRIMARY KEY (match_id, player_id)
);

CREATE INDEX idx_rounds_match_id       ON rounds(match_id);
CREATE INDEX idx_round_scores_round_id ON round_scores(round_id);
CREATE INDEX idx_match_players_player  ON match_players(player_id);
CREATE INDEX idx_match_winners_player  ON match_winners(player_id);
