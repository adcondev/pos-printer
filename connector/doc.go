/*
Package connector proporciona interfaces para la comunicación con impresoras ESC/POS
en múltiples sistemas operativos. La interfaz principal extiende io.WriteCloser para
facilitar la escritura y el cierre seguro de conexiones. En Windows, se implementa un
conector nativo utilizando la API (winspool.drv); en otros sistemas, se retorna un stub
que notifica la indisponibilidad de la funcionalidad.
*/
package connector
