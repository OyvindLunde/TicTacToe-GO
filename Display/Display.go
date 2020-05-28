package display

/* 
This module creates and displays the game, with its components and status
*/

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"
	"strconv"
	"time"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/imageutil"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/font"
	"golang.org/x/image/font/inconsolata"
	"golang.org/x/image/math/fixed"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/key"

	"../game"
)

var (
	black     = color.RGBA{0x00, 0x00, 0x00, 0x00}
	blue0     = color.RGBA{0x00, 0x00, 0x1f, 0xff}
	blue1     = color.RGBA{0x00, 0x00, 0x3f, 0xff}
	darkGray  = color.RGBA{0x3f, 0x3f, 0x3f, 0xff}
	lightGray = color.RGBA{0xd8, 0xd8, 0xd8, 0x7f}
	green     = color.RGBA{0x16, 0xee, 0x50, 0x7f}
	red       = color.RGBA{0xff, 0x00, 0x00, 0x7f}
	yellow    = color.RGBA{0xef, 0xff, 0x00, 0x3f}
	white     = color.RGBA{0xff, 0xff, 0xff, 0xff}
)

const numRows = 3
const numCols = 3

var gameBoardX0 = 150 						// "Start" position of Button Panel (x and y coordinate, top left corner)
var gameBoardY0 = 70
var gameBoardWidth = 240                    // Width of Button Panel
var gameBoardHeight = 240                   // Height of Button Panel
var tileSizeX = gameBoardWidth / numCols 	// Width of button in the Button Panel
var tileSizeY = gameBoardHeight / numRows  	// Height of button in the Button Panel


