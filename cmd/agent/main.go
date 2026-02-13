package main

import (
	"agent/internal/config"
	"agent/internal/repository/bigquery"
	"context"
	"fmt"
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

	datasetRows, err := client.Query(context.Background(),
		bigquery.DatasetIdQuery,
	)
	if err != nil {
		log.Fatalf("failed to query bigquery: %v", err)
	}
	log.Printf("datasetRows: %v\n", datasetRows)

	for _, datasetRow := range datasetRows {
		datasetID := datasetRow["dataset_id"].(string)
		rows, err := client.Query(context.Background(),
			fmt.Sprintf(bigquery.ComparisonQuery, datasetID),
		)
		if err != nil {
			log.Fatalf("failed to query bigquery: %v", err)
		}
		log.Printf("rows: %v\n", rows)
	}

	// if err != nil {
	// 	log.Fatal("failed to load config: %w", err)
	// }
	// log.Println("config loaded: %v", cfg)

	// log.Println("test data")
	// Args := os.Args
	// log.Println(strings.Join(Args[1:], "-"))
}
