// +build integration

package service

import (
	"context"
	"database/sql"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"

	"github.com/go-sink/sink/internal/app/datastruct"
	"github.com/go-sink/sink/internal/app/repository"
	"github.com/go-sink/sink/internal/pkg/urlgen"
)

func TestURLShortener(t *testing.T) {
	testShortener := setUpTestShortener(t)

	testLink := datastruct.Link{
		Original:    "somelink.com",
		Shortened:   "undefined",
	}

	t.Run("it encodes and decodes a link", func(t *testing.T) {
		encoded, err := testShortener.Shorten(context.TODO(), testLink.Original)
		if err != nil {
			t.Fatalf("could not encode link: %v", err)
		}

		original, err := testShortener.Unshort(context.TODO(), encoded)
		if err != nil {
			t.Fatalf("could not decode a link: %v", err)
		}

		err = testShortener.Repository.DeleteLink(context.TODO(), testLink)
		if err != nil {
			t.Fatalf("could not delete link: %v", err)
		}

		assert.Equal(t, testLink.Original, original)
	})

}

func setUpTestShortener(t *testing.T) *RandomStringGenerator {
	t.Helper()

	random := rand.New(rand.NewSource(time.Now().Unix())) //nolint:gosec

	repo := setupTestLinkRepository(t)
	return NewGenerator(urlgen.NewRandomURLGenerator(random), repo)
}

func setupTestLinkRepository(t testing.TB) (linkRepository repository.LinkRepository) {
	t.Helper()

	DSN, ok := os.LookupEnv("TEST_DSN")
	if !ok {
		log.Warn().Msg("TEST_DSN environment variable is required")
	}

	conn, err := sql.Open("postgres", DSN)
	if err != nil {
		t.Fatalf("could not establish db connection: %v", err)
	}

	return repository.NewLinkRepository(conn)
}
