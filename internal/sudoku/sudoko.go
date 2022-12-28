// Package sudoku contains types and functions for representing Sudoku boards
package sudoku

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Symbol int

const (
	None Symbol = iota // Not really a symbol
	One
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
)

// Set is a set of symbols
type Set struct {
	elements map[Symbol]int
}

func Empty() *Set {
	return &Set{make(map[Symbol]int)}
}

func All() *Set {
	s := Empty();
	for i := 1; i <= 9; i++ {
		s.Add(Symbol(i))
	}
	return s
}

func (s *Set) Add(sym Symbol) {
	s.elements[sym] = 0
}

func (s *Set) Has(sym Symbol) bool {
	_, found := s.elements[sym]
	return found
}

// Board is a representation of a 9x9 sudoku board
type Board struct {
	Square [9][9]Symbol
}

func New() *Board {
	b := &Board{}

	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			b.Square[r][c] = None
		}
	}

	return b
}

// Parse from a reader. Input is expected to be 9 lines separated by '\n'
// with 9 characters per line. Each character is 1-9 for the 9 symbols or
// "." for empty/None.
func Parse(r io.Reader) (*Board, error) {
	board := New()
	scanner := bufio.NewScanner(r)
	for r := 0; r < 9; r++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("unexpected failure after parsing %d lines: %v", r, scanner.Err())
		}
		line := scanner.Bytes()
		if len(line) != 9 {
			return nil, fmt.Errorf("expecting 9 symbols in line %d", r)
		}
		for c := 0; c < 9; c++ {
			b := line[c]
			if b == byte('.') {
				board.Square[r][c] = None
				continue
			}
			if b < byte('1') {
				return nil, fmt.Errorf("Unexpected symbol %s", []byte{b})
			}
			if b > byte('9') {
				return nil, fmt.Errorf("Unexpected symbol %s", []byte{b})
			}
			board.Square[r][c] = Symbol(b - '0')
		}
	}
	return board, nil
}

func (b Board) Text() string {
	var builder strings.Builder
	for r := 0; r < 9; r++ {
		for c := 0 ; c <9; c++ {
			sym := b.Square[r][c]
			if sym == None {
				builder.WriteByte('.')
				continue
			}
			builder.WriteByte(byte('0' + sym))
		}
		builder.WriteByte('\n')
	}
	return builder.String()
}

func (b Board) rowValid(r int) bool {
	symbols := Empty()
	for c := 0; c < 9; c++ {
		s := b.Square[r][c]
		if s == None {
			continue
		}
		if symbols.Has(s) {
			return false
		}
		symbols.Add(s)
	}
	return true
}

func (b Board) columnValid(c int) bool {
	symbols := Empty()
	for r := 0; r < 9; r++ {
		s := b.Square[r][c]
		if s == None {
			continue
		}
		if symbols.Has(s) {
			return false
		}
		symbols.Add(s)
	}
	return true
}

func (b Board) subBoardValid(sbr, sbc int) bool {
	symbols := Empty()
	for rd := 0; rd < 3; rd++ {
		r := (3 * sbr) + rd
		for cd := 0; cd < 3; cd++ {
			c := (3 * sbc) + cd
			s := b.Square[r][c]
			if s == None {
				continue
			}
			if symbols.Has(s) {
				return false
			}
			symbols.Add(s)
		}
	}
	return true
}

// Valid returns true if the partial board is valid
func (b Board) Valid() bool {
	for r := 0; r < 9; r++ {
		if !b.rowValid(r) {
			return false
		}
	}
	for c := 0; c < 9; c++ {
		if !b.columnValid(c) {
			return false
		}
	}
	for sbr := 0; sbr < 3; sbr++ {
		for sbc := 0; sbc < 3; sbc++ {
			if !b.subBoardValid(sbr, sbc) {
				return false
			}
		}
	}
	return true
}

func (b Board) Clone() *Board {
	board := b
	return &board
}

func (b Board) Full() bool {
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if b.Square[r][c] == None {
				return false
			}
		}
	}
	return true
}
