package main_test

import (
	"log/slog"
	"testing"

	"github.com/illbjorn/skal/pkg/clog"
)

func BenchmarkSlog(b *testing.B) {
	for range b.N {
		slog.Info("Hello", "key1", "key2", "key3", 4)
	}
}

func BenchmarkClog(b *testing.B) {
	for range b.N {
		clog.Info("Hello", "key1", "key2", "key3", 4)
	}
}
