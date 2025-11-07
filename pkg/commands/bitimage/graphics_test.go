package bitimage_test

import (
	"testing"

	"github.com/adcondev/pos-printer/internal/testutils"
	bitimage2 "github.com/adcondev/pos-printer/pkg/commands/bitimage"
	"github.com/adcondev/pos-printer/pkg/commands/common"
)

// ============================================================================
// Graphics Commands Tests
// ============================================================================

func TestGraphicsCommands_SetGraphicsDotDensity(t *testing.T) {
	cmd := bitimage2.NewGraphicsCommands()

	tests := []struct {
		name    string
		fn      bitimage2.FunctionCode
		x       bitimage2.DotDensity
		y       bitimage2.DotDensity
		want    []byte
		wantErr error
	}{
		{
			name:    "180x180 dpi with function code 1",
			fn:      bitimage2.FunctionCodeDensity1,
			x:       bitimage2.Density180x180,
			y:       bitimage2.Density180x180,
			want:    []byte{common.GS, '(', 'L', 0x04, 0x00, 0x30, 1, 50, 50},
			wantErr: nil,
		},
		{
			name:    "180x180 dpi with function code 49",
			fn:      bitimage2.FunctionCodeDensity49,
			x:       bitimage2.Density180x180,
			y:       bitimage2.Density180x180,
			want:    []byte{common.GS, '(', 'L', 0x04, 0x00, 0x30, 49, 50, 50},
			wantErr: nil,
		},
		{
			name:    "360x360 dpi with function code 1",
			fn:      bitimage2.FunctionCodeDensity1,
			x:       bitimage2.Density360x360,
			y:       bitimage2.Density360x360,
			want:    []byte{common.GS, '(', 'L', 0x04, 0x00, 0x30, 1, 51, 51},
			wantErr: nil,
		},
		{
			name:    "360x360 dpi with function code 49",
			fn:      bitimage2.FunctionCodeDensity49,
			x:       bitimage2.Density360x360,
			y:       bitimage2.Density360x360,
			want:    []byte{common.GS, '(', 'L', 0x04, 0x00, 0x30, 49, 51, 51},
			wantErr: nil,
		},
		{
			name:    "invalid function code 2",
			fn:      bitimage2.FunctionCodePrint2,
			x:       bitimage2.Density180x180,
			y:       bitimage2.Density180x180,
			want:    nil,
			wantErr: bitimage2.ErrInvalidFunctionCode,
		},
		{
			name:    "invalid function code 50",
			fn:      bitimage2.FunctionCodePrint50,
			x:       bitimage2.Density180x180,
			y:       bitimage2.Density180x180,
			want:    nil,
			wantErr: bitimage2.ErrInvalidFunctionCode,
		},
		{
			name:    "invalid density value",
			fn:      bitimage2.FunctionCodeDensity1,
			x:       52,
			y:       52,
			want:    nil,
			wantErr: bitimage2.ErrInvalidDensityValue,
		},
		{
			name:    "mismatched density values",
			fn:      bitimage2.FunctionCodeDensity1,
			x:       bitimage2.Density180x180,
			y:       bitimage2.Density360x360,
			want:    nil,
			wantErr: bitimage2.ErrInvalidDensityCombination,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.SetGraphicsDotDensity(tt.fn, tt.x, tt.y)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "SetGraphicsDotDensity") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			testutils.AssertBytes(t, got, tt.want, "SetGraphicsDotDensity(%v, %v, %v)", tt.fn, tt.x, tt.y)
		})
	}
}

