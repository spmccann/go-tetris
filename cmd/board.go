package main

import (
	"fmt"
	"strings"
)

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