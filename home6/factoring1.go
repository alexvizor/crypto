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

	x := new(big.Int).Mul(A, A)
	x = mathutil.SqrtBig(x.Sub(x, N))

	p := new(big.Int).Sub(A, x)
	q := new(big.Int).Add(A, x)

	if new(big.Int).Mul(p, q).Cmp(N) == 0 {
		fmt.Println("Result:", p)
	}
}
