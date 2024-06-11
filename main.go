// package main

// import (
// 	"compress/gzip"
// 	"encoding/csv"
// 	"fmt"
// 	"io"
// 	"os"
// 	"strings"
// 	"time"

// 	"github.com/klauspost/compress/zstd"
// )

// func main() {
// 	// Replace with the path to your input file
// 	inputFile := "input.txt"

// 	// Step 1: Compress using gzip
// 	start := time.Now()
// 	gzipFile := "gzip_compressed.gz"
// 	gzipTime, gzipSize, err := compressWithGzip(inputFile, gzipFile)
// 	if err != nil {
// 		fmt.Printf("Error compressing with gzip: %v\n", err)
// 		return
// 	}
// 	gzipDuration := time.Since(start)

// 	// Step 2: Compress using zstd with various levels
// 	zstdTimes := make(map[int]time.Duration)
// 	zstdSizes := make(map[int]int64)

// 	zstdLevels := []int{-1, 1, 3, 5, 10, 15, 20}
// 	for _, level := range zstdLevels {
// 		start := time.Now()
// 		zstdFile := fmt.Sprintf("zstd_compressed_level_%d.zst", level)
// 		zstdTime, zstdSize, err := compressWithZstd(inputFile, zstdFile, level)
// 		if err != nil {
// 			fmt.Printf("Error compressing with zstd (level %d): %v\n", level, err)
// 			return
// 		}
// 		zstdTimes[level] = time.Since(start)
// 		zstdSizes[level] = zstdSize
// 	}

// 	// Step 3: Decompress gzip file
// 	start = time.Now()
// 	decompressedFile := "decompressed_gzip.txt"
// 	gunzipTime, err := decompressGzip(gzipFile, decompressedFile)
// 	if err != nil {
// 		fmt.Printf("Error decompressing gzip: %v\n", err)
// 		return
// 	}
// 	gunzipDuration := time.Since(start)

// 	// Step 4: Decompress zstd files (various levels)
// 	zstdDecompressTimes := make(map[int]time.Duration)
// 	for _, level := range zstdLevels {
// 		start := time.Now()
// 		decompressedFile := fmt.Sprintf("decompressed_zstd_level_%d.txt", level)
// 		zstdDecompressTime, err := decompressZstd(fmt.Sprintf("zstd_compressed_level_%d.zst", level), decompressedFile)
// 		if err != nil {
// 			fmt.Printf("Error decompressing zstd (level %d): %v\n", level, err)
// 			return
// 		}
// 		zstdDecompressTimes[level] = time.Since(start)
// 	}

// 	// Record original file size
// 	originalSize, err := getFileSize(inputFile)
// 	if err != nil {
// 		fmt.Printf("Error getting original file size: %v\n", err)
// 		return
// 	}

// 	// Print results
// 	fmt.Printf("Original file size: %d bytes\n", originalSize)
// 	fmt.Printf("Step 1: Gzip compression took %v, compressed size: %d bytes\n", gzipDuration, gzipSize)
// 	for _, level := range zstdLevels {
// 		fmt.Printf("Step 2: Zstd compression (level %d) took %v, compressed size: %d bytes\n", level, zstdTimes[level], zstdSizes[level])
// 	}
// 	fmt.Printf("Step 3: Gzip decompression took %v\n", gunzipDuration)
// 	for _, level := range zstdLevels {
// 		fmt.Printf("Step 4: Zstd decompression (level %d) took %v\n", level, zstdDecompressTimes[level])
// 	}
// }

// // CompressionResult holds the results of the compression
// type CompressionResult struct {
// 	FilePath          string
// 	Algorithm         string
// 	OriginalSize      int64
// 	CompressedSize    int64
// 	CompressionTime   time.Duration
// 	DecompressionTime time.Duration
// }

// type FinalCompressionResult struct {
// 	FilePath              string
// 	OriginalSize          int64
// 	GzipCompressedSize    int64
// 	ZstdCompressedSize    int64
// 	GzipCompressionTime   time.Duration
// 	GzipDecompressionTime time.Duration
// 	ZstdCompressionTime   time.Duration
// 	ZstdDecompressionTime time.Duration
// }

// func writeResultsToCSV(filename string, results []FinalCompressionResult) error {
// 	var writeHeader bool

