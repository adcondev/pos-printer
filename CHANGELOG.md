# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

## [1.0.0](https://github.com/AdConDev/pos-daemon/compare/v0.2.0...v1.0.0) (2025-08-26)


### ⚠ BREAKING CHANGES

* **protocols:** add new architecture for command chaining and 2-layered commands.
* **escpos:** The protocol interface has been modified
to support multiple protocols and may require updates to
existing implementations.

Signed-off-by: Adrián Constante <ad_con.reload@proton.me>

### 🐛 Bug Fixes

* **errors:** standardize error variable names ([b865304](https://github.com/AdConDev/pos-daemon/commit/b865304a04e7079ca09a08bbafbb4fb00528995e))
* **escpos:** improve comments and code clarity ([bbf654c](https://github.com/AdConDev/pos-daemon/commit/bbf654c2e3d705af8e7825f1836bd39fd51c5673))


### ✅ Tests

* **escpos:** add tests for dependency injection functionality ([3d7c680](https://github.com/AdConDev/pos-daemon/commit/3d7c680391865a44ea79235cc70edd331b17ea72))
* **escpos:** add tests for line spacing functionality ([6bbfb2a](https://github.com/AdConDev/pos-daemon/commit/6bbfb2a9d8c7ace7b691e005d6a62b606ba5e9c0))
* **escpos:** add unit tests for command functionalities ([d759a33](https://github.com/AdConDev/pos-daemon/commit/d759a33be3c6cffd3e74e9b8601ec0ea86da894e))


### ✨ Features

* **barcode:** update barcode handling functions ([85558e7](https://github.com/AdConDev/pos-daemon/commit/85558e790d037f74c84e06a0aa8aa1ca0d213c30))
* **escpos:** add line spacing capabilities and refactor commands ([b0a0f84](https://github.com/AdConDev/pos-daemon/commit/b0a0f84e499d90ac6769ebc7491916553120202a))
* **escpos:** enhance printer command structure and add comments ([dd1c733](https://github.com/AdConDev/pos-daemon/commit/dd1c7333bf0e9b8dc58c7cc9136108d031ed0b58))
* **escpos:** refactor printer protocol handling ([9ce0903](https://github.com/AdConDev/pos-daemon/commit/9ce09039be23b004c1e282e8d09efd522c6d1129))
* **escpos:** refactor protocol structure and update imports ([f1840f8](https://github.com/AdConDev/pos-daemon/commit/f1840f87ef9b3cedb1f519184b53cb84bcc1dd30))
* **printer:** enhance printer configuration structure ([53c4c9c](https://github.com/AdConDev/pos-daemon/commit/53c4c9ccc93b1e33c6c5e2e27a8626af95a156bd))
* **protocols:** add new architecture for command chaining and 2-layered commands. ([599214f](https://github.com/AdConDev/pos-daemon/commit/599214f87e55896323056e47aa919776b2513d36)), closes [#4](https://github.com/adcondev/pos-printer/issues/4)
* **protocol:** update import paths for escpos types ([599aca6](https://github.com/AdConDev/pos-daemon/commit/599aca6982e17ce3a83902b7ddf449a3c34b1d18))

## 0.2.0 (2025-08-15)

### 🤖 Continuous Integration

* bump amannn/action-semantic-pull-request from 5 to
  6 ([648be79](https://github.com/AdConDev/pos-daemon/commit/648be7999f29327db7bee9bbad30874ae27cbc64))
* bump codecov/codecov-action from 4 to
  5 ([3ce0298](https://github.com/AdConDev/pos-daemon/commit/3ce0298273748a58a796e0c90382bb9e3bc585e5))

### ✨ Features

* **escpos:** add initial implementation for ESC/POS
  commands ([f9772b4](https://github.com/AdConDev/pos-daemon/commit/f9772b47c1e4e2f8cd11910817250ef45ac472ca))
* **github:** add initial github workflows and
  files ([812b851](https://github.com/AdConDev/pos-daemon/commit/812b8513d31c12bb2eb240eb551d68bf9708c8e6))
