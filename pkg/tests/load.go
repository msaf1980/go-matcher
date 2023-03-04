package tests

import (
	"bufio"
	"os"
)

func LoadPatterns(filename string) []string {
	patternsFile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	patterns := make([]string, 0, 12)
	patternsReader := bufio.NewReader(patternsFile)
	for {
		pattern, err1 := patternsReader.ReadString('\n')
		if err1 != nil {
			break
		}
		patterns = append(patterns, pattern[:len(pattern)-1])
	}
	return patterns
}
