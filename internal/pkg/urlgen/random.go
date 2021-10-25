package urlgen

import (
	mathRand "math/rand"
	"strings"
)

// NumberSystemConverter data structure
type NumberSystemConverter struct {
	random *mathRand.Rand
}

// NewRandomURLGenerator creates a new instance of NumberSystemConverter
func NewRandomURLGenerator(rand *mathRand.Rand) *NumberSystemConverter {
	return &NumberSystemConverter{random: rand}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// RandomString implements generation of random string using math package
func (r *NumberSystemConverter) RandomString(length int) string {
	randomNumber := r.random.Uint64()

	sb := strings.Builder{}
	sb.Grow(length)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := length-1, randomNumber, letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randomNumber, letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}
