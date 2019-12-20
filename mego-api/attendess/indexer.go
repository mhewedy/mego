package attendess

import (
	"github.com/mhewedy/go-conf"
	"sort"
	"strings"
)

type token string

var indexDB map[token][]Attendee

func index(attendees []Attendee) {
	indexDB = make(map[token][]Attendee)

	for _, aa := range attendees {
		doOnToken(aa.DisplayName, func(t token) []Attendee {
			atts, found := indexDB[t]
			if !found {
				atts = make([]Attendee, 0)
			}
			if !contains(atts, aa) {
				atts = append(atts, aa)
				indexDB[t] = atts
			}
			return nil
		})
	}
}

func search(input string) []Attendee {

	temp := doOnToken(input, func(t token) []Attendee {
		return indexDB[t]
	})

	return sortByOccurrence(temp)
}

func sortByOccurrence(temp [][]Attendee) []Attendee {
	occurrenceMap := make(map[Attendee]int)
	for _, aa := range temp {
		for _, aaa := range aa {
			_, found := occurrenceMap[aaa]
			if !found {
				occurrenceMap[aaa] = 1
			} else {
				occurrenceMap[aaa] = occurrenceMap[aaa] + 1
			}
		}
	}
	type kv struct {
		Key   Attendee
		Value int
	}
	var kvs []kv
	for k, v := range occurrenceMap {
		kvs = append(kvs, kv{k, v})
	}
	sort.Slice(kvs, func(i, j int) bool {
		return kvs[i].Value > kvs[j].Value
	})

	result := make([]Attendee, 0)
	for _, s := range kvs {
		if s.Value > conf.GetInt("indexer.min_occurrence", 0) {
			result = append(result, s.Key)
		}
	}
	return result
}

func doOnToken(input string, fn func(t token) []Attendee) [][]Attendee {

	lower := strings.ToLower(input)
	clear := substituteVowels(lower)
	fields := strings.Fields(clear)

	ii := make([][]Attendee, 0)

	for _, field := range fields {
		tokens := tokenize(field, conf.GetInt("indexer.token_size", 2))

		for _, t := range tokens {
			ii = append(ii, fn(t))
		}
	}
	return ii
}

func substituteVowels(s string) string {

	rr := make([]rune, 0)
	for _, r := range s {
		if strings.ContainsRune("aie", r) {
			rr = append(rr, 'e')
		} else if strings.ContainsRune("ou", r) {
			rr = append(rr, 'o')
		} else {
			rr = append(rr, r)
		}
	}
	return string(rr)
}

func tokenize(s string, tokenSize int) []token {
	tokens := make([]token, 0)

	for i := range s {
		if i+tokenSize > len(s) {
			break
		}
		tokens = append(tokens, token(s[i:i+tokenSize]))
	}
	if len(tokens) == 0 {
		tokens = append(tokens, token(s))
	}
	return tokens
}

func contains(attendees []Attendee, attendee Attendee) bool {
	for _, att := range attendees {
		if att.DisplayName == attendee.DisplayName {
			return true
		}
	}
	return false
}
