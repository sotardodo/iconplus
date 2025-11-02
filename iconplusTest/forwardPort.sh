#!/bin/bash
kubectl port-forward svc/frontend-service 3000:3000 -n devops &
kubectl port-forward svc/goapi-service 8080:8080 -n devops &
kubectl port-forward svc/laravel-service 8001:8001 -n devops &
kubectl port-forward svc/mysql-service 3306:3306 -n devops &
kubectl port-forward svc/prometheus-service 9090:9090 -n devops &
kubectl port-forward svc/grafana-service 3001:3001 -n devops &
kubectl port-forward svc/node-exporter 9100:9100 -n devops &
kubectl port-forward svc/kube-state-metrics 8081:8080 -n devops &
wait

