package utils

import "GolangPractice/Http_Learn/Parsing/vyzhenercipher"

var Key = "BanaNa"

func EncodeData(key, val string) string {
	cipher := &vyzhenercipher.VizhenerCipher{
		Key:  key,
		Text: val,
	}
	cipher.Encrypt()
	return cipher.ChangedText
}

func DecodeData(key, val string) string {
	cipher := &vyzhenercipher.VizhenerCipher{
		Key:  key,
		Text: val,
	}
	cipher.Decrypt()
	return cipher.ChangedText
}
