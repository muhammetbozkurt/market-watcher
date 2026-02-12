package main

import (
	"agent/internal/config"
	"agent/internal/repository/bigquery"
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

	rows, err := client.Query(context.Background(), "SELECT install_date, round(AVG(trial_cr), 2) trial_cr FROM `bi-data-430410.comaichatbotassistantapp2024.trial_cr_base` WHERE install_date >= DATE_SUB(CURRENT_DATE(), INTERVAL 2 week) group by 1 order by 1 desc LIMIT 2")
	if err != nil {
		log.Fatalf("failed to query bigquery: %v", err)
	}
	log.Printf("rows: %v", rows)

	// if err != nil {
	// 	log.Fatal("failed to load config: %w", err)
	// }
	// log.Println("config loaded: %v", cfg)

	// log.Println("test data")
	// Args := os.Args
	// log.Println(strings.Join(Args[1:], "-"))
}
