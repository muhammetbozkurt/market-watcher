package service

import (
	"agent/internal/domain"
	"context"
	"fmt"
)

type WatcherService struct {
	repo domain.MetricRepository
}

func calculateCostPerInstall(installCount int64, cost float64) float64 {

	if cost == 0 {
		return  0
	}

	return float64(installCount) / cost
}

func NewWatcherService(repo domain.MetricRepository) *WatcherService {
	return &WatcherService{
		repo: repo,
	}
}

func (r *WatcherService) CheckForAnomalies(ctx context.Context) ([]domain.Anomaly, error) {
	var anomalies []domain.Anomaly
	metricSnapshots, err := r.repo.GetMetricSnapshots(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get metric snapshots: %v\n", err)
	}

	for _, snap := range metricSnapshots {
		fmt.Printf("---\n")
		fmt.Printf("snap for %s:\n", snap.DatasetID)
		fmt.Printf("\tLast3DaysInstalls:\t%d\n", snap.Last3DaysInstalls)
		fmt.Printf("\tPrevious3DaysInstalls:\t%d\n", snap.Previous3DaysInstalls)
		fmt.Printf("***\n")

		last3DaysCostPerInstall := calculateCostPerInstall(snap.Last3DaysInstalls, snap.Last3DaysCost)
		previous3DaysCostPerInstall := calculateCostPerInstall(snap.Previous3DaysInstalls, snap.Previous3DaysCost)

		if last3DaysCostPerInstall < previous3DaysCostPerInstall * 0.85 {
			anomalies = append(anomalies, domain.Anomaly{
				DatasetID: snap.DatasetID,
			})
		}
		
	}

	return anomalies, nil
}
