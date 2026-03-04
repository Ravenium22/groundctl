package detector

import (
	"testing"
)

func BenchmarkDetectAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DetectAll()
	}
}

func BenchmarkDetectByNames_Single(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DetectByNames([]string{"git"})
	}
}

func BenchmarkDetectByNames_Five(b *testing.B) {
	names := []string{"git", "node", "python", "go", "docker"}
	for i := 0; i < b.N; i++ {
		DetectByNames(names)
	}
}

func BenchmarkDetectByNames_All(b *testing.B) {
	names := ListKnownTools()
	for i := 0; i < b.N; i++ {
		DetectByNames(names)
	}
}
