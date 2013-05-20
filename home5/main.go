package main

import "log"
import "fmt"
import "math/big"
import "os"
import "io/ioutil"
import "strings"

func main() {
	raw, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalln(err)
	}

	raw_strings := strings.Fields(strings.Replace(string(raw), "\n", " ", -1))

	p, _ := new(big.Int).SetString(raw_strings[0], 0)
	g, _ := new(big.Int).SetString(raw_strings[1], 0)
	h, _ := new(big.Int).SetString(raw_strings[2], 0)

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
