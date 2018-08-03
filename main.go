package main

import (
	"bufio"
	"fmt"
	"os"
    "strconv"
	"baumfalk/reversi/game"
)

type Point struct {
    x int
    y int
}


func getNumberFromIO(scanner * bufio.Scanner, min int, max int) int {
    scanner.Scan()
    text := scanner.Text()
    num, error := strconv.Atoi(text)
    for; error != nil;{
        fmt.Println(error)
        scanner.Scan()
        text = scanner.Text()
        num, error = strconv.Atoi(text)
    }
    return num
}

func checkLegalMove(y int, x int, board [][] string, bluesTurn bool, dir_y int, dir_x int) (bool, [] Point) {
    otherSymbol := "O"
    ownSymbol := "X"
    if !bluesTurn {
        otherSymbol = "X"
        ownSymbol = "O"
    }

    yCur := y + dir_y
    if yCur >= len(board) || yCur < 0 {
        return false, nil
    }

    xCur := x + dir_x
    if xCur >= len(board[y]) || xCur < 0 {
        return false, nil
    }

    validStreak := false
    var list []Point
    for ; board[yCur][xCur] == otherSymbol; yCur, xCur = yCur + dir_y, xCur + dir_x {
        validStreak = true
        
        if yCur+dir_y >= len(board) || yCur+dir_y < 0 {
            return false, nil
        }
        if xCur+dir_x >= len(board[yCur]) || xCur+dir_x < 0 {
            return false, nil
        }
        list = append(list,Point{y:yCur, x:xCur })
    } 
    valid := validStreak && board[yCur][xCur] == ownSymbol
    if !valid {
        list = nil
    }
    return valid, list
}

func horizontallyLegal(y int, x int, board [][] string, bluesTurn bool) (bool,[]Point) {
    l2r, l1 := checkLegalMove(y, x, board, bluesTurn, 0, 1)
    r2l, l2 := checkLegalMove(y, x, board, bluesTurn, 0, -1)

    return l2r || r2l, append(l1,l2...)
}
func verticallyLegal(y int, x int, board [][] string, bluesTurn bool) (bool,[]Point) {
    t2b, l1 := checkLegalMove(y, x, board, bluesTurn, 1, 0)
    b2t, l2 := checkLegalMove(y, x, board, bluesTurn, -1, 0)

    return t2b || b2t, append(l1,l2...)
}

func diagonallyLegal(y int, x int, board [][] string, bluesTurn bool) (bool,[]Point) {
    t2bl2r, l1 := checkLegalMove(y, x, board, bluesTurn, -1, 1)
    t2br2l, l2 := checkLegalMove(y, x, board, bluesTurn, -1, -1)
    b2tl2r, l3 := checkLegalMove(y, x, board, bluesTurn, 1, 1)
    b2tr2l, l4 := checkLegalMove(y, x, board, bluesTurn, 1, -1)

    totalList := append(append(l1,l2...),append(l3,l4...)...)
    return t2bl2r || t2br2l || b2tl2r || b2tr2l, totalList
}


func generateLegalMoves(board [][] string, bluesTurn bool) ([][] bool, int) {
    legalMoveBoard := make([][]bool, len(board), len(board)) 
    legalMoveCount := 0
    for y, row := range board {
        legalMoveBoard[y] = make([]bool, len(board[y]), len(board[y])) 

        for x, _ := range row {
            legalMoveBoard[y][x], _, _ = isLegalMove(y, x, board, bluesTurn)
            if legalMoveBoard[y][x] {
                legalMoveCount++
            }
        }
    }
    return legalMoveBoard, legalMoveCount
}

func isLegalMove(y int, x int, board [][]string, bluesTurn bool) (bool, string, [] Point) {
    outOfBounds := y < 0 || y >= len(board) || x <0 || x >= len(board[y])
    if outOfBounds {
        return false, "Out of bounds", nil
    }
    cellIsOccupied := board[y][x] != " "
    if cellIsOccupied || outOfBounds{
        return false, "Cell is occupied", nil
    }

    h, l1 := horizontallyLegal(y, x, board, bluesTurn)
    v, l2 := verticallyLegal(y, x, board, bluesTurn)
    d, l3 := diagonallyLegal(y, x, board, bluesTurn)
    totalList := append(append(l1,l2...),l3...)
    return h || v || d, "no streak", totalList
}

