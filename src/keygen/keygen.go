package keygen

import (
	"math/rand"
	"strings"
	"time"

	"github.com/asad1123/url-shortener/src/config"
)

// we shall allow any alphanumeric character in the shortened URL
const characterSpace = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type KeyGen struct {
	config *config.AppConfig
}

func NewKeyGen(config *config.AppConfig) *KeyGen {
	return &KeyGen{config}
}

func (k *KeyGen) RandomString() string {
	length := k.config.ShortenedIdLength

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
