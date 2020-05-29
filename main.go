package main

// Developed by Øyvind Reppen Lunde, May 2020.

import (
	"./game"
	"./display"
)

func main() {
	TileClicked := make(chan game.Position)
	ResetChannel := make(chan bool)

	go game.PlayTicTacToe(TileClicked, ResetChannel)
	go display.DisplayGame(TileClicked, ResetChannel)

	select{}
}