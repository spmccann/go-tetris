package main

import (
	"fmt"
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

func readKeyboard(keyPresses chan keyboard.Key, piece tetromino, db [][]*cell, dest [][]int) {
	// Check for key press (non-blocking)
	select {
	case key := <-keyPresses:
		if key == keyboard.KeyEsc {
			routeKeyPress("esc", piece, db, dest)
		} else if key == keyboard.KeyArrowLeft {
			routeKeyPress("left", piece, db, dest)
		} else if key == keyboard.KeyArrowRight {
			routeKeyPress("right", piece, db, dest)
		} else if key == keyboard.KeyArrowUp {
			routeKeyPress("up", piece, db, dest)
		} else if key == keyboard.KeySpace {
			routeKeyPress("space", piece, db, dest)
		} 
	default:
		routeKeyPress("", piece, db, dest)
		// No key pressed, continue with game loop
	}
}

func startGame() {
	_, _, err := keyboard.GetSingleKey()
	if err != nil {
		panic(err)
	}
}