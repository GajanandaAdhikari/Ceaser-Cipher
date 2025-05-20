package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// applyCipher applies a substitution cipher with the given shift factor to the plaintext
func applyCipher(plaintext string, shift int) string {
	var result strings.Builder
	result.Grow(len(plaintext)) // Pre-allocate space for efficiency
	
	// Handle negative shifts and large shifts (wraparound)
	shift = shift % 26
	if shift < 0 {
		shift += 26
	}
	
	// Process each character
	for _, char := range plaintext {
		if (char >= 'A' && char <= 'Z') {
			// Handle uppercase letters
			shifted := 'A' + (char - 'A' + rune(shift)) % 26
			result.WriteRune(shifted)
		} else if (char >= 'a' && char <= 'z') {
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

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	
	// Get plaintext input
	fmt.Print("Enter plaintext: ")
	scanner.Scan()
	plaintext := scanner.Text()
	
	// Get shift factor
	var shift int
	fmt.Print("Enter shift factor (integer): ")
	fmt.Scanln(&shift)
	
	// Apply cipher and output result
	ciphertext := applyCipher(plaintext, shift)
	fmt.Println("Ciphertext:", ciphertext)
}