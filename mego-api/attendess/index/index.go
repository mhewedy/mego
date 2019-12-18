package index

import (
	"fmt"
	"github.com/mhewedy/go-conf"
	"sort"
	"strings"
)

type Input struct {
	Field string
	Ref   interface{}
}

var tokenSize = conf.GetInt("indexer.token_algo.token_size", 2)

var index map[string][]interface{}

func Index(inputs []Input) {
	index = make(map[string][]interface{})

	for _, input := range inputs {

		lower := strings.ToLower(input.Field)
		clear := removeVowels(lower)
		fields := strings.Fields(clear)

		for _, ff := range fields {
			tokens := tokenize(ff, tokenSize)

			for _, tt := range tokens {
				ii, found := index[tt]
				if !found {
					ii = make([]interface{}, 0)
				}

				if !contains(ii, input.Ref) {
					ii = append(ii, input.Ref)
					index[tt] = ii
				}
			}
		}
	}
	fmt.Println("Done indexing")
}

func Search(input string) []interface{} {

	lower := strings.ToLower(input)
	clear := removeVowels(lower)
	fields := strings.Fields(clear)

	temp := make([][]interface{}, 0)

	for _, ff := range fields {
		tokens := tokenize(ff, tokenSize)

		for _, tt := range tokens {
			ii := index[tt]
			temp = append(temp, ii)
		}
	}


	fmt.Println("temporary", temp)

	sort.Slice(temp, func(i, j int) bool {
		return len(temp[i]) < len(temp[j])
	})


	fmt.Println("temporary", temp)

	result := make([]interface{}, 0)

	for _, ii := range temp {
		for _, iii := range ii {
			if !contains(result, iii) {
				result = append(result, iii)
			}
		}
	}


	fmt.Println("result", temp)

	return result
}

func removeVowels(s string) string {

	rr := make([]rune, 0)
	for _, r := range s {
		if !strings.ContainsRune("aieou", r) {
			rr = append(rr, r)
		}
	}
	return string(rr)
}

func tokenize(s string, tokenSize int) []string {
	tokens := make([]string, 0)

	for i := range s {
		if i+tokenSize > len(s) {
			break
		}
		tokens = append(tokens, s[i:i+tokenSize])
	}
	return tokens
}

func contains(s []interface{}, t interface{}) bool {
	for _, ss := range s {
		if ss == t {
			return true
		}
	}
	return false
}