func DisplayGame(TileClicked chan game.Position, ResetChannel chan bool) {
	driver.Main(func(s screen.Screen) {
		// Create a window of desired size
		w, err := s.NewWindow(&screen.NewWindowOptions{ 
			Width:  600,
			Height: 500,
			Title:  "Tic Tac Toe",
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Let's play Tic Tac Toe")
		defer w.Release()

		go updateDisplay(w) 

		// Static components
		gameBoardShape := createEmptyGameBoard(s)
		noughtShape := createNought(s, tileSizeX, tileSizeY, red, white)
		crossShape := createCross(s, tileSizeX, tileSizeY, blue1, white)

		var sz size.Event
		for {
			e := w.NextEvent()
			switch e := e.(type) {

			case paint.Event:
				paintScreen(w, sz, lightGray, blue0)
				drawGameBoard(w, gameBoardShape, noughtShape, crossShape, gameBoardX0, gameBoardY0)
				drawGameStatusInfo(w, s)
				
			case size.Event: // Event that occurs when the window is resized by the user
				sz = e

			case mouse.Event:
				if e.Button == mouse.ButtonLeft && e.Direction == mouse.DirPress {
					row, col := findClickedTile(int(e.X), int(e.Y))
					TileClicked <- game.Position{Row: row, Col: col}
				}

			case key.Event:
				if e.Code == key.CodeR && e.Direction == key.DirPress {
					ResetChannel <- true
				}

			case error:
				log.Print(e)
			}
		}
	})
}

// Generate a paint event if the game board og game status has changed
func updateDisplay(w screen.EventDeque) {
	for {
		time.Sleep(20 * time.Millisecond)
		if game.UpdateDisplay {
			w.Send(paint.Event{})
			game.UpdateDisplay = false
		}
	}
}

func findClickedTile(x int, y int) (int, int) {
	xAdjusted := x - gameBoardX0
	yAdjusted := y - gameBoardY0
	for row := 0; row < numRows; row++ {
		for col := 0; col < numCols; col++ {
			if xAdjusted >= game.GameBoardMatrix[row][col].XMin && xAdjusted <= game.GameBoardMatrix[row][col].XMax {
				if yAdjusted >= game.GameBoardMatrix[row][col].YMin && yAdjusted <= game.GameBoardMatrix[row][col].YMax {
					return row, col
				}
			}
		}
	}
	
	return -1, -1
}

func drawGameStatusInfo(w screen.Window, s screen.Screen) {
	player1Info := drawText(s, "Player 1 is naughts (O)", 200, 20) 
	player2Info := drawText(s, "Player 2 is crosses (X)", 200, 20)
	w.Copy(image.Point{gameBoardX0, gameBoardY0+numRows*tileSizeY+10}, player1Info, player1Info.Bounds(), screen.Src, nil)
	w.Copy(image.Point{gameBoardX0, gameBoardY0+numRows*tileSizeY+30}, player2Info, player2Info.Bounds(), screen.Src, nil)

	if game.Winner != -1 {
		winnerInfo := drawText(s, "Player " + strconv.Itoa(game.Winner) + " has won!", 200, 20) 
		restartInfo := drawText(s, "Press 'R' to restart", 200, 20)
		w.Copy(image.Point{gameBoardX0, gameBoardY0-50}, winnerInfo, winnerInfo.Bounds(), screen.Src, nil)
		w.Copy(image.Point{gameBoardX0, gameBoardY0-30}, restartInfo, restartInfo.Bounds(), screen.Src, nil)
	} else if game.CheckForDraw() {
		drawInfo := drawText(s, "Game ended in a draw", 200, 20) 
		restartInfo := drawText(s, "Press 'R' to restart", 200, 20)
		w.Copy(image.Point{gameBoardX0, gameBoardY0-50}, drawInfo, drawInfo.Bounds(), screen.Src, nil)
		w.Copy(image.Point{gameBoardX0, gameBoardY0-30}, restartInfo, restartInfo.Bounds(), screen.Src, nil)
	} else {
		info := drawText(s, "It's player " + strconv.Itoa(game.CurrentPlayer) + "'s turn", 200, 20)
		w.Copy(image.Point{gameBoardX0, gameBoardY0-30}, info, info.Bounds(), screen.Src, nil)
	}
}

func drawCross(w screen.Window, cross screen.Texture, start_x int) {
	w.Copy(image.Point{start_x+1, gameBoardY0+1}, cross, cross.Bounds(), screen.Src, nil)
}

func createCross(s screen.Screen, width int, length int, color color.RGBA, backgroundColor color.RGBA) screen.Texture {
	crossRectangle := image.Point{width-2, length-2}
	crossBuffer, _ := s.NewBuffer(crossRectangle)
	pixelBuffer := crossBuffer.RGBA()
	bounds := crossBuffer.Bounds()

	deltaX := (bounds.Max.X - bounds.Min.X) / 10
	deltaY := (bounds.Max.Y - bounds.Min.Y) / 10

	// Paint the entire tile in the background color
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			pixelBuffer.SetRGBA(x, y, backgroundColor)
		}
	}

	// Paint a cross in the tile
	for x := bounds.Min.X + deltaX; x < bounds.Max.X - deltaX; x++ {
		for y := bounds.Min.Y + deltaY; y < bounds.Max.Y - deltaY; y++ {
			if y >=  x - deltaY && y <= x + deltaY {
				pixelBuffer.SetRGBA(x, y, color)
			} else if y >=  -x - deltaY + bounds.Max.Y && y <= -x + deltaY + bounds.Max.Y {
				pixelBuffer.SetRGBA(x, y, color)
			} else {
				pixelBuffer.SetRGBA(x, y, backgroundColor)
			}
		}
	}

	cross, _ := s.NewTexture(crossRectangle)
	cross.Upload(image.Point{}, crossBuffer, crossBuffer.Bounds())
	defer crossBuffer.Release()
	return cross
}

func drawNought(w screen.Window, nought screen.Texture, xTile int, yTile int) {
	w.Copy(image.Point{xTile+1, yTile+1}, nought, nought.Bounds(), screen.Src, nil)
}

// Nought is the Tic Tac Toe term for the circle draw on the playing board 
func createNought(s screen.Screen, width int, length int, color color.RGBA, backgroundColor color.RGBA) screen.Texture {
	noughtRectangle := image.Point{width-2, length-2}
	noughtBuffer, _ := s.NewBuffer(noughtRectangle)
	pixelBuffer := noughtBuffer.RGBA()
	bounds := noughtBuffer.Bounds()

	x0 := (bounds.Max.X - bounds.Min.X) / 2
	y0 := (bounds.Max.Y - bounds.Min.Y) / 2
	radiusInner := (bounds.Max.X - bounds.Min.X) * 3/10
	radiusOuter := (bounds.Max.X - bounds.Min.X) * 4/10

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			radius := int(math.Sqrt(math.Pow(float64(x-x0), 2) + math.Pow(float64(y-y0), 2)))
			if radius >= radiusInner && radius <= radiusOuter {
				pixelBuffer.SetRGBA(x, y, color)
			} else {
				pixelBuffer.SetRGBA(x, y, backgroundColor)
			}
		}
	}

	nought, _ := s.NewTexture(noughtRectangle)
	nought.Upload(image.Point{}, noughtBuffer, noughtBuffer.Bounds())
	defer noughtBuffer.Release()
	return nought
}

