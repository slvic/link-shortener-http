package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/rs/zerolog/log"

	"github.com/go-sink/sink/internal/app/datastruct"
)

// LinkRepository data structure.
type LinkRepository struct {
	database *sql.DB
}

// NewLinkRepository creates new LinkRepository instance.
func NewLinkRepository(db *sql.DB) LinkRepository {
	return LinkRepository{database: db}
}

func (r LinkRepository) GetAllLinks(ctx context.Context) ([]datastruct.Link, error) {
	var links []datastruct.Link

	rows, err := r.database.QueryContext(ctx, "SELECT id, original, shortened, follow_qtty from links")
	if err != nil {
		return nil, fmt.Errorf("could not execute a query %w", err)
	}
	defer func(rows *sql.Rows) {
		closeErr := rows.Close()
		if closeErr != nil {
			log.Error().Err(closeErr).Msg("could not close cursor after getting link")
		}
	}(rows)

	for rows.Next() {
		link := datastruct.Link{}
		err = rows.Scan(&link.ID, &link.Original, &link.Shortened, &link.FollowQuantity)
		if err != nil {
			return nil, fmt.Errorf("could not scan a row: %w", err)
		}
		links = append(links, link)
	}

	return links, nil
}

// GetLink from database.
func (r LinkRepository) GetLink(ctx context.Context, short string) (datastruct.Link, error) {
	var link datastruct.Link

	rows, err := r.database.QueryContext(ctx, "SELECT id, original, shortened, follow_qtty from links where shortened = $1", short)
	if err != nil {
		return link, fmt.Errorf("could not execute a query %w", err)
	}
	defer func(rows *sql.Rows) {
		closeErr := rows.Close()
		if closeErr != nil {
			log.Error().Err(closeErr).Msg("could not close cursor after getting link")
		}
	}(rows)

	if !rows.Next() {
		return link, fmt.Errorf("coud not find original link for %v", short)
	}

	err = rows.Scan(&link.ID, &link.Original, &link.Shortened, &link.FollowQuantity)
	if err != nil {
		return link, fmt.Errorf("could not scan a row: %w", err)
	}

	return link, nil
}

// SetLink to database.
func (r LinkRepository) SetLink(ctx context.Context, link datastruct.Link) error {
	rows, err := r.database.QueryContext(ctx, "INSERT INTO links(original, shortened, follow_qtty) VALUES ($1, $2, $3)", link.Original, link.Shortened, link.FollowQuantity)
	if err != nil {
		return err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Error().Err(err).Msg("could not close cursor after setting link")
		}
	}(rows)
	return nil
}

// UpdateLink in database.
func (r LinkRepository) UpdateLink(ctx context.Context, link datastruct.Link) error {
	rows, err := r.database.QueryContext(ctx, "UPDATE links SET original = $1, shortened = $2, follow_qtty = $3 WHERE id = $4", link.Original, link.Shortened, link.FollowQuantity, link.ID)
	if err != nil {
		return err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Error().Err(err).Msg("could not close cursor after updating link")
		}
	}(rows)
	return nil
}

func (r LinkRepository) UpdateLinks(ctx context.Context, links []datastruct.Link) error {
	groupedLinks := make(map[int][]int64)
	for _, link := range links {
		groupedLinks[link.LastStatus] = append(groupedLinks[link.LastStatus], link.ID)
	}

	for statusCode, ids := range groupedLinks {
		rows, err := r.database.QueryContext(ctx, "UPDATE links SET last_status = $1 WHERE id = ANY($2)", statusCode, pq.Int64Array(ids))

		if err != nil {
			return fmt.Errorf("could not update links: %w", err)
		}
		_ = rows.Close()
	}
	return nil
}

// DeleteLink from database.
func (r LinkRepository) DeleteLink(ctx context.Context, link datastruct.Link) error {
	rows, err := r.database.QueryContext(ctx, "DELETE FROM links WHERE id = $1", link.ID)
	if err != nil {
		return err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Error().Err(err).Msg("could not close cursor after deleting link")
		}
	}(rows)
	return nil
}
