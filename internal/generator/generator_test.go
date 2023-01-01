package generator

import (
	"testing"
)

func benchmarkGenerate(b *testing.B, blanks int) {
	for i := 0; i < b.N; i++ {
		Generate(blanks, 1)
	}
}

func BenchmarkGenerate10(b *testing.B) {
	benchmarkGenerate(b, 10)
}

func BenchmarkGenerate40(b *testing.B) {
	benchmarkGenerate(b, 40)
}

func BenchmarkGenerate50(b *testing.B) {
	benchmarkGenerate(b, 50)
}
