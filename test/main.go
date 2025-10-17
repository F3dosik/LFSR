package main

import (
	"fmt"
	"strings"
)

// ShiftRightBytes: возвращает новый слайс, представляющий B(x) * x^shift
func ShiftRightBytes(data []byte, shift int) []byte {
	if shift <= 0 {
		// возвращаем копию, чтобы не менять исходный слайс
		out := make([]byte, len(data))
		copy(out, data)
		return out
	}
	out := make([]byte, len(data)+shift)
	copy(out[shift:], data)
	return out
}

// SumPolynomials: XOR (сложение в GF(2)) двух многочленов (младшие коэффициенты в начале)
func SumPolynomials(coefs1, coefs2 []byte) []byte {
	// хотим результат длины max(len(coefs1), len(coefs2))
	if len(coefs1) < len(coefs2) {
		coefs1, coefs2 = coefs2, coefs1
	}

	result := make([]byte, len(coefs1))
	copy(result, coefs1)
	for i := 0; i < len(coefs2); i++ {
		result[i] ^= coefs2[i]
	}

	// обрежем старшие нули (чтобы полином не "длиннеть" без причин)
	// но оставим хотя бы один коэффициент (1) — обычно ведущий коэффициент ненулевой,
	// но тут возможны нулевые — выполнем trim.
	for len(result) > 1 && result[len(result)-1] == 0 {
		result = result[:len(result)-1]
	}
	return result
}

// Convolve: вычисляет d = s[n] + sum_{i=1..L} c[i] * s[n-i]  (в GF(2) это XOR/AND)
func Convolve(c, s []byte, L, n int) byte {
	var result byte = 0
	for i := 1; i <= L; i++ {
		ci := byte(0)
		if i < len(c) {
			ci = c[i]
		}
		sidx := n - i
		if sidx >= 0 && sidx < len(s) {
			result ^= ci & s[sidx]
		}
	}
	return result
}

func PolynomialToString(c []byte) string {
	parts := []string{}
	for i := len(c) - 1; i >= 0; i-- {
		if c[i] != 0 {
			if i == 0 {
				parts = append(parts, "1")
			} else if i == 1 {
				parts = append(parts, "x")
			} else {
				parts = append(parts, fmt.Sprintf("x^%d", i))
			}
		}
	}
	if len(parts) == 0 {
		return "0"
	}
	return strings.Join(parts, " + ")
}

// Berlekamp-Massey для GF(2)
func BerlekampMassey(seq []byte) (string, int) {
	C := []byte{1} // C(x) = 1
	B := []byte{1} // B(x) = 1 (копия предыдущего C)
	L := 0
	m := 1
	N := len(seq)

	for n := 0; n < N; n++ {
		// Убедимся, что C имеет длину хотя бы L+1 (многочлен степени L)
		for len(C) <= L {
			C = append(C, 0)
		}

		d := seq[n] ^ Convolve(C, seq, L, n)
		if d == 0 {
			m++
			continue
		}

		// Шаг 5 (и 6) — сначала сохраняем T := C
		T := make([]byte, len(C))
		copy(T, C)

		// C := C - d*b^{-1} x^m B(x)
		// В GF(2) d = 1, b = 1, поэтому просто C = C XOR (B << m)
		shiftedB := ShiftRightBytes(B, m)
		C = SumPolynomials(C, shiftedB)

		if 2*L <= n {
			// обновляем L, B, b, m
			L = n + 1 - L
			B = T
			m = 1
		} else {
			m++
		}
	}

	return PolynomialToString(C), L
}

// Генератор последовательности по рекурренту s_n = s_{n-1} + s_{n-7} (пример f(x)=1+x+x^7)
func GenerateSequenceFromPolyL7(initial []byte, length int) []byte {
	if len(initial) < 7 {
		panic("initial must have length 7")
	}
	out := make([]byte, length)
	copy(out, initial)
	for n := len(initial); n < length; n++ {
		out[n] = out[n-1] ^ out[n-7] // GF(2): сложение = XOR
	}
	return out
}

func main() {
	// Пример: f(x) = 1 + x + x^7, L=7
	init := []byte{1, 0, 0, 0, 0, 0, 0} // s0..s6
	seq := GenerateSequenceFromPolyL7(init, 100)
	poly, L := BerlekampMassey(seq) // проверим по первым 150 битам
	fmt.Println("Found poly:", poly)
	fmt.Println("L =", L)
}
