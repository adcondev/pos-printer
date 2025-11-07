# Análisis Estratégico y Arquitectura para una Biblioteca ESC/POS Superior en Go

Este documento es un análisis de la situación actual de las bibliotecas ESC/POS, con un enfoque en Go, y una propuesta
de arquitectura para una nueva biblioteca diseñada para ser superior en todos los aspectos clave.

## Parte 1: El Panorama de la Traducción de JSON a ESC/POS

Tu interés en la capa JSON -> Comandos es el núcleo del problema en los sistemas POS modernos. Los desarrolladores de
aplicaciones (web, móviles, de escritorio) no quieren (ni deben) construir arrays de bytes. Quieren definir un recibo en
un formato de datos universal.

### Enfoques Existentes para la Traducción

Existen dos enfoques fundamentales para mapear JSON a comandos de impresión:

#### Enfoque 1: Mapeo Directo de Comandos (El Enfoque de "Bajo Nivel")

Este enfoque utiliza JSON para describir la secuencia de comandos a ejecutar. Es una simple serialización de la API de
la biblioteca.

**Ejemplo de JSON:**

```json
[
  {
    "command": "init"
  },
  {
    "command": "setTextSize",
    "width": 2,
    "height": 2
  },
  {
    "command": "write",
    "text": "Total:"
  },
  {
    "command": "align",
    "align": "right"
  },
  {
    "command": "write",
    "text": "$10.00"
  },
  {
    "command": "feed",
    "lines": 3
  },
  {
    "command": "cut"
  }
]
```

* **Ventajas:** Flexible. Expones 1:1 toda la funcionalidad de la impresora.
* **Desventajas:** Horrible para el desarrollador de la aplicación. Está totalmente acoplado a los comandos ESC/POS. No
  hay abstracción. El desarrollador necesita saber ESC/POS.

#### Enfoque 2: Formato Declarativo / Semántico (El Enfoque "Preferido")

Este enfoque utiliza JSON para describir la estructura y semántica del documento (el recibo), no los comandos de la
impresora. La biblioteca hace el trabajo pesado de traducir esa intención a los comandos de bytes óptimos.

**Ejemplo de JSON:**

```json
{
  "header": {
    "type": "image",
    "path": "logo.png",
    "align": "center"
  },
  "body": [
    {
      "type": "text",
      "content": "Concepto",
      "style": [
        "bold"
      ],
      "width": 24
    },
    {
      "type": "text",
      "content": "Precio",
      "style": [
        "bold"
      ],
      "align": "right",
      "width": 16
    },
    {
      "type": "divider",
      "char": "-"
    },
    {
      "type": "text",
      "content": "Producto 1",
      "width": 24
    },
    {
      "type": "text",
      "content": "$10.00",
      "align": "right",
      "width": 16
    },
    {
      "type": "feed",
      "lines": 2
    },
    {
      "type": "barcode",
      "format": "QR",
      "data": "[https://mi.dominio.com/factura/123](https://mi.dominio.com/factura/123)",
      "align": "center"
    }
  ],
  "footer": [
    {
      "type": "cut",
      "mode": "partial"
    }
  ]
}
```

* **Ventajas:**
* **Desacoplamiento total:** El desarrollador de la aplicación solo describe qué quiere imprimir.
    * **Portabilidad:** El mismo JSON podría (en teoría) ser renderizado a HTML, a PDF o a otro protocolo de impresora.
    * **Inteligencia de la Biblioteca:** La biblioteca puede manejar lógicas complejas, como el ajuste de palabras (
      word-wrap), el diseño de columnas (como en el ejemplo de "Concepto" y "Precio"), y la codificación de caracteres.
* **Desventajas:** Requiere un trabajo de diseño de esquema (schema) mucho mayor.

### ¿Qué se Prefiere en los Sistemas POS?

El **Enfoque 2 (Declarativo)** es abrumadoramente preferido.

* Los sistemas POS modernos (Odoo, Shopify POS, Square) exponen una API de "plantillas" o "recibos" a sus
  desarrolladores, no una API de "impresora".
* **Odoo** usa un sistema de plantillas XML (QWeb) que se renderiza a HTML/PDF y, para la impresión directa, a un
  formato intermedio (a menudo basado en XML o JSON) que su "IoT Box" o cliente web traduce a ESC/POS.
* **Star Micronics** fue pionero en esto con su "StarPRNT JSON," que es un formato declarativo/semántico.
* **Epson** tiene "ePOS-Print," que es un formato XML que cumple el mismo propósito.

### Conclusión

Tu biblioteca debe centrarse en el **Enfoque 2**. El mapeo JSON debe ser semántico.