func TestGraphicsCommands_PrintBufferedGraphics(t *testing.T) {
	cmd := bitimage2.NewGraphicsCommands()

	tests := []struct {
		name    string
		fn      bitimage2.FunctionCode
		want    []byte
		wantErr error
	}{
		{
			name:    "function code 2",
			fn:      bitimage2.FunctionCodePrint2,
			want:    []byte{common.GS, '(', 'L', 0x02, 0x00, 0x30, 2},
			wantErr: nil,
		},
		{
			name:    "function code 50",
			fn:      bitimage2.FunctionCodePrint50,
			want:    []byte{common.GS, '(', 'L', 0x02, 0x00, 0x30, 50},
			wantErr: nil,
		},
		{
			name:    "invalid function code 1",
			fn:      bitimage2.FunctionCodeDensity1,
			want:    nil,
			wantErr: bitimage2.ErrInvalidFunctionCode,
		},
		{
			name:    "invalid function code 49",
			fn:      bitimage2.FunctionCodeDensity49,
			want:    nil,
			wantErr: bitimage2.ErrInvalidFunctionCode,
		},
		{
			name:    "invalid function code 99",
			fn:      99,
			want:    nil,
			wantErr: bitimage2.ErrInvalidFunctionCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.PrintBufferedGraphics(tt.fn)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "PrintBufferedGraphics") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			testutils.AssertBytes(t, got, tt.want, "PrintBufferedGraphics(%v)", tt.fn)
		})
	}
}

