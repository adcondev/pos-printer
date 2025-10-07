package barcode_test

import (
	"bytes"
	"testing"

	"github.com/adcondev/pos-printer/escpos/barcode"
	"github.com/adcondev/pos-printer/escpos/common"
	"github.com/adcondev/pos-printer/utils/test"
)

func TestIntegration_Barcode_CompleteWorkflow(t *testing.T) {
	cmd := barcode.NewCommands()
	builder := test.NewBufferBuilder()

	t.Run("complete barcode printing workflow", func(t *testing.T) {

		// Configure HRI position
		hriPosCmd, err := cmd.SelectHRICharacterPosition(barcode.HRIBelow)
		if err != nil {
			t.Fatalf("SelectHRICharacterPosition: %v", err)
		}
		builder.Append(hriPosCmd)

		// Configure HRI font
		hriFontCmd, err := cmd.SelectFontForHRI(barcode.HRIFontB)
		if err != nil {
			t.Fatalf("SelectFontForHRI: %v", err)
		}
		builder.Append(hriFontCmd)

		// Configure dimensions
		heightCmd, err := cmd.SetBarcodeHeight(100)
		if err != nil {
			t.Fatalf("SetBarcodeHeight: %v", err)
		}
		builder.Append(heightCmd)

		widthCmd, err := cmd.SetBarcodeWidth(3)
		if err != nil {
			t.Fatalf("SetBarcodeWidth: %v", err)
		}
		builder.Append(widthCmd)

		// Print barcode
		barcodeCmd, err := cmd.PrintBarcode(barcode.CODE39, []byte("*TEST123*"))
		if err != nil {
			t.Fatalf("PrintBarcode: %v", err)
		}
		builder.Append(barcodeCmd)

		// Get final buffer
		buffer := builder.GetBuffer()

		// Verify buffer contains all commands
		if len(buffer) == 0 {
			t.Error("Buffer should contain commands")
		}

		// Verify specific command sequences
		expectedHRIPos := []byte{common.GS, 'H', 2}
		if !bytes.Equal(buffer[:3], expectedHRIPos) {
			t.Errorf("Buffer should start with HRI position command")
		}

		// Verify barcode command at end
		if !bytes.Contains(buffer, []byte("*TEST123*")) {
			t.Error("Buffer should contain barcode data")
		}
	})

	t.Run("multiple barcode types", func(t *testing.T) {
		barcodes := []struct {
			name      string
			symbology barcode.Symbology
			data      []byte
		}{
			{"UPC-A", barcode.UPCA, []byte("12345678901")},
			{"EAN-13", barcode.JAN13, []byte("4901234567890")},
			{"CODE39", barcode.CODE39, []byte("*CODE39*")},
			{"ITF", barcode.ITF, []byte("123456")},
		}

		for _, bc := range barcodes {
			t.Run(bc.name, func(t *testing.T) {
				builder := test.NewBufferBuilder()

				// Set height for this barcode
				heightCmd, _ := cmd.SetBarcodeHeight(80)
				builder.Append(heightCmd)

				// Print barcode
				barcodeCmd, err := cmd.PrintBarcode(bc.symbology, bc.data)
				if err != nil {
					t.Errorf("PrintBarcode(%s): %v", bc.name, err)
					return
				}
				builder.Append(barcodeCmd)

				// Verify command was generated
				if len(barcodeCmd) == 0 {
					t.Errorf("%s: empty barcode command", bc.name)
				}

				// Get final buffer
				buffer := builder.GetBuffer()

				// Verify data is in command
				if !bytes.Contains(buffer, bc.data) {
					t.Errorf("%s: buffer should contain barcode data", bc.name)
				}
			})
		}
	})

	t.Run("CODE128 with different code sets", func(t *testing.T) {
		builder := test.NewBufferBuilder()

		// Configure barcode appearance
		heightCmd, _ := cmd.SetBarcodeHeight(60)
		builder.Append(heightCmd)
		widthCmd, _ := cmd.SetBarcodeWidth(2)
		builder.Append(widthCmd)

		// CODE128 Set A (uppercase and control chars)
		codeACmd, err := cmd.PrintBarcodeWithCodeSet(
			barcode.CODE128,
			barcode.Code128SetA,
			[]byte("HELLO"),
		)
		if err != nil {
			t.Fatalf("CODE128 Set A: %v", err)
		}
		builder.Append(codeACmd)

		// CODE128 Set B (mixed case)
		codeBCmd, err := cmd.PrintBarcodeWithCodeSet(
			barcode.CODE128,
			barcode.Code128SetB,
			[]byte("Hello123"),
		)
		if err != nil {
			t.Fatalf("CODE128 Set B: %v", err)
		}
		builder.Append(codeBCmd)

		// CODE128 Set C (numeric pairs)
		codeCCmd, err := cmd.PrintBarcodeWithCodeSet(
			barcode.CODE128,
			barcode.Code128SetC,
			[]byte("123456"),
		)
		if err != nil {
			t.Fatalf("CODE128 Set C: %v", err)
		}
		builder.Append(codeCCmd)

		// Get final buffer
		buffer := builder.GetBuffer()

		// Verify all code sets are in buffer
		if !bytes.Contains(buffer, []byte{'{', 'A'}) {
			t.Error("Buffer should contain CODE128 Set A prefix")
		}
		if !bytes.Contains(buffer, []byte{'{', 'B'}) {
			t.Error("Buffer should contain CODE128 Set B prefix")
		}
		if !bytes.Contains(buffer, []byte{'{', 'C'}) {
			t.Error("Buffer should contain CODE128 Set C prefix")
		}
	})
}

