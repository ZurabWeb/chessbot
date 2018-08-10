// Package game contains all the logic for creating and manipulating a game
package game

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/notnil/chess"
)

// Challenge represents a challenge between two players
type Challenge struct {
	ChallengerID string
	ChallengedID string
	GameID       string
	ChannelID    string
}

type Color string

const (
	White Color = "White"
	Black Color = "Black"
)

// Player represents a human Chess player
type Player struct {
	ID string
}

// Game is the state of a game (active or not)
type Game struct {
	game        *chess.Game
	Players     map[Color]Player
	started     bool
	lastMove    *chess.Move
	checkedTile *chess.Square
}

// NewGame will create a new game with typical starting positions
func NewGame(players ...Player) *Game {
	gm := &Game{
		game: chess.NewGame(chess.UseNotation(chess.LongAlgebraicNotation{})),
	}
	attachPlayers(gm, players...)
	return gm
}

func attachPlayers(game *Game, players ...Player) {
	randomOrder := make([]Player, len(players))
	perm := rand.Perm(len(players))
	for i, v := range perm {
		randomOrder[v] = players[i]
	}
	game.Players = make(map[Color]Player)
	game.Players[White] = randomOrder[0]
	game.Players[Black] = randomOrder[1]
}

// NewGameFromFEN will create a new game with a given FEN starting position
func NewGameFromFEN(fen string, players ...Player) (*Game, error) {
	gameState, err := chess.FEN(fen)
	if err != nil {
		return &Game{}, err
	}
	game := &Game{
		game:    chess.NewGame(gameState, chess.UseNotation(chess.LongAlgebraicNotation{})),
		started: true,
	}
	attachPlayers(game, players...)
	return game, nil
}

func NewGameFromPGN(pgn string, white Player, black Player) (*Game, error) {
	reader := strings.NewReader(pgn)
	gameState, err := chess.PGN(reader)
	if err != nil {
		return &Game{}, err
	}
	game := &Game{
		game: chess.NewGame(gameState, chess.UseNotation(chess.LongAlgebraicNotation{})),
	}
	game.Players = make(map[Color]Player)
	game.Players[White] = white
	game.Players[Black] = black
	return game, nil
}

// TurnPlayer returns which player should move next
func (g *Game) TurnPlayer() Player {
	return g.Players[g.Turn()]
}

// Turn returns which color should move next
func (g *Game) Turn() Color {
	switch g.game.Position().Turn() {
	case chess.White:
		return White
	case chess.Black:
		return Black
	default:
		return White
	}
}

// FEN serializer
func (g *Game) FEN() string {
	return g.game.FEN()
}

// PGN serializer
func (g *Game) PGN() string {
	return g.game.String()
}

// Outcome determines the outcome of the game (or no outcome)
func (g *Game) Outcome() chess.Outcome {
	return g.game.Outcome()
}

// ResultText will show the outcome of the game in textual format
func (g *Game) ResultText() string {
	return fmt.Sprintf("Game completed. %s by %s.", g.Outcome(), g.game.Method())
}

// LastMove returns the last move done of the game
func (g *Game) LastMove() *chess.Move {
	moves := g.game.Moves()
	if len(moves) == 0 {
		return nil
	}
	return moves[len(moves)-1]
}

// Move a Chess piece based on standard algabreic notation (d2d4, etc)
func (g *Game) Move(san string) (*chess.Move, error) {
	err := g.game.MoveStr(san)
	if err != nil {
		return nil, err
	}
	g.started = true
	return g.LastMove(), nil
}

// Start indicates the game has been started
func (g *Game) Start() {
	g.started = true
}

// Started determines if the game has been started
func (g *Game) Started() bool {
	return g.started
}

// ValidMoves returns a list of all moves available to the current player's turn
func (g *Game) ValidMoves() []*chess.Move {
	return g.game.ValidMoves()
}

func (g *Game) CheckedKing() chess.Square {
	squareMap := g.game.Position().Board().SquareMap()
	lastMovePiece := squareMap[g.LastMove().S2()]
	for square, piece := range squareMap {
		if piece.Type() == chess.King && piece.Color() == lastMovePiece.Color().Other() {
			return square
		}
	}
	return chess.NoSquare
}

// String representation of the current game state (draws an ascii board)
func (g *Game) String() string {
	return g.game.Position().Board().Draw()
}
