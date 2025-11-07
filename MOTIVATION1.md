# Reporte Estratégico: Arquitectura de una Librería ESC/POS de Siguiente Generación en Go

## 1\. Análisis de Enfoques: Mapeo de Formatos Estructurados a ESC/POS

La transformación de un formato de datos estructurado (como JSON) en un flujo de bytes ESC/POS es un requisito
fundamental para cualquier sistema POS moderno. Esta capa de "plantillas" (templating) define la experiencia del
desarrollador, la mantenibilidad del sistema y la flexibilidad para que usuarios no técnicos modifiquen los recibos. El
análisis de los repositorios de código abierto revela tres patrones de traducción principales, además de la API
programática fundamental sobre la que se construyen.

### 1.1. Panorama de Patrones de Traducción

#### 1.1.1. Enfoque 1: Mapeo Directo 1:1 (JSON-a-Comando)

Este es el enfoque más simplista, donde la estructura JSON actúa como una lista serializada de llamadas a métodos de la
librería. Proyectos como `php-json-escpos` 1 y `js-escpos-builder` 3 sugieren este patrón.

Un documento JSON en este enfoque se vería así:

```json
 }[
[
{"cmd": "init"},
{"cmd": "text", "args": ["Texto de ejemplo"] },
{"cmd": "center"}, lo"] },
{"cmd": "text", "args": ["Bienvenido"] },
{"cmd": "barcode", "args": ["123456789", "EAN13"] },
{"cmd": "cut"}
]
```

**Ventajas:**

* **Implementación Simple:** El parser es trivial; se reduce a un `switch` o un `map[string]func` que itera sobre el
  array y ejecuta los comandos.
* **Generación Fácil:** Un frontend (aplicación web) puede generar este JSON sin necesidad de conocer la lógica de un
  recibo, simplemente exponiendo los comandos de la impresora.

**Desventajas:**

* **Alto Acoplamiento:** La plantilla está íntimamente acoplada a la API de la librería de impresión. Si la librería
  cambia el nombre de un método (ej., de `setTextSize` a `textSize`), todas las plantillas JSON existentes se rompen.
* **Baja Mantenibilidad:** La plantilla es imperativa (dice "cómo" hacerlo), no declarativa (dice "qué" se quiere). Esto
  la hace difícil de leer y mantener, un problema común cuando se usan lenguajes de marcado para definir lógica de
  aplicación.4

#### 1.1.2. Enfoque 2: Capa de Abstracción Semántica (JSON-a-Abstracción)

Este es un patrón significativamente más robusto. El JSON define un documento semántico, describiendo qué debe
imprimirse, no cómo se implementan los comandos de la impresora.

El repositorio `grandchef/escpos-template` 6 es un caso de estudio perfecto de este enfoque. El análisis de su
implementación 6 confirma que es una "capa de abstracción". El JSON no contiene nombres de comandos ESC/POS, sino claves
semánticas de alto nivel:

```json
[
  {
    "items": "coupon.title",
    "align": "center",
    "style": "bold+",
    "width": "2x"
  },
  {
    "type": "qrcode",
    "data": "https://github.com/...",
    "align": "right"
  },
  {
    "type": "image",
    "data": "picture.image",
    "align": "center"
  }
]
```

**Ventajas:**

* **Bajo Acoplamiento:** La plantilla es declarativa y legible por humanos. Describe la intención del documento.
* **Alta Mantenibilidad:** La librería subyacente (en el caso de `grandchef`, es `escpos-buffer` 6) puede ser
  reemplazada o actualizada. El procesador de plantillas simplemente necesita ser actualizado para mapear las claves
  semánticas (`align: 'center'`) a los nuevos comandos de la librería, sin que ninguna plantilla de usuario final se vea
  afectada.

**Desventajas:**

* **Mayor Complejidad:** Requiere un procesador de plantillas más sofisticado que entienda la semántica y gestione la
  traducción de "documento" a "comandos".

#### 1.1.3. Enfoque 3: Lenguaje Específico de Dominio (DSL)

Este es el nivel más alto de abstracción. En lugar de JSON, la entrada es un lenguaje de marcado personalizado, diseñado
específicamente para la creación de recibos.

El caso de estudio principal es `receiptline` 7, que implementa un "Markdown for receipts". Este proyecto no es un
simple parser de JSON; es un procesador de lenguaje completo.

```
{width: 2x; style: bold}
| Título Centrado
{align: right}
| {qrcode: https://github.com/...}
```

**Ventajas:**

