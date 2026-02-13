package bigquery

import (
	"agent/internal/domain"
	"context"
	"fmt"
)

type MetricRepository struct {
	client *Client
}

func (r *MetricRepository) NewMetricRepository(client *Client) *MetricRepository {
	return &MetricRepository{client: client}
}

func (r *MetricRepository) GetMetricSnapshot(ctx context.Context, datasetID string) (*domain.MetricSnapshot, error) {
	rows, err := r.client.Query(ctx, fmt.Sprintf(ComparisonQuery, datasetID))
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("No row found for %s", datasetID)
	}
	row := rows[0]

	return &domain.MetricSnapshot{
		Last3DaysInstalls:     row["last_3_days_installs"].(int64),
		Last3DaysCost:         row["last_3_days_cost"].(float64),
		Previous3DaysInstalls: row["previous_3_days_installs"].(int64),
		Previous3DaysCost:     row["previous_3_days_cost"].(float64),
	}, nil
}
