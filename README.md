# TicTacToe-GO

The create-functions makes an object with -2 subtracted from the height and width, and the draw-functions adds 1 to the position. This is to avoid painting over the lines that separate the tiles in the game board.

createGameBoard() has a +1 to height and width to have all the edge lines visible

defer "xxx"Buffer.Release() is to avoid running out of memory (ish?)

The static components, such as a cross, nought and gameboard has one function for creation and one for drawing. Dynamic components has one function that does all of this as they have to be recreated regularly.