* **Máxima Abstracción y Potencia:** El DSL puede ser diseñado para ser extremadamente conciso y expresivo.
* **Independencia del Backend:** Como se detalla en el análisis de `receiptline` 7, su método `transform()` puede
  generar no solo comandos `escpos`, sino también imágenes vectoriales (`svg`). Esto desacopla completamente el diseño
  del recibo del hardware de salida.
* **Accesibilidad:** Con la documentación adecuada, un usuario no técnico (como un gerente de tienda) podría editar el
  diseño de un recibo sin necesidad de un desarrollador.

**Desventajas:**

* **Complejidad de Implementación:** Es, con diferencia, el enfoque más complejo. Requiere escribir y mantener un
  parser (con su léxico y gramática) y un evaluador que traduzca el árbol de sintaxis abstracta (AST) a comandos de
  impresora.

#### 1.1.4. Enfoque 4: API Programática Fluida (El "Builder")

Este enfoque no es un formato de plantilla, sino la base programática sobre la cual se construyen los otros tres
enfoques. Es el método preferido por librerías maduras como `mike42/escpos-php` 8 y `python-escpos` 9 para la
construcción de recibos directamente en el código.

Un ejemplo en Go (basado en el estilo de `escpos-php` 8) se vería así:

```go
p.Init()
p.SetAlign(escpos.AlignCenter).
SetTextSize(2, 2).
SetBold(true).
Text("Título del Recibo\n")
p.Barcode("123456789", escpos.BarcodeEAN13)
p.Cut()
p.End()
```

**Ventajas:**

* **Máximo Control y Rendimiento:** No hay overhead de parsing de JSON o DSL.
* **Seguridad de Tipos:** En Go, esto proporciona verificación en tiempo de compilación de los comandos y sus
  argumentos.
* **Alta Mantenibilidad (Técnica):** El flujo del recibo es código explícito, fácil de depurar y mantener para un
  desarrollador.

**Desventajas:**

* **Flexibilidad Nula para el Usuario Final:** Cualquier cambio en el recibo, por pequeño que sea (incluso un error
  tipográfico), requiere una modificación del código, una recompilación y un redespliegue de la aplicación.

### 1.2. Análisis Comparativo: Robustez y Mantenibilidad

La elección del enfoque depende directamente de quién tiene la responsabilidad de mantener las plantillas de recibos.
Los sistemas POS personalizables 10 deben equilibrar la facilidad de uso para el usuario final con la robustez técnica
para el desarrollador.

**Tabla 1: Matriz de Comparación de Enfoques de Plantillas**

| Enfoque                               | Nivel de Abstracción | Mantenibilidad (Técnica)                       | Flexibilidad (Usuario Final)                         | Robustez (Manejo de Errores)                                 | Caso de Uso Ideal                                                 |
|:--------------------------------------|:---------------------|:-----------------------------------------------|:-----------------------------------------------------|:-------------------------------------------------------------|:------------------------------------------------------------------|
| **Mapeo Directo (JSON-Comando)**      | Muy Baja             | Pobre. Acoplado a la API de la librería.4      | Pobre. Requiere que el usuario conozca los comandos. | Pobre. Errores tipográficos en "cmd" fallan en runtime.      | Prototipos rápidos; GUIs que generan JSON.                        |
| **Abstracción Semántica (JSON-Doc)**  | Media                | Excelente. Desacoplado de la API.6             | Buena. El JSON es legible y semántico.               | Buena. El procesador puede validar la semántica.             | Sistemas POS (SaaS) donde las plantillas se almacenan en BBDD.    |
| **DSL (ej. `receiptline`)**           | Alta                 | Alta. Desacoplado de la API.7                  | Excelente. El DSL es conciso y potente.              | Excelente. Un parser proporciona errores de sintaxis claros. | Sistemas que necesitan backends múltiples (ej. `svg` y `escpos`). |
| **API Programática (Fluent Builder)** | Nula (Es la API)     | Excelente. Código Go nativo, seguro en tipos.8 | Nula. Requiere recompilación.                        | Excelente. Errores en tiempo de compilación.                 | Backend de la librería; POS para un cliente único.                |

### 1.3. Recomendación Estratégica: El Enfoque Preferido en POS del Mundo Real

El análisis de los sistemas POS del mundo real 10 muestra una tensión: los desarrolladores que construyen un sistema
para un cliente específico prefieren la velocidad y seguridad de la API Programática 8; los sistemas POS SaaS (Software
as a Service) deben ofrecer plantillas dinámicas, donde la Abstracción Semántica (JSON) 6 es el enfoque más robusto y
mantenible.

