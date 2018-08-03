package main

import (
	"baumfalk/reversi/game"
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func getNumberFromIO(scanner *bufio.Scanner, min int, max int) int {
	scanner.Scan()
	text := scanner.Text()
	num, error := strconv.Atoi(text)
	for error != nil {
		fmt.Println(error)
		scanner.Scan()
		text = scanner.Text()
		num, error = strconv.Atoi(text)
	}
	return num
}

func getCoordinates(scanner *bufio.Scanner, boardSize int) (y int, x int) {
	fmt.Println("Type y coordinate")
	y = getNumberFromIO(scanner, 0, boardSize-1)
	fmt.Println("Type x coordinate")
	x = getNumberFromIO(scanner, 0, boardSize-1)
	return
}

func printInfo(curPlayer int) {
	fmt.Println("It's player", curPlayer, "'s turn!")
}

func main() {
	fmt.Println("Starting reversi")
	boardSize := 6
	board := game.NewReversiBoard(boardSize)
	scanner := bufio.NewScanner(os.Stdin)

	curPlayer := 1
	passes := 0
	for gameDone := false; !gameDone; {
		legalMoveBoard, legalMoveCount := game.GenerateLegalMoves(board, curPlayer)
		if legalMoveCount == 0 {
			game.HandlePass(&passes, &curPlayer, &gameDone)
			continue
		}
		passes = 0
		game.PrintBoard(board, legalMoveBoard)
		printInfo(curPlayer)
		y, x := getCoordinates(scanner, boardSize)

		game.HandleMove(&y, &x, &legalMoveBoard, &curPlayer, &board)
	}

	winner := game.DetermineWinner(board)
	fmt.Println("winner is: ", winner)

	fmt.Println("Quitting reversi")
}
