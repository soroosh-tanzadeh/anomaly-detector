package detectors

import (
	"math"
)

type DetectPoint struct {
	Detected bool
	Traffic  float64
}

func SMA(data []float64, k int) []float64 {
	sma := make([]float64, len(data)-k+1)
	for i := k; i <= len(data); i++ {
		sum := 0.0
		for _, v := range data[i-k : i] {
			sum += v
		}
		sma[i-k] = sum / float64(k)
	}
	return sma
}

func DetectAnomalyWithSMA(traffics []float64, period int, threshold float64) []DetectPoint {
	sma := SMA(traffics, period)
	anomalies := make([]DetectPoint, len(traffics))
	for i, v := range traffics {
		if i < period {
			continue
		}
		diff := math.Abs(v - sma[i-period])
		if diff > threshold {
			anomalies[i] = DetectPoint{Traffic: v, Detected: true}
		}
	}
	return anomalies
}
