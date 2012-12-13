package main

import "log"
import "fmt"
import "math/big"
import "github.com/cznic/mathutil"

func main() {
	ctext, _ := new(big.Int).SetString("22096451867410381776306561134883418017410069787892831071731839143676135600120538004282329650473509424343946219751512256465839967942889460764542040581564748988013734864120452325229320176487916666402997509188729971690526083222067771600019329260870009579993724077458967773697817571267229951148662959627934791540", 0)
	N, _ := new(big.Int).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805505620686985379449212982959585501387537164015710139858647833778606925583497541085196591615128057575940752635007475935288710823649949940771895617054361149474865046711015101563940680527540071584560878577663743040086340742855278549092581", 0)
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
