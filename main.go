package main

import (
	"bufio"
	"fmt"
	"os"
    "strconv"
	"baumfalk/reversi/game"
)


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



func hor_legal(y int, x int, board [][] string, bluesTurn bool, dir int) bool{
    otherSymbol := "O"
    ownSymbol := "X"
    if !bluesTurn {
        otherSymbol = "X"
        ownSymbol = "O"
    }

    xCur := x + dir
    if xCur >= len(board[y]) || xCur < 0 {
        return false
    }
    validStreak := false
    for ; board[y][xCur] == otherSymbol; xCur += dir {
        if xCur+dir >= len(board[y]) || xCur+dir < 0 {
            return false
        }
        validStreak = true
    } 
    return validStreak && board[y][xCur] == ownSymbol
}

func horizontallyLegal(y int, x int, board [][] string, bluesTurn bool) bool {
    //left to right
    //[X]OOOX
    l2r := hor_legal(y, x, board, bluesTurn, 1)
    r2l := hor_legal(y, x, board, bluesTurn, -1)

    return l2r || r2l
}

func ver_legal(y int, x int, board [][] string, bluesTurn bool, dir int) bool{
    otherSymbol := "O"
    ownSymbol := "X"
    if !bluesTurn {
        otherSymbol = "X"
        ownSymbol = "O"
    }

    yCur := y + dir
    if yCur >= len(board) || yCur < 0 {
        return false
    }
    validStreak := false
    for ; board[yCur][x] == otherSymbol; yCur += dir {
        if yCur+dir >= len(board) || yCur +dir< 0 {
            return false
        }
        validStreak = true
    } 
    return validStreak && board[yCur][x] == ownSymbol
}

func verticallyLegal(y int, x int, board [][] string, bluesTurn bool) bool {
    //left to right
    //[X]OOOX
    t2b := ver_legal(y, x, board, bluesTurn, 1)
    b2t := ver_legal(y, x, board, bluesTurn, -1)

    return t2b || b2t
}

func diag_legal(y int, x int, board [][] string, bluesTurn bool, dir_y int, dir_x int) bool{
    otherSymbol := "O"
    ownSymbol := "X"
    if !bluesTurn {
        otherSymbol = "X"
        ownSymbol = "O"
    }

    yCur := y + dir_y
    if yCur >= len(board) || yCur < 0 {
        return false
    }

    xCur := x + dir_x
    if xCur >= len(board[y]) || xCur < 0 {
        return false
    }

    validStreak := false
    for ; board[yCur][xCur] == otherSymbol; yCur, xCur = yCur + dir_y, xCur + dir_x {
        validStreak = true
        if yCur+dir_y >= len(board) || yCur+dir_y < 0 {
            return false
        }
        if xCur+dir_x >= len(board[yCur]) || xCur+dir_x < 0 {
            return false
        }
    } 
    return validStreak && board[yCur][xCur] == ownSymbol
}

func diagonallyLegal(y int, x int, board [][] string, bluesTurn bool) bool {
    //left to right
    //[X]OOOX
    t2bl2r := diag_legal(y, x, board, bluesTurn, -1, 1)
    t2br2l := diag_legal(y, x, board, bluesTurn, -1, -1)
    b2tl2r := diag_legal(y, x, board, bluesTurn, 1, 1)
    b2tr2l := diag_legal(y, x, board, bluesTurn, 1, -1)

    return t2bl2r || t2br2l || b2tl2r || b2tr2l
}


func generateLegalMoves(board [][] string, bluesTurn bool) ([][] bool, int) {
    legalMoveBoard := make([][]bool, len(board), len(board)) 
    legalMoveCount := 0
    for y, row := range board {
        legalMoveBoard[y] = make([]bool, len(board[y]), len(board[y])) 

        for x, _ := range row {
            legalMoveBoard[y][x], _ = isLegalMove(y, x, board, bluesTurn)
            if legalMoveBoard[y][x] {
                legalMoveCount++
            }
        }
    }
    return legalMoveBoard, legalMoveCount
}

func isLegalMove(y int, x int, board [][]string, bluesTurn bool) (bool, string) {
    outOfBounds := y < 0 || y >= len(board) || x <0 || x >= len(board[y])
    if outOfBounds {
        return false, "Out of bounds"
    }
    cellIsOccupied := board[y][x] != " "
    if cellIsOccupied || outOfBounds{
        return false, "Cell is occupied"
    }

    h := horizontallyLegal(y, x, board, bluesTurn)
    v := verticallyLegal(y, x, board, bluesTurn)
    d := diagonallyLegal(y, x, board, bluesTurn)

    return h || v || d, "no streak"
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

func changeBoard(y int, x int, board [][]string, bluesTurn bool) {
    otherSymbol := "O"
    ownSymbol := "X"
    if !bluesTurn {
        otherSymbol = "X"
        ownSymbol = "O"
    }
    if horizontallyLegal(y, x, board, bluesTurn) {
        dir := 1
        if hor_legal(y, x, board, bluesTurn, -1) {
            dir = -1
        } 
        xCur := x + dir
        stillCool := board[y][xCur] == otherSymbol
        for ; stillCool; xCur += dir {
            stillCool = board[y][xCur] == otherSymbol
            board[y][xCur] = ownSymbol 
        }
    } 
    if verticallyLegal(y, x, board, bluesTurn) {
        dir := 1
        if ver_legal(y, x, board, bluesTurn, -1) {
            dir = -1
        } 
        yCur := y + dir
        stillCool := board[yCur][x] == otherSymbol
        for ; stillCool; yCur += dir {
            stillCool = board[yCur][x] == otherSymbol
            board[yCur][x] = ownSymbol 
        }
    } 
    if diagonallyLegal(y, x, board, bluesTurn) {
        y_dir := 1
        x_dir := 1
        if diag_legal(y, x, board, bluesTurn, -1, -1) {
            y_dir = -1
            x_dir = -1
        } else if diag_legal(y, x, board, bluesTurn, 1, -1) {
            x_dir = -1
        } else if diag_legal(y, x, board, bluesTurn, -1, 1) {
            y_dir = -1
        }

        yCur := y + y_dir
        xCur := x + x_dir
        stillCool := board[yCur][xCur] == otherSymbol
        for ; stillCool; yCur, xCur = yCur + y_dir, xCur + x_dir {
            stillCool = board[yCur][xCur] == otherSymbol
            board[yCur][xCur] = ownSymbol 
        }
    }

    board[y][x] = ownSymbol
}

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
        legalMove, reason := isLegalMove(y,x,board, bluesTurn)

        if legalMove {
            changeBoard(y,x,board, bluesTurn)    
            bluesTurn = !bluesTurn
        } else {
            fmt.Println("This (", y,", ", x, ") is not a legal move:", reason)
        }
	}

    winner := determineWinner(board)
    fmt.Println("winner is: ", winner)

    fmt.Println("Quitting reversi")
}
