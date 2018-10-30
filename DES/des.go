package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	BLOCK_SIZE = 6
)

//Processes the given string and key through DES
func main() {
	message := strings.ToUpper(os.Args[1])
	prettyPrint("Plaintext", message)

	binPlain := hexToBin(message)
	prettyPrint("Binary Plaintext", binPlain)

	// In binary,
	key := strings.ToUpper(os.Args[2])
	prettyPrint("Key", key)

	binKey := hexToBin(key)
	prettyPrint("Binary Key", binKey)

	fmt.Println("//////////// Phase 1 ////////////")

	// Apply PC1 to the initial key
	keyPermute := permute(binKey, PC1)
	prettyPrint("Key Post-PC1", keyPermute)

	C := make([]string, 17)
	D := make([]string, 17)
	K := make([]string, 17)

	C[0] = keyPermute[:len(keyPermute)/2]
	prettyPrint("C0", C[0])
	D[0] = keyPermute[len(keyPermute)/2:]
	prettyPrint("D0", D[0])
	K[0] = binKey
	prettyPrint("K0", K[0])

	for i := 1; i <= 16; i++ {
		C[i] = shift(C[i-1], i)
		prettyPrint(fmt.Sprintf("C%v", i), C[i])

		D[i] = shift(D[i-1], i)
		prettyPrint(fmt.Sprintf("D%v", i), D[i])
	}

	//16 rounds of key generation
	for i := 1; i <= 16; i++ {
		// Apply PC2 to CD
		K[i] = permute(C[i]+D[i], PC2)
		prettyPrint(fmt.Sprintf("K%v", i), K[i])
	}

	fmt.Println("//////////// Phase 2 ////////////")

	// Initial permutation of message
	initPermute := permute(binPlain, IP)
	// prettyPrint("Message", binPlain)
	prettyPrint("Message post-IP", initPermute)

	L := make([]string, 17)
	R := make([]string, 17)

	L[0] = initPermute[:len(initPermute)/2]
	prettyPrint("L0", L[0])
	R[0] = initPermute[len(initPermute)/2:]
	prettyPrint("R0", R[0])

	// Go through 16 rounds of L, R, f
	for i := 1; i <= 16; i++ {
		L[i] = R[i-1]
		prettyPrint(fmt.Sprintf("L%v", i), L[i])

		fRes := f(R[i-1], K[i])
		prettyPrint(fmt.Sprintf("f%v", i), fRes)
		R[i] = xor(L[i-1], fRes)
		prettyPrint(fmt.Sprintf("R%v", i), R[i])
	}

	prettyPrint("R16 L16", R[16]+L[16])

	// Apply IP_INV to R(16)L(16)
	finalBin := permute(R[16]+L[16], IP_INV)
	prettyPrint("Post IP_INV in Binary", finalBin)

	// to hex
	finalHex := binToHex(finalBin)
	prettyPrint("Ciphertext in Hex", finalHex)
}

//Shift the given string a number of bits to the left
//bit shift given in the table ShiftCount, determined by which iteration of
//shift is being performed
func shift(s string, iter int) string {
	toShift := ShiftCount[iter]

	LHS := s[toShift:]
	RHS := s[:toShift]

	return LHS + RHS
}

//Convert the given hex string to binary
func hexToBin(s string) (binString string) {
	for _, c := range s {
		binString = fmt.Sprintf("%s%s", binString, binaryTable[string(c)])
	}
	return
}

//Binary string to hex representation
func binToHex(s string) (hexString string) {
	for i := 0; i < len(s); i += 4 {
		intVal, _ := strconv.ParseInt(s[i:i+4], 2, 64)
		toHex := fmt.Sprintf("%X", intVal)
		hexString = hexString + toHex
	}
	return
}

//Run the given string through a permutation matrix ref
func permute(s string, ref []int) (permuted string) {
	for _, target := range ref {
		permuted = permuted + string(s[target-1])
	}
	return
}

//Bitwise XOR between two binary strings of the same length
func xor(l, r string) (xor string) {
	for i, _ := range l {
		li, _ := strconv.Atoi(string(l[i]))
		ri, _ := strconv.Atoi(string(r[i]))

		xor = xor + strconv.Itoa((li+ri)%2)
	}
	return
}

