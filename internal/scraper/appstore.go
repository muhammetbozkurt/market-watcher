package scraper

import (
	"agent/internal/domain"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"
)

type AppStoreScraper struct {
	client *http.Client
}

func NewAppStoreScraper() *AppStoreScraper {
	return &AppStoreScraper{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// XML yapıları (Apple RSS formatına göre)
type feed struct {
	Entry []entry `xml:"entry"`
}

type entry struct {
	Author  author `xml:"author"`
	Updated string `xml:"updated"`
	Rating  string `xml:"im:rating"`
	Version string `xml:"im:version"`
	Title   string `xml:"title"`
	Content string `xml:"content"` // Bazı durumlarda content array olabilir, basitleştiriyoruz
}

type author struct {
	Name string `xml:"name"`
}

func (s *AppStoreScraper) FetchReviews(appID string, lang string, limit int) ([]domain.Review, error) {
	// Örnek URL: https://itunes.apple.com/de/rss/customerreviews/id=123456789/sortBy=mostRecent/xml
	url := fmt.Sprintf("https://itunes.apple.com/%s/rss/customerreviews/id=%s/sortBy=mostRecent/xml", lang, appID)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch RSS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	var f feed
	if err := xml.NewDecoder(resp.Body).Decode(&f); err != nil {
		return nil, fmt.Errorf("failed to decode XML: %w", err)
	}

	var reviews []domain.Review
	for _, e := range f.Entry {
		// Tarih parse etme (Apple formatı: 2023-10-27T09:30:00-07:00)
		t, _ := time.Parse(time.RFC3339, e.Updated)

		// Rating string gelir, int'e çevirmek gerekir (basitlik için 0 varsayalım)
		var rating int
		fmt.Sscanf(e.Rating, "%d", &rating)

		reviews = append(reviews, domain.Review{
			Source:   "AppStore",
			Author:   e.Author.Name,
			Rating:   rating,
			Title:    e.Title,
			Content:  e.Content,
			Version:  e.Version,
			Date:     t,
			Language: lang,
		})

		if len(reviews) >= limit {
			break
		}
	}

	return reviews, nil
}