-----

## Parte 2: Análisis de las Bibliotecas ESC/POS en Go (Los Rivales)

He investigado los repositorios más relevantes en GitHub para "escpos golang". Los dos principales contendientes (y
tus "rivales a vencer") son:

* `github.com/qplC/escpos`
* `github.com/knq/escpos`

### Rival 1: github.com/qplC/escpos (El Rival Principal)

* **Actividad:** Alta. Es el más activo y mantenido. Los PRs se revisan y los issues se discuten.
* **Enfoque:** Es una biblioteca de "bajo nivel". Proporciona una API fluida (fluent API) para construir comandos.

    ```go
    p := escpos.New(conn)
    p.Init().
      SetFont("A").
      SetAlign("center").
      Println("Hello World").
      Cut().
      End()
    ```

* **Ventajas:**
    * Buena cobertura de comandos básicos y algunos avanzados (códigos de barras).
    * La API fluida es cómoda para tareas simples.
    * Comunidad activa (para este nicho).
* **Desventajas y "Drawbacks" (Críticos para tu objetivo):**
    * **Acoplamiento de Conexión:** El constructor `escpos.New()` toma un `io.ReadWriteCloser` (la conexión). Esto es un
      anti-patrón en Go. La biblioteca no debería gestionar la conexión, solo generar los bytes. El usuario debería ser
      responsable de escribir esos bytes en `net.Conn`, `os.File` (para USB), etc.
    * **Sin Abstracción Semántica:** Es 100% Enfoque 1. El desarrollador debe saber ESC/POS. No hay un traductor JSON o
      de structs.
    * **Gestión de Gráficos Primitiva:** El `p.Image()` es básico. La impresión de imágenes es un dolor de cabeza en
      ESC/POS (dithering, algoritmos, modos de impresión). Los issues a menudo mencionan problemas con las imágenes.
    * **Codificación de Caracteres:** La gestión de páginas de códigos (code pages) y UTF-8 es un problema recurrente en
      todas las bibliotecas ESC/POS. Los issues sobre caracteres chinos, japoneses o acentos (ñ, á, é) son comunes.
    * **Sin Comunicación Bidireccional:** La biblioteca está diseñada para escribir, no para leer. No hay una API para
      consultar el estado de la impresora (sin papel, tapa abierta), que es una característica "profesional" vital.

### Rival 2: <https://github.com/knq/escpos>

* **Actividad:** Baja. Parece estar en modo de mantenimiento o archivado.
* **Enfoque:** Similar al de qplC, es un constructor de comandos de bajo nivel.
* **Desventajas:**
    * Menos completo que `qplC/escpos`.
    * Sufre de los mismos problemas: acoplamiento de conexión, sin abstracción, gestión de gráficos débil.

### Conclusión

`knq/escpos` fue un buen esfuerzo inicial, pero `qplC/escpos` es tu verdadero rival a vencer.

### ¿Dónde se Usan y Quién es el Rival a Vencer?

Se usan en aplicaciones de Go de back-end que necesitan imprimir directamente en una impresora de red (por ejemplo, en
una cocina o en un servidor de POS).

Tu rival a vencer es `github.com/qplC/escpos`. Tu biblioteca será superior si resuelves sus *drawbacks* fundamentales.

-----

## Parte 3: El Plan para una Biblioteca "Superior"

Para ser "superior en todos los aspectos", tu biblioteca debe tener una arquitectura de dos capas (Core + Abstracción) y
abordar los puntos débiles históricos (gráficos, codificación, estado).

### Arquitectura Propuesta: El Modelo de Dos Capas

Tu biblioteca `github.com/tu-usuario/go-escpos` debería tener dos paquetes principales:

#### Capa 1: `package escpos` (El Núcleo de Comandos)

Esta es tu versión mejorada de `qplC/escpos`. Es una API de bajo nivel, fluida y **sin estado**.

* **Punto Clave: Sin Gestión de Conexión.**
    * **MAL (El rival):** `p := escpos.New(conn)`
    * **BIEN (Tu biblioteca):** `p := escpos.NewPrinter(profile)` (opcionalmente un perfil de impresora).
* Todas las funciones (ej. `p.Text("hola")`, `p.Cut()`) no escriben en una conexión. Internamente, escriben en un
  `bytes.Buffer`.
* El usuario obtiene los bytes finales llamando a `p.Bytes()`.
* El usuario es responsable de la E/S: `conn.Write(p.Bytes())`.

