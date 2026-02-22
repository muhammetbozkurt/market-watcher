package main

import (
	"agent/internal/config"
	"agent/internal/domain"
	"agent/internal/intelligence"
	"agent/internal/scraper"
	"context" // Added context import
	"fmt"
	"log"
	"time"

	"github.com/tmc/langchaingo/embeddings"          // Changed from embeddings/googleai
	"github.com/tmc/langchaingo/llms/googleai"       // Added for LLM
	"github.com/tmc/langchaingo/vectorstores/chroma" // Kept for vector store
	// "agent/internal/repository/bigquery"
	// "agent/internal/service"
	// "google.golang.org/api/option"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	/*

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

	fmt.Printf("%d revies fetched succesfullly\n", len(reviews))

	ctx := context.Background()                                                                                                      // Added context initialization
	llm, err := googleai.New(ctx, googleai.WithAPIKey(cfg.GeminiAPIKey), googleai.WithDefaultEmbeddingModel("gemini-embedding-001")) // Initialized LLM
	if err != nil {
		log.Fatalf("failed to create gemini llm: %v", err) // Updated error message
	}

	embedder, err := embeddings.NewEmbedder(llm) // Initialized embedder using the LLM
	if err != nil {
		log.Fatalf("failed to create gemini embedder: %v", err)
	}

	store, err := chroma.New(
		chroma.WithChromaURL("http://localhost:8000"),
		chroma.WithEmbedder(embedder),
		chroma.WithDistanceFunction("cosine"),
		chroma.WithNameSpace("app_reviews_collection"), // Vektörlerin saklanacağı tablo (collection) adı
	)
	if err != nil {
		log.Fatalf("failed to connect to ChromaDB: %v", err)
	}

	brain, err := intelligence.NewBrain(ctx, cfg.GeminiAPIKey, store)
	if err != nil {
		log.Fatalf("failed to initialize brain: %v", err)
	}

	err = brain.IngestReviews(ctx, reviews)
	if err != nil {
		log.Fatalf("failed to ingest reviews: %v", err)
	}

	anomaly := domain.Anomaly{
		MetricName: "intro_cr_base",
		Country:    "US",
		DiffRatio:  -0.15, // %15'lik bir CR düşüşü senaryosu
		Date:       time.Now(),
	}

	analysis, err := brain.AnalyzeAnomaly(ctx, anomaly)
	if err != nil {
		log.Fatalf("failed to analyze anomaly: %v", err)
	}

	fmt.Println("\n================= AJAN RAPORU =================")
	fmt.Println(analysis)
	fmt.Println("===============================================")

	/*

		playScraper := scraper.NewPlayStoreScraper()

		fmt.Println("Fetching Play Store Reviews...")
		playReviews, err := playScraper.FetchReviews("com.maps.radar.navigation.android2023", "us", 10)
		if err != nil {
			log.Fatal(err)
		}

		for _, r := range playReviews {
			fmt.Printf("[%d stars] %s: %s\n", r.Rating, r.Author, r.Content)
		}
	*/
}
