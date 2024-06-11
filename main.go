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

	tetrominos := tetrominos()
	db := board(26, 12, tetrominos)
	logo()
	fmt.Println(printBoard((db)))
	fmt.Println("  Press Any Key to Play")
	startGame()
	keyPresses := keyboardChannel()
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
			actives := findActives(db)
			tetrominoPlaced(db, actives, *tetrominos[randomBlock], newRandomNumber, keyPresses)
		}
		fmt.Print("\033[H\033[2J")
		logo()
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
		if key == keyboard.KeyEsc {
			runGame = false
		} else if key == keyboard.KeyArrowLeft {
			insertBlock(dest, piece, db, 0, -1)
		} else if key == keyboard.KeyArrowRight {
			insertBlock(dest, piece, db, 0, 1)
		} else if key == keyboard.KeyArrowUp {
			rotateBlock(dest, piece, db)
		} else if key == keyboard.KeySpace {
			hardDrop(dest, piece, db)
		} else {
			insertBlock(dest, piece, db, 0, 0)
		}
	default:
		insertBlock(dest, piece, db, 0, 0)
		// No key pressed, continue with game loop
	}
}

func startGame() {
	char, _, err := keyboard.GetSingleKey()
	if err != nil {
		panic(err)
	}
	fmt.Printf("You pressed: %q\r\n", char)
}


var orientations = map[string][][][]int{
	"I": {
			{{0, 0}, {1, 0}, {2, 0}, {3, 0}},
			{{1, 0}, {1, 1}, {1, 2}, {1, 3}},
	},
	"O": {
			{{0, 0}, {0, 1}, {1, 0}, {1, 1}},
	},
	"T": {
			{{0, 1}, {1, 0}, {1, 1}, {1, 2}},
			{{0, 1}, {1, 1}, {1, 2}, {2, 1}},
			{{0, 0}, {0, 1}, {0, 2}, {1, 1}},
			{{0, 1}, {1, 0}, {1, 1}, {2, 1}},
	},
	"S": {
			{{0, 1}, {0, 2}, {1, 0}, {1, 1}},
			{{0, 0}, {1, 0}, {1, 1}, {2, 1}},
	},
	"Z": {
			{{0, 0}, {0, 1}, {1, 1}, {1, 2}},
			{{0, 1}, {1, 0}, {1, 1}, {2, 0}},
	},
	"J": {
			{{0, 0}, {1, 0}, {1, 1}, {1, 2}},
			{{0, 1}, {1, 1}, {2, 0}, {2, 1}},
			{{0, 0}, {0, 1}, {0, 2}, {1, 2}},
			{{0, 0}, {0, 1}, {1, 0}, {2, 0}},
	},
	"L": {
			{{0, 2}, {1, 0}, {1, 1}, {1, 2}},
			{{0, 0}, {0, 1}, {1, 1}, {2, 1}},
			{{0, 0}, {0, 1}, {0, 2}, {1, 0}},
			{{0, 0}, {1, 0}, {2, 0}, {2, 1}},
	},
}

