package main

import (
	"bufio"
	"fmt"
	"os"
)

func readIn(f *os.File) any {
	r := bufio.NewReader(f)
	return parseObject(r)
}

func main() {
	// f, _ := os.Open("test.json")
	// dict := readIn(f)
	dict := readIn(os.Stdin)
	program := parseGQ(os.Args[1])

	result := transform(dict, program)

	fmt.Println(result)
}
