package main

import (
	"log"
	"math"
	"sync"
)

// AnalyticsData - хранилище для окна данных
type AnalyticsData struct {
	mu     sync.Mutex
	window []float64
	sum    float64
}

// NewAnalytics создает новую структуру
func NewAnalytics() *AnalyticsData {
	return &AnalyticsData{
		window: make([]float64, 0, WindowSize),
	}
}

// AddAndAnalyze добавляет значение и возвращает, является ли оно аномалией
func (a *AnalyticsData) AddAndAnalyze(val float64) (bool, float64, float64) {
	a.mu.Lock()
	defer a.mu.Unlock()

	// 1. Если окно не полное — просто копим
	if len(a.window) < WindowSize {
		a.window = append(a.window, val)
		a.sum += val
		return false, 0, 0
	}

	// 2. Считаем среднее
	mean := a.sum / float64(len(a.window))

	// 3. Считаем StdDev
	var varianceSum float64
	for _, v := range a.window {
		varianceSum += math.Pow(v-mean, 2)
	}
	stdDev := math.Sqrt(varianceSum / float64(len(a.window)))

	// 4. Проверяем аномалию
	isAnomaly := false
	zScore := 0.0
	if stdDev > 0 {
		zScore = (val - mean) / stdDev
		if math.Abs(zScore) > ZScoreThreshold {
			isAnomaly = true
			log.Printf("[ANOMALY] Val: %.2f, Mean: %.2f, Z-Score: %.2f", val, mean, zScore)
		}
	}

	// 5. Сдвигаем окно (удаляем старый, добавляем новый)
	oldest := a.window[0]
	a.window = a.window[1:]
	a.window = append(a.window, val)
	a.sum = a.sum - oldest + val

	return isAnomaly, mean, zScore
}