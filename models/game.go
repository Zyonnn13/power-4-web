package models

import (
	"fmt"
)

const (
	Rows = 6
	Cols = 7
)

// Structure Player
type Player struct {
	Name  string
	Color string
}

// Structure Game avec "Players" (la liste)
type Game struct {
	Grid    [Rows][Cols]string
	Players [2]Player
	Turn    int
	Status  string
	Winner  string
}

func NewGame() *Game {
	return &Game{
		Status: "playing",
	}
}

func (g *Game) ConfigurePlayers(p1Name, p2Name, p1Color, p2Color string) {
	g.Players[0] = Player{Name: p1Name, Color: p1Color}
	g.Players[1] = Player{Name: p2Name, Color: p2Color}
}

func (g *Game) DropToken(col int) bool {
	if g.Status != "playing" || col < 0 || col >= Cols {
		return false
	}

	for r := Rows - 1; r >= 0; r-- {
		if g.Grid[r][col] == "" {
			currentColor := g.Players[g.Turn].Color
			g.Grid[r][col] = currentColor

			if g.CheckWin(r, col, currentColor) {
				g.Status = "win"
				g.Winner = g.Players[g.Turn].Name
				AddRecord(g) // Sauvegarde
			} else if g.CheckDraw() {
				g.Status = "draw"
				AddRecord(g) // Sauvegarde
			} else {
				g.Turn = (g.Turn + 1) % 2
			}
			return true
		}
	}
	return false
}

func (g *Game) CheckDraw() bool {
	for c := 0; c < Cols; c++ {
		if g.Grid[0][c] == "" {
			return false
		}
	}
	return true
}

func (g *Game) CheckWin(row, col int, color string) bool {
	directions := [][2]int{{0, 1}, {1, 0}, {1, 1}, {1, -1}}
	for _, d := range directions {
		count := 1
		for i := 1; i < 4; i++ {
			r, c := row+d[0]*i, col+d[1]*i
			if r >= 0 && r < Rows && c >= 0 && c < Cols && g.Grid[r][c] == color {
				count++
			} else {
				break
			}
		}
		for i := 1; i < 4; i++ {
			r, c := row-d[0]*i, col-d[1]*i
			if r >= 0 && r < Rows && c >= 0 && c < Cols && g.Grid[r][c] == color {
				count++
			} else {
				break
			}
		}
		if count >= 4 {
			return true
		}
	}
	return false
}

func (g *Game) String() string {
	return fmt.Sprintf("Tour de %s (%s)", g.Players[g.Turn].Name, g.Players[g.Turn].Color)
}