func TestIntegration_Barcode_HRIConfiguration(t *testing.T) {
	cmd := barcode.NewCommands()

	t.Run("all HRI positions", func(t *testing.T) {
		positions := []struct {
			name     string
			position barcode.HRIPosition
		}{
			{"not printed", barcode.HRINotPrinted},
			{"above", barcode.HRIAbove},
			{"below", barcode.HRIBelow},
			{"both", barcode.HRIBoth},
		}

		for _, p := range positions {
			t.Run(p.name, func(t *testing.T) {
				builder := test.NewBufferBuilder()

				// Set HRI position
				posCmd, err := cmd.SelectHRICharacterPosition(p.position)
				if err != nil {
					t.Errorf("SelectHRICharacterPosition(%s): %v", p.name, err)
					return
				}
				builder.Append(posCmd)

				// Print a barcode to test HRI
				barcodeCmd, _ := cmd.PrintBarcode(barcode.UPCA, []byte("12345678901"))
				builder.Append(barcodeCmd)

				// Verify position command is correct
				expected := []byte{common.GS, 'H', byte(p.position)}
				if !bytes.Equal(posCmd, expected) {
					t.Errorf("%s: command = %#v, want %#v", p.name, posCmd, expected)
				}

				// Get final buffer
				buffer := builder.GetBuffer()

				// Verify all commands were added
				if len(buffer) < 18 {
					t.Error("Buffer seems too small for all commands")
				}
			})
		}
	})

	t.Run("all HRI fonts", func(t *testing.T) {
		fonts := []struct {
			name string
			font barcode.HRIFont
		}{
			{"Font A", barcode.HRIFontA},
			{"Font B", barcode.HRIFontB},
			{"Font C", barcode.HRIFontC},
			{"Font D", barcode.HRIFontD},
			{"Font E", barcode.HRIFontE},
		}

		for _, f := range fonts {
			t.Run(f.name, func(t *testing.T) {
				fontCmd, err := cmd.SelectFontForHRI(f.font)
				if err != nil {
					t.Errorf("SelectFontForHRI(%s): %v", f.name, err)
					return
				}

				expected := []byte{common.GS, 'f', byte(f.font)}
				if !bytes.Equal(fontCmd, expected) {
					t.Errorf("%s: command = %#v, want %#v", f.name, fontCmd, expected)
				}
			})
		}
	})
}

