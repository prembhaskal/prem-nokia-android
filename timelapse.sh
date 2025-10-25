#!/data/data/com.termux/files/usr/bin/bash
set -uo pipefail

CAM="${CAM:-0}"                 # 0=rear on many phones
INTERVAL="${INTERVAL:-10}"      # seconds
OUT_DIR="/data/data/com.termux/files/home/timelapse"
RAW_DIR="$OUT_DIR/raw"
mkdir -p "$RAW_DIR"
export TMPDIR=~/tmp && mkdir -p "$TMPDIR"

# Keep CPU awake while running
termux-wake-lock
trap 'termux-wake-unlock; exit 0' INT TERM EXIT

echo "Starting timelapse: cam=$CAM interval=${INTERVAL}s -> $RAW_DIR"

while :; do
  ts=$(date +%Y%m%d-%H%M%S)
  termux-torch on
  if termux-camera-photo -c "$CAM" "$RAW_DIR/$ts.jpg"; then
    termux-torch off
    f="$RAW_DIR/$ts.jpg"
    jpegtran -copy none -optimize -progressive -outfile half.jpg $f
    jpegoptim --size=250k --strip-all -q half.jpg
    mv half.jpg $f
  else
    echo "$(date -Is) camera failed; retry in 5s" >&2
    sleep 5
    continue
  fi
  termux-torch off
  sleep "$INTERVAL"
done

# TODO
# move files to sdcard after usage
# run ffmpeg parallely to keep compressing images to videos
# add logging
# run as a background service , add check only one instance running at a time.