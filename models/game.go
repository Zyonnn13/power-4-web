package models

import "fmt"

type Game struct {
	Grid          [6][7]int
	Player1       Player
	Player2       Player
	CurrentPlayer int
	Winner        int
	TurnCount     int
}

func NewGame(p1, p2 Player) Game {
	return Game{
		Grid:          [6][7]int{},
		Player1:       p1,
		Player2:       p2,
		CurrentPlayer: 1,
		Winner:        0,
		TurnCount:     0,
	}
}

func (g *Game) AddToken(column int) bool {
	if column < 0 || column > 6 || g.Winner != 0 {
		return false
	}
	symbol := g.CurrentPlayer
	for row := len(g.Grid) - 1; row >= 0; row-- {
		if g.Grid[row][column] == 0 {
			g.Grid[row][column] = symbol
			g.TurnCount++
			if g.checkWin(row, column, symbol) {
				g.Winner = symbol
			} else if g.TurnCount == 42 {
				g.Winner = 3
			} else {
				g.CurrentPlayer = 3 - g.CurrentPlayer
			}
			return true
		}
	}
	return false
}

func (g *Game) checkWin(row, col, symbol int) bool {
	directions := [][2]int{{0, 1}, {1, 0}, {1, 1}, {1, -1}}
	for _, d := range directions {
		count := 1
		for i := 1; i < 4; i++ {
			r, c := row+d[0]*i, col+d[1]*i
			if r < 0 || r > 5 || c < 0 || c > 6 || g.Grid[r][c] != symbol {
				break
			}
			count++
		}
		for i := 1; i < 4; i++ {
			r, c := row-d[0]*i, col-d[1]*i
			if r < 0 || r > 5 || c < 0 || c > 6 || g.Grid[r][c] != symbol {
				break
			}
			count++
		}
		if count >= 4 {
			return true
		}
	}
	return false
}

func (g *Game) Reset() {
	g.Grid = [6][7]int{}
	g.CurrentPlayer = 1
	g.Winner = 0
	g.TurnCount = 0
}

func (g *Game) PrintGrid() {
	for _, row := range g.Grid {
		for _, cell := range row {
			if cell == 0 {
				fmt.Print(".")
			} else if cell == 1 {
				fmt.Print("R")
			} else {
				fmt.Print("J")
			}
			fmt.Print(" ")
		}
		fmt.Println()
	}
	fmt.Println()
}
