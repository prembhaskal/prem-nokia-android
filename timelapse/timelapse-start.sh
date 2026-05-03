#!/data/data/com.termux/files/usr/bin/bash
echo "Starting timelapse service..."
sv up timelapse
sleep 1
sv status timelapse

