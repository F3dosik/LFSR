package main

import (
	"fmt"
	"log"
	"strings"
)

func PolynomialToString(coefs []byte) string {
	var sliceResult []string
	for i := len(coefs) - 1; i >= 0; i-- {
		if coefs[i] != 0 {
			switch i {
			case 0:
				sliceResult = append(sliceResult, "1")
			case 1:
				sliceResult = append(sliceResult, "x")
			default:
				sliceResult = append(sliceResult, fmt.Sprintf("x^%d", i))
			}
		}
	}
	result := strings.Join(sliceResult, " + ")
	return result
}

func ShiftRightBytes(data []byte, shift int) []byte {
	if shift <= 0 {
		return data
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
	for i := 1; i <= L; i++ {
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
	return PolynomialToString(C), L
}

func GenerateSequenceFromPoly(initial, coefs []byte, length int) []byte {
	if len(initial) != len(coefs) {
		log.Fatalf("initial must have length %d", len(coefs))
	}

	regs := initial
	seq := make([]byte, length)
	lenReg := len(regs)
	for n := 0; n < length; n++ {
		seq[n] = regs[0]
		res := byte(0)
		for i := 0; i < lenReg; i++ {
			res ^= regs[i] & coefs[lenReg-1-i]
		}
		for i := 1; i < lenReg; i++ {
			regs[i-1] = regs[i]
		}
		regs[lenReg-1] = res
	}
	return seq
}

func main() {
	seq := GenerateSequenceFromPoly([]byte{0, 0, 0, 1}, []byte{1, 1, 0, 1}, 30)
	fmt.Println(seq)

	// coefs1 := []byte{1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 1, 0, 1, 0, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 0, 1, 1, 1, 0, 1, 0, 0, 1, 0, 1, 1, 0, 0, 0, 1, 1, 0, 1, 1, 1, 1, 0, 1, 1, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 1, 1, 1, 1, 1, 0, 0, 1, 0, 1, 0, 1, 1, 1, 0, 0, 1, 1, 0, 1}
	res, L := BerlekampMassey(seq)
	fmt.Printf("%s, L=%d\n", res, L)
}
