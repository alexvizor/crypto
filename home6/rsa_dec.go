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

	raw_strs := strings.Fields(strings.Replace(string(raw), "\n", " ", -1))

	ctext, _ := new(big.Int).SetString(raw_strs[0], 0)
	N, _ := new(big.Int).SetString(raw_strs[1], 0)
	e := big.NewInt(65537)

	p, q := get_factor(N)
	if p == nil || q == nil {
		log.Fatalln("Didn't get factor")
	}

	f := new(big.Int).Sub(N, p)
	f.Sub(f, q).Add(f, big.NewInt(1))

	d := new(big.Int).ModInverse(e, f)

	pkcs1 := new(big.Int).Exp(ctext, d, N).Bytes()
	for i, bt := range pkcs1 {
		if bt == 0x00 {
			fmt.Printf("%s\n", pkcs1[i+1:])
			break
		}
	}
}

func get_factor(N *big.Int) (*big.Int, *big.Int) {
	A := mathutil.SqrtBig(N)
	A.Add(A, big.NewInt(1))

	x := new(big.Int).Mul(A, A)
	x = mathutil.SqrtBig(x.Sub(x, N))

	p := new(big.Int).Sub(A, x)
	q := new(big.Int).Add(A, x)

	if new(big.Int).Mul(p, q).Cmp(N) == 0 {
		return p, q
	}

	return nil, nil
}
