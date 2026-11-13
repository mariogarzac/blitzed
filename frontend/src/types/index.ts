export interface Player {
  id: number
  name: string
  created_at: string
}

export interface RoundScore {
  player_id: number
  delta: number
  cumulative: number
}

export interface Round {
  id: number
  match_id: number
  round_num: number
  created_at: string
  scores: RoundScore[]
}

export interface Match {
  id: number
  winning_score: number
  is_standard: boolean
  hard_mode: boolean
  status: 'in_progress' | 'finished' | 'draw'
  started_at: string
  finished_at: string | null
  players: Player[]
  winners: Player[]
  rounds: Round[]
}

export interface MatchSummary {
  id: number
  winning_score: number
  is_standard: boolean
  status: 'in_progress' | 'finished' | 'draw'
  started_at: string
  finished_at: string | null
  winners: Player[]
  player_count: number
}

export interface HeadToHeadResult {
  player_a: Player
  player_b: Player
  matches_played: number
  player_a_wins: number
  player_b_wins: number
  draws: number
}

export interface RoundScoreInput {
  player_id: number
  delta: number
}
