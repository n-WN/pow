package pow

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"

	"github.com/ncw/gmp"
)

const version = "s"

var (
	mod = gmp.NewInt(0)
	exp = gmp.NewInt(0)
	one = gmp.NewInt(1)
	two = gmp.NewInt(2)
)

func init() {
	mod.Lsh(one, 1279)
	mod.Sub(mod, one)
	exp.Lsh(one, 1277)
}

type Challenge struct {
	d uint32
	x *gmp.Int
}

// DecodeChallenge decodes a redpwnpow challenge produced by String.
func DecodeChallenge(v string) (*Challenge, error) {
	parts := strings.SplitN(v, ".", 3)
	if len(parts) != 3 || parts[0] != version {
		return nil, errors.New("incorrect version")
	}
	dBytes, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}
	if len(dBytes) > 4 {
		return nil, errors.New("difficulty too long")
	}
	// pad start with 0s to 4 bytes
	dBytes = append(make([]byte, 4-len(dBytes)), dBytes...)
	xBytes, err := base64.StdEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, err
	}
	d := binary.BigEndian.Uint32(dBytes)
	x := gmp.NewInt(0).SetBytes(xBytes)
	return &Challenge{d: d, x: x}, nil
}

// GenerateChallenge creates a new random challenge.
func GenerateChallenge(d uint32) *Challenge {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return &Challenge{
		x: gmp.NewInt(0).SetBytes(b),
		d: d,
	}
}

// String encodes the challenge in a format that can be decoded by DecodeChallenge.
func (c *Challenge) String() string {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, c.d)
	return fmt.Sprintf("%s.%s.%s", version, base64.StdEncoding.EncodeToString(b), base64.StdEncoding.EncodeToString(c.x.Bytes()))
}

// Solve solves the challenge and returns a solution proof that can be checked by Check.
func (c *Challenge) Solve() string {
	x := gmp.NewInt(0).Set(c.x) // dont mutate c.x
	
	// Fast path for edge cases (though rare in practice)
	if x.Sign() == 0 {
		// 0 -> 1 -> 0 -> 1 ... alternating pattern
		if c.d%2 == 0 {
			// Even number of iterations: 0 -> 1 -> 0 -> ... -> 0
			return fmt.Sprintf("%s.%s", version, base64.StdEncoding.EncodeToString(gmp.NewInt(0).Bytes()))
		} else {
			// Odd number of iterations: 0 -> 1 -> 0 -> ... -> 1
			return fmt.Sprintf("%s.%s", version, base64.StdEncoding.EncodeToString(one.Bytes()))
		}
	}
	
	if x.Cmp(one) == 0 {
		// 1 -> 0 -> 1 -> 0 ... alternating pattern
		if c.d%2 == 0 {
			// Even number of iterations: 1 -> 0 -> 1 -> ... -> 1
			return fmt.Sprintf("%s.%s", version, base64.StdEncoding.EncodeToString(one.Bytes()))
		} else {
			// Odd number of iterations: 1 -> 0 -> 1 -> ... -> 0
			return fmt.Sprintf("%s.%s", version, base64.StdEncoding.EncodeToString(gmp.NewInt(0).Bytes()))
		}
	}
	
	// Optimization: Unroll loop for small difficulties to reduce loop overhead
	if c.d <= 4 {
		switch c.d {
		case 1:
			x.Exp(x, exp, mod)
			x.Xor(x, one)
		case 2:
			x.Exp(x, exp, mod)
			x.Xor(x, one)
			x.Exp(x, exp, mod)
			x.Xor(x, one)
		case 3:
			x.Exp(x, exp, mod)
			x.Xor(x, one)
			x.Exp(x, exp, mod)
			x.Xor(x, one)
			x.Exp(x, exp, mod)
			x.Xor(x, one)
		case 4:
			x.Exp(x, exp, mod)
			x.Xor(x, one)
			x.Exp(x, exp, mod)
			x.Xor(x, one)
			x.Exp(x, exp, mod)
			x.Xor(x, one)
			x.Exp(x, exp, mod)
			x.Xor(x, one)
		}
	} else {
		// General case: perform the computation
		for i := uint32(0); i < c.d; i++ {
			x.Exp(x, exp, mod)
			x.Xor(x, one)
		}
	}
	
	return fmt.Sprintf("%s.%s", version, base64.StdEncoding.EncodeToString(x.Bytes()))
}

func decodeSolution(s string) (*gmp.Int, error) {
	parts := strings.SplitN(s, ".", 2)
	if len(parts) != 2 || parts[0] != version {
		return nil, errors.New("incorrect version")
	}
	yBytes, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}
	return gmp.NewInt(0).SetBytes(yBytes), nil
}

// Check verifies that a solution proof from Solve is correct.
func (c *Challenge) Check(s string) (bool, error) {
	y, err := decodeSolution(s)
	if err != nil {
		return false, fmt.Errorf("decode solution: %w", err)
	}
	
	// Fast path for edge cases
	if c.d == 0 {
		return y.Cmp(c.x) == 0, nil
	}
	
	// Apply the inverse transformation d times
	for i := uint32(0); i < c.d; i++ {
		y.Xor(y, one)
		y.Exp(y, two, mod)
	}
	
	x := gmp.NewInt(0).Set(c.x) // dont mutate c.x
	if x.Cmp(y) == 0 {
		return true, nil
	}
	x.Sub(mod, c.x)
	return x.Cmp(y) == 0, nil
}
