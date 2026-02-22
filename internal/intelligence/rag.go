package intelligence

import (
	"agent/internal/domain"
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
)

type Brain struct {
	llm    *googleai.GoogleAI
	vector vectorstores.VectorStore
}

func NewBrain(ctx context.Context, apiKey string, store vectorstores.VectorStore) (*Brain, error) {
	llm, err := googleai.New(
		ctx,
		googleai.WithAPIKey(apiKey),
		googleai.WithDefaultModel("gemini-2.5-flash"),
	)

	if err != nil {
		return nil, fmt.Errorf("gemini couldn't start: %w\n", err)
	}

	return &Brain{
		llm:    llm,
		vector: store,
	}, nil
}

func (b *Brain) IngestReviews(ctx context.Context, reviews []domain.Review) error {
	var docs []schema.Document

	for _, r := range reviews {
		doc := schema.Document{
			PageContent: fmt.Sprintf("Title: %s\nReview:%s", r.Title, r.Content),
			Metadata: map[string]any{
				"source":   r.Source,
				"rating":   r.Rating,
				"author":   r.Author,
				"language": r.Language,
				"date":     r.Date.Format("2006-01-02"),
			},
		}
		docs = append(docs, doc)

	}

	_, err := b.vector.AddDocuments(ctx, docs)
	if err != nil {
		return fmt.Errorf("vectordb insert error: %w", err)
	}

	return nil
}

func (b *Brain) AnalyzeAnomaly(ctx context.Context, anomaly domain.Anomaly) (string, error) {
	searchQuery := fmt.Sprintf("In %s, the metric %s changed by %.2f%%. I am looking for price issues, bugs, complaints, or competitor moves based on the following context.",
		anomaly.Country, anomaly.MetricName, anomaly.DiffRatio*100)

	docs, err := b.vector.SimilaritySearch(ctx, searchQuery, 15) // Gemini'ın context limiti yüksek olduğu için sayıyı artırabiliriz
	if err != nil {
		return "", fmt.Errorf("semantic similarty search error: %w", err)
	}

	var contextText string
	for i, doc := range docs {
		contextText += fmt.Sprintf("[Yorum %d]: %s\n", i+1, doc.PageContent)
	}

	template := `You are a skilled Product and Marketing Analyst.
The following anomaly has been detected in our internal data:
Metric: {{.metric}}
Country: {{.country}}
Change: {{.diff}}%

Below are the store reviews received from users during this period:
{{.context}}

Your Tasks:
Explain the reason for the decline in the metric based on these reviews.
Suggest 3 actionable steps for the Product or Marketing teams to take.
If the reviews do not contain sufficient information to explain the decline, simply state: "Insufficient signals found in the reviews" and do not speculate.`

	prompt := prompts.NewPromptTemplate(template, []string{"metric", "country", "diff", "context"})

	promptValue, err := prompt.Format(map[string]any{
		"metric":  anomaly.MetricName,
		"country": anomaly.Country,
		"diff":    anomaly.DiffRatio * 100,
		"context": contextText,
	})
	if err != nil {
		return "", err
	}

	// langchaingo'da modele metin göndermenin standart ve en güvenli yolu
	completion, err := llms.GenerateFromSinglePrompt(ctx, b.llm, promptValue)
	if err != nil {
		return "", fmt.Errorf("gemini response generation: %w", err)
	}

	return completion, nil
}
