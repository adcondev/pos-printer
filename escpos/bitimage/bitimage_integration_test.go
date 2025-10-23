package bitimage_test

import (
	"bytes"
	"testing"

	"github.com/adcondev/pos-printer/escpos/bitimage"
	"github.com/adcondev/pos-printer/escpos/shared"
	"github.com/adcondev/pos-printer/utils/test"
)

func TestIntegration_BitImage_LogoWorkflow(t *testing.T) {
	cmd := bitimage.NewCommands()

	t.Run("complete logo printing workflow", func(t *testing.T) {
		var buffer []byte

		// Step 1: Set graphics dot density for high quality
		densityCmd, err := cmd.Graphics.SetGraphicsDotDensity(
			bitimage.FunctionCodeDensity1,
			bitimage.Density360x360,
			bitimage.Density360x360,
		)
		if err != nil {
			t.Fatalf("SetGraphicsDotDensity failed: %v", err)
		}
		buffer = append(buffer, densityCmd...)

		// Step 2: Store raster graphics in buffer (small logo)
		logoWidth := uint16(200)
		logoHeight := uint16(50)
		widthBytes := (int(logoWidth) + 7) / 8
		logoData := test.RepeatByte(widthBytes*int(logoHeight), 0xAA)

		storeCmd, err := cmd.Graphics.StoreRasterGraphicsInBuffer(
			bitimage.Monochrome,
			bitimage.NormalScale,
			bitimage.NormalScale,
			bitimage.Color1,
			logoWidth,
			logoHeight,
			logoData,
		)
		if err != nil {
			t.Fatalf("StoreRasterGraphicsInBuffer failed: %v", err)
		}
		buffer = append(buffer, storeCmd...)

		// Step 3: Print the buffered graphics
		printCmd, err := cmd.Graphics.PrintBufferedGraphics(bitimage.FunctionCodePrint2)
		if err != nil {
			t.Fatalf("PrintBufferedGraphics failed: %v", err)
		}
		buffer = append(buffer, printCmd...)

		// Verify commands were generated
		if !bytes.Contains(buffer, []byte{shared.GS, '(', 'L'}) {
			t.Error("Buffer should contain graphics commands")
		}

		if len(buffer) < 100 {
			t.Error("Buffer should contain substantial graphics data")
		}
	})

	t.Run("multi-color graphics workflow", func(t *testing.T) {
		var buffer []byte

		// Store multiple color layers
		width := uint16(100)
		height := uint16(100)
		widthBytes := (int(width) + 7) / 8
		dataSize := widthBytes * int(height)

		colors := []bitimage.GraphicsColor{
			bitimage.Color1,
			bitimage.Color2,
			bitimage.Color3,
		}

		for _, color := range colors {
			data := test.RepeatByte(dataSize, byte(color))
			storeCmd, err := cmd.Graphics.StoreRasterGraphicsInBuffer(
				bitimage.Monochrome,
				bitimage.NormalScale,
				bitimage.NormalScale,
				color,
				width,
				height,
				data,
			)
			if err != nil {
				t.Fatalf("StoreRasterGraphicsInBuffer for color %v failed: %v", color, err)
			}
			buffer = append(buffer, storeCmd...)
		}

		// Print the combined graphics
		printCmd, err := cmd.Graphics.PrintBufferedGraphics(bitimage.FunctionCodePrint50)
		if err != nil {
			t.Fatalf("PrintBufferedGraphics failed: %v", err)
		}
		buffer = append(buffer, printCmd...)

		if len(buffer) < len(colors)*dataSize {
			t.Error("Buffer should contain data for all color layers")
		}
	})

	t.Run("column format graphics workflow", func(t *testing.T) {
		var buffer []byte

		// Use column format for vertical patterns
		width := uint16(256)
		height := uint16(64)
		heightBytes := (int(height) + 7) / 8
		dataSize := int(width) * heightBytes

		// Create vertical stripe pattern
		columnData := make([]byte, dataSize)
		for i := 0; i < dataSize; i++ {
			if (i/heightBytes)%2 == 0 {
				columnData[i] = 0xFF
			} else {
				columnData[i] = 0x00
			}
		}

		storeCmd, err := cmd.Graphics.StoreColumnGraphicsInBuffer(
			bitimage.NormalScale,
			bitimage.DoubleScale,
			bitimage.Color1,
			width,
			height,
			columnData,
		)
		if err != nil {
			t.Fatalf("StoreColumnGraphicsInBuffer failed: %v", err)
		}
		buffer = append(buffer, storeCmd...)

		printCmd, err := cmd.Graphics.PrintBufferedGraphics(bitimage.FunctionCodePrint2)
		if err != nil {
			t.Fatalf("PrintBufferedGraphics failed: %v", err)
		}
		buffer = append(buffer, printCmd...)

		if len(buffer) != dataSize+22 { // dataSize + command headers (15 for StoreColumn + 7 for Print)
			t.Errorf("Buffer length = %d, expected %d", len(buffer), dataSize+22)
		}
	})

	t.Run("legacy bit image mode workflow", func(t *testing.T) {
		var buffer []byte

		// Use legacy 8-dot single density mode
		width := uint16(100)
		data := test.RepeatByte(int(width), 0x55)

		legacyCmd, err := cmd.SelectBitImageMode(
			bitimage.SingleDensity8,
			width,
			data,
		)
		if err != nil {
			t.Fatalf("SelectBitImageMode failed: %v", err)
		}
		buffer = append(buffer, legacyCmd...)

		// Use legacy 24-dot double density mode
		width24 := uint16(50)
		data24 := test.RepeatByte(int(width24)*3, 0xAA)

		legacy24Cmd, err := cmd.SelectBitImageMode(
			bitimage.DoubleDensity24,
			width24,
			data24,
		)
		if err != nil {
			t.Fatalf("SelectBitImageMode 24-dot failed: %v", err)
		}
		buffer = append(buffer, legacy24Cmd...)

		if !bytes.Contains(buffer, []byte{shared.ESC, '*'}) {
			t.Error("Buffer should contain legacy bit image commands")
		}

		expectedLen := 5 + int(width) + 5 + int(width24)*3
		if len(buffer) != expectedLen {
			t.Errorf("Buffer length = %d, expected %d", len(buffer), expectedLen)
		}
	})
}

