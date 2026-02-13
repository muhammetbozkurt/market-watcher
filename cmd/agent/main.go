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

	worker.CheckForAnomalies(context.Background())
}