//Run the given RHS and key through f within DES
func f(rhs string, key string) string {
	// E(R)
	e := permute(rhs, E)
	prettyPrint("E", e)

	// 8 sets of 6 bits
	blockSet := xor(key, e)
	prettyPrint("K xor E", blockSet)

	total := ""
	for i := 0; i < 8; i++ {
		blockStart := i * BLOCK_SIZE
		blockEnd := blockStart + BLOCK_SIZE
		block := blockSet[blockStart:blockEnd]
		prettyPrint(fmt.Sprintf("B%v", i+1), block)

		sRes := s(block, i)
		prettyPrint(fmt.Sprintf("S%v(B%v)", i+1, i+1), sRes)

		total = total + sRes
	}

	prettyPrint("S(B)", total)

	// apply P to S(B)
	res := permute(total, P)

	return res
}

func s(b string, n int) string {
	row := string(b[0]) + string(b[5])
	col := b[1:5]

	rowInt, _ := strconv.ParseInt(row, 2, 64)
	colInt, _ := strconv.ParseInt(col, 2, 64)

	table := S[n]
	val := table[(rowInt*16)+colInt]

	return hexToBin(fmt.Sprintf("%X", val))
}

//Print the given string with the given label on two seperate lines
func prettyPrint(label, s string) {
	fmt.Printf("%s:\n\t%s\n", label, s)
}

//////////////////////////////////////////////////////////
////////////////// MATRIX CONSTANTS //////////////////////
//////////////////////////////////////////////////////////

var (
	PC1 = []int{
		57, 49, 41, 33, 25, 17, 9,
		1, 58, 50, 42, 34, 26, 18,
		10, 2, 59, 51, 43, 35, 27,
		19, 11, 3, 60, 52, 44, 36,
		63, 55, 47, 39, 31, 23, 15,
		7, 62, 54, 46, 38, 30, 22,
		14, 6, 61, 53, 45, 37, 29,
		21, 13, 5, 28, 20, 12, 4}

	PC2 = []int{
		14, 17, 11, 24, 1, 5,
		3, 28, 15, 6, 21, 10,
		23, 19, 12, 4, 26, 8,
		16, 7, 27, 20, 13, 2,
		41, 52, 31, 37, 47, 55,
		30, 40, 51, 45, 33, 48,
		44, 49, 39, 56, 34, 53,
		46, 42, 50, 36, 29, 32}

	ShiftCount = map[int]int{
		1:  1,
		2:  1,
		3:  2,
		4:  2,
		5:  2,
		6:  2,
		7:  2,
		8:  2,
		9:  1,
		10: 2,
		11: 2,
		12: 2,
		13: 2,
		14: 2,
		15: 2,
		16: 1}

	IP = []int{
		58, 50, 42, 34, 26, 18, 10, 2,
		60, 52, 44, 36, 28, 20, 12, 4,
		62, 54, 46, 38, 30, 22, 14, 6,
		64, 56, 48, 40, 32, 24, 16, 8,
		57, 49, 41, 33, 25, 17, 9, 1,
		59, 51, 43, 35, 27, 19, 11, 3,
		61, 53, 45, 37, 29, 21, 13, 5,
		63, 55, 47, 39, 31, 23, 15, 7}

	E = []int{
		32, 1, 2, 3, 4, 5,
		4, 5, 6, 7, 8, 9,
		8, 9, 10, 11, 12, 13,
		12, 13, 14, 15, 16, 17,
		16, 17, 18, 19, 20, 21,
		20, 21, 22, 23, 24, 25,
		24, 25, 26, 27, 28, 29,
		28, 29, 30, 31, 32, 1}

	binaryTable = map[string]string{
		"0": "0000",
		"1": "0001",
		"2": "0010",
		"3": "0011",
		"4": "0100",
		"5": "0101",
		"6": "0110",
		"7": "0111",
		"8": "1000",
		"9": "1001",
		"A": "1010",
		"B": "1011",
		"C": "1100",
		"D": "1101",
		"E": "1110",
		"F": "1111"}
)

