package pow

import (
	"context"
	"fmt"
	"testing"
	"time"
	"github.com/ncw/gmp"
)

// TestSpecificChallengePerformanceSmarter measures performance with timeout protection
func TestSpecificChallengePerformanceSmarter(t *testing.T) {
	challengeStr := "s.AAFfkA==.wxZVoJ86n1h9CNavECXG4w=="
	c, err := DecodeChallenge(challengeStr)
	if err != nil {
		t.Fatalf("Failed to decode challenge: %v", err)
	}
	
	fmt.Printf("\n=== SMART PERFORMANCE COMPARISON ===\n")
	fmt.Printf("Challenge difficulty: %d\n", c.d)
	fmt.Printf("Challenge value: %s\n", c.x.String())
	
	// Measure optimized version
	fmt.Printf("\n--- Testing Optimized Implementation ---\n")
	var optimizedTime time.Duration
	optimizedSolution := ""
	{
		start := time.Now()
		optimizedSolution = c.Solve()
		optimizedTime = time.Since(start)
		fmt.Printf("Optimized completed in: %v\n", optimizedTime)
	}
	
	// Test original with timeout
	fmt.Printf("\n--- Testing Original Implementation (with timeout) ---\n")
	var originalTime time.Duration
	originalSolution := ""
	
	// Use context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	done := make(chan bool)
	go func() {
		start := time.Now()
		originalSolution = c.solveOriginal()
		originalTime = time.Since(start)
		done <- true
	}()
	
	select {
	case <-done:
		fmt.Printf("Original completed in: %v\n", originalTime)
		
		// Verify results match
		if optimizedSolution != originalSolution {
			t.Errorf("Solutions don't match!")
		} else {
			fmt.Printf("Solutions match\n")
		}
		
		// Calculate improvement
		if originalTime > 0 && optimizedTime > 0 {
			improvement := float64(originalTime) / float64(optimizedTime)
			fmt.Printf("Performance improvement: %.2fx\n", improvement)
		}
		
	case <-ctx.Done():
		fmt.Printf("Original implementation timed out after 30 seconds\n")
		fmt.Printf("Optimized vs Original comparison:\n")
		fmt.Printf("   Optimized: %v (completed)\n", optimizedTime)
		fmt.Printf("   Original:  >30s (timed out)\n")
		fmt.Printf("Performance improvement: >%.0fx\n", float64(30*time.Second)/float64(optimizedTime))
	}
}

// TestOptimizationEffectiveness tests different optimization scenarios
func TestOptimizationEffectiveness(t *testing.T) {
	fmt.Printf("\n=== OPTIMIZATION EFFECTIVENESS ANALYSIS ===\n")
	
	testCases := []struct {
		name       string
		difficulty uint32
		value      int64
		maxTime    time.Duration
	}{
		{"Low difficulty", 10, 12345, 5 * time.Second},
		{"Medium difficulty", 100, 12345, 10 * time.Second},
		{"High difficulty", 1000, 12345, 30 * time.Second},
		{"Edge case - Zero (high diff)", 1000, 0, 5 * time.Second},
		{"Edge case - One (high diff)", 1000, 1, 5 * time.Second},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := &Challenge{d: tc.difficulty, x: gmp.NewInt(tc.value)}
			
			// Test optimized
			start := time.Now()
			optimizedSolution := c.Solve()
			optimizedTime := time.Since(start)
			
			fmt.Printf("\n%s (d=%d, x=%d):\n", tc.name, tc.difficulty, tc.value)
			fmt.Printf("  Optimized: %v\n", optimizedTime)
			
			// Test original with timeout
			ctx, cancel := context.WithTimeout(context.Background(), tc.maxTime)
			defer cancel()
			
			done := make(chan bool)
			var originalTime time.Duration
			var originalSolution string
			
			go func() {
				start := time.Now()
				originalSolution = c.solveOriginal()
				originalTime = time.Since(start)
				done <- true
			}()
			
			select {
			case <-done:
				if optimizedSolution != originalSolution {
					t.Errorf("Solutions don't match for %s", tc.name)
				}
				improvement := float64(originalTime) / float64(optimizedTime)
				fmt.Printf("  Original:  %v\n", originalTime)
				fmt.Printf("  Improvement: %.2fx\n", improvement)
				
			case <-ctx.Done():
				fmt.Printf("  Original:  >%v (timed out)\n", tc.maxTime)
				improvement := float64(tc.maxTime) / float64(optimizedTime)
				fmt.Printf("  Improvement: >%.0fx\n", improvement)
			}
		})
	}
}

// BenchmarkOptimizedVsOriginal provides benchmark comparison for different scenarios
func BenchmarkOptimizedVsOriginal(b *testing.B) {
	// Benchmark scenarios that can complete reasonably quickly
	scenarios := []struct {
		name       string
		difficulty uint32
		value      int64
	}{
		{"Low_d10", 10, 12345},
		{"Medium_d50", 50, 12345},
		{"EdgeCase_Zero_d100", 100, 0},
		{"EdgeCase_One_d100", 100, 1},
		{"LoopUnroll_d3", 3, 12345},
	}
	
	for _, scenario := range scenarios {
		c := &Challenge{d: scenario.difficulty, x: gmp.NewInt(scenario.value)}
		
		b.Run(scenario.name+"_Optimized", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				c.Solve()
			}
		})
		
		// Only benchmark original for lower difficulties
		if scenario.difficulty <= 50 {
			b.Run(scenario.name+"_Original", func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					c.solveOriginal()
				}
			})
		}
	}
}