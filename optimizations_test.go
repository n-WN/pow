package pow

import (
	"encoding/base64"
	"fmt"
	"testing"
	"time"
	"github.com/ncw/gmp"
)

// solveOriginal implements the original unoptimized version of Solve for performance comparison
func (c *Challenge) solveOriginal() string {
	x := gmp.NewInt(0).Set(c.x) // dont mutate c.x
	for i := uint32(0); i < c.d; i++ {
		x.Exp(x, exp, mod)
		x.Xor(x, one)
	}
	return fmt.Sprintf("%s.%s", version, base64.StdEncoding.EncodeToString(x.Bytes()))
}

// BenchmarkPerformanceComparison compares optimized vs unoptimized performance
func BenchmarkPerformanceComparison(b *testing.B) {
	// Edge case with zero - should show massive improvement
	b.Run("EdgeCase_Zero_d1000", func(b *testing.B) {
		c := &Challenge{d: 1000, x: gmp.NewInt(0)}
		
		b.Run("Optimized", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				c.Solve()
			}
		})
		
		b.Run("Original", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				c.solveOriginal()
			}
		})
	})
	
	// Edge case with one - should show massive improvement
	b.Run("EdgeCase_One_d1000", func(b *testing.B) {
		c := &Challenge{d: 1000, x: gmp.NewInt(1)}
		
		b.Run("Optimized", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				c.Solve()
			}
		})
		
		b.Run("Original", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				c.solveOriginal()
			}
		})
	})
	
	// Small difficulty with loop unrolling - should show minor improvement
	b.Run("SmallDifficulty_d3", func(b *testing.B) {
		c := &Challenge{d: 3, x: gmp.NewInt(12345)}
		
		b.Run("Optimized", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				c.Solve()
			}
		})
		
		b.Run("Original", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				c.solveOriginal()
			}
		})
	})
	
	// Regular case - should show no significant difference
	b.Run("Regular_d10", func(b *testing.B) {
		c := &Challenge{d: 10, x: gmp.NewInt(12345)}
		
		b.Run("Optimized", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				c.Solve()
			}
		})
		
		b.Run("Original", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				c.solveOriginal()
			}
		})
	})
}

// TestPerformanceImprovements validates the performance claims from the PR description
func TestPerformanceImprovements(t *testing.T) {
	t.Run("EdgeCasePerformance", func(t *testing.T) {
		difficulties := []uint32{100, 1000, 5000}
		
		for _, d := range difficulties {
			t.Run(fmt.Sprintf("Zero_d%d", d), func(t *testing.T) {
				c := &Challenge{d: d, x: gmp.NewInt(0)}
				
				// Measure optimized version
				start := time.Now()
				optimizedResult := c.Solve()
				optimizedTime := time.Since(start)
				
				// Measure original version (but only for smaller difficulties to avoid timeouts)
				var originalTime time.Duration
				if d <= 1000 {
					start = time.Now()
					originalResult := c.solveOriginal()
					originalTime = time.Since(start)
					
					// Verify both produce same result
					if optimizedResult != originalResult {
						t.Errorf("Results don't match for d=%d", d)
					}
					
					improvement := float64(originalTime) / float64(optimizedTime)
					t.Logf("Difficulty %d: Original=%v, Optimized=%v, Improvement=%.0fx", 
						d, originalTime, optimizedTime, improvement)
					
					// Should show significant improvement
					if improvement < 100 && d >= 1000 {
						t.Logf("Warning: Expected >100x improvement for d=%d, got %.0fx", d, improvement)
					}
				} else {
					t.Logf("Difficulty %d: Optimized=%v (original too slow to measure)", d, optimizedTime)
				}
			})
		}
	})
	
	t.Run("LoopUnrollingPerformance", func(t *testing.T) {
		for d := uint32(1); d <= 4; d++ {
			t.Run(fmt.Sprintf("d%d", d), func(t *testing.T) {
				c := &Challenge{d: d, x: gmp.NewInt(12345)}
				
				// Measure optimized version
				start := time.Now()
				optimizedResult := c.Solve()
				optimizedTime := time.Since(start)
				
				// Measure original version
				start = time.Now()
				originalResult := c.solveOriginal()
				originalTime := time.Since(start)
				
				// Verify both produce same result
				if optimizedResult != originalResult {
					t.Errorf("Results don't match for d=%d", d)
				}
				
				improvement := float64(originalTime) / float64(optimizedTime)
				t.Logf("Difficulty %d: Original=%v, Optimized=%v, Improvement=%.2fx", 
					d, originalTime, optimizedTime, improvement)
			})
		}
	})
}

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
		
		// Test optimized version
		optimizedSolution := c.Solve()
		valid, err := c.Check(optimizedSolution)
		if err != nil {
			t.Errorf("Optimized Check failed for d=%d, x=%d: %v", tc.d, tc.x, err)
			continue
		}
		if !valid {
			t.Errorf("Optimized solution invalid for d=%d, x=%d", tc.d, tc.x)
		}
		
		// Test original version produces same result
		originalSolution := c.solveOriginal()
		if optimizedSolution != originalSolution {
			t.Errorf("Solutions don't match for d=%d, x=%d: optimized=%s, original=%s", 
				tc.d, tc.x, optimizedSolution, originalSolution)
		}
	}
}