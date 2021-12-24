package keygen

import (
	"math/rand"
	"strings"
	"time"
)

// we shall allow any alphanumeric character in the shortened URL
const characterSpace = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(length int) string {
	// ensure that we are as truly random as we can be
	rand.Seed(time.Now().UnixNano())

	sb := strings.Builder{}
	sb.Grow(length)

	i := 0
	for i < length {
		index := rand.Intn(61)
		sb.WriteByte(characterSpace[index])
		i++
	}

	return sb.String()
}
