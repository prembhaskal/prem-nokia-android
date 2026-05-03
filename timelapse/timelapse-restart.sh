#!/data/data/com.termux/files/usr/bin/bash
echo "Restarting timelapse service..."
sv restart timelapse
sleep 1
sv status timelapse

