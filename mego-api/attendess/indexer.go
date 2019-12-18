package attendess

import (
	"github.com/mhewedy/go-conf"
	"sort"
	"strings"
)

var indexDB map[string][]Attendee

func index(attendees []Attendee) {
	indexDB = make(map[string][]Attendee)

	for _, aa := range attendees {
		doOnToken(aa.DisplayName, func(token string) []Attendee {
			atts, found := indexDB[token]
			if !found {
				atts = make([]Attendee, 0)
			}
			if !contains(atts, aa) {
				atts = append(atts, aa)
				indexDB[token] = atts
			}
			return nil
		})
	}
}

func search(input string) []Attendee {

	temp := doOnToken(input, func(token string) []Attendee {
		return indexDB[token]
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
		if s.Value > conf.GetInt("indexer.min_occurrence", 1) {
			result = append(result, s.Key)
		}
	}
	return result
}

func doOnToken(input string, fn func(token string) []Attendee) [][]Attendee {

	lower := strings.ToLower(input)
	clear := removeVowels(lower)
	fields := strings.Fields(clear)

	ii := make([][]Attendee, 0)

	for _, field := range fields {
		tokens := tokenize(field, conf.GetInt("indexer.token_size", 2))

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

func contains(attendees []Attendee, attendee Attendee) bool {
	for _, att := range attendees {
		if att.DisplayName == attendee.DisplayName {
			return true
		}
	}
	return false
}