func tetrominos() []*tetromino {
	tetrominos := []*tetromino{
		newTetromino("I", "🟦", orientations["I"][0], "⬜", "  "),
		newTetromino("T", "🟪", orientations["T"][0], "⬜", "  "),
		newTetromino("Z", "🟥", orientations["Z"][0], "⬜", "  "),
		newTetromino("S", "🟩", orientations["S"][0], "⬜", "  "),
		newTetromino("O", "🟨", orientations["O"][0], "⬜", "  "),
		newTetromino("J", "🟫", orientations["J"][0], "⬜", "  "),
		newTetromino("L", "🟧", orientations["L"][0], "⬜", "  "),
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

func tetrominoPlaced(db [][]*cell, actives [][]int, piece tetromino, _ int, keyPresses chan keyboard.Key) {
	dest := nextLocations(actives)
	//dest = dest[len(dest)-4:]
	if isGameOver(dest, db) {
		fmt.Println("Game Over")
		runGame = false
	}
	if isFloor(dest) || isOccupancy(db, dest) {
		setOccupied(db)
		setInactive(db, piece, actives)
		needBlock = true
	} else {
		setInactive(db, piece, actives)
		readKeyboard(keyPresses, db, piece, dest)
	}
}

func isOccupancy(db [][]*cell, dest [][]int) bool {
	for block := 0; block < len(dest); block++ {
		if db[dest[block][0]][dest[block][1]].occupied {
			return true
		}
	}
	return false
}

func isFloor(dest [][]int) bool {
	floor := 26
	var rows []int
	for block := 0; block < len(dest); block++ {
		rows = append(rows, dest[block][0])
	}
	return slices.Max(rows) == floor
}

func setOccupied(db [][]*cell) {
	actives := findActives(db)
	for block := 0; block < len(actives); block++ {
		db[actives[block][0]][actives[block][1]].occupied = true
	}
}

func setInactive(db [][]*cell, piece tetromino, actives [][]int) {
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

func nextLocations(actives [][]int) [][]int {
	var dest [][]int
	for block := 0; block < len(actives); block++ {
		move := []int{actives[block][0] + 1, actives[block][1]}
		dest = append(dest, move)
	}
	return dest
}

func insertBlock(dest [][]int, piece tetromino, db [][]*cell, xValMod int, yValMod int) {
	obCounter := boundaries(dest, yValMod)
	for loc := 0; loc < len(dest); loc++ {
		x_val := dest[loc][0] + xValMod
		y_val := dest[loc][1] + yValMod + obCounter
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

func boundaries(dest [][]int, yValMod int) int {
	obCounter := 0 //resets block coordinate values if user input would push them off the left or right side of board
	for loc := 0; loc < len(dest); loc++ {
		y_val := dest[loc][1] + yValMod
		if y_val < 1 {
			obCounter = 1
		}
		if y_val > 10 {
			obCounter = -1
		}
	}
	return obCounter
}

// func rotateBlock(dest [][]int, piece tetromino, db [][]*cell) {
// 	cords := piece.coords
// 	var new_cords [][]int 
// 	for cord := 0; cord<len(cords); cord++ {
// 		new_x := cords[cord][1] * -1
// 		new_y := cords[cord][0]
// 		new_cords = append(new_cords, []int{new_x, new_y})
// 	}
// 	for loc := 0; loc < len(dest); loc++ {
		
// 		x_val := dest[loc][0]
// 		y_val := dest[loc][1]
// 		db[x_val][y_val].block = piece.reset
// 		db[x_val][y_val].active = false
// 		db[x_val + new_cords[loc][0]][y_val + new_cords[loc][1]].block = piece.block
// 		db[x_val + new_cords[loc][0]][y_val + new_cords[loc][1]].active = true
// 	}
// }


func rotateBlock(dest [][]int, piece tetromino, db [][]*cell) {
  originalCoords := piece.coords
  rotatedCoords := make([][]int, len(originalCoords))
  for i := range originalCoords {
    rotatedCoords[i] = []int{-originalCoords[i][1], originalCoords[i][0]} // Apply rotation formula
  }

  for i := range originalCoords {
    x := dest[i][0]
    y := dest[i][1]
    db[x][y].block = piece.reset
    db[x][y].active = false
    db[x + rotatedCoords[i][0]][y + rotatedCoords[i][1]].block = piece.block
    db[x + rotatedCoords[i][0]][y + rotatedCoords[i][1]].active = true
  
  }
}

func hardDrop(dest [][]int, piece tetromino, db [][]*cell) {
}

func isGameOver(dest [][]int, db [][]*cell) bool {
	for block := 0; block < len(dest); block++ {
		if dest[block][0] == 4 && db[dest[block][0]][dest[block][1]].occupied {
			return true
		}
	}
	return false
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
	padding := "  "
	for j := 4; j < len(board); j++ {
		sb.WriteString(padding)
		for i := 0; i < len(board[j]); i++ {
			sb.WriteString(board[j][i].block)
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func logo() {
	fmt.Print(`
  _______   _        _     
 |__   __| | |      (_)    
    | | ___| |_ _ __ _ ___ 
    | |/ _ \ __| '__| / __|
    | |  __/ |_| |  | \__ \
    |_|\___|\__|_|  |_|___/
                           
                           
`)
}
