package protocol

import (
	"github.com/AdConDev/pos-printer/imaging"
	"github.com/AdConDev/pos-printer/protocol/escpos/types"
)

// Protocol define una interfaz para cualquier protocolo de impresión.
// Esta interfaz devuelve comandos en bytes que el conector enviará a la impresora
type Protocol interface {
	// === Comandos básicos ===

	InitializePrinter() []byte
	Close() []byte

	// === Manipulación de texto ===

	SetJustification(justification types.Alignment) []byte
	SelectCharacterFont(n types.Font) ([]byte, error)
	TurnEmphasizedMode(n types.EmphasizedMode) ([]byte, error)
	SetDoubleStrike(on bool) []byte
	TurnUnderlineMode(n types.UnderlineMode) ([]byte, error)
	SetTextSize(widthMultiplier int, heightMultiplier int) []byte
	SetLineSpacing(n types.LineSpace) []byte
	SetPrintLeftMargin(margin int) []byte
	SetPrintWidth(width int) []byte

	// === Manejo de Character Tables/Code Pages ===

	SelectCharacterTable(table types.CharacterSet) []byte
	CancelKanjiMode() []byte

	// === Comandos de texto ===

	Text(str string) []byte
	TextLn(str string) []byte
	TextRaw(str string) []byte

	// === Códigos de barras ===

	SetBarcodeHeight(height int) []byte
	SetBarcodeWidth(width int) []byte
	SetBarcodeTextPosition(position types.BarcodeTextPosition) []byte
	Barcode(content string, barType types.BarcodeType) ([]byte, error)

	// === Impresión de códigos QR ===

	PrintQR(string, types.QRModel, types.QRModuleSize, types.QRErrorCorrection) ([][]byte, error)

	// === Imágenes ===

	// PrintRasterBitImage recibe una imagen genérica y la convierte a comandos del protocolo
	PrintRasterBitImage(img *imaging.PrintImage, density types.Density) ([]byte, error)

	// GetMaxImageWidth devuelve el ancho máximo de imagen soportado
	GetMaxImageWidth(paperWidth int, dpi int) int

	// === Control de papel ===

	Cut(mode types.CutMode, lines int) []byte
	Feed(lines int) []byte
	FeedReverse(lines int) []byte
	FeedForm() []byte

	// === Hardware ===

	Pulse(pin int, onMS int, offMS int) []byte
	Release() []byte

	// === Información del protocolo ===

	Name() string

	// TODO: Agregar más métodos según necesites:
	// - PrintQRCode(data string, size int) ([]byte, error)
	// - SetPrintSpeed(speed int) []byte
	// - GetStatus() []byte
	// etc.
}

// Factory es una función que crea una instancia de un protocolo
type Factory func() Protocol
