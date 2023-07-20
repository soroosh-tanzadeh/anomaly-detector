package detectors

import (
	"fmt"
	"math"
)

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

func DetectAnomalyWithSMA(traffics []float64, period int, threshold float64) []bool {
	sma := SMA(traffics, period)
	anomalies := make([]bool, len(traffics))
	for i, v := range traffics {
		if i < period {
			continue
		}
		diff := math.Abs(v - sma[i-period])
		if diff > threshold {
			fmt.Printf("Average moved from %f to %f \n", v, sma[i-period])
			anomalies[i] = true
		}
	}
	return anomalies
}