var (
	S1 = []int{
		14, 4, 13, 1, 2, 15, 11, 8, 3, 10, 6, 12, 5, 9, 0, 7,
		0, 15, 7, 4, 14, 2, 13, 1, 10, 6, 12, 11, 9, 5, 3, 8,
		4, 1, 14, 8, 13, 6, 2, 11, 15, 12, 9, 7, 3, 10, 5, 0,
		15, 12, 8, 2, 4, 9, 1, 7, 5, 11, 3, 14, 10, 0, 6, 13}

	S2 = []int{
		15, 1, 8, 14, 6, 11, 3, 4, 9, 7, 2, 13, 12, 0, 5, 10,
		3, 13, 4, 7, 15, 2, 8, 14, 12, 0, 1, 10, 6, 9, 11, 5,
		0, 14, 7, 11, 10, 4, 13, 1, 5, 8, 12, 6, 9, 3, 2, 15,
		13, 8, 10, 1, 3, 15, 4, 2, 11, 6, 7, 12, 0, 5, 14, 9}

	S3 = []int{
		10, 0, 9, 14, 6, 3, 15, 5, 1, 13, 12, 7, 11, 4, 2, 8,
		13, 7, 0, 9, 3, 4, 6, 10, 2, 8, 5, 14, 12, 11, 15, 1,
		13, 6, 4, 9, 8, 15, 3, 0, 11, 1, 2, 12, 5, 10, 14, 7,
		1, 10, 13, 0, 6, 9, 8, 7, 4, 15, 14, 3, 11, 5, 2, 12}

	S4 = []int{
		7, 13, 14, 3, 0, 6, 9, 10, 1, 2, 8, 5, 11, 12, 4, 15,
		13, 8, 11, 5, 6, 15, 0, 3, 4, 7, 2, 12, 1, 10, 14, 9,
		10, 6, 9, 0, 12, 11, 7, 13, 15, 1, 3, 14, 5, 2, 8, 4,
		3, 15, 0, 6, 10, 1, 13, 8, 9, 4, 5, 11, 12, 7, 2, 14}

	S5 = []int{
		2, 12, 4, 1, 7, 10, 11, 6, 8, 5, 3, 15, 13, 0, 14, 9,
		14, 11, 2, 12, 4, 7, 13, 1, 5, 0, 15, 10, 3, 9, 8, 6,
		4, 2, 1, 11, 10, 13, 7, 8, 15, 9, 12, 5, 6, 3, 0, 14,
		11, 8, 12, 7, 1, 14, 2, 13, 6, 15, 0, 9, 10, 4, 5, 3}

	S6 = []int{
		12, 1, 10, 15, 9, 2, 6, 8, 0, 13, 3, 4, 14, 7, 5, 11,
		10, 15, 4, 2, 7, 12, 9, 5, 6, 1, 13, 14, 0, 11, 3, 8,
		9, 14, 15, 5, 2, 8, 12, 3, 7, 0, 4, 10, 1, 13, 11, 6,
		4, 3, 2, 12, 9, 5, 15, 10, 11, 14, 1, 7, 6, 0, 8, 13}

	S7 = []int{
		4, 11, 2, 14, 15, 0, 8, 13, 3, 12, 9, 7, 5, 10, 6, 1,
		13, 0, 11, 7, 4, 9, 1, 10, 14, 3, 5, 12, 2, 15, 8, 6,
		1, 4, 11, 13, 12, 3, 7, 14, 10, 15, 6, 8, 0, 5, 9, 2,
		6, 11, 13, 8, 1, 4, 10, 7, 9, 5, 0, 15, 14, 2, 3, 12}

	S8 = []int{
		13, 2, 8, 4, 6, 15, 11, 1, 10, 9, 3, 14, 5, 0, 12, 7,
		1, 15, 13, 8, 10, 3, 7, 4, 12, 5, 6, 11, 0, 14, 9, 2,
		7, 11, 4, 1, 9, 12, 14, 2, 0, 6, 10, 13, 15, 3, 5, 8,
		2, 1, 14, 7, 4, 10, 8, 13, 15, 12, 9, 0, 3, 5, 6, 11}

	S = [][]int{S1, S2, S3, S4, S5, S6, S7, S8}

	P = []int{
		16, 7, 20, 21,
		29, 12, 28, 17,
		1, 15, 23, 26,
		5, 18, 31, 10,
		2, 8, 24, 14,
		32, 27, 3, 9,
		19, 13, 30, 6,
		22, 11, 4, 25}

	IP_INV = []int{
		40, 8, 48, 16, 56, 24, 64, 32,
		39, 7, 47, 15, 55, 23, 63, 31,
		38, 6, 46, 14, 54, 22, 62, 30,
		37, 5, 45, 13, 53, 21, 61, 29,
		36, 4, 44, 12, 52, 20, 60, 28,
		35, 3, 43, 11, 51, 19, 59, 27,
		34, 2, 42, 10, 50, 18, 58, 26,
		33, 1, 41, 9, 49, 17, 57, 25}
)
