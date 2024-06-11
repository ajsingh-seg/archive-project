package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/klauspost/compress/zstd"
)

func main() {
	// Replace with the path to your input file
	inputFile := "input.txt"

	// Step 1: Compress using gzip
	start := time.Now()
	gzipFile := "gzip_compressed.gz"
	gzipTime, gzipSize, err := compressWithGzip(inputFile, gzipFile)
	if err != nil {
		fmt.Printf("Error compressing with gzip: %v\n", err)
		return
	}
	gzipDuration := time.Since(start)

	// Step 2: Compress using zstd with various levels
	zstdTimes := make(map[int]time.Duration)
	zstdSizes := make(map[int]int64)

	zstdLevels := []int{-1, 1, 3, 5, 10, 15, 20}
	for _, level := range zstdLevels {
		start := time.Now()
		zstdFile := fmt.Sprintf("zstd_compressed_level_%d.zst", level)
		zstdTime, zstdSize, err := compressWithZstd(inputFile, zstdFile, level)
		if err != nil {
			fmt.Printf("Error compressing with zstd (level %d): %v\n", level, err)
			return
		}
		zstdTimes[level] = time.Since(start)
		zstdSizes[level] = zstdSize
	}

	// Step 3: Decompress gzip file
	start = time.Now()
	decompressedFile := "decompressed_gzip.txt"
	gunzipTime, err := decompressGzip(gzipFile, decompressedFile)
	if err != nil {
		fmt.Printf("Error decompressing gzip: %v\n", err)
		return
	}
	gunzipDuration := time.Since(start)

	// Step 4: Decompress zstd files (various levels)
	zstdDecompressTimes := make(map[int]time.Duration)
	for _, level := range zstdLevels {
		start := time.Now()
		decompressedFile := fmt.Sprintf("decompressed_zstd_level_%d.txt", level)
		zstdDecompressTime, err := decompressZstd(fmt.Sprintf("zstd_compressed_level_%d.zst", level), decompressedFile)
		if err != nil {
			fmt.Printf("Error decompressing zstd (level %d): %v\n", level, err)
			return
		}
		zstdDecompressTimes[level] = time.Since(start)
	}

	// Record original file size
	originalSize, err := getFileSize(inputFile)
	if err != nil {
		fmt.Printf("Error getting original file size: %v\n", err)
		return
	}

	// Print results
	fmt.Printf("Original file size: %d bytes\n", originalSize)
	fmt.Printf("Step 1: Gzip compression took %v, compressed size: %d bytes\n", gzipDuration, gzipSize)
	for _, level := range zstdLevels {
		fmt.Printf("Step 2: Zstd compression (level %d) took %v, compressed size: %d bytes\n", level, zstdTimes[level], zstdSizes[level])
	}
	fmt.Printf("Step 3: Gzip decompression took %v\n", gunzipDuration)
	for _, level := range zstdLevels {
		fmt.Printf("Step 4: Zstd decompression (level %d) took %v\n", level, zstdDecompressTimes[level])
	}
}

func compressWithGzip(inputFile, outputFile string) (time.Duration, int64, error) {
	start := time.Now()

	// Open input file
	inFile, err := os.Open(inputFile)
	if err != nil {
		return 0, 0, err
	}
	defer inFile.Close()

	// Create output file
	outFile, err := os.Create(outputFile)
	if err != nil {
		return 0, 0, err
	}
	defer outFile.Close()

	// Create gzip writer
	gzipWriter := gzip.NewWriter(outFile)
	defer gzipWriter.Close()

	// Copy content from input to gzip writer
	_, err = io.Copy(gzipWriter, inFile)
	if err != nil {
		return 0, 0, err
	}

	// Calculate compressed file size
	fileInfo, err := outFile.Stat()
	if err != nil {
		return 0, 0, err
	}
	compressedSize := fileInfo.Size()

	return time.Since(start), compressedSize, nil
}

func compressWithZstd(inputFile, outputFile string, level int) (time.Duration, int64, error) {
	start := time.Now()

	// Open input file
	inFile, err := os.Open(inputFile)
	if err != nil {
		return 0, 0, err
	}
	defer inFile.Close()

	// Create output file
	outFile, err := os.Create(outputFile)
	if err != nil {
		return 0, 0, err
	}
	defer outFile.Close()

	// Create zstd writer with specified level
	zstdWriter, err := zstd.NewWriter(outFile, zstd.WithEncoderLevel(zstd.EncoderLevel(level)))
	if err != nil {
		return 0, 0, err
	}
	defer zstdWriter.Close()

	// Copy content from input to zstd writer
	_, err = io.Copy(zstdWriter, inFile)
	if err != nil {
		return 0, 0, err
	}

	// Calculate compressed file size
	fileInfo, err := outFile.Stat()
	if err != nil {
		return 0, 0, err
	}
	compressedSize := fileInfo.Size()

	return time.Since(start), compressedSize, nil
}

func decompressGzip(inputFile, outputFile string) (time.Duration, error) {
	start := time.Now()

	// Open input gzip file
	gzipFile, err := os.Open(inputFile)
	if err != nil {
		return 0, err
	}
	defer gzipFile.Close()

	// Create output file
	outFile, err := os.Create(outputFile)
	if err != nil {
		return 0, err
	}
	defer outFile.Close()

	// Create gzip reader
	gzipReader, err := gzip.NewReader(gzipFile)
	if err != nil {
		return 0, err
	}
	defer gzipReader.Close()

	// Copy content from gzip reader to output file
	_, err = io.Copy(outFile, gzipReader)
	if err != nil {
		return 0, err
	}

	return time.Since(start), nil
}

func decompressZstd(inputFile, outputFile string) (time.Duration, error) {
	start := time.Now()

	// Open input zstd file
	zstdFile, err := os.Open(inputFile)
	if err != nil {
		return 0, err
	}
	defer zstdFile.Close()

	// Create output file
	outFile, err := os.Create(outputFile)
	if err != nil {
		return 0, err
	}
	defer outFile.Close()

	// Create zstd reader
	zstdReader, err := zstd.NewReader(zstdFile)
	if err != nil {
		return 0, err
	}
	defer zstdReader.Close()

	// Copy content from zstd reader to output file
	_, err = io.Copy(outFile, zstdReader)
	if err != nil {
		return 0, err
	}

	return time.Since(start), nil
}

func getFileSize(filename string) (int64, error) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}
