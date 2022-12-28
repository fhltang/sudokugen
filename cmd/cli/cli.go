package main

import (
	"fmt"
	"log"

	"github.com/fhltang/sudokugen/internal/generator"
)

func main() {
	board, err := generator.Generate(50, 1)

	if err != nil {
		log.Fatalf("Failed to find board: %v", err)
	}

	fmt.Print(board.Text())
}
