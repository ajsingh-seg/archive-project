package main

import (
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/klauspost/compress/zstd"
)

func main() {
	fmt.Println("Starting the compression engine")

	// Input file
	inputFile := "random_json_files/1-10MB.txt"

	// Compressing the file using gzip
	start := time.Now()
	gzipResult, err := compressFileGzip(inputFile)
	if err != nil {
		fmt.Printf("Error compressing the file with gzip: %v\n", err)
		return
	}
	gzipResult.CompressionTime = time.Since(start)
	fmt.Printf("Gzip compression result: %+v\n", gzipResult)

	// Compressing the file using zstd
	start = time.Now()
	zstdResult, err := compressFileZstd(inputFile, 50) // Example using compression level 3
	if err != nil {
		fmt.Printf("Error compressing the file with zstd: %v\n", err)
		return
	}
	zstdResult.CompressionTime = time.Since(start)
	fmt.Printf("Zstd compression result: %+v\n", zstdResult)

	// Getting the size of the original file
	originalSize, err := getFileSize(inputFile)
	if err != nil {
		fmt.Printf("Error getting the size of the original file: %v\n", err)
		return
	}

	// Creating the final result
	finalResult := FinalCompressionResult{
		FilePath:              inputFile,
		OriginalSize:          originalSize,
		GzipCompressedSize:    gzipResult.CompressedSize,
		GzipCompressionTime:   gzipResult.CompressionTime,
		GzipDecompressionTime: gzipResult.DecompressionTime,
		ZstdCompressedSize:    zstdResult.CompressedSize,
		ZstdCompressionTime:   zstdResult.CompressionTime,
		ZstdDecompressionTime: zstdResult.DecompressionTime,
	}

	// Writing the results to CSV
	results := []FinalCompressionResult{finalResult}
	err = writeResultsToCSV("compression_results.csv", results)
	if err != nil {
		fmt.Printf("Error writing results to CSV: %v\n", err)
		return
	}
	fmt.Println("Results written to CSV successfully")
}

type FinalCompressionResult struct {
	FilePath              string
	OriginalSize          int64
	GzipCompressedSize    int64
	ZstdCompressedSize    int64
	GzipCompressionTime   time.Duration
	GzipDecompressionTime time.Duration
	ZstdCompressionTime   time.Duration
	ZstdDecompressionTime time.Duration
}

func writeResultsToCSV(filename string, results []FinalCompressionResult) error {
	var writeHeader bool

	// Check if the file exists and if it's empty
	if _, err := os.Stat(filename); err == nil {
		// File exists, check if it's empty
		fi, err := os.Stat(filename)
		if err != nil {
			return err
		}
		writeHeader = fi.Size() == 0
	} else if os.IsNotExist(err) {
		// File doesn't exist, we need to write the header
		writeHeader = true
	} else {
		return err
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if writeHeader {
		// Write the header
		header := []string{"FilePath", "Original Size", "Gzip Compressed Size(B)", "Gzip Compression Time(s)", "Gzip Decompression Time(s)", "Zstd Compressed Size(B)", "Zstd Compression Time(s)", "Zstd Decompression Time(s)"}
		err = writer.Write(header)
		if err != nil {
			return err
		}
	}

	// Write the data
	for _, result := range results {
		record := []string{
			result.FilePath,
			fmt.Sprintf("%d", result.OriginalSize),
			fmt.Sprintf("%d", result.GzipCompressedSize),
			fmt.Sprintf("%f", result.GzipCompressionTime.Seconds()),
			fmt.Sprintf("%f", result.GzipDecompressionTime.Seconds()),
			fmt.Sprintf("%d", result.ZstdCompressedSize),
			fmt.Sprintf("%f", result.ZstdCompressionTime.Seconds()),
			fmt.Sprintf("%f", result.ZstdDecompressionTime.Seconds()),
		}
		err = writer.Write(record)
		if err != nil {
			return err
		}
	}
	log.Printf("Results written to CSV successfully: %v", filename)
	return nil
}

func getFileSize(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

type CompressionResult struct {
	FilePath          string
	Algorithm         string
	OriginalSize      int64
	CompressedSize    int64
	CompressionTime   time.Duration
	DecompressionTime time.Duration
}

func compressFileGzip(inputFileName string) (CompressionResult, error) {
	start := time.Now()

	// Open the input file
	inputFile, err := os.Open(inputFileName)
	if err != nil {
		return CompressionResult{}, fmt.Errorf("error opening input file: %v", err)
	}
	defer inputFile.Close()

	// Get the original file size
	inputFileInfo, err := inputFile.Stat()
	if err != nil {
		return CompressionResult{}, fmt.Errorf("error getting input file info: %v", err)
	}
	originalSize := inputFileInfo.Size()

	// Create the output file
	outputFileName := inputFileName + ".gz"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		return CompressionResult{}, fmt.Errorf("error creating output file: %v", err)
	}
	defer outputFile.Close()

	// Create a gzip writer
	gzipWriter := gzip.NewWriter(outputFile)
	defer gzipWriter.Close()

	// Compress the data
	_, err = io.Copy(gzipWriter, inputFile)
	if err != nil {
		return CompressionResult{}, fmt.Errorf("error compressing file: %v", err)
	}

	// Ensure all data is flushed
	err = gzipWriter.Close()
	if err != nil {
		return CompressionResult{}, fmt.Errorf("error closing gzip writer: %v", err)
	}

	// Get the compressed file size
	compressedFileInfo, err := os.Stat(outputFileName)
	if err != nil {
		return CompressionResult{}, fmt.Errorf("error getting compressed file size: %v", err)
	}
	compressedSize := compressedFileInfo.Size()

	compressionTime := time.Since(start)

	// Decompress the file to measure decompression time
	start = time.Now()
	err = decompressFileGzip(outputFileName)
	if err != nil {
		return CompressionResult{}, fmt.Errorf("error decompressing file: %v", err)
	}
	decompressionTime := time.Since(start)

	return CompressionResult{
		FilePath:          outputFileName,
		Algorithm:         "gzip",
		OriginalSize:      originalSize,
		CompressedSize:    compressedSize,
		CompressionTime:   compressionTime,
		DecompressionTime: decompressionTime,
	}, nil
}

func decompressFileGzip(compressedFileName string) error {
	compressedFile, err := os.Open(compressedFileName)
	if err != nil {
		return fmt.Errorf("error opening compressed file: %v", err)
	}
	defer compressedFile.Close()

	gzipReader, err := gzip.NewReader(compressedFile)
	if err != nil {
		return fmt.Errorf("error creating gzip reader: %v", err)
	}
	defer gzipReader.Close()

	outputFileName := compressedFileName + ".decompressed"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, gzipReader)
	if err != nil {
		return fmt.Errorf("error decompressing file: %v", err)
	}

	return nil
}

