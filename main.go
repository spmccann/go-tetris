package main

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"math/rand/v2"
	"slices"
	"strings"
	"time"
)

var needBlock bool = false
var runGame bool = true

func main() {
	keyPresses := keyboardChannel()
	tetrominos := tetrominos()
	db := board(26, 12, tetrominos)
	fmt.Println(printBoard((db)))
	randomBlock := rand.IntN(7)
	newGame := true

	for runGame {
		newRandomNumber := rand.IntN(7)
		if newGame || needBlock {
			randomBlock = newRandomNumber
			dropTetromino(*tetrominos[randomBlock], db, 2, 6)
			newGame = false
			needBlock = false
		} else {
			tetrominoPlaced(db, *tetrominos[randomBlock], newRandomNumber, keyPresses)
		}
		fmt.Print("\033[H\033[2J")
		fmt.Println(printBoard((db)))
		time.Sleep(300 * time.Millisecond)
	}
}

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
		fmt.Printf("You pressed: %q\r\n", key)
		if key == keyboard.KeyEsc {
			runGame = false
		}
		if key == keyboard.KeyArrowLeft {
			fmt.Println("Block goes left")
			insertBlock(dest, piece, db, 0, -1)
		}
		if key == keyboard.KeyArrowRight {
			fmt.Println("Block goes right")
			insertBlock(dest, piece, db, 0, 1)
		}
		if key == keyboard.KeyArrowUp {
			fmt.Println("Block rotates")
			rotateBlock(dest, piece, db)
		}
		if key == keyboard.KeySpace {
			fmt.Println("Block hard drops")
			hardDrop(dest, piece, db)
		}
	default:
		insertBlock(dest, piece, db, 0, 0)
		// No key pressed, continue with game loop
	}
}

func tetrominos() []*tetromino {
	tetrominos := []*tetromino{
		newTetromino("I", "🟦", [][]int{{0, -2}, {0, -1}, {0, 0}, {0, 1}}, "⬜", "  "),
		newTetromino("T", "🟪", [][]int{{1, 0}, {0, -1}, {0, 0}, {0, 1}}, "⬜", "  "),
		newTetromino("Z", "🟥", [][]int{{-1, -1}, {-1, 0}, {0, 0}, {0, 1}}, "⬜", "  "),
		newTetromino("S", "🟩", [][]int{{0, -1}, {-1, 0}, {0, 0}, {1, -1}}, "⬜", "  "),
		newTetromino("O", "🟨", [][]int{{0, -1}, {-1, -1}, {0, 0}, {-1, 0}}, "⬜", "  "),
		newTetromino("J", "🟫", [][]int{{1, 1}, {0, -1}, {0, 0}, {0, 1}}, "⬜", "  "),
		newTetromino("L", "🟧", [][]int{{-1, -1}, {0, -1}, {0, 0}, {0, 1}}, "⬜", "  "),
		newTetromino("Background", "⬛", nil, "⬜", "  "),
	}
	return tetrominos
}

func dropTetromino(piece tetromino, db [][]*cell, start_x int, start_y int) {
	for i := 0; i < len(piece.coords); i++ {
		db[start_x+piece.coords[i][0]][start_y+piece.coords[i][1]].block = piece.invisible
		db[start_x+piece.coords[i][0]][start_y+piece.coords[i][1]].active = true
	}
}

func findActives(db [][]*cell) [][]int {
	var actives [][]int
	for row := 0; row < len(db); row++ {
		for cell := 0; cell < len(db[row]); cell++ {
			if db[row][cell].active {
				actives = append(actives, db[row][cell].location)
			}
		}
	}
	return actives
}

func tetrominoPlaced(db [][]*cell, piece tetromino, _ int, keyPresses chan keyboard.Key) {
	dest := nextLocations(db)
	dest = dest[len(dest)-4:]
	if isGameOver(dest, db) {
		runGame = false
	}
	if isFloor(db) || isOccupancy(db) {
		setOccupied(db)
		setInactive(db, piece)
		needBlock = true
	} else {
		setInactive(db, piece)
		readKeyboard(keyPresses, db, piece, dest)
	}
}

func isOccupancy(db [][]*cell) bool {
	dest := nextLocations(db)
	for block := 0; block < len(dest); block++ {
		if db[dest[block][0]][dest[block][1]].occupied {
			return true
		}
	}
	return false
}

func setOccupied(db [][]*cell) {
	actives := findActives(db)
	for block := 0; block < len(actives); block++ {
		db[actives[block][0]][actives[block][1]].occupied = true
	}
}

func setInactive(db [][]*cell, piece tetromino) {
	actives := findActives(db)
	for block := 0; block < len(actives); block++ {
		cell := db[actives[block][0]][actives[block][1]]
		cell.active = false
		if !cell.occupied {
			if actives[block][0] > 4 {
				cell.block = piece.reset
			} else if actives[block][0] == 4 {
				cell.block = "⬛"
			} else {
				cell.block = piece.invisible
			}

		}
	}
}

