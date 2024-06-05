package main

import (
	"fmt"
	//"github.com/eiannone/keyboard"
	"strings"
	"time"
	//"math/rand/v2"

	
)

func main() {
	tetrominos := tetrominos() 
	db := board(22, 12, tetrominos)
	fmt.Println(printBoard((db)))
	
	for gameLoop() {
		random_block := 0
		if findActives(db) == nil {
			dropTetromino(*tetrominos[random_block], db)
		}
		if tetrominoGravity(db, *tetrominos[random_block]) {
			dropTetromino(*tetrominos[random_block], db)
		} 
		fmt.Print("\033[H\033[2J")
		fmt.Println(printBoard((db)))
		// char, _, err := keyboard.GetSingleKey()
		// if (err != nil) {
		// 		panic(err)
		// }
		// fmt.Printf("You pressed: %q\r\n", char)
		time.Sleep(1 * time.Second)
	}
}
		
func gameLoop() bool {
	return true
}

func tetrominos() []*tetromino {
	tetrominos := []*tetromino {
		newTetromino("I", "🟦", [][]int{{0, -2}, {0, -1}, {0, 0}, {0, 1}}, "⬜"),
		newTetromino("T", "🟪", [][]int{{0, 1}, {0, -1}, {0, 0}, {1, 0}}, "⬜"),
		newTetromino("Z", "🟥", [][]int{{-1, -1}, {-1, 0}, {0, 0}, {0, 1}}, "⬜"),
		newTetromino("S", "🟩", [][]int{{0, -1}, {-1, 0}, {0, 0}, {1, -1}}, "⬜"),
		newTetromino("O", "🟨", [][]int{{0, -1}, {-1, -1}, {0, 0}, {-1, 0}}, "⬜"),
		newTetromino("J", "🟫", [][]int{{1, 1}, {0, -1}, {0, 0}, {0, 1}}, "⬜"),
		newTetromino("L", "🟧", [][]int{{-1, -1}, {0, -1}, {0, 0}, {0, 1}}, "⬜"),
		newTetromino("Background", "⬛", nil,"⬜"),
	}
	return tetrominos
}

func dropTetromino(piece tetromino, db [][]*cell) {
	start_row := 1
	start_col := 6
	for i := 0; i < len(piece.coords); i++ {
		db[start_row+piece.coords[i][0]][start_col+piece.coords[i][1]].block = piece.block
		db[start_row+piece.coords[i][0]][start_col+piece.coords[i][1]].active = true
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

func tetrominoGravity(db[][]*cell, piece tetromino) bool {
	actives := findActives(db)
	grav := 1
	floor := 22
	for block:=0; block<len(actives); block++ {
		next_cell := db[actives[block][0] + grav][actives[block][1]]
		if next_cell.location[0] == floor || next_cell.occupied {
			resetCells(db, true, piece, actives)
			return true
		} else {
			resetCells(db, false, piece, actives)
			next_cell.block = piece.block
			next_cell.active = true
		}
	}
	return false
}

func resetCells(db[][]*cell, hit bool, piece tetromino, actives[][]int) {
	for block:=0; block<len(actives); block++ {
		cell := db[actives[block][0]][actives[block][1]]
		if !hit {
			cell.block = piece.reset
		} else {
			cell.occupied = true
		}
		cell.active = false
	}
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