func TestGraphicsCommands_StoreRasterGraphicsInBuffer(t *testing.T) {
	cmd := bitimage2.NewGraphicsCommands()

	// Helper to create testutils data for raster format
	createRasterData := func(width, height uint16) []byte {
		widthBytes := (int(width) + 7) / 8
		return testutils.RepeatByte(widthBytes*int(height), 0xFF)
	}

	tests := []struct {
		name            string
		tone            bitimage2.GraphicsTone
		horizontalScale bitimage2.GraphicsScale
		verticalScale   bitimage2.GraphicsScale
		color           bitimage2.GraphicsColor
		width           uint16
		height          uint16
		data            []byte
		wantErr         error
	}{
		{
			name:            "monochrome normal scale",
			tone:            bitimage2.Monochrome,
			horizontalScale: bitimage2.NormalScale,
			verticalScale:   bitimage2.NormalScale,
			color:           bitimage2.Color1,
			width:           100,
			height:          50,
			data:            createRasterData(100, 50),
			wantErr:         nil,
		},
		{
			name:            "monochrome double width",
			tone:            bitimage2.Monochrome,
			horizontalScale: bitimage2.DoubleScale,
			verticalScale:   bitimage2.NormalScale,
			color:           bitimage2.Color2,
			width:           200,
			height:          100,
			data:            createRasterData(200, 100),
			wantErr:         nil,
		},
		{
			name:            "monochrome double height",
			tone:            bitimage2.Monochrome,
			horizontalScale: bitimage2.NormalScale,
			verticalScale:   bitimage2.DoubleScale,
			color:           bitimage2.Color3,
			width:           150,
			height:          1200,
			data:            createRasterData(150, 1200),
			wantErr:         nil,
		},
		{
			name:            "multiple tone normal scale",
			tone:            bitimage2.MultipleTone,
			horizontalScale: bitimage2.NormalScale,
			verticalScale:   bitimage2.NormalScale,
			color:           bitimage2.Color4,
			width:           100,
			height:          600,
			data:            createRasterData(100, 600),
			wantErr:         nil,
		},
		{
			name:            "multiple tone double scale",
			tone:            bitimage2.MultipleTone,
			horizontalScale: bitimage2.DoubleScale,
			verticalScale:   bitimage2.DoubleScale,
			color:           bitimage2.Color1,
			width:           100,
			height:          300,
			data:            createRasterData(100, 300),
			wantErr:         nil,
		},
		{
			name:            "maximum width",
			tone:            bitimage2.Monochrome,
			horizontalScale: bitimage2.NormalScale,
			verticalScale:   bitimage2.NormalScale,
			color:           bitimage2.Color1,
			width:           2400,
			height:          10,
			data:            createRasterData(2400, 10),
			wantErr:         nil,
		},
		{
			name:            "invalid tone",
			tone:            49,
			horizontalScale: bitimage2.NormalScale,
			verticalScale:   bitimage2.NormalScale,
			color:           bitimage2.Color1,
			width:           100,
			height:          100,
			data:            createRasterData(100, 100),
			wantErr:         bitimage2.ErrInvalidTone,
		},
		{
			name:            "invalid horizontal scale",
			tone:            bitimage2.Monochrome,
			horizontalScale: 0,
			verticalScale:   bitimage2.NormalScale,
			color:           bitimage2.Color1,
			width:           100,
			height:          100,
			data:            createRasterData(100, 100),
			wantErr:         bitimage2.ErrInvalidScale,
		},
		{
			name:            "invalid vertical scale",
			tone:            bitimage2.Monochrome,
			horizontalScale: bitimage2.NormalScale,
			verticalScale:   3,
			color:           bitimage2.Color1,
			width:           100,
			height:          100,
			data:            createRasterData(100, 100),
			wantErr:         bitimage2.ErrInvalidScale,
		},
		{
			name:            "invalid color",
			tone:            bitimage2.Monochrome,
			horizontalScale: bitimage2.NormalScale,
			verticalScale:   bitimage2.NormalScale,
			color:           48,
			width:           100,
			height:          100,
			data:            createRasterData(100, 100),
			wantErr:         bitimage2.ErrInvalidColor,
		},
		{
			name:            "width exceeds limit",
			tone:            bitimage2.Monochrome,
			horizontalScale: bitimage2.NormalScale,
			verticalScale:   bitimage2.NormalScale,
			color:           bitimage2.Color1,
			width:           2401,
			height:          100,
			data:            createRasterData(2401, 100),
			wantErr:         bitimage2.ErrInvalidWidth,
		},
		{
			name:            "height exceeds limit for monochrome normal",
			tone:            bitimage2.Monochrome,
			horizontalScale: bitimage2.NormalScale,
			verticalScale:   bitimage2.NormalScale,
			color:           bitimage2.Color1,
			width:           100,
			height:          2401,
			data:            createRasterData(100, 2401),
			wantErr:         bitimage2.ErrInvalidHeight,
		},
		{
			name:            "height exceeds limit for monochrome double",
			tone:            bitimage2.Monochrome,
			horizontalScale: bitimage2.NormalScale,
			verticalScale:   bitimage2.DoubleScale,
			color:           bitimage2.Color1,
			width:           100,
			height:          1201,
			data:            createRasterData(100, 1201),
			wantErr:         bitimage2.ErrInvalidHeight,
		},
		{
			name:            "height exceeds limit for multiple tone normal",
			tone:            bitimage2.MultipleTone,
			horizontalScale: bitimage2.NormalScale,
			verticalScale:   bitimage2.NormalScale,
			color:           bitimage2.Color1,
			width:           100,
			height:          601,
			data:            createRasterData(100, 601),
			wantErr:         bitimage2.ErrInvalidHeight,
		},
		{
			name:            "height exceeds limit for multiple tone double",
			tone:            bitimage2.MultipleTone,
			horizontalScale: bitimage2.NormalScale,
			verticalScale:   bitimage2.DoubleScale,
			color:           bitimage2.Color1,
			width:           100,
			height:          301,
			data:            createRasterData(100, 301),
			wantErr:         bitimage2.ErrInvalidHeight,
		},
		{
			name:            "invalid data length",
			tone:            bitimage2.Monochrome,
			horizontalScale: bitimage2.NormalScale,
			verticalScale:   bitimage2.NormalScale,
			color:           bitimage2.Color1,
			width:           100,
			height:          100,
			data:            []byte{0xFF}, // Should be more bytes
			wantErr:         bitimage2.ErrInvalidDataLength,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.StoreRasterGraphicsInBuffer(tt.tone, tt.horizontalScale, tt.verticalScale,
				tt.color, tt.width, tt.height, tt.data)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "StoreRasterGraphicsInBuffer") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			// Verify the command structure
			if got[0] != common.GS || got[1] != '(' || got[2] != 'L' {
				t.Errorf("StoreRasterGraphicsInBuffer: invalid command prefix")
			}
		})
	}
}

