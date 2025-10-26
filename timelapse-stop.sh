#!/data/data/com.termux/files/usr/bin/bash
echo "Stopping timelapse service..."
sv down timelapse
sleep 1
sv status timelapse

