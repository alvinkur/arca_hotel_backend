@echo off
echo ========================================
echo  Hotel Arca - ML Recommendation Service
echo ========================================
echo.

cd /d "%~dp0"

if not exist venv (
    echo [1/3] Creating Python virtual environment...
    python -m venv venv
)

echo [2/3] Installing dependencies...
venv\Scripts\pip install -q -r requirements.txt

echo [3/3] Starting ML service...
venv\Scripts\python app.py
