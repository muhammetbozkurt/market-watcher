package main

import (
	"agent/internal/config"
	"agent/internal/repository/bigquery"
	"agent/internal/service"
	"context"
	"log"

	"google.golang.org/api/option"
)

func main() {
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
}
