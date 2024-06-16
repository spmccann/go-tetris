package main

import (
	"fmt"
	"os"

	"github.com/eiannone/keyboard"
)

func keyboardChannel() chan keyboard.Key {
	keyPresses := make(chan keyboard.Key)
	go func() {
		defer close(keyPresses)
		for {
			_, key, err := keyboard.GetSingleKey()
			if err != nil {
				fmt.Println("Error reading key:", err)
				return
			}
			keyPresses <- key
		}
	}()
	return keyPresses
}

func readKeyboard(keyPresses chan keyboard.Key, db [][]*cell, piece tetromino, dest [][]int) {
	// Check for key press (non-blocking)
	select {
	case key := <-keyPresses:
		if key == keyboard.KeyEsc {
			fmt.Println("quitting game...")
			os.Exit(0)
		} else if key == keyboard.KeyArrowLeft {
			insertBlock(dest, piece, db, 0, -1)
		} else if key == keyboard.KeyArrowRight {
			insertBlock(dest, piece, db, 0, 1)
		} else if key == keyboard.KeyArrowUp {
			newDest := rotateBlock(dest, piece)
			insertBlock(newDest, piece, db, 0, 0)
		} else if key == keyboard.KeySpace {
			newDest := hardDrop(dest, db)
			insertBlock(newDest, piece, db, 0, 0)
		} else {
			insertBlock(dest, piece, db, 0, 0)
		}
	default:
		insertBlock(dest, piece, db, 0, 0)
		// No key pressed, continue with game loop
	}
}

func startGame() {
	_, _, err := keyboard.GetSingleKey()
	if err != nil {
		panic(err)
	}
}