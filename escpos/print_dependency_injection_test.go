package escpos

import (
	"bytes"
	"testing"
)

// Naming Convention: Test{Pattern}_{Function}_{Scenario}

// ============================================================================
// Helper Functions for Dependency Injection Testing
// ============================================================================

// formatDocument demonstrates dependency injection by accepting any PrinterCapability
func formatDocument(printer PrinterCapability, title, body string) ([]byte, error) {
	var output []byte

	// Print title
	titleData, err := printer.Text(title)
	if err != nil {
		return nil, err
	}
	output = append(output, titleData...)
	output = append(output, printer.PrintAndLineFeed()...)

	// Print body
	bodyData, err := printer.Text(body)
	if err != nil {
		return nil, err
	}
	output = append(output, bodyData...)
	output = append(output, printer.FormFeed()...)

	return output, nil
}

// ============================================================================
// Dependency Injection Tests
// ============================================================================

func TestDependencyInjection_FormatDocument_RealImplementation(t *testing.T) {
	// Create real implementation
	realPrinter := &PrintCommands{Page: &PagePrint{}}

	output, err := formatDocument(realPrinter, "Title", "Body")

	if err != nil {
		t.Fatalf("formatDocument() unexpected error: %v", err)
	}

	// Verify output structure
	expectedLen := len("Title") + 1 + len("Body") + 1 // +1 for LF, +1 for FF
	if len(output) != expectedLen {
		t.Errorf("formatDocument() output length = %d, want %d", len(output), expectedLen)
	}

	// Verify content
	if !bytes.Contains(output, []byte("Title")) {
		t.Error("formatDocument() output should contain 'Title'")
	}
	if !bytes.Contains(output, []byte("Body")) {
		t.Error("formatDocument() output should contain 'Body'")
	}
	if !bytes.Contains(output, []byte{LF}) {
		t.Error("formatDocument() output should contain line feed")
	}
	if !bytes.Contains(output, []byte{FF}) {
		t.Error("formatDocument() output should contain form feed")
	}
}

func TestDependencyInjection_FormatDocument_MockImplementation(t *testing.T) {
	// Create mock with predefined responses
	mock := &MockPrinterCapability{
		TextReturn: []byte("processed"),
	}

	_, err := formatDocument(mock, "Title", "Body")

	if err != nil {
		t.Fatalf("formatDocument() unexpected error: %v", err)
	}

	// Verify that all required methods were called
	if !mock.TextCalled {
		t.Error("formatDocument() should call Text()")
	}
	if !mock.PrintAndLineFeedCalled {
		t.Error("formatDocument() should call PrintAndLineFeed()")
	}
	if !mock.FormFeedCalled {
		t.Error("formatDocument() should call FormFeed()")
	}

	// Verify Text was called twice (for title and body)
	// Note: This would require a call counter in the mock
}

func TestDependencyInjection_SwappableImplementations(t *testing.T) {
	// This test demonstrates that the same function works with different implementations

	testCases := []struct {
		name    string
		printer PrinterCapability
		wantErr bool
	}{
		{
			name:    "real PrintCommands implementation",
			printer: &PrintCommands{Page: &PagePrint{}},
			wantErr: false,
		},
		{
			name:    "mock implementation with success",
			printer: &MockPrinterCapability{TextReturn: []byte("test")},
			wantErr: false,
		},
		{
			name:    "fake implementation",
			printer: NewFakePrinter(),
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := formatDocument(tc.printer, "Test", "Content")

			if (err != nil) != tc.wantErr {
				t.Errorf("formatDocument() with %s error = %v, wantErr %v",
					tc.name, err, tc.wantErr)
			}
		})
	}
}
