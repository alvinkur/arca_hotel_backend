@echo off
echo ============================================
echo   Hotel Arca — Starting All Services
echo ============================================
echo.

REM Make sure .env is loaded (services read it from CWD)
cd /d %~dp0

echo Starting Auth Service (port 8001)...
start "Auth Service" cmd /k "go run services/auth-service/main.go"
timeout /t 2 >nul

echo Starting Room Service (port 8002)...
start "Room Service" cmd /k "go run services/room-service/main.go"
timeout /t 1 >nul

echo Starting Booking Service (port 8003)...
start "Booking Service" cmd /k "go run services/booking-service/main.go"
timeout /t 1 >nul

echo Starting Payment Service (port 8004)...
start "Payment Service" cmd /k "go run services/payment-service/main.go"
timeout /t 1 >nul

echo Starting Chat Service (port 8005)...
start "Chat Service" cmd /k "go run services/chat-service/main.go"
timeout /t 1 >nul

echo Starting Review Service (port 8006)...
start "Review Service" cmd /k "go run services/review-service/main.go"
timeout /t 1 >nul

echo Starting Report Service (port 8007)...
start "Report Service" cmd /k "go run services/report-service/main.go"
timeout /t 1 >nul

echo Starting AI Service (port 8008)...
start "AI Service" cmd /k "go run services/ai-service/main.go"
timeout /t 2 >nul

echo Starting API Gateway (port 8080)...
start "API Gateway" cmd /k "go run gateway/main.go"
timeout /t 2 >nul

echo.
echo ============================================
echo   All services started!
echo   Open: http://localhost:8080/login
echo.
echo   Login credentials:
echo     Owner:    owner@arca.com    / password123
echo     Staff:    staff@arca.com    / password123
echo     Customer: customer@arca.com / password123
echo.
echo   Run seed.sql first to create initial data:
echo     psql -U postgres -d db_hotel_arca -f seed.sql
echo ============================================
pause
