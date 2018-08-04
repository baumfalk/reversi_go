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

func posIsOutOfBounds(yPos int, yLength int, xPos int, xLength int) bool {
	return yPos >= yLength || yPos < 0 || xPos >= xLength || xPos < 0
}

func checkLegalMoveHelper(y int, x int, board [][]int, curPlayer int, dirY int, dirX int) []Point {
	otherPlayer := 3 - curPlayer
	yCur := y + dirY
	xCur := x + dirX
	var list []Point
	var validStreak bool

	if posIsOutOfBounds(yCur, len(board), xCur, len(board[y])) {
		return nil
	}

	for board[yCur][xCur] == otherPlayer {
		validStreak = true
		if posIsOutOfBounds(yCur+dirY, len(board), xCur+dirX, len(board[yCur])) {
			return nil
		}
		list = append(list, Point{y: yCur, x: xCur})
		yCur, xCur = yCur+dirY, xCur+dirX
	}
	validStreak = validStreak && board[yCur][xCur] == curPlayer

	if !validStreak {
		list = nil
	}
	return list
}

func horizontallyLegal(y int, x int, board [][]int, curPlayer int, ch chan []Point) {
	l2r := checkLegalMoveHelper(y, x, board, curPlayer, 0, 1)
	r2l := checkLegalMoveHelper(y, x, board, curPlayer, 0, -1)

	ch <- append(l2r, r2l...)
}
func verticallyLegal(y int, x int, board [][]int, curPlayer int, ch chan []Point) {
	t2b := checkLegalMoveHelper(y, x, board, curPlayer, 1, 0)
	b2t := checkLegalMoveHelper(y, x, board, curPlayer, -1, 0)

	ch <- append(t2b, b2t...)
}

func diagonallyLegal(y int, x int, board [][]int, curPlayer int, ch chan []Point) {
	t2bl2r := checkLegalMoveHelper(y, x, board, curPlayer, -1, 1)
	t2br2l := checkLegalMoveHelper(y, x, board, curPlayer, -1, -1)
	b2tl2r := checkLegalMoveHelper(y, x, board, curPlayer, 1, 1)
	b2tr2l := checkLegalMoveHelper(y, x, board, curPlayer, 1, -1)

	totalList := append(append(t2bl2r, t2br2l...), append(b2tl2r, b2tr2l...)...)
	ch <- totalList
}

func GenerateLegalMoves(board [][]int, curPlayer int) ([][][]Point, int) {
	legalMoveBoard := make([][][]Point, len(board), len(board))
	legalMoveCount := 0

	for y, row := range board {
		legalMoveBoard[y] = make([][]Point, len(board[y]), len(board[y]))

		for x := range row {
			moveResult, _ := IsLegalMove(y, x, board, curPlayer)
			legalMoveBoard[y][x] = moveResult
			if moveResult != nil {
				legalMoveCount++
			}
		}
	}

	return legalMoveBoard, legalMoveCount
}

func IsLegalMove(y int, x int, board [][]int, curPlayer int) ([]Point, string) {
	outOfBounds := y < 0 || y >= len(board) || x < 0 || x >= len(board[y])
	if outOfBounds {
		return nil, "Out of bounds"
	}
	cellIsOccupied := board[y][x] != 0
	if cellIsOccupied {
		return nil, "Cell is occupied"
	}

	ch1, ch2, ch3 := make(chan []Point), make(chan []Point), make(chan []Point)

	go horizontallyLegal(y, x, board, curPlayer, ch1)
	go verticallyLegal(y, x, board, curPlayer, ch2)
	go diagonallyLegal(y, x, board, curPlayer, ch3)

	h, v, d := <-ch1, <-ch2, <-ch3

	totalList := append(append(h, v...), d...)

	return totalList, "no streak"
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
		fmt.Println("This (", *y, ", ", *x, ") is not a legal move")
	}
}
