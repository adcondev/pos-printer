package print_test

import (
	"bytes"
	"errors"
	"testing"

	common2 "github.com/adcondev/pos-printer/pkg/commands/common"
	print2 "github.com/adcondev/pos-printer/pkg/commands/print"
)

// ============================================================================
// Utility Functions Tests
// ============================================================================

// Naming Convention: TestUtility_{Function}_{Optional Scenario}

func TestUtility_Formatting_CharacterReplacement(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want []byte
	}{
		{
			name: "replaces newline with LF",
			data: []byte("Hello\nWorld"),
			want: []byte{'H', 'e', 'l', 'l', 'o', print2.LF, 'W', 'o', 'r', 'l', 'd'},
		},
		{
			name: "replaces tab with HT",
			data: []byte("Col1\tCol2"),
			want: []byte{'C', 'o', 'l', '1', common2.HT, 'C', 'o', 'l', '2'},
		},
		{
			name: "replaces carriage return with CR",
			data: []byte("Line1\rLine2"),
			want: []byte{'L', 'i', 'n', 'e', '1', print2.CR, 'L', 'i', 'n', 'e', '2'},
		},
		{
			name: "handles multiple replacements",
			data: []byte("A\nB\tC\rD"),
			want: []byte{'A', print2.LF, 'B', common2.HT, 'C', print2.CR, 'D'},
		},
		{
			name: "preserves regular characters",
			data: []byte("NoSpecialChars123!@#"),
			want: []byte("NoSpecialChars123!@#"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clone input to avoid modifying original
			data := append([]byte(nil), tt.data...)
			got := print2.Formatting(data)
			if !bytes.Equal(got, tt.want) {
				t.Errorf("Formatting(%q) = %#v, want %#v", tt.data, got, tt.want)
			}
		})
	}
}

// ============================================================================
// Commands Tests
// ============================================================================

// Naming Convention: Test{Struct}_{Method}_{Optional Scenario}

func TestCommands_Text(t *testing.T) {
	cmd := print2.NewCommands()

	tests := []struct {
		name    string
		text    string
		want    []byte
		wantErr bool
	}{
		{
			name:    "simple text",
			text:    "Hello",
			want:    []byte("Hello"),
			wantErr: false,
		},
		{
			name:    "text with newline",
			text:    "Line1\nLine2",
			want:    []byte{'L', 'i', 'n', 'e', '1', print2.LF, 'L', 'i', 'n', 'e', '2'},
			wantErr: false,
		},
		{
			name:    "text with tab",
			text:    "A\tB",
			want:    []byte{'A', common2.HT, 'B'},
			wantErr: false,
		},
		{
			name:    "text with carriage return",
			text:    "A\rB",
			want:    []byte{'A', print2.CR, 'B'},
			wantErr: false,
		},
		{
			name:    "empty buffer",
			text:    "",
			want:    nil,
			wantErr: true,
		},
		{
			name: "buffer overflow",
			// FIXME: change anonymous func to utils helpers
			text: func() string {
				overflow := make([]byte, common2.MaxBuf+1)
				return string(overflow)
			}(),
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.Text(tt.text)

			// Standardized error checking
			if (err != nil) != tt.wantErr {
				t.Errorf("Text(%s) error = %v, wantErr %v", tt.text, err, tt.wantErr)
				return
			}

			var baseErr error
			switch tt.name {
			case "empty buffer":
				baseErr = print2.ErrEmptyText
			case "buffer overflow":
				baseErr = print2.ErrTextTooLarge
			default:
				baseErr = nil
			}

			// Check specific error type if expecting error
			if tt.wantErr && err != nil {
				if !errors.Is(err, baseErr) {
					t.Errorf("Text(%s) error = %v, want %v", tt.text, err, baseErr)
				}
				if !errors.Is(err, baseErr) {
					t.Errorf("Text(%s) error = %v, want %v", tt.text, err, baseErr)
				}
				return
			}

			// Check result if no error expected
			if !tt.wantErr && !bytes.Equal(got, tt.want) {
				t.Errorf("Text(%s) = %#v, want %#v", tt.text, got, tt.want)
			}
		})
	}
}

func TestCommands_PrintAndFeedPaper(t *testing.T) {
	cmd := print2.NewCommands()

	tests := []struct {
		name  string
		units byte
		want  []byte
	}{
		{
			name:  "minimum feed (0 units)",
			units: 0,
			want:  []byte{common2.ESC, 'J', 0},
		},
		{
			name:  "typical feed (30 units)",
			units: 30,
			want:  []byte{common2.ESC, 'J', 30},
		},
		{
			name:  "maximum feed (255 units)",
			units: 255,
			want:  []byte{common2.ESC, 'J', 255},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cmd.PrintAndFeedPaper(tt.units)
			if !bytes.Equal(got, tt.want) {
				t.Errorf("PrintAndFeedPaper(%d) = %#v, want %#v", tt.units, got, tt.want)
			}
		})
	}
}

func TestCommands_FormFeed(t *testing.T) {
	cmd := print2.NewCommands()
	got := cmd.FormFeed()
	want := []byte{print2.FF}

	if !bytes.Equal(got, want) {
		t.Errorf("FormFeed() = %#v, want %#v", got, want)
	}
}

func TestCommands_PrintAndCarriageReturn(t *testing.T) {
	cmd := print2.NewCommands()
	got := cmd.PrintAndCarriageReturn()
	want := []byte{print2.CR}

	if !bytes.Equal(got, want) {
		t.Errorf("PrintAndCarriageReturn() = %#v, want %#v", got, want)
	}
}

