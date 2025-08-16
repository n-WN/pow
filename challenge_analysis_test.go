package pow

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"testing"
	"time"
)

// analyzeChallengeString decodes and analyzes a challenge string
func analyzeChallengeString(challengeStr string) (*Challenge, error) {
	c, err := DecodeChallenge(challengeStr)
	if err != nil {
		return nil, err
	}
	
	fmt.Printf("Challenge Analysis:\n")
	fmt.Printf("  Difficulty: %d\n", c.d)
	fmt.Printf("  Value: %s\n", c.x.String())
	fmt.Printf("  Value (hex): %s\n", fmt.Sprintf("%x", c.x.Bytes()))
	
	// Check if it's an edge case
	if c.x.Sign() == 0 {
		fmt.Printf("  Type: EDGE CASE - Zero value\n")
	} else if c.x.Cmp(one) == 0 {
		fmt.Printf("  Type: EDGE CASE - One value\n")
	} else {
		fmt.Printf("  Type: Regular value\n")
	}
	
	return c, nil
}

// TestSpecificChallengeAnalysis analyzes the specific challenge provided by the user
func TestSpecificChallengeAnalysis(t *testing.T) {
	challengeStr := "s.AAFfkA==.wxZVoJ86n1h9CNavECXG4w=="
	
	fmt.Printf("\n=== ANALYZING SPECIFIC CHALLENGE ===\n")
	c, err := analyzeChallengeString(challengeStr)
	if err != nil {
		t.Fatalf("Failed to decode challenge: %v", err)
	}
	
	// Test correctness first
	fmt.Printf("\n=== CORRECTNESS TEST ===\n")
	optimizedSolution := c.Solve()
	originalSolution := c.solveOriginal()
	
	fmt.Printf("Optimized solution: %s\n", optimizedSolution)
	fmt.Printf("Original solution:  %s\n", originalSolution)
	
	if optimizedSolution != originalSolution {
		t.Errorf("Solutions don't match!")
		t.Errorf("Optimized: %s", optimizedSolution)
		t.Errorf("Original:  %s", originalSolution)
	} else {
		fmt.Printf("Both implementations produce identical results\n")
	}
	
	// Verify solution is correct
	valid, err := c.Check(optimizedSolution)
	if err != nil {
		t.Errorf("Check failed: %v", err)
	} else if !valid {
		t.Errorf("Solution is invalid")
	} else {
		fmt.Printf("Solution is valid\n")
	}
}

// BenchmarkSpecificChallenge benchmarks the specific challenge provided
func BenchmarkSpecificChallenge(b *testing.B) {
	challengeStr := "s.AAFfkA==.wxZVoJ86n1h9CNavECXG4w=="
	c, err := DecodeChallenge(challengeStr)
	if err != nil {
		b.Fatalf("Failed to decode challenge: %v", err)
	}
	
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
}

// TestSpecificChallengePerformance measures actual performance for this specific challenge
func TestSpecificChallengePerformance(t *testing.T) {
	challengeStr := "s.AAFfkA==.wxZVoJ86n1h9CNavECXG4w=="
	c, err := DecodeChallenge(challengeStr)
	if err != nil {
		t.Fatalf("Failed to decode challenge: %v", err)
	}
	
	fmt.Printf("\n=== PERFORMANCE COMPARISON ===\n")
	
	// Measure optimized version multiple times
	var optimizedTimes []time.Duration
	for i := 0; i < 10; i++ {
		start := time.Now()
		c.Solve()
		optimizedTimes = append(optimizedTimes, time.Since(start))
	}
	
	// Measure original version multiple times (with timeout protection)
	var originalTimes []time.Duration
	for i := 0; i < 10; i++ {
		start := time.Now()
		c.solveOriginal()
		elapsed := time.Since(start)
		originalTimes = append(originalTimes, elapsed)
		
		// If taking too long, break early
		if elapsed > time.Second {
			fmt.Printf("Original implementation too slow, stopping after %d iterations\n", i+1)
			break
		}
	}
	
	// Calculate averages
	var optimizedAvg, originalAvg time.Duration
	for _, d := range optimizedTimes {
		optimizedAvg += d
	}
	optimizedAvg /= time.Duration(len(optimizedTimes))
	
	for _, d := range originalTimes {
		originalAvg += d
	}
	originalAvg /= time.Duration(len(originalTimes))
	
	fmt.Printf("Optimized average: %v (%d samples)\n", optimizedAvg, len(optimizedTimes))
	fmt.Printf("Original average:  %v (%d samples)\n", originalAvg, len(originalTimes))
	
	if len(originalTimes) > 0 {
		improvement := float64(originalAvg) / float64(optimizedAvg)
		fmt.Printf("Performance improvement: %.2fx\n", improvement)
		
		if improvement > 2 {
			fmt.Printf("Significant performance improvement detected!\n")
		} else {
			fmt.Printf("Minimal performance difference (expected for non-edge cases)\n")
		}
	}
}

// TestDifficultyLevelAnalysis analyzes what difficulty level the challenge represents
func TestDifficultyLevelAnalysis(t *testing.T) {
	// Manually decode the difficulty to understand it better
	parts := []string{"s", "AAFfkA==", "wxZVoJ86n1h9CNavECXG4w=="}
	
	dBytes, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatalf("Failed to decode difficulty: %v", err)
	}
	
	// Pad to 4 bytes if needed
	if len(dBytes) < 4 {
		dBytes = append(make([]byte, 4-len(dBytes)), dBytes...)
	}
	
	difficulty := binary.BigEndian.Uint32(dBytes)
	
	fmt.Printf("\n=== DIFFICULTY ANALYSIS ===\n")
	fmt.Printf("Difficulty bytes: %v\n", dBytes)
	fmt.Printf("Difficulty value: %d\n", difficulty)
	
	// Categorize the difficulty
	if difficulty == 0 {
		fmt.Printf("Category: No work required\n")
	} else if difficulty <= 4 {
		fmt.Printf("Category: Very low (loop unrolling optimization applies)\n")
	} else if difficulty <= 20 {
		fmt.Printf("Category: Low\n")
	} else if difficulty <= 100 {
		fmt.Printf("Category: Medium\n")
	} else if difficulty <= 1000 {
		fmt.Printf("Category: High\n")
	} else {
		fmt.Printf("Category: Very high (edge case optimizations would be crucial)\n")
	}
	
	// Estimate expected performance without optimizations
	if difficulty > 100 {
		fmt.Printf("Without optimizations, this would take considerable time\n")
	}
}