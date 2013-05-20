package main

import "log"
import "fmt"
import "os"
import "io/ioutil"
import "strings"
import "math/big"
import "github.com/cznic/mathutil"

func main() {
	raw, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalln(err)
	}

	N, _ := new(big.Int).SetString(strings.Trim(string(raw), " \n"), 0)

	A := mathutil.SqrtBig(N)
	A.Add(A, big.NewInt(1))

	res := -1
	one := big.NewInt(1)
	for A.Cmp(N) == -1 && res != 0 {
		res = check_factor(A, N)
		A.Add(A, one)
	}
}

func check_factor(A, N *big.Int) int {
	x := new(big.Int).Mul(A, A)
	x = mathutil.SqrtBig(x.Sub(x, N))

	p := new(big.Int).Sub(A, x)
	q := new(big.Int).Add(A, x)

	res := new(big.Int).Mul(p, q).Cmp(N)
	if res == 0 {
		fmt.Println("Result:", p)
	}

	return res
}
