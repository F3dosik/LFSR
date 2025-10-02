package main

import (
	"fmt"
	"math/big"
	"sort"
)

func main() {
	var Reg = []byte{1, 0, 1}
	var Res []byte
	var Func = []byte{1, 1, 0}
	for range 10 {
		var FuncRes byte = 0
		Res = append(Res, Reg[0])
		for i, val := range Reg {
			FuncRes ^= val & Func[i]
		}
		for i := 1; i < len(Reg); i++ {
			Reg[i-1] = Reg[i]
		}
		Reg[len(Reg)-1] = byte(FuncRes)
		fmt.Println(Reg)
	}
	fmt.Println(Res)
}


func GLOrderOptimized(n int, q int64) *big.Int {
    result := big.NewInt(1)
    qBig := big.NewInt(q)
    
    // Предварительно вычисляем qⁿ
    qn := new(big.Int).Exp(qBig, big.NewInt(int64(n)), nil)
    
    // qi будет накапливать степени q
    qi := big.NewInt(1)
    
    for i := 0; i < n; i++ {
        // Вычисляем (qⁿ - qⁱ)
        factor := new(big.Int).Sub(qn, qi)
        
        // Умножаем результат
        result.Mul(result, factor)
        
        // Умножаем qi на q для следующей итерации
        qi.Mul(qi, qBig)
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