package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/bowei/cilium-bpf-hack/pkg/llvmp/srcnote"
)

var (
	conditionalRe = regexp.MustCompile(`^(#ifdef|#ifndef|#if|#else|#elif|#endif).*$`)
	noteRe        = regexp.MustCompile(`^.*// !note:(.*)$`)
)

func main() {
	stripPrefix := flag.String("strip", "", "Prefix to strip from the filename")

	flag.Parse()

	fileName := flag.Args()[0]
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if *stripPrefix != "" {
		fileName, _ = strings.CutPrefix(fileName, *stripPrefix)
	}

	scanner := bufio.NewScanner(f)

	var curLine int
	for scanner.Scan() {
		curLine++
		line := scanner.Text()

		switch {
		case conditionalRe.MatchString(line):
			fmt.Printf("%s:%d:%s:%s\n", fileName, curLine, srcnote.KindConditional, sanitize(line))
		case noteRe.MatchString(line):
			matches := noteRe.FindStringSubmatch(line)
			fmt.Printf("%s:%d:%s:%s\n", fileName, curLine, srcnote.KindNote, matches[1])
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func sanitize(s string) string {
	var ret string
	for _, c := range []byte(s) {
		switch {
		case ('a' <= c && c <= 'z'), ('A' <= c && c <= 'Z'), ('0' <= c && c <= '9'):
			ret = ret + string(c)
		case c == '#', c == '/', c == '*', c == '!', c == '|', c == '_':
			ret = ret + string(c)
		case c == '&':
			ret = ret + "&amp;"
		case c == '<':
			ret = ret + "&lt;"
		case c == '>':
			ret = ret + "&gt;"
		default:
			ret += " "
		}
	}
	return ret
}