func TestGraphicsCommands_StoreRasterGraphicsInBufferLarge(t *testing.T) {
	cmd := bitimage2.NewGraphicsCommands()

	// Helper to create testutils data for raster format
	createRasterData := func(width, height uint16) []byte {
		widthBytes := (int(width) + 7) / 8
		return testutils.RepeatByte(widthBytes*int(height), 0xFF)
	}

	tests := []struct {
		name            string
		tone            bitimage2.GraphicsTone
		horizontalScale bitimage2.GraphicsScale
		verticalScale   bitimage2.GraphicsScale
		color           bitimage2.GraphicsColor
		width           uint16
		height          uint16
		data            []byte
		wantErr         error
	}{
		{
			name:            "large monochrome data",
			tone:            bitimage2.Monochrome,
			horizontalScale: bitimage2.NormalScale,
			verticalScale:   bitimage2.NormalScale,
			color:           bitimage2.Color1,
			width:           2000,
			height:          2000,
			data:            createRasterData(2000, 2000),
			wantErr:         nil,
		},
		{
			name:            "invalid parameters same as standard",
			tone:            49,
			horizontalScale: bitimage2.NormalScale,
			verticalScale:   bitimage2.NormalScale,
			color:           bitimage2.Color1,
			width:           100,
			height:          100,
			data:            createRasterData(100, 100),
			wantErr:         bitimage2.ErrInvalidTone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.StoreRasterGraphicsInBufferLarge(tt.tone, tt.horizontalScale, tt.verticalScale,
				tt.color, tt.width, tt.height, tt.data)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "StoreRasterGraphicsInBufferLarge") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			// Verify the command structure for large format
			if got[0] != common.GS || got[1] != '8' || got[2] != 'L' {
				t.Errorf("StoreRasterGraphicsInBufferLarge: invalid command prefix")
			}
		})
	}
}

func TestGraphicsCommands_StoreColumnGraphicsInBuffer(t *testing.T) {
	cmd := bitimage2.NewGraphicsCommands()

	// Helper to create testutils data for column format
	createColumnData := func(width, height uint16) []byte {
		heightBytes := (int(height) + 7) / 8
		return testutils.RepeatByte(int(width)*heightBytes, 0xFF)
	}

	tests := []struct {
		name            string
		horizontalScale bitimage2.GraphicsScale
		verticalScale   bitimage2.GraphicsScale
		color           bitimage2.GraphicsColor
		width           uint16
		height          uint16
		data            []byte
		wantErr         error
	}{
		{
			name:            "normal scale color 1",
			horizontalScale: bitimage2.NormalScale,
			verticalScale:   bitimage2.NormalScale,
			color:           bitimage2.Color1,
			width:           100,
			height:          50,
			data:            createColumnData(100, 50),
			wantErr:         nil,
		},
		{
			name:            "double width color 2",
			horizontalScale: bitimage2.DoubleScale,
			verticalScale:   bitimage2.NormalScale,
			color:           bitimage2.Color2,
			width:           200,
			height:          100,
			data:            createColumnData(200, 100),
			wantErr:         nil,
		},
		{
			name:            "double height color 3",
			horizontalScale: bitimage2.NormalScale,
			verticalScale:   bitimage2.DoubleScale,
			color:           bitimage2.Color3,
			width:           150,
			height:          128,
			data:            createColumnData(150, 128),
			wantErr:         nil,
		},
		{
			name:            "maximum dimensions",
			horizontalScale: bitimage2.NormalScale,
			verticalScale:   bitimage2.NormalScale,
			color:           bitimage2.Color1,
			width:           2048,
			height:          128,
			data:            createColumnData(2048, 128),
			wantErr:         nil,
		},
		{
			name:            "color 4 not supported",
			horizontalScale: bitimage2.NormalScale,
			verticalScale:   bitimage2.NormalScale,
			color:           bitimage2.Color4,
			width:           100,
			height:          100,
			data:            createColumnData(100, 100),
			wantErr:         bitimage2.ErrInvalidColor,
		},
		{
			name:            "invalid horizontal scale",
			horizontalScale: 0,
			verticalScale:   bitimage2.NormalScale,
			color:           bitimage2.Color1,
			width:           100,
			height:          100,
			data:            createColumnData(100, 100),
			wantErr:         bitimage2.ErrInvalidScale,
		},
		{
			name:            "invalid vertical scale",
			horizontalScale: bitimage2.NormalScale,
			verticalScale:   3,
			color:           bitimage2.Color1,
			width:           100,
			height:          100,
			data:            createColumnData(100, 100),
			wantErr:         bitimage2.ErrInvalidScale,
		},
		{
			name:            "width exceeds limit",
			horizontalScale: bitimage2.NormalScale,
			verticalScale:   bitimage2.NormalScale,
			color:           bitimage2.Color1,
			width:           2049,
			height:          100,
			data:            createColumnData(2049, 100),
			wantErr:         bitimage2.ErrInvalidWidth,
		},
		{
			name:            "height exceeds limit",
			horizontalScale: bitimage2.NormalScale,
			verticalScale:   bitimage2.NormalScale,
			color:           bitimage2.Color1,
			width:           100,
			height:          129,
			data:            createColumnData(100, 129),
			wantErr:         bitimage2.ErrInvalidHeight,
		},
		{
			name:            "invalid data length",
			horizontalScale: bitimage2.NormalScale,
			verticalScale:   bitimage2.NormalScale,
			color:           bitimage2.Color1,
			width:           100,
			height:          100,
			data:            []byte{0xFF}, // Should be more bytes
			wantErr:         bitimage2.ErrInvalidDataLength,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.StoreColumnGraphicsInBuffer(tt.horizontalScale, tt.verticalScale,
				tt.color, tt.width, tt.height, tt.data)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "StoreColumnGraphicsInBuffer") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			// Verify the command structure
			if got[0] != common.GS || got[1] != '(' || got[2] != 'L' {
				t.Errorf("StoreColumnGraphicsInBuffer: invalid command prefix")
			}
		})
	}
}

