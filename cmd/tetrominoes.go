package main

type tetromino struct {
	name      string
	block     string
	coords    [][][]int
	reset     string
}

func newTetromino(name string, block string, coords [][][]int, reset string) *tetromino {
	t := tetromino{name: name}
	t.block = block
	t.coords = coords
	t.reset = reset
	return &t
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
		newTetromino("I", "🟦", orientations["I"], "⬛"),
		newTetromino("T", "🟪", orientations["T"], "⬛"),
		newTetromino("Z", "🟥", orientations["Z"], "⬛"),
		newTetromino("S", "🟩", orientations["S"], "⬛"),
		newTetromino("O", "🟨", orientations["O"], "⬛"),
		newTetromino("J", "🟫", orientations["J"], "⬛"),
		newTetromino("L", "🟧", orientations["L"], "⬛"),
		newTetromino("Background", "🔳", nil, "⬛"),
	}
	return tetrominos
}