package vyzhenercipher

import (
	"unicode"
)

type VizhenerCipher struct {
	Key         string
	Text        string
	ChangedText string
}

// Encrypt - encrypting text to Vigener cipher
func (vc *VizhenerCipher) Encrypt() *VizhenerCipher {
	vc.ChangedText = processText(vc.Text, vc.Key, 1)
	return vc
}

// Decrypt - decrypting text to Vigener cipher
func (vc *VizhenerCipher) Decrypt() *VizhenerCipher {
	vc.ChangedText = processText(vc.Text, vc.Key, -1)
	return vc
}

// processText - process text and shift him
func processText(text, key string, dire int) string {
	var (
		processedText []rune
		keyLen        = len(key)
		idx           = 0
	)
	for _, r := range text {
		var (
			keyIdx        = rune(key[idx%keyLen])
			processedRune rune
		)
		if isEn(r) {
			shift := int(keyIdx-'A') % 26
			if unicode.IsLower(r) {
				shift = int(keyIdx-'a') % 26
			}

			processedRune = shiftRune(r, shift*dire, 26)
			idx++
		} else if isRu(r) {
			shift := int(keyIdx-'А') % 32
			if unicode.IsLower(r) {
				shift = int(keyIdx-'а') % 32
			}

			processedRune = shiftRune(r, shift*dire, 32)
			idx++
		} else {
			processedText = append(processedText, r)
			continue
		}

		processedText = append(processedText, processedRune)
	}

	return string(processedText)
}

// shiftRune - shift rune for alphabet on key rune
func shiftRune(r rune, shift int, alphabetSize int) rune {
	var base rune
	if unicode.IsLower(r) {
		base = 'a'
		if alphabetSize == 32 {
			base = 'а'
		}
	} else {
		base = 'A'
		if alphabetSize == 32 {
			base = 'А'
		}
	}

	return base + (r-base+rune(shift)+rune(alphabetSize))%rune(alphabetSize)
}

// isRu - check letter is Ru?
func isRu(r rune) bool {
	return (r >= 'а' && r <= 'я') || (r >= 'А' && r <= 'Я')
}

// isEn - check letter is En?
func isEn(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}