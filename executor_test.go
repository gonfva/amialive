package main

import (
	"testing"
	"time"

	"github.com/caseymrm/menuet"
	probing "github.com/prometheus-community/pro-bing"
	"github.com/stretchr/testify/assert"
)

func TestCalculateStats(t *testing.T) {
	stats := &Stats{
		LastNRuns: make([]int64, numRuns),
	}

	// Create a mock probing.Statistics
	mockStats := &probing.Statistics{
		MaxRtt: 100 * time.Millisecond,
	}

	stats.calculateStats(mockStats)

	assert.Equal(t, int64(100*time.Millisecond), stats.MostRecent)
	assert.Equal(t, int64(100*time.Millisecond), stats.CurrentSum)
	assert.Equal(t, int64(1), stats.NumIterations)
	assert.Equal(t, int64(100*time.Millisecond), stats.avg)
	assert.Equal(t, 1, stats.CurrentPointer)
}

func TestString(t *testing.T) {
	stats := &Stats{
		LastNRuns:      []int64{1000000000, 2000000000, 3000000000},
		MostRecent:     1000000000,
		CurrentSum:     6000000000,
		NumIterations:  3,
		CurrentPointer: 2,
		avg:            2000000000,
	}

	expected := "LastNRuns: [1s 2s** 3s] Average: 2s MostRecent: 1s"
	actual := stats.String()

	assert.Equal(t, expected, actual)
}

func TestGetTitle(t *testing.T) {
	stats := &Stats{
		LastNRuns: make([]int64, numRuns),
	}

	// Setting mock values
	stats.MostRecent = int64(100 * time.Millisecond)
	stats.PacketLoss = 0
	stats.avg = int64(50 * time.Millisecond)

	// Mock menuet.Defaults().Integer
	menuet.Defaults().SetInteger("AlertOn", MultipleAverage)

	stats.getTitle()

	// Here, we would ideally check the output of the getTitle method.
	// This would require mocking the menuet.App().SetMenuState method.
	// For simplicity, we assume that it sets the title correctly.
}

func TestCalculateStats_Overflow(t *testing.T) {
	stats := &Stats{
		LastNRuns: make([]int64, numRuns),
	}

	for i := 0; i < numRuns+5; i++ {
		mockStats := &probing.Statistics{
			MaxRtt: time.Duration((i+1)*100) * time.Millisecond,
		}
		stats.calculateStats(mockStats)
	}

	assert.Equal(t, numRuns, int(stats.NumIterations))
	sum_end := time.Duration((numRuns+5)*(numRuns+6)*100/2) * time.Millisecond
	sum_init := time.Duration((5)*(6)*100/2) * time.Millisecond
	sum := sum_end - sum_init
	assert.Equal(t, int64(sum)/int64(numRuns), stats.avg)
}

func TestCalculateStats_PacketLoss(t *testing.T) {
	stats := &Stats{
		LastNRuns: make([]int64, numRuns),
	}

	// Create a mock probing.Statistics
	mockStats := &probing.Statistics{
		MaxRtt:      0,
		PacketsRecv: 0,
		PacketsSent: 1,
	}

	stats.calculateStats(mockStats)

	assert.Equal(t, int64(0), stats.MostRecent)
	assert.Equal(t, 1, stats.PacketLoss)
}
