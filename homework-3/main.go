package main

import (
	"fmt"
	"math/big"
	"os"
	"sort"
)

const (
	SUBSTITUTION = "EMGLOSUDCGDNCUSWYSFHNSFCYKDPUMLWGYICOXYSIPJCKQPKUGWGOLICGINCGACKSNISACYKZSCKXECJCKSHYSXCGOIDPKZCNKSHICGIWYGKKGKGOLDSILKGOIUSIGLEDSPWZUGFZCCNDGYYSFUSZCNXEOJNCGYEOWEUPXEZGACGNFGLKNSACIGOIYCKXCJUCIUZCFZCCNDGYYSFEUEKUZCSOCFZCCNCIACZEJNCSHFZEJZEGMXCYHCJUMGKUCY"
	VIGENERE     = "KCCPKBGUFDPHQTYAVINRRTMVGRKDNBVFDETDGILTXRGUDDKOTFMBPVGEGLTGCKQRACQCWDNAWCRXIZAKFTLEWRPTYCQKYVXCHKFTPONCQQRHJVAJUWETMCMSPKQDYHJVDAHCTRLSVSKCGCZQQDZXGSFRLSWCWSJTBHAFSIASPRJAHKJRJUMVGKMITZHFPDISPZLVLGWTFPLKKEBDPGCEBSHCTJRWXBAFSPEZQNRWXCVYCGAONWDDKACKAWBBIKFTIOVKCGGHJVLNHIFFSQESVYCLACNVRWBBIREPBBVFEXOSCDYGZWPFDTKFQIYCWHJVLNHIQIBTKHJVNPIST"
	AFFINE       = "KQEREJEBCPPCJCRKIEACUZBKRVPKRBCIBQCARBJCVFCUPKRIOFKPACUZQEPBKRXPEIIEABDKPBCPFCDCCAFIEABDKPBCPFEQPKAZBKRHAIBKAPCCIBURCCDKDCCJCIDFUIXPAFFERBICZDFKABICBBENEFCUPJCVKABPCYDCCDPKBCOCPERKIVKSCPICBRKIJPKABI"
	UNSPECIFIED  = "BNVSNSIHQCEELSSKKYERIFJKXUMBGYKAMQLJTYAVFBKVTDVBPVVRJYYLAOKYMPQSCGDLFSRLLPROYGESEBUUALRWXMMASAZLGLEDFJBZAVVPXWICGJXASCBYEHOSNMULKCEAHTQOKMFLEBKFXLRRFDTZXCIWBJSICBGAWDVYDHAVFJXZIBKCGJIWEAHTTOEWTUHKRQVVRGZBXYIREMMASCSPBNLHJMBLRFFJELHWEYLWISTFVVYFJCMHYUYRUFSFMGESIGRLWALSWMNUHSIMYYITCCQPZSICEHBCCMZFEGVJYOCDEMMPGHVAAUMELCMOEHVLTIPSUYILVGFLMVWDVYDBTHFRAYISYSGKVSUUHYHGGCKTMBLRX"
	ENGLISH      = "ETAOINSHRDLCUMWFGYPBVKJXQZ"
)

type Pair struct {
	Key   rune
	Value int
}

type PairList []Pair

func main() {
	if len(os.Args) < 2 {
		panic("Must provide Affine ciphertext")
	}

	ciphertext := os.Args[1]

	frequency, _ := frequencyTable(ciphertext)
	sorted := rankByWordCount(frequency)

	keys := make([]rune, 0)
	for i, pair := range sorted {
		keys = append(keys, pair.Key)
		fmt.Printf("%v: %s with %v occurrences\n", i, string(pair.Key), pair.Value)
	}

	plaintext, a, b := solveAffine(ciphertext, keys)

	fmt.Printf("Plain: %s\n(A, B) = (%v, %v)\n", plaintext, a, b)

	// str, key := substitute(SUBSTITUTION, strings.Join(organized, ""))
	// fmt.Println("String: ", str, "Key: ", key)

}

func solveAffine(cipher string, sorted []rune) (string, int, int) {
	FIRST_SUB := 24
	SECOND_SUB := 23

	// For a certain substitution:
	cipher_first, cipher_second := runeToInt(sorted[0]), runeToInt(sorted[1])
	fmt.Printf("Cipher first: %s int %v\nCipher second: %s int %v\n", string(sorted[0]), cipher_first, string(sorted[1]), cipher_second)

	// Create system of equations
	plain_first, plain_second := runeToInt(rune(ENGLISH[FIRST_SUB])), runeToInt(rune(ENGLISH[SECOND_SUB]))
	fmt.Printf("Plain first: %s int %v\nPlain second: %s int %v\n", string(ENGLISH[FIRST_SUB]), plain_first, string(ENGLISH[SECOND_SUB]), plain_second)

	big26 := big.NewInt(26)

	// Solve system of equations for (a)
	var coeff *big.Int
	var plain *big.Int
	if cipher_first > cipher_second {
		coeff = big.NewInt(int64(cipher_first - cipher_second))
		plain = big.NewInt(int64(plain_first - plain_second))
	} else {
		coeff = big.NewInt(int64(cipher_second - cipher_first))
		plain = big.NewInt(int64(plain_second - plain_first))
	}
	inverse := coeff.ModInverse(coeff, big26)
	plain.Mod(plain, big26)
	if inverse == nil {
		panic("NO INVERSE AHH")
	}

	bigA := plain.Mul(plain, inverse)
	bigA = bigA.Mod(bigA, big26)

	a := int(bigA.Int64())
	b := int(plain_first) - (int(a) * int(cipher_first))
	bigB := big.NewInt(int64(b))
	bigB = bigB.Mod(bigB, big26)
	b = int(bigB.Int64())

	text := decryptAffine(cipher, a, b)

	return text, a, b
}

func decryptAffine(cipher string, inverse, b int) string {
	plaintext := make([]rune, 0)
	for _, char := range cipher {
		charValue := big.NewInt(runeToInt(char))

		// charEquivalent := SwapRune(char)
		charValue.Sub(charValue, big.NewInt(int64(b)))
		charValue.Mul(charValue, big.NewInt(int64(inverse)))
		charValue.Mod(charValue, big.NewInt(26))

		plaintext = append(plaintext, rune(intToChar(charValue.Int64())))
	}
	return string(plaintext)
}

func intToChar(i int64) rune {
	return rune(i + 97)
}

func runeToInt(r rune) int64 {
	switch {
	case 97 <= r && r <= 122:
		return int64(r - 97)
	case 65 <= r && r <= 90:
		return int64(r - 65)
	default:
		return 0
	}
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
	sort.Sort(sort.Reverse(pl))
	return pl
}

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
