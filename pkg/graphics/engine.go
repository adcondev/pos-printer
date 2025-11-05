package graphics

import (
	"fmt"
	"image"
	"image/color"

	"golang.org/x/image/draw"
)

// ProcessingMode defines how images are converted to monochrome
type ProcessingMode int

const (
	// Threshold applies simple threshold conversion
	Threshold ProcessingMode = iota
	// Atkinson applies Atkinson dithering algorithm
	Atkinson
	// FloydSteinberg applies Floyd-Steinberg dithering (future)
	FloydSteinberg
	// Ordered applies ordered dithering with Bayer matrix (future)
	Ordered
)

// Options configures the graphics processing pipeline
type Options struct {
	Width          int            // Target width in pixels
	Threshold      uint8          // Threshold for black/white (0-255)
	Mode           ProcessingMode // Processing algorithm
	AutoRotate     bool           // Auto-rotate for best fit
	PreserveAspect bool           // Maintain aspect ratio
}

// DefaultOptions returns sensible defaults for 80mm printers
func DefaultOptions() *Options {
	return &Options{
		Width:          512,
		Threshold:      128,
		Mode:           Threshold,
		AutoRotate:     false,
		PreserveAspect: true,
	}
}

// Pipeline represents the image processing pipeline
type Pipeline struct {
	opts *Options
}

// NewPipeline creates a new processing pipeline with given options
func NewPipeline(opts *Options) *Pipeline {
	if opts == nil {
		opts = DefaultOptions()
	}
	return &Pipeline{opts: opts}
}

// Process transforms an image through the complete pipeline
func (p *Pipeline) Process(img image.Image) (*MonochromeBitmap, error) {
	if img == nil {
		return nil, fmt.Errorf("input image cannot be nil")
	}

	// Step 1: Resize if needed
	if p.opts.Width > 0 && img.Bounds().Dx() != p.opts.Width {
		img = p.resize(img)
	}

	// Step 2: Convert to grayscale
	gray := p.toGrayscale(img)

	// Step 3: Apply processing mode
	var mono *MonochromeBitmap
	switch p.opts.Mode {
	case Atkinson:
		mono = p.applyAtkinson(gray)
	case Threshold:
		fallthrough
	default:
		mono = p.applyThreshold(gray)
	}

	return mono, nil
}

// resize scales the image to target width maintaining aspect ratio
func (p *Pipeline) resize(img image.Image) image.Image {
	bounds := img.Bounds()
	srcW, srcH := bounds.Dx(), bounds.Dy()

	targetW := p.opts.Width
	targetH := srcH

	if p.opts.PreserveAspect {
		targetH = (srcH * targetW) / srcW
	}

	// Simple nearest-neighbor scaling for now
	// TODO: Implement better scaling algorithms
	dst := image.NewRGBA(image.Rect(0, 0, targetW, targetH))

	// Bilinear scaling
	draw.BiLinear.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)

	return dst
}

// toGrayscale converts any image to grayscale
func (p *Pipeline) toGrayscale(img image.Image) *image.Gray {
	bounds := img.Bounds()
	gray := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			gray.Set(x, y, c)
		}
	}

	return gray
}

// applyThreshold applies simple threshold conversion
func (p *Pipeline) applyThreshold(gray *image.Gray) *MonochromeBitmap {
	bounds := gray.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	mono := NewMonochromeBitmap(width, height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := gray.GrayAt(x, y).Y
			// Set pixel to black (true) if below threshold
			if pixel < p.opts.Threshold {
				mono.SetPixel(x, y, true)
			}
		}
	}

	return mono
}

// TODO: Implement Floyd-Steinberg and Ordered dithering methods in future

// applyAtkinson implements Atkinson dithering algorithm
func (p *Pipeline) applyAtkinson(gray *image.Gray) *MonochromeBitmap {
	bounds := gray.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	mono := NewMonochromeBitmap(width, height)

	// Create a working copy for error diffusion
	work := make([][]int, height)
	for y := 0; y < height; y++ {
		work[y] = make([]int, width)
		for x := 0; x < width; x++ {
			work[y][x] = int(gray.GrayAt(x, y).Y)
		}
	}

	// Atkinson dithering pattern:
	//     *  1  1
	//  1  1  1
	//     1
	// Error is distributed as 1/8 to each neighbor (total 6/8 = 3/4)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			oldPixel := work[y][x]
			newPixel := 0
			if oldPixel > int(p.opts.Threshold) {
				newPixel = 255
			}

			// Set the monochrome pixel
			mono.SetPixel(x, y, newPixel == 0)

			// Calculate error
			err := oldPixel - newPixel

			// Atkinson only diffuses 3/4 (6/8) of the error
			// Each of the 6 neighbors gets 1/8 of the original error
			diffusedError := err / 8

			// Distribute to neighbors
			if x+1 < width {
				work[y][x+1] += diffusedError
			}
			if x+2 < width {
				work[y][x+2] += diffusedError
			}
			if y+1 < height {
				if x-1 >= 0 {
					work[y+1][x-1] += diffusedError
				}
				work[y+1][x] += diffusedError
				if x+1 < width {
					work[y+1][x+1] += diffusedError
				}
			}
			if y+2 < height {
				work[y+2][x] += diffusedError
			}
		}
	}

	return mono
}
