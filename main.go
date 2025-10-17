package main

import (
	"fmt"
	"log"
	"strings"
)

func PolynomialToString(coefs []byte) string {
	monoms := []string{"1"}
	for i := 0; i < len(coefs); i++ {
		if coefs[i] != 0 {
			if i == 0 {
				monoms = append(monoms, "x")
				continue
			}
			monoms = append(monoms, fmt.Sprintf("x^%d", i+1))

		}
	}
	result := strings.Join(monoms, " ⊕ ")
	return result
}

func ShiftRightBytes(data []byte, shift int) []byte {
	if shift <= 0 {
		out := make([]byte, len(data))
		copy(out, data)
		return out
	}
	result := make([]byte, len(data)+shift)
	copy(result[shift:], data)
	return result
}

func SumPolynomials(coefs1, coefs2 []byte) []byte {
	if len(coefs1) < len(coefs2) {
		coefs1, coefs2 = coefs2, coefs1
	}

	result := make([]byte, len(coefs1))
	copy(result, coefs1)

	for i := 0; i < len(coefs2); i++ {
		result[i] ^= coefs2[i]
	}
	return result
}

func Convolve(c, s []byte, L, n int) byte {
	result := byte(0)
	for i := 1; i <= L && n-i >= 0; i++ {
		result ^= c[i] & s[n-i]
	}
	return result
}

func BerlekampMassey(seq []byte) (string, int) {
	C := []byte{1}
	B := []byte{1}
	L := 0
	m := 1
	N := len(seq)

	for n := 0; n < N; n++ {
		delta := seq[n] ^ Convolve(C, seq, L, n)
		if delta == 0 {
			m++
		} else if 2*L <= n {
			T := make([]byte, len(C))
			copy(T, C)
			C = SumPolynomials(C, ShiftRightBytes(B, m))
			L = n + 1 - L
			B = T
			m = 1
		} else {
			C = SumPolynomials(C, ShiftRightBytes(B, m))
			m++
		}
	}
	return PolynomialToString(C[1:]), L
}

func GenerateSequenceFromPoly(initial, coefs []byte, length int) []byte {
	if len(initial) != len(coefs) {
		log.Fatalf("initial must have length %d", len(coefs))
	}

	regs := append([]byte(nil), initial...)
	fmt.Printf("Initialization: %v\n", regs)
	seq := make([]byte, length)
	lenReg := len(regs)

	for n := 0; n < length; n++ {
		seq[n] = regs[0]
		feedback := byte(0)
		for i := 0; i < lenReg; i++ {
			feedback ^= regs[i] & coefs[i]
		}
		fmt.Printf("step %2d: %v → fb=%d\n", n, regs, feedback)
		for i := 1; i < lenReg; i++ {
			regs[i-1] = regs[i]
		}
		regs[lenReg-1] = feedback
	}
	fmt.Printf("step %2d: %v\n", length, regs)

	return seq
}

func main() {
	// seq := GenerateSequenceFromPoly([]byte{0, 0, 1}, []byte{1, 0, 1}, 30)
	// fmt.Println(seq)
	// seq2 := []byte{0, 1, 0, 1, 1, 1, 1, 0, 0, 0, 1, 0, 0, 1, 1, 0, 1, 0, 1, 1, 1, 1}
	// res, L := BerlekampMassey(seq2)
	// fmt.Printf("%s, L=%d\n", res, L)

}

func PolyMod(a, b uint32) uint32 {
	for deg(a) >= deg(b) {
		a ^= b << (deg(a) - deg(b))
	}
	return a
}

func deg(x uint32) int {
	d := -1
	for x > 0 {
		x >>= 1
		d++
	}
	return d
}

func GCD(a, b uint32) uint32 {
	for b != 0 {
		a, b = b, PolyMod(a, b)
	}
	return a
}

func IsIrreducible(f uint32) bool {
	n := deg(f)
	x := uint32(2) // x

	for k := 1; k <= n/2; k++ {
		// x = x^(2) mod f(x)
		x = PolyMod(MulMod(x, x, f), f)

		// gcd(x - 2, f)
		diff := x ^ 2
		if GCD(f, diff) != 1 {
			return false
		}
	}
	return true
}

// умножение многочленов mod 2 и mod f
func MulMod(a, b, f uint32) uint32 {
	res := uint32(0)
	for b > 0 {
		if b&1 != 0 {
			res ^= a
		}
		b >>= 1
		a <<= 1
		if deg(a) >= deg(f) {
			a ^= f << (deg(a) - deg(f))
		}
	}
	return res
}
