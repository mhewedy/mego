package index

import (
	"github.com/mhewedy/go-conf"
	"sort"
	"strings"
)

type Input struct {
	Field string
	Ref   interface{}
}

var index map[string][]interface{}

func Index(inputs []Input) {
	index = make(map[string][]interface{})

	for _, input := range inputs {

		doOnToken(input.Field, func(token string) []interface{} {
			refs, found := index[token]
			if !found {
				refs = make([]interface{}, 0)
			}

			if !contains(refs, input.Ref) {
				refs = append(refs, input.Ref)
				index[token] = refs
			}
			return nil
		})
	}
}

func Search(input string) []interface{} {

	temp := doOnToken(input, func(token string) []interface{} {
		return index[token]
	})

	sort.Slice(temp, func(i, j int) bool {
		return len(temp[i]) < len(temp[j])
	})

	result := make([]interface{}, 0)
	for _, ii := range temp {
		for _, iii := range ii {
			if !contains(result, iii) {
				result = append(result, iii)
			}
		}
	}
	return result
}

func doOnToken(input string, fn func(token string) []interface{}) [][]interface{} {

	lower := strings.ToLower(input)
	clear := removeVowels(lower)
	fields := strings.Fields(clear)

	ii := make([][]interface{}, 0)

	for _, field := range fields {
		tokens := tokenize(field, conf.GetInt("indexer.token_algo.token_size", 2))

		for _, token := range tokens {
			ii = append(ii, fn(token))
		}
	}
	return ii
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
	if len(tokens) == 0 {
		tokens = append(tokens, s)
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
