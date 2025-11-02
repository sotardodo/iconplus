# 🚀 Deployment Guide — iconplusTest (Laravel + Go + Frontend + MySQL)

Panduan ini menjelaskan langkah-langkah lengkap untuk:
1. Build Docker image untuk setiap komponen aplikasi.
2. Push image ke registry (contoh: Docker Hub).
3. Deploy semua komponen ke Kubernetes.

---

## 📂 Struktur Project
```bash
rootProject/
├── laravel/
│ ├── Dockerfile
│ └── (source code Laravel)
├── go/
│ ├── Dockerfile
│ └── (source code Go API)
├── frontend/
│ ├── Dockerfile
│ └── (source code Frontend)
├── k8s/
```

> 📁 Folder `k8s/` berisi semua file manifest Kubernetes.

---

## ⚙️ Prasyarat

Pastikan hal-hal berikut sudah siap:

- Akun **Docker Hub** (atau registry lain seperti GitLab, ECR, GCR)
- Tools:
  - `docker`
  - `kubectl`
  - `helm` (opsional)
- Kubernetes cluster aktif (misal: Minikube, EKS, GKE, k3s)
- Login ke Docker Hub:
  ```bash
  docker login

## 🏗️ Tahapan Build & Push Image

Ubah DOCKER_USERNAME di bawah sesuai username Docker kamu.

```bash
export DOCKER_USERNAME=sotar
export PROJECT_NAME=iconplusTest
```

## 1️⃣ Build dan Push Laravel Image
```bash
cd laravel
docker build -t $DOCKER_USERNAME/$PROJECT_NAME-laravel:latest .
docker push $DOCKER_USERNAME/$PROJECT_NAME-laravel:latest
```
## 2️⃣ Build dan Push Go API Image
```bash
cd ../go
docker build -t $DOCKER_USERNAME/$PROJECT_NAME-goapi:latest .
docker push $DOCKER_USERNAME/$PROJECT_NAME-goapi:latest
```
## 3️⃣ Build dan Push Frontend Image
```bash
cd ../frontend
docker build -t $DOCKER_USERNAME/$PROJECT_NAME-frontend:latest .
docker push $DOCKER_USERNAME/$PROJECT_NAME-frontend:latest
```

## 🧩 Update Deployment YAML
Pastikan file deployment di folder k8s/ sudah memakai image yang kamu push.


---

## ⚙️ Prasyarat

Sebelum memulai deployment, pastikan:

1. Kubernetes cluster sudah **running** (misalnya di Minikube, k3s, EKS, GKE, atau AKS).
2. Sudah terinstall tools berikut:
   - `kubectl`
   - `helm` (opsional untuk monitoring)
   - `docker` (untuk build image)
3. Image Docker untuk semua service sudah **tersedia di registry** (misal Docker Hub atau private registry).
4. Pastikan akses internet ke registry tidak diblokir dari cluster.

---

## 🏗️ Tahapan Deployment

### 1️⃣ Buat Namespace

Namespace digunakan agar resource terisolasi dalam cluster.

```bash
kubectl apply -f namespace.yaml
kubectl apply -f docker-reg-secret.yaml
kubectl apply -f mysql-secret.yaml
kubectl apply -f configmapLaravel.yaml
kubectl apply -f mysql-deployment.yaml
kubectl apply -f laravel-deployment.yaml
kubectl apply -f laravel-init-job.yaml
kubectl apply -f go-deployment.yaml
kubectl apply -f frontend-deployment.yaml
kubectl apply -f ingress.yaml
```
Pastikan Ingress Controller (misalnya NGINX Ingress) sudah aktif di cluster.
Untuk Minikube:

```bash
minikube addons enable ingress
```

## ✅ Verifikasi Deployment

Cek semua pod:
```bash
kubectl get pods -n devops
```

# Monitoring Kubernetes menggunakan Prometheus & Grafana

Dokumentasi ini menjelaskan cara melakukan instalasi **Prometheus** dan **Grafana** pada cluster Kubernetes menggunakan Helm Chart **kube-prometheus-stack**.

## Prasyarat
Pastikan sudah terinstall:
- Kubernetes Cluster (Minikube / K3s / Kind / EKS / GKE / AKS / dll)
- `kubectl`
- `helm`

## 1 Buat Namespace Monitoring
```bash
kubectl create namespace monitoring
```
## 2. Tambahkan Helm Repository Prometheus dan Grafana
```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update
```
## 3. Install Prometheus + Grafana (kube-prometheus-stack)
```bash
helm install prometheus prometheus-community/kube-prometheus-stack -n monitoring
```
### 4. Akses Grafana
Cek password default Grafana:
```bash
kubectl get secret prometheus-grafana -n monitoring -o jsonpath="{.data.admin-password}" | base64 --decode ; echo
```
## 5. Port Forward untuk Akses dari Local:
```bash
kubectl port-forward svc/prometheus-grafana 3000:80 -n monitoring
```

## forward semua port service nya untuk diakses dilokal
## ➡️  Forwarding Laravel (8001)
```bash
kubectl port-forward -n devops svc/laravel-service 8001:8001
```
## ➡️  Forwarding Frontend (3000)
```bash
kubectl port-forward -n devops svc/frontend-service 3000:3000
```
## ➡️  Forwarding Go API (8080)
```bash
kubectl port-forward -n devops svc/goapi-service 8080:8080
```
## Namespace monitoring
## ➡️  Forwarding Grafana (3001)
```bash
kubectl port-forward -n monitoring svc/prometheus-grafana 3001:80
```
## ➡️  Forwarding Prometheus (9090)
```bash
kubectl port-forward -n monitoring svc/prometheus-kube-prometheus-prometheus 9090:9090
```


saya buat diluar installer kube-prometheus-stack.tgz agar bisa diinstall secara offline