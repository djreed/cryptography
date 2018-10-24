package main

import (
	"fmt"
	"os"
	"sort"
)

const (
// SUB1    = "BBRPLWEZQBSIQCRECQ"
// SUB2    = "CCHSCXSMWCFZWBUGKS"
// SUB3    = "RQKLWXPPCBHHVRLDR"
// SUB4    = "ROESRDEYEOCSHWNRR"
// ENGLISH = "ETAOINSHRDLCUMWFGYPBVKJXQZ"

// ENGLISH = "ABCDEFGHIJKLMNOPQRSTUVQXYZ"
)

type Pair struct {
	Key   rune
	Value int
}

type PairList []Pair

func main() {
	if len(os.Args) < 2 {
		panic("Must provide ciphertext")
	}

	ciphertext := os.Args[1]

	frequency, _ := frequencyTable(ciphertext)
	sorted := rankByWordCount(frequency)

	keys := make([]rune, 0)
	sum := 0
	for i, pair := range sorted {
		keys = append(keys, pair.Key)
		sum += pair.Value
		fmt.Printf("%v: %s with %v occurrences\n", i, string(pair.Key), pair.Value)
	}
	fmt.Println("Total: ", sum)
}

func frequencyTable(str string) (map[rune]int, int) {
	m := make(map[rune]int)
	for _, char := range str {
		m[char]++
	}
	return m, len(str)
}

func rankByWordCount(wordFrequencies map[rune]int) PairList {
	pl := make(PairList, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(pl)
	return pl
}

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Key < p[j].Key }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
