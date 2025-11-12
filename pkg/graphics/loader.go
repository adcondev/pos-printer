package graphics

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/jpeg" // Register JPEG decoder
	_ "image/png"  // Register PNG decoder
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ImgFromFile loads an image from a file path within baseDir.
//
// It validates the path to prevent directory traversal attacks.
func ImgFromFile(baseDir, relPath string) (image.Image, error) {
	baseAbs, err := filepath.Abs(baseDir)
	if err != nil {
		return nil, fmt.Errorf("invalid base directory: %w", err)
	}

	target := filepath.Join(baseAbs, relPath)
	targetAbs, err := filepath.Abs(target)
	if err != nil {
		return nil, fmt.Errorf("invalid target path: %w", err)
	}

	// Resolve symlinks to prevent escaping
	if eval, err := filepath.EvalSymlinks(baseAbs); err == nil {
		baseAbs = eval
	}
	if eval, err := filepath.EvalSymlinks(targetAbs); err == nil {
		targetAbs = eval
	}

	// Ensure target is inside base directory
	if targetAbs != baseAbs && !strings.HasPrefix(targetAbs, baseAbs+string(os.PathSeparator)) {
		return nil, fmt.Errorf("file %s is outside allowed directory", relPath)
	}

	// Open and decode the image file
	file, err := os.Open(targetAbs)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("failed to close file: %v\n", err)
		}
	}(file)

	// Decode image
	img, format, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	log.Printf("Loaded image format: %s\n", format)
	return img, nil
}

// ImgFromBase64 converts a base64-encoded string to an image.Image.
func ImgFromBase64(data string) (image.Image, error) {
	// Decode base64 to byte slice
	imgBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 string: %w", err)
	}

	// Wrap byte slice in an io.Reader for image.Decode
	reader := bytes.NewReader(imgBytes)

	// Convert to image.Image
	img, format, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	log.Printf("Decoded image format: %s\n", format)
	return img, nil
}
