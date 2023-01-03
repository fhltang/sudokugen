package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/fhltang/sudokugen/internal/generator"
	"github.com/fhltang/sudokugen/internal/sudoku"
	"github.com/fhltang/sudokugen/internal/web"
)

func generateBoard(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	blanks := 0
	if blanksString, ok := r.PostForm["blanks"]; ok && len(blanksString) == 1 {
		v, err := strconv.Atoi(blanksString[0])
		if err != nil {
			log.Printf("Failed to parse field value: %v", blanksString[0])
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if v < 0 || v > 81 {
			log.Printf("Request blanks %d out of range", v)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		blanks = v
	}
	board, err := generator.Generate(blanks, 1)
	if err != nil {
		log.Printf("Board generation failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	q := web.Query{*board}
	encoded, err := q.Encode()
	if err != nil {
		log.Fatalf("failed to encode state: %v", err)
	}
	http.Redirect(w, r, fmt.Sprintf("?q=%s", encoded), http.StatusSeeOther)
}

func cellStyle(r, c int) string {
	styles := [][]string{
		[]string{"tl", "t", "tr"},
		[]string{"l", "c", "r"},
		[]string{"bl", "b", "br"},
	}
	return styles[r % 3][c % 3]
}

func renderBoard(query string, w http.ResponseWriter, r *http.Request) {
	q, err := web.Decode(query)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	table := &strings.Builder{}
	fmt.Fprintln(table, "<table class='board'>")
	for r := 0; r < 9; r++ {
		fmt.Fprintln(table, "<tr>")
		for c := 0; c < 9; c++ {
			fmt.Fprintf(table, "<td class='%s'>", cellStyle(r, c))
			sym := q.Board.Square[r][c]
			if sym != sudoku.None {
				table.WriteByte(sym.Byte())
			}
			table.WriteString("</td>")
		}
		fmt.Fprintln(table, "</tr>")
	}
	fmt.Fprintln(table, "</table>")

	qrcode := fmt.Sprintf("/qrcode?q=%s", query)
	
	fmt.Fprintf(w, `<html>
<head>
<title>Sudoku Generator</title>
<style type="text/css">
table.board {
    border-collapse: collapse;
}
.board td {
    width: 2em;
    height: 2em;
    border-style: solid;
    padding: 0px;
    text-align: center;
}
.board td.tl {
    border-width: medium thin thin medium;
}
.board td.t {
    border-width: medium thin thin thin;
}
.board td.tr {
    border-width: medium medium thin thin;
}
.board td.r {
    border-width: thin medium thin thin;
}
.board td.br {
    border-width: thin medium medium thin;
}
.board td.b {
    border-width: thin thin medium thin;
}
.board td.bl {
    border-width: thin thin medium medium;
}
.board td.l {
    border-width: thin thin thin medium;
}
.board td.c {
    border-width: thin thin thin thin;
}
</style>
</head>
<body>
<a href="/">New board</a>
%s

<img src="%s"/>
</body>
</html>`, table.String(), qrcode)
}

func renderForm(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<html>
<head>
<title>Sudoku Generator</title>
</head>
<body>
<h1>Sudoku Generator</h1>
<form method="POST">
<label for="blanks">Blanks</label>
<input type="text" name="blanks" value="45"/>
<input type="submit" value="Generate">
</form>

<p>
Source code: <a href="https://github.com/fhltang/sudokugen">github.com/fhltang/sudokugen</a>
</p>
</body>
</html>`)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		generateBoard(w, r)
		return
	}

	params := r.URL.Query()
	if query, ok := params["q"]; ok && len(query) > 0 {
		renderBoard(query[0], w, r)
		return
	}

	renderForm(w, r)
}

func main() {
	http.HandleFunc("/qrcode", web.QrcodeHandler)
	http.HandleFunc("/", handler)
	log.Println("Starting web server...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
