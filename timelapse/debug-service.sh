#!/data/data/com.termux/files/usr/bin/bash

echo "=========================================="
echo "Timelapse Service Debug Information"
echo "=========================================="
echo ""

echo "1. Service directory structure:"
echo "================================"
if [ -d "$HOME/.termux/sv/timelapse" ]; then
    find "$HOME/.termux/sv/timelapse" -ls
else
    echo "❌ Service directory does not exist: $HOME/.termux/sv/timelapse"
fi
echo ""

echo "2. Binary check:"
echo "================"
if [ -f "$HOME/timelapser" ]; then
    ls -lh "$HOME/timelapser"
    file "$HOME/timelapser"
else
    echo "❌ Binary not found: $HOME/timelapser"
fi
echo ""

echo "3. Service symlink:"
echo "==================="
if [ -L "$PREFIX/var/service/timelapse" ]; then
    ls -l "$PREFIX/var/service/timelapse"
    echo "Points to: $(readlink $PREFIX/var/service/timelapse)"
else
    echo "⚠️  No symlink at: $PREFIX/var/service/timelapse"
fi
echo ""

echo "4. Service status:"
echo "=================="
sv status timelapse 2>&1 || echo "Service not running or not accessible"
echo ""

echo "5. runsvdir status:"
echo "==================="
ps aux | grep -E "(runsvdir|runsv)" | grep -v grep
echo ""

echo "6. Environment:"
echo "==============="
echo "PREFIX=$PREFIX"
echo "HOME=$HOME"
echo ""

echo "7. termux-services package:"
echo "==========================="
pkg list-installed | grep termux-services
echo ""

echo "8. Service run scripts:"
echo "======================="
if [ -f "$HOME/.termux/sv/timelapse/run" ]; then
    echo "Main run script exists:"
    ls -l "$HOME/.termux/sv/timelapse/run"
    echo "Content:"
    cat "$HOME/.termux/sv/timelapse/run"
else
    echo "❌ Main run script missing"
fi
echo ""

if [ -f "$HOME/.termux/sv/timelapse/log/run" ]; then
    echo "Log run script exists:"
    ls -l "$HOME/.termux/sv/timelapse/log/run"
else
    echo "❌ Log run script missing"
fi
echo ""

echo "9. Recent logs (if any):"
echo "========================"
if [ -f "$HOME/.termux/sv/timelapse/logs/current" ]; then
    echo "Last 20 lines:"
    tail -20 "$HOME/.termux/sv/timelapse/logs/current"
else
    echo "No logs found yet"
fi

