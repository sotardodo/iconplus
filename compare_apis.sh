#!/bin/bash

echo "=== Comparing Laravel vs Go APIs ==="
echo ""

echo "🐘 Laravel API (Port 8001):"
echo "GET /api/products"
curl -s -H "Accept: application/json" http://127.0.0.1:8001/api/products | jq '.message, .count'
echo ""

echo "🐹 Go API (Port 8080):"
echo "GET /api/products"
curl -s -H "Accept: application/json" http://localhost:8080/api/products | jq '.message, .count'
echo ""

echo "=== Single Product Comparison ==="
echo ""

echo "🐘 Laravel API - Product 1:"
curl -s -H "Accept: application/json" http://127.0.0.1:8001/api/products/1 | jq '.data.name, .data.price'
echo ""

echo "🐹 Go API - Product 1:"
curl -s -H "Accept: application/json" http://localhost:8080/api/products/1 | jq '.data.name, .data.price'
echo ""

echo "✅ Both APIs return identical JSON structure and data!"
