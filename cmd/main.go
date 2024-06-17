package main

import (
	"fmt"
	"math/rand/v2"
	"time"
)

var needBlock bool = false
var blockPos int = rand.IntN(3)

func main() {
	tetrominos := tetrominos()
	db := board(22, 12, tetrominos)
	logo()
	//fmt.Println(printBoard((db)))
	fmt.Println("   Press Any Key to Play")
	startGame()
	keyPresses := keyboardChannel()
	randomBlock := rand.IntN(7)
	newGame := true

	for {
		newRandomNumber := rand.IntN(7)
		if newGame || needBlock {
			randomBlock = newRandomNumber
			dropTetromino(*tetrominos[randomBlock], db, 1, 6)
			newGame = false
			needBlock = false
		} else {
			actives := currentActives(db)
			tetrominoPlaced(db, actives, *tetrominos[randomBlock], newRandomNumber, keyPresses)
		}
		fmt.Print("\033[H\033[2J") //clears the terminal
		logo()
		fmt.Println(printBoard((db)))
		time.Sleep(300 * time.Millisecond)
	}
}
