package scraper

import (
	"agent/internal/domain"
	"fmt"
	"strconv"
	"time"

	"github.com/gocolly/colly"
)

type PlayStoreScraper struct {
	collector *colly.Collector
}

func NewPlayStoreScraper() *PlayStoreScraper {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36"),
	)
	return &PlayStoreScraper{collector: c}
}

func (s *PlayStoreScraper) FetchReviews(appID string, lang string, limit int) ([]domain.Review, error) {
	var reviews []domain.Review

	url := fmt.Sprintf("https://play.google.com/store/apps/details?id=%s&hl=%s&gl=US", appID, lang)

	fmt.Println("URL:", url)

	s.collector.OnHTML("div.EGFGHd", func(e *colly.HTMLElement) {
		if len(reviews) >= limit {
			return
		}

		author := e.ChildText("div.X5PpBb")
		content := e.ChildText("div.h3YV2d")
		dateStr := e.ChildText("span.bp9Aid")

		fmt.Println("Author:", author)
		fmt.Println("Content:", content)
		fmt.Println("Date:", dateStr)

		ratingAttr := e.ChildAttr("div.iXRFPc", "aria-label")
		rating := extractRating(ratingAttr)

		reviews = append(reviews, domain.Review{
			Source:   "GooglePlay",
			Author:   author,
			Content:  content,
			Rating:   rating,
			Date:     parsePlayDate(dateStr),
			Language: lang,
		})
	})

	err := s.collector.Visit(url)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func extractRating(text string) int {
	if len(text) > 6 {
		val, _ := strconv.Atoi(string(text[6]))
		return val
	}
	return 0
}

func parsePlayDate(dateStr string) time.Time {
	// Tarih parse işi Google Play'de zordur çünkü "12 Kasım 2023" veya "Nov 12, 2023" gelir.
	// Şimdilik dummy zaman dönelim.
	return time.Now()
}
