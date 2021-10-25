package cron

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/go-sink/sink/internal/app/datastruct"
	"github.com/go-sink/sink/internal/app/repository"
)

const (
	workersQtty     = 2<<7
	httpGetTimeOut  = 20 * time.Second
)

// LinkChecker interface for cron job
type LinkChecker interface {
	CheckLinks() error
}

// StatusChecker structure for cron job that checks if links are working
type StatusChecker struct {
	linkRepo repository.LinkRepository
}

// NewLinkStatusChecker constructor for StatusChecker
func NewLinkStatusChecker(linkRepo repository.LinkRepository) *StatusChecker {
	return &StatusChecker{
		linkRepo: linkRepo,
	}
}

// CheckLinks method for cron job that checks if links are working
func (c *StatusChecker) CheckLinks(ctx context.Context) {

	links, err := c.linkRepo.GetAllLinks(ctx)
	if err != nil {
		log.Error().Err(err).Msg("cron job failed")
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(len(links))

	jobs := make(chan datastruct.Link, len(links))
	updatedLinks := make(chan datastruct.Link, len(links))
	done := make(chan struct{})
	for _, link := range links {
		jobs <- link
	}

	httpClient := http.Client{}
	httpClient.Timeout = httpGetTimeOut

	for w := 0; w < workersQtty; w++ {
		go worker(httpClient, &wg, jobs, updatedLinks)
	}
	newLinks := make([]datastruct.Link, len(links))
	go func() {
		for updatedLink := range updatedLinks {
			newLinks = append(newLinks, updatedLink)
		}
		done<- struct{}{}
	}()
	wg.Wait()
	close(updatedLinks)
	close(jobs)
	<-done

	err = c.linkRepo.UpdateLinks(ctx, newLinks)
	if err != nil {
		log.Error().Err(err).Msg("cron job failed")
		return
	}

	log.Info().Msg("cron job works")
}

func worker(httpClient http.Client, wg *sync.WaitGroup, jobs <-chan datastruct.Link, updatedLinks chan<- datastruct.Link) {
	for link := range jobs {

		res, err := httpClient.Get(link.Original)

		if err != nil {
			log.Error().Err(err).Msg("could not http get")
			res.StatusCode = http.StatusInternalServerError
		}

		if link.LastStatus != res.StatusCode {
			link.LastStatus = res.StatusCode
			updatedLinks <- link
		}

		_ = res.Body.Close()

		wg.Done()
	}
}
