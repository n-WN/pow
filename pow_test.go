package pow

import (
	"fmt"
	"testing"
	"time"
)

func TestBasicFunctionality(t *testing.T) {
	// Test with a small difficulty for quick verification
	c := GenerateChallenge(5)
	solution := c.Solve()
	
	valid, err := c.Check(solution)
	if err != nil {
		t.Fatalf("Check failed: %v", err)
	}
	if !valid {
		t.Fatal("Solution should be valid")
	}
}

func TestChallengeEncodeDecode(t *testing.T) {
	original := GenerateChallenge(10)
	encoded := original.String()
	
	decoded, err := DecodeChallenge(encoded)
	if err != nil {
		t.Fatalf("Failed to decode challenge: %v", err)
	}
	
	if decoded.d != original.d {
		t.Errorf("Difficulty mismatch: got %d, want %d", decoded.d, original.d)
	}
	
	if decoded.x.Cmp(original.x) != 0 {
		t.Error("x value mismatch after encode/decode")
	}
}

func BenchmarkSolveSmall(b *testing.B) {
	c := GenerateChallenge(10)
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		c.Solve()
	}
}

func BenchmarkSolveMedium(b *testing.B) {
	c := GenerateChallenge(100)
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		c.Solve()
	}
}

func BenchmarkSolveLarge(b *testing.B) {
	c := GenerateChallenge(1000)
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		c.Solve()
	}
}

func TestSolvePerformance(t *testing.T) {
	difficulties := []uint32{10, 50, 100}
	
	for _, d := range difficulties {
		t.Run(fmt.Sprintf("difficulty_%d", d), func(t *testing.T) {
			c := GenerateChallenge(d)
			
			start := time.Now()
			solution := c.Solve()
			duration := time.Since(start)
			
			t.Logf("Difficulty %d took %v", d, duration)
			
			// Verify solution is correct
			valid, err := c.Check(solution)
			if err != nil {
				t.Fatalf("Check failed: %v", err)
			}
			if !valid {
				t.Fatal("Solution should be valid")
			}
		})
	}
}