package TicTacToe

type Tile struct {
	xMin      int
	xMax 	  int
	yMin      int
	yMax 	  int
	status 	  TileStatus
}

type TileStatus int

const (
	inactive TileStatus    = -1
	nought                 = 0
	cross                  = 1
)

var gameBoardMatrix [3][3]Tile

