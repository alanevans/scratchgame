package main

import (
	"log"

	"scratchgame/internal/game"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("2D Platform Game")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := game.NewGame(screenWidth, screenHeight)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
