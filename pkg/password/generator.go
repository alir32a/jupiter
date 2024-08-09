package password

import "math/rand/v2"

var charSet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const DefaultLength = 8

func NewRandomPassword(length int) string {
	var password string

	for i := 0; i < length; i++ {
		password += string(charSet[rand.IntN(len(charSet))])
	}

	return password
}
