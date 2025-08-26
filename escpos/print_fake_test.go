package escpos

import (
	"bytes"
	"testing"
)

// ============================================================================
// Fake Implementation
// ============================================================================

// FakePrinter simulates a real printer with internal state tracking
type FakePrinter struct {
	buffer   []byte
	position int
}

// NewFakePrinter creates a new fake printer instance
func NewFakePrinter() *FakePrinter {
	return &FakePrinter{
		buffer:   make([]byte, 0),
		position: 0,
	}
}

func (fp *FakePrinter) Text(n string) ([]byte, error) {
	if n == "" {
		return nil, errEmptyBuffer
	}
	data := format([]byte(n))
	fp.buffer = append(fp.buffer, data...)
	fp.position += len(data)
	return data, nil
}

func (fp *FakePrinter) PrintAndFeedPaper(n byte) []byte {
	cmd := []byte{ESC, 'J', n}
	fp.buffer = append(fp.buffer, cmd...)
	fp.position = 0 // Reset position after feed
	return cmd
}

func (fp *FakePrinter) FormFeed() []byte {
	fp.buffer = append(fp.buffer, FF)
	fp.position = 0 // Reset position after form feed
	return []byte{FF}
}

func (fp *FakePrinter) PrintAndCarriageReturn() []byte {
	fp.buffer = append(fp.buffer, CR)
	fp.position = 0 // Reset to beginning of line
	return []byte{CR}
}

func (fp *FakePrinter) PrintAndLineFeed() []byte {
	fp.buffer = append(fp.buffer, LF)
	fp.position = 0 // Reset to beginning of next line
	return []byte{LF}
}

// GetBuffer returns the accumulated buffer for verification
func (fp *FakePrinter) GetBuffer() []byte {
	return fp.buffer
}

// GetPosition returns current print position for testing
func (fp *FakePrinter) GetPosition() int {
	return fp.position
}

// ============================================================================
// Fake Implementation Tests
// ============================================================================

func TestFakePrinter_Text_StateManagement(t *testing.T) {
	fake := NewFakePrinter()

	// Test text accumulation
	text1 := "Hello"
	result1, err := fake.Text(text1)
	if err != nil {
		t.Fatalf("FakePrinter.Text(%q) unexpected error: %v", text1, err)
	}

	// Verify return value
	if !bytes.Equal(result1, []byte(text1)) {
		t.Errorf("FakePrinter.Text(%q) = %#v, want %#v", text1, result1, []byte(text1))
	}

	// Verify buffer accumulation
	if !bytes.Equal(fake.GetBuffer(), []byte(text1)) {
		t.Errorf("FakePrinter buffer after Text(%q) = %#v, want %#v",
			text1, fake.GetBuffer(), []byte(text1))
	}

	// Verify position tracking
	if fake.GetPosition() != len(text1) {
		t.Errorf("FakePrinter position = %d, want %d", fake.GetPosition(), len(text1))
	}
}

func TestFakePrinter_Integration_PrintSequence(t *testing.T) {
	fake := NewFakePrinter()
	cmd := &Commands{
		Print: fake,
	}

	// Simulate printing a receipt
	receiptHeader := "RECEIPT"
	_, err := cmd.Print.Text(receiptHeader)
	if err != nil {
		t.Fatalf("Commands.Print.Text(%q) unexpected error: %v", receiptHeader, err)
	}

	cmd.Print.PrintAndLineFeed()

	itemLine := "Item 1: $10.00"
	_, err = cmd.Print.Text(itemLine)
	if err != nil {
		t.Fatalf("Commands.Print.Text(%q) unexpected error: %v", itemLine, err)
	}

	cmd.Print.PrintAndLineFeed()
	cmd.Print.FormFeed()

	// Verify the complete buffer
	buffer := fake.GetBuffer()

	// Check content presence
	expectedParts := []string{receiptHeader, itemLine}
	for _, part := range expectedParts {
		if !bytes.Contains(buffer, []byte(part)) {
			t.Errorf("FakePrinter buffer should contain %q", part)
		}
	}

	// Check control characters count
	lfCount := bytes.Count(buffer, []byte{LF})
	if lfCount != 2 {
		t.Errorf("FakePrinter buffer LF count = %d, want 2", lfCount)
	}

	ffCount := bytes.Count(buffer, []byte{FF})
	if ffCount != 1 {
		t.Errorf("FakePrinter buffer FF count = %d, want 1", ffCount)
	}

	// Verify final position is reset
	if fake.GetPosition() != 0 {
		t.Errorf("FakePrinter final position = %d, want 0 (after FormFeed)",
			fake.GetPosition())
	}
}
