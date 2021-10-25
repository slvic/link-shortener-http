package urlgen_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/go-sink/sink/internal/pkg/urlgen"
)

const wantedLength = 6

func TestFunction(t *testing.T) {
	t.Run("it generates a random string of fixed length", func(t *testing.T) {
		random := rand.New(rand.NewSource(time.Now().Unix())) //nolint:gosec
		randomStringGenerator := urlgen.NewRandomURLGenerator(random)
		someString := randomStringGenerator.RandomString(wantedLength)
		assert.Len(t, someString, wantedLength)
	})
}
