package print_test

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/adcondev/pos-printer/escpos/common"
	"github.com/adcondev/pos-printer/escpos/print"
)

func TestIntegration_Print_CompleteDocument(t *testing.T) {
	cmd := print.NewCommands()

	t.Run("receipt printing workflow", func(t *testing.T) {
		var buffer []byte

		// Header
		text, err := cmd.Text("STORE NAME")
		if err != nil {
			t.Fatalf("Text header: %v", err)
		}
		buffer = append(buffer, text...)
		buffer = append(buffer, cmd.PrintAndLineFeed()...)

		// Items
		for i := 1; i <= 3; i++ {
			item, err := cmd.Text(fmt.Sprintf("Item %d", i))
			if err != nil {
				t.Fatalf("Text item %d: %v", i, err)
			}
			buffer = append(buffer, item...)
			buffer = append(buffer, cmd.PrintAndLineFeed()...)
		}

		// Feed and cut
		buffer = append(buffer, cmd.PrintAndFeedPaper(50)...)
		buffer = append(buffer, cmd.FormFeed()...)

		// Verify document structure
		if !bytes.Contains(buffer, []byte("STORE NAME")) {
			t.Error("Missing header")
		}
		if !bytes.Contains(buffer, []byte("Item 1")) {
			t.Error("Missing items")
		}
		if bytes.Count(buffer, []byte{print.LF}) < 3 {
			t.Error("Missing line feeds")
		}
	})
}

func TestIntegration_PageMode_Workflow(t *testing.T) {
	cmd := print.NewCommands()

	t.Run("page mode operations", func(t *testing.T) {
		var buffer []byte

		// Enter page mode (would be in escpos package)
		// ...

		// Add text
		text, _ := cmd.Text("Page Mode Text")
		buffer = append(buffer, text...)

		// Print page
		buffer = append(buffer, cmd.PrintDataInPageMode()...)

		// Cancel if needed
		buffer = append(buffer, cmd.CancelData()...)

		// Verify page mode commands
		if !bytes.Contains(buffer, []byte{print.CAN}) {
			t.Error("Missing CAN command")
		}
	})
}

func TestIntegration_Print_ErrorHandling(t *testing.T) {
	cmd := print.NewCommands()

	t.Run("handles empty text error", func(t *testing.T) {
		_, err := cmd.Text("")
		if !errors.Is(err, common.ErrEmptyBuffer) {
			t.Errorf("Text(\"\") error = %v, want %v", err, common.ErrEmptyBuffer)
		}
	})

	t.Run("handles page mode errors", func(t *testing.T) {
		_, err := cmd.PrintAndReverseFeed(100) // Exceeds max
		if !errors.Is(err, print.ErrPrintReverseFeed) {
			t.Errorf("PrintAndReverseFeed(100) error = %v, want %v",
				err, print.ErrPrintReverseFeed)
		}

		_, err = cmd.PrintAndReverseFeedLines(10) // Exceeds max
		if !errors.Is(err, print.ErrPrintReverseFeedLines) {
			t.Errorf("PrintAndReverseFeedLines(10) error = %v, want %v",
				err, print.ErrPrintReverseFeedLines)
		}
	})
}

func TestIntegration_Print_FormattingWorkflow(t *testing.T) {
	cmd := print.NewCommands()

	t.Run("formats complex text correctly", func(t *testing.T) {
		input := "Header\n\tItem1\n\tItem2\rTotal"
		result, err := cmd.Text(input)

		if err != nil {
			t.Fatalf("Text() unexpected error: %v", err)
		}

		// Verify formatting
		if !bytes.Contains(result, []byte{print.LF}) {
			t.Error("Should contain LF")
		}
		if !bytes.Contains(result, []byte{common.HT}) {
			t.Error("Should contain HT")
		}
		if !bytes.Contains(result, []byte{print.CR}) {
			t.Error("Should contain CR")
		}
	})
}

func TestIntegration_Print_PageModeComplete(t *testing.T) {
	cmd := print.NewCommands()

	t.Run("complete page mode workflow", func(t *testing.T) {
		var commands []byte

		// Build page
		text, _ := cmd.Text("Page Content")
		commands = append(commands, text...)

		// Reverse operations
		reverse, _ := cmd.PrintAndReverseFeed(10)
		commands = append(commands, reverse...)

		reverseLine, _ := cmd.PrintAndReverseFeedLines(1)
		commands = append(commands, reverseLine...)

		// Forward feed
		forward := cmd.PrintAndFeedLines(5)
		commands = append(commands, forward...)

		// Print page
		printPage := cmd.PrintDataInPageMode()
		commands = append(commands, printPage...)

		// Clear if needed
		clean := cmd.CancelData()
		commands = append(commands, clean...)

		// Verify sequence integrity
		if len(commands) == 0 {
			t.Error("No commands generated")
		}
		if !bytes.Contains(commands, []byte("Page Content")) {
			t.Error("Missing page content")
		}
		if !bytes.Contains(commands, []byte{common.ESC, print.FF}) {
			t.Error("Missing print page command")
		}
		if !bytes.Contains(commands, []byte{print.CAN}) {
			t.Error("Missing cancel command")
		}
	})
}
