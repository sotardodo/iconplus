#!/bin/bash
# =========================================
# Port-forward Prometheus dan Grafana
# Namespace: monitoring
# =========================================

echo "🚀 Starting port-forwarding for monitoring namespace..."

# Jalankan Prometheus di background
kubectl port-forward svc/prometheus-service -n monitoring 9090:9090 >/dev/null 2>&1 &
PROM_PID=$!

# Jalankan Grafana di background
kubectl port-forward svc/grafana-service -n monitoring 3030:3000 >/dev/null 2>&1 &
GRAFANA_PID=$!

echo "✅ Prometheus -> http://localhost:9090"
echo "✅ Grafana    -> http://localhost:3030"
echo ""
echo "Tekan [CTRL+C] untuk stop port-forwarding."

# Tunggu hingga user stop script
wait