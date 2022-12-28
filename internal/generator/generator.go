// Package generator generates sudoku puzzles
package generator

import (
	"fmt"

	"github.com/fhltang/sudokugen/internal/solver"
	"github.com/fhltang/sudokugen/internal/sudoku"
)

// Generate generates a sudoku puzzle with the given number of blank squares
// and maximum number of solutions.
func Generate(blanks int, maxSolutions int) (*sudoku.Board, error) {
	// Step 1: use solver to generate a sudoku board with no blanks.
	sol := solver.New()
	sol.Solve(sudoku.New(), 1)

	solutions := sol.Solutions()

	if len(solutions) < 1 {
		return nil, fmt.Errorf("Failed to generate a complete sudoku board")
	}

	// Step 2: iteratively delete symbols from randomly selected squares
	// ensuring that the resulting board does not have too many solutions
	board := &solutions[0]
	count := 0
	pathGenerator := &solver.PathGen{Shuffle: true}
	for coord := pathGenerator.Path(board); count < blanks && coord != nil; coord = coord.Next {
		sym := board.Square[coord.Row][coord.Col]

		if sym == sudoku.None {
			continue
		}

		board.Square[coord.Row][coord.Col] = sudoku.None
		sol.Reset()
		sol.Solve(board, maxSolutions + 1)
		if len(sol.Solutions()) <= maxSolutions {
			count++
			continue
		}
		board.Square[coord.Row][coord.Col] = sym
	}

	if count < blanks {
		return nil, fmt.Errorf("Failed to delete %d symbols, only deleted %d", blanks, count)
	}

	return board, nil
}