func TestIntegration_BitImage_ErrorHandling(t *testing.T) {
	cmd := bitimage.NewCommands()

	t.Run("invalid parameters cascade", func(t *testing.T) {
		// Invalid density combination
		_, err := cmd.Graphics.SetGraphicsDotDensity(
			bitimage.FunctionCodeDensity1,
			bitimage.Density180x180,
			bitimage.Density360x360,
		)
		if err == nil {
			t.Error("Mismatched densities should return error")
		}

		// Invalid tone
		_, err = cmd.Graphics.StoreRasterGraphicsInBuffer(
			99,
			bitimage.NormalScale,
			bitimage.NormalScale,
			bitimage.Color1,
			100,
			100,
			[]byte{},
		)
		if err == nil {
			t.Error("Invalid tone should return error")
		}

		// Invalid dimensions
		_, err = cmd.Graphics.StoreColumnGraphicsInBuffer(
			bitimage.NormalScale,
			bitimage.NormalScale,
			bitimage.Color1,
			0,
			0,
			[]byte{},
		)
		if err == nil {
			t.Error("Zero dimensions should return error")
		}
	})

	t.Run("data size validation", func(t *testing.T) {
		// Raster format with wrong data size
		width := uint16(100)
		height := uint16(50)
		wrongData := []byte{0xFF} // Should be much larger

		_, err := cmd.Graphics.StoreRasterGraphicsInBuffer(
			bitimage.Monochrome,
			bitimage.NormalScale,
			bitimage.NormalScale,
			bitimage.Color1,
			width,
			height,
			wrongData,
		)
		if err == nil {
			t.Error("Wrong data size should return error")
		}

		// Column format with wrong data size
		_, err = cmd.Graphics.StoreColumnGraphicsInBuffer(
			bitimage.NormalScale,
			bitimage.NormalScale,
			bitimage.Color1,
			width,
			height,
			wrongData,
		)
		if err == nil {
			t.Error("Wrong data size should return error")
		}
	})

	t.Run("height limits based on tone and scale", func(t *testing.T) {
		// Monochrome with normal scale - max 2400
		_, err := cmd.Graphics.StoreRasterGraphicsInBuffer(
			bitimage.Monochrome,
			bitimage.NormalScale,
			bitimage.NormalScale,
			bitimage.Color1,
			100,
			2401, // Exceeds limit
			test.RepeatByte(1000, 0xFF),
		)
		if err == nil {
			t.Error("Height exceeding monochrome normal limit should return error")
		}

		// Multiple tone with double scale - max 300
		_, err = cmd.Graphics.StoreRasterGraphicsInBuffer(
			bitimage.MultipleTone,
			bitimage.NormalScale,
			bitimage.DoubleScale,
			bitimage.Color1,
			100,
			301, // Exceeds limit
			test.RepeatByte(1000, 0xFF),
		)
		if err == nil {
			t.Error("Height exceeding multiple tone double limit should return error")
		}
	})
}

