#!/bin/bash
# ============================================================
# deploy.sh - Script untuk deploy seluruh stack ke Kubernetes
# ============================================================
# Usage:
#   chmod +x deploy.sh
#   ./deploy.sh          # Deploy semua
#   ./deploy.sh delete   # Hapus semua
# ============================================================

set -e

NAMESPACE="sotar"
K8S_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Warna output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo -e "${CYAN}============================================================${NC}"
echo -e "${CYAN}  Sotar Kubernetes Deployment Script                         ${NC}"
echo -e "${CYAN}============================================================${NC}"

if [ "$1" == "delete" ]; then
    echo -e "${RED}⚠️  Menghapus semua resources di namespace: ${NAMESPACE}${NC}"
    read -p "Apakah Anda yakin? (y/N): " confirm
    if [ "$confirm" == "y" ] || [ "$confirm" == "Y" ]; then
        kubectl delete namespace $NAMESPACE --ignore-not-found=true
        kubectl delete clusterrolebinding vault-auth-delegator --ignore-not-found=true
        kubectl delete clusterrole prometheus --ignore-not-found=true
        kubectl delete clusterrolebinding prometheus --ignore-not-found=true
        echo -e "${GREEN}✅ Semua resources berhasil dihapus.${NC}"
    else
        echo -e "${YELLOW}Dibatalkan.${NC}"
    fi
    exit 0
fi

echo ""
echo -e "${GREEN}[1/8]${NC} 📁 Membuat Namespace..."
kubectl apply -f "$K8S_DIR/00-namespace.yaml"

echo -e "${GREEN}[2/8]${NC} 🔑 Membuat Secrets..."
kubectl apply -f "$K8S_DIR/01-secrets.yaml"

echo -e "${GREEN}[3/8]${NC} ⚙️  Membuat ConfigMaps..."
kubectl apply -f "$K8S_DIR/02-configmaps.yaml"

echo -e "${GREEN}[4/8]${NC} 🗄️  Deploying MySQL Database..."
kubectl apply -f "$K8S_DIR/10-mysql.yaml"
echo "    ⏳ Menunggu MySQL ready..."
kubectl wait --for=condition=available deployment/mysql -n $NAMESPACE --timeout=120s 2>/dev/null || true

echo -e "${GREEN}[5/8]${NC} 🚀 Deploying Applications..."
kubectl apply -f "$K8S_DIR/20-frontend.yaml"
kubectl apply -f "$K8S_DIR/21-laravel.yaml"
kubectl apply -f "$K8S_DIR/22-golang.yaml"

echo -e "${GREEN}[6/8]${NC} 🔐 Deploying HashiCorp Vault..."
kubectl apply -f "$K8S_DIR/30-vault.yaml"

echo -e "${GREEN}[7/8]${NC} 📊 Deploying Monitoring Stack..."
kubectl apply -f "$K8S_DIR/40-prometheus.yaml"
kubectl apply -f "$K8S_DIR/41-grafana.yaml"

echo -e "${GREEN}[8/8]${NC} 🌐 Applying Network Policies & Ingress..."
kubectl apply -f "$K8S_DIR/50-network-policies.yaml"
kubectl apply -f "$K8S_DIR/60-ingress.yaml"

echo ""
echo -e "${CYAN}============================================================${NC}"
echo -e "${GREEN}✅ Deployment selesai!${NC}"
echo -e "${CYAN}============================================================${NC}"
echo ""
echo -e "${YELLOW}📋 Service URLs (NodePort):${NC}"
echo -e "   Frontend   : http://<NODE_IP>:30080"
echo -e "   Vault UI   : http://<NODE_IP>:30820"
echo -e "   Prometheus  : http://<NODE_IP>:30909"
echo -e "   Grafana     : http://<NODE_IP>:30300"
echo ""
echo -e "${YELLOW}📋 Service URLs (Ingress - tambahkan ke /etc/hosts):${NC}"
echo -e "   Frontend   : http://app.sotar.local"
echo -e "   Laravel API : http://api.sotar.local"
echo -e "   Golang API  : http://go.sotar.local"
echo -e "   Grafana     : http://grafana.sotar.local"
echo -e "   Prometheus  : http://prometheus.sotar.local"
echo -e "   Vault       : http://vault.sotar.local"
echo ""
echo -e "${YELLOW}🔑 Default Credentials:${NC}"
echo -e "   Grafana  : admin / sotar123"
echo -e "   Vault    : Token = s.sotar-vault-token-2026"
echo -e "   MySQL    : root / rootP@ssw0rd"
echo -e "   MySQL    : sotar / sotar123 (DB: sotar_db)"
echo ""
echo -e "${CYAN}📌 Cek status pods:${NC}"
echo -e "   kubectl get pods -n $NAMESPACE"
echo -e "   kubectl get svc -n $NAMESPACE"
echo ""