func TestCommands_PrintAndLineFeed(t *testing.T) {
	cmd := print2.NewCommands()
	got := cmd.PrintAndLineFeed()
	want := []byte{print2.LF}

	if !bytes.Equal(got, want) {
		t.Errorf("PrintAndLineFeed() = %#v, want %#v", got, want)
	}
}

// ============================================================================
// PagePrint Tests
// ============================================================================

func TestCommands_PrintDataInPageMode(t *testing.T) {
	pp := print2.NewCommands()
	got := pp.PrintDataInPageMode()
	want := []byte{common2.ESC, print2.FF}

	if !bytes.Equal(got, want) {
		t.Errorf("PrintDataInPageMode() = %#v, want %#v", got, want)
	}
}

func TestCommands_CancelData(t *testing.T) {
	pp := print2.NewCommands()
	got := pp.CancelData()
	want := []byte{print2.CAN}

	if !bytes.Equal(got, want) {
		t.Errorf("CancelData() = %#v, want %#v", got, want)
	}
}

func TestCommands_PrintAndReverseFeed(t *testing.T) {
	pp := print2.NewCommands()

	tests := []struct {
		name    string
		reverse byte
		want    []byte
		wantErr bool
	}{
		{
			name:    "minimum reverse feed (0 units)",
			reverse: 0,
			want:    []byte{common2.ESC, 'K', 0},
			wantErr: false,
		},
		{
			name:    "typical reverse feed (10 units)",
			reverse: 10,
			want:    []byte{common2.ESC, 'K', 10},
			wantErr: false,
		},
		{
			name:    "maximum allowed reverse feed",
			reverse: print2.MaxReverseMotionUnits,
			want:    []byte{common2.ESC, 'K', print2.MaxReverseMotionUnits},
			wantErr: false,
		},
		{
			name:    "exceeds maximum returns error",
			reverse: print2.MaxReverseMotionUnits + 1,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pp.PrintAndReverseFeed(tt.reverse)

			// Standardized error checking
			if (err != nil) != tt.wantErr {
				t.Errorf("PrintAndReverseFeed(%d) error = %v, wantErr %v", tt.reverse, err, tt.wantErr)
				return
			}

			// Check specific error type if expecting error
			if tt.wantErr && err != nil {
				if !errors.Is(err, print2.ErrReverseUnits) {
					t.Errorf("PrintAndReverseFeed(%v) error = %v, want %v",
						tt.reverse, err, print2.ErrReverseUnits)
				}
				return
			}

			// Check result if no error expected
			if !tt.wantErr && !bytes.Equal(got, tt.want) {
				t.Errorf("PrintAndReverseFeed(%v) = %#v, want %#v", tt.reverse, got, tt.want)
			}
		})
	}
}

func TestCommands_PrintAndReverseFeedLines(t *testing.T) {
	pp := print2.NewCommands()

	tests := []struct {
		name    string
		lines   byte
		want    []byte
		wantErr bool
	}{
		{
			name:    "minimum lines (0)",
			lines:   0,
			want:    []byte{common2.ESC, 'e', 0},
			wantErr: false,
		},
		{
			name:    "single line reverse",
			lines:   1,
			want:    []byte{common2.ESC, 'e', 1},
			wantErr: false,
		},
		{
			name:    "maximum allowed lines",
			lines:   print2.MaxReverseFeedLines,
			want:    []byte{common2.ESC, 'e', print2.MaxReverseFeedLines},
			wantErr: false,
		},
		{
			name:    "exceeds maximum returns error",
			lines:   print2.MaxReverseFeedLines + 1,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pp.PrintAndReverseFeedLines(tt.lines)

			// Standardized error checking
			if (err != nil) != tt.wantErr {
				t.Errorf("PrintAndReverseFeedLines(%d) error = %v, wantErr %v", tt.lines, err, tt.wantErr)
				return
			}

			// Check specific error type if expecting error
			if tt.wantErr && err != nil {
				if !errors.Is(err, print2.ErrReverseLines) {
					t.Errorf("PrintAndReverseFeedLines(%v) error = %v, want %v",
						tt.lines, err, print2.ErrReverseLines)
				}
				return
			}

			// Check result if no error expected
			if !tt.wantErr && !bytes.Equal(got, tt.want) {
				t.Errorf("PrintAndReverseFeedLines(%v) = %#v, want %#v", tt.lines, got, tt.want)
			}
		})
	}
}

func TestPagePrint_PrintAndFeedLines(t *testing.T) {
	pp := print2.NewCommands()

	tests := []struct {
		name    string
		lines   byte
		want    []byte
		wantErr bool
	}{
		{
			name:  "no feed (0 lines)",
			lines: 0,
			want:  []byte{common2.ESC, 'd', 0},
		},
		{
			name:  "typical feed (5 lines)",
			lines: 5,
			want:  []byte{common2.ESC, 'd', 5},
		},
		{
			name:  "maximum feed (255 lines)",
			lines: 255,
			want:  []byte{common2.ESC, 'd', 255},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pp.PrintAndFeedLines(tt.lines)

			// Check result if no error expected
			if !tt.wantErr && !bytes.Equal(got, tt.want) {
				t.Errorf("PrintAndFeedLines(%v) = %#v, want %#v", tt.lines, got, tt.want)
			}
		})
	}
}
