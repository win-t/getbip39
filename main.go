package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var raw = flag.Bool("raw", false, "use input as raw 32 byte entropy (hex encoded also supported)")

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Read stdin and print out bip39 24 words\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	var entropy []byte
	if *raw {
		var buf bytes.Buffer
		n, err := io.CopyN(&buf, os.Stdin, 70)
		if errors.Is(err, io.EOF) {
			err = nil
		}
		if err != nil {
			exitMsg("error on reading input: " + err.Error())
		}
		if n == 32 {
			entropy = buf.Bytes()
		} else if n >= 64 {
			entropy, err = hex.DecodeString(string(buf.Bytes()[:64]))
			if err != nil {
				exitMsg("error parsing input as hex: " + err.Error())
			}
		} else {
			exitMsg("raw input must be exactly 32 bytes")
		}
	} else {
		hash := sha256.New()
		_, err := io.Copy(hash, os.Stdin)
		if err != nil {
			exitMsg("error on reading input: " + err.Error())
		}
		entropy = hash.Sum(nil)
	}

	if len(entropy) != 32 {
		panic("invalid entropy length")
	}

	checksum := sha256.Sum256(entropy[:])

	mnemonic := append([]byte(nil), entropy[:]...)
	mnemonic = append(mnemonic, checksum[0])

	var words []string
	for i := 0; i < 24; i++ {
		low := i * 11
		high := low + 11
		wordIndex := bitSlice(mnemonic, low, high)
		words = append(words, wordlist[wordIndex])
	}
	fmt.Println(strings.Join(words, " "))
}

func bitSlice(data []byte, low, high int) int {
	var ret int
	for i := low; i < high; i++ {
		ret <<= 1
		b := data[i/8]
		b &= 1 << (7 - (i % 8))
		if b != 0 {
			ret |= 1
		}
	}
	return ret
}

func exitMsg(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
