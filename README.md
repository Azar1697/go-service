# 🚀 Highload Go Service with K8s & Observability

![Go Version](https://img.shields.io/badge/Go-1.22-blue) ![Kubernetes](https://img.shields.io/badge/Kubernetes-Kind-326ce5) ![Prometheus](https://img.shields.io/badge/Monitoring-Prometheus%20%26%20Grafana-orange)

Высоконагруженный микросервис для обработки потоковых данных (IoT-метрик) с реализацией паттерна Worker Pool и статистической аналитикой аномалий (Z-Score). Проект развернут в Kubernetes с полным стеком мониторинга и автоматическим масштабированием (HPA).

> Особенность проекта: Инфраструктура оптимизирована для работы в условиях жестких ограничений ресурсов (Edge/VPS 2GB RAM). Реализован Fine-Grained Resource Tuning для компонентов Ingress и Monitoring.

---

## 🏗 Архитектура

Система построена на микросервисной архитектуре с использованием асинхронной обработки данных.

* Core: Go (Golang) + Goroutines (Fan-Out pattern).
* Storage: Redis Cluster (кэширование и дедупликация).
* Orchestration: Kubernetes (Kind).
* Traffic: NGINX Ingress Controller.
* Observability: Prometheus + Grafana + Alertmanager.

### Схема потоков данных
Client -> Ingress (L7) -> Service -> Go Pods -> Buffered Channel -> Workers <-> Redis

---

## 📊 Основные возможности

* High Performance: Обработка >1000 RPS на одном узле благодаря буферизированным каналам.
* Anomaly Detection: "Налету" вычисляет скользящую среднюю и Z-Score (threshold > 2σ) без использования тяжелых ML-библиотек.
* Auto-Scaling: Настроен HPA (Horizontal Pod Autoscaler), который автоматически увеличивает число реплик при CPU load > 50%.
* Robust Monitoring: Дашборды Grafana для отслеживания памяти, сети и состояния Redis.

---

## 🛠 Установка и запуск

### Предварительные требования
* Docker
* Kind (Kubernetes in Docker)
* Kubectl
* Helm

### 1. Запуск кластера
```bash
kind create cluster --name highload-cluster
```

### 2. Установка зависимостей (Helm)
Устанавливаем Redis и Prometheus с оптимизированными флагами (для экономии памяти):

# Redis
```bash
helm install redis oci://registry-1.docker.io/bitnami/charts/redis
```
# Prometheus (Lite version without node-exporter)
```bash
helm install prometheus prometheus-community/prometheus \
  --set nodeExporter.enabled=false \
  --set alertmanager.resources.limits.memory=32Mi
```

## 3. Развертывание приложения
# Сборка образа (если локально)
```bash
docker build -t go-service:latest .
kind load docker-image go-service:latest --name highload-cluster
```
# Применение манифестов
```bash
kubectl apply -f k8s/
```
## 🧩 Структура проекта

```PLAINTEXT
.
├── k8s/                   # Kubernetes Manifests
│   ├── deployment.yaml    # App Deployment config
│   ├── service.yaml       # ClusterIP Service
│   ├── ingress.yaml       # NGINX Ingress rules
│   └── hpa.yaml           # Auto-scaling rules
├── src/                   # Go Source Code
│   ├── main.go            # Entry point
│   ├── workers.go         # Worker Pool logic
│   └── analytics.go       # Math logic (Z-Score)
├── Dockerfile             # Multi-stage build
└── README.md              # Documentation
```