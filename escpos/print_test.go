package escpos

import (
	"bytes"
	"errors"
	"testing"
)

// Naming Convention: Test{Struct}_{Method}_{Scenario}

// ============================================================================
// PrintCommands Tests
// ============================================================================

func TestPrintCommands_Text_ValidInput(t *testing.T) {
	pc := &PrintCommands{}

	tests := []struct {
		name    string
		input   string
		want    []byte
		wantErr bool
	}{
		{
			name:    "simple text",
			input:   "Hello",
			want:    []byte("Hello"),
			wantErr: false,
		},
		{
			name:    "text with control characters",
			input:   "A\nB\tC\rD",
			want:    []byte{'A', LF, 'B', HT, 'C', CR, 'D'},
			wantErr: false,
		},
		{
			name:    "empty text returns error",
			input:   "",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pc.Text(tt.input)

			// Check error state
			if (err != nil) != tt.wantErr {
				t.Errorf("PrintCommands.Text(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}

			// Check result only if no error expected
			if !tt.wantErr && !bytes.Equal(got, tt.want) {
				t.Errorf("PrintCommands.Text(%q) = %#v, want %#v", tt.input, got, tt.want)
			}
		})
	}
}

func TestPrintCommands_PrintAndFeedPaper_ByteSequence(t *testing.T) {
	pc := &PrintCommands{}

	tests := []struct {
		name string
		n    byte
		want []byte
	}{
		{
			name: "minimum feed (0 units)",
			n:    0,
			want: []byte{ESC, 'J', 0},
		},
		{
			name: "typical feed (5 units)",
			n:    5,
			want: []byte{ESC, 'J', 5},
		},
		{
			name: "maximum feed (255 units)",
			n:    255,
			want: []byte{ESC, 'J', 255},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pc.PrintAndFeedPaper(tt.n)
			if !bytes.Equal(got, tt.want) {
				t.Errorf("PrintCommands.PrintAndFeedPaper(%d) = %#v, want %#v", tt.n, got, tt.want)
			}
		})
	}
}

func TestPrintCommands_FormFeed_SingleByte(t *testing.T) {
	pc := &PrintCommands{}

	got := pc.FormFeed()
	want := []byte{FF}

	if !bytes.Equal(got, want) {
		t.Errorf("PrintCommands.FormFeed() = %#v, want %#v", got, want)
	}
}

func TestPrintCommands_PrintAndCarriageReturn_SingleByte(t *testing.T) {
	pc := &PrintCommands{}

	got := pc.PrintAndCarriageReturn()
	want := []byte{CR}

	if !bytes.Equal(got, want) {
		t.Errorf("PrintCommands.PrintAndCarriageReturn() = %#v, want %#v", got, want)
	}
}

func TestPrintCommands_PrintAndLineFeed_SingleByte(t *testing.T) {
	pc := &PrintCommands{}

	got := pc.PrintAndLineFeed()
	want := []byte{LF}

	if !bytes.Equal(got, want) {
		t.Errorf("PrintCommands.PrintAndLineFeed() = %#v, want %#v", got, want)
	}
}

// ============================================================================
// PagePrint Tests
// ============================================================================

func TestPagePrint_PrintDataInPageMode_ByteSequence(t *testing.T) {
	pp := &PagePrint{}

	got := pp.PrintDataInPageMode()
	want := []byte{ESC, FF}

	if !bytes.Equal(got, want) {
		t.Errorf("PagePrint.PrintDataInPageMode() = %#v, want %#v", got, want)
	}
}

func TestPagePrint_PrintAndReverseFeed_Validation(t *testing.T) {
	pp := &PagePrint{}

	t.Run("valid range", func(t *testing.T) {
		tests := []struct {
			name string
			n    byte
			want []byte
		}{
			{
				name: "minimum reverse feed (0 units)",
				n:    0,
				want: []byte{ESC, 'K', 0},
			},
			{
				name: "typical reverse feed (10 units)",
				n:    10,
				want: []byte{ESC, 'K', 10},
			},
			{
				name: "maximum allowed reverse feed",
				n:    MaxReverseMotionUnits,
				want: []byte{ESC, 'K', MaxReverseMotionUnits},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := pp.PrintAndReverseFeed(tt.n)

				if err != nil {
					t.Errorf("PagePrint.PrintAndReverseFeed(%d) unexpected error: %v", tt.n, err)
				}
				if !bytes.Equal(got, tt.want) {
					t.Errorf("PagePrint.PrintAndReverseFeed(%d) = %#v, want %#v", tt.n, got, tt.want)
				}
			})
		}
	})

	t.Run("overflow error", func(t *testing.T) {
		n := MaxReverseMotionUnits + 1
		_, err := pp.PrintAndReverseFeed(n)

		if !errors.Is(err, errPrintReverseFeed) {
			t.Errorf("PagePrint.PrintAndReverseFeed(%d) error = %v, want %v", n, err, errPrintReverseFeed)
		}
	})
}

func TestPagePrint_PrintAndReverseFeedLines_Validation(t *testing.T) {
	pp := &PagePrint{}

	t.Run("valid range", func(t *testing.T) {
		tests := []struct {
			name string
			n    byte
			want []byte
		}{
			{
				name: "minimum lines (0)",
				n:    0,
				want: []byte{ESC, 'e', 0},
			},
			{
				name: "single line reverse",
				n:    1,
				want: []byte{ESC, 'e', 1},
			},
			{
				name: "maximum allowed lines",
				n:    MaxReverseFeedLines,
				want: []byte{ESC, 'e', MaxReverseFeedLines},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := pp.PrintAndReverseFeedLines(tt.n)

				if err != nil {
					t.Errorf("PagePrint.PrintAndReverseFeedLines(%d) unexpected error: %v", tt.n, err)
				}
				if !bytes.Equal(got, tt.want) {
					t.Errorf("PagePrint.PrintAndReverseFeedLines(%d) = %#v, want %#v", tt.n, got, tt.want)
				}
			})
		}
	})

	t.Run("overflow error", func(t *testing.T) {
		n := MaxReverseFeedLines + 1
		_, err := pp.PrintAndReverseFeedLines(n)

		if !errors.Is(err, errPrintReverseFeedLines) {
			t.Errorf("PagePrint.PrintAndReverseFeedLines(%d) error = %v, want %v", n, err, errPrintReverseFeedLines)
		}
	})
}

func TestPagePrint_PrintAndFeedLines_ByteSequence(t *testing.T) {
	pp := &PagePrint{}

	tests := []struct {
		name string
		n    byte
		want []byte
	}{
		{
			name: "no feed (0 lines)",
			n:    0,
			want: []byte{ESC, 'd', 0},
		},
		{
			name: "typical feed (7 lines)",
			n:    7,
			want: []byte{ESC, 'd', 7},
		},
		{
			name: "maximum feed (255 lines)",
			n:    255,
			want: []byte{ESC, 'd', 255},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pp.PrintAndFeedLines(tt.n)

			if err != nil {
				t.Errorf("PagePrint.PrintAndFeedLines(%d) unexpected error: %v", tt.n, err)
			}
			if !bytes.Equal(got, tt.want) {
				t.Errorf("PagePrint.PrintAndFeedLines(%d) = %#v, want %#v", tt.n, got, tt.want)
			}
		})
	}
}
