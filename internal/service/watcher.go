package service

import (
	"agent/internal/domain"
	"context"
	"fmt"
)

type WatcherService struct {
	repo domain.MetricRepository
}

func NewWatcherService(repo domain.MetricRepository) *WatcherService {
	return &WatcherService{
		repo: repo,
	}
}

func (r *WatcherService) CheckForAnomalies(ctx context.Context) {
	metricSnapshots, err := r.repo.GetMetricSnapshots(ctx)
	if err != nil {
		fmt.Printf("failed to get metric snapshots: %v\n", err)
		return
	}

	for _, snap := range metricSnapshots {
		fmt.Printf("---\n")
		fmt.Printf("snap for %s:\n", snap.DatasetID)
		fmt.Printf("\tLast3DaysInstalls:\t%d\n", snap.Last3DaysInstalls)
		fmt.Printf("\tPrevious3DaysInstalls:\t%d\n", snap.Previous3DaysInstalls)
		fmt.Printf("***\n")
	}

}