func TestIntegration_Barcode_DifferentSizes(t *testing.T) {
	cmd := barcode.NewCommands()

	t.Run("various heights", func(t *testing.T) {
		heights := []barcode.Height{
			barcode.MinHeight,
			50,
			100,
			barcode.DefaultHeight,
			200,
			barcode.MaxHeight,
		}

		for _, height := range heights {
			t.Run(string(rune(height)), func(t *testing.T) {
				builder := test.NewBufferBuilder()

				heightCmd, err := cmd.SetBarcodeHeight(height)
				if err != nil {
					t.Errorf("SetBarcodeHeight(%d): %v", height, err)
					return
				}
				builder.Append(heightCmd)

				// Print a barcode with this height
				barcodeCmd, _ := cmd.PrintBarcode(barcode.CODE39, []byte("TEST"))
				builder.Append(barcodeCmd)

				// Verify height was set
				if heightCmd[2] != byte(height) {
					t.Errorf("Height byte = %d, want %d", heightCmd[2], height)
				}

				// Get final buffer
				buffer := builder.GetBuffer()

				// Verify all commands were added
				if len(buffer) < 11 {
					t.Error("Buffer seems too small for all commands")
				}
			})
		}
	})

	t.Run("various widths", func(t *testing.T) {
		widths := []barcode.Width{
			barcode.MinWidth,
			barcode.DefaultWidth,
			4,
			5,
			barcode.MaxWidth,
		}

		for _, width := range widths {
			t.Run(string(rune(width)), func(t *testing.T) {
				widthCmd, err := cmd.SetBarcodeWidth(width)
				if err != nil {
					t.Errorf("SetBarcodeWidth(%d): %v", width, err)
					return
				}

				// Verify width was set
				if widthCmd[2] != byte(width) {
					t.Errorf("Width byte = %d, want %d", widthCmd[2], width)
				}
			})
		}
	})

	t.Run("extended width range", func(t *testing.T) {
		extendedWidths := []barcode.Width{
			barcode.ExtendedMinWidth,
			70,
			72,
			74,
			barcode.ExtendedMaxWidth,
		}

		for _, width := range extendedWidths {
			widthCmd, err := cmd.SetBarcodeWidth(width)
			if err != nil {
				t.Errorf("SetBarcodeWidth(%d): %v", width, err)
				continue
			}

			if widthCmd[2] != byte(width) {
				t.Errorf("Width byte = %d, want %d", widthCmd[2], width)
			}
		}
	})
}

func TestIntegration_Barcode_ErrorScenarios(t *testing.T) {
	cmd := barcode.NewCommands()

	t.Run("ITF with odd length", func(t *testing.T) {
		oddData := []byte("12345")

		// Assert validation
		test.AssertNumeric(t, oddData, "ITF should be numeric")
		if test.IsEvenLength(oddData) {
			t.Error("Test data should be odd length for this test")
		}

		// Function A version
		_, err := cmd.PrintBarcode(barcode.ITF, oddData)
		if err == nil {
			t.Error("ITF with odd length should return error")
		}
	})

	t.Run("CODE128 without code set", func(t *testing.T) {
		data := []byte("Hello")

		// Verify data is valid but missing code set prefix
		test.AssertPrintableASCII(t, data, "Data should be printable")

		_, err := cmd.PrintBarcode(barcode.CODE128, data)
		if err == nil {
			t.Error("CODE128 without code set should return error")
		}
	})

	t.Run("data length limits", func(t *testing.T) {
		// Empty data
		var emptyData []byte
		test.AssertInvalidLength(t, emptyData, 1, 255, "Empty data validation")

		_, err := cmd.PrintBarcode(barcode.CODE39, emptyData)
		if err == nil {
			t.Error("Empty data should return error")
		}

		// Data too long (>255 bytes)
		longData := test.RepeatByte(256, 'A')
		if test.ValidateLength(longData, 1, 255) {
			t.Error("Test data should be too long")
		}

		_, err = cmd.PrintBarcode(barcode.CODE39B, longData)
		if err == nil {
			t.Error("Data >255 bytes should return error")
		}
	})
}

