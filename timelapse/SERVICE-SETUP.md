# Timelapse Service Setup Guide

This guide will help you set up your timelapse program as a persistent Termux service that runs in the background and survives device reboots.

## Quick Setup

### Step 1: Compile the Go Binary

```bash
cd ~/Downloads/temp/mobile
GOOS=android GOARCH=arm64 go build -o timelapser timelapse.go
mv timelapser ~/timelapser
```

### Step 2: Run the Setup Script

```bash
chmod +x setup-service.sh
./setup-service.sh
```

**Important:** If this is your first time installing termux-services, you'll need to:
1. Close Termux completely (swipe away from recent apps)
2. Open Termux again
3. Run `./setup-service.sh` again

### Step 3: Optional - Setup Auto-Start on Boot

```bash
chmod +x setup-boot.sh
./setup-boot.sh
```

Then install the **Termux:Boot** app from F-Droid or Google Play Store.

## Service Management

### Control Scripts

Make all scripts executable first:
```bash
chmod +x timelapse-*.sh
```

Available commands:

| Script | Description |
|--------|-------------|
| `./timelapse-start.sh` | Start the service |
| `./timelapse-stop.sh` | Stop the service |
| `./timelapse-restart.sh` | Restart the service |
| `./timelapse-status.sh` | Show service status and recent logs |
| `./timelapse-logs.sh` | View live logs (Ctrl+C to exit) |

### Direct Commands

You can also use `sv` commands directly:

```bash
# Start service
sv up timelapse

# Stop service
sv down timelapse

# Restart service
sv restart timelapse

# Check status
sv status timelapse

# View logs
tail -f ~/.termux/sv/timelapse/logs/current
```

## Configuration

To change timelapse settings, set environment variables in the service run script:

```bash
nano ~/.termux/sv/timelapse/run
```

Example with custom settings:
```bash
#!/data/data/com.termux/files/usr/bin/bash
exec 2>&1
cd $HOME

# Set environment variables
export CAM=0
export INTERVAL=15
export VIDEO_INTERVAL=10
export CLEANUP_INTERVAL=5
export IMAGE_RETENTION=30

exec ./timelapser
```

Then restart the service:
```bash
sv restart timelapse
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `CAM` | 0 | Camera ID (0=rear, 1=front) |
| `INTERVAL` | 10 | Capture interval in seconds |
| `VIDEO_INTERVAL` | 5 | Video generation interval in minutes |
| `CLEANUP_INTERVAL` | 5 | Cleanup check interval in minutes |
| `IMAGE_RETENTION` | 30 | How long to keep processed images (minutes) |

## Log Files

Logs are stored at: `~/.termux/sv/timelapse/logs/current`

View logs:
```bash
# Live logs
tail -f ~/.termux/sv/timelapse/logs/current

# Or use the helper script
./timelapse-logs.sh
```

## Preventing Android from Killing the App

Even with the service running, Android may kill Termux. To prevent this:

### 1. Disable Battery Optimization
- Go to **Settings → Apps → Termux**
- **Battery → Unrestricted** (or "Don't optimize")
- Enable **Allow background activity**

### 2. Manufacturer-Specific Settings

**Xiaomi:**
- Settings → Apps → Manage apps → Termux
- Enable **Autostart**
- Battery saver → **No restrictions**

**Huawei:**
- Settings → Battery → App launch → Termux
- **Manage manually** → Enable all options

**Samsung:**
- Settings → Battery → Background usage limits
- Add Termux to **Never sleeping apps**

**OnePlus/Oppo:**
- Settings → Battery → Battery Optimization → Termux → **Don't optimize**
- Settings → Apps → Termux → **Allow auto-start**

## Troubleshooting

### Service won't start
```bash
# Check if termux-services is installed
pkg list-installed | grep termux-services

# Check service status
sv status timelapse

# View error logs
cat ~/.termux/sv/timelapse/logs/current
```

### Binary not found
```bash
# Check if binary exists
ls -l ~/timelapser

# Recompile if needed
cd ~/Downloads/temp/mobile
GOOS=android GOARCH=arm64 go build -o timelapser timelapse.go
mv timelapser ~/timelapser
```

### Service starts but crashes immediately
```bash
# View logs
tail -50 ~/.termux/sv/timelapse/logs/current

# Test binary manually
~/timelapser
```

### Check disk space
```bash
df -h
du -sh ~/timelapse/*
```

## Uninstall

To remove the service completely:

```bash
chmod +x uninstall-service.sh
./uninstall-service.sh
```

This will:
- Stop the service
- Disable the service
- Remove service files
- Remove boot scripts

The timelapse binary itself will NOT be removed (do manually if needed).

## Directory Structure

After setup, you'll have:

```
~/.termux/sv/timelapse/
├── run                    # Service run script
├── log/
│   └── run               # Log service script
└── logs/
    └── current           # Current log file

~/.termux/boot/
└── start-timelapse.sh    # Auto-start script (if setup-boot.sh was run)

~/timelapse/
├── raw/
│   ├── temp/            # Images being captured
│   └── ready/           # Images ready for video
├── videos/              # Local videos (temporary)
├── processed/           # Processed images (cleaned after 30min)
└── storage/external-1/timelapse/
    └── videos/          # Final videos on SD card
```

## Monitoring

Keep an eye on your timelapse:

```bash
# Quick status check
./timelapse-status.sh

# Watch resource usage
top | grep timelapse

# Check storage usage
du -sh ~/timelapse/* ~/storage/external-1/timelapse/*

# Count images
ls ~/timelapse/raw/ready/*.jpg 2>/dev/null | wc -l
```

## Tips

1. **Test first**: Run `~/timelapser` manually to make sure it works before setting up as service
2. **Monitor logs**: Check logs regularly, especially in the first few hours
3. **Storage**: Make sure you have enough space on both internal and SD card
4. **Wake lock**: The program automatically acquires a wake lock to prevent sleep
5. **Permissions**: Make sure Termux has camera and storage permissions

## Support

If you encounter issues:
1. Check the logs: `./timelapse-logs.sh`
2. Check service status: `./timelapse-status.sh`
3. Test binary manually: `~/timelapser`
4. Check Termux permissions in Android settings
5. Verify SD card is mounted properly