func compressFileZstd(inputFileName string, level int) (CompressionResult, error) {
	start := time.Now()

	// Open the input file
	inputFile, err := os.Open(inputFileName)
	if err != nil {
		return CompressionResult{}, fmt.Errorf("error opening input file: %v", err)
	}
	defer inputFile.Close()

	// Get the original file size
	inputFileInfo, err := inputFile.Stat()
	if err != nil {
		return CompressionResult{}, fmt.Errorf("error getting input file info: %v", err)
	}
	originalSize := inputFileInfo.Size()

	// Create the output file
	outputFileName := inputFileName + ".zst"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		return CompressionResult{}, fmt.Errorf("error creating output file: %v", err)
	}
	defer outputFile.Close()

	// Create a zstd writer with the specified compression level
	encoderLevel := zstd.EncoderLevelFromZstd(level)
	zstdWriter, err := zstd.NewWriter(outputFile, zstd.WithEncoderLevel(encoderLevel))
	if err != nil {
		return CompressionResult{}, fmt.Errorf("error creating zstd writer: %v", err)
	}
	defer zstdWriter.Close()

	// Compress the data
	_, err = io.Copy(zstdWriter, inputFile)
	if err != nil {
		return CompressionResult{}, fmt.Errorf("error compressing file: %v", err)
	}

	// Ensure all data is flushed
	err = zstdWriter.Close()
	if err != nil {
		return CompressionResult{}, fmt.Errorf("error closing zstd writer: %v", err)
	}

	// Get the compressed file size
	compressedFileInfo, err := os.Stat(outputFileName)
	if err != nil {
		return CompressionResult{}, fmt.Errorf("error getting compressed file size: %v", err)
	}
	compressedSize := compressedFileInfo.Size()

	compressionTime := time.Since(start)

	// Decompress the file to measure decompression time
	start = time.Now()
	err = decompressFileZstd(outputFileName)
	if err != nil {
		return CompressionResult{}, fmt.Errorf("error decompressing file: %v", err)
	}
	decompressionTime := time.Since(start)

	return CompressionResult{
		FilePath:          outputFileName,
		Algorithm:         "zstd",
		OriginalSize:      originalSize,
		CompressedSize:    compressedSize,
		CompressionTime:   compressionTime,
		DecompressionTime: decompressionTime,
	}, nil
}

func decompressFileZstd(compressedFileName string) error {
	compressedFile, err := os.Open(compressedFileName)
	if err != nil {
		return fmt.Errorf("error opening compressed file: %v", err)
	}
	defer compressedFile.Close()

	zstdReader, err := zstd.NewReader(compressedFile)
	if err != nil {
		return fmt.Errorf("error creating zstd reader: %v", err)
	}
	defer zstdReader.Close()

	outputFileName := compressedFileName + ".decompressed"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, zstdReader)
	if err != nil {
		return fmt.Errorf("error decompressing file: %v", err)
	}

	return nil
}