func TestIntegration_Barcode_RealWorldScenarios(t *testing.T) {
	cmd := barcode.NewCommands()

	t.Run("retail receipt with UPC", func(t *testing.T) {
		builder := test.NewBufferBuilder()

		// Configure for retail receipt
		hriCmd, _ := cmd.SelectHRICharacterPosition(barcode.HRIBelow)
		builder.Append(hriCmd)

		heightCmd, _ := cmd.SetBarcodeHeight(50)
		builder.Append(heightCmd)

		widthCmd, _ := cmd.SetBarcodeWidth(2)
		builder.Append(widthCmd)

		// Print multiple product UPCs
		products := []string{
			"12345678901", // Product 1
			"98765432109", // Product 2
			"55555555555", // Product 3
		}

		for _, upc := range products {
			barcodeCmd, err := cmd.PrintBarcode(barcode.UPCA, []byte(upc))
			if err != nil {
				t.Errorf("Failed to print UPC %s: %v", upc, err)
				continue
			}
			builder.Append(barcodeCmd)
		}

		// Get final buffer
		buffer := builder.GetBuffer()

		// Verify all UPCs are in buffer
		for _, upc := range products {
			if !bytes.Contains(buffer, []byte(upc)) {
				t.Errorf("Buffer should contain UPC %s", upc)
			}
		}
	})

	t.Run("inventory label with CODE128", func(t *testing.T) {
		builder := test.NewBufferBuilder()

		// Configure for inventory labels
		hriCmd, _ := cmd.SelectHRICharacterPosition(barcode.HRIBoth)
		builder.Append(hriCmd)

		fontCmd, _ := cmd.SelectFontForHRI(barcode.HRIFontA)
		builder.Append(fontCmd)

		heightCmd, _ := cmd.SetBarcodeHeight(80)
		builder.Append(heightCmd)

		// Print inventory codes using different code sets
		inventoryCodes := []struct {
			codeSet barcode.Code128Set
			data    []byte
		}{
			{barcode.Code128SetB, []byte("INV-2024-001")},
			{barcode.Code128SetC, []byte("123456789012")},
			{barcode.Code128SetA, []byte("WAREHOUSE-A")},
		}

		for _, inv := range inventoryCodes {
			barcodeCmd, err := cmd.PrintBarcodeWithCodeSet(
				barcode.CODE128,
				inv.codeSet,
				inv.data,
			)
			if err != nil {
				t.Errorf("Failed to print inventory code: %v", err)
				continue
			}
			builder.Append(barcodeCmd)
		}

		// Get final buffer
		buffer := builder.GetBuffer()

		// Verify all commands were added
		if len(buffer) < 62 {
			t.Error("Buffer seems too small for all commands")
		}
	})

	t.Run("shipping label with mixed barcodes", func(t *testing.T) {
		builder := test.NewBufferBuilder()

		// Tracking number as CODE128
		heightCmd, _ := cmd.SetBarcodeHeight(60)
		builder.Append(heightCmd)

		trackingCmd, err := cmd.PrintBarcodeWithCodeSet(
			barcode.CODE128,
			barcode.Code128SetB,
			[]byte("TRK-123456789"),
		)
		if err != nil {
			t.Fatalf("Failed to print tracking code: %v", err)
		}
		builder.Append(trackingCmd)

		// ZIP code as CODABAR
		zipCmd, err := cmd.PrintBarcode(barcode.CODABAR, []byte("A12345B"))
		if err != nil {
			t.Fatalf("Failed to print ZIP code: %v", err)
		}
		builder.Append(zipCmd)

		// Package ID as CODE39
		pkgCmd, err := cmd.PrintBarcode(barcode.CODE39, []byte("*PKG001*"))
		if err != nil {
			t.Fatalf("Failed to print package ID: %v", err)
		}
		builder.Append(pkgCmd)

		buffer := builder.GetBuffer()

		// Verify all barcode types are present
		if !bytes.Contains(buffer, []byte("TRK-123456789")) {
			t.Error("Buffer should contain tracking number")
		}
		if !bytes.Contains(buffer, []byte("A12345B")) {
			t.Error("Buffer should contain ZIP code")
		}
		if !bytes.Contains(buffer, []byte("*PKG001*")) {
			t.Error("Buffer should contain package ID")
		}
	})
}

