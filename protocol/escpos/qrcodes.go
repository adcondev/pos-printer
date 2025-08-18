package escpos

import (
	"fmt"

	"github.com/AdConDev/pos-printer/protocol/escpos/types"
	"github.com/AdConDev/pos-printer/utils"
)

// modelMap mapea los modelos de QR a sus valores ESC/POS
var modelMap = map[types.QRModel]byte{
	types.Model1: '1', // Modelo 1
	types.Model2: '2', // Modelo 2
}

// ecMap mapea los niveles de corrección de errores a sus valores ESC/POS
var ecMap = map[types.QRErrorCorrection]byte{
	types.ECLow:     '0', // 7% de corrección
	types.ECMedium:  '1', // 15% de corrección
	types.ECHigh:    '2', // 25% de corrección
	types.ECHighest: '3', // 30% de corrección
}

// PrintQR implementa el comando ESC Z para imprimir códigos QR
func (p *Commands) PrintQR(
	data string,
	model types.QRModel,
	moduleSize types.QRModuleSize,
	ecLevel types.QRErrorCorrection,
) ([][]byte, error) {
	// Validación de modelo
	if model < types.Model1 || model > types.Model2 {
		return nil, fmt.Errorf("modelo de QR inválida(0-1): %d", model)
	}

	// Comando para seleccionar tamaño del módulo
	mdl, err := p.SelectQRModel(model)
	if err != nil {
		return nil, fmt.Errorf("error al seleccionar modelo de QR: %w", err)
	}

	// Comando para seleccionar tamaño del módulo
	mdlSz, err := p.SelectQRSize(moduleSize)
	if err != nil {
		return nil, fmt.Errorf("error al seleccionar tamaño de módulo de QR: %w", err)
	}

	// Obtener el byte correspondiente al nivel de corrección
	ec, err := p.SelectQRErrorCorrection(ecLevel)
	if err != nil {
		return nil, fmt.Errorf("error al seleccionar nivel de corrección de QR: %w", err)
	}

	// Almacenamiento de datos para QR
	ct, err := p.SetQRData(data)
	if err != nil {
		return nil, fmt.Errorf("error al preparar datos de QR: %w", err)
	}

	// Comando para imprimir QR
	prnt, err := p.PrintQRData()
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
func (p *Commands) SelectQRModel(model types.QRModel) ([]byte, error) {
	// Validación de modelo
	if model < types.Model1 || model > types.Model2 {
		return nil, fmt.Errorf("modelo de QR inválida(0-1): %d", model)
	}

	pL, pH, err := utils.LengthLowHigh(4)
	if err != nil {
		return nil, fmt.Errorf("error al calcular longitud de parametros QR: %w", err)
	}
	cn, fn := byte('1'), byte('A')
	n1 := modelMap[model]
	n2 := byte(0) // Siempre 0, reservado

	cmd := make([]byte, 0, 9)
	cmd = append(cmd, GS, '(', 'k') // Comando QR
	cmd = append(cmd, pL, pH, cn, fn, n1, n2)

	return cmd, nil
}

// SelectQRSize selecciona el tamaño del módulo del código QR
func (p *Commands) SelectQRSize(moduleSize types.QRModuleSize) ([]byte, error) {
	// Validar tamaño del módulo
	if moduleSize < types.MinType || moduleSize > types.MaxType {
		return nil, fmt.Errorf("tamaño de módulo QR inválido(1-16): %d", moduleSize)
	}

	pL, pH, err := utils.LengthLowHigh(3)
	if err != nil {
		return nil, fmt.Errorf("error al calcular longitud de parametros QR: %w", err)
	}
	cn, fn := byte('1'), byte('C')
	n := byte(moduleSize)

	cmd := make([]byte, 0, 8)
	cmd = append(cmd, GS, '(', 'k') // Comando QR
	cmd = append(cmd, pL, pH, cn, fn, n)

	return cmd, nil
}

// SelectQRErrorCorrection selecciona el nivel de corrección de errores del código QR
func (p *Commands) SelectQRErrorCorrection(level types.QRErrorCorrection) ([]byte, error) {
	// Validar nivel de corrección
	ec, ok := ecMap[level]
	if !ok {
		return nil, fmt.Errorf("nivel de corrección de QR inválido(0-3): %d", level)
	}

	pL, pH, err := utils.LengthLowHigh(3)
	if err != nil {
		return nil, fmt.Errorf("error al calcular longitud de parametros QR: %w", err)
	}
	cn, fn := byte('1'), byte('E')

	cmd := make([]byte, 0, 8)
	cmd = append(cmd, GS, '(', 'k') // Comando QR
	cmd = append(cmd, pL, pH, cn, fn, ec)

	return cmd, nil
}

// SetQRData prepara los datos para el código QR
func (p *Commands) SetQRData(data string) ([]byte, error) {
	// Validar longitud de datos
	if len(data) == 0 || len(data) > 7089 {
		return nil, fmt.Errorf("longitud de datos de QR inválida (1-7089): %d", len(data))
	}

	pL, pH, err := utils.LengthLowHigh(len(data) + 3)
	if err != nil {
		return nil, fmt.Errorf("error al calcular longitud de parametros QR: %w", err)
	}
	cn, fn := byte('1'), byte('P')
	m := byte('0') // Siempre 0, reservado

	cmd := make([]byte, 0, 7+len(data))
	cmd = append(cmd, GS, '(', 'k') // Comando QR
	cmd = append(cmd, pL, pH, cn, fn, m)
	cmd = append(cmd, data...)

	return cmd, nil
}

// PrintQRData genera el comando para imprimir el código QR
func (p *Commands) PrintQRData() ([]byte, error) {
	// Comando para imprimir QR
	pL, pH, err := utils.LengthLowHigh(3)
	if err != nil {
		return nil, fmt.Errorf("error al calcular longitud de parametros QR: %w", err)
	}
	cn, fn := byte('1'), byte('Q')
	m := byte('0') // Siempre 0 para impresion estandard

	cmd := make([]byte, 0, 8)
	cmd = append(cmd, GS, '(', 'k') // Comando QR
	cmd = append(cmd, pL, pH, cn, fn, m)

	return cmd, nil
}
