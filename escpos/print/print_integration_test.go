package print_test

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/adcondev/pos-printer/escpos/common"
	"github.com/adcondev/pos-printer/escpos/print"
)

func TestIntegration_Print_StandardWorkflow(t *testing.T) {
	cmd := print.NewCommands()

	t.Run("complete print workflow", func(t *testing.T) {
		// Setup
		var buffer []byte

		// Execute
		text1, err := cmd.Text("First line")
		if err != nil {
			t.Fatalf("Text: %v", err)
		}
		buffer = append(buffer, text1...)

		buffer = append(buffer, cmd.PrintAndLineFeed()...)

		text2, err := cmd.Text("Second line")
		if err != nil {
			t.Fatalf("Text: %v", err)
		}
		buffer = append(buffer, text2...)

		buffer = append(buffer, cmd.PrintAndCarriageReturn()...)

		buffer = append(buffer, cmd.PrintAndFeedPaper(50)...)

		buffer = append(buffer, cmd.PrintAndFeedLines(3)...)

		buffer = append(buffer, cmd.FormFeed()...)

		// Verify
		if !bytes.Contains(buffer, []byte("First line")) {
			t.Error("Missing first line")
		}
		if !bytes.Contains(buffer, []byte("Second line")) {
			t.Error("Missing second line")
		}
		if !bytes.Contains(buffer, []byte{print.LF}) {
			t.Error("Missing line feed")
		}
		if !bytes.Contains(buffer, []byte{print.CR}) {
			t.Error("Missing carriage return")
		}
		if !bytes.Contains(buffer, []byte{print.FF}) {
			t.Error("Missing form feed")
		}
	})

	t.Run("text formatting workflow", func(t *testing.T) {
		// Setup
		testCases := []struct {
			name     string
			input    string
			contains []byte
		}{
			{"newline", "Line1\nLine2", []byte{print.LF}},
			{"carriage return", "Line1\rLine2", []byte{print.CR}},
			{"tab", "Col1\tCol2", []byte{common.HT}},
			{"mixed", "A\tB\nC\rD", []byte{common.HT, print.LF, print.CR}},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Execute
				result, err := cmd.Text(tc.input)
				if err != nil {
					t.Fatalf("Text(%s): %v", tc.name, err)
				}

				// Verify
				for _, expected := range tc.contains {
					if !bytes.Contains(result, []byte{expected}) {
						t.Errorf("Result should contain %#v", expected)
					}
				}
			})
		}
	})

	t.Run("feed operations sequence", func(t *testing.T) {
		// Setup
		var buffer []byte
		feeds := []struct {
			name   string
			units  byte
			opType string
		}{
			{"small", 10, "paper"},
			{"medium", 50, "paper"},
			{"large", 100, "paper"},
			{"1 line", 1, "lines"},
			{"5 lines", 5, "lines"},
			{"10 lines", 10, "lines"},
		}

		// Execute
		for _, f := range feeds {
			if f.opType == "paper" {
				buffer = append(buffer, cmd.PrintAndFeedPaper(f.units)...)
			} else {
				buffer = append(buffer, cmd.PrintAndFeedLines(f.units)...)
			}
		}

		// Verify
		paperFeedCount := bytes.Count(buffer, []byte{common.ESC, 'J'})
		if paperFeedCount != 3 {
			t.Errorf("Paper feed count = %d, want 3", paperFeedCount)
		}
		lineFeedCount := bytes.Count(buffer, []byte{common.ESC, 'd'})
		if lineFeedCount != 3 {
			t.Errorf("Line feed count = %d, want 3", lineFeedCount)
		}
	})
}

