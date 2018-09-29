package markov

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func LoadSource(path string, toLower bool) [][]string {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}

	st := strings.Split(string(raw), "\n")
	source := [][]string{}

	for _, word := range st {
		if word != "" && string(word[0]) != "#" {
			w := strings.TrimSpace(word)
			source = append(source, StringToSequence(w, toLower))
		}
	}

	return source
}

func LoadSources(paths []string, toLower bool) [][]string {
	source := [][]string{}
	for _, path := range paths {
		source = append(source, LoadSource(path, toLower)...)
	}
	return source
}

func StringToSequence(s string, toLower bool) []string {
	sequence := []string{}
	for _, c := range s {
		char := string(c)
		if toLower {
			char = strings.ToLower(char)
		}
		sequence = append(sequence, char)
	}
	return sequence
}
