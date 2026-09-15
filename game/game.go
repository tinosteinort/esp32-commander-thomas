package game

import (
	"blinky/hardware"
)

type GameState int

const (
	ShowInitScreen = iota
	PlayGame       = iota
	GameOverScreen = iota
)

type Game struct {
	hardware                hardware.Hardware
	initScreenController    *InitScreenController
	playScreenController    *PlayScreenController
	gameOverScreeController *GamveOverScreenController
	state                   GameState
}

type spaceShip struct {
	x      int16
	y      int16
	w      int16
	h      int16
	deltaY int16
}

type bullet struct {
	x      int16
	y      int16
	w      int16
	h      int16
	deltaX int16
}

type enemy struct {
	x      int16
	y      int16
	w      int16
	h      int16
	deltaX int16
}

type playScreenResult struct {
	screenDone bool
	points     int16
}

func New(h hardware.Hardware) *Game {
	g := &Game{
		hardware: h,
	}
	g.ShowInitScreen()
	return g
}

func (g *Game) Control() {

	switch g.state {
	case ShowInitScreen:
		done := g.initScreenController.Control()
		if done {
			g.Play()
		}
	case PlayGame:
		result := g.playScreenController.Control()
		if result.screenDone {
			g.ShowGameOverScreen(result.points)
		}
	case GameOverScreen:
		done := g.gameOverScreeController.Control()
		if done {
			g.ShowInitScreen()
		}
	}
}

func (g *Game) ShowInitScreen() {
	g.state = ShowInitScreen
	g.initScreenController = NewInitScreenController(g.hardware)
	g.playScreenController = nil
	g.gameOverScreeController = nil
}

func (g *Game) Play() {
	g.state = PlayGame
	g.initScreenController = nil
	g.playScreenController = NewPlayScreenController(g.hardware)
	g.gameOverScreeController = nil
}

func (g *Game) ShowGameOverScreen(points int16) {
	g.state = GameOverScreen
	g.initScreenController = nil
	g.playScreenController = nil
	g.gameOverScreeController = NewGameOverController(g.hardware, points)
}
