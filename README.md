# POS Printer: A Modular Go Library for Thermal Printers

<div align="center">

![CI Status](https://github.com/adcondev/pos-printer/actions/workflows/ci.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/adcondev/pos-printer)](https://goreportcard.com/report/github.com/adcondev/pos-printer)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/adcondev/pos-printer)

<br>

<img src="assets/images/pos-printer.jpg" alt="POS Printer Logo" width="200" height="auto">

**A modular, extensible library for thermal printer communication across multiple protocols**

</div>

## 📝 Overview

**POS Printer** is a Go library that provides a unified, developer-friendly interface for communicating with a wide range of POS (Point of Sale) printer models. Designed with a clean and modular architecture, it simplifies the process of sending commands to thermal printers, whether you need to print text, generate QR codes, or render images. With comprehensive test coverage and support for multiple connection types, this library is built for reliability in production environments.

## ✨ Key Features

- **Multi-Protocol Support**: Native support for **ESC/POS** and **ZPL**, with an extensible architecture to add more.
- **Flexible Connection Options**: Connect to printers via **Serial**, **USB**, **Network (TCP/IP)**, or **Bluetooth**.
- **Protocol-Agnostic Image Printing**: A powerful imaging package that handles direct bitmap generation from image files.
- **Centralized Printer Registry**: Manage multiple printer configurations in a centralized system, making it easy to switch between devices.
- **Comprehensive Testing**: Includes a full suite of unit tests, mocks, and fakes to ensure stable and predictable behavior.
- **Well-Documented and Idiomatic Go**: Clean, commented code that follows modern Go best practices.

## 🏗️ Architecture Diagram

The library is designed with a decoupled architecture, where each core functionality is separated into its own package. This allows for easy extension and maintenance.

```mermaid
graph TD
    A[Application] --> B{pos Printer API}

    subgraph "Core Library"
        B --> C[Printer Registry]
        B --> D[Connection Manager]
        B --> E{Protocols}
    end

    subgraph "Connection Layer"
        D --> F[USB]
        D --> G[Network]
        D --> H[Serial]
        D --> I[Bluetooth]
    end

    subgraph "Protocols Layer"
        E --> J[ESC/POS]
        E --> K[ZPL]
        E --> L[Imaging]
    end

    J --> M[Physical Printer]
    K --> M
    L --> M
```

## 🚀 Installation

To get started, add the library to your Go project using `go get`:

```bash
go get github.com/adcondev/pos-printer
```

## 💡 Usage Example

Here’s a simple example of how to connect to a printer and print a "Hello, World!" message.

```go
package main

import (
	"log"

	"github.com/adcondev/pos-printer/escpos"
	"github.com/adcondev/pos-printer/pos"
)

func main() {
	// 1. Create a new printer profile in the registry
	registry := pos.NewPrinterRegistry()
	printerProfile := &pos.Printer{
		Name: "MyReceiptPrinter",
		Device: &pos.Device{
			Connector: &pos.NetworkConnector{Address: "192.168.1.100:9100"},
		},
	}
	registry.Add("my_printer", printerProfile)

	// 2. Get the printer from the registry
	p, err := registry.Get("my_printer")
	if err != nil {
		log.Fatalf("Failed to get printer: %v", err)
	}

	// 3. Connect to the printer
	if err := p.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer p.Close()

	// 4. Send a command
	cmd := escpos.NewPrinter(p.Device.Connector)
	if err := cmd.Print("Hello, World!\n"); err != nil {
		log.Fatalf("Failed to print: %v", err)
	}

	log.Println("Successfully printed!")
}
```

## 🖨️ Supported Protocols

| Protocol | Status         | Description                                    |
|----------|----------------|------------------------------------------------|
| ESC/POS  | ✅ Stable       | Epson Standard Code for Point of Sale Printers |
| ZPL      | 🔄 In Progress | Zebra Programming Language for label printers  |
| Image    | ✅ Stable       | Direct bitmap generation for any printer       |

## 🤝 Contributing

Contributions are welcome! If you'd like to help improve the library, please feel free to fork the repository, make your changes, and submit a pull request. For major changes, please open an issue first to discuss what you would like to change.

1.  **Fork** the repository.
2.  Create your feature branch (`git checkout -b feature/AmazingFeature`).
3.  Commit your changes (`git commit -m 'Add some AmazingFeature'`).
4.  Push to the branch (`git push origin feature/AmazingFeature`).
5.  Open a **Pull Request**.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<div align="center">
  <sub>Built with ❤️ by adcondev</sub>
</div>
