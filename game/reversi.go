package game

import (
    "fmt"
    "github.com/fatih/color"
)
/*
 * To create an empty reversiboard
 */
func NewReversiBoard(size int) [][]string {
    board := make([][]string, size, size)
    for i := 0; i < size; i++ {
        board[i] = make([]string, size, size)
    }
    color.Cyan("Prints text in cyan.")

    for i := 0; i < size; i++ {
        for j := 0; j < size; j++ {
            board[i][j] = " "
        }
    }

    board[2][2] ="X"
    board[2][3] ="O"
    board[3][2] ="O"
    board[3][3] ="X"

	return board
}

func PrintBoard(board [][] string, legalMoveBoard [][] bool) {
    
    fmt.Print("  ")
    for i, _ := range board {
        fmt.Print(i)
    }
    fmt.Println()
    for y, row := range board {
        fmt.Print(y," ")
        for x, cell := range row {
            if cell == "X" {
                color.Set(color.FgBlue)
                color.Set(color.BgWhite)
                fmt.Print("X")
                color.Unset()
            } else if cell == "O" {
                color.Set(color.FgRed)
                color.Set(color.BgWhite)
                fmt.Print("O")    
                color.Unset()    
            } else if legalMoveBoard[y][x]{
                color.Set(color.BgWhite)
                fmt.Print("*")
                color.Unset()
            } else {
                color.Set(color.BgWhite)
                fmt.Print(cell)
                color.Unset()
            }
        }
        
        fmt.Println()
        //fmt.Println(row)
    }
}