func nextLocations(db [][]*cell) [][]int {
	actives := findActives(db)
	var dest [][]int
	for block := 0; block < len(actives); block++ {
		move := []int{actives[block][0] + 1, actives[block][1]}
		dest = append(dest, move)
	}
	return dest
}

func insertBlock(dest [][]int, piece tetromino, db [][]*cell, xValMod int, yValMod int) {
	var x_val int
	var y_val int
	obCounter := 0 //resets block coordinate values if user input would push them off the left or right side of board
	for loc := 0; loc < len(dest); loc++ {
		y_val = dest[loc][1] + yValMod
		if y_val < 1 {
			obCounter = 1
		}
		if y_val > 10 {
			obCounter = -1
		}
	}
	for loc := 0; loc < len(dest); loc++ {
		x_val = dest[loc][0] + xValMod
		y_val = dest[loc][1] + yValMod + obCounter
		if dest[loc][0] > 5 {
			db[x_val][y_val].block = piece.block
		} else if dest[loc][0] == 5 {
			db[x_val][y_val].block = "⬛"
		} else {
			db[x_val][y_val].block = piece.invisible
		}
		db[x_val][y_val].active = true
	}
}

func rotateBlock(dest [][]int, piece tetromino, db [][]*cell) {
	for c := 0; c < len(piece.coords); c++ {
		for loc := 0; loc < len(dest); loc++ {
			x_val := dest[loc][0] - piece.coords[c][1]
			y_val := dest[loc][1] - piece.coords[c][0]
			db[x_val][y_val].block = piece.block
			db[x_val][y_val].active = true
		}
	}
}

func hardDrop(dest [][]int, piece tetromino, db [][]*cell) {
	// need to implement drop to block
	high := 0
	for loc := 0; loc < len(dest); loc++ {
		if dest[loc][0] > high {
			high = dest[loc][0]
		}
	}
	for loc := 0; loc < len(dest); loc++ {
		x_val := dest[loc][0] + (25 - high)
		y_val := dest[loc][1]
		db[x_val][y_val].block = piece.block
		db[x_val][y_val].active = true
	}
}

func isGameOver(dest [][]int, db [][]*cell) bool {
	for block := 0; block < len(dest); block++ {
		fmt.Println(dest[block][0])
		if dest[block][0] == 4 && db[dest[block][0]][dest[block][1]].occupied {
			return true
		}
	}
	return false
}

func isFloor(db [][]*cell) bool {
	floor := 26
	dest := nextLocations(db)
	var rows []int
	for block := 0; block < len(dest); block++ {
		rows = append(rows, dest[block][0])
	}
	return slices.Max(rows) == floor
}

type tetromino struct {
	name      string
	block     string
	coords    [][]int
	reset     string
	invisible string
}

func newTetromino(name string, block string, coords [][]int, reset string, invisible string) *tetromino {
	t := tetromino{name: name}
	t.block = block
	t.coords = coords
	t.reset = reset
	t.invisible = invisible
	return &t
}

type cell struct {
	location []int
	occupied bool
	active   bool
	block    string
}

func newCell(location []int, occupied bool, active bool, block string) *cell {
	c := cell{location: location}
	c.occupied = occupied
	c.active = active
	c.block = block
	return &c
}

func board(height int, width int, tetrominos []*tetromino) [][]*cell {
	var grid [][]*cell
	var row []*cell

	for f := 0; f <= height; f++ {
		// top and bottom border
		if f == 0 || f == 1 || f == 2 || f == 3 {
			for i := 0; i < width; i++ {
				if i == width-1 {
					row = append(row, newCell([]int{f, i}, false, false, tetrominos[7].invisible+"\n"))
				} else {
					row = append(row, newCell([]int{f, i}, false, false, tetrominos[7].invisible))
				}
			}
			grid = append(grid, row)
			row = nil
		} else if f == 4 || f == height {
			for i := 0; i < width; i++ {
				row = append(row, newCell([]int{f, i}, false, false, tetrominos[7].block))
			}
			grid = append(grid, row)
			row = nil
		} else {
			// middle cells
			for j := 0; j < width; j++ {
				if j == 0 || j == width-1 {
					row = append(row, newCell([]int{f, j}, true, false, tetrominos[7].block))
				} else {
					row = append(row, newCell([]int{f, j}, false, false, tetrominos[7].reset))
				}
			}
			grid = append(grid, row)
			row = nil
		}
	}
	return grid
}

func printBoard(board [][]*cell) string {
	var sb strings.Builder
	for j := 0; j < len(board); j++ {
		for i := 0; i < len(board[j]); i++ {
			sb.WriteString(board[j][i].block)
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

/*
TODO
hard drop onto block
rotate
other rotate
completed line(s)
game speed / responsiveness
score
logo / other presentation text

*/