Una librería "definitiva" no debe forzar una elección. Debe soportar ambos casos de uso mediante una arquitectura de
capas:

* **Núcleo (Core):** El núcleo de la librería debe exponer la API Programática Fluida (Enfoque 4). Esta es la base de
  alto rendimiento, segura en tipos y mantenible.
* **Paquetes Opcionales:** Los motores de plantillas deben implementarse como paquetes separados y opcionales que
  consumen el núcleo.
    * `parsers/json`: Implementaría el enfoque de Abstracción Semántica (Enfoque 2). Su función
      `Parse(jsonData, dataObject)` importaría la librería core y traduciría el JSON en una serie de llamadas a la API
      fluida (ej. `core.Text()`, `core.QR()`).
    * `parsers/dsl`: Implementaría el enfoque de DSL (Enfoque 3).

Esta arquitectura ofrece la "mantenibilidad" de la API para los desarrolladores y la "flexibilidad" de las plantillas
para los sistemas SaaS, cumpliendo con todos los requisitos del mercado POS.

-----

## 2\. Análisis Competitivo: Librerías ESC/POS en Go

Para construir una librería "definitiva", primero debemos identificar las debilidades estratégicas y los vacíos
arquitectónicos en el ecosistema de Go existente.

### 2.1. Identificación del "Rival a Vencer": kenshaw/escpos

Con 254 estrellas y 84 forks 14, `kenshaw/escpos` es el líder de facto en el ecosistema Go. Es la librería más visible y
la que, con mayor probabilidad, un desarrollador de Go encontrará primero.

* **Actividad y Mantenimiento:** El repositorio incluye una lista TODO prominente en su README, lo que indica que el
  mantenedor es consciente de sus deficiencias, pero que estas persisten.14
* **Análisis Arquitectónico (E/S):** El diseño de E/S de la librería es idiomático y correcto para Go. El análisis de su
  código de ejemplo 14 muestra que el constructor `escpos.New(w)` se inicializa con un `*bufio.Writer`.14 Dado que
  `*bufio.Writer` acepta cualquier tipo que implemente `io.Writer`, la librería depende correctamente de la interfaz
  estándar de E/S de Go. Esto permite una excelente inyección de dependencias para la E/S, permitiendo que la librería
  escriba en archivos (`*os.File`), conexiones de red (`net.Conn`) o buffers en memoria (`bytes.Buffer`) sin
  modificación.

**Crítica y Vulnerabilidad Estratégica:**

La principal vulnerability de `kenshaw/escpos` es explícita: su propia lista TODO declara la necesidad de "Fix
barcode/image support".14

Esto no es un bug menor; es un fallo en dos de las características más esenciales de un sistema POS moderno. Un recibo
que no puede mostrar de forma fiable el logotipo de la empresa o el código de barras de una transacción está incompleto.

Esta debilidad ha forzado la fragmentación del ecosistema. Los desarrolladores que usan `kenshaw/escpos` se ven
obligados a buscar soluciones de terceros para llenar estos vacíos, como `boombuler/barcode` para códigos de barras 15 o
`skip2/go-qrcode` para QR. Otros pueden incluso abandonar Go para esta tarea, optando por soluciones más maduras como
`node-escpos`.16

`kenshaw/escpos` es un "rey" vulnerable. Su popularidad le da visibilidad, pero su funcionalidad incompleta crea una
oportunidad clara y definida para un retador que ofrezca una solución "completa en características" desde el primer
momento.

### 2.2. Los Retadores Emergentes

La debilidad de `kenshaw/escpos` ha creado un espacio para retadores.

#### 2.2.1. hennedo/escpos

* **Análisis:** Con 98 estrellas y 42 forks 17, este es el "retador directo". Su existencia parece motivada precisamente
  por las deficiencias de `kenshaw`.
* **Fortaleza:** Su README 17 y su conjunto de características 17 lista explícitamente "Image Printing", "Barcodes" (
  UPC-A, UPC-E, EAN13, EAN8) y "QR Codes" como características funcionales. Resuelve directamente la vulnerabilidad
  estratégica de `kenshaw`.
* **Debilidad Arquitectónica:** Aunque el constructor `New()` probablemente acepta una interfaz genérica
  `io.ReadWriteCloser`, el único ejemplo de uso proporcionado en su README 17 utiliza `net.Dial("tcp",...)`directamente.
  Esto sugiere un enfoque menos idiomático en su documentación, lo que podría llevar a los usuarios a un acoplamiento
  más fuerte con la E/S de red.

