#!/data/data/com.termux/files/usr/bin/bash
set -e

echo "=========================================="
echo "Timelapse Service Setup"
echo "=========================================="
echo ""

# Check if termux-services is installed
if ! command -v sv &> /dev/null; then
    echo "Installing termux-services..."
    pkg install termux-services -y
    echo ""
    echo "⚠️  IMPORTANT: You need to RESTART Termux completely for services to work!"
    echo "   1. Close Termux completely (swipe away from recent apps)"
    echo "   2. Open Termux again"
    echo "   3. Run this script again"
    echo ""
    exit 0
fi

# Check if timelapse binary exists
if [ ! -f "$HOME/timelapser" ]; then
    echo "❌ Error: timelapser binary not found at $HOME/timelapser"
    echo "   Please compile your Go program first:"
    echo "   GOOS=android GOARCH=arm64 go build -o timelapser timelapse.go"
    exit 1
fi

# Make sure timelapse is executable
chmod +x "$HOME/timelapser"

# Create service directory
echo "Creating service directory..."
SERVICE_DIR="$HOME/.termux/sv/timelapse"
mkdir -p "$SERVICE_DIR"
mkdir -p "$SERVICE_DIR/log"
mkdir -p "$SERVICE_DIR/logs"

# Create run script
echo "Creating service run script..."
cat > "$SERVICE_DIR/run" << 'EOF'
#!/data/data/com.termux/files/usr/bin/bash
exec 2>&1
cd $HOME
exec ./timelapser
EOF

chmod +x "$SERVICE_DIR/run"

# Create log run script
echo "Creating log service..."
cat > "$SERVICE_DIR/log/run" << 'EOF'
#!/data/data/com.termux/files/usr/bin/bash
mkdir -p $HOME/.termux/sv/timelapse/logs
exec svlogd -tt $HOME/.termux/sv/timelapse/logs
EOF

chmod +x "$SERVICE_DIR/log/run"

# Verify service directory structure
echo "Verifying service structure..."
if [ ! -x "$SERVICE_DIR/run" ]; then
    echo "❌ Error: Service run script is not executable"
    exit 1
fi

if [ ! -x "$SERVICE_DIR/log/run" ]; then
    echo "❌ Error: Log run script is not executable"
    exit 1
fi

# Enable service
echo "Enabling timelapse service..."

# First try sv-enable
sv-enable timelapse 2>/dev/null || true

# Always ensure symlink exists (in case sv-enable failed silently)
echo "Creating service symlink..."
ln -sf "$SERVICE_DIR" "$PREFIX/var/service/timelapse"

# Verify symlink was created
if [ -L "$PREFIX/var/service/timelapse" ]; then
    echo "✓ Service symlink created: $PREFIX/var/service/timelapse -> $SERVICE_DIR"
else
    echo "❌ Failed to create service symlink"
    exit 1
fi

# Wait a moment for runsvdir to detect the new service
sleep 3

echo ""
echo "✅ Service setup complete!"
echo ""
echo "Available commands:"
echo "  Start:   sv up timelapse"
echo "  Stop:    sv down timelapse"
echo "  Restart: sv restart timelapse"
echo "  Status:  sv status timelapse"
echo "  Logs:    tail -f ~/.termux/sv/timelapse/logs/current"
echo ""
echo "Starting service now..."
sv up timelapse

sleep 2
sv status timelapse

echo ""
echo "✅ Timelapse service is running!"
echo "View logs with: tail -f ~/.termux/sv/timelapse/logs/current"

