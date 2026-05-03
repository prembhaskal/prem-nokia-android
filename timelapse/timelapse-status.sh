#!/data/data/com.termux/files/usr/bin/bash
echo "Timelapse service status:"
echo "========================"
sv status timelapse
echo ""
echo "Process info:"
ps aux | grep -E "(timelapse|PID)" | grep -v grep
echo ""
echo "Recent log entries:"
echo "-------------------"
if [ -f "$HOME/.termux/sv/timelapse/logs/current" ]; then
    tail -20 "$HOME/.termux/sv/timelapse/logs/current"
else
    echo "No logs found yet"
fi