// 	// Check if the file exists and if it's empty
// 	if _, err := os.Stat(filename); err == nil {
// 		// File exists, check if it's empty
// 		fi, err := os.Stat(filename)
// 		if err != nil {
// 			return err
// 		}
// 		writeHeader = fi.Size() == 0
// 	} else if os.IsNotExist(err) {
// 		// File doesn't exist, we need to write the header
// 		writeHeader = true
// 	} else {
// 		return err
// 	}

// 	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	writer := csv.NewWriter(file)
// 	defer writer.Flush()

// 	if writeHeader {
// 		// Write the header
// 		header := []string{"FilePath", "Original Size", "Gzip Compressed Size(B)", "Gzip Compression Time(s)", "Gzip Decompression Time(s)", "Zstd Compressed Size(B)", "Zstd Compression Time(s)", "Zstd Decompression Time(s)"}
// 		err = writer.Write(header)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	// Write the data
// 	for _, result := range results {
// 		record := []string{
// 			result.FilePath,
// 			fmt.Sprintf("%d", result.OriginalSize),
// 			fmt.Sprintf("%d", result.GzipCompressedSize),
// 			fmt.Sprintf("%f", result.GzipCompressionTime.Seconds()),
// 			fmt.Sprintf("%f", result.GzipDecompressionTime.Seconds()),
// 			fmt.Sprintf("%d", result.ZstdCompressedSize),
// 			fmt.Sprintf("%f", result.ZstdCompressionTime.Seconds()),
// 			fmt.Sprintf("%f", result.ZstdDecompressionTime.Seconds()),
// 		}
// 		err = writer.Write(record)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// // compressGzip compresses a file using gzip and returns the result
// func compressGzip(inputFileName string) (CompressionResult, error) {
// 	start := time.Now()

// 	// Open the input file
// 	inputFile, err := os.Open(inputFileName)
// 	if err != nil {
// 		return CompressionResult{}, fmt.Errorf("error opening input file: %v", err)
// 	}
// 	defer inputFile.Close()

// 	// Get the original file size
// 	inputFileInfo, err := inputFile.Stat()
// 	if err != nil {
// 		return CompressionResult{}, fmt.Errorf("error getting input file info: %v", err)
// 	}
// 	originalSize := inputFileInfo.Size()

// 	// Create the output file
// 	outputFileName := inputFileName + ".gz"
// 	outputFile, err := os.Create(outputFileName)
// 	if err != nil {
// 		return CompressionResult{}, fmt.Errorf("error creating output file: %v", err)
// 	}
// 	defer outputFile.Close()

// 	// Create a gzip writer
// 	gzipWriter := gzip.NewWriter(outputFile)
// 	defer gzipWriter.Close()

// 	// Compress the data
// 	_, err = io.Copy(gzipWriter, inputFile)
// 	if err != nil {
// 		return CompressionResult{}, fmt.Errorf("error compressing file: %v", err)
// 	}

// 	// Ensure all data is flushed
// 	err = gzipWriter.Close()
// 	if err != nil {
// 		return CompressionResult{}, fmt.Errorf("error closing gzip writer: %v", err)
// 	}

// 	// Get the compressed file size
// 	compressedFileInfo, err := os.Stat(outputFileName)
// 	if err != nil {
// 		return CompressionResult{}, fmt.Errorf("error getting compressed file size: %v", err)
// 	}
// 	compressedSize := compressedFileInfo.Size()

// 	compressionTime := time.Since(start)

// 	// Decompress the file to measure decompression time
// 	start = time.Now()
// 	err = decompressGzip(outputFileName)
// 	if err != nil {
// 		return CompressionResult{}, fmt.Errorf("error decompressing file: %v", err)
// 	}
// 	decompressionTime := time.Since(start)

// 	return CompressionResult{
// 		FilePath:          outputFileName,
// 		Algorithm:         "gzip",
// 		OriginalSize:      originalSize,
// 		CompressedSize:    compressedSize,
// 		CompressionTime:   compressionTime,
// 		DecompressionTime: decompressionTime,
// 	}, nil
// }

// // decompressGzip decompresses a gzip file
// func decompressGzip(compressedFileName string) error {
// 	compressedFile, err := os.Open(compressedFileName)
// 	if err != nil {
// 		return fmt.Errorf("error opening compressed file: %v", err)
// 	}
// 	defer compressedFile.Close()