func TestIntegration_Barcode_DataValidation(t *testing.T) {
	cmd := barcode.NewCommands()

	t.Run("UPC barcodes require numeric", func(t *testing.T) {
		validUPC := []byte("12345678901")
		test.AssertNumeric(t, validUPC, "UPC-A data")
		test.AssertValidLength(t, validUPC, 11, 12, "UPC-A length")

		_, err := cmd.PrintBarcode(barcode.UPCA, validUPC)
		if err != nil {
			t.Errorf("Valid UPC should not error: %v", err)
		}
	})

	t.Run("CODE39 validation", func(t *testing.T) {
		validCode39 := []byte("ABC-123")

		// Extract alphabetic part for case check
		alphaOnly := test.FilterBytes(validCode39, test.IsAlphanumericByte)
		test.AssertUppercase(t, alphaOnly, "CODE39 letters should be uppercase")

		// Check allowed character set
		allowed := test.Code39Charset
		test.AssertContainsOnly(t, validCode39, allowed, "CODE39 character set")

		test.AssertNotEmpty(t, validCode39, "CODE39 data should not be empty")

		_, err := cmd.PrintBarcode(barcode.CODE39, validCode39)
		if err != nil {
			t.Errorf("Valid CODE39 should not error: %v", err)
		}
	})

	t.Run("ITF requires even numeric", func(t *testing.T) {
		validITF := []byte("123456")

		test.AssertNumeric(t, validITF, "ITF data should be numeric")
		test.AssertEvenLength(t, validITF, "ITF data should be even length")
		test.AssertValidLength(t, validITF, 2, 254, "ITF length range")

		_, err := cmd.PrintBarcode(barcode.ITF, validITF)
		if err != nil {
			t.Errorf("Valid ITF should not error: %v", err)
		}
	})

	t.Run("Function A barcodes have null terminator", func(t *testing.T) {
		result, err := cmd.PrintBarcode(barcode.UPCA, []byte("12345678901"))
		if err != nil {
			t.Fatalf("PrintBarcode failed: %v", err)
		}

		test.AssertHasNullTerminator(t, result, "Function A command format")
		test.AssertHasSuffix(t, result, []byte{common.NUL}, "NUL terminator check")
	})

	t.Run("CODABAR start/stop validation", func(t *testing.T) {
		validCodabar := []byte("A12345B")

		// Check start/stop are in allowed range
		startStopChars := []byte{'A', 'B', 'C', 'D', 'a', 'b', 'c', 'd'}
		if !test.ContainsAny([]byte{validCodabar[0]}, startStopChars) {
			t.Error("Invalid CODABAR start character")
		}
		if !test.ContainsAny([]byte{validCodabar[len(validCodabar)-1]}, startStopChars) {
			t.Error("Invalid CODABAR stop character")
		}

		_, err := cmd.PrintBarcode(barcode.CODABAR, validCodabar)
		if err != nil {
			t.Errorf("Valid CODABAR should not error: %v", err)
		}
	})
}
