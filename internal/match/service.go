package match

import (
	"context"
	"errors"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type winResult struct {
	status  string
	winners []int
}

func checkWinCondition(winningScore int, cumulatives map[int]int) winResult {
	var crossed []int
	for pid, score := range cumulatives {
		if score >= winningScore {
			crossed = append(crossed, pid)
		}
	}
	if len(crossed) == 0 {
		return winResult{status: "in_progress"}
	}

	maxScore := cumulatives[crossed[0]]
	for _, pid := range crossed[1:] {
		if cumulatives[pid] > maxScore {
			maxScore = cumulatives[pid]
		}
	}

	var top []int
	for _, pid := range crossed {
		if cumulatives[pid] == maxScore {
			top = append(top, pid)
		}
	}

	if len(top) == 1 {
		return winResult{status: "finished", winners: top}
	}
	return winResult{status: "draw", winners: top}
}

func (s *Service) CreateMatch(ctx context.Context, winningScore int, hardMode bool, playerIDs []int) (Match, error) {
	if len(playerIDs) < 2 || len(playerIDs) > 8 {
		return Match{}, errors.New("a match requires between 2 and 8 players")
	}
	return s.repo.Create(ctx, winningScore, hardMode, playerIDs)
}

func (s *Service) SubmitRound(ctx context.Context, matchID int, inputs []ScoreInput) (Match, error) {
	_, cumulatives, err := s.repo.InsertRound(ctx, matchID, inputs)
	if err != nil {
		return Match{}, err
	}

	m, err := s.repo.GetDetail(ctx, matchID)
	if err != nil {
		return Match{}, err
	}

	result := checkWinCondition(m.WinningScore, cumulatives)
	if result.status != "in_progress" {
		tx, err := s.repo.Pool().Begin(ctx)
		if err != nil {
			return Match{}, err
		}
		defer tx.Rollback(ctx)

		if err := s.repo.UpdateFinished(ctx, tx, matchID, result.status, time.Now(), result.winners); err != nil {
			return Match{}, err
		}
		if err := tx.Commit(ctx); err != nil {
			return Match{}, err
		}
	}

	return s.repo.GetDetail(ctx, matchID)
}

func (s *Service) UndoLastRound(ctx context.Context, matchID int) (Match, error) {
	if err := s.repo.UndoLastRound(ctx, matchID); err != nil {
		return Match{}, err
	}
	return s.repo.GetDetail(ctx, matchID)
}

func (s *Service) EndMatch(ctx context.Context, matchID int) (Match, error) {
	m, err := s.repo.GetDetail(ctx, matchID)
	if err != nil {
		return Match{}, err
	}
	if m.Status != "in_progress" {
		return Match{}, errors.New("match is not in progress")
	}

	playerIDs, err := s.repo.GetPlayerIDs(ctx, matchID)
	if err != nil {
		return Match{}, err
	}

	cumulatives, err := s.repo.GetCurrentCumulatives(ctx, matchID, playerIDs)
	if err != nil {
		return Match{}, err
	}

	maxScore := 0
	for _, score := range cumulatives {
		if score > maxScore {
			maxScore = score
		}
	}

	var top []int
	for pid, score := range cumulatives {
		if score == maxScore {
			top = append(top, pid)
		}
	}

	status := "finished"
	if len(top) > 1 {
		status = "draw"
	}

	tx, err := s.repo.Pool().Begin(ctx)
	if err != nil {
		return Match{}, err
	}
	defer tx.Rollback(ctx)

	if err := s.repo.UpdateFinished(ctx, tx, matchID, status, time.Now(), top); err != nil {
		return Match{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return Match{}, err
	}

	return s.repo.GetDetail(ctx, matchID)
}
