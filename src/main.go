package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os" 

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
)

// Конфигурация
const (
	WindowSize      = 50
	ZScoreThreshold = 2.0
)

// Глобальные переменные
var (
	rdb      *redis.Client
	dataChan chan float64

	// Метрики Prometheus
	opsProcessed = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "Total processed events",
	})
	anomaliesDetected = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "myapp_anomalies_total",
		Help: "Total anomalies detected",
	})
	currentLoad = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "myapp_current_load_value",
		Help: "Current value of the metric being processed",
	})
)

func init() {
	prometheus.MustRegister(opsProcessed)
	prometheus.MustRegister(anomaliesDetected)
	prometheus.MustRegister(currentLoad)
}

// Воркер
func startWorker(analytics *AnalyticsData) {
	for val := range dataChan {
		// 1. Считаем
		isAnomaly, _, _ := analytics.AddAndAnalyze(val)

		// 2. Обновляем метрики
		opsProcessed.Inc()
		currentLoad.Set(val) // Добавил обновление Gauge
		if isAnomaly {
			anomaliesDetected.Inc()
		}

		// 3. Пишем в Redis (Fire and forget)
		go func(v float64, anom bool) {
			ctx := context.Background()
			// Игнорируем ошибку для теста, чтобы не засорять логи, если Redis недоступен
			_ = rdb.Set(ctx, "last_value", fmt.Sprintf("%f", v), 0)
			if anom {
				_ = rdb.Incr(ctx, "anomaly_count")
			}
		}(val, isAnomaly)
	}
}

func main() {
	// Читаем адрес из переменной окружения.
	// Если переменной нет (мы запускаем локально), используем localhost.
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	
	log.Printf("Connecting to Redis at: %s", redisAddr) 

	// 1. Инит Redis с динамическим адресом
	rdb = redis.NewClient(&redis.Options{Addr: redisAddr})

	// 2. Инит каналов и аналитики
	dataChan = make(chan float64, 1000)
	analytics := NewAnalytics()

	// 3. Запуск воркера
	go startWorker(analytics)

	// 4. Роутинг
	http.HandleFunc("/data", handleInput)
	http.HandleFunc("/health", handleHealth)
	http.Handle("/metrics", promhttp.Handler())

	log.Println("Service started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}