#### 2.2.2. conejoninja/go-escpos

* **Análisis:** Esta es una librería de nicho (16 estrellas 18) pero con una arquitectura conceptualmente superior.
* **Fortalezas 19:**
    * **Abstracción de E/S Perfecta:** Su constructor base es `NewPrinterByRW(rwc io.ReadWriteCloser)`.19 Este es el
      diseño de E/S idiomático de Go ideal, construido sobre la interfaz estándar `io.ReadWriteCloser`.
    * **Soporte de TinyGo:** La librería está diseñada para funcionar con TinyGo, lo que la hace ideal para
      microcontroladores y sistemas embedded.19
    * **Descubrimiento de USB:** Ofrece `NewUSBPrinterByPath("")` para el descubrimiento automático de dispositivos USB
      19, una característica de conveniencia muy atractiva.
* **Debilidad:** El descubrimiento de USB es una "fuga de abstracción"; depende de `udev` y, por lo tanto, solo funciona
  en Linux, limitando la portabilidad de esta característica.19

### 2.3. Síntesis y Oportunidad Estratégica: El Vacío en el Ecosistema Go

El ecosistema actual de Go para ESC/POS está fragmentado y en un ciclo de "parcheo de características". `kenshaw` tiene
la popularidad pero le faltan características críticas. `hennedo` añade esas características pero tiene una
documentación de E/S más débil. `conejoninja` tiene la mejor arquitectura de E/S pero la menor adopción.

Sin embargo, todas estas librerías de Go cometen el mismo error fundamental: ignoran el problema real y más difícil de
la impresión POS: la **compatibilidad de modelos**. Sus repositorios están repletos de hacks y comandos codificados (
hardcoding) para un modelo específico de Epson (ej. TM-T82, TM-T20II 14).

Este es el verdadero vacío en el mercado. La ventaja decisiva no vendrá de simplemente implementar códigos QR (lo cual
es necesario, pero no suficiente). La superioridad objetiva se alcanzará implementando un sistema de **Perfiles de
Capacidad (Capability Profiles)**, un concepto del que el ecosistema Go carece por completo, pero que es fundamental
para las soluciones maduras en otros lenguajes.

-----

## 3\. Benchmarking: Sistemas POS de Código Abierto (Multi-lenguaje)

Al expandir la investigación más allá de Go, podemos aprender de los ecosistemas maduros (Python, PHP, Java) que han
estado resolviendo estos problemas durante más de una década.

### 3.1. Abstracción de Impresión en Sistemas POS (Odoo, OSPOS, Floreant)

#### 3.1.1. Odoo (Python-based)

Odoo 20 es un ERP/POS masivo basado en Python. Su arquitectura de impresión 21 es muy reveladora de los desafíos del
mundo real. Odoo distingue entre dos tipos de impresoras:

* **Impresoras ePOS:** Impresoras de red que exponen un SDK de Javascript/XML. La impresión ocurre en el cliente (el
  navegador).21
* **Impresoras ESC/POS:** Impresoras tradicionales (USB, Serial, Red).

Para que el POS basado en web (servidor Odoo) se comunique con una impresora ESC/POS local (USB/Serial), Odoo requiere
un intermediario: el "IoT system" o "PosBox".21

Esta "PosBox" 24 es un dispositivo de hardware (típicamente un Raspberry Pi) que actúa como un puente de microservicio.
El POS web envía una solicitud de impresión (probablemente JSON semántico) al PosBox. El PosBox, que ejecuta Python o un
binario similar, recibe esta solicitud, la traduce a comandos ESC/POS y la reenvía al hardware físico (USB, Serial) al
que está conectado.

* **Implicación:** Este modelo de "puente de impresión" es un caso de uso perfecto para una librería Go. Un binario Go
  compilado estáticamente, de alto rendimiento y bajo consumo de memoria, es el ejecutable ideal para desplegar en un "
  IoT Box".

#### 3.1.2. Open Source POS (OSPOS) (PHP-based)

OSPOS es una aplicación POS basada en web escrita en PHP.20 Enfrenta el mismo desafío arquitectónico: el servidor PHP no
puede acceder directamente al `/dev/usb/lp0` del cliente.26

La solución, al igual que con Odoo, es un servicio de impresión intermediario o el uso de una librería del lado del
servidor (como `mike42/escpos-php` 8) para generar los comandos que luego se envían al cliente para su impresión local.

