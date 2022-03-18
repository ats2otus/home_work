package hw03frequencyanalysis

import (
	"sort"
	"strings"
	"unicode"
)

func Top10(src string) []string {
	top := make(map[string]int)
	for _, word := range strings.Fields(src) {
		if word = normalize(word); word != "" {
			top[word]++
		}
	}

	words := make([]string, 0, len(top))
	for word := range top {
		words = append(words, word)
	}
	sort.Slice(words, func(i, j int) bool {
		worda, wordb := words[i], words[j]
		if top[worda] == top[wordb] {
			return worda < wordb
		}
		return top[worda] > top[wordb]
	})

	if len(words) > 10 {
		words = words[:10]
	}
	return words
}

// убирает пунктуацию и переводит в нижний регистр
func normalize(word string) string {
	word = strings.TrimFunc(word, unicode.IsPunct)
	word = strings.ToLower(word)

	return word
}
