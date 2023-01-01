package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/fhltang/sudokugen/internal/generator"
)

func main() {
	blanks := 50

	if len(os.Args) == 2 {
		arg, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatalf("Cannot parse blanks count %s: %v", os.Args[1], err)
		}
		blanks = arg
	}

	board, err := generator.Generate(blanks, 1)

	if err != nil {
		log.Fatalf("Failed to find board: %v", err)
	}

	fmt.Print(board.Text())
}
