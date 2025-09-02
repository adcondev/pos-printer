package escpos

import (
	"fmt"

	"github.com/adcondev/pos-printer/escpos/common"
)

// modelMap mapea los modelos de QR a sus valores ESC/POS
var modelMap = map[QRModel]byte{
	Model1: '1', // Modelo 1
	Model2: '2', // Modelo 2
}

// ecMap mapea los niveles de corrección de errores a sus valores ESC/POS
var ecMap = map[QRErrorCorrection]byte{
	ECLow:     '0', // 7% de corrección
	ECMedium:  '1', // 15% de corrección
	ECHigh:    '2', // 25% de corrección
	ECHighest: '3', // 30% de corrección
}

// PrintQR implementa el comando ESC Z para imprimir códigos QR
func (c *Protocol) PrintQR(
	data string,
	model QRModel,
	moduleSize QRModuleSize,
	ecLevel QRErrorCorrection,
) ([][]byte, error) {
	// Validación de modelo
	if model < Model1 || model > Model2 {
		return nil, fmt.Errorf("modelo de QR inválida(0-1): %d", model)
	}

	// Comando para seleccionar tamaño del módulo
	mdl, err := c.SelectQRModel(model)
	if err != nil {
		return nil, fmt.Errorf("error al seleccionar modelo de QR: %w", err)
	}

	// Comando para seleccionar tamaño del módulo
	mdlSz, err := c.SelectQRSize(moduleSize)
	if err != nil {
		return nil, fmt.Errorf("error al seleccionar tamaño de módulo de QR: %w", err)
	}

	// Obtener el byte correspondiente al nivel de corrección
	ec, err := c.SelectQRErrorCorrection(ecLevel)
	if err != nil {
		return nil, fmt.Errorf("error al seleccionar nivel de corrección de QR: %w", err)
	}

	// Almacenamiento de datos para QR
	ct, err := c.SetQRData(data)
	if err != nil {
		return nil, fmt.Errorf("error al preparar datos de QR: %w", err)
	}

	// Comando para imprimir QR
	prnt, err := c.PrintQRData()
	if err != nil {
		return nil, fmt.Errorf("error al generar comando de impresión de QR: %w", err)
	}

	cmdLines := [][]byte{mdl, mdlSz, ec, ct, prnt}
	if len(cmdLines) == 0 {
		return nil, fmt.Errorf("no se generaron comandos para imprimir QR")
	}

	return cmdLines, nil
}

// SelectQRModel selecciona el modelo de código QR a utilizar
func (c *Protocol) SelectQRModel(model QRModel) ([]byte, error) {
	// Validación de modelo
	if model < Model1 || model > Model2 {
		return nil, fmt.Errorf("modelo de QR inválida(0-1): %d", model)
	}

	pL, pH := common.LengthLowHigh(4)

	cn, fn := byte('1'), byte('A')
	n1 := modelMap[model]
	n2 := byte(0) // Siempre 0, reservado

	cmd := make([]byte, 0, 9)
	cmd = append(cmd, common.GS, '(', 'k') // Comando QR
	cmd = append(cmd, pL, pH, cn, fn, n1, n2)

	return cmd, nil
}

// SelectQRSize selecciona el tamaño del módulo del código QR
func (c *Protocol) SelectQRSize(moduleSize QRModuleSize) ([]byte, error) {
	// Validar tamaño del módulo
	if moduleSize < MinType || moduleSize > MaxType {
		return nil, fmt.Errorf("tamaño de módulo QR inválido(1-16): %d", moduleSize)
	}

	pL, pH := common.LengthLowHigh(3)

	cn, fn := byte('1'), byte('C')
	n := byte(moduleSize)

	cmd := make([]byte, 0, 8)
	cmd = append(cmd, common.GS, '(', 'k') // Comando QR
	cmd = append(cmd, pL, pH, cn, fn, n)

	return cmd, nil
}

// SelectQRErrorCorrection selecciona el nivel de corrección de errores del código QR
func (c *Protocol) SelectQRErrorCorrection(level QRErrorCorrection) ([]byte, error) {
	// Validar nivel de corrección
	ec, ok := ecMap[level]
	if !ok {
		return nil, fmt.Errorf("nivel de corrección de QR inválido(0-3): %d", level)
	}

	pL, pH := common.LengthLowHigh(3)

	cn, fn := byte('1'), byte('E')

	cmd := make([]byte, 0, 8)
	cmd = append(cmd, common.GS, '(', 'k') // Comando QR
	cmd = append(cmd, pL, pH, cn, fn, ec)

	return cmd, nil
}

// SetQRData prepara los datos para el código QR
func (c *Protocol) SetQRData(data string) ([]byte, error) {
	// Validar longitud de datos
	if len(data) == 0 || len(data) > 7089 {
		return nil, fmt.Errorf("longitud de datos de QR inválida (1-7089): %d", len(data))
	}

	// Secure, it is validated before.
	pL, pH := common.LengthLowHigh(uint16(len(data) + 3)) // nolint:gosec

	cn, fn := byte('1'), byte('P')
	m := byte('0') // Siempre 0, reservado

	cmd := make([]byte, 0, 7+len(data))
	cmd = append(cmd, common.GS, '(', 'k') // Comando QR
	cmd = append(cmd, pL, pH, cn, fn, m)
	cmd = append(cmd, data...)

	return cmd, nil
}

// PrintQRData genera el comando para imprimir el código QR
func (c *Protocol) PrintQRData() ([]byte, error) {
	// Comando para imprimir QR
	pL, pH := common.LengthLowHigh(3)

	cn, fn := byte('1'), byte('Q')
	m := byte('0') // Siempre 0 para impresion estandard

	cmd := make([]byte, 0, 8)
	cmd = append(cmd, common.GS, '(', 'k') // Comando QR
	cmd = append(cmd, pL, pH, cn, fn, m)

	return cmd, nil
}
