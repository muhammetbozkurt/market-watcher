package main

import (
	"agent/internal/config"
	"agent/internal/domain"
	"agent/internal/intelligence"
	"agent/internal/repository/bigquery"
	"agent/internal/scraper"
	"agent/internal/service"
	"context" // Added context import
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/tmc/langchaingo/embeddings"          // Changed from embeddings/googleai
	"github.com/tmc/langchaingo/llms/googleai"       // Added for LLM
	"github.com/tmc/langchaingo/vectorstores/chroma" // Kept for vector store
	"google.golang.org/api/option"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	ctx := context.Background()

	llm, err := googleai.New(ctx, googleai.WithAPIKey(cfg.GeminiAPIKey), googleai.WithDefaultEmbeddingModel("gemini-embedding-001"))
	if err != nil {
		log.Fatalf("failed to create gemini llm: %v", err)
	}

	embedder, err := embeddings.NewEmbedder(llm)
	if err != nil {
		log.Fatalf("failed to create gemini embedder: %v", err)
	}

	store, err := chroma.New(
		chroma.WithChromaURL("http://localhost:8000"),
		chroma.WithEmbedder(embedder),
		chroma.WithDistanceFunction("cosine"),
		chroma.WithNameSpace("app_reviews_collection"),
	)
	if err != nil {
		log.Fatalf("failed to connect to ChromaDB: %v", err)
	}

	brain, err := intelligence.NewBrain(ctx, cfg.GeminiAPIKey, store)
	if err != nil {
		log.Fatalf("failed to initialize brain: %v", err)
	}

	appScraper := scraper.NewAppStoreScraper()
	playScraper := scraper.NewPlayStoreScraper()

	client, err := bigquery.NewClient(ctx, cfg.BigQueryProjectID, option.WithCredentialsJSON([]byte(cfg.BigQueryServiceAccountJSON)))
	if err != nil {
		log.Fatalf("failed to create bigquery client: %v", err)
	}
	defer client.Close()

	bigqueryRepo := bigquery.NewMetricRepository(client)
	worker := service.NewWatcherService(bigqueryRepo)

	anomalies, err := worker.CheckForAnomalies(ctx)
	if err != nil {
		log.Fatalf("failed to get anomalies: %v", err)
	}

	log.Printf("---------Anomalies---------")
	for _, anomaly := range anomalies {
		log.Printf("Processing anomaly for dataset: %s", anomaly.DatasetID)

		var reviews []domain.Review
		var fetchErr error

		if _, errInt := strconv.Atoi(anomaly.AppID); errInt == nil {
			fmt.Printf("Fetching App Store Reviews for %s...\n", anomaly.AppID)
			reviews, fetchErr = appScraper.FetchReviews(anomaly.AppID, "us", 10)
		} else {
			fmt.Printf("Fetching Play Store Reviews for %s...\n", anomaly.AppID)
			reviews, fetchErr = playScraper.FetchReviews(anomaly.AppID, "us", 10)
		}

		if fetchErr != nil {
			log.Printf("failed to fetch reviews for app_id %s (dataset: %s): %v", anomaly.AppID, anomaly.DatasetID, fetchErr)
			continue
		}

		fmt.Printf("%d reviews fetched successfully for %s\n", len(reviews), anomaly.AppID)

		err = brain.IngestReviews(ctx, reviews)
		if err != nil {
			log.Printf("failed to ingest reviews for %s: %v", anomaly.DatasetID, err)
			continue
		}

		anomaly.MetricName = "CostPerInstall"
		anomaly.Date = time.Now()

		analysis, err := brain.AnalyzeAnomaly(ctx, anomaly)
		if err != nil {
			log.Printf("failed to analyze anomaly for %s: %v", anomaly.DatasetID, err)
			continue
		}

		fmt.Println("\n================= AGENT REPORT =================")
		fmt.Printf("DatasetID: %s\n", anomaly.DatasetID)
		fmt.Println(analysis)
		fmt.Println("===============================================")
	}
}
