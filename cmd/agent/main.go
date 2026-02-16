package main

import (
	// "agent/internal/config"
	// "agent/internal/repository/bigquery"
	// "agent/internal/service"
	"agent/internal/scraper"
	// "context"
	"fmt"
	"log"
	// "google.golang.org/api/option"
)

func main() {
	/*
		cfg, err := config.LoadConfig()
		if err != nil {
			log.Fatalf("failed to load config: %v", err)
		}
		client, err := bigquery.NewClient(context.Background(), cfg.BigQueryProjectID, option.WithCredentialsJSON(cfg.BigQueryServiceAccountJSON))
		if err != nil {
			log.Fatalf("failed to create bigquery client: %v", err)
		}
		defer client.Close()

		bigqueryRepo := bigquery.NewMetricRepository(client)
		worker := service.NewWatcherService(bigqueryRepo)


		anomalies, err := worker.CheckForAnomalies(context.Background())
		if err != nil {
			log.Fatalf("failed to get anomalies: %v", err)
		}

		log.Printf("---------Anomalies---------")
		for _, anomaly := range anomalies {
			log.Printf("anomaly: %w\n", anomaly)
		}
	*/
	appScraper := scraper.NewAppStoreScraper()

	fmt.Println("Fetching App Store Reviews...")
	reviews, err := appScraper.FetchReviews("6463052823", "us", 10)
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range reviews {
		fmt.Printf("[%d stars] %s: %s\n", r.Rating, r.Author, r.Title)
	}

	playScraper := scraper.NewPlayStoreScraper()

	fmt.Println("Fetching Play Store Reviews...")
	playReviews, err := playScraper.FetchReviews("com.maps.radar.navigation.android2023", "us", 10)
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range playReviews {
		fmt.Printf("[%d stars] %s: %s\n", r.Rating, r.Author, r.Content)
	}
}