func TestIntegration_Print_PageModeWorkflow(t *testing.T) {
	cmd := print.NewCommands()

	t.Run("complete page mode workflow", func(t *testing.T) {
		// Setup
		var buffer []byte

		// Execute
		text, err := cmd.Text("Page content")
		if err != nil {
			t.Fatalf("Text: %v", err)
		}
		buffer = append(buffer, text...)

		buffer = append(buffer, cmd.PrintDataInPageMode()...)

		text2, err := cmd.Text("More content")
		if err != nil {
			t.Fatalf("Text: %v", err)
		}
		buffer = append(buffer, text2...)

		buffer = append(buffer, cmd.CancelData()...)

		text3, err := cmd.Text("New content")
		if err != nil {
			t.Fatalf("Text: %v", err)
		}
		buffer = append(buffer, text3...)

		buffer = append(buffer, cmd.FormFeed()...)

		// Verify
		if !bytes.Contains(buffer, []byte{common.ESC, print.FF}) {
			t.Error("Missing PrintDataInPageMode command")
		}
		if !bytes.Contains(buffer, []byte{print.CAN}) {
			t.Error("Missing CancelData command")
		}
	})

	t.Run("reverse feed operations", func(t *testing.T) {
		// Setup
		var buffer []byte

		// Execute
		buffer = append(buffer, cmd.PrintAndFeedPaper(30)...)

		reverseCmd, err := cmd.PrintAndReverseFeed(20)
		if err != nil {
			t.Fatalf("PrintAndReverseFeed: %v", err)
		}
		buffer = append(buffer, reverseCmd...)

		reverseLines, err := cmd.PrintAndReverseFeedLines(1)
		if err != nil {
			t.Fatalf("PrintAndReverseFeedLines: %v", err)
		}
		buffer = append(buffer, reverseLines...)

		buffer = append(buffer, cmd.PrintAndFeedLines(2)...)

		// Verify
		if !bytes.Contains(buffer, []byte{common.ESC, 'K'}) {
			t.Error("Missing reverse feed command")
		}
		if !bytes.Contains(buffer, []byte{common.ESC, 'e'}) {
			t.Error("Missing reverse lines command")
		}
	})

	t.Run("page mode print cycles", func(t *testing.T) {
		// Setup
		var buffer []byte

		// Execute
		for i := 0; i < 3; i++ {
			// Add content
			text, _ := cmd.Text(fmt.Sprintf("Cycle %d", i+1))
			buffer = append(buffer, text...)

			// Print page
			buffer = append(buffer, cmd.PrintDataInPageMode()...)
		}

		buffer = append(buffer, cmd.FormFeed()...)

		// Verify
		printCount := bytes.Count(buffer, []byte{common.ESC, print.FF})
		if printCount != 3 {
			t.Errorf("Print page count = %d, want 3", printCount)
		}
	})
}