// 	gzipReader, err := gzip.NewReader(compressedFile)
// 	if err != nil {
// 		return fmt.Errorf("error creating gzip reader: %v", err)
// 	}
// 	defer gzipReader.Close()

// 	outputFileName := strings.TrimSuffix(compressedFileName, ".gz")
// 	outputFile, err := os.Create(outputFileName)
// 	if err != nil {
// 		return fmt.Errorf("error creating output file: %v", err)
// 	}
// 	defer outputFile.Close()

// 	_, err = io.Copy(outputFile, gzipReader)
// 	if err != nil {
// 		return fmt.Errorf("error decompressing file: %v", err)
// 	}

// 	return nil
// }

// // compressZstd compresses a file using zstd and returns the result
// func compressZstd(inputFileName string) (CompressionResult, error) {
// 	start := time.Now()

// 	// Open the input file
// 	inputFile, err := os.Open(inputFileName)
// 	if err != nil {
// 		return CompressionResult{}, fmt.Errorf("error opening input file: %v", err)
// 	}
// 	defer inputFile.Close()

// 	// Get the original file size
// 	inputFileInfo, err := inputFile.Stat()
// 	if err != nil {
// 		return CompressionResult{}, fmt.Errorf("error getting input file info: %v", err)
// 	}
// 	originalSize := inputFileInfo.Size()

// 	// Create the output file
// 	outputFileName := inputFileName + ".zst"
// 	outputFile, err := os.Create(outputFileName)
// 	if err != nil {
// 		return CompressionResult{}, fmt.Errorf("error creating output file: %v", err)
// 	}
// 	defer outputFile.Close()

// 	// Create a zstd writer
// 	zstdWriter, err := zstd.NewWriter(outputFile, zstd.WithEncoderLevel(zstd.SpeedBestCompression))
// 	if err != nil {
// 		return CompressionResult{}, fmt.Errorf("error creating zstd writer: %v", err)
// 	}
// 	defer zstdWriter.Close()

// 	// Compress the data
// 	_, err = io.Copy(zstdWriter, inputFile)
// 	if err != nil {
// 		return CompressionResult{}, fmt.Errorf("error compressing file: %v", err)
// 	}

// 	// Ensure all data is flushed
// 	err = zstdWriter.Close()
// 	if err != nil {
// 		return CompressionResult{}, fmt.Errorf("error closing zstd writer: %v", err)
// 	}

// 	// Get the compressed file size
// 	compressedFileInfo, err := os.Stat(outputFileName)
// 	if err != nil {
// 		return CompressionResult{}, fmt.Errorf("error getting compressed file size: %v", err)
// 	}
// 	compressedSize := compressedFileInfo.Size()

// 	compressionTime := time.Since(start)

// 	// Decompress the file to measure decompression time
// 	start = time.Now()
// 	err = decompressZstd(outputFileName)
// 	if err != nil {
// 		return CompressionResult{}, fmt.Errorf("error decompressing file: %v", err)
// 	}
// 	decompressionTime := time.Since(start)

// 	return CompressionResult{
// 		FilePath:          outputFileName,
// 		Algorithm:         "zstd",
// 		OriginalSize:      originalSize,
// 		CompressedSize:    compressedSize,
// 		CompressionTime:   compressionTime,
// 		DecompressionTime: decompressionTime,
// 	}, nil
// }

// // decompressZstd decompresses a zstd file
// func decompressZstd(compressedFileName string) error {
// 	compressedFile, err := os.Open(compressedFileName)
// 	if err != nil {
// 		return fmt.Errorf("error opening compressed file: %v", err)
// 	}
// 	defer compressedFile.Close()

// 	zstdReader, err := zstd.NewReader(compressedFile)
// 	if err != nil {
// 		return fmt.Errorf("error creating zstd reader: %v", err)
// 	}
// 	defer zstdReader.Close()

// 	outputFileName := compressedFileName + ".decompressed"
// 	outputFile, err := os.Create(outputFileName)
// 	if err != nil {
// 		return fmt.Errorf("error creating output file: %v", err)
// 	}
// 	defer outputFile.Close()

// 	_, err = io.Copy(outputFile, zstdReader)
// 	if err != nil {
// 		return fmt.Errorf("error decompressing file: %v", err)
// 	}

// 	return nil
// }

// // get the size of file
// func getFileSize(filename string) (int64, error) {
// 	fileInfo, err := os.Stat(filename)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return fileInfo.Size(), nil
// }
