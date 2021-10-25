package service

import (
	"context"
	"fmt"

	"github.com/go-sink/sink/internal/app/datastruct"
)

const minLinkLength = 6

// Repository interface.
type Repository interface {
	SetLink(ctx context.Context, link datastruct.Link) error
	GetLink(ctx context.Context, shortenedLink string) (datastruct.Link, error)
	GetAllLinks(ctx context.Context) ([]datastruct.Link, error)
	UpdateLink(ctx context.Context, link datastruct.Link) error
	UpdateLinks(ctx context.Context, links []datastruct.Link) error
	DeleteLink(ctx context.Context, link datastruct.Link) error
}

// RandomAlgorithm interface.
type RandomAlgorithm interface {
	RandomString(linkLength int) string
}

// RandomStringGenerator data structure.
type RandomStringGenerator struct {
	algorithm  RandomAlgorithm
	Repository Repository
}

// NewGenerator creates new instance of RandomStringGenerator.
func NewGenerator(algorithm RandomAlgorithm, repository Repository) *RandomStringGenerator {
	return &RandomStringGenerator{
		algorithm:  algorithm,
		Repository: repository,
	}
}

// Shorten the given link.
func (r *RandomStringGenerator) Shorten(ctx context.Context, link string) (string, error) {
	linkLength := minLinkLength
	encodedLink := r.algorithm.RandomString(linkLength)
	_, err := r.Repository.GetLink(ctx, encodedLink)

	for err == nil {
		linkLength = linkLength + 1
		encodedLink = r.algorithm.RandomString(linkLength)
		_, err = r.Repository.GetLink(ctx, encodedLink)
	}

	newLink := datastruct.Link{
		Original:       link,
		Shortened:      encodedLink,
		FollowQuantity: 0,
	}

	if err := r.Repository.SetLink(ctx, newLink); err != nil {
		return encodedLink, fmt.Errorf("error encoding link while setting link: %v", err)
	}

	return encodedLink, nil
}

// Unshort the given shortened link.
func (r *RandomStringGenerator) Unshort(ctx context.Context, link string) (string, error) {
	decodedLink, err := r.Repository.GetLink(ctx, link)
	if err != nil {
		return "", fmt.Errorf("error getting link from repository: %v", err)
	}

	newFollowQtty := decodedLink.FollowQuantity + 1

	newLink := datastruct.Link{
		ID:             decodedLink.ID,
		Original:       decodedLink.Original,
		Shortened:      decodedLink.Shortened,
		FollowQuantity: newFollowQtty,
		LastStatus:     decodedLink.LastStatus,
	}

	err = r.Repository.UpdateLink(ctx, newLink)
	if err != nil {
		return "", fmt.Errorf("error updating link: %v", err)
	}

	return decodedLink.Original, nil
}
