package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

// English letter frequency from most common to least common
var englishFrequency = "ETAOINSHRDLUCMFWYPVBGKJQXZ"

// decipherWithShift attempts to decipher text with a specific shift value
func decipherWithShift(ciphertext string, shift int) string {
	var result strings.Builder
	result.Grow(len(ciphertext))
	
	// Ensure shift is in the valid range (0-25)
	shift = shift % 26
	if shift < 0 {
		shift += 26
	}
	
	// Reverse the shift to decrypt
	shift = 26 - shift
	
	// Process each character
	for _, char := range ciphertext {
		if char >= 'A' && char <= 'Z' {
			// Handle uppercase letters
			shifted := 'A' + (char - 'A' + rune(shift)) % 26
			result.WriteRune(shifted)
		} else if char >= 'a' && char <= 'z' {
			// Handle lowercase letters
			shifted := 'a' + (char - 'a' + rune(shift)) % 26
			result.WriteRune(shifted)
		} else {
			// Non-alphabetic characters remain unchanged
			result.WriteRune(char)
		}
	}
	
	return result.String()
}

// calculateFrequencies counts letter frequencies in the text
func calculateFrequencies(text string) map[rune]int {
	freq := make(map[rune]int)
	
	// Count occurrences of each letter
	for _, char := range strings.ToUpper(text) {
		if char >= 'A' && char <= 'Z' {
			freq[char]++
		}
	}
	
	return freq
}

// getFrequencyOrder returns letters ordered by frequency (most to least common)
func getFrequencyOrder(freq map[rune]int) string {
	// Create slice of letter-frequency pairs
	type letterFreq struct {
		letter rune
		count  int
	}
	
	pairs := make([]letterFreq, 0, len(freq))
	for letter, count := range freq {
		pairs = append(pairs, letterFreq{letter, count})
	}
	
	// Sort by frequency (descending)
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].count > pairs[j].count
	})
	
	// Extract just the letters in order
	var result strings.Builder
	for _, pair := range pairs {
		result.WriteRune(pair.letter)
	}
	
	return result.String()
}

// scoreDecipheredText scores how likely the text is to be English
func scoreDecipheredText(text string) float64 {
	// Simple scoring: count common English words
	commonWords := map[string]bool{
		"THE": true, "BE": true, "TO": true, "OF": true, "AND": true,
		"A": true, "IN": true, "THAT": true, "HAVE": true, "I": true,
		"IT": true, "FOR": true, "NOT": true, "ON": true, "WITH": true,
		"HE": true, "AS": true, "YOU": true, "DO": true, "AT": true,
	}
	
	score := 0.0
	words := strings.Fields(strings.ToUpper(text))
	
	for _, word := range words {
		// Clean word of non-letters
		word = strings.Map(func(r rune) rune {
			if r >= 'A' && r <= 'Z' {
				return r
			}
			return -1
		}, word)
		
		if commonWords[word] {
			score += 1.0
		}
	}
	
	// Add bonus for text containing space distribution similar to English
	spaceCount := strings.Count(text, " ")
	spaceRatio := float64(spaceCount) / float64(len(text))
	if spaceRatio > 0.1 && spaceRatio < 0.25 {
		score += 2.0
	}
	
	return score
}

// breakCipherBruteForce tries all possible shifts and returns the best candidate
func breakCipherBruteForce(ciphertext string) (string, int) {
	bestScore := -1.0
	bestShift := 0
	bestPlaintext := ""
	
	// Try all possible shift values (0-25)
	for shift := 0; shift < 26; shift++ {
		plaintext := decipherWithShift(ciphertext, shift)
		score := scoreDecipheredText(plaintext)
		
		if score > bestScore {
			bestScore = score
			bestShift = shift
			bestPlaintext = plaintext
		}
	}
	
	return bestPlaintext, bestShift
}

// breakCipherFrequencyAnalysis uses letter frequency analysis to estimate the shift
func breakCipherFrequencyAnalysis(ciphertext string) (string, int) {
	// Only analyze letters (remove spaces, punctuation)
	lettersOnly := strings.Map(func(r rune) rune {
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
			return r
		}
		return -1
	}, ciphertext)
	
	if len(lettersOnly) < 5 {
		// Too short for reliable frequency analysis, use brute force instead
		return breakCipherBruteForce(ciphertext)
	}
	
	// Get frequency order of letters in ciphertext
	freq := calculateFrequencies(lettersOnly)
	freqOrder := getFrequencyOrder(freq)
	
	bestShift := 0
	bestScore := -1.0
	bestPlaintext := ""
	
	// Try the most likely shifts based on most common letters
	// In English, 'E' is most common, so we try aligning the most common letter with 'E' first
	potentialShifts := make([]int, 0)
	
	if len(freqOrder) > 0 {
		// Get shift if most common letter in cipher is mapped to 'E'
		mostCommon := rune(freqOrder[0])
		eShift := (mostCommon - 'E') % 26
		if eShift < 0 {
			eShift += 26
		}
		potentialShifts = append(potentialShifts, int(eShift))
	}
	
	// Add all other possible shifts
	for shift := 0; shift < 26; shift++ {
		if len(potentialShifts) == 0 || shift != potentialShifts[0] {
			potentialShifts = append(potentialShifts, shift)
		}
	}
	
	// Try potential shifts and score results
	for _, shift := range potentialShifts {
		plaintext := decipherWithShift(ciphertext, shift)
		score := scoreDecipheredText(plaintext)
		
		if score > bestScore {
			bestScore = score
			bestShift = shift
			bestPlaintext = plaintext
		}
	}
	
	return bestPlaintext, bestShift
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	
	// Get ciphertext input
	fmt.Print("Enter ciphertext to break: ")
	scanner.Scan()
	ciphertext := scanner.Text()

	// Break the cipher using both methods
	bruteForceResult, bruteForceShift := breakCipherBruteForce(ciphertext)
	freqAnalysisResult, freqAnalysisShift := breakCipherFrequencyAnalysis(ciphertext)
	
	// Display results
	fmt.Println("\nResults from brute force method:")
	fmt.Printf("Shift used: %d\n", bruteForceShift)
	fmt.Printf("Plaintext: %s\n", bruteForceResult)
	
	fmt.Println("\nResults from frequency analysis method:")
	fmt.Printf("Shift used: %d\n", freqAnalysisShift)
	fmt.Printf("Plaintext: %s\n", freqAnalysisResult)
	
	// If both methods agree, we're more confident in the result
	if bruteForceShift == freqAnalysisShift {
		fmt.Println("\nBoth methods found the same shift value, which increases confidence in the result.")
	} else {
		fmt.Println("\nThe methods found different shift values. Review both results to determine which is correct.")
	}
}