package print_test

import (
	"bytes"
	"testing"

	"github.com/adcondev/pos-printer/escpos/common"
	"github.com/adcondev/pos-printer/escpos/print"
)

// ============================================================================
// Interface Composition Tests
// ============================================================================

func TestPagePrint_Implements_PageCapability(t *testing.T) {
	pp := &print.PagePrint{}

	// Verify PagePrint can be used as PageCapability
	var pageCapability print.PageCapability = pp

	// Test PrintDataInPageMode through interface
	result := pageCapability.PrintDataInPageMode()
	expectedBytes := []byte{common.ESC, common.FF}

	if !bytes.Equal(result, expectedBytes) {
		t.Errorf("PageCapability.PrintDataInPageMode() = %#v, want %#v",
			result, expectedBytes)
	}

	// Test PrintAndFeedLines through interface
	lines := byte(5)
	feedResult, err := pageCapability.PrintAndFeedLines(lines)
	if err != nil {
		t.Errorf("PageCapability.PrintAndFeedLines(%d) unexpected error: %v", lines, err)
	}

	expectedFeed := []byte{common.ESC, 'd', lines}
	if !bytes.Equal(feedResult, expectedFeed) {
		t.Errorf("PageCapability.PrintAndFeedLines(%d) = %#v, want %#v",
			lines, feedResult, expectedFeed)
	}
}

func TestPagePrint_Implements_ReverseCapability(t *testing.T) {
	pp := &print.PagePrint{}

	// Verify PagePrint can be used as ReverseCapability
	var reverseCapability print.ReverseCapability = pp

	t.Run("PrintAndReverseFeed", func(t *testing.T) {
		units := byte(10)
		result, err := reverseCapability.PrintAndReverseFeed(units)

		if err != nil {
			t.Errorf("ReverseCapability.PrintAndReverseFeed(%d) unexpected error: %v",
				units, err)
		}

		expected := []byte{common.ESC, 'K', units}
		if !bytes.Equal(result, expected) {
			t.Errorf("ReverseCapability.PrintAndReverseFeed(%d) = %#v, want %#v",
				units, result, expected)
		}
	})

	t.Run("PrintAndReverseFeedLines", func(t *testing.T) {
		lines := byte(1)
		result, err := reverseCapability.PrintAndReverseFeedLines(lines)

		if err != nil {
			t.Errorf("ReverseCapability.PrintAndReverseFeedLines(%d) unexpected error: %v",
				lines, err)
		}

		expected := []byte{common.ESC, 'e', lines}
		if !bytes.Equal(result, expected) {
			t.Errorf("ReverseCapability.PrintAndReverseFeedLines(%d) = %#v, want %#v",
				lines, result, expected)
		}
	})
}

func TestPagePrint_Implements_PageModeCapability(t *testing.T) {
	pp := &print.PagePrint{}

	// Verify PagePrint can be used as complete PageModeCapability
	var pageModeCapability print.PageModeCapability = pp

	// PageModeCapability embeds both PageCapability and ReverseCapability
	// Test through the composite interface

	t.Run("PageCapability methods", func(t *testing.T) {
		// Test PrintDataInPageMode
		pageResult := pageModeCapability.PrintDataInPageMode()
		if len(pageResult) != 2 {
			t.Errorf("PageModeCapability.PrintDataInPageMode() returned %d bytes, want 2",
				len(pageResult))
		}

		// Test PrintAndFeedLines
		_, err := pageModeCapability.PrintAndFeedLines(3)
		if err != nil {
			t.Errorf("PageModeCapability.PrintAndFeedLines(3) unexpected error: %v", err)
		}
	})

	t.Run("ReverseCapability methods", func(t *testing.T) {
		// Test PrintAndReverseFeed
		_, err := pageModeCapability.PrintAndReverseFeed(5)
		if err != nil {
			t.Errorf("PageModeCapability.PrintAndReverseFeed(5) unexpected error: %v", err)
		}

		// Test PrintAndReverseFeedLines
		_, err = pageModeCapability.PrintAndReverseFeedLines(1)
		if err != nil {
			t.Errorf("PageModeCapability.PrintAndReverseFeedLines(1) unexpected error: %v", err)
		}
	})
}

func TestInterfaceComposition_Polymorphism(t *testing.T) {
	// This test demonstrates that the same struct can be used
	// through different interface views

	pp := &print.PagePrint{}

	// Function that accepts PageCapability
	testPageCapability := func(pc print.PageCapability) bool {
		result := pc.PrintDataInPageMode()
		return len(result) == 2
	}

	// Function that accepts ReverseCapability
	testReverseCapability := func(rc print.ReverseCapability) bool {
		_, err := rc.PrintAndReverseFeed(10)
		return err == nil
	}

	// Function that accepts PageModeCapability
	testPageModeCapability := func(pmc print.PageModeCapability) bool {
		result := pmc.PrintDataInPageMode()
		_, err := pmc.PrintAndReverseFeed(5)
		return len(result) == 2 && err == nil
	}

	// Same struct works with all interface types
	if !testPageCapability(pp) {
		t.Error("PagePrint should work as PageCapability")
	}

	if !testReverseCapability(pp) {
		t.Error("PagePrint should work as ReverseCapability")
	}

	if !testPageModeCapability(pp) {
		t.Error("PagePrint should work as PageModeCapability")
	}
}
