package main

import (
	"fmt"
	"os"
	"slices"
	"github.com/eiannone/keyboard"
)
func routeKeyPress(keyInput string, piece tetromino, db [][]*cell, dest [][]int) {
	if keyInput == "esc" {
		fmt.Println("quitting game...")
		os.Exit(0)
	} else if keyInput == "left" {
		insertBlock(dest, piece, db, 0, -1)
	} else if keyInput == "right" {
		insertBlock(dest, piece, db, 0, 1)
	} else if keyInput == "up" {
		newDest := rotateBlock(dest, piece)
    insertBlock(newDest, piece, db, 0, 0)
	} else if keyInput == "space" {
		newDest := hardDrop(dest, db)
		insertBlock(newDest, piece, db, 0, 0)
	} else {
		insertBlock(dest, piece, db, 0, 0)
	}
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

func currentActives(db [][]*cell) [][]int {
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

func nextLocations(actives [][]int) [][]int {
	var dest [][]int
	for block := 0; block < len(actives); block++ {
		move := []int{actives[block][0] + 1, actives[block][1]}
		dest = append(dest, move)
	}
	return dest
}

func tetrominoPlaced(db [][]*cell, actives [][]int, piece tetromino, _ int, keyPresses chan keyboard.Key) {
	dest := nextLocations(actives)
	checkGameOver(dest, db)
	if isFloor(dest) || isOccupancy(db, dest) {
		setOccupied(db)
		setInactive(db, piece, actives)
		needBlock = true
	} else {
		completedLines(db)
		setInactive(db, piece, actives)
		readKeyboard(keyPresses, piece, db, dest)
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
	floor := 22
	var rows []int
	for block := 0; block < len(dest); block++ {
		rows = append(rows, dest[block][0])
	}
	return slices.Max(rows) == floor
}

func setOccupied(db [][]*cell) {
	actives := currentActives(db)
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

func insertBlock(dest [][]int, piece tetromino, db [][]*cell, xValMod int, yValMod int) {
	sideHit := blockBoundaries(dest, yValMod, db)
	for loc := 0; loc < len(dest); loc++ {
		x_val := dest[loc][0] + xValMod
		y_val := dest[loc][1] + yValMod + sideHit
		db[x_val][y_val].block = piece.block
		db[x_val][y_val].active = true
	}
}

func blockBoundaries(dest [][]int, yValMod int, db [][]*cell) int {
	sideHit := 0
	for loc := 0; loc < len(dest); loc++ {
		x_val := dest[loc][0]
		y_val := dest[loc][1] + yValMod
		if db[x_val][y_val].occupied && yValMod == -1 {
			sideHit = 1
		}
		if db[x_val][y_val].occupied && yValMod == 1 {
			sideHit = -1
		}
	}
	return sideHit
}

func rotateBlock(dest [][]int, piece tetromino) [][]int {
	originalCoords := piece.coords
	var rotPointX int
	var rotPointY int
	var newDest [][]int
	for i := range originalCoords[blockPos] {
		x := dest[i][0]
		y := dest[i][1]
		if i == 0 {
			rotPointX = x
			rotPointY = y
		}
		newDest = append(newDest, []int{rotPointX + originalCoords[blockPos][i][0] -1, rotPointY + originalCoords[blockPos][i][1]})	
	}
	newDest = rotateInBounds(newDest)
	if blockPos < 3 {
		blockPos += 1
	} else {
		blockPos = 0
	}
	return newDest
}

func rotateInBounds(newDest [][]int) [][]int{
	var amountRight int
	var amountBottom int
	for i := range(newDest) {
		if newDest[i][0] - 21 > amountBottom {
			amountBottom = newDest[i][0] - 21
		}
		if newDest[i][1] - 10 > amountRight {
			amountRight = newDest[i][1] - 10
		}
	}
	for j := range(newDest) {
		newDest[j][0] -= amountBottom
		newDest[j][1] -= amountRight
	}
	return newDest
}

func hardDrop(dest [][]int, db [][]*cell) [][]int {
	var newDest [][]int
	var rows []int
	var columns []int
	floor := 21
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
				if db[rows][cells].location[0] <= floor {
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

func completedLines(db [][]*cell) {
	for line := len(db) - 2; line > 0; line-- { 
		if isLineComplete(db[line]) {
			moveBlocksDown(db, line)
			line++ 
		}
	}
}

func isLineComplete(line []*cell) bool {
	for _, c := range line {
		if !c.occupied {
			return false
		}
	}
	return true
}

func moveBlocksDown(db [][]*cell, mark int) {
	for line := mark; line > 2; line-- { 
		for cell := 1; cell < len(db[line])-1; cell++ {
			db[line][cell].occupied = db[line-1][cell].occupied
			db[line][cell].block = db[line-1][cell].block
		}
	}

	for cell := 1; cell < len(db[1])-1; cell++ {
		db[1][cell].occupied = false
		db[1][cell].block = "⬛"
	}
}

func checkGameOver(dest [][]int, db [][]*cell) {
	for block := 0; block < len(dest); block++ {
		if dest[block][0] == 2 && db[dest[block][0]][dest[block][1]].occupied {
			fmt.Println("Game Over")
			os.Exit(0)
		}
	}
}
