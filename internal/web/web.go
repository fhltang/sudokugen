// Package web provides utilities for the web ui
package web

import (
	"encoding/json"
	"encoding/base64"
	"fmt"
	"log"
	
	"github.com/fhltang/sudokugen/internal/sudoku"
)

type Query struct {
	Board sudoku.Board
}

func Decode(encoded string) (*Query, error) {
	marshalled, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode base64 string %s: %v", encoded, err)
	}
	s := Query{}
	err = json.Unmarshal(marshalled, &s)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal json %s: %v", marshalled, err)
	}
	return &s, nil
}

func (s *Query) Encode() (string, error) {
	marshalled, err := json.Marshal(s)
	if err != nil {
		return "", fmt.Errorf("Failed to marshal state %v: %v", s, err)
	}

	log.Printf("Marshalled: %s", marshalled)

	return base64.StdEncoding.EncodeToString(marshalled), nil
}
