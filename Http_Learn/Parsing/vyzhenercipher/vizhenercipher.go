package vyzhenercipher

import "unicode"

type VizhenerCipher struct {
	Key         string
	Text        string
	ChangedText string
}

// Encrypt - ciphering input text (only for Ru or En letters)
func (vc *VizhenerCipher) Encrypt() *VizhenerCipher {
	var (
		keylen        = len(vc.Key)
		encryptedText []rune
	)

	// Check witch alphabet exist letter and start encrypting
	for i, r := range vc.Text {
		keyIdx := rune(vc.Key[i%keylen])
		var (
			firstChar      rune
			abcSize, shift int
		)
		if s, ok, ct := isEn(r); ok {
			abcSize = s
			if ct == 0 {
				firstChar = 'a'
			} else if ct == 1 {
				firstChar = 'A'
			}

			shift = int(keyIdx - firstChar)
		} else if s, ok, ct := isRu(r); ok {
			abcSize = s
			if ct == 0 {
				firstChar = 'а'
			} else if ct == 1 {
				firstChar = 'А'
			}

			shift = int(keyIdx - firstChar)
		} else {
			encryptedText = append(encryptedText, r)
			continue
		}

		encLetter := (int(r)-int(firstChar)+shift)%abcSize + int(firstChar)
		encryptedText = append(encryptedText, rune(encLetter))
	}

	vc.ChangedText = string(encryptedText)
	return vc
}

func isRu(r rune) (int, bool, int) {
	if unicode.IsLower(r) {
		return 33, true, 0
	} else if unicode.IsUpper(r) {
		return 34, true, 1
	} else {
		return 0, false, 0
	}
}

func isEn(r rune) (int, bool, int) {
	if unicode.IsLower(r) {
		return 26, true, 0
	} else if unicode.IsUpper(r) {
		return 26, true, 1
	} else {
		return 0, false, 0
	}
}
