package game

import (
	"fmt"

	"github.com/fatih/color"
)

type Point struct {
	x int
	y int
}

/*
 * To create an empty reversiboard
 */
func NewReversiBoard(size int) [][]int {
	board := make([][]int, size, size)
	for i := 0; i < size; i++ {
		board[i] = make([]int, size, size)
	}

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			board[i][j] = 0
		}
	}

	board[2][2] = 1
	board[2][3] = 2
	board[3][2] = 2
	board[3][3] = 1

	return board
}

/*
   To print a reversi board
*/

func PrintBoard(board [][]int, legalMoveBoard [][][]Point) {

	fmt.Print("  ")
	for i := range board {
		fmt.Print(i)
	}
	fmt.Println()
	for y, row := range board {
		fmt.Print(y, " ")
		for x, cell := range row {
			if cell == 1 {
				color.Set(color.FgBlue)
				color.Set(color.BgWhite)
				fmt.Print("X")
				color.Unset()
			} else if cell == 2 {
				color.Set(color.FgRed)
				color.Set(color.BgWhite)
				fmt.Print("O")
				color.Unset()
			} else if legalMoveBoard[y][x] != nil {
				color.Set(color.BgWhite)
				fmt.Print("*")
				color.Unset()
			} else {
				color.Set(color.BgWhite)
				fmt.Print(" ")
				color.Unset()
			}
		}

		fmt.Println()
	}
}

func checkLegalMoveHelper(y int, x int, board [][]int, curPlayer int, dirY int, dirX int) (bool, []Point) {
	otherPlayer := 3 - curPlayer

	yCur := y + dirY
	if yCur >= len(board) || yCur < 0 {
		return false, nil
	}

	xCur := x + dirX
	if xCur >= len(board[y]) || xCur < 0 {
		return false, nil
	}

	validStreak := false
	var list []Point
	for ; board[yCur][xCur] == otherPlayer; yCur, xCur = yCur+dirY, xCur+dirX {
		validStreak = true

		if yCur+dirY >= len(board) || yCur+dirY < 0 || xCur+dirX >= len(board[yCur]) || xCur+dirX < 0 {
			validStreak = false
			break
		}
		list = append(list, Point{y: yCur, x: xCur})
	}
	valid := validStreak && board[yCur][xCur] == curPlayer
	if !valid {
		list = nil
	}
	return valid, list
}

func horizontallyLegal(y int, x int, board [][]int, curPlayer int) (bool, []Point) {
	l2r, l1 := checkLegalMoveHelper(y, x, board, curPlayer, 0, 1)
	r2l, l2 := checkLegalMoveHelper(y, x, board, curPlayer, 0, -1)

	return l2r || r2l, append(l1, l2...)
}
func verticallyLegal(y int, x int, board [][]int, curPlayer int) (bool, []Point) {
	t2b, l1 := checkLegalMoveHelper(y, x, board, curPlayer, 1, 0)
	b2t, l2 := checkLegalMoveHelper(y, x, board, curPlayer, -1, 0)

	return t2b || b2t, append(l1, l2...)
}

func diagonallyLegal(y int, x int, board [][]int, curPlayer int) (bool, []Point) {
	t2bl2r, l1 := checkLegalMoveHelper(y, x, board, curPlayer, -1, 1)
	t2br2l, l2 := checkLegalMoveHelper(y, x, board, curPlayer, -1, -1)
	b2tl2r, l3 := checkLegalMoveHelper(y, x, board, curPlayer, 1, 1)
	b2tr2l, l4 := checkLegalMoveHelper(y, x, board, curPlayer, 1, -1)

	totalList := append(append(l1, l2...), append(l3, l4...)...)
	return t2bl2r || t2br2l || b2tl2r || b2tr2l, totalList
}

func GenerateLegalMoves(board [][]int, curPlayer int) ([][][]Point, int) {
	legalMoveBoard := make([][][]Point, len(board), len(board))

	legalMoveCount := 0

	for y, row := range board {
		legalMoveBoard[y] = make([][]Point, len(board[y]), len(board[y]))

		for x := range row {
			legal, _, changes := IsLegalMove(y, x, board, curPlayer)
			legalMoveBoard[y][x] = changes
			if legal {
				legalMoveCount++
			}
		}
	}
	return legalMoveBoard, legalMoveCount
}

func IsLegalMove(y int, x int, board [][]int, curPlayer int) (bool, string, []Point) {
	outOfBounds := y < 0 || y >= len(board) || x < 0 || x >= len(board[y])
	if outOfBounds {
		return false, "Out of bounds", nil
	}
	cellIsOccupied := board[y][x] != 0
	if cellIsOccupied || outOfBounds {
		return false, "Cell is occupied", nil
	}

	h, l1 := horizontallyLegal(y, x, board, curPlayer)
	v, l2 := verticallyLegal(y, x, board, curPlayer)
	d, l3 := diagonallyLegal(y, x, board, curPlayer)
	totalList := append(append(l1, l2...), l3...)
	return h || v || d, "no streak", totalList
}

func ChangeBoard(y int, x int, board [][]int, changedCells []Point, curPlayer int) {
	for _, pt := range changedCells {
		board[pt.y][pt.x] = curPlayer
	}
	board[y][x] = curPlayer
}

func HandlePass(passes *int, curPlayer *int, gameDone *bool) {
	*passes++
	fmt.Println("No possible options for player", *curPlayer, ", passing")
	*curPlayer = 3 - *curPlayer

	if *passes == 2 {
		*gameDone = true
		fmt.Println("Game over, determining winner!")
	}
}

func DetermineWinner(board [][]int) string {
	countP1 := 0
	countP2 := 0
	for _, row := range board {
		for _, cell := range row {
			if cell == 1 {
				countP1++
			} else if cell == 2 {
				countP2++
			}
		}
	}

	if countP1 > countP2 {
		return "player 1"
	} else if countP1 < countP2 {
		return "player 2"
	} else {
		return "draw"
	}
}

func HandleMove(y *int, x *int, legalMoveBoard *[][][]Point, curPlayer *int, board *[][]int) {
	legalMove := ((*legalMoveBoard)[*y][*x] != nil)
	if legalMove {
		changes := (*legalMoveBoard)[*y][*x]
		ChangeBoard(*y, *x, *board, changes, *curPlayer)
		*curPlayer = 3 - *curPlayer
	} else {
		fmt.Println("This (", y, ", ", x, ") is not a legal move")
	}
}
