package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
	"syscall"
	"time"
)

type Config struct {
	cam                string
	captureInterval    time.Duration
	videoInterval      time.Duration
	cleanupInterval    time.Duration
	imageRetentionTime time.Duration
	outDir             string
	tempDir            string // Temporary directory for images being captured/optimized
	readyDir           string // Ready directory for images ready to be processed into video
	videoDir           string
	externalVideoDir   string // External SD card location for videos
	processedDir       string
}

func main() {
	// Environment variables with defaults
	cam := getEnvDefault("CAM", "0")
	captureIntervalStr := getEnvDefault("INTERVAL", "10")
	videoIntervalStr := getEnvDefault("VIDEO_INTERVAL", "5")      // minutes
	cleanupIntervalStr := getEnvDefault("CLEANUP_INTERVAL", "5")  // minutes
	retentionStr := getEnvDefault("IMAGE_RETENTION", "30")        // minutes

	captureInterval, err := strconv.Atoi(captureIntervalStr)
	if err != nil {
		log.Fatalf("Invalid INTERVAL value: %s", captureIntervalStr)
	}

	videoInterval, err := strconv.Atoi(videoIntervalStr)
	if err != nil {
		log.Fatalf("Invalid VIDEO_INTERVAL value: %s", videoIntervalStr)
	}

	cleanupInterval, err := strconv.Atoi(cleanupIntervalStr)
	if err != nil {
		log.Fatalf("Invalid CLEANUP_INTERVAL value: %s", cleanupIntervalStr)
	}

	retention, err := strconv.Atoi(retentionStr)
	if err != nil {
		log.Fatalf("Invalid IMAGE_RETENTION value: %s", retentionStr)
	}

	// Setup directories
	outDir := "/data/data/com.termux/files/home/timelapse"
	tempDir := filepath.Join(outDir, "raw", "temp")
	readyDir := filepath.Join(outDir, "raw", "ready")
	videoDir := filepath.Join(outDir, "videos")
	externalVideoDir := "/data/data/com.termux/files/home/storage/external-1/timelapse/videos"
	processedDir := filepath.Join(outDir, "processed")

	for _, dir := range []string{tempDir, readyDir, videoDir, processedDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Create external directory if possible (don't fail if SD card not available)
	if err := os.MkdirAll(externalVideoDir, 0755); err != nil {
		log.Printf("Warning: Failed to create external video directory %s: %v", externalVideoDir, err)
		log.Println("Videos will only be saved locally")
	}

	// Setup TMPDIR
	tmpDir := filepath.Join(os.Getenv("HOME"), "tmp")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		log.Fatalf("Failed to create tmp directory %s: %v", tmpDir, err)
	}
	os.Setenv("TMPDIR", tmpDir)

	config := &Config{
		cam:                cam,
		captureInterval:    time.Duration(captureInterval) * time.Second,
		videoInterval:      time.Duration(videoInterval) * time.Minute,
		cleanupInterval:    time.Duration(cleanupInterval) * time.Minute,
		imageRetentionTime: time.Duration(retention) * time.Minute,
		outDir:             outDir,
		tempDir:            tempDir,
		readyDir:           readyDir,
		videoDir:           videoDir,
		externalVideoDir:   externalVideoDir,
		processedDir:       processedDir,
	}

	// Keep CPU awake while running
	if err := exec.Command("termux-wake-lock").Run(); err != nil {
		log.Fatalf("Failed to acquire wake lock: %v", err)
	}

	// Setup signal handling for cleanup
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cleanup()
		os.Exit(0)
	}()

	fmt.Printf("Starting timelapse system:\n")
	fmt.Printf("  Camera: %s\n", config.cam)
	fmt.Printf("  Capture interval: %v\n", config.captureInterval)
	fmt.Printf("  Video generation: every %v\n", config.videoInterval)
	fmt.Printf("  Cleanup: every %v (retention: %v)\n", config.cleanupInterval, config.imageRetentionTime)
	fmt.Printf("  Directories:\n")
	fmt.Printf("    - temp: %s\n", config.tempDir)
	fmt.Printf("    - ready: %s\n", config.readyDir)
	fmt.Printf("    - videos (local): %s\n", config.videoDir)
	fmt.Printf("    - videos (external): %s\n", config.externalVideoDir)

	var wg sync.WaitGroup

	// Flow 1: Image capture loop
	wg.Add(1)
	go func() {
		defer wg.Done()
		captureLoop(config)
	}()

	// Flow 2: Video generation loop
	wg.Add(1)
	go func() {
		defer wg.Done()
		videoLoop(config)
	}()

	// Flow 3: Cleanup loop
	wg.Add(1)
	go func() {
		defer wg.Done()
		cleanupLoop(config)
	}()

	wg.Wait()
}

