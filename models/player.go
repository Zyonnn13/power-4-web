package models

type Player struct {
	Name   string
	Color  string
	Symbol int
	Wins   int
}

func NewPlayer(name string, color string, symbol int) Player {
	return Player{
		Name:   name,
		Color:  color,
		Symbol: symbol,
		Wins:   0,
	}
}
