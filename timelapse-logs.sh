#!/data/data/com.termux/files/usr/bin/bash
LOG_FILE="$HOME/.termux/sv/timelapse/logs/current"

if [ ! -f "$LOG_FILE" ]; then
    echo "❌ Log file not found: $LOG_FILE"
    echo "The service may not be running yet."
    exit 1
fi

echo "Showing timelapse logs (Ctrl+C to exit)..."
echo "==========================================="
tail -f "$LOG_FILE"

