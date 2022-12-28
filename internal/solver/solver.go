// Package solver has a Sudoku solver
package solver

import (
	"math/rand"
	"time"

	"github.com/fhltang/sudokugen/internal/sudoku"
)

// SymbolIt iterates over a set of symbols in random order
type SymbolIt struct {
	r *rand.Rand
}

func (it *SymbolIt) Symbols(syms *sudoku.Set) chan sudoku.Symbol {
	symbols := []sudoku.Symbol{}
	for i := 1; i <= 9; i++ {
		s := sudoku.Symbol(i)
		if syms.Has(s) {
			symbols = append(symbols, s)
		}
	}
	it.r.Shuffle(len(symbols), func(i, j int) {
		symbols[i], symbols[j] = symbols[j], symbols[i]
	})

	ch := make(chan sudoku.Symbol)

	go func() {
		for _, s := range symbols {
			ch <- s
		}
		close(ch)
	}()

	return ch
}

// Coord is a linked list node
type Coord struct {
	Row, Col int
	Next *Coord
}

// PathGen generates a path of squares through a sudoku board
type PathGen struct {
	Shuffle bool
}

func (gen *PathGen) Path(b *sudoku.Board) *Coord {
	coords := []*Coord{}
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			coords = append(coords, &Coord{r, c, nil})
		}
	}

	if gen.Shuffle {
		rand.Shuffle(len(coords), func(i, j int) {
			coords[i], coords[j] = coords[j], coords[i]
		})
	}

	for i, coord := range coords {
		if i + 1 < len(coords) {
			coord.Next = coords[i+1]
		}
	}
	return coords[0]
}

type Solver struct {
	iterator *SymbolIt
	pathGenerator *PathGen

	// Solutions could be incomplete because of solver limits
	solutions []sudoku.Board
}

func New() *Solver {
	seed := time.Now().UnixMicro()
	solver := &Solver{
		iterator: &SymbolIt{r: rand.New(rand.NewSource(seed))},
		pathGenerator: &PathGen{Shuffle: false},
		solutions: []sudoku.Board{},
	}

	return solver
}

func (s *Solver) Solve(board *sudoku.Board, limit int) {
	coord := s.pathGenerator.Path(board)

	s.internalSolve(board, limit, coord)
}

func (s *Solver) internalSolve(b *sudoku.Board, limit int, coord *Coord) {
	if len(s.solutions) >= limit {
		return
	}

	if coord == nil {
		if b.Full() {
			s.solutions = append(s.solutions, *b.Clone())
		}
		return
	}

	if b.Square[coord.Row][coord.Col] != sudoku.None {
		s.internalSolve(b, limit, coord.Next)
		return
	}

	symbols := sudoku.All()
	for sym := range s.iterator.Symbols(symbols) {
		b.Square[coord.Row][coord.Col] = sym
		if !b.Valid() {
			continue
		}

		s.internalSolve(b, limit, coord.Next)
	}
	b.Square[coord.Row][coord.Col] = sudoku.None
}

func (s *Solver) Solutions() []sudoku.Board {
	return s.solutions
}

func (s *Solver) Reset() {
	s.solutions = []sudoku.Board{}
}
