package print_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/adcondev/pos-printer/escpos/common"
	"github.com/adcondev/pos-printer/escpos/print"
)

// ============================================================================
// PagePrint Tests
// ============================================================================

func TestPagePrint_PrintDataInPageMode_ByteSequence(t *testing.T) {
	pp := &print.PagePrint{}

	got := pp.PrintDataInPageMode()
	want := []byte{common.ESC, print.FF}

	if !bytes.Equal(got, want) {
		t.Errorf("PagePrint.PrintDataInPageMode() = %#v, want %#v", got, want)
	}
}

func TestPagePrint_CancelData_SingleByte(t *testing.T) {
	pp := &print.PagePrint{}

	got := pp.CancelData()
	want := []byte{print.CAN}

	if !bytes.Equal(got, want) {
		t.Errorf("PagePrint.PrintDataInPageMode() = %#v, want %#v", got, want)
	}
}

func TestPagePrint_PrintAndReverseFeed_Validation(t *testing.T) {
	pp := &print.PagePrint{}

	t.Run("valid range", func(t *testing.T) {
		tests := []struct {
			name string
			n    byte
			want []byte
		}{
			{
				name: "minimum reverse feed (0 units)",
				n:    0,
				want: []byte{common.ESC, 'K', 0},
			},
			{
				name: "typical reverse feed (10 units)",
				n:    10,
				want: []byte{common.ESC, 'K', 10},
			},
			{
				name: "maximum allowed reverse feed",
				n:    common.MaxReverseMotionUnits,
				want: []byte{common.ESC, 'K', common.MaxReverseMotionUnits},
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
		n := common.MaxReverseMotionUnits + 1
		_, err := pp.PrintAndReverseFeed(n)

		if !errors.Is(err, common.ErrPrintReverseFeed) {
			t.Errorf("PagePrint.PrintAndReverseFeed(%d) error = %v, want %v", n, err, common.ErrPrintReverseFeed)
		}
	})
}

func TestPagePrint_PrintAndReverseFeedLines_Validation(t *testing.T) {
	pp := &print.PagePrint{}

	t.Run("valid range", func(t *testing.T) {
		tests := []struct {
			name string
			n    byte
			want []byte
		}{
			{
				name: "minimum lines (0)",
				n:    0,
				want: []byte{common.ESC, 'e', 0},
			},
			{
				name: "single line reverse",
				n:    1,
				want: []byte{common.ESC, 'e', 1},
			},
			{
				name: "maximum allowed lines",
				n:    common.MaxReverseFeedLines,
				want: []byte{common.ESC, 'e', common.MaxReverseFeedLines},
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
		n := common.MaxReverseFeedLines + 1
		_, err := pp.PrintAndReverseFeedLines(n)

		if !errors.Is(err, common.ErrPrintReverseFeedLines) {
			t.Errorf("PagePrint.PrintAndReverseFeedLines(%d) error = %v, want %v", n, err, common.ErrPrintReverseFeedLines)
		}
	})
}

func TestPagePrint_PrintAndFeedLines_ByteSequence(t *testing.T) {
	pp := &print.PagePrint{}

	tests := []struct {
		name string
		n    byte
		want []byte
	}{
		{
			name: "no feed (0 lines)",
			n:    0,
			want: []byte{common.ESC, 'd', 0},
		},
		{
			name: "typical feed (7 lines)",
			n:    7,
			want: []byte{common.ESC, 'd', 7},
		},
		{
			name: "maximum feed (255 lines)",
			n:    255,
			want: []byte{common.ESC, 'd', 255},
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
