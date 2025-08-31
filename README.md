# POS Printer

<div align="center">

![POS Printer Logo](https://via.placeholder.com/150?text=POS+Printer)
![CI Status](https://github.com/AdConDev/pos-daemon/actions/workflows/ci.yml/badge.svg)

**A modular, extensible library for thermal printer communication across multiple protocols**

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

</div>

## 📝 Overview

POS Printer is a Go library that provides a unified interface for communicating with various POS (Point of Sale) printer models using multiple protocols. It implements a clean, modular architecture with comprehensive test coverage for maximum reliability in production environments.

## ✨ Features

- **Multi-Protocol Support**: ESC/POS, ZPL, and more
- **Modular Architecture**: Each printer capability in its own package
- **Flexible Connection Options**: Serial, USB, Network, Bluetooth
- **Image Printing**: Direct bitmap generation from images
- **Registry System**: Centralized management of multiple printer configurations
- **Comprehensive Testing**: Unit tests, mocks, and fakes for reliable operation
- **Well-Documented Code**: Detailed comments for all commands and functions

## 🚀 Installation

```bash
go get github.com/adcondev/pos-printer
```

## 🏗️ Architecture

POS Printer is organized into modular packages:

```
pos-printer/
├── escpos/          # ESC/POS protocol implementation
│   ├── common/      # Shared constants and utilities
│   ├── print/       # Print commands
│   ├── linespacing/ # Line spacing commands
│   └── character/   # Character formatting commands
├── zpl/             # ZPL (Zebra) protocol implementation
├── image/           # Protocol-agnostic image printing
├── connection/      # Connection handling (USB, Network, etc.)
└── pos/             # Main package with printer registry
```

Each capability follows the same pattern:
- Interface definition for the capability
- Implementation struct with command methods
- Comprehensive tests with mocks and fakes

## 🖨️ Supported Protocols

| Protocol | Status | Description |
|----------|--------|-------------|
| ESC/POS  | ✅ Stable | Epson Standard Code for Point of Sale Printers |
| ZPL      | 🔄 In Progress | Zebra Programming Language for label printers |
| Image    | 🔄 In Progress | Direct bitmap generation for any printer |

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

---

<div align="center">
  <sub>Built with ❤️ by adcondev</sub>
</div>