> **¿Por qué es superior?** Es idiomático en Go. Permite *composability*. El usuario puede escribir a un archivo, a una
> conexión de red, a `stdout` o a un buffer para pruebas. Desacopla la generación de comandos del transporte.

#### Capa 2: `package receipt` (El Traductor Declarativo)

Este es el paquete que te hará ganar. Implementa el Enfoque 2 (Declarativo) y traduce JSON (o structs de Go) a la Capa
1.

El usuario define un recibo en Go:

```go
doc := receipt.Document{
Body: []receipt.Block{
receipt.TextBlock{Content: "Total:", Style: receipt.StyleBold},
receipt.TextBlock{Content: "$10.00", Align: receipt.AlignRight},
receipt.QRBlock{Data: "...", Align: receipt.AlignCenter},
},
Footer: receipt.CutBlock{Mode: receipt.CutPartial},
}
```

* El usuario puede hacer `unmarshal` de JSON a este struct `receipt.Document`.
* El paquete `receipt` tiene un `Render()` que toma el `Document` y usa el `package escpos` (Capa 1) para generar el
  `[]byte` final.

El "sueño" del desarrollador de la aplicación:

```go
// El "sueño" del desarrollador de la aplicación:
jsonData, _ := os.ReadFile("receipt.json")
doc, _ := receipt.Unmarshal(jsonData)

printer := escpos.NewPrinter() // Usando la Capa 1
bytes, err := doc.Render(printer)
if err != nil {
// Manejar error de renderizado (ej. imagen no encontrada)
}

conn.Write(bytes) // El usuario maneja la conexión
```

### Características Clave para la Superioridad

Para superar a `qplC/escpos` y otros:

* **Arquitectura de Dos Capas (Core + Receipt):** Como se describió. Esto por sí solo ya es una victoria masiva.
* **Motor de Gráficos Avanzado:**
    * La impresión de imágenes es el issue nº 1.
    * Tu `p.Image(img image.Image)` debe ser brillante.
    * Debe incluir algoritmos de *dithering* (ej. Floyd-Steinberg, Bayer) y diferentes modos de umbral (*thresholding*).
    * Debe manejar la impresión de "columnas" y "páginas" de gráficos automáticamente.
* **Gestión de Codificación y UTF-8 Robusta:**
    * Este es el issue nº 2.
    * La biblioteca debe (idealmente) tomar texto UTF-8 y hacer lo correcto, ya sea:
        * Cambiando a la página de códigos correcta (ej. CP850 para español).
        * Usando los comandos de impresión de UTF-8 si la impresora lo soporta.
    * Debe ser automático o con una API muy simple: `p.SetEncoding("UTF-8")` o `p.SetCodePage(escpos.CP_LATIN1)`.
* **Comunicación Bidireccional (Estado de la Impresora):**
    * Este es el "killer feature" que nadie tiene bien.
    * Proporciona funciones que escriben un comando de solicitud de estado (DLE EOT) y leen/parsean la respuesta.
    * Tu API: `func (p *Printer) GetStatus(r io.Reader) (Status, error)`
    * Struct de Retorno: `type Status struct { IsOnline bool; IsPaperOut bool; IsCoverOpen bool; ... }`
* **Perfiles de Impresora (Capabilities):**
    * No todas las impresoras son iguales. Epson, Star y Bixolon tienen ligeras diferencias.
    * Tu `escpos.NewPrinter(profile)` debería tomar un perfil (cargado desde un JSON/YAML) que defina las capacidades (
      ej. `hasCutter`, `supportsUTF8`, `maxCharsPerLine`).
    * Esto permite a tu capa `receipt` tomar decisiones inteligentes (ej. no intentar cortar si `hasCutter: false`).
* **Documentación y Ejemplos (Recetas):**
    * La documentación de las bibliotecas rivales es escasa.
    * Tu README y `godoc` deben ser impecables, con "recetas" para cada tarea común: "Cómo imprimir un recibo de
      cocina", "Cómo imprimir un logo y un QR", "Cómo leer el estado de la impresora".

-----

## Parte 4: Solución a Discutir (El Esquema Declarativo)

Aquí hay un borrador de los structs de Go (y por extensión, el JSON) para tu paquete `receipt`. Se basa en un diseño
de "bloques", similar a los editores de contenido modernos.

### receipt.go (El esquema de la Capa 2)