#### 3.1.3. Floreant POS (Java-based)

Floreant es una aplicación de escritorio Java.28 Al ser una aplicación local, evita el problema del puente de red.

Su arquitectura distingue entre tipos de impresoras (Recibo, Cocina) y tamaños de papel (80mm, 76mm) 31, lo que implica
una capa de configuración de impresión que maneja diferentes drivers y formatos.

### 3.2. La Característica Inspiradora N° 1: El Ecosistema `escpos-printer-db`

Este es el descubrimiento más crítico de esta investigación. Las librerías más maduras y robustas, `python-escpos` 9 y
`escpos-php` 8, no son proyectos monolíticos. Son clientes de un ecosistema de compatibilidad compartido.

**El Eslabón Perdido: `receipt-print-hq/escpos-printer-db`.33**

* Este es un repositorio independiente y agnóstico al lenguaje.

* Contiene una base de datos mantenida por la comunidad sobre las capacidades de cientos de modelos de impresoras
  ESC/POS.33

* Es consumido activamente por las librerías de Python y PHP.8

* **Formato:** La base de datos está en formato JSON/YAML.39 El código fuente de `python-escpos` 41 confirma que carga y
  parsea un archivo `capabilities.json` o `capabilities.yaml`.

* **Esquema:** Aunque una inspección directa del esquema no fue posible 34, el análisis de los issues 43 y los commits
  34 revela las claves que define:

    * `codePages`: Mapea nombres de codificación (ej. Cirílico) a los comandos ESC/POS necesarios para activarlos.43
    * `paperFullCut`, `paperPartCut`: Define los comandos específicos para diferentes tipos de corte.34
    * Variantes de comandos para códigos de barras, imágenes y otras características específicas del fabricante.

* **Acción Estratégica:**

    * **No Reinventar:** La nueva librería Go no debe intentar crear su propia base de datos de compatibilidad.
    * **Consumir:** La librería debe ser diseñada para consumir el archivo `capabilities.json/yaml` de
      `escpos-printer-db`.33
    * Esto requiere un nuevo módulo (`escpos-capabilities`) que parsee este archivo de base de datos en `structs` de Go
      nativos.

* **Ventaja Competitiva:** Al hacer esto, la librería heredará instantáneamente la compatibilidad con cientos de modelos
  de impresoras, superando a todas las demás librerías de Go en el mercado desde el día de su lanzamiento.

### 3.3. La Característica Inspiradora N° 2: Robustez de Características (Imagen/QR)

* **Procesamiento de Imágenes:** `python-escpos` depende de Pillow (PIL) 9, una librería de procesamiento de imágenes de
  nivel industrial. Esto valida la estrategia de crear un módulo `escpos-graphics` superior con dithering avanzado; es
  el enfoque probado por las librerías maduras.

* **Códigos QR y Fallbacks:** La librería de Dart `esc_pos_utils` 40 demuestra una robustez de siguiente nivel. Su
  generador de QR intenta dos métodos:

    1. **Nativo:** Intenta usar los comandos ESC/POS nativos de la impresora para generar QR:
       `generator.qrcode('example.com')`.
    2. **Fallback (Imagen):** Si la impresora no soporta QR nativo (una capacidad que se conocería a través de un
       perfil), la librería tiene un fallback: usa un paquete externo (`qr_flutter`) para renderizar el QR como una
       imagen y luego imprime esa imagen.40

* **Implicación:** Una librería "definitiva" debe manejar estos fallbacks de forma elegante. Esto solo es posible si un
  `CapabilityProfile` (cargado desde `escpos-printer-db`) informa a la librería si la impresora soporta `QR_NATIVE` o si
  debe usar `QR_FALLBACK_IMAGE`.

-----

## 4\. Síntesis y Propuesta Estratégica

Basado en el análisis de los enfoques de plantillas, las debilidades del ecosistema Go y las fortalezas de los
ecosistemas maduros, se define el siguiente plan arquitectónico para la nueva librería.

### 4.1. Definición de "Superioridad" Objetiva

Para que esta librería sea objetivamente superior, debe ser:

1. **Arquitectónicamente Desacoplada:** Diseñada como un meta-paquete con módulos interoperables y una clara segregación
   de interfaces.45
2. **Impulsada por Compatibilidad:** Ser la primera y única librería Go en consumir la base de datos `escpos-printer-db`
   33, resolviendo el problema central de la compatibilidad de modelos 50 y la codificación de caracteres.43
