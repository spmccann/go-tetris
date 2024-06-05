package main

import (
	"fmt"
	//"github.com/eiannone/keyboard"
	"slices"
	"strings"
	"time"
	"math/rand/v2"

	
)

func main() {
	tetrominos := tetrominos() 
	db := board(22, 12, tetrominos)
	fmt.Println(printBoard((db)))	
	randomBlock := rand.IntN(7)
	newGame := true
	for {
		newRandomNumber := rand.IntN(7)
		if newGame || tetrominoPlaced(db, *tetrominos[randomBlock], newRandomNumber) {
			randomBlock = newRandomNumber
			dropTetromino(*tetrominos[randomBlock], db, 2, 6)
			newGame = false
		}  
		fmt.Print("\033[H\033[2J")
		fmt.Println(printBoard((db)))
		time.Sleep(200 * time.Millisecond)
	}
}

func tetrominos() []*tetromino {
	tetrominos := []*tetromino {
		newTetromino("I", "🟦", [][]int{{0, -2}, {0, -1}, {0, 0}, {0, 1}}, "⬜"),
		newTetromino("T", "🟪", [][]int{{1, 0}, {0, -1}, {0, 0}, {0, 1}}, "⬜"),
		newTetromino("Z", "🟥", [][]int{{-1, -1}, {-1, 0}, {0, 0}, {0, 1}}, "⬜"),
		newTetromino("S", "🟩", [][]int{{0, -1}, {-1, 0}, {0, 0}, {1, -1}}, "⬜"),
		newTetromino("O", "🟨", [][]int{{0, -1}, {-1, -1}, {0, 0}, {-1, 0}}, "⬜"),
		newTetromino("J", "🟫", [][]int{{1, 1}, {0, -1}, {0, 0}, {0, 1}}, "⬜"),
		newTetromino("L", "🟧", [][]int{{-1, -1}, {0, -1}, {0, 0}, {0, 1}}, "⬜"),
		newTetromino("Background", "⬛", nil,"⬜"),
	}
	return tetrominos
}

func dropTetromino(piece tetromino, db [][]*cell, start_x int, start_y int) {
	for i := 0; i < len(piece.coords); i++ {
		db[start_x+piece.coords[i][0]][start_y+piece.coords[i][1]].block = piece.block
		db[start_x+piece.coords[i][0]][start_y+piece.coords[i][1]].active = true
	}
}

func findActives(db[][]*cell) [][]int{
	var actives [][]int
	for row:=0; row<len(db); row++ {
		for cell:=0; cell<len(db[row]); cell++ {
			if db[row][cell].active {
				actives = append(actives, db[row][cell].location)
			}
		}
	}
	return actives
}

func tetrominoPlaced(db[][]*cell, piece tetromino, _ int) bool{
	dest := nextLocations(db)
	dest = dest[len(dest)-4:]
	if isFloor(db) || isOccupancy(db) {
		setOccupied(db)
		setInactive(db, piece)
		return true
	} else {
		setInactive(db, piece)
		insertBlock(dest, piece, db)
	}
	return false
}

func isOccupancy(db [][]*cell) bool {
	dest := nextLocations(db)
	for block:=0; block<len(dest); block++ {
		if db[dest[block][0]][dest[block][1]].occupied {
			return true
		}
	}
	return false
}

func setOccupied(db [][]*cell) {
	actives := findActives(db)
	for block:=0; block<len(actives); block++ {
		db[actives[block][0]][actives[block][1]].occupied = true
	}
}

func setInactive(db [][]*cell, piece tetromino) {
	actives := findActives(db)
	for block:=0; block<len(actives); block++ {
		db[actives[block][0]][actives[block][1]].active = false
		if !db[actives[block][0]][actives[block][1]].occupied {
			db[actives[block][0]][actives[block][1]].block = piece.reset
		}
	}
}

func nextLocations(db [][]*cell) [][]int{
	actives := findActives(db)
	var dest [][]int
	for block:=0; block<len(actives); block++ {
		move := []int{actives[block][0] + 1, actives[block][1]}
		dest = append(dest, move)
		}
	return dest
}

func insertBlock(dest [][]int, piece tetromino, db[][]*cell) {
	for loc:=0; loc<len(dest); loc++ {
		x_val := dest[loc][0]
		y_val := dest[loc][1]
		db[x_val][y_val].block = piece.block
		db[x_val][y_val].active = true
	}
}

func isFloor(db [][]*cell) bool {
	floor := 22
	dest := nextLocations(db)
	var rows []int
	for block:=0; block<len(dest); block++ {
		rows = append(rows, dest[block][0])	
	}
	return slices.Max(rows) == floor
}

type tetromino struct {
	name   string
	block  string
	coords [][]int
	reset string
}

func newTetromino(name string, block string, coords [][]int, reset string) *tetromino {
	t := tetromino{name: name}
	t.block = block
	t.coords = coords
	t.reset = reset
	return &t
}

type cell struct {
	location []int
	occupied bool
	active bool
	block string
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
		if f == 0 || f == height {
			for i := 0; i < width; i++ {
				row = append(row, newCell([]int{f,i}, true, false, tetrominos[7].block))
			}
			grid = append(grid, row)
			row = nil
		} else {
			// middle cells
			for j := 0; j < width; j++ {
				if j == 0 || j == width -1 {
					row = append(row, newCell([]int{f,j}, true, false, tetrominos[7].block))
				} else {
					row = append(row, newCell([]int{f,j}, false, false, tetrominos[7].reset))
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
