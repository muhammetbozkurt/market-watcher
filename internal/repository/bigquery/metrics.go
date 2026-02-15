package bigquery

import (
	"agent/internal/domain"
	"context"
	"fmt"
	"log"
	"sync"
)

type MetricRepository struct {
	client *Client
}

func NewMetricRepository(client *Client) *MetricRepository {
	return &MetricRepository{client: client}
}

func (r *MetricRepository) GetMetricSnapshots(ctx context.Context) ([]domain.MetricSnapshot, error) {
	var wg sync.WaitGroup

	datasetRows, err := r.client.Query(context.Background(),
		DatasetIdQuery,
	)
	if err != nil {
		log.Fatalf("failed to query bigquery: %v", err)
	}

	results := make(chan domain.MetricSnapshot, len(datasetRows))

	for _, datasetRow := range datasetRows {
		wg.Add(1)

		go func(datasetID string) { // lambda
			defer wg.Done()
			rows, err := r.client.Query(ctx, fmt.Sprintf(ComparisonQuery, datasetID))

			if err != nil {
				return
			}
			if len(rows) == 0 {
				return
			}

			results <- domain.MetricSnapshot{
				DatasetID:             datasetID,
				Last3DaysInstalls:     rows[0]["last_3_days_installs"].(int64),
				Last3DaysCost:         rows[0]["last_3_days_cost"].(float64),
				Previous3DaysInstalls: rows[0]["previous_3_days_installs"].(int64),
				Previous3DaysCost:     rows[0]["previous_3_days_cost"].(float64),
			}
		}(datasetRow["dataset_id"].(string))
	}
	wg.Wait()
	close(results)


	var metrics []domain.MetricSnapshot
	for result := range results {
		metrics = append(metrics, result)
	}


	return metrics, nil
}
