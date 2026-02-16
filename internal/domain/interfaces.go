package domain

import "context"
import "time"


type MetricSnapshot struct {
	DatasetID             string // will be app_id
	Last3DaysInstalls     int64
	Last3DaysCost         float64
	Previous3DaysInstalls int64
	Previous3DaysCost     float64
}

type MetricRepository interface {
	GetMetricSnapshots(ctx context.Context) ([]MetricSnapshot, error)
}

type Anomaly struct {
	DatasetID             string // will be app_id
	Country				string
	CampaignID			string
	CampaignName		string
}



type Review struct {
	Source   string    // "AppStore" veya "GooglePlay"
	Author   string    // Kullanıcı adı
	Rating   int       // 1-5 arası puan
	Title    string    // Yorum başlığı
	Content  string    // Yorumun kendisi
	Version  string    // App versiyonu (varsa)
	Date     time.Time // Yorum tarihi
	Language string    // "de", "en", "tr" vb.
}

// Scraper, farklı kaynaklardan yorum çekmek için ortak arayüz
type Scraper interface {
	FetchReviews(appID string, lang string, limit int) ([]Review, error)
}