3. **Completa en Características:** Implementar de forma robusta y con fallbacks (al estilo de 40) las características
   que faltan en `kenshaw/escpos` 14: imágenes avanzadas, códigos de barras y códigos QR.
4. **Verificablemente Interoperable:** Cumplir con el requisito principal de que el módulo de gráficos sea un paquete Go
   puro, utilizable por cualquier otra librería.
5. **Robusta (Goroutine-Safe):** Ser segura para su uso en servidores Go concurrentes, manejando el bloqueo de
   dispositivos 53 y la gestión de estado de la impresora.17

### 4.2. Propuesta de Arquitectura Modular (Meta-Paquete)

Se recomienda una estructura de monorepo/meta-paquete 49 que segregue las responsabilidades.

**`github.com/your-org/escpos` (el `go.mod` raíz)**

* **`core/` (o `escpos/`):** El corazón de la librería.

    * Define las interfaces clave (el "contrato" de la librería 45).
      ```go
      // Connection define el transporte de E/S.
      type Connection interface {
      io.ReadWriteCloser
      // Lock/Unlock para control de concurrencia.
      Lock()
      Unlock()
      }

          // GraphicsProcessor define el pipeline de imágenes.
          type GraphicsProcessor interface {
              Process(img image.Image) (*graphics.MonochromeImage, error)
          }
          ```
    * Contiene la API Fluida principal: `type Printer struct {... }`.
    * Contiene los métodos de la API: `func (p *Printer) Text(s string) *Printer`, `func (p* Printer) Cut() error`.
    * Maneja la lógica de estado de la impresora (ver 4.4.1).

* **`capabilities/`**

    * El parser de `escpos-printer-db`.33
    * `func LoadProfile(name string) (*Profile, error)`
    * `type Profile struct {... }` (contiene los comandos de corte, páginas de códigos, formatos de imagen, etc.,
      parseados desde el YAML/JSON 41).
    * El `core.Printer` será inyectado con este perfil: `p.SetProfile(*Profile)`.

* **`graphics/` (El Módulo Interoperable)**

    * Un paquete Go puro. No importa `core/`.
    * Contiene algoritmos de dithering (ej. Floyd-Steinberg, Atkinson) y conversión.
    * `func Dither(img image.Image) (*MonochromeImage, error)`
    * Cumple el requisito de interoperabilidad (ver 4.3).

* **`drivers/`**

    * Implementaciones de la interfaz `core.Connection`.
    * `drivers/usb`: (usando `gousb` o `karalabe/hid`).
    * `drivers/net`: Un wrapper sobre `net.Conn`.
    * `drivers/serial`: (usando `go.bug.st/serial`).
    * `drivers/file`: Un wrapper sobre `*os.File` que implementa `Lock()`.
    * `drivers/mock`: Para pruebas, escribe en un `bytes.Buffer`.

* **`parsers/`**

    * Implementaciones de los motores de plantillas (Sección 1).
    * `parsers/json`: Importa `core/`. Parsea el JSON semántico 6 y llama a la API fluida.
    * `parsers/receiptline`: Importa `core/`. Parsea el DSL 7 y llama a la API fluida.

### 4.3. Caso de Uso Específico: Diseño de `escpos-graphics` para Interoperabilidad

Este es el requisito de interoperabilidad de alta prioridad. Se logra mediante el desacoplamiento total, siguiendo un
patrón de "utilidad pura".45 El error a evitar es hacer que `graphics/` dependa de `core/`. El módulo de gráficos
resuelve un problema matemático (conversión de imagen), no un problema de protocolo (comandos ESC/POS).

#### 4.3.1. El Contrato (Definido en `graphics/image.go`)

```go
package graphics

import "image"

// MonochromeImage es la representación pura de 1-bit por píxel
// que una impresora térmica entiende.
type MonochromeImage struct {
	Pix []byte // Buffer de 1-bit por píxel
	Stride int // Bytes por fila
	Rect image.Rectangle
}

// DitherOptions permite configurar el algoritmo.
type DitherOptions struct {
	Threshold uint8
	Algorithm DitherMode // ej. FloydSteinberg, Atkinson, etc.
}

// Process es una función pura. No tiene estado y no tiene
// dependencias de ningún otro paquete de nuestra librería.
// Solo depende del paquete 'image' estándar de Go.
func Process(img image.Image, opts DitherOptions) *MonochromeImage {
	// 1. Convertir a escala de grises.
	// 2. Aplicar el algoritmo de dithering (ej. Floyd-Steinberg).
	// 3. Empaquetar los bits en el formato MonochromeImage.
	//...
	return monoImage
}
```