func TestIntegration_BitImage_ScalingCombinations(t *testing.T) {
	cmd := bitimage.NewCommands()

	scales := []struct {
		name       string
		horizontal bitimage.GraphicsScale
		vertical   bitimage.GraphicsScale
	}{
		{"normal", bitimage.NormalScale, bitimage.NormalScale},
		{"double width", bitimage.DoubleScale, bitimage.NormalScale},
		{"double height", bitimage.NormalScale, bitimage.DoubleScale},
		{"quadruple", bitimage.DoubleScale, bitimage.DoubleScale},
	}

	for _, scale := range scales {
		t.Run(scale.name, func(t *testing.T) {
			width := uint16(50)
			height := uint16(50)
			widthBytes := (int(width) + 7) / 8
			data := test.RepeatByte(widthBytes*int(height), 0xFF)

			storeCmd, err := cmd.Graphics.StoreRasterGraphicsInBuffer(
				bitimage.Monochrome,
				scale.horizontal,
				scale.vertical,
				bitimage.Color1,
				width,
				height,
				data,
			)
			if err != nil {
				t.Fatalf("StoreRasterGraphicsInBuffer with %s failed: %v", scale.name, err)
			}

			// Verify command structure includes scale parameters
			if len(storeCmd) < 15 {
				t.Error("Command should include scale parameters")
			}
			if storeCmd[8] != byte(scale.horizontal) {
				t.Errorf("Horizontal scale = %d, want %d", storeCmd[8], scale.horizontal)
			}
			if storeCmd[9] != byte(scale.vertical) {
				t.Errorf("Vertical scale = %d, want %d", storeCmd[9], scale.vertical)
			}
		})
	}
}

func TestIntegration_BitImage_LargeDataHandling(t *testing.T) {
	cmd := bitimage.NewCommands()

	t.Run("large raster graphics exceeding standard size", func(t *testing.T) {
		// Create data larger than 65535 bytes
		width := uint16(2000)
		height := uint16(2000)
		widthBytes := (int(width) + 7) / 8
		dataSize := widthBytes * int(height)
		largeData := test.RepeatByte(dataSize, 0x99)

		// Should use large format command
		storeCmd, err := cmd.Graphics.StoreRasterGraphicsInBufferLarge(
			bitimage.Monochrome,
			bitimage.NormalScale,
			bitimage.NormalScale,
			bitimage.Color1,
			width,
			height,
			largeData,
		)
		if err != nil {
			t.Fatalf("StoreRasterGraphicsInBufferLarge failed: %v", err)
		}

		// Verify it uses extended format (GS 8 L)
		if storeCmd[0] != shared.GS || storeCmd[1] != '8' || storeCmd[2] != 'L' {
			t.Error("Large data should use extended command format")
		}

		// Verify 32-bit size encoding
		if len(storeCmd) < 7 {
			t.Error("Command should have 32-bit size parameters")
		}
	})

	t.Run("large column graphics", func(t *testing.T) {
		width := uint16(2048)
		height := uint16(128)
		heightBytes := (int(height) + 7) / 8
		dataSize := int(width) * heightBytes
		largeData := test.RepeatByte(dataSize, 0xEE)

		storeCmd, err := cmd.Graphics.StoreColumnGraphicsInBufferLarge(
			bitimage.NormalScale,
			bitimage.NormalScale,
			bitimage.Color2,
			width,
			height,
			largeData,
		)
		if err != nil {
			t.Fatalf("StoreColumnGraphicsInBufferLarge failed: %v", err)
		}

		// Verify extended format
		if storeCmd[0] != shared.GS || storeCmd[1] != '8' || storeCmd[2] != 'L' {
			t.Error("Large column data should use extended command format")
		}
	})
}