func getCoordinates(scanner * bufio.Scanner, boardSize int) (y int, x int) {
    fmt.Println("Type y coordinate")
    y = getNumberFromIO(scanner,0, boardSize-1)
    fmt.Println("Type x coordinate")
    x = getNumberFromIO(scanner,0, boardSize-1)
    return
}

func printInfo(bluesTurn bool) {
    if bluesTurn {
        fmt.Println("It's blue's turn!")
    } else {
        fmt.Println("It's reds turn!")
    }
}

func changeBoard(y int, x int, board [][]string, changedCells [] Point, bluesTurn bool) {
    ownSymbol := "X"
    if !bluesTurn {
        ownSymbol = "O"
    }

    for _, pt := range changedCells {
        board[pt.y][pt.x] = ownSymbol
    }
        board[y][x] = ownSymbol

}

// func changeBoard(y int, x int, board [][]string, bluesTurn bool) {
//     otherSymbol := "O"
//     ownSymbol := "X"
//     if !bluesTurn {
//         otherSymbol = "X"
//         ownSymbol = "O"
//     }
//     if horizontallyLegal(y, x, board, bluesTurn) {
//         dir := 1
//         if checkLegalMove(y, x, board, bluesTurn, 0, -1) {
//             dir = -1
//         } 
//         xCur := x + dir
//         stillCool := board[y][xCur] == otherSymbol
//         for ; stillCool; xCur += dir {
//             stillCool = board[y][xCur] == otherSymbol
//             board[y][xCur] = ownSymbol 
//         }
//     } 
//     if verticallyLegal(y, x, board, bluesTurn) {
//         dir := 1
//         if checkLegalMove(y, x, board, bluesTurn, -1, 0) {
//             dir = -1
//         } 
//         yCur := y + dir
//         stillCool := board[yCur][x] == otherSymbol
//         for ; stillCool; yCur += dir {
//             stillCool = board[yCur][x] == otherSymbol
//             board[yCur][x] = ownSymbol 
//         }
//     } 
//     if diagonallyLegal(y, x, board, bluesTurn) {
//         y_dir := 1
//         x_dir := 1
//         if checkLegalMove(y, x, board, bluesTurn, -1, -1) {
//             y_dir = -1
//             x_dir = -1
//         } else if checkLegalMove(y, x, board, bluesTurn, 1, -1) {
//             x_dir = -1
//         } else if checkLegalMove(y, x, board, bluesTurn, -1, 1) {
//             y_dir = -1
//         }

//         yCur := y + y_dir
//         xCur := x + x_dir
//         stillCool := board[yCur][xCur] == otherSymbol
//         for ; stillCool; yCur, xCur = yCur + y_dir, xCur + x_dir {
//             stillCool = board[yCur][xCur] == otherSymbol
//             board[yCur][xCur] = ownSymbol 
//         }
//     }

//     board[y][x] = ownSymbol
// }

func determineWinner(board [][]string) string {
    countX := 0
    countO := 0
    for _, row := range board {
        for _, cell := range row {
            if cell == "X" {
                countX++
            } else if cell == "O" {
                countO++
            }
        }
    }

    if countX > countO {
        return "player 1"
    } else if countX < countO {
        return "player 2"
    } else {
        return "draw"
    }
}

func main() {
	fmt.Println("Starting reversi")
	boardSize := 6
	board := game.NewReversiBoard(boardSize)
	scanner := bufio.NewScanner(os.Stdin)

    bluesTurn := true
    passes := 0
	for gameDone := false; !gameDone; {
        legalMoveBoard, legalMoveCount := generateLegalMoves(board, bluesTurn)
        if legalMoveCount == 0 {
            passes++
            fmt.Println("No possible options, passing")
            bluesTurn = !bluesTurn

            if passes == 2 {
                gameDone = true
                fmt.Println("Game over, determining winner!")
            }
            continue
        }
        passes = 0
        game.PrintBoard(board, legalMoveBoard)
        printInfo(bluesTurn)               
        y, x := getCoordinates(scanner, boardSize)
        legalMove, reason, changes := isLegalMove(y,x,board, bluesTurn)

        if legalMove {
            changeBoard(y,x,board, changes, bluesTurn)    
            bluesTurn = !bluesTurn
        } else {
            fmt.Println("This (", y,", ", x, ") is not a legal move:", reason)
        }
	}

    winner := determineWinner(board)
    fmt.Println("winner is: ", winner)

    fmt.Println("Quitting reversi")
}
