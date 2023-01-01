// Package sudoku contains types and functions for representing Sudoku boards
package sudoku

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type Symbol int

const (
	None Symbol = 0 // Not really a symbol
	One Symbol = 1
	Two Symbol = 2
	Three Symbol = 4
	Four Symbol = 8
	Five Symbol = 16
	Six Symbol = 32
	Seven Symbol = 64
	Eight Symbol = 128
	Nine Symbol = 256
)

func FromByte(b byte) Symbol {
	if b < byte('1') {
		return None
	}
	if b > byte('9') {
		return None
	}
	syms := []Symbol{One, Two, Three, Four, Five, Six, Seven, Eight, Nine}
	return syms[b - '1']
}

func (s Symbol) Byte() byte {
	switch s {
	case One:
		return '1'
	case Two:
		return '2'
	case Three:
		return '3'
	case Four:
		return '4'
	case Five:
		return '5'
	case Six:
		return '6'
	case Seven:
		return '7'
	case Eight:
		return '8'
	case Nine:
		return '9'
	default:
		return '.'
	}
}

// Set is a set of symbols
type Set struct {
	bitmap int
}

func Empty() *Set {
	return &Set{0}
}

func All() *Set {
	return &Set{511}
}

func (s *Set) Add(sym Symbol) bool {
	before := s.bitmap
	s.bitmap = s.bitmap | int(sym)
	return before == s.bitmap
}

func (s *Set) Has(sym Symbol) bool {
	return (s.bitmap & int(sym)) == int(sym)
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
			board.Square[r][c] = FromByte(b)
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
			builder.WriteByte(sym.Byte())
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
		if symbols.Add(s) {
			return false
		}
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
		if symbols.Add(s) {
			return false
		}
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
			if symbols.Add(s) {
				return false
			}
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

// StillValid determines if a valid board is still valid after changing square (r, c)
func (b Board) StillValid(r, c int) bool {
	if !b.rowValid(r) {
		return false
	}
	if !b.columnValid(c) {
		return false
	}
	sbr, sbc := r / 3, c / 3
	if !b.subBoardValid(sbr, sbc) {
		return false
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

func (b *Board) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return fmt.Errorf("failed to unmarshal data %s: %v", data, err)
	}
	parsed, err := Parse(strings.NewReader(s))
	if err != nil {
		return fmt.Errorf("Failed to unmarshal sudoku board: %v", err)
	}
	b.Square = parsed.Square
	return nil
}

func (b Board) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.Text())
}

