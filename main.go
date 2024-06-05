package main

import (
	"fmt"
	"strings"
	"time"
	"math/rand/v2"
	
)

func main() {
	for {
		time.Sleep(1 * time.Second)
		db := board(22, 12)

		tetrominos := []*tetromino {
			newTetromino("I", "🟦", [][]int{{0, -2}, {0, -1}, {0, 0}, {0, 1}}),
			newTetromino("T", "🟪", [][]int{{-1, 0}, {0, -1}, {0, 0}, {0, 1}}),
			newTetromino("Z", "🟥", [][]int{{-1, -1}, {-1, 0}, {0, 0}, {0, 1}}),
			newTetromino("S", "🟩", [][]int{{0, -1}, {-1, 0}, {0, 0}, {1, -1}}),
			newTetromino("O", "🟨", [][]int{{0, -1}, {-1, -1}, {0, 0}, {-1, 0}}),
			newTetromino("J", "🟫", [][]int{{1, 1}, {0, -1}, {0, 0}, {0, 1}}),
			newTetromino("L", "🟧", [][]int{{-1, -1}, {0, -1}, {0, 0}, {0, 1}}),
		}

		picker := rand.IntN(7)
		dropTetromino(*tetrominos[picker], db)
		fmt.Println(printBoard((db)))
	}
}

func dropTetromino(piece tetromino, db [][]*cell) {
	for i := 0; i < len(piece.coords); i++ {
		db[2+piece.coords[i][0]][5+piece.coords[i][1]].block = piece.block
		db[2+piece.coords[i][0]][5+piece.coords[i][1]].active = true
	}
}


type tetromino struct {
	name   string
	block  string
	coords [][]int
}

func newTetromino(name string, block string, coords [][]int) *tetromino {
	t := tetromino{name: name}
	t.block = block
	t.coords = coords
	return &t
}

type cell struct {
	location int
	occupied bool
	active bool
	block    string
}

func newCell(location int, occupied bool, active bool, block string) *cell {
	c := cell{location: location}
	c.occupied = occupied
	c.active = active
	c.block = block
	return &c
}

func board(height int, width int) [][]*cell {
	var grid [][]*cell
	var row []*cell
	block := "⬛"
	space := "⬜"

	for f := 0; f <= height; f++ {
		// top and bottom border
		if f == 0 || f == height {
			for i := 0; i < width; i++ {
				row = append(row, newCell(f, true, false, block))
			}
			grid = append(grid, row)
			row = nil
		} else {
			// middle cells
			row = append(row, newCell(f, true, false, block))
			for j := 0; j < width-2; j++ {
				row = append(row, newCell(f, false, false, space))
			}
			row = append(row, newCell(f, true, false, block))
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
