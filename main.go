package main

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"math/rand/v2"
	"os"
	"slices"
	"strings"
	"time"
)

var needBlock bool = false
var runGame bool = true
var blockPos int = rand.IntN(3)

func main() {

	tetrominos := tetrominos()
	db := board(26, 12, tetrominos)
	logo()
	fmt.Println(printBoard((db)))
	fmt.Println("   Press Any Key to Play")
	startGame()
	keyPresses := keyboardChannel()
	randomBlock := rand.IntN(7)
	newGame := true

	for runGame {
		newRandomNumber := rand.IntN(7)
		if newGame || needBlock {
			randomBlock = newRandomNumber
			dropTetromino(*tetrominos[randomBlock], db, 5, 6)
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

var orientations = map[string][][][]int{
	"I": {
		{{0, 0}, {0, 1}, {0, 2}, {0, 3}},
		{{0, 1}, {1, 1}, {2, 1}, {3, 1}},
		{{0, 0}, {0, 1}, {0, 2}, {0, 3}},
		{{0, 1}, {1, 1}, {2, 1}, {3, 1}},
	},
	"O": {
		{{0, 0}, {0, 1}, {1, 0}, {1, 1}},
		{{0, 0}, {0, 1}, {1, 0}, {1, 1}},
		{{0, 0}, {0, 1}, {1, 0}, {1, 1}},
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
		{{0, 1}, {0, 2}, {1, 0}, {1, 1}},
		{{0, 0}, {1, 0}, {1, 1}, {2, 1}},
	},
	"Z": {
		{{0, 0}, {0, 1}, {1, 1}, {1, 2}},
		{{0, 1}, {1, 0}, {1, 1}, {2, 0}},
		{{0, 0}, {0, 1}, {1, 1}, {1, 2}},
		{{0, 1}, {1, 0}, {1, 1}, {2, 0}},
	},
	"J": {
		{{0, 0}, {1, 0}, {1, 1}, {1, 2}},
		{{0, 0}, {0, 1}, {1, 0}, {2, 0}},
		{{0, 0}, {0, 1}, {0, 2}, {1, 2}},
		{{0, 1}, {1, 1}, {2, 0}, {2, 1}},
	},
	"L": {
		{{0, 2}, {1, 0}, {1, 1}, {1, 2}},
		{{0, 0}, {1, 0}, {2, 0}, {2, 1}},
		{{0, 0}, {0, 1}, {0, 2}, {1, 0}},
		{{0, 0}, {0, 1}, {1, 1}, {2, 1}},
	},
}

func tetrominos() []*tetromino {
	tetrominos := []*tetromino{
		newTetromino("I", "🟦", orientations["I"], "⬜", "  "),
		newTetromino("T", "🟪", orientations["T"], "⬜", "  "),
		newTetromino("Z", "🟥", orientations["Z"], "⬜", "  "),
		newTetromino("S", "🟩", orientations["S"], "⬜", "  "),
		newTetromino("O", "🟨", orientations["O"], "⬜", "  "),
		newTetromino("J", "🟫", orientations["J"], "⬜", "  "),
		newTetromino("L", "🟧", orientations["L"], "⬜", "  "),
		newTetromino("Background", "⬛", nil, "⬜", "  "),
	}
	return tetrominos
}

func dropTetromino(piece tetromino, db [][]*cell, start_x int, start_y int) {
	var pos int
	if blockPos < 3 {
		pos = blockPos + 1
	} else {
		pos = 0
	}
	for i := 0; i < len(piece.coords[pos]); i++ {
		db[start_x+piece.coords[pos][i][0]][start_y+piece.coords[pos][i][1]].block = piece.block
		db[start_x+piece.coords[pos][i][0]][start_y+piece.coords[pos][i][1]].active = true
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
	completedLines(db, piece)
	//dest = dest[len(dest)-4:]
	if isGameOver(dest, db) {
		fmt.Println("Game Over")
		os.Exit(0)
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
			cell.block = piece.reset
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
		db[x_val][y_val].block = piece.block
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

func rotateBlock(dest [][]int, piece tetromino) [][]int {
	originalCoords := piece.coords
	var rotPointX int
	var rotPointY int
	var newDest [][]int
	avoidOBX := 0
	avoidOBY := 0
	for i := range originalCoords[0] {
		x := dest[i][0]
		y := dest[i][1]
		if i == 0 {
			rotPointX = x
			rotPointY = y
		}
		if rotPointY+originalCoords[blockPos][i][1] > 11 {
			avoidOBY = rotPointY + originalCoords[blockPos][i][1] - 11
		}
		if rotPointX+originalCoords[blockPos][i][0] > 25 {
			avoidOBX = rotPointX + originalCoords[blockPos][i][0] - 25
		}
		newDest = append(newDest, []int{rotPointX + originalCoords[blockPos][i][0] - avoidOBX, rotPointY + originalCoords[blockPos][i][1] - avoidOBY})
	}
	if blockPos < 3 {
		blockPos += 1
	} else {
		blockPos = 0
	}
	return newDest
}

func hardDrop(dest [][]int, db [][]*cell) [][]int {
	var newDest [][]int
	var rows []int
	var columns []int
	floor := 25
	for block := 0; block < len(dest); block++ {
		rows = append(rows, dest[block][0])
		columns = append(columns, dest[block][1])
	}
	rowMax := slices.Max(rows)
	colMin := slices.Min(columns)
	colMax := slices.Max(columns)
	for rows := rowMax; rows < len(db)-1; rows++ {
		for cells := colMin; cells <= colMax; cells++ {
			if db[rows][cells].occupied {
				if db[rows][cells].location[0] < floor {
					floor = db[rows][cells].location[0] - 1
				}
			}
		}
	}
	for i := range dest {
		newDest = append(newDest, []int{dest[i][0] + floor - rowMax, dest[i][1]})
	}
	return newDest
}

func completedLines(db [][]*cell, piece tetromino) {
	var full []int
	var total int
	for line := range db {
		for cell := range db[line] {
			if db[line][cell].occupied {
				total += 1
				full = append(full, cell)
			}
		}
		if total == 12 {
			for cell := 1; cell < len(full)-1; cell++ {
				db[line][cell].occupied = false
				db[line][cell].block = piece.reset
			}
			//moveBlocksDown(db, piece)
		}
		total = 0
		full = nil
	}

}

// func moveBlocksDown(db [][]*cell, piece tetromino) {
// 	for line:=0; line<len(db)-1; line++ {
// 		for cell:=1; cell<len(db[line])-1; cell++ {
// 			if db[line][cell].occupied {
// 				db[line+1][cell].block = piece.block
// 				db[line+1][cell].occupied = true
// 				db[line][cell].block = piece.reset
// 				db[line][cell].occupied = false
// 			}
// 		}
// 	}
// }

func isGameOver(dest [][]int, db [][]*cell) bool {
	for block := 0; block < len(dest); block++ {
		if dest[block][0] == 6 && db[dest[block][0]][dest[block][1]].occupied {
			return true
		}
	}
	return false
}

type tetromino struct {
	name      string
	block     string
	coords    [][][]int
	reset     string
	invisible string
}

func newTetromino(name string, block string, coords [][][]int, reset string, invisible string) *tetromino {
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
					row = append(row, newCell([]int{f, i}, false, false, tetrominos[7].block+"\n"))
				} else {
					row = append(row, newCell([]int{f, i}, false, false, tetrominos[7].block))
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
	padding := "   "
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