#### 4.3.2. Consumo por un Competidor (ej. `kenshaw/escpos`)

El mantenedor de `kenshaw/escpos` ahora puede solucionar su bug de imágenes 14 importando este módulo:

```go
package kenshaw_escpos

import (
	"image"

	"github.com/your-org/escpos/graphics" // 1. Importa SU módulo
)

// PrintImage es el método hipotético de kenshaw
func (p *KenshawPrinter) PrintImage(img image.Image) {
	// 2. Llama a su función pura de dithering
	opts := graphics.DitherOptions{Algorithm: graphics.FloydSteinberg}
	monoImg := graphics.Process(img, opts)

	// 3. El mantenedor de kenshaw usa SU PROPIA lógica de comandos
	//    ESC/POS (GS v 0) para formatear los bytes de monoImg.Pix
	//    y enviarlos a la impresora.
	p.write(p.formatImageAs_GS_v_0(monoImg))
}
```

#### 4.3.3. Consumo por su `core.Printer`

Su librería hará lo mismo, pero utilizando la interfaz `GraphicsProcessor` para mantenerlo conectable:

```go
// en core/printer.go
import "github.com/your-org/escpos/graphics"

// DefaultProcessor es la implementación estándar
type DefaultProcessor struct{}

func (gp *DefaultProcessor) Process(img image.Image) (*graphics.MonochromeImage, error) {
opts := graphics.DitherOptions{Algorithm: graphics.FloydSteinberg}
return graphics.Process(img, opts), nil
}

//... en el constructor de NewPrinter()
p.gfx = &DefaultProcessor{} // Se establece el procesador por defecto

// Image es el método de la API fluida
func (p *Printer) Image(img image.Image) error {
// 1. Llama a la interfaz
monoImg, err := p.gfx.Process(img)
if err != nil { return err }

// 2. Llama al comando de formato de imagen
//    definido por el PERFIL DE CAPACIDAD cargado.
cmdBytes, err := p.profile.FormatImage(monoImg)
if err != nil { return err }

// 3. Escribe en la conexión
return p.conn.Write(cmdBytes)
}
```

Este diseño cumple el requisito de interoperabilidad al 100%. `escpos-graphics` se convierte en una librería Go líder en
su clase para el procesamiento de imágenes de recibos, que cualquier proyecto puede importar y utilizar.

### 4.4. Puntos Ciegos: Desafíos Críticos a Anticipar

La implementación de una librería ESC/POS está llena de trampas. Una arquitectura robusta debe anticiparlas.

#### 4.4.1. El Infierno de la Gestión de Estado

* **Problema:** ESC/POS es un protocolo con estado. Si se llama a `SetBold(true)`, la impresora permanece en negrita
  indefinidamente. Si otro proceso (u otra goroutine) imprime en la misma impresora, heredará un "estado sucio",
  resultando en recibos mal impresos.
* **Análisis:** El comando `ESC @` (0x1B 0x40) es el comando de "Initialize printer".55 Borra el búfer y resetea todos
  los modos (negrita, tamaño, etc.) a sus valores de fábrica.59
* **Recomendación Arquitectónica (Semántica Sin Estado):** La librería debe forzar una semántica sin estado (stateless)
  por defecto para garantizar trabajos de impresión atómicos y predecibles. La librería `hennedo/escpos` 17 sugiere esto
  al "establecer todos los parámetros de estilo nuevamente para cada llamada de Write".
* **Implementación:** El `core.Printer` no debe enviar comandos de estilo (`SetBold`, `SetAlign`) inmediatamente. Debe
  acumularlos como estado interno (ej. `p.bold`, `p.align`). Cuando se llama a un método de finalización (`Print()`,
  `Cut()`, `End()`), la librería debe:
    1. Crear un búfer de comandos temporal.
    2. Anteponer (prepend) a este búfer la secuencia de inicialización completa: `ESC @`, seguido de los comandos para
       establecer el estado deseado (ej. `ESC E 1` si `p.bold == true`, `ESC a 1` si `p.align == AlignCenter`).
    3. Añadir el contenido del usuario (texto, imágenes) al búfer.
    4. Enviar el búfer completo a la `Connection` de una sola vez.
* **Ventaja:** Esto garantiza que cada trabajo de impresión sea atómico, predecible y no se vea afectado por el estado
  anterior de la impresora.