func drawGameBoard(w screen.Window, gameBoard screen.Texture, noughtShape screen.Texture, crossShape screen.Texture, xPos int, yPos int) {
	w.Copy(image.Point{xPos, yPos}, gameBoard, gameBoard.Bounds(), screen.Src, nil)

	for row := 0; row < numRows; row++ {
		for col := 0; col < numCols; col++ {
			if game.GameBoardMatrix[row][col].Status == game.Nought {
				drawNought(w, noughtShape, gameBoardX0 + row*tileSizeX, gameBoardY0 + col*tileSizeY)
			} else if game.GameBoardMatrix[row][col].Status == game.Cross {
				drawNought(w, crossShape, gameBoardX0 + row*tileSizeX, gameBoardY0 + col*tileSizeY)
			}
		}
	}
}

func createEmptyGameBoard(s screen.Screen) screen.Texture {
	gameBoardRectangle := image.Point{gameBoardWidth+1, gameBoardHeight+1}
	gameBoardBuffer, _ := s.NewBuffer(gameBoardRectangle)
	pixelBuffer := gameBoardBuffer.RGBA()
	bounds := gameBoardBuffer.Bounds()

	for i := bounds.Min.X; i < bounds.Max.X; i++ {
		for j := bounds.Min.Y; j < bounds.Max.Y; j++ {
			pixelBuffer.SetRGBA(i, j, white)
		}
	}

	drawHorizontalLines(pixelBuffer, numRows-1, black)
	drawVerticalLines(pixelBuffer, numCols-1, black)

	gameBoard, _ := s.NewTexture(gameBoardRectangle)
	gameBoard.Upload(image.Point{}, gameBoardBuffer, gameBoardBuffer.Bounds())
	defer gameBoardBuffer.Release()
	return gameBoard
}

func paintScreen(w screen.Window, sz size.Event, backgroundColor color.RGBA, borderColor color.RGBA) {
	const inset = 10
	for _, r := range imageutil.Border(sz.Bounds(), inset) {
		w.Fill(r, borderColor, screen.Src) // Paint border of screen
	}
	w.Fill(sz.Bounds().Inset(inset), backgroundColor, screen.Src) // Paint screen
}

// The most basic functions for drawing text and lines

func drawText(s screen.Screen, text string, x_size int, y_size int) screen.Texture {
	floor0 := image.Point{x_size, y_size}
	f0, err := s.NewBuffer(floor0)

	drawRGBA(f0.RGBA(), text)

	f01, err := s.NewTexture(floor0)
	if err != nil {
		log.Fatal(err)
	}
	f01.Upload(image.Point{}, f0, f0.Bounds())
	defer f0.Release()
	return f01
}

func drawRGBA(m *image.RGBA, str string) {
	draw.Draw(m, m.Bounds(), image.White, image.Point{}, draw.Src)

	d := font.Drawer{
		Dst:  m,
		Src:  image.Black,
		Face: inconsolata.Regular8x16,
		Dot: fixed.Point26_6{
			Y: inconsolata.Regular8x16.Metrics().Ascent,
		},
	}
	d.DrawString(str)
}

func drawHorizontalLines(m *image.RGBA, num int, color color.RGBA) {
	b := m.Bounds()
	intervall := (b.Max.Y - b.Min.Y) / (num + 1)
	for i := 0; i <= num+1; i++ {
		drawHorizontalLine(m, intervall*i, color)
	}
}

func drawHorizontalLine(m *image.RGBA, y int, color color.RGBA) {
	b := m.Bounds()
	for x := b.Min.X; x < b.Max.X; x++ {
		m.SetRGBA(x, y, color)
	}
}

func drawVerticalLines(m *image.RGBA, num int, color color.RGBA) {
	b := m.Bounds()
	intervall := (b.Max.X - b.Min.X) / (num + 1)
	for i := 0; i <= num+1; i++ {
		drawVerticalLine(m, intervall*i, color)
	}
}

func drawVerticalLine(m *image.RGBA, x int, color color.RGBA) {
	b := m.Bounds()

	for y := b.Min.Y; y < b.Max.Y; y++ {
		m.SetRGBA(x, y, color)
	}
}
