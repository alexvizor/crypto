// home1 project main.go
package main

import "encoding/hex"
import "log"
import "fmt"
import "sort"
import "flag"
import "os"
import "io"
import "bufio"
import "strings"

var path = flag.String("path", "data.txt", "Path to file with data")

type weighted_guesses [][]byte

func (wg weighted_guesses) Len() int {
	return len(wg)
}

func (wg weighted_guesses) Less(i, j int) bool {
	return wg[j][1] < wg[i][1]
}

func (wg weighted_guesses) Swap(i, j int) {
	wg[i], wg[j] = wg[j], wg[i]
}

// This program should deduce key from several ciphertexts and print last message decripted
func main() {
	ciphers := parseData(*path)

	bciphers := make([][]byte, len(ciphers))
	for i, cipher := range ciphers {
		bcipher, err := hex.DecodeString(cipher)
		if err != nil {
			log.Fatal(err)
		}
		bciphers[i] = bcipher
	}

	key := guess_key(bciphers)

	fmt.Println(string(xor(bciphers[len(bciphers)-1], key)))
}

func guess_key(ciphers [][]byte) []byte {
	max_len := 0
	for _, c := range ciphers {
		ln := len(c)
		if max_len < ln {
			max_len = ln
		}
	}

	guesses := make([][]byte, max_len)
	for i := 0; i < max_len; i++ {
		guesses[i] = make([]byte, 256)
	}

	for i := 1; i < len(ciphers); i++ {
		for k := 0; k < i; k++ {
			x := xor(ciphers[k], ciphers[i])

			for l, b := range x {
				if (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') {
					guesses[l][ciphers[i][l]^' '] += 1
					guesses[l][ciphers[k][l]^' '] += 1
				}
			}
		}
	}

	fguess := make([]weighted_guesses, max_len)
	for i, guess := range guesses {
		fguess[i] = make(weighted_guesses, 0)
		for k, w := range guess {
			if w != 0 {
				fguess[i] = append(fguess[i], []byte{byte(k), w})
			}
		}

		sort.Sort(fguess[i])
	}

	key := make([]byte, max_len)
	for pos, c := range fguess {
		if len(c) > 0 && len(c[0]) > 0 {
			key[pos] = c[0][0]
		}
	}

	return key
}

func xor(left, right []byte) []byte {
	ln := len(left)
	if len(right) < ln {
		ln = len(right)
	}

	result := make([]byte, ln)

	for i := 0; i < ln; i++ {
		result[i] = left[i] ^ right[i]
	}

	return result
}

func parseData(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	eof := false
	ciphers := make([]string, 0)
	reader := bufio.NewReader(file)
	for !eof {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			eof = true
			err = nil
		}
		if err != nil {
			log.Fatal(err)
		}
		line = strings.Trim(line, " \n")

		if len(line) > 0 {
			ciphers = append(ciphers, line)
		}
	}

	return ciphers
}
