package tool

import (
	"github.com/google/uuid"
	"testing"
)

func BenchmarkUUID(b *testing.B)  {
	for i:=0; i < b.N; i ++ {
		b.Log(uuid.Must(uuid.NewRandom()).String())
	}
}
