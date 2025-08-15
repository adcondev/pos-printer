package connector

import "io"

// Connector define la interfaz para cualquier tipo de conexión con la impresora
type Connector interface {
	io.WriteCloser // Write([]byte) (int, error) y Close() error

	// TODO: Agregar más métodos si necesitas:
	// - Read([]byte) (int, error) para status
	// - IsConnected() bool
	// - Reset() error
}

// El WindowsPrintConnector ya implementa Write y Close,
// así que debería cumplir con esta interfaz
