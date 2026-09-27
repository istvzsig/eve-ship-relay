package assetscanner

import (
	"context"
	"time"

	"github.com/istvzsig/eve-ship-relay/internal/shipment"
)

// FakeAssetScanner simulates the latency of an external EVE
// integration while keeping the demo deterministic.
type FakeAssetScanner struct{}

func NewFakeAssetScanner() *FakeAssetScanner {

	return &FakeAssetScanner{}
}

func (f FakeAssetScanner) Scan(
	ctx context.Context,
	s shipment.Shipment,
) (shipment.AssetScan, error) {

	select {
	case <-time.After(1500 * time.Millisecond):
		// Continue.
	case <-ctx.Done():
		return shipment.AssetScan{}, ctx.Err()
	}

	return shipment.AssetScan{
		ExpectedModules: s.AssetScan.ExpectedModules,
		ActualModules:   s.AssetScan.ActualModules,
		Passed: modulesMatch(
			s.AssetScan.ExpectedModules,
			s.AssetScan.ActualModules,
		),
	}, nil
}

func modulesMatch(expected, actual []string) bool {
	if len(expected) != len(actual) {
		return false
	}

	expectedSet := make(map[string]struct{}, len(expected))

	for _, module := range expected {
		expectedSet[module] = struct{}{}
	}

	for _, module := range actual {
		if _, ok := expectedSet[module]; !ok {
			return false
		}
	}

	return true
}
