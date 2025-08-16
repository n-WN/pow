package pow

import (
	"testing"
	"github.com/ncw/gmp"
)

// BenchmarkEdgeCasePerformance demonstrates the performance improvement for edge cases
func BenchmarkEdgeCasePerformance(b *testing.B) {
	// Test the optimization for values that hit the fast path
	b.Run("Zero_HighDifficulty", func(b *testing.B) {
		c := &Challenge{d: 10000, x: gmp.NewInt(0)}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			c.Solve()
		}
	})
	
	b.Run("One_HighDifficulty", func(b *testing.B) {
		c := &Challenge{d: 10000, x: gmp.NewInt(1)}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			c.Solve()
		}
	})
	
	// Compare with a regular case for the same difficulty
	b.Run("Regular_HighDifficulty", func(b *testing.B) {
		c := &Challenge{d: 10000, x: gmp.NewInt(12345)}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			c.Solve()
		}
	})
}

// TestCorrectness ensures all optimizations maintain correctness
func TestCorrectness(t *testing.T) {
	testCases := []struct {
		d uint32
		x int64
	}{
		{0, 0}, {0, 1}, {0, 12345},
		{1, 0}, {1, 1}, {1, 12345},
		{10, 0}, {10, 1}, {10, 12345},
		{100, 0}, {100, 1}, {100, 12345},
	}
	
	for _, tc := range testCases {
		c := &Challenge{d: tc.d, x: gmp.NewInt(tc.x)}
		solution := c.Solve()
		
		valid, err := c.Check(solution)
		if err != nil {
			t.Errorf("Check failed for d=%d, x=%d: %v", tc.d, tc.x, err)
			continue
		}
		if !valid {
			t.Errorf("Solution invalid for d=%d, x=%d", tc.d, tc.x)
		}
	}
}