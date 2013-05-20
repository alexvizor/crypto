package main

import "os"
import "log"
import "fmt"
import "strings"
import "net/http"
import "encoding/hex"
import "io/ioutil"

const base_url = "http://crypto-class.appspot.com/po?er="
const BLOCK_SIZE = 16

func main() {
	raw, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalln(err)
	}
	secret, err := hex.DecodeString(strings.Trim(string(raw), " \n"))
	if err != nil {
		log.Fatal("Error while decoding secret")
	}

	sln := len(secret)
	result := make([]byte, sln-BLOCK_SIZE)
	log.Println("Start guessing")

	chns, schn, chunks := make([]chan []byte, sln/BLOCK_SIZE), make(chan int), 0
	for bt := 1; bt < sln/BLOCK_SIZE; bt++ {
		chns[bt-1] = make(chan []byte)
		go oracle(schn, chns[bt-1], secret[:BLOCK_SIZE*(bt+1)])
		chunks += 1
	}

	for i := 0; i < chunks; i++ {
		bt := (<-schn) - 1
		result = append(result[:bt*BLOCK_SIZE], append(<-chns[bt], result[(bt+1)*BLOCK_SIZE:]...)...)
	}

	fmt.Printf("Result: %s\n(%x)", result, result)
}

func oracle(sch chan int, ch chan []byte, cipher []byte) {
	cln, prev_guess := len(cipher), 0
	result, fsecret := make([]byte, BLOCK_SIZE), make([]byte, cln)
	block := cln/BLOCK_SIZE - 1

	copy(fsecret, cipher) // Create forged ciphertext
	for pad := 1; pad <= BLOCK_SIZE; pad++ {
		padbyte := make([]byte, pad)
		for i := 0; i < pad; i++ {
			padbyte[i] = byte(pad)
		}

		for guess := 0x00; guess <= 0xff; guess++ {
			fsecret = append(
				fsecret[:cln-(BLOCK_SIZE+pad)],
				append(
					xor(cipher[cln-(BLOCK_SIZE+pad):cln-BLOCK_SIZE],
						xor(padbyte,
							append(
								[]byte{byte(guess)},
								result[BLOCK_SIZE-pad+1:]...,
							),
						),
					),
					fsecret[cln-BLOCK_SIZE:]...,
				)...,
			)

			resp, err := http.Get(base_url + fmt.Sprintf("%x", fsecret))
			if err != nil {
				log.Fatalln("Block", block, "Error when requesting remote host:", err)
			}

			if resp.StatusCode == 404 || // Wrong MAC
				(resp.StatusCode == 200 && pad == guess && pad == prev_guess) { // Forged pad matches original pad
				prev_guess, result[BLOCK_SIZE-pad] = guess, byte(guess)
				log.Printf("Block %d guessed %x", block, result)
				break
			} // else if resp.StatusCode == 403 { // Wrong pad - do nothing }
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
