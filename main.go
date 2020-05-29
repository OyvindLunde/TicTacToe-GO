package main

// Developed by Øyvind Reppen Lunde, May 2020.

import (
	"./game"
	"./display"
)

func main() {
	TileClicked := make(chan game.Position)
	ResetChannel := make(chan bool)
	UpdateDisplay := make(chan bool)

	go game.PlayTicTacToe(TileClicked, ResetChannel, UpdateDisplay)
	go display.DisplayGame(TileClicked, ResetChannel, UpdateDisplay)

	select{}
}