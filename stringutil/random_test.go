package stringutil

import (
	"fmt"
	"testing"
)

func TestByte2Int(t *testing.T) {
	b := 0xff
	u := uint(b)
	i := int(b)
	n := int(uint(b))
	fmt.Println(u, i, n)
}

func TestRandom(t *testing.T) {
	for i := 0; i < 16; i++ {
		fmt.Println(RandomFull(128))
	}
}

func BenchmarkByte2Int(b *testing.B) {
	for j := 0; j < b.N; j++ {
		b := 0xff
		u := uint(b)
		i := int(b)
		n := int(uint(b))

		u, i, n = uint(n), int(u), i
	}
}

func BenchmarkRandom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomFull(128)
	}
}
