#!/bin/bash

# Terminate child processes on exit (CTRL+C)
trap 'kill 0' EXIT

echo "=========================================================="
echo "⚡ VOLTHUB - MEMULAI SELURUH LAYANAN SPKLU (BACKEND + FRONTEND)"
echo "=========================================================="

# Directory paths
BASE_DIR="$(cd "$(dirname "$0")" && pwd)"
BACKEND_DIR="$BASE_DIR/backend"
FRONTEND_DIR="$BASE_DIR/frontend"

# 1. Start User Service (Port 8086)
echo "🚀 Starting user-service on port 8086..."
(cd "$BACKEND_DIR/user-service" && go run main.go) &

# 2. Start Station Service (Port 8082)
echo "🚀 Starting station-service on port 8082..."
(cd "$BACKEND_DIR/station-service" && go run main.go) &

# 3. Start Booking Service (Port 8083)
echo "🚀 Starting booking-service on port 8083..."
(cd "$BACKEND_DIR/booking-service" && go run main.go) &

# 4. Start Session Service (Port 8084)
echo "🚀 Starting session-service on port 8084..."
(cd "$BACKEND_DIR/session-service" && go run main.go) &

# 5. Start Billing Service (Port 8085)
echo "🚀 Starting billing-service on port 8085..."
(cd "$BACKEND_DIR/billing-service" && go run main.go) &

# 6. Start React Frontend (Port 5173)
echo "⚡ Starting React Frontend on port 5173..."
(cd "$FRONTEND_DIR" && npm run dev) &

# Wait for background services
wait