func TestGraphicsCommands_StoreColumnGraphicsInBufferLarge(t *testing.T) {
	cmd := bitimage2.NewGraphicsCommands()

	// Helper to create testutils data for column format
	createColumnData := func(width, height uint16) []byte {
		heightBytes := (int(height) + 7) / 8
		return testutils.RepeatByte(int(width)*heightBytes, 0xFF)
	}

	tests := []struct {
		name            string
		horizontalScale bitimage2.GraphicsScale
		verticalScale   bitimage2.GraphicsScale
		color           bitimage2.GraphicsColor
		width           uint16
		height          uint16
		data            []byte
		wantErr         error
	}{
		{
			name:            "large column data",
			horizontalScale: bitimage2.NormalScale,
			verticalScale:   bitimage2.NormalScale,
			color:           bitimage2.Color1,
			width:           2048,
			height:          128,
			data:            createColumnData(2048, 128),
			wantErr:         nil,
		},
		{
			name:            "color 4 not supported",
			horizontalScale: bitimage2.NormalScale,
			verticalScale:   bitimage2.NormalScale,
			color:           bitimage2.Color4,
			width:           100,
			height:          100,
			data:            createColumnData(100, 100),
			wantErr:         bitimage2.ErrInvalidColor,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.StoreColumnGraphicsInBufferLarge(tt.horizontalScale, tt.verticalScale,
				tt.color, tt.width, tt.height, tt.data)

			if !testutils.AssertErrorOccurred(t, err, tt.wantErr != nil, "StoreColumnGraphicsInBufferLarge") {
				return
			}
			if tt.wantErr != nil {
				testutils.AssertError(t, err, tt.wantErr)
				return
			}

			// Verify the command structure for large format
			if got[0] != common.GS || got[1] != '8' || got[2] != 'L' {
				t.Errorf("StoreColumnGraphicsInBufferLarge: invalid command prefix")
			}
		})
	}
}

// ============================================================================
// Validation Functions Tests
// ============================================================================

func TestValidateDensityFunctionCode(t *testing.T) {
	tests := []struct {
		name    string
		fn      bitimage2.FunctionCode
		wantErr bool
	}{
		{"valid code 1", bitimage2.FunctionCodeDensity1, false},
		{"valid code 49", bitimage2.FunctionCodeDensity49, false},
		{"invalid code 2", bitimage2.FunctionCodePrint2, true},
		{"invalid code 50", bitimage2.FunctionCodePrint50, true},
		{"invalid code 0", 0, true},
		{"invalid code 99", 99, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bitimage2.ValidateDensityFunctionCode(tt.fn)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDensityFunctionCode(%v) error = %v, wantErr %v", tt.fn, err, tt.wantErr)
			}
		})
	}
}

func TestValidatePrintFunctionCode(t *testing.T) {
	tests := []struct {
		name    string
		fn      bitimage2.FunctionCode
		wantErr bool
	}{
		{"valid code 2", bitimage2.FunctionCodePrint2, false},
		{"valid code 50", bitimage2.FunctionCodePrint50, false},
		{"invalid code 1", bitimage2.FunctionCodeDensity1, true},
		{"invalid code 49", bitimage2.FunctionCodeDensity49, true},
		{"invalid code 0", 0, true},
		{"invalid code 99", 99, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bitimage2.ValidatePrintFunctionCode(tt.fn)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePrintFunctionCode(%v) error = %v, wantErr %v", tt.fn, err, tt.wantErr)
			}
		})
	}
}

func TestValidateDotDensity(t *testing.T) {
	tests := []struct {
		name    string
		x       bitimage2.DotDensity
		y       bitimage2.DotDensity
		wantErr bool
	}{
		{"valid 180x180", bitimage2.Density180x180, bitimage2.Density180x180, false},
		{"valid 360x360", bitimage2.Density360x360, bitimage2.Density360x360, false},
		{"invalid x value", 49, bitimage2.Density180x180, true},
		{"invalid y value", bitimage2.Density180x180, 52, true},
		{"mismatched values", bitimage2.Density180x180, bitimage2.Density360x360, true},
		{"both invalid", 49, 52, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bitimage2.ValidateDotDensity(tt.x, tt.y)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDotDensity(%v, %v) error = %v, wantErr %v", tt.x, tt.y, err, tt.wantErr)
			}
		})
	}
}

func TestValidateGraphicsTone(t *testing.T) {
	tests := []struct {
		name    string
		tone    bitimage2.GraphicsTone
		wantErr bool
	}{
		{"valid monochrome", bitimage2.Monochrome, false},
		{"valid multiple tone", bitimage2.MultipleTone, false},
		{"invalid 49", 49, true},
		{"invalid 51", 51, true},
		{"invalid 0", 0, true},
		{"invalid 99", 99, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bitimage2.ValidateGraphicsTone(tt.tone)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGraphicsTone(%v) error = %v, wantErr %v", tt.tone, err, tt.wantErr)
			}
		})
	}
}

func TestValidateGraphicsScale(t *testing.T) {
	tests := []struct {
		name    string
		scale   bitimage2.GraphicsScale
		wantErr bool
	}{
		{"valid normal", bitimage2.NormalScale, false},
		{"valid double", bitimage2.DoubleScale, false},
		{"invalid 0", 0, true},
		{"invalid 3", 3, true},
		{"invalid 99", 99, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bitimage2.ValidateGraphicsScale(tt.scale)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGraphicsScale(%v) error = %v, wantErr %v", tt.scale, err, tt.wantErr)
			}
		})
	}
}

func TestValidateGraphicsColor(t *testing.T) {
	tests := []struct {
		name    string
		color   bitimage2.GraphicsColor
		wantErr bool
	}{
		{"valid color 1", bitimage2.Color1, false},
		{"valid color 2", bitimage2.Color2, false},
		{"valid color 3", bitimage2.Color3, false},
		{"valid color 4", bitimage2.Color4, false},
		{"invalid 48", 48, true},
		{"invalid 53", 53, true},
		{"invalid 0", 0, true},
		{"invalid 99", 99, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bitimage2.ValidateGraphicsColor(tt.color)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGraphicsColor(%v) error = %v, wantErr %v", tt.color, err, tt.wantErr)
			}
		})
	}
}

func TestValidateRasterDimensions(t *testing.T) {
	tests := []struct {
		name          string
		width         uint16
		height        uint16
		tone          bitimage2.GraphicsTone
		verticalScale bitimage2.GraphicsScale
		wantErr       bool
	}{
		{"valid monochrome normal", 100, 100, bitimage2.Monochrome, bitimage2.NormalScale, false},
		{"valid monochrome double", 100, 1200, bitimage2.Monochrome, bitimage2.DoubleScale, false},
		{"valid multiple tone normal", 100, 600, bitimage2.MultipleTone, bitimage2.NormalScale, false},
		{"valid multiple tone double", 100, 300, bitimage2.MultipleTone, bitimage2.DoubleScale, false},
		{"max width", 2400, 100, bitimage2.Monochrome, bitimage2.NormalScale, false},
		{"max monochrome height normal", 100, 2400, bitimage2.Monochrome, bitimage2.NormalScale, false},
		{"max monochrome height double", 100, 1200, bitimage2.Monochrome, bitimage2.DoubleScale, false},
		{"max multiple tone height normal", 100, 600, bitimage2.MultipleTone, bitimage2.NormalScale, false},
		{"max multiple tone height double", 100, 300, bitimage2.MultipleTone, bitimage2.DoubleScale, false},
		{"width zero", 0, 100, bitimage2.Monochrome, bitimage2.NormalScale, true},
		{"height zero", 100, 0, bitimage2.Monochrome, bitimage2.NormalScale, true},
		{"width exceeds", 2401, 100, bitimage2.Monochrome, bitimage2.NormalScale, true},
		{"monochrome height exceeds normal", 100, 2401, bitimage2.Monochrome, bitimage2.NormalScale, true},
		{"monochrome height exceeds double", 100, 1201, bitimage2.Monochrome, bitimage2.DoubleScale, true},
		{"multiple tone height exceeds normal", 100, 601, bitimage2.MultipleTone, bitimage2.NormalScale, true},
		{"multiple tone height exceeds double", 100, 301, bitimage2.MultipleTone, bitimage2.DoubleScale, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bitimage2.ValidateRasterDimensions(tt.width, tt.height, tt.tone, tt.verticalScale)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRasterDimensions(%v, %v, %v, %v) error = %v, wantErr %v",
					tt.width, tt.height, tt.tone, tt.verticalScale, err, tt.wantErr)
			}
		})
	}
}

func TestValidateColumnDimensions(t *testing.T) {
	tests := []struct {
		name    string
		width   uint16
		height  uint16
		wantErr bool
	}{
		{"minimum valid", 1, 1, false},
		{"typical dimensions", 100, 100, false},
		{"maximum width", 2048, 100, false},
		{"maximum height", 100, 128, false},
		{"maximum both", 2048, 128, false},
		{"width zero", 0, 100, true},
		{"height zero", 100, 0, true},
		{"width exceeds", 2049, 100, true},
		{"height exceeds", 100, 129, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bitimage2.ValidateColumnDimensions(tt.width, tt.height)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateColumnDimensions(%v, %v) error = %v, wantErr %v",
					tt.width, tt.height, err, tt.wantErr)
			}
		})
	}
}

// ============================================================================
// Helper Functions Tests
// ============================================================================