```go
package receipt

import (
	"image"

	"[github.com/tu-usuario/go-escpos/escpos](https://github.com/tu-usuario/go-escpos/escpos)" // Importa tu Capa 1

)

// Constantes para estilos y alineación
type Align string

const (
	AlignLeft   Align = "left"
	AlignCenter Align = "center"
	AlignRight  Align = "right"
)

type Style string

const (
	StyleBold      Style = "bold"
	StyleUnderline Style = "underline"
	StyleInverse   Style = "inverse"
)

type FontSize string

const (
	SizeSmall  FontSize = "small"
	SizeNormal FontSize = "normal"
	SizeLarge  FontSize = "large" // (ej. 2x2)
	SizeTall   FontSize = "tall"  // (ej. 1x2)
	SizeWide   FontSize = "wide"  // (ej. 2x1)
)

// Document es la raíz del recibo.
type Document struct {
	Header []Block `json:"header,omitempty"`
	Body   []Block `json:"body"`
	Footer []Block `json:"footer,omitempty"`
}

// Block es la interfaz que implementa cada tipo de contenido.
// `json:"-"` asegura que la lógica de renderizado no sea parte del JSON.
type Block interface {
	Render(p *escpos.Printer) error
`json:"-"`
}

// --- Implementaciones de Bloques ---

// TextBlock: Para texto simple o formateado.
type TextBlock struct {
	Content string   `json:"content"`
	Style   []Style  `json:"style,omitempty"`
	Size    FontSize `json:"size,omitempty"`
	Align   Align    `json:"align,omitempty"`
}

// Implementa Block.Render(p)...

// ColumnsBlock: Para diseño de columnas (clave para recibos).
type ColumnsBlock struct {
	Columns []Column `json:"columns"`
}
type Column struct {
	Content string  `json:"content"`
	Width   int     `json:"width"` // Ancho en caracteres
	Style   []Style `json:"style,omitempty"`
	Align   Align   `json:"align,omitempty"`
}

// Implementa Block.Render(p)...

// ImageBlock: Para logos e imágenes.
type ImageBlock struct {
	// Podrías aceptar path, URL, o bytes en base64
	Base64Data string `json:"base64,omitempty"`
	Align      Align  `json:"align,omitempty"`

	// Campo privado para la imagen decodificada
	img image.Image `json:"-"`
}

// Implementa Block.Render(p)...

// BarcodeBlock: Para códigos de barras 1D y 2D.
type BarcodeBlock struct {
	Format string `json:"format"` // ej: "QR", "CODE128", "EAN13"
	Data   string `json:"data"`
	Align  Align  `json:"align,omitempty"`
}

// Implementa Block.Render(p)...

// FeedBlock: Para saltos de línea.
type FeedBlock struct {
	Lines int `json:"lines"`
}

// Implementa Block.Render(p)...

// CutBlock: Para cortar el papel.
type CutBlock struct {
	Mode string `json:"mode"` // "full" o "partial"
}

// Implementa Block.Render(p)...

// DividerBlock: Para líneas separadoras (ej. '---' o '===')
type DividerBlock struct {
	Char string `json:"char,omitempty"` // Default: "-"
}

// Implementa Block.Render(p)...

// Unmarshal y Render serían las funciones principales de este paquete.
```

### JSON de Ejemplo (basado en los structs anteriores)

```json
{
  "header": [
    {
      "base64": "iVBORw0KGgoAAAANSUhEUg...",
      "align": "center"
    }
  ],
  "body": [
    {
      "content": "RECIBO DE VENTA",
      "size": "wide",
      "align": "center"
    },
    {
      "lines": 1
    },
    {
      "columns": [
        {
          "content": "Producto",
          "width": 24,
          "style": [
            "bold"
          ]
        },
        {
          "content": "Total",
          "width": 16,
          "style": [
            "bold"
          ],
          "align": "right"
        }
      ]
    },
    {
      "char": "-"
    },
    {
      "columns": [
        {
          "content": "1x Tacos al Pastor",
          "width": 24
        },
        {
          "content": "$120.00",
          "width": 16,
          "align": "right"
        }
      ]
    },
    {
      "columns": [
        {
          "content": "1x Agua de Horchata",
          "width": 24
        },
        {
          "content": "$30.00",
          "width": 16,
          "align": "right"
        }
      ]
    },
    {
      "lines": 2
    },
    {
      "content": "¡Gracias por su visita!",
      "align": "center"
    },
    {
      "format": "QR",
      "data": "[https://mi.negocio.com/factura/12345](https://mi.negocio.com/factura/12345)",
      "align": "center"
    }
  ],
  "footer": [
    {
      "lines": 3
    },
    {
      "mode": "partial"
    }
  ]
}
```

Este plan te da una hoja de ruta clara para construir una biblioteca que no solo es "una más", sino que resuelve
problemas reales que los desarrolladores de Go enfrentan hoy en día con las soluciones existentes.