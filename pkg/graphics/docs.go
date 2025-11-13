// Package graphics provides an advanced image processing engine for ESC/POS printers
//
// # Guía de Tamaños para QR Codes
//
// La siguiente tabla muestra los tamaños TOTALES (incluyendo quiet zone de 4 modules)
// para diferentes configuraciones:
//
// Grid Size | Module 8px  | Module 12px | Module 16px
// ----------|-------------|-------------|-------------
// 21×21 (V1)| 232px total | 348px total | 464px total
//
//	| (168+64)    | (252+96)    | (336+128)
//
// 25×25 (V2)| 264px total | 396px total | 528px total
//
//	| (200+64)    | (300+96)    | (400+128)
//
// 29×29 (V3)| 296px total | 444px total | 592px total
//
//	| (232+64)    | (348+96)    | (464+128)
//
// 33×33 (V4)| 328px total | 492px total | 656px total
//
//	| (264+64)    | (396+96)    | (528+128)
//
// Formato: TotalSize (DataSize + BorderSize)
// BorderSize = 8 modules × moduleSize (4 modules por cada lado)
//
// Para calcular el PixelWidth necesario:
//
//	totalModules = gridSize + 8
//	moduleSize = deseado para cada módulo del QR
//	PixelWidth = totalModules × moduleSize
//
// Para calcular el tamaño de datos real:
//
//	dataSize = gridSize × moduleSize
//	borderSize = 8 × moduleSize
package graphics
