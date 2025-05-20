import collections
import string
import re

# English letter frequency from most common to least common
ENGLISH_FREQUENCY = "ETAOINSHRDLUCMFWYPVBGKJQXZ"

def decipher_with_shift(ciphertext, shift):
    """Deciphers text with a specific shift value."""
    result = []
    
    # Ensure shift is in valid range (0-25)
    shift = shift % 26
    
    # Reverse the shift to decrypt
    shift = (26 - shift) % 26
    
    for char in ciphertext:
        if 'A' <= char <= 'Z':
            # Handle uppercase letters
            shifted = chr((ord(char) - ord('A') + shift) % 26 + ord('A'))
            result.append(shifted)
        elif 'a' <= char <= 'z':
            # Handle lowercase letters
            shifted = chr((ord(char) - ord('a') + shift) % 26 + ord('a'))
            result.append(shifted)
        else:
            # Non-alphabetic characters remain unchanged
            result.append(char)
    
    return ''.join(result)

def calculate_frequencies(text):
    """Counts letter frequencies in the text."""
    # Convert to uppercase and count only letters
    text = text.upper()
    return collections.Counter(char for char in text if 'A' <= char <= 'Z')

def get_frequency_order(freq):
    """Returns letters ordered by frequency (most to least common)."""
    # Sort by count (descending) and return just the letters
    return ''.join(letter for letter, count in freq.most_common())

def score_deciphered_text(text):
    """Scores how likely the text is to be English."""
    # Common English words for scoring
    common_words = {
        "THE", "BE", "TO", "OF", "AND", 
        "A", "IN", "THAT", "HAVE", "I",
        "IT", "FOR", "NOT", "ON", "WITH",
        "HE", "AS", "YOU", "DO", "AT"
    }
    
    score = 0.0
    
    # Convert to uppercase for word matching
    text_upper = text.upper()
    
    # Remove non-letters for word analysis
    words = re.findall(r'\b[A-Z]+\b', text_upper)
    
    # Score based on common English words
    for word in words:
        if word in common_words:
            score += 1.0
    
    # Check for English-like space distribution
    space_count = text.count(' ')
    if len(text) > 0:  # Avoid division by zero
        space_ratio = space_count / len(text)
        if 0.1 < space_ratio < 0.25:  # Typical English text ratio
            score += 2.0
    
    # Add character frequency scoring
    if len(text) >= 5:  # Only for longer texts
        freq = calculate_frequencies(text)
        freq_order = get_frequency_order(freq)
        
        # Compare most common letters with English patterns
        for i, char in enumerate(freq_order[:5]):  # Check top 5 most common
            if char in ENGLISH_FREQUENCY[:6]:  # If among top 6 English letters
                score += 0.5
    
    # Detect common English letter patterns
    for pattern in ["TH", "HE", "IN", "ER", "AN", "RE"]:
        pattern_count = text_upper.count(pattern)
        score += pattern_count * 0.2
    
    return score

def break_cipher_brute_force(ciphertext):
    """Tries all possible shifts and returns the best candidate."""
    best_score = -1.0
    best_shift = 0
    best_plaintext = ""
    
    # Try all possible shift values (0-25)
    for shift in range(26):
        plaintext = decipher_with_shift(ciphertext, shift)
        score = score_deciphered_text(plaintext)
        
        if score > best_score:
            best_score = score
            best_shift = shift
            best_plaintext = plaintext
    
    return best_plaintext, best_shift

def break_cipher_frequency_analysis(ciphertext):
    """Uses letter frequency analysis to estimate the shift."""
    # Extract letters only (remove spaces, punctuation)
    letters_only = ''.join(char for char in ciphertext if char.isalpha())
    
    if len(letters_only) < 5:
        # Too short for reliable frequency analysis, use brute force instead
        return break_cipher_brute_force(ciphertext)
    
    # Get frequency order of letters in ciphertext
    freq = calculate_frequencies(letters_only)
    freq_order = get_frequency_order(freq)
    
    best_shift = 0
    best_score = -1.0
    best_plaintext = ""
    potential_shifts = []
    
    # If we have frequency data, start with most likely shift
    if freq_order:
        # Get shift if most common letter in cipher maps to 'E' (most common in English)
        most_common = freq_order[0]
        e_shift = (ord(most_common) - ord('E')) % 26
        potential_shifts.append(e_shift)
    
    # Add all other possible shifts
    for shift in range(26):
        if not potential_shifts or shift != potential_shifts[0]:
            potential_shifts.append(shift)
    
    # Try potential shifts and score results
    for shift in potential_shifts:
        plaintext = decipher_with_shift(ciphertext, shift)
        score = score_deciphered_text(plaintext)
        
        if score > best_score:
            best_score = score
            best_shift = shift
            best_plaintext = plaintext
    
    return best_plaintext, best_shift

def main():
    """Main function for the cipher breaker program."""
    # Get ciphertext input
    ciphertext = input("Enter ciphertext to break: ")
    
    # Break the cipher using both methods
    brute_force_result, brute_force_shift = break_cipher_brute_force(ciphertext)
    freq_analysis_result, freq_analysis_shift = break_cipher_frequency_analysis(ciphertext)
    
    # Display results
    print("\nResults from brute force method:")
    print(f"Shift used: {brute_force_shift}")
    print(f"Plaintext: {brute_force_result}")
    
    print("\nResults from frequency analysis method:")
    print(f"Shift used: {freq_analysis_shift}")
    print(f"Plaintext: {freq_analysis_result}")
    
    # If both methods agree, we're more confident
    if brute_force_shift == freq_analysis_shift:
        print("\nBoth methods found the same shift value, which increases confidence in the result.")
    else:
        print("\nThe methods found different shift values. Review both results to determine which is correct.")

if __name__ == "__main__":
    main()