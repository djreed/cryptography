# Code taken with permission from github.com/JacksBurden

# Probability distribution of the english language
probabilities = {'A': 0.082, 'B': 0.015, 'C': 0.028, 'D': 0.043, 'E': 0.0127, 'F': 0.022, 'G': 0.02, 'H': 0.061,
                    'I': 0.07, 'J': 0.002, 'K': 0.008, 'L': 0.04, 'M': .024, 'N': 0.067, 'O': .075, 'P': 0.019,
                    'Q': 0.001, 'R': 0.06, 'S': 0.063, 'T': 0.091, 'U' : 0.028, 'V': 0.01, 'W': 0.023,
                    'X': 0.001, 'Y': 0.02, 'Z': 0.001 }

# The four following functions were used to determine key length by manual inspection

def gcd(a, b):
    if not b:
        return a
    else:
        return gcd(b, a % b)

def gcd_list(numbers):
    return reduce(lambda x, y: gcd(x, y), numbers)


# Returns the number of times each trigram occurred in a string and where it occurred
def trigramOccurrences(string):
    trigramMap = {}
    for i in range(0, len(string) - 3):
        trigram = string[i] + string[i + 1] + string[i + 2]
        if trigramMap.get(trigram, False):
            data = trigramMap[trigram]

            count = data['count']
            indices = data['indices']
            count +=1
            indices.append(i)
            trigramMap[trigram] = { 'count': count, 'indices': indices }
        else:
            trigramMap[trigram] = { 'count': 1, 'indices': [i] }

    return trigramMap

# Gets trigrams that appear frequently in a text.
def getFrequentTrigrams(string):
    occurrences = trigramOccurrences(string)
    for key in occurrences:
        value = occurrences[key]
        if value['count'] > 1:
            print key, value

# Creates letter frequency tables for strings.
def letterOccurrences(string):
    letterMap = {}
    for i in range(0, len(string)):
        letter = string[i]
        if letterMap.get(letter, False):
            letterMap[letter] = letterMap[letter] + 1
        else:
            letterMap[letter] = 1
    return letterMap

# Constructs a Yi based on key length
def constructYi(string, i, m):
    position = i
    yi = ""
    while position < len(string):
        yi += string[position]
        position += m

    return yi

# Computes Mg values for a Yi
def Mg(string, g, nprime):
    total = 0
    letterMap = letterOccurrences(string)
    for i in range(0, 26):
        letter = chr(i + 65)
        shiftedLetter = chr(((i + g) % 26) + 65)
        value = (probabilities[letter] * letterMap.get(shiftedLetter, 0)) / nprime
        total += value

    return total


cipherText="KBPYUBACDMLRQNMGOMLGVETQVPXUQZLRQNMGOMLGVETQVPXYIMHDYQLBQUBRYILRJMTEGWYDQWEGUPGCUABRYILRJMXNQKAMHJXJKMYGVETQVPXCRWVFQNBLEZXBWTBRAQMUCAMFGAXYUWGMHTBEJBBRYILRJMLCCAHLQNWYTSGCUABRYILRJMLNTQGEQNAMRMBRYILRJMPGPBXPQNWCUXTGT"
# cipherText2="KBPYU BACDM LRQNM GOMLG VETQV PXUQZ LRQNM GOMLG VETQV PXYIM HDYQL BQUBR YILRJ MTEGW YDQWE GUPGC UABRY ILRJM XNQKA MHJXJ KMYGV ETQVP XCRWV FQNBL EZXBW TBRAQ MUCAM FGAXY UWGMH TBEJB BRYIL RJMLC CAHLQ NWYTS GCUAB RYILR JMLNT QGEQN AMRMB RYILR JMPGP BXPQN WCUXT GT"

# Helper to get a single key letter. This is not perfect as it doesn't handle ties and will
# need to be improved.
def getKeyLetter(yi, nprime):
    mg = 0
    firstLetter = None
    for i in range(0, 26):
        value = Mg(yi, i, nprime)
        if value > mg:
            mg = value
            firstLetter = i

    return chr(firstLetter + 65)

# Returns the key for a ciphertext with key length m
def getKey(string, m):
    key = ""
    nprime = len(string) / m
    for i in range(0, m):
        yi = constructYi(string, i, m)
        key += getKeyLetter(yi, nprime)

    return key

print getFrequentTrigrams(cipherText)

print getKey(cipherText, 1)
print getKey(cipherText, 2)
print getKey(cipherText, 3)
print getKey(cipherText, 4)
print getKey(cipherText, 5)
print getKey(cipherText, 6)
print getKey(cipherText, 8)
print getKey(cipherText, 9)
print getKey(cipherText, 10)