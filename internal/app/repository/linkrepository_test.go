// +build integration

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"

	"github.com/go-sink/sink/internal/app/datastruct"
)

func TestRepository(t *testing.T) {
	linkRepository := setupTestLinkRepository(t)

	var origTestValue = fmt.Sprintf("orig-%v", time.Now().UnixNano())
	var shortTestValue = fmt.Sprintf("short-%v", time.Now().UnixNano())

	t.Run("it writes a link to a database", func(t *testing.T) {
		link := datastruct.Link{Original: origTestValue, Shortened: shortTestValue}

		err := linkRepository.SetLink(context.Background(), link)
		if err != nil {
			t.Fatalf("couldnt write a link to a database: %v", err)
		}
	})

	t.Run("it gets corresponding link", func(t *testing.T) {
		want := datastruct.Link{ID: -1, Original: origTestValue, Shortened: shortTestValue}

		got, err := linkRepository.GetLink(context.Background(), shortTestValue)
		assert.Nil(t, err)

		// because we don't know ID for sure
		got.ID = -1

		assert.Equal(t, want, got)
	})
}

func setupTestLinkRepository(t testing.TB) (linkRepository LinkRepository) {
	t.Helper()

	DSN, ok := os.LookupEnv("TEST_DSN")
	if !ok {
		log.Warn().Msg("TEST_DSN environment variable is required")
	}

	conn, err := sql.Open("postgres", DSN)
	if err != nil {
		t.Fatalf("could not establish db connection: %v", err)
	}

	return NewLinkRepository(conn)
}
