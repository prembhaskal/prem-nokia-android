#!/data/data/com.termux/files/usr/bin/bash

echo "=========================================="
echo "Uninstall Timelapse Service"
echo "=========================================="
echo ""

# Stop and disable service
if command -v sv &> /dev/null; then
    echo "Stopping service..."
    sv down timelapse 2>/dev/null || true
    
    echo "Disabling service..."
    sv-disable timelapse 2>/dev/null || true
fi

# Remove service directory
if [ -d "$HOME/.termux/sv/timelapse" ]; then
    echo "Removing service directory..."
    rm -rf "$HOME/.termux/sv/timelapse"
fi

# Remove boot script
if [ -f "$HOME/.termux/boot/start-timelapse.sh" ]; then
    echo "Removing boot script..."
    rm -f "$HOME/.termux/boot/start-timelapse.sh"
fi

echo ""
echo "✅ Timelapse service uninstalled!"
echo ""
echo "Note: The timelapser binary at $HOME/timelapser was NOT removed"
echo "      Remove it manually if needed: rm $HOME/timelapser"

