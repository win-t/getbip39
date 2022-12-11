package main

import (
	_ "embed"
	"strings"
)

//go:embed wordlist.txt
var wordlistRaw string

var wordlist []string

func init() {
	wordlist = strings.Split(wordlistRaw, "\n")[:2048]
}
