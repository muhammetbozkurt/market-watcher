package domain

import "context"

type MetricSnapshot struct {
	DatasetID			  string	// will be app_id
	Last3DaysInstalls     int64
	Last3DaysCost         float64
	Previous3DaysInstalls int64
	Previous3DaysCost     float64
}

type MetricRepository interface {
	NewMetricRepository(client any) *MetricRepository
	GetMetricSnapshots(ctx context.Context, datasetID string) ([]MetricSnapshot, error)
}