func TestCalculateRasterDataSize(t *testing.T) {
	tests := []struct {
		name   string
		width  uint16
		height uint16
		want   int
	}{
		{"8x8 pixels", 8, 8, 8},
		{"16x16 pixels", 16, 16, 32},
		{"100x50 pixels", 100, 50, 650},
		{"7x1 pixels", 7, 1, 1},
		{"9x1 pixels", 9, 1, 2},
		{"2400x1 pixels", 2400, 1, 300},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test through the command to ensure data validation works
			cmd := bitimage2.NewGraphicsCommands()

			// Calculate expected size
			widthBytes := (int(tt.width) + 7) / 8
			expectedSize := widthBytes * int(tt.height)

			// Create data with expected size
			data := testutils.RepeatByte(expectedSize, 0xFF)

			// Should succeed with correct data size
			_, err := cmd.StoreRasterGraphicsInBuffer(bitimage2.Monochrome, bitimage2.NormalScale,
				bitimage2.NormalScale, bitimage2.Color1, tt.width, tt.height, data)
			if err != nil {
				t.Errorf("calculateRasterDataSize validation failed: %v", err)
			}

			// Should fail with incorrect data size
			if expectedSize > 0 {
				wrongData := testutils.RepeatByte(expectedSize-1, 0xFF)
				_, err = cmd.StoreRasterGraphicsInBuffer(bitimage2.Monochrome, bitimage2.NormalScale,
					bitimage2.NormalScale, bitimage2.Color1, tt.width, tt.height, wrongData)
				if err == nil {
					t.Errorf("calculateRasterDataSize should have failed for incorrect data length")
				}
			}
		})
	}
}

func TestCalculateColumnDataSize(t *testing.T) {
	tests := []struct {
		name   string
		width  uint16
		height uint16
		want   int
	}{
		{"8x8 pixels", 8, 8, 8},
		{"16x16 pixels", 16, 16, 32},
		{"100x50 pixels", 100, 50, 700},
		{"1x7 pixels", 1, 7, 1},
		{"1x9 pixels", 1, 9, 2},
		{"1x128 pixels", 1, 128, 16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test through the command to ensure data validation works
			cmd := bitimage2.NewGraphicsCommands()

			// Calculate expected size
			heightBytes := (int(tt.height) + 7) / 8
			expectedSize := int(tt.width) * heightBytes

			// Create data with expected size
			data := testutils.RepeatByte(expectedSize, 0xFF)

			// Should succeed with correct data size
			_, err := cmd.StoreColumnGraphicsInBuffer(bitimage2.NormalScale, bitimage2.NormalScale,
				bitimage2.Color1, tt.width, tt.height, data)
			if err != nil {
				t.Errorf("calculateColumnDataSize validation failed: %v", err)
			}

			// Should fail with incorrect data size
			if expectedSize > 0 {
				wrongData := testutils.RepeatByte(expectedSize-1, 0xFF)
				_, err = cmd.StoreColumnGraphicsInBuffer(bitimage2.NormalScale, bitimage2.NormalScale,
					bitimage2.Color1, tt.width, tt.height, wrongData)
				if err == nil {
					t.Errorf("calculateColumnDataSize should have failed for incorrect data length")
				}
			}
		})
	}
}

func TestGetMaxHeight(t *testing.T) {
	tests := []struct {
		name          string
		tone          bitimage2.GraphicsTone
		verticalScale bitimage2.GraphicsScale
		want          uint16
	}{
		{"monochrome normal", bitimage2.Monochrome, bitimage2.NormalScale, 2400},
		{"monochrome double", bitimage2.Monochrome, bitimage2.DoubleScale, 1200},
		{"multiple tone normal", bitimage2.MultipleTone, bitimage2.NormalScale, 600},
		{"multiple tone double", bitimage2.MultipleTone, bitimage2.DoubleScale, 300},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test through validation
			// Should succeed with max height
			err := bitimage2.ValidateRasterDimensions(100, tt.want, tt.tone, tt.verticalScale)
			if err != nil {
				t.Errorf("getMaxHeight: should accept height %v for tone %v scale %v: %v",
					tt.want, tt.tone, tt.verticalScale, err)
			}

			// Should fail with max height + 1
			err = bitimage2.ValidateRasterDimensions(100, tt.want+1, tt.tone, tt.verticalScale)
			if err == nil {
				t.Errorf("getMaxHeight: should reject height %v for tone %v scale %v",
					tt.want+1, tt.tone, tt.verticalScale)
			}
		})
	}
}