// Flow 1: Capture images at regular intervals
func captureLoop(config *Config) {
	log.Println("Starting image capture loop")
	ticker := time.NewTicker(config.captureInterval)
	defer ticker.Stop()

	// Take first photo immediately
	captureAndOptimize(config)

	// Then continue with ticker
	for range ticker.C {
		captureAndOptimize(config)
	}
}

func captureAndOptimize(config *Config) {
	timestamp := time.Now().Format("20060102-150405")
	filename := timestamp + ".jpg"
	torchOn()

	// Write to temp directory first
	tempPath := filepath.Join(config.tempDir, filename)
	if err := capturePhoto(config.cam, tempPath); err != nil {
		torchOff()
		log.Printf("Camera failed: %v", err)
		return
	}

	torchOff()

	// Optimize the image in temp directory
	if err := optimizeImage(tempPath); err != nil {
		log.Printf("Optimization failed: %v", err)
		// Still move to ready even if optimization fails
	}

	// Atomically move to ready directory
	readyPath := filepath.Join(config.readyDir, filename)
	if err := os.Rename(tempPath, readyPath); err != nil {
		log.Printf("Failed to move image to ready: %v", err)
		return
	}

	torchOff() // Extra safety to ensure torch is off
}

// Flow 2: Generate videos from captured images
func videoLoop(config *Config) {
	log.Println("Starting video generation loop")
	ticker := time.NewTicker(config.videoInterval)
	defer ticker.Stop()

	for range ticker.C {
		if err := generateVideo(config); err != nil {
			log.Printf("Video generation failed: %v", err)
		}
	}
}

// Flow 3: Cleanup old processed images
func cleanupLoop(config *Config) {
	log.Println("Starting cleanup loop")
	ticker := time.NewTicker(config.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		if err := cleanupOldImages(config); err != nil {
			log.Printf("Cleanup failed: %v", err)
		}
	}
}

func getEnvDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func torchOn() {
	if err := exec.Command("termux-torch", "on").Run(); err != nil {
		log.Printf("Warning: failed to turn torch on: %v", err)
	}
}

func torchOff() {
	if err := exec.Command("termux-torch", "off").Run(); err != nil {
		log.Printf("Warning: failed to turn torch off: %v", err)
	}
}

func capturePhoto(cam, outputPath string) error {
	cmd := exec.Command("termux-camera-photo", "-c", cam, outputPath)
	return cmd.Run()
}

func optimizeImage(photoPath string) error {
	// Run jpegtran
	cmd := exec.Command("jpegtran", "-copy", "none", "-optimize", "-progressive", "-outfile", "half.jpg", photoPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("jpegtran failed: %w", err)
	}

	// Run jpegoptim
	cmd = exec.Command("jpegoptim", "--size=250k", "--strip-all", "-q", "half.jpg")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("jpegoptim failed: %w", err)
	}

	// Move optimized image to original location
	if err := os.Rename("half.jpg", photoPath); err != nil {
		return fmt.Errorf("failed to move optimized image: %w", err)
	}

	return nil
}

