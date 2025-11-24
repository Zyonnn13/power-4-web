package models

const (
	Rows = 6
	Cols = 7

	emptyCell   = 0
	playerOne   = 1
	playerTwo   = 2
	winTarget   = 4
	totalSpaces = Rows * Cols
)

// Position résume la dernière action effectuée sur la grille.
type Position struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

// Game maintient l'état complet d'une partie de Puissance 4.
type Game struct {
	Board         [Rows][Cols]int `json:"board"`
	CurrentPlayer int             `json:"currentPlayer"`
	Winner        int             `json:"winner"` // 0 si pas de gagnant / match nul
	Finished      bool            `json:"finished"`
	Moves         int             `json:"moves"`
	LastMove      Position        `json:"lastMove"`
	Player1Name   string          `json:"player1Name"`
	Player2Name   string          `json:"player2Name"`
	Player1Color  string          `json:"player1Color"`
	Player2Color  string          `json:"player2Color"`
}

// NewGame crée une nouvelle partie avec une grille vide et le joueur 1 qui commence.
func NewGame() *Game {
	return &Game{
		CurrentPlayer: playerOne,
		Player1Name:   "Joueur 1",
		Player2Name:   "Joueur 2",
		Player1Color:  "red",
		Player2Color:  "yellow",
	}
}

// DropToken tente de déposer un jeton dans la colonne spécifiée.
// Retourne false si la colonne est pleine ou si la partie est déjà terminée.
func (g *Game) DropToken(col int) bool {
	if g == nil || g.Finished || col < 0 || col >= Cols {
		return false
	}

	row := g.findAvailableRow(col)
	if row == -1 {
		return false
	}

	g.Board[row][col] = g.CurrentPlayer
	g.Moves++
	g.LastMove = Position{Row: row, Col: col}

	if g.hasConnectFour(row, col) {
		g.Winner = g.CurrentPlayer
		g.Finished = true
		return true
	}

	if g.Moves == totalSpaces {
		g.Finished = true // match nul
		return true
	}

	g.switchPlayer()
	return true
}

func (g *Game) findAvailableRow(col int) int {
	for row := Rows - 1; row >= 0; row-- {
		if g.Board[row][col] == emptyCell {
			return row
		}
	}
	return -1
}

func (g *Game) switchPlayer() {
	if g.CurrentPlayer == playerOne {
		g.CurrentPlayer = playerTwo
	} else {
		g.CurrentPlayer = playerOne
	}
}

func (g *Game) hasConnectFour(row, col int) bool {
	player := g.CurrentPlayer
	directions := [][2]int{
		{0, 1},  // horizontal
		{1, 0},  // vertical
		{1, 1},  // diagonale descendante
		{1, -1}, // diagonale ascendante
	}

	for _, dir := range directions {
		count := 1
		count += g.countDirection(row, col, dir[0], dir[1], player)
		count += g.countDirection(row, col, -dir[0], -dir[1], player)
		if count >= winTarget {
			return true
		}
	}
	return false
}

func (g *Game) countDirection(row, col, dRow, dCol, player int) int {
	total := 0
	r, c := row+dRow, col+dCol
	for r >= 0 && r < Rows && c >= 0 && c < Cols && g.Board[r][c] == player {
		total++
		r += dRow
		c += dCol
	}
	return total
}

// ConfigurePlayers permet d'associer noms et couleurs aux joueurs.
func (g *Game) ConfigurePlayers(p1Name, p2Name, p1Color, p2Color string) {
	if g == nil {
		return
	}
	g.Player1Name = sanitizeName(p1Name, "Joueur 1")
	g.Player2Name = sanitizeName(p2Name, "Joueur 2")
	g.Player1Color = sanitizeColor(p1Color, "red")
	g.Player2Color = sanitizeColor(p2Color, "yellow")
}

func sanitizeName(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func sanitizeColor(value, fallback string) string {
	switch value {
	case "red", "yellow":
		return value
	default:
		return fallback
	}
}

// PlayerName renvoie le nom du joueur correspondant à l'identifiant.
func (g *Game) PlayerName(id int) string {
	if id == playerTwo {
		return g.Player2Name
	}
	return g.Player1Name
}

// PlayerColor renvoie la couleur associée au joueur.
func (g *Game) PlayerColor(id int) string {
	if id == playerTwo {
		return g.Player2Color
	}
	return g.Player1Color
}
