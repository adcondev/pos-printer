// Package posprinter proporciona una interfaz unificada para interactuar con impresoras
// térmicas de puntos de venta (POS), utilizando una arquitectura modular que separa
// protocolos, conectores y perfiles de impresoras.
//
// # Arquitectura
//
// El paquete está diseñado con tres componentes principales:
//
//   - Protocol: Implementa comandos específicos del protocolo (ESC/POS, ZPL, etc.)
//   - Connector: Maneja la comunicación con la impresora (Windows, Red, USB, etc.)
//   - Profile: Define las características y capacidades de cada modelo de impresora
//
// Esta separación permite cambiar cualquier componente sin afectar los demás,
// facilitando la compatibilidad con diferentes impresoras y métodos de conexión.
//
// Uso básico
//
//	conn, err := connector.NewWindowsPrintConnector("Nombre de Impresora")
//	proto := escpos.NewESCPOSProtocol()
//	prof := profile.CreateProfile58mm()
//	printer, err := posprinter.NewGenericPrinter(proto, conn, prof)
//
//	// Imprimir texto
//	printer.TextLn("Hola Mundo")
//
//	// Formateo
//	printer.TurnEmphasizedMode(true)
//	printer.TextLn("Texto en negrita")
//
//	// Imprimir imagen
//	printer.PrintImageFromFile("logo.png")
//
//	// Cortar papel
//	printer.Cut(types.CutFeed, 3)
//
// Características principales
//
//   - Impresión de texto con múltiples juegos de caracteres
//   - Formateo de texto (negrita, subrayado, alineación)
//   - Impresión de imágenes con algoritmos de dithering
//   - Generación e impresión de códigos QR y códigos de barras
//   - Soporte para diferentes modelos de impresoras térmicas
//
// Para más detalles, consulte los ejemplos en la carpeta 'examples'.
package internal