#### 4.4.2. Concurrencia y Bloqueo de Dispositivos

* **Problema:** Un servidor Go maneja miles de goroutines. Si dos goroutines intentan escribir en el mismo descriptor de
  archivo (ej. `/dev/usb/lp0` 54 o `COM1`) simultáneamente, el sistema operativo puede devolver un error de "device or
  resource busy" 53, o, peor aún, los datos de los dos recibos se intercalarán, produciendo basura ilegible.
* **Recomendación Arquitectónica:** La librería debe ser goroutine-safe a nivel de E/S.
* **Implementación:** Las implementaciones concretas de la interfaz `core.Connection` (en `drivers/`) deben contener un
  `sync.Mutex` interno. La interfaz `Connection` debe exponer los métodos `Lock()` y `Unlock()` para que el
  `core.Printer` pueda gestionar transacciones de impresión completas.

    ```go
    type UsbConnection struct {
        mu sync.Mutex
        f  *os.File // o el descriptor de dispositivo
    }

    func (c *UsbConnection) Write(data []byte) (int, error) {
        c.mu.Lock()
        defer c.mu.Unlock()
        return c.f.Write(data)
    }

    // El método Print del core.Printer usaría esto:
    func (p *Printer) Print() error {
        p.conn.Lock()
        defer p.conn.Unlock()

        buffer := p.buildBuffer() // Construye el búfer con estado (ver 4.4.1)
        _, err := p.conn.Write(buffer.Bytes())
        return err
    }
    ```

* **Ventaja:** Esto protege al usuario de la librería de los errores de concurrencia a nivel de dispositivo, un
  diferenciador crucial para una librería de nivel de servidor.

#### 4.4.3. El Infierno de la Codificación de Caracteres

* **Problema:** Go utiliza UTF-8. Las impresoras ESC/POS utilizan "páginas de códigos" (codepages) de 8 bits (ej. CP437,
  CP1250, Katakana).63 Imprimir caracteres fuera de ASCII (ej. "€", "ñ", "ü") o texto en Árabe 52 o Cirílico 43 fallará,
  imprimiendo caracteres incorrectos.
* **Solución:** Este problema se resuelve únicamente mediante la estrategia de `escpos-printer-db` (Sección 3.2). El
  `Profile` cargado 41 le dirá a la librería qué páginas de códigos están disponibles y los comandos ESC/POS para
  activarlas.
* **Implementación:** El método `core.Printer.Text(s string)` debe:
    1. Consultar al `p.profile` la página de códigos activa (ej. CP858).
    2. Utilizar un paquete de transcodificación (como `golang.org/x/text/encoding`) para convertir el string UTF-8 de Go
       a los bytes de la página de códigos de destino.
    3. Enviar los comandos para activar esa página de códigos, seguidos de los bytes transcodificados.
* **Fallback:** Si un carácter (ej. un emoji o un ideograma Chino) no existe en las páginas de códigos de la impresora,
  la librería debe implementar el fallback robusto 52: renderizar ese texto como una imagen y enviarlo a la impresora.

#### 4.4.4. Compatibilidad de Protocolo (Epson vs. Star vs. Bixolon)

* **Problema:** No todo lo que se autodenomina "ESC/POS" es 100% compatible con Epson. Los fabricantes como Star 50,
  Bixolon 51 y otros a menudo usan sus propios dialectos (ej. "Star Mode" 50) para funciones avanzadas como códigos de
  barras o corte de papel.
* **Solución:** Este problema ya está resuelto si se sigue la estrategia de diseño. El `CapabilityProfile` es la
  solución.
* **Implementación:** La API fluida (`p.Cut()`) permanece idéntica para el usuario. Internamente, `p.Cut()` no contiene
  bytes codificados. En su lugar, hace:

    ```go
    func (p *Printer) Cut() error {
        // Obtiene la secuencia de bytes para "corte"
        // desde el perfil cargado.
        cmdBytes, err := p.profile.GetCommand("CutPartial")
        if err != nil {
            // El perfil de esta impresora no soporta corte.
            return err 
        }
        return p.conn.Write(cmdBytes)
    }
    ```

* **Ventaja:** El perfil de una Epson TM-T88 devolverá `[0x1D, 0x56, 0x42, 0x00]` para `CutPartial`. El perfil de una
  impresora Star podría devolver una secuencia completamente diferente. El usuario no lo sabe y no le importa. La
  librería simplemente "hace lo correcto", logrando la máxima compatibilidad de hardware con un mínimo esfuerzo del
  desarrollador.
