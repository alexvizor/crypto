package main

import "log"
import "fmt"
import "net/http"
import "encoding/hex"

var base_url = "http://crypto-class.appspot.com/po?er="

func main() {
	secret, err := hex.DecodeString("f20bdba6ff29eed7b046d1df9fb7000058b1ffb4210a580f748b4ac714c001bd4a61044426fb515dad3f21f18aa577c0bdf302936266926ff37dbf7035d5eeb4")
	if err != nil {
		log.Fatal("Error while decoding secret")
	}

	sln := len(secret)
	result := make([]byte, sln-16)
	log.Println("Start guessing")

	chns, schn := make([]chan []byte, sln/16), make(chan int)
	for bt := 1; bt < sln/16; bt++ {
		chns[bt-1] = make(chan []byte)
		go oracle(schn, chns[bt-1], secret[:16*(bt+1)])
	}

	for i := 0; i < 3; i++ {
		select {
		case bt := <-schn:
			bt -= 1
			result = append(result[:bt*16], append(<-chns[bt], result[(bt+1)*16:]...)...)
		}
	}

	fmt.Printf("Result %s\n(%x)", result, result)
}

func oracle(sch chan int, ch chan []byte, cipher []byte) {
	cln, prev_guess := len(cipher), 0
	result, fsecret := make([]byte, 16), make([]byte, cln)
	block := cln/16 - 1

	copy(fsecret, cipher) // Create forged ciphertext
	for pad := 1; pad <= 16; pad++ {
		padbyte := make([]byte, pad)
		for i := 0; i < pad; i++ {
			padbyte[i] = byte(pad)
		}

		for guess := 0x00; guess <= 0xff; guess++ {
			fsecret = append(
				fsecret[:cln-(16+pad)],
				append(
					xor(cipher[cln-(16+pad):cln-16],
						xor(padbyte,
							append(
								[]byte{byte(guess)},
								result[16-pad+1:]...,
							),
						),
					),
					fsecret[cln-16:]...,
				)...,
			)

			resp, err := http.Get(base_url + fmt.Sprintf("%x", fsecret))
			if err != nil {
				log.Fatalln("Block", block, "Error when requesting remote host:", err)
			}

			if resp.StatusCode == 403 { // Wrong pad - do nothing
			} else if resp.StatusCode == 404 || // Wrong MAC
				(resp.StatusCode == 200 && pad == guess && pad == prev_guess) { // Forged pad matches original pad
				prev_guess, result[16-pad] = guess, byte(guess)
				log.Printf("Block %d guessed %x", block, result)
				break
			}
		}
	}

	sch <- block
	ch <- result
}

func xor(left, right []byte) []byte {
	result := make([]byte, len(left))
	for i, l := range left {
		result[i] = l ^ right[i]
	}

	return result
}
