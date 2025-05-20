def apply_cipher(plaintext, shift):
    """
    Applies a substitution cipher with the given shift factor to the plaintext.
    
    Args:
        plaintext (str): The original text to encrypt
        shift (int): The number of positions to shift each letter
    
    Returns:
        str: The encrypted ciphertext
    """
    # Normalize the shift value to handle negative shifts and large values
    shift = shift % 26
    
    result = []
    for char in plaintext:
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

def main():
    """Main function to run the substitution cipher program."""
    # Get user input
    plaintext = input("Enter plaintext: ")
    
    # Get shift factor with error handling
    while True:
        try:
            shift = int(input("Enter shift factor (integer): "))
            break
        except ValueError:
            print("Please enter a valid integer.")
    
    # Apply cipher and display result
    ciphertext = apply_cipher(plaintext, shift)
    print(f"Ciphertext: {ciphertext}")

if __name__ == "__main__":
    main()