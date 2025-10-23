package models

type Scoreboard struct {
	Player1Wins int
	Player2Wins int
	Draws       int
	TotalGames  int
}

func NewScoreboard() *Scoreboard {
	return &Scoreboard{
		Player1Wins: 0,
		Player2Wins: 0,
		Draws:       0,
		TotalGames:  0,
	}
}

func (s *Scoreboard) AddWin(playerNumber int) {
	s.TotalGames++
	if playerNumber == 1 {
		s.Player1Wins++
	} else if playerNumber == 2 {
		s.Player2Wins++
	}
}

func (s *Scoreboard) AddDraw() {
	s.TotalGames++
	s.Draws++
}

func (s *Scoreboard) Reset() {
	s.Player1Wins = 0
	s.Player2Wins = 0
	s.Draws = 0
	s.TotalGames = 0
}