func TestIntegration_Print_EdgeCases(t *testing.T) {
	cmd := print.NewCommands()

	t.Run("maximum values", func(t *testing.T) {
		// Setup
		var buffer []byte

		// Execute
		buffer = append(buffer, cmd.PrintAndFeedPaper(255)...)
		buffer = append(buffer, cmd.PrintAndFeedLines(255)...)

		maxReverse, err := cmd.PrintAndReverseFeed(print.MaxReverseMotionUnits)
		if err != nil {
			t.Fatalf("PrintAndReverseFeed max: %v", err)
		}
		buffer = append(buffer, maxReverse...)

		maxReverseLines, err := cmd.PrintAndReverseFeedLines(print.MaxReverseFeedLines)
		if err != nil {
			t.Fatalf("PrintAndReverseFeedLines max: %v", err)
		}
		buffer = append(buffer, maxReverseLines...)

		// Verify
		if buffer[2] != 255 { // PrintAndFeedPaper value
			t.Error("Should accept maximum paper feed")
		}
		if buffer[5] != 255 { // PrintAndFeedLines value
			t.Error("Should accept maximum line feed")
		}
	})

	t.Run("minimum values", func(t *testing.T) {
		// Setup
		var buffer []byte

		// Execute
		buffer = append(buffer, cmd.PrintAndFeedPaper(0)...)
		buffer = append(buffer, cmd.PrintAndFeedLines(0)...)

		zeroReverse, err := cmd.PrintAndReverseFeed(0)
		if err != nil {
			t.Fatalf("PrintAndReverseFeed(0): %v", err)
		}
		buffer = append(buffer, zeroReverse...)

		zeroReverseLines, err := cmd.PrintAndReverseFeedLines(0)
		if err != nil {
			t.Fatalf("PrintAndReverseFeedLines(0): %v", err)
		}
		buffer = append(buffer, zeroReverseLines...)

		// Verify
		if len(buffer) != 12 { // 4 commands × 3 bytes
			t.Errorf("Buffer length = %d, want 12", len(buffer))
		}
	})

	t.Run("text size limits", func(t *testing.T) {
		// Setup
		maxText := strings.Repeat("A", 1024)

		// TODO: Correct the test from max/min values perspective. Check values defined in Print package.
		// Execute
		result, err := cmd.Text(maxText)
		if err != nil {
			t.Fatalf("Text with 1024 chars: %v", err)
		}
		if len(result) != 1024 {
			t.Errorf("Result length = %d, want 1024", len(result))
		}
		// TODO: Check if config values apply here
		// Very large text (should handle or error appropriately)
		largeText := strings.Repeat("B", 0xFFFF+1)
		_, err = cmd.Text(largeText)

		// Verify
		if err == nil {
			// If no error, verify it processed
			t.Log("Large text accepted")
		} else if !errors.Is(err, print.ErrInvalidTextTooLarge) {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}

func TestIntegration_Print_ErrorConditions(t *testing.T) {
	cmd := print.NewCommands()

	t.Run("invalid parameters", func(t *testing.T) {
		// Empty text
		_, err := cmd.Text("")
		if !errors.Is(err, print.ErrInvalidEmptyText) {
			t.Errorf("Text('') error = %v, want %v", err, print.ErrInvalidEmptyText)
		}

		// Reverse feed beyond limit
		_, err = cmd.PrintAndReverseFeed(print.MaxReverseMotionUnits + 1)
		if !errors.Is(err, print.ErrInvalidReverseUnits) {
			t.Errorf("PrintAndReverseFeed overflow error = %v, want %v",
				err, print.ErrInvalidReverseUnits)
		}

		// Reverse lines beyond limit
		_, err = cmd.PrintAndReverseFeedLines(print.MaxReverseFeedLines + 1)
		if !errors.Is(err, print.ErrInvalidReverseLines) {
			t.Errorf("PrintAndReverseFeedLines overflow error = %v, want %v",
				err, print.ErrInvalidReverseLines)
		}
	})

	t.Run("boundary conditions", func(t *testing.T) {
		// Exactly at limits
		_, err := cmd.PrintAndReverseFeed(print.MaxReverseMotionUnits)
		if err != nil {
			t.Errorf("Should accept maximum reverse units: %v", err)
		}

		_, err = cmd.PrintAndReverseFeedLines(print.MaxReverseFeedLines)
		if err != nil {
			t.Errorf("Should accept maximum reverse lines: %v", err)
		}

		// Just over limits
		_, err = cmd.PrintAndReverseFeed(49)
		if err == nil {
			t.Error("Should reject reverse units over 48")
		}

		_, err = cmd.PrintAndReverseFeedLines(3)
		if err == nil {
			t.Error("Should reject reverse lines over 2")
		}
	})
}

func TestIntegration_Print_RealWorldScenarios(t *testing.T) {
	cmd := print.NewCommands()

	t.Run("receipt printing", func(t *testing.T) {
		// Setup
		var buffer []byte

		// Execute
		// Header
		header, _ := cmd.Text("========== RECEIPT ==========")
		buffer = append(buffer, header...)
		buffer = append(buffer, cmd.PrintAndLineFeed()...)
		buffer = append(buffer, cmd.PrintAndLineFeed()...)

		// Date/Time
		datetime, _ := cmd.Text("2025-01-01 12:00:00")
		buffer = append(buffer, datetime...)
		buffer = append(buffer, cmd.PrintAndLineFeed()...)

		// Items
		items := []string{
			"Item 1              $10.00",
			"Item 2              $15.00",
			"Item 3              $20.00",
		}
		for _, item := range items {
			itemText, _ := cmd.Text(item)
			buffer = append(buffer, itemText...)
			buffer = append(buffer, cmd.PrintAndLineFeed()...)
		}

		// Separator
		buffer = append(buffer, cmd.PrintAndFeedLines(1)...)
		separator, _ := cmd.Text("----------------------------")
		buffer = append(buffer, separator...)
		buffer = append(buffer, cmd.PrintAndLineFeed()...)

		// Total
		total, _ := cmd.Text("TOTAL:              $45.00")
		buffer = append(buffer, total...)
		buffer = append(buffer, cmd.PrintAndLineFeed()...)

		// Feed and cut
		buffer = append(buffer, cmd.PrintAndFeedPaper(100)...)
		buffer = append(buffer, cmd.FormFeed()...)

		// Verify receipt structure
		if !bytes.Contains(buffer, []byte("RECEIPT")) {
			t.Error("Missing header")
		}
		if !bytes.Contains(buffer, []byte("TOTAL")) {
			t.Error("Missing total")
		}
		lineCount := bytes.Count(buffer, []byte{print.LF})
		if lineCount < 8 {
			t.Error("Insufficient line feeds for receipt")
		}
	})

	t.Run("label printing with alignment", func(t *testing.T) {
		var buffer []byte

		// Product name (would be centered)
		name, _ := cmd.Text("PRODUCT NAME")
		buffer = append(buffer, name...)
		buffer = append(buffer, cmd.PrintAndLineFeed()...)

		// Small feed
		buffer = append(buffer, cmd.PrintAndFeedPaper(20)...)

		// Barcode area (would be positioned)
		barcode, _ := cmd.Text("||||||||||||||||")
		buffer = append(buffer, barcode...)
		buffer = append(buffer, cmd.PrintAndLineFeed()...)

		// Price (would be emphasized)
		price, _ := cmd.Text("$99.99")
		buffer = append(buffer, price...)
		buffer = append(buffer, cmd.PrintAndLineFeed()...)

		// Feed to cut position
		buffer = append(buffer, cmd.PrintAndFeedPaper(80)...)

		// Verify label structure
		if !bytes.Contains(buffer, []byte("PRODUCT")) {
			t.Error("Missing product name")
		}
		if !bytes.Contains(buffer, []byte("$")) {
			t.Error("Missing price")
		}
	})

	t.Run("multi-page document", func(t *testing.T) {
		var buffer []byte

		pages := []string{
			"Page 1 Content",
			"Page 2 Content",
			"Page 3 Content",
		}

		for i, content := range pages {
			// Page content
			text, _ := cmd.Text(content)
			buffer = append(buffer, text...)
			buffer = append(buffer, cmd.PrintAndLineFeed()...)

			// Page break
			if i < len(pages)-1 {
				buffer = append(buffer, cmd.PrintAndFeedPaper(200)...)
				// Could add form feed for page boundary
				buffer = append(buffer, cmd.FormFeed()...)
			}
		}

		// Final form feed
		buffer = append(buffer, cmd.FormFeed()...)

		// Verify multi-page structure
		pageCount := bytes.Count(buffer, []byte{print.FF})
		if pageCount != 3 {
			t.Errorf("Form feed count = %d, want 3", pageCount)
		}
	})

	t.Run("columnar data with tabs", func(t *testing.T) {
		var buffer []byte

		// Header with tabs
		header, _ := cmd.Text("QTY\tITEM\tPRICE")
		buffer = append(buffer, header...)
		buffer = append(buffer, cmd.PrintAndLineFeed()...)

		// Data rows
		rows := []string{
			"1\tApple\t$1.00",
			"2\tOrange\t$2.00",
			"3\tBanana\t$1.50",
		}

		for _, row := range rows {
			rowText, _ := cmd.Text(row)
			buffer = append(buffer, rowText...)
			buffer = append(buffer, cmd.PrintAndCarriageReturn()...)
		}

		// Verify tab formatting
		tabCount := bytes.Count(buffer, []byte{common.HT})
		if tabCount != 8 { // 2 tabs per row × 3 rows + 2 in header
			t.Errorf("Tab count = %d, want 8", tabCount)
		}
	})
}

func TestIntegration_Print_FormattingConsistency(t *testing.T) {
	cmd := print.NewCommands()

	t.Run("formatting preservation", func(t *testing.T) {
		inputs := []struct {
			name     string
			input    string
			expected []byte
		}{
			{"simple", "Hello", []byte("Hello")},
			{"newline", "A\nB", []byte{'A', print.LF, 'B'}},
			{"return", "A\rB", []byte{'A', print.CR, 'B'}},
			{"tab", "A\tB", []byte{'A', common.HT, 'B'}},
			{"mixed", "A\n\r\tB", []byte{'A', print.LF, print.CR, common.HT, 'B'}},
		}

		for _, tc := range inputs {
			t.Run(tc.name, func(t *testing.T) {
				result, err := cmd.Text(tc.input)
				if err != nil {
					t.Fatalf("Text error: %v", err)
				}

				if !bytes.Equal(result, tc.expected) {
					t.Errorf("Text(%q) = %#v, want %#v",
						tc.input, result, tc.expected)
				}
			})
		}
	})

	t.Run("control character handling", func(t *testing.T) {
		// Test that control characters are properly converted
		input := "Line1\nLine2\rLine3\tColumn"
		result, _ := cmd.Text(input)

		// Count control characters
		lfCount := bytes.Count(result, []byte{print.LF})
		crCount := bytes.Count(result, []byte{print.CR})
		htCount := bytes.Count(result, []byte{common.HT})

		if lfCount != 1 {
			t.Errorf("LF count = %d, want 1", lfCount)
		}
		if crCount != 1 {
			t.Errorf("CR count = %d, want 1", crCount)
		}
		if htCount != 1 {
			t.Errorf("HT count = %d, want 1", htCount)
		}
	})
}
