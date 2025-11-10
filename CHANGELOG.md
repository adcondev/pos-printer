# Changelog

All notable changes to the POS Printer library will be documented in this file.

### [3.0.2](https://github.com/adcondev/pos-printer/compare/v3.0.1...v3.0.2) (2025-11-10)


### 📦 Dependencies

* **npm:** bump @commitlint/cli from 19.8.1 to 20.1.0 ([#63](https://github.com/adcondev/pos-printer/issues/63)) ([1d19593](https://github.com/adcondev/pos-printer/commit/1d1959395768b428b075502cd21c2168e21918b5))

### [3.0.1](https://github.com/adcondev/pos-printer/compare/v3.0.0...v3.0.1) (2025-11-10)


### 🤖 Continuous Integration

* **dependabot:** enhance auto-merge workflow with timeout and success criteria ([5f74df6](https://github.com/adcondev/pos-printer/commit/5f74df60c908f8361b777fe9b4337200538a3052))
* **github:** add new scopes for gomod, npm, and gh-actions ([39746f4](https://github.com/adcondev/pos-printer/commit/39746f4dc5f920e163582d627415ad1fdcbf3285))


### 📦 Dependencies

* **npm:** bump @commitlint/config-conventional from 19.8.1 to 20.0.0 ([#62](https://github.com/adcondev/pos-printer/issues/62)) ([277a12c](https://github.com/adcondev/pos-printer/commit/277a12ca66553a20cbd94a61edca87cba494f1f9))

## [3.0.0](https://github.com/adcondev/pos-printer/compare/v2.3.0...v3.0.0) (2025-11-07)


### ⚠ BREAKING CHANGES

* **commands:** rename controllers to commands (#61)
* **escpos:** rename controllers to commands

### 📝 Documentation

* add LEARNING.md and update README.md ([21bc0f5](https://github.com/adcondev/pos-printer/commit/21bc0f525faee77b6a32df265fd6451622fec28c))


### 🐛 Bug Fixes

* update LEARNING.md ([99a7a85](https://github.com/adcondev/pos-printer/commit/99a7a858cce47007610f21c23448e355ced51204))


### ♻️ Code Refactoring

* **base64:** rename example file for clarity ([8d6533d](https://github.com/adcondev/pos-printer/commit/8d6533db418a2a58ede14fc245a3e285a3fb0d16))
* **commands:** improve command documentation and comments ([764a58f](https://github.com/adcondev/pos-printer/commit/764a58feafe8f4b68a87740fe77e2c608c5b4363))


### 🤖 Continuous Integration

* **github:** update workflow to dynamically find example directories ([cc643b8](https://github.com/adcondev/pos-printer/commit/cc643b88009eb569ed2e745d9c384827daf866d2))


### ✨ Features

* **commands:** rename controllers to commands ([#61](https://github.com/adcondev/pos-printer/issues/61)) ([fa786f3](https://github.com/adcondev/pos-printer/commit/fa786f3e5e445935ae2a193ffe15ac74f24c6272)), closes [/#diff-8f12d5f5467d5a9e7d05741bfaa3151879b70dc9b12ed2da94e21fcbbd80a986L14-L16](https://github.com/adcondev///issues/diff-8f12d5f5467d5a9e7d05741bfaa3151879b70dc9b12ed2da94e21fcbbd80a986L14-L16) [/#diff-8f12d5f5467d5a9e7d05741bfaa3151879b70dc9b12ed2da94e21fcbbd80a986L27-R69](https://github.com/adcondev///issues/diff-8f12d5f5467d5a9e7d05741bfaa3151879b70dc9b12ed2da94e21fcbbd80a986L27-R69) [/#diff-b803fcb7f17ed9235f1e5cb1fcd2f5d3b2838429d4368ae4c57ce4436577f03fL48-L50](https://github.com/adcondev///issues/diff-b803fcb7f17ed9235f1e5cb1fcd2f5d3b2838429d4368ae4c57ce4436577f03fL48-L50) [/#diff-b803fcb7f17ed9235f1e5cb1fcd2f5d3b2838429d4368ae4c57ce4436577f03fL61-R70](https://github.com/adcondev///issues/diff-b803fcb7f17ed9235f1e5cb1fcd2f5d3b2838429d4368ae4c57ce4436577f03fL61-R70) [/#diff-b803fcb7f17ed9235f1e5cb1fcd2f5d3b2838429d4368ae4c57ce4436577f03fR103-R118](https://github.com/adcondev///issues/diff-b803fcb7f17ed9235f1e5cb1fcd2f5d3b2838429d4368ae4c57ce4436577f03fR103-R118) [/#diff-89708b194914fd416481a81522862b10b3cd5f123be3eb550e6a1de67a01765aL6-R6](https://github.com/adcondev///issues/diff-89708b194914fd416481a81522862b10b3cd5f123be3eb550e6a1de67a01765aL6-R6) [/#diff-89708b194914fd416481a81522862b10b3cd5f123be3eb550e6a1de67a01765aL240-R242](https://github.com/adcondev///issues/diff-89708b194914fd416481a81522862b10b3cd5f123be3eb550e6a1de67a01765aL240-R242) [/#diff-89708b194914fd416481a81522862b10b3cd5f123be3eb550e6a1de67a01765aL269-R269](https://github.com/adcondev///issues/diff-89708b194914fd416481a81522862b10b3cd5f123be3eb550e6a1de67a01765aL269-R269) [/#diff-5c25f228ee796b995a56de464ea307b5cc1eeaaffb990f9329fe91a063fc3429L6-R6](https://github.com/adcondev///issues/diff-5c25f228ee796b995a56de464ea307b5cc1eeaaffb990f9329fe91a063fc3429L6-R6) [/#diff-cd2d359855d0301ce190f1ec3b4c572ea690c83747f6df61c9340720e3d2425eL26-R67](https://github.com/adcondev///issues/diff-cd2d359855d0301ce190f1ec3b4c572ea690c83747f6df61c9340720e3d2425eL26-R67) [/#diff-6179837f7df53a6f05c522b6b7bb566d484d5465d9894fb04910dd08bb40dcc9R8-R62](https://github.com/adcondev///issues/diff-6179837f7df53a6f05c522b6b7bb566d484d5465d9894fb04910dd08bb40dcc9R8-R62) [/#diff-6179837f7df53a6f05c522b6b7bb566d484d5465d9894fb04910dd08bb40dcc9R71-R83](https://github.com/adcondev///issues/diff-6179837f7df53a6f05c522b6b7bb566d484d5465d9894fb04910dd08bb40dcc9R71-R83) [/#diff-7a68c862ed13ecb99f59c4f61a92bbbc265afe66afa76b17bb9739a8cce7cab1L2-R2](https://github.com/adcondev///issues/diff-7a68c862ed13ecb99f59c4f61a92bbbc265afe66afa76b17bb9739a8cce7cab1L2-R2) [/#diff-7a68c862ed13ecb99f59c4f61a92bbbc265afe66afa76b17bb9739a8cce7cab1L12-R27](https://github.com/adcondev///issues/diff-7a68c862ed13ecb99f59c4f61a92bbbc265afe66afa76b17bb9739a8cce7cab1L12-R27)
* **escpos:** rename controllers to commands ([1aa49fc](https://github.com/adcondev/pos-printer/commit/1aa49fc332c9a369d524d8f15e63dd9cfa3f1ab1))

## [2.3.0](https://github.com/adcondev/pos-printer/compare/v2.2.1...v2.3.0) (2025-11-07)


### ✨ Features

* **escpos:** add document building and execution features ([#58](https://github.com/adcondev/pos-printer/issues/58)) ([3e89d73](https://github.com/adcondev/pos-printer/commit/3e89d734d6f86221088cb7ce9e0c64f2cf842864)), closes [#59](https://github.com/adcondev/pos-printer/issues/59)

### [2.2.1](https://github.com/adcondev/pos-printer/compare/v2.2.0...v2.2.1) (2025-11-06)


### 🐛 Bug Fixes

* **profile:** improve logging for unsupported code tables ([#57](https://github.com/adcondev/pos-printer/issues/57)) ([b886820](https://github.com/adcondev/pos-printer/commit/b886820fa250fa03a20189c54ed1da80fefb674d))

## [2.2.0](https://github.com/adcondev/pos-printer/compare/v2.1.0...v2.2.0) (2025-11-06)


### 🐛 Bug Fixes

* **profile:** enhance encoding support and error handling ([b28539b](https://github.com/adcondev/pos-printer/commit/b28539b1a312af0bdaf8431bbe8f8c985aaea6e4))


### ✨ Features

* **character:** add encoding support for character tables ([b6a2a6f](https://github.com/adcondev/pos-printer/commit/b6a2a6f98c248898ff09090eb683cba73556b40d))
* **escpos:** add graphics base64 image printing and autoencoding for printer code tables ([#56](https://github.com/adcondev/pos-printer/issues/56)) ([d636895](https://github.com/adcondev/pos-printer/commit/d63689555d53a42710014cf0db064108192a4e13)), closes [/#diff-562156e83a675c98d4982e84e8336971ca88db93fecfcceb69cd0fc1ca6fca18L28-R45](https://github.com/adcondev///issues/diff-562156e83a675c98d4982e84e8336971ca88db93fecfcceb69cd0fc1ca6fca18L28-R45) [/#diff-562156e83a675c98d4982e84e8336971ca88db93fecfcceb69cd0fc1ca6fca18L64-R56](https://github.com/adcondev///issues/diff-562156e83a675c98d4982e84e8336971ca88db93fecfcceb69cd0fc1ca6fca18L64-R56) [/#diff-562156e83a675c98d4982e84e8336971ca88db93fecfcceb69cd0fc1ca6fca18L84-R76](https://github.com/adcondev///issues/diff-562156e83a675c98d4982e84e8336971ca88db93fecfcceb69cd0fc1ca6fca18L84-R76) [/#diff-562156e83a675c98d4982e84e8336971ca88db93fecfcceb69cd0fc1ca6fca18L116-L125](https://github.com/adcondev///issues/diff-562156e83a675c98d4982e84e8336971ca88db93fecfcceb69cd0fc1ca6fca18L116-L125) [/#diff-a1c3ae8fd8f99a288574b99ee6b6d1b0fcfc7b6291679c91a9b2cd2b26631f74L47-R50](https://github.com/adcondev///issues/diff-a1c3ae8fd8f99a288574b99ee6b6d1b0fcfc7b6291679c91a9b2cd2b26631f74L47-R50) [/#diff-a1c3ae8fd8f99a288574b99ee6b6d1b0fcfc7b6291679c91a9b2cd2b26631f74L71-R74](https://github.com/adcondev///issues/diff-a1c3ae8fd8f99a288574b99ee6b6d1b0fcfc7b6291679c91a9b2cd2b26631f74L71-R74) [/#diff-a1c3ae8fd8f99a288574b99ee6b6d1b0fcfc7b6291679c91a9b2cd2b26631f74L80-R87](https://github.com/adcondev///issues/diff-a1c3ae8fd8f99a288574b99ee6b6d1b0fcfc7b6291679c91a9b2cd2b26631f74L80-R87) [/#diff-a1c3ae8fd8f99a288574b99ee6b6d1b0fcfc7b6291679c91a9b2cd2b26631f74R204-R213](https://github.com/adcondev///issues/diff-a1c3ae8fd8f99a288574b99ee6b6d1b0fcfc7b6291679c91a9b2cd2b26631f74R204-R213) [/#diff-af506b9fe4fcc35d7d61e3b6aba087ea5f4187c8a4a7025e40d5248dd0302abbL128-R128](https://github.com/adcondev///issues/diff-af506b9fe4fcc35d7d61e3b6aba087ea5f4187c8a4a7025e40d5248dd0302abbL128-R128) [/#diff-af506b9fe4fcc35d7d61e3b6aba087ea5f4187c8a4a7025e40d5248dd0302abbR138-R151](https://github.com/adcondev///issues/diff-af506b9fe4fcc35d7d61e3b6aba087ea5f4187c8a4a7025e40d5248dd0302abbR138-R151)
* **graphics:** add base64 image loading functionality and example ([aa5b68a](https://github.com/adcondev/pos-printer/commit/aa5b68aaf4475a9f8d859fc85466106dea5c03a7))
* **graphics:** update version to 2.1.0 and changelog ([32b7757](https://github.com/adcondev/pos-printer/commit/32b7757fb53bbda55a0d87f0c358442d84a5ddb7))

## [2.1.0](https://github.com/adcondev/pos-printer/compare/v2.0.1...v2.1.0) (2025-11-05)

### ✨ Features

* **graphics:** add base64 image loading functionality and
  example ([#55](https://github.com/adcondev/pos-printer/issues/55)) ([8abc948](https://github.com/adcondev/pos-printer/commit/8abc948da3d68f90792f796228ec8b861be42019))

### [2.0.1](https://github.com/adcondev/pos-printer/compare/v2.0.0...v2.0.1) (2025-11-05)


### 📝 Documentation

* **readme:** update LEARNING.md and README.md ([#53](https://github.com/adcondev/pos-printer/issues/53)) ([54d39a4](https://github.com/adcondev/pos-printer/commit/54d39a4ae9136ad3ff548d43a5285e3bb221b6a7))


### 🤖 Continuous Integration

* **escpos:** update CI configuration and add pre-commit hooks ([#54](https://github.com/adcondev/pos-printer/issues/54)) ([6ee16d6](https://github.com/adcondev/pos-printer/commit/6ee16d6be8e410d287b2c232b3146a2a461060cb))


### 📦 Dependencies

* **ci:** bump actions/setup-node from 5 to 6 ([#49](https://github.com/adcondev/pos-printer/issues/49)) ([925e1be](https://github.com/adcondev/pos-printer/commit/925e1bef25466a41aaec64a4899fa5f5aa12a19d))

## [2.0.0](https://github.com/AdConDev/pos-daemon/compare/v1.8.0...v2.0.0) (2025-11-04)


### ⚠ BREAKING CHANGES

* **escpos:** refactor Protocol to Commands and enhance functionality

### ✨ Features

* **arq:** introduce modular escpos architecture and printer profiles ([86d877b](https://github.com/AdConDev/pos-daemon/commit/86d877b64813613cbd437f24f9e1c8dd2c30371f))
* **escpos:** add QR code capability to ESC/POS protocol ([2d1f7f3](https://github.com/AdConDev/pos-daemon/commit/2d1f7f3a691bd7b915c96fc3b6b860eddb38a9cd))
* **escpos:** refactor Protocol to Commands and enhance functionality ([eec9582](https://github.com/AdConDev/pos-daemon/commit/eec9582bc030b068941a20d8daabb0668ca49b95))
* **graphics:** add advanced image processing engine ([0148e8c](https://github.com/AdConDev/pos-daemon/commit/0148e8c5a8e04987d384e3b5beaef4c136b1eeb9))
* **printposition:** introduce composer and refactor print position ([9b4aeb0](https://github.com/AdConDev/pos-daemon/commit/9b4aeb0104c09c2910f603570fa1005df5f9ce78))
* **qrcode:** implement QR Code generation commands ([57571fa](https://github.com/AdConDev/pos-daemon/commit/57571fa2c4bbe4a2135d200743841b4a373226de))


### 📦 Dependencies

* remove unused QR code dependency ([f479762](https://github.com/AdConDev/pos-daemon/commit/f4797628613309fb42f644ae289a1638f905790d))

## [1.8.0](https://github.com/AdConDev/pos-daemon/compare/v1.7.0...v1.8.0) (2025-10-28)


### ✨ Features

* **taskfiles:** add Task files to automate and improve some processes ([#48](https://github.com/adcondev/pos-printer/issues/48)) ([54de86c](https://github.com/AdConDev/pos-daemon/commit/54de86c6bde378206dca8c8d2f86627e45d3eff5))

## [1.7.0](https://github.com/AdConDev/pos-daemon/compare/v1.6.0...v1.7.0) (2025-10-23)


### ✨ Features

* **bitimage:** implement ESC/POS commands for bit images ([#45](https://github.com/adcondev/pos-printer/issues/45)) ([1b98b1f](https://github.com/AdConDev/pos-daemon/commit/1b98b1f9587ee1121961676f516ae946db0fc4e0))

## [1.6.0](https://github.com/AdConDev/pos-daemon/compare/v1.5.0...v1.6.0) (2025-10-22)


### ✨ Features

* **mechanismcontrol:** add commands for printer control ([#41](https://github.com/adcondev/pos-printer/issues/41)) ([53ef71c](https://github.com/AdConDev/pos-daemon/commit/53ef71c7f59ee935a7763a3d508ccb2ac95dca3c))

## [1.5.0](https://github.com/AdConDev/pos-daemon/compare/v1.4.0...v1.5.0) (2025-10-12)


### ✨ Features

* **escpos:** implemented SetTextSize(widthMultiplier, heightMultiplier int) method to adjust text dimensions. ([#36](https://github.com/adcondev/pos-printer/issues/36)) ([29c6e4e](https://github.com/AdConDev/pos-daemon/commit/29c6e4ece6ffa2767f8a19d8c51a327667ca5743))

## [1.4.0](https://github.com/AdConDev/pos-daemon/compare/v1.3.1...v1.4.0) (2025-10-07)


### 🤖 Continuous Integration

* bump actions/stale from 9 to 10 ([#27](https://github.com/adcondev/pos-printer/issues/27)) ([ee65c43](https://github.com/AdConDev/pos-daemon/commit/ee65c43cf84a99616a0ac5eb1934ec3483107ca3))
* bump lewagon/wait-on-check-action from 1.4.0 to 1.4.1 ([#30](https://github.com/adcondev/pos-printer/issues/30)) ([45bf780](https://github.com/AdConDev/pos-daemon/commit/45bf780aebbe70d48188b8f324aab391fcbf81f1))


### ✅ Tests

* **barcode:** add byte slice validation helpers ([2830158](https://github.com/AdConDev/pos-daemon/commit/283015888a1c1c2812e5b93cd370c2e41dced96b))


### ✨ Features

* **barcode:** add barcode commands and integration tests for barcode commands ([464a741](https://github.com/AdConDev/pos-daemon/commit/464a741c730267e5c6de1256b03ffc2cd8da2d0c))


### 🐛 Bug Fixes

* **barcode:** update utils/test/validation_helpers.go ([07fa11a](https://github.com/AdConDev/pos-daemon/commit/07fa11a2de8992c8bd929db1d9b21a18148b751d))

### [1.3.1](https://github.com/AdConDev/pos-daemon/compare/v1.3.0...v1.3.1) (2025-09-15)


### 🤖 Continuous Integration

* bump actions/checkout from 4 to 5 ([1e9499f](https://github.com/AdConDev/pos-daemon/commit/1e9499f554f7e4ff92a4b530e59607b54efb79bd))
* bump actions/github-script from 7 to 8 ([71bad7d](https://github.com/AdConDev/pos-daemon/commit/71bad7dd4ce67078fff9c54b17a178db0f71a023))
* bump actions/labeler from 5 to 6 ([e965a54](https://github.com/AdConDev/pos-daemon/commit/e965a5492638a9a97927630e38b83675e6916a18))
* bump actions/setup-go from 5 to 6 ([e62e0e0](https://github.com/AdConDev/pos-daemon/commit/e62e0e00eba21e2af109e7016bed19e2776d998a))
* bump actions/setup-node from 4 to 5 ([07802e7](https://github.com/AdConDev/pos-daemon/commit/07802e78c30ba8e3f66006d9fb2328ce76e2c771))
* bump lewagon/wait-on-check-action from 1.3.1 to 1.4.0 ([b39aad5](https://github.com/AdConDev/pos-daemon/commit/b39aad5ad15bb25245b36d7e55f993581d8644f0))


### 📦 Dependencies

* **deps:** bump golang.org/x/image from 0.30.0 to 0.31.0 in the golang-x group ([#22](https://github.com/adcondev/pos-printer/issues/22)) ([3c7131d](https://github.com/AdConDev/pos-daemon/commit/3c7131d6b5f338f248aaa9b19b1b6559f8698ddb))

## [1.3.0](https://github.com/AdConDev/pos-daemon/compare/v1.2.1...v1.3.0) (2025-09-15)


### ✨ Features

* **udchars:** add user-defined character test example ([f1e53cc](https://github.com/AdConDev/pos-daemon/commit/f1e53cc2799228917bfb476a96443794ddb78f81))

### [1.2.1](https://github.com/AdConDev/pos-daemon/compare/v1.2.0...v1.2.1) (2025-09-15)

### 📦 Dependencies

* **deps:** bump golang.org/x/text in the golang-x
  group ([cc79970](https://github.com/AdConDev/pos-daemon/commit/cc79970113218c838146bc75c0bac88c8a624c05))

### 🤖 Continuous Integration

* **dependabot:** enhance auto-merge workflow and add PR status
  dashboard ([1add7a1](https://github.com/AdConDev/pos-daemon/commit/1add7a13707c7835ad0b0ba5616daee9003d527a))

## [1.2.0](https://github.com/AdConDev/pos-daemon/compare/v1.1.0...v1.2.0) (2025-09-08)


### ✨ Features

* **printposition:** add print position management commands ([bc25641](https://github.com/AdConDev/pos-daemon/commit/bc256411ee830abdfd4757097942acb5c68dcabb))


### ✅ Tests

* **printposition:** update tests for print position functionality ([ab7705f](https://github.com/AdConDev/pos-daemon/commit/ab7705f279c4af5845b0145ac81463baa460fe5a))
* **print:** update error messages and assertions in tests ([3932bc0](https://github.com/AdConDev/pos-daemon/commit/3932bc0b14cc08b0b99fc3a2bb52052bed0b0b4a))
* **print:** update tests for print command integration ([74d8e58](https://github.com/AdConDev/pos-daemon/commit/74d8e5835311ce7d3b61843f2f5bc6df365249ce))
* **test:** add assertion helpers for testing utilities ([e382f4f](https://github.com/AdConDev/pos-daemon/commit/e382f4f36446f4d602de02499c14857fbc68f3e7))

## [1.1.0](https://github.com/AdConDev/pos-daemon/compare/v1.0.1...v1.1.0) (2025-09-02)


### ✨ Features

* **character:** add character handling capabilities ([#12](https://github.com/adcondev/pos-printer/issues/12)) ([c62dd8e](https://github.com/AdConDev/pos-daemon/commit/c62dd8eb651ee69cf9c7c92cadb2b3676bc2a344))

### [1.0.1](https://github.com/AdConDev/pos-daemon/compare/v1.0.0...v1.0.1) (2025-08-26)


### 📦 Dependencies

* **go.mod:** add missing golang.org/x/text dependency ([83d8598](https://github.com/AdConDev/pos-daemon/commit/83d859877d9b6a46d7ae9c6f65862ff6d7d09d9e))
* **go.mod:** fix go.mod file ([9d5b966](https://github.com/AdConDev/pos-daemon/commit/9d5b966795494d94b3fdd651fcbd03379de9da9e)), closes [#7](https://github.com/adcondev/pos-printer/issues/7)

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