func generateVideo(config *Config) error {
	// Get all images in ready directory
	files, err := filepath.Glob(filepath.Join(config.readyDir, "*.jpg"))
	if err != nil {
		return fmt.Errorf("failed to list images: %w", err)
	}

	if len(files) == 0 {
		log.Println("No images to process for video generation")
		return nil
	}

	// Sort files by name (which are timestamps)
	sort.Strings(files)

	log.Printf("Generating video from %d images", len(files))

	// Generate video filename with timestamp
	videoTimestamp := time.Now().Format("20060102-150405")
	videoPath := filepath.Join(config.videoDir, fmt.Sprintf("timelapse-%s.mp4", videoTimestamp))

	// Create ffmpeg command
	// Using pattern matching: ffmpeg expects files in sequential order
	// We'll use a concat demuxer for better control
	listFile := filepath.Join(config.outDir, "ffmpeg_input.txt")
	f, err := os.Create(listFile)
	if err != nil {
		return fmt.Errorf("failed to create input list: %w", err)
	}

	for _, file := range files {
		fmt.Fprintf(f, "file '%s'\n", file)
	}
	f.Close()
	defer os.Remove(listFile)

	// Run ffmpeg
	cmd := exec.Command("ffmpeg",
		"-f", "concat",
		"-safe", "0",
		"-i", listFile,
		"-vf", "fps=30",
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		"-preset", "medium",
		"-crf", "23",
		"-y",
		videoPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg failed: %w\nOutput: %s", err, string(output))
	}

	log.Printf("Video created: %s", videoPath)

	// Copy video to external SD card location
	externalPath := filepath.Join(config.externalVideoDir, filepath.Base(videoPath))
	if err := copyFile(videoPath, externalPath); err != nil {
		log.Printf("Warning: failed to copy video to external storage: %v", err)
		log.Printf("Video is still available locally at: %s", videoPath)
	} else {
		log.Printf("Video copied to external storage: %s", externalPath)
		// Remove local video after successful copy
		if err := os.Remove(videoPath); err != nil {
			log.Printf("Warning: failed to remove local video: %v", err)
		} else {
			log.Printf("Local video removed after successful copy to external storage")
		}
	}

	// Move processed images to processed directory
	for _, file := range files {
		basename := filepath.Base(file)
		destPath := filepath.Join(config.processedDir, basename)
		if err := os.Rename(file, destPath); err != nil {
			log.Printf("Warning: failed to move %s to processed: %v", basename, err)
		}
	}

	log.Printf("Moved %d images to processed directory", len(files))
	return nil
}

func cleanupOldImages(config *Config) error {
	// Get all images in processed directory
	files, err := filepath.Glob(filepath.Join(config.processedDir, "*.jpg"))
	if err != nil {
		return fmt.Errorf("failed to list processed images: %w", err)
	}

	if len(files) == 0 {
		log.Println("No processed images to clean up")
		return nil
	}

	now := time.Now()
	deletedCount := 0

	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			log.Printf("Warning: failed to stat %s: %v", file, err)
			continue
		}

		// Check if file is older than retention time
		age := now.Sub(info.ModTime())
		if age > config.imageRetentionTime {
			if err := os.Remove(file); err != nil {
				log.Printf("Warning: failed to delete %s: %v", file, err)
			} else {
				deletedCount++
			}
		}
	}

	if deletedCount > 0 {
		log.Printf("Deleted %d old processed images (retention: %v)", deletedCount, config.imageRetentionTime)
	}

	return nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	if _, err := destFile.ReadFrom(sourceFile); err != nil {
		return fmt.Errorf("failed to copy data: %w", err)
	}

	// Sync to ensure data is written to disk
	if err := destFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync file: %w", err)
	}

	return nil
}

func cleanup() {
	fmt.Println("\nCleaning up...")
	exec.Command("termux-wake-unlock").Run()
	torchOff()
}

