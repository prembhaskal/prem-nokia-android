#!/data/data/com.termux/files/usr/bin/bash
set -e

echo "=========================================="
echo "Setup Termux:Boot for Timelapse"
echo "=========================================="
echo ""

# Check if service is set up
if [ ! -d "$HOME/.termux/sv/timelapse" ]; then
    echo "❌ Error: Timelapse service not set up yet"
    echo "   Please run setup-service.sh first"
    exit 1
fi

# Create boot directory
echo "Creating boot script..."
mkdir -p "$HOME/.termux/boot"

# Create boot script
cat > "$HOME/.termux/boot/start-timelapse.sh" << 'EOF'
#!/data/data/com.termux/files/usr/bin/bash

# Acquire wake lock to prevent system from sleeping
termux-wake-lock

# Wait a bit for system to stabilize
sleep 10

# Start the timelapse service
sv-enable timelapse 2>/dev/null || true
sv up timelapse
EOF

chmod +x "$HOME/.termux/boot/start-timelapse.sh"

echo ""
echo "✅ Boot script created!"
echo ""
echo "📱 IMPORTANT: Install Termux:Boot app"
echo "   1. Install from F-Droid or Google Play Store"
echo "   2. Open Termux:Boot app at least once"
echo "   3. Grant any requested permissions"
echo "   4. Reboot your device to test"
echo ""
echo "The timelapse service will now start automatically on boot!"

