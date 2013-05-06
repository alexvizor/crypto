package main

import "fmt"
import "math/big"

func main() {
	p, _ := new(big.Int).SetString("13407807929942597099574024998205846127479365820592393377723561443721764030073546976801874298166903427690031858186486050853753882811946569946433649006084171", 0)
	g, _ := new(big.Int).SetString("11717829880366207009516117596335367088558084999998952205599979459063929499736583746670572176471460312928594829675428279466566527115212748467589894601965568", 0)
	h, _ := new(big.Int).SetString("3239475104050450443565264378728065788649097520952449527834792452971981976143292558073856937958553180532878928001494706097394108577585732452307673444020333", 0)

	B := big.NewInt(2)
	B.Exp(B, big.NewInt(20), p)
	mitms := make(map[string]*big.Int)
	one := big.NewInt(1)

	fmt.Println("p:", p)
	fmt.Println("g:", g)
	fmt.Println("h:", h)

	fmt.Println("Started calculating")

	for i := big.NewInt(0); i.Cmp(B) != 0; i.Add(i, one) {
		ex := new(big.Int).Exp(g, i, p)
		ex.ModInverse(ex, p)

		res := ex.Mul(ex, h)
		res.Mod(res, p)

		mitms[res.String()] = new(big.Int).Set(i)
	}

	fmt.Println("Calculated x1s")

	var x1, x0 *big.Int
	var ok bool
	for i := big.NewInt(0); i.Cmp(B) != 0; i.Add(i, one) {
		r := new(big.Int).Exp(g, B, p)
		r.Exp(r, i, p)

		if x1, ok = mitms[r.String()]; ok {
			x0 = new(big.Int).Set(i)
			fmt.Println("Found x0", x0, "x1", x1)
			break
		}
	}

	if x1 == nil || x0 == nil {
		fmt.Println("Something went wrong - no solution")
	} else {
		x := x0.Mul(x0, B)
		x.Add(x, x1)

		fmt.Println(x)
	}
}
