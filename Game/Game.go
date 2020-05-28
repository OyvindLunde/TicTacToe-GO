package game

import (
	"fmt"
	"strconv"
)

type Tile struct {
	XMin      int
	XMax 	  int
	YMin      int
	YMax 	  int
	Status 	  TileStatus
}

type TileStatus int

const (
	Inactive TileStatus    = -1
	Nought                 = 1
	Cross                  = 2
)

type Position struct {
    Row int
    Col int
}

const numRows = 3
const numCols = 3

var gameBoardX0 = 150 						// "Start" position of Button Panel (x and y coordinate, top left corner)
var gameBoardY0 = 70
var gameBoardWidth = 240                    // Width of Button Panel
var gameBoardHeight = 240                   // Height of Button Panel
var tileSizeX = gameBoardWidth / numCols 	// Width of button in the Button Panel
var tileSizeY = gameBoardHeight / numRows  	// Height of button in the Button Panel

var GameBoardMatrix [numRows][numCols]Tile
var CurrentPlayer = 1
var Winner = -1
var UpdateDisplay = false

func initializeGameBoardMatrix() {
	for row := 0; row < numRows; row++ {
		for col := 0; col < numCols; col++ {
			GameBoardMatrix[row][col].XMin = row*tileSizeX + 1
			GameBoardMatrix[row][col].XMax = (row+1)*tileSizeX - 1
			GameBoardMatrix[row][col].YMin = col*tileSizeY + 1
			GameBoardMatrix[row][col].YMax = (col+1)*tileSizeY - 1
			GameBoardMatrix[row][col].Status = Inactive
		}
	}
}

func resetGame() {
	Winner = -1
	initializeGameBoardMatrix()
	UpdateDisplay = true
}

func executePlayerTurn(row int, col int) {
	if Winner != -1 || CheckForDraw() {
		return
	} 
	// Invalid move made (clicked outside the gameBoard, or on a non-empty tile)
	if row == -1 || col == -1 || GameBoardMatrix[row][col].Status != Inactive {
		fmt.Println("Invalid move, please click on an empty tile")
		return

	} else { // Valid move made
		GameBoardMatrix[row][col].Status = TileStatus(CurrentPlayer)
		changeCurrentPlayer()
		Winner = int(checkForVictory())
		if Winner != -1 {
			fmt.Println("Player " + strconv.Itoa(Winner) + " has won!")
		} else if CheckForDraw() {
			fmt.Println("Game ended in a draw")
		}
		// w.Send(paint.Event{})
		UpdateDisplay = true
	}
}

func changeCurrentPlayer() {
	if CurrentPlayer == 1 {
		CurrentPlayer = 2
	} else {
		CurrentPlayer = 1
	}
}

func checkForVictory() TileStatus {
	for row := 0; row < numRows; row++ {
		if GameBoardMatrix[row][0].Status == GameBoardMatrix[row][1].Status && GameBoardMatrix[row][0].Status == GameBoardMatrix[row][2].Status && GameBoardMatrix[row][0].Status != -1 {
			return GameBoardMatrix[row][0].Status
		}
	}

	for col := 0; col < numCols; col++ {
		if GameBoardMatrix[0][col].Status == GameBoardMatrix[1][col].Status && GameBoardMatrix[0][col].Status == GameBoardMatrix[2][col].Status && GameBoardMatrix[0][col].Status != -1 {
			return GameBoardMatrix[0][col].Status
		}
	}

	if GameBoardMatrix[0][0].Status == GameBoardMatrix[1][1].Status && GameBoardMatrix[0][0].Status == GameBoardMatrix[2][2].Status && GameBoardMatrix[0][0].Status != -1 {
		return GameBoardMatrix[0][0].Status
	}

	if GameBoardMatrix[0][2].Status == GameBoardMatrix[1][1].Status && GameBoardMatrix[0][2].Status == GameBoardMatrix[2][0].Status && GameBoardMatrix[0][2].Status != -1 {
		return GameBoardMatrix[0][2].Status
	}

	return -1
}

func CheckForDraw() bool {
	if Winner != -1 {
		return false
	}

	for row := 0; row < numRows; row++ {
		for col := 0; col < numCols; col++ {
			if GameBoardMatrix[row][col].Status == -1 {
				return false
			}
		}
	}
	return true
}

func PlayTicTacToe(tileClicked chan Position, ResetChannel chan bool) {
	initializeGameBoardMatrix()

	for {
		select {
		case position := <-tileClicked:
			executePlayerTurn(position.Row, position.Col)

		case <- ResetChannel:
			if Winner != -1 || CheckForDraw() {
				resetGame()
			}
		}
	}
	
}