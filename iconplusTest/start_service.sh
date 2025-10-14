#!/bin/bash
# =====================================
# 🔁 Auto Port Forward All Services
# =====================================
# Menjalankan port-forward semua service di namespace tertentu
# secara paralel & otomatis
# =====================================

NAMESPACE="devops"

echo "🔍 Mendeteksi semua service di namespace: $NAMESPACE ..."
SERVICES=$(kubectl get svc -n "$NAMESPACE" --no-headers | awk '{print $1}')

if [ -z "$SERVICES" ]; then
    echo "❌ Tidak ada service ditemukan di namespace $NAMESPACE"
    exit 1
fi

# Tangani Ctrl+C agar semua port-forward dimatikan bersih
trap 'echo; echo "🛑 Menghentikan semua port-forward..."; kill 0' SIGINT

for svc in $SERVICES; do
    # Ambil port pertama dari service
    PORT_INFO=$(kubectl get svc "$svc" -n "$NAMESPACE" -o jsonpath='{.spec.ports[0].port}')
    TARGET_PORT=$(kubectl get svc "$svc" -n "$NAMESPACE" -o jsonpath='{.spec.ports[0].targetPort}')

    # Jika tidak ada port, skip
    if [ -z "$PORT_INFO" ]; then
        echo "⚠️  Service $svc tidak punya port, dilewati."
        continue
    fi

    echo "🚀 Forwarding service $svc: localhost:$PORT_INFO → $svc:$TARGET_PORT"
    kubectl port-forward svc/"$svc" "$PORT_INFO":"$PORT_INFO" -n "$NAMESPACE" >/dev/null 2>&1 &
done

echo "✅ Semua port-forward dijalankan di background."
echo "Gunakan 'ps aux | grep kubectl' untuk melihat proses aktif."
echo "Tekan Ctrl+C untuk menghentikan semuanya."

# Biarkan script tetap hidup agar trap bekerja
wait

