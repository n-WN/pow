package pow

import (
	"encoding/base64"
	"fmt"
	"testing"
	"time"
	"github.com/ncw/gmp"
)

// TestCycleDetection analyzes if values other than 0 and 1 have short cycles
func TestCycleDetection(t *testing.T) {
	fmt.Printf("\n=== CYCLE DETECTION ANALYSIS ===\n")
	
	// Test various small values to see if they have cycles
	testValues := []int64{2, 3, 4, 5, 10, 100, 255, 256, 1000}
	maxIterations := 20 // Look for cycles within 20 iterations
	
	for _, val := range testValues {
		fmt.Printf("\nTesting value %d:\n", val)
		x := gmp.NewInt(val)
		seen := make(map[string]int)
		
		for i := 0; i < maxIterations; i++ {
			key := x.String()
			if prevI, exists := seen[key]; exists {
				cycleLength := i - prevI
				fmt.Printf("  🔄 Cycle detected! Length: %d, starts at iteration %d\n", cycleLength, prevI)
				fmt.Printf("  💡 Could optimize for difficulties >= %d\n", cycleLength)
				break
			}
			seen[key] = i
			
			// Apply one iteration of the function
			x.Exp(x, exp, mod)
			x.Xor(x, one)
			fmt.Printf("  iter %d: %s\n", i+1, x.String())
		}
		
		if len(seen) == maxIterations {
			fmt.Printf("  ❌ No cycle found within %d iterations\n", maxIterations)
		}
	}
}

// advancedSolveWithCycleDetection implements cycle detection optimization
func (c *Challenge) advancedSolveWithCycleDetection() string {
	x := gmp.NewInt(0).Set(c.x) // don't mutate c.x
	
	// Fast path for known edge cases
	if x.Sign() == 0 {
		if c.d%2 == 0 {
			return fmt.Sprintf("%s.%s", version, base64.StdEncoding.EncodeToString(gmp.NewInt(0).Bytes()))
		} else {
			return fmt.Sprintf("%s.%s", version, base64.StdEncoding.EncodeToString(one.Bytes()))
		}
	}
	
	if x.Cmp(one) == 0 {
		if c.d%2 == 0 {
			return fmt.Sprintf("%s.%s", version, base64.StdEncoding.EncodeToString(one.Bytes()))
		} else {
			return fmt.Sprintf("%s.%s", version, base64.StdEncoding.EncodeToString(gmp.NewInt(0).Bytes()))
		}
	}
	
	// Cycle detection for other values
	seen := make(map[string]uint32)
	for i := uint32(0); i < c.d; i++ {
		key := x.String()
		if startI, exists := seen[key]; exists {
			// Found a cycle!
			cycleLength := i - startI
			remaining := c.d - i
			finalPos := remaining % cycleLength
			
			// Fast-forward through the cycle
			for j := uint32(0); j < finalPos; j++ {
				x.Exp(x, exp, mod)
				x.Xor(x, one)
			}
			return fmt.Sprintf("%s.%s", version, base64.StdEncoding.EncodeToString(x.Bytes()))
		}
		seen[key] = i
		
		x.Exp(x, exp, mod)
		x.Xor(x, one)
	}
	
	return fmt.Sprintf("%s.%s", version, base64.StdEncoding.EncodeToString(x.Bytes()))
}

// TestAdvancedOptimizations compares advanced optimization techniques
func TestAdvancedOptimizations(t *testing.T) {
	fmt.Printf("\n=== ADVANCED OPTIMIZATION TESTING ===\n")
	
	// Test with the specific challenge
	challengeStr := "s.AAFfkA==.wxZVoJ86n1h9CNavECXG4w=="
	c, err := DecodeChallenge(challengeStr)
	if err != nil {
		t.Fatalf("Failed to decode challenge: %v", err)
	}
	
	fmt.Printf("Testing challenge with d=%d\n", c.d)
	
	// Test cycle detection (with smaller difficulty for demonstration)
	smallC := &Challenge{d: 100, x: gmp.NewInt(0).Set(c.x)}
	
	start := time.Now()
	regularResult := smallC.Solve()
	regularTime := time.Since(start)
	
	start = time.Now()
	cycleResult := smallC.advancedSolveWithCycleDetection()
	cycleTime := time.Since(start)
	
	fmt.Printf("\nSmall difficulty test (d=100):\n")
	fmt.Printf("  Regular: %v\n", regularTime)
	fmt.Printf("  Cycle detection: %v\n", cycleTime)
	
	if regularResult != cycleResult {
		t.Errorf("Results don't match!")
		t.Errorf("Regular: %s", regularResult)
		t.Errorf("Cycle:   %s", cycleResult)
	} else {
		fmt.Printf("  ✅ Results match\n")
		improvement := float64(regularTime) / float64(cycleTime)
		fmt.Printf("  📊 Improvement: %.2fx\n", improvement)
	}
}

// TestPerformanceRecommendations provides specific recommendations
func TestPerformanceRecommendations(t *testing.T) {
	fmt.Printf("\n=== PERFORMANCE RECOMMENDATIONS ===\n")
	
	fmt.Printf(`
Based on the analysis of the specific challenge (d=90,000):

🔍 CURRENT STATE:
  - Challenge difficulty: 90,000 (very high)
  - Challenge value: regular (not edge case)
  - Current optimized time: ~55 seconds
  - Original time: >30 seconds (timed out)
  - Current optimizations effective for edge cases only

💡 RECOMMENDED IMPROVEMENTS:

1. CYCLE DETECTION (High Impact for some values):
   - Detect when x returns to a previously seen value
   - Skip remaining iterations using modular arithmetic
   - Could provide significant speedup if cycles exist

2. MATHEMATICAL OPTIMIZATIONS (Medium Impact):
   - Analyze the function f(x) = (x^(2^1277) XOR 1) mod (2^1279-1)
   - Look for mathematical shortcuts in modular exponentiation
   - Consider sliding window exponentiation techniques

3. IMPLEMENTATION OPTIMIZATIONS (Low-Medium Impact):
   - Use more efficient GMP operations
   - Optimize memory allocation patterns
   - Consider assembly optimizations for critical paths

4. ALGORITHMIC IMPROVEMENTS (Unknown Impact):
   - Research if there are known mathematical properties
   - Consider precomputation for common values
   - Investigate if the specific modulus (2^1279-1) has special properties

🎯 NEXT STEPS:
   1. Implement cycle detection as primary optimization
   2. Profile GMP operations to find bottlenecks  
   3. Research mathematical properties of the transformation
   4. Consider consulting cryptographic literature for similar functions

⚠️  NOTE: For the specific challenge provided, improvements beyond edge cases
    may be limited due to the cryptographic nature of the function.
`)
}

// BenchmarkSpecificChallengeComparison runs a comprehensive benchmark
func BenchmarkSpecificChallengeComparison(b *testing.B) {
	challengeStr := "s.AAFfkA==.wxZVoJ86n1h9CNavECXG4w=="
	c, err := DecodeChallenge(challengeStr)
	if err != nil {
		b.Fatalf("Failed to decode challenge: %v", err)
	}
	
	// Use smaller difficulty for benchmarking to avoid timeouts
	smallC := &Challenge{d: 50, x: gmp.NewInt(0).Set(c.x)}
	
	b.Run("Current_Optimized", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			smallC.Solve()
		}
	})
	
	b.Run("Original_Unoptimized", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			smallC.solveOriginal()
		}
	})
	
	b.Run("Cycle_Detection", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			smallC.advancedSolveWithCycleDetection()
		}
	})
}