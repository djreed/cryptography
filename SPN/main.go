package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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

	N, _ := strconv.Atoi(os.Args[3])
	prettyPrint("N", strconv.Itoa(N))

	fmt.Println("//////////// Phase 1 ////////////")

	K := make([]string, N+1)
	for i := 0; i < N; i++ {
		shift := 4 * i
		K[i+1] = binKey[shift : shift+16]
		prettyPrint(fmt.Sprintf("K%v", i+1), K[i+1])
	}

	W := make([]string, N+1)
	W[0] = binPlain
	U := make([]string, N+1)
	V := make([]string, N+1)

	for r := 1; r <= N-1; r++ {
		U[r] = xor(W[r-1], K[r])
		prettyPrint(fmt.Sprintf("U%v", r), U[r])

		uHex := binToHex(U[r])
		prettyPrint(fmt.Sprintf("U%v (hex)", r), uHex)

		V[r] = convert(uHex, SBOX)
		prettyPrint(fmt.Sprintf("V%v", r), V[r])

		vBin := hexToBin(V[r])
		prettyPrint(fmt.Sprintf("V%v (binary)", r), vBin)

		W[r] = permute(vBin, PERM)
		prettyPrint(fmt.Sprintf("W%v", r), W[r])
	}

	// U[N] = xor(W[N-1], K[N])
	// prettyPrint(fmt.Sprintf("U%v", N), U[N])
	//
	// uHex := binToHex(U[N])
	// prettyPrint(fmt.Sprintf("U%v (hex)", N), uHex)
	//
	// V[N] = convert(uHex, SBOX)
	// prettyPrint(fmt.Sprintf("V%v", N), V[N])

	vBin := hexToBin(V[N-1])
	// prettyPrint(fmt.Sprintf("V%v (binary)", N), vBin)

	Y := xor(vBin, K[N])
	prettyPrint("Y", Y)

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
	// fmt.Println("strlen", len(s))
	// fmt.Println("reflen", len(ref))

	for _, target := range ref {
		permuted = permuted + string(s[target])
	}
	return
}

func convert(s string, m map[string]string) (converted string) {
	for _, char := range s {
		converted = converted + m[string(char)]
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

//Print the given string with the given label on two seperate lines
func prettyPrint(label, s string) {
	fmt.Printf("%s:\n\t%s\n", label, s)
}

//////////////////////////////////////////////////////////
////////////////// MATRIX VARIABLES //////////////////////
//////////////////////////////////////////////////////////

var (
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

	SBOX = map[string]string{
		"0": "E",
		"1": "A",
		"2": "0",
		"3": "5",
		"4": "2",
		"5": "C",
		"6": "F",
		"7": "3",
		"8": "7",
		"9": "9",
		"A": "8",
		"B": "D",
		"C": "1",
		"D": "B",
		"E": "6",
		"F": "4"}

	PERM = []int{12, 9, 0, 1, 15, 8, 11, 13, 7, 2, 3, 14, 4, 5, 6, 10}
)
