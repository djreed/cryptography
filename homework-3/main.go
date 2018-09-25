package main

import (
	"fmt"
	"sort"
	"strings"
)

const (
	SUBSTITUTION        = "EMGLOSUDCGDNCUSWYSFHNSFCYKDPUMLWGYICOXYSIPJCKQPKUGWGOLICGINCGACKSNISACYKZSCKXECJCKSHYSXCGOIDPKZCNKSHICGIWYGKKGKGOLDSILKGOIUSIGLEDSPWZUGFZCCNDGYYSFUSZCNXEOJNCGYEOWEUPXEZGACGNFGLKNSACIGOIYCKXCJUCIUZCFZCCNDGYYSFEUEKUZCSOCFZCCNCIACZEJNCSHFZEJZEGMXCYHCJUMGKUCY"
	VIGENERE            = "KCCPKBGUFDPHQTYAVINRRTMVGRKDNBVFDETDGILTXRGUDDKOTFMBPVGEGLTGCKQRACQCWDNAWCRXIZAKFTLEWRPTYCQKYVXCHKFTPONCQQRHJVAJUWETMCMSPKQDYHJVDAHCTRLSVSKCGCZQQDZXGSFRLSWCWSJTBHAFSIASPRJAHKJRJUMVGKMITZHFPDISPZLVLGWTFPLKKEBDPGCEBSHCTJRWXBAFSPEZQNRWXCVYCGAONWDDKACKAWBBIKFTIOVKCGGHJVLNHIFFSQESVYCLACNVRWBBIREPBBVFEXOSCDYGZWPFDTKFQIYCWHJVLNHIQIBTKHJVNPIST"
	AFFINE              = "KQEREJEBCPPCJCRKIEACUZBKRVPKRBCIBQCARBJCVFCUPKRIOFKPACUZQEPBKRXPEIIEABDKPBCPFCDCCAFIEABDKPBCPFEQPKAZBKRHAIBKAPCCIBURCCDKDCCJCIDFUIXPAFFERBICZDFKABICBBENEFCUPJCVKABPCYDCCDPKBCOCPERKIVKSCPICBRKIJPKABI"
	UNSPECIFIED         = "BNVSNSIHQCEELSSKKYERIFJKXUMBGYKAMQLJTYAVFBKVTDVBPVVRJYYLAOKYMPQSCGDLFSRLLPROYGESEBUUALRWXMMASAZLGLEDFJBZAVVPXWICGJXASCBYEHOSNMULKCEAHTQOKMFLEBKFXLRRFDTZXCIWBJSICBGAWDVYDHAVFJXZIBKCGJIWEAHTTOEWTUHKRQVVRGZBXYIREMMASCSPBNLHJMBLRFFJELHWEYLWISTFVVYFJCMHYUYRUFSFMGESIGRLWALSWMNUHSIMYYITCCQPZSICEHBCCMZFEGVJYOCDEMMPGHVAAUMELCMOEHVLTIPSUYILVGFLMVWDVYDBTHFRAYISYSGKVSUUHYHGGCKTMBLRX"
	FREQUENCY_ENGLISH   = "etaoinshrdlcumwfgypbvkjxqz"
	FREQUENCY_ENGLISH_2 = "eatoinshrdlcumwfgypbvkjxqz"
)

type CharData struct {
	char      string
	count     int
	frequency float64
}

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func main() {
	frequency, size := frequencyTable(SUBSTITUTION)
	fmt.Println(frequency, size)

	tuples := frequencyToTuples(frequency, size)
	fmt.Println(tuples)

	sorted := rankByWordCount(frequency)

	organized := pairCharacters(sorted)

	str, key := substitute(SUBSTITUTION, strings.Join(organized, ""))

	fmt.Println("String: ", str, "Key: ", key)

}

func frequencyTable(str string) (map[string]int, int) {
	m := make(map[string]int)
	for _, char := range str {
		m[string(char)]++
	}
	return m, len(str)
}

func rankByWordCount(wordFrequencies map[string]int) PairList {
	pl := make(PairList, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

func frequencyToTuples(m map[string]int, total int) []CharData {
	charData := make([]CharData, 0)
	for char, count := range m {
		charData = append(charData, CharData{char, count, float64(count) / float64(total)})
	}
	return charData
}

func substitute(cipher string, organized string) (string, map[string]string) {
	key := make(map[string]string, 0)
	for i, char := range organized {
		key[string(char)] = string(FREQUENCY_ENGLISH[i])
	}

	plaintext := cipher
	for _, char := range cipher {
		plaintext = plaintext + key[string(char)]
	}

	return plaintext, key
}

func pairCharacters(p PairList) []string {
	chars := make([]string, 0)
	for _, v := range p {
		chars = append(chars, v.Key)
	}
	return chars
}

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
