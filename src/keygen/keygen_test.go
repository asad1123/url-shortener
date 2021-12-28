package keygen

import (
	"os"
	"strings"
	"testing"

	"github.com/asad1123/url-shortener/src/config"
)

const idLength = 4

var testConfig config.AppConfig
var keygen KeyGen

func TestMain(m *testing.M) {
	testConfig.ShortenedIdLength = idLength
	keygen = KeyGen{&testConfig}
	os.Exit(m.Run())
}

func TestValidRandomString(t *testing.T) {
	str := keygen.RandomString()

	if len(str) != idLength {
		t.Errorf("ID length should be %d, but is %d", idLength, len(str))
	}

	for _, c := range str {
		if !strings.Contains(characterSpace, string(c)) {
			t.Errorf("Invalid string in ID: %c", c)
		}
	}
}

func TestUniqueRandomString(t *testing.T) {
	set := map[string]bool{}

	for i := 0; i < 1000; i++ {
		str := keygen.RandomString()
		if set[str] {
			t.Errorf("Strings are not unique")
		} else {
			set[str] = true
		}
	}
}
