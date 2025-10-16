package main

import (
	"fmt"
	"sort"
	"strings"
)

func GLOrderOptimized(n int, q int64) int64 {
	result := int64(1)

	// Предварительно вычисляем qⁿ
	qn := int64(1)
	for i := 0; i < n; i++ {
		qn *= q
		// Проверка на переполнение
		if qn < 0 {
			panic("Overflow: q^n exceeds int64")
		}
	}

	// qi будет накапливать степени q
	qi := int64(1)

	for i := 0; i < n; i++ {
		// Вычисляем (qⁿ - qⁱ)
		factor := qn - qi

		// Умножаем результат
		result *= factor

		// Проверка на переполнение
		if result < 0 {
			panic("Overflow: result exceeds int64")
		}

		// Умножаем qi на q для следующей итерации
		qi *= q

		// Проверка на переполнение
		if qi < 0 {
			panic("Overflow: q^i exceeds int64")
		}
	}

	return result
}

func getDivisors(n int) []int {
	if n == 0 {
		return []int{}
	}

	divisors := []int{}
	i := 1
	for i*i <= n {
		if n%i == 0 {
			divisors = append(divisors, i)
			if i != n/i {
				divisors = append(divisors, n/i)
			}
		}
		i++
	}

	sort.Ints(divisors)
	return divisors
}

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

func Convolve(c, s []byte, L, N int) byte {
	result := byte(0)
	for i := 1; i <= L; i++ {
		result ^= c[i] & s[N-i]
	}
	return result
}

func BerlekampMassey(seq []byte) (string, int) {
	C := []byte{1}
	L := 0
	m := -1
	B := []byte{1}
	n := len(seq)
	delta := byte(0)
	for N := 0; N < n; N++ {
		delta = seq[N] ^ Convolve(C, seq, L, N)
		if delta == 1 {
			T := make([]byte, len(C))
			copy(T, C)
			C = SumPolynomials(C, ShiftRightBytes(B, N-m))
			if L <= N/2 {
				L = N + 1 - L
				m = N
				B = T
			}
		}
	}
	return PolynomialToString(C), L
}

func main() {
	//var Reg = []byte{1, 0, 1}
	//var Res []byte
	//var Func = []byte{1, 1, 0}
	//for range 10 {
	//	var FuncRes byte = 0
	//	Res = append(Res, Reg[0])
	//	for i, val := range Reg {
	//		FuncRes ^= val & Func[i]
	//	}
	//	for i := 1; i < len(Reg); i++ {
	//		Reg[i-1] = Reg[i]
	//	}
	//	Reg[len(Reg)-1] = byte(FuncRes)
	//	fmt.Println(Reg)
	//}
	//fmt.Println(Res)

	coefs1 := []byte{1, 1, 1, 1, 1, 0, 0, 1, 0, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 0}
	res, L := BerlekampMassey(coefs1)
	fmt.Printf("%s, L=%d", res, L)
}
