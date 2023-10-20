# Changelog

## [0.113.1](https://github.com/chanzuckerberg/happy/compare/shared-v0.113.0...shared-v0.113.1) (2023-10-18)


### Bug Fixes

* Allow image src aws role arn to be provided for cross-account image promotion ([#2611](https://github.com/chanzuckerberg/happy/issues/2611)) ([2c69389](https://github.com/chanzuckerberg/happy/commit/2c693897054b03d530c4d23a1969da7c8558e5d1))

## [0.113.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.112.1...shared-v0.113.0) (2023-10-17)


### Features

* start using the appconfig 'source' column ([#2597](https://github.com/chanzuckerberg/happy/issues/2597)) ([2c63967](https://github.com/chanzuckerberg/happy/commit/2c639678867e1514105e00d06eb2c5ea007861e4))

## [0.112.1](https://github.com/chanzuckerberg/happy/compare/shared-v0.112.0...shared-v0.112.1) (2023-10-17)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.112.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.111.0...shared-v0.112.0) (2023-10-16)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.111.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.110.1...shared-v0.111.0) (2023-10-13)


### Features

* replace Gorm with Ent ORM ([#2530](https://github.com/chanzuckerberg/happy/issues/2530)) ([fa87b1a](https://github.com/chanzuckerberg/happy/commit/fa87b1a0bbd2c6b41ac4e9f013c8c60ff5409913))

## [0.110.1](https://github.com/chanzuckerberg/happy/compare/shared-v0.110.0...shared-v0.110.1) (2023-10-05)


### Bug Fixes

* Suppress validation errors if module cannot be downloaded ([#2528](https://github.com/chanzuckerberg/happy/issues/2528)) ([fdc8e18](https://github.com/chanzuckerberg/happy/commit/fdc8e18fbaa2556fe8b5a39520173a22473279d7))

## [0.110.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.109.0...shared-v0.110.0) (2023-10-04)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.109.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.108.0...shared-v0.109.0) (2023-10-03)


### Features

* autogenerate default_tag field in aws provider when using happy infra generate command ([#2522](https://github.com/chanzuckerberg/happy/issues/2522)) ([c1143d6](https://github.com/chanzuckerberg/happy/commit/c1143d64937aaa4ffaabd15c63e353da6f1fe83e))
* Validate happy configuration on every happy operation ([#2511](https://github.com/chanzuckerberg/happy/issues/2511)) ([c1084f2](https://github.com/chanzuckerberg/happy/commit/c1084f2eca552f76e4010f5f1673e47f5981fa15))

## [0.108.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.107.0...shared-v0.108.0) (2023-09-25)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.107.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.106.2...shared-v0.107.0) (2023-09-22)


### Bug Fixes

* Fix integration test (mismatched parameter type) ([#2491](https://github.com/chanzuckerberg/happy/issues/2491)) ([9af5bb7](https://github.com/chanzuckerberg/happy/commit/9af5bb7efc055b5d32bc9c1ca562dcccc5db1650))

## [0.106.2](https://github.com/chanzuckerberg/happy/compare/shared-v0.106.1...shared-v0.106.2) (2023-09-21)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.106.1](https://github.com/chanzuckerberg/happy/compare/shared-v0.106.0...shared-v0.106.1) (2023-09-20)


### Bug Fixes

* Improve ECR scanning messaging ([#2480](https://github.com/chanzuckerberg/happy/issues/2480)) ([1d58703](https://github.com/chanzuckerberg/happy/commit/1d587039606ecf36212f65d24489cff811ca3588))

## [0.106.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.105.0...shared-v0.106.0) (2023-09-19)


### Features

* Warn when ECR scans fail before deployment ([#2477](https://github.com/chanzuckerberg/happy/issues/2477)) ([772d6c1](https://github.com/chanzuckerberg/happy/commit/772d6c1fafa7fbda4f12d42ab852e043bac8eed0))


### Bug Fixes

* Docker compose env file discovery doesn't work ([#2479](https://github.com/chanzuckerberg/happy/issues/2479)) ([d8003d6](https://github.com/chanzuckerberg/happy/commit/d8003d626ddb40059f04fc22013026bdd265ddbd))

## [0.105.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.104.1...shared-v0.105.0) (2023-09-15)


### Features

* Allow execution of shell commands non-interactively ([#2457](https://github.com/chanzuckerberg/happy/issues/2457)) ([cbbc2a5](https://github.com/chanzuckerberg/happy/commit/cbbc2a5bc4fe3803901465d5da6fc29386937d04))

## [0.104.1](https://github.com/chanzuckerberg/happy/compare/shared-v0.104.0...shared-v0.104.1) (2023-09-14)


### Bug Fixes

* better error reporting for happy cli and happy tf provider ([#2445](https://github.com/chanzuckerberg/happy/issues/2445)) ([894b4bd](https://github.com/chanzuckerberg/happy/commit/894b4bd804558e956e12e51b91304bb6ff12053d))

## [0.104.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.103.0...shared-v0.104.0) (2023-09-07)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.103.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.102.2...shared-v0.103.0) (2023-09-07)


### Features

* CCIE-1662 hvm GitHub PAT support ([#2387](https://github.com/chanzuckerberg/happy/issues/2387)) ([bcc3def](https://github.com/chanzuckerberg/happy/commit/bcc3def9783de6bb4f84a97a20e007c93559fbbe))

## [0.102.2](https://github.com/chanzuckerberg/happy/compare/shared-v0.102.1...shared-v0.102.2) (2023-09-01)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.102.1](https://github.com/chanzuckerberg/happy/compare/shared-v0.102.0...shared-v0.102.1) (2023-09-01)


### Bug Fixes

* On bootstrap, prompt the user if docker-compose.yml already exists ([#2392](https://github.com/chanzuckerberg/happy/issues/2392)) ([5cefe53](https://github.com/chanzuckerberg/happy/commit/5cefe53bd543eedfe886df5d33cf280682ef4717))

## [0.102.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.101.0...shared-v0.102.0) (2023-09-01)


### Bug Fixes

* using the wrong AWS profile when promoting images ([#2393](https://github.com/chanzuckerberg/happy/issues/2393)) ([43330fc](https://github.com/chanzuckerberg/happy/commit/43330fc37dadf7458f5ba4806b2d19deff12859e))

## [0.101.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.100.0...shared-v0.101.0) (2023-08-29)


### Features

* allow the generated terraform files to use app variable  ([#2381](https://github.com/chanzuckerberg/happy/issues/2381)) ([adbda5b](https://github.com/chanzuckerberg/happy/commit/adbda5b6d3f7cc2aaf33d28cfa9b56d8eeab1c43))
* Support for ECR tag immutability ([#2376](https://github.com/chanzuckerberg/happy/issues/2376)) ([c1d5f5b](https://github.com/chanzuckerberg/happy/commit/c1d5f5b6e6a093c19ba2a092111842cc0e4f195f))

## [0.100.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.99.0...shared-v0.100.0) (2023-08-28)


### Bug Fixes

* Made workspace run message consistent ([#2374](https://github.com/chanzuckerberg/happy/issues/2374)) ([c478184](https://github.com/chanzuckerberg/happy/commit/c478184ebe03372aeff230fd4b94a3871723afdb))
* Notify user to restart docker engine, allow pre-release for docker compose ([#2377](https://github.com/chanzuckerberg/happy/issues/2377)) ([48745e6](https://github.com/chanzuckerberg/happy/commit/48745e66116b0c5a6e82be71b0ec2f3653f36606))

## [0.99.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.98.0...shared-v0.99.0) (2023-08-25)


### Features

* Added app name field to be included in auto-generated field ([#2351](https://github.com/chanzuckerberg/happy/issues/2351)) ([42aac44](https://github.com/chanzuckerberg/happy/commit/42aac449515235843b1b19ef588acba1269101cb))

## [0.98.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.97.1...shared-v0.98.0) (2023-08-22)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.97.1](https://github.com/chanzuckerberg/happy/compare/shared-v0.97.0...shared-v0.97.1) (2023-08-16)


### Bug Fixes

* only include stacks for the requested app in api response ([#2269](https://github.com/chanzuckerberg/happy/issues/2269)) ([4491496](https://github.com/chanzuckerberg/happy/commit/4491496f8d81f9e4c002aef2901fbd59bc173494))

## [0.97.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.96.1...shared-v0.97.0) (2023-08-08)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.96.1](https://github.com/chanzuckerberg/happy/compare/shared-v0.96.0...shared-v0.96.1) (2023-08-04)


### Bug Fixes

* Bootstrap doesn't generate code for all configured environments ([#2204](https://github.com/chanzuckerberg/happy/issues/2204)) ([a54f67d](https://github.com/chanzuckerberg/happy/commit/a54f67d8448f800efa8f77f4145323e62854acf2))

## [0.96.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.95.0...shared-v0.96.0) (2023-08-04)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.95.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.94.0...shared-v0.95.0) (2023-08-04)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.94.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.93.0...shared-v0.94.0) (2023-08-02)


### Features

* CCIE-1619: happy version manager v1 ([#2066](https://github.com/chanzuckerberg/happy/issues/2066)) ([816447b](https://github.com/chanzuckerberg/happy/commit/816447b5255f22cafd3795ef244e628b1af4ea4a))
* consolidate stack service in shared pkg ([#2096](https://github.com/chanzuckerberg/happy/issues/2096)) ([24d885c](https://github.com/chanzuckerberg/happy/commit/24d885cd8a8845d1e1d1934c1c3e345cfb0e951e))

## [0.93.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.92.0...shared-v0.93.0) (2023-07-28)


### Features

* happy restart stack ([#2127](https://github.com/chanzuckerberg/happy/issues/2127)) ([975ad28](https://github.com/chanzuckerberg/happy/commit/975ad28d547c2a5c8b784736af1883adfc6f0f43))

## [0.92.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.91.1...shared-v0.92.0) (2023-07-10)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.91.1](https://github.com/chanzuckerberg/happy/compare/shared-v0.91.0...shared-v0.91.1) (2023-07-03)


### Bug Fixes

* broken filepath in shared stack package ([#1998](https://github.com/chanzuckerberg/happy/issues/1998)) ([dd7e714](https://github.com/chanzuckerberg/happy/commit/dd7e714b06247d97e4a9785f2dd238474f8cca58))

## [0.91.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.90.1...shared-v0.91.0) (2023-06-30)


### Features

* fix duplicates returning from API ([#1990](https://github.com/chanzuckerberg/happy/issues/1990)) ([58a0aa7](https://github.com/chanzuckerberg/happy/commit/58a0aa745a9646d34fc7adc418001d8f63d65047))

## [0.90.1](https://github.com/chanzuckerberg/happy/compare/shared-v0.90.0...shared-v0.90.1) (2023-06-28)


### Bug Fixes

* reuse datastructure ([#1982](https://github.com/chanzuckerberg/happy/issues/1982)) ([4bd98db](https://github.com/chanzuckerberg/happy/commit/4bd98db581e8a72a8ad9c6032126215eac220cc3))

## [0.90.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.89.1...shared-v0.90.0) (2023-06-22)


### Features

* reuse happy client ([#1960](https://github.com/chanzuckerberg/happy/issues/1960)) ([fc3991d](https://github.com/chanzuckerberg/happy/commit/fc3991d0670579e34013e854e6a5a4f3fc4e189e))

## [0.89.1](https://github.com/chanzuckerberg/happy/compare/shared-v0.89.0...shared-v0.89.1) (2023-06-21)


### Bug Fixes

* docker-compose.yml doesn't allow for a name attribute on services ([#1959](https://github.com/chanzuckerberg/happy/issues/1959)) ([4e18e5e](https://github.com/chanzuckerberg/happy/commit/4e18e5e082c9277348f3cff31ca85f8db7fdd66a))

## [0.89.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.88.0...shared-v0.89.0) (2023-06-21)


### Features

* consolidate stack operations on cli and hapi ([#1867](https://github.com/chanzuckerberg/happy/issues/1867)) ([a4a8b5d](https://github.com/chanzuckerberg/happy/commit/a4a8b5db6ce01811592278107da58cb0aba5fc5b))


### Bug Fixes

* Fix breaking change in the tfe api mock ([#1957](https://github.com/chanzuckerberg/happy/issues/1957)) ([a9a372d](https://github.com/chanzuckerberg/happy/commit/a9a372dbe942efa0c3cad9ff619a2555f9381bc6))
* happy infra ingest fails for the tasks example ([#1956](https://github.com/chanzuckerberg/happy/issues/1956)) ([99d74e0](https://github.com/chanzuckerberg/happy/commit/99d74e097b7b56c91410b7b61dbe420174843a78))

## [0.88.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.87.0...shared-v0.88.0) (2023-06-20)


### Features

* Implement "happy bootstrap" to init happy configuration on the existing GitHub repo ([#1866](https://github.com/chanzuckerberg/happy/issues/1866)) ([6cd3084](https://github.com/chanzuckerberg/happy/commit/6cd3084fd720f972f4434e82db2112b225230ee3))


### Bug Fixes

* Service port numbers are not added to docker-compose past happy bootstrap ([#1943](https://github.com/chanzuckerberg/happy/issues/1943)) ([e0e9603](https://github.com/chanzuckerberg/happy/commit/e0e960310699c4dde08d5ee234ebcebedaaf4798))

## [0.87.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.86.0...shared-v0.87.0) (2023-06-07)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.86.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.85.1...shared-v0.86.0) (2023-06-06)


### Features

* Consume and produce docker-compose.yml via happy ingest and happy generate ([#1852](https://github.com/chanzuckerberg/happy/issues/1852)) ([addb506](https://github.com/chanzuckerberg/happy/commit/addb506505db527e6c08c71a33717cb38fd1b570))

## [0.85.1](https://github.com/chanzuckerberg/happy/compare/shared-v0.85.0...shared-v0.85.1) (2023-06-05)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.85.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.84.0...shared-v0.85.0) (2023-06-05)


### Features

* Add sidecar support to happy logs ([#1844](https://github.com/chanzuckerberg/happy/issues/1844)) ([12949d7](https://github.com/chanzuckerberg/happy/commit/12949d7b027721b69b0acf4e2b0f71dc5c4b1fb9))


### Bug Fixes

* happy infra generate doesn't work with additional_envs_from_secret ([#1845](https://github.com/chanzuckerberg/happy/issues/1845)) ([d48fef7](https://github.com/chanzuckerberg/happy/commit/d48fef77bcc129f93a1f3e1984664d3fb59acb7d))

## [0.84.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.83.0...shared-v0.84.0) (2023-06-02)


### Features

* Add sidecar support to happy shell ([#1842](https://github.com/chanzuckerberg/happy/issues/1842)) ([9c52320](https://github.com/chanzuckerberg/happy/commit/9c5232066acebd6562541de03b91028bca1fc8bb))


### Bug Fixes

* Display a link to failed TFE plan in case an error occurs in TFE, even when -v flag is not passed ([#1838](https://github.com/chanzuckerberg/happy/issues/1838)) ([c86b96b](https://github.com/chanzuckerberg/happy/commit/c86b96b3348e6e2f1c9c4421d5b1838ce08b063c))
* Gracefully fail if module invocation doesn't have the source specified ([#1841](https://github.com/chanzuckerberg/happy/issues/1841)) ([196135f](https://github.com/chanzuckerberg/happy/commit/196135f922778205bd3a1a413d7e7f8d51eb3e28))

## [0.83.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.82.1...shared-v0.83.0) (2023-06-01)


### Features

* Implement "happy infra refresh" to refresh terraform scripts ([#1832](https://github.com/chanzuckerberg/happy/issues/1832)) ([52fc23d](https://github.com/chanzuckerberg/happy/commit/52fc23dc3517c7fbe209aa82ac95ee9cf41c7e9f))
* multistack destroy; refactor destroy ([#1833](https://github.com/chanzuckerberg/happy/issues/1833)) ([7c37665](https://github.com/chanzuckerberg/happy/commit/7c3766504521025b4b8bfc8d07264b723ac5a4f6))

## [0.82.1](https://github.com/chanzuckerberg/happy/compare/shared-v0.82.0...shared-v0.82.1) (2023-05-31)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.82.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.81.0...shared-v0.82.0) (2023-05-31)


### Features

* Implement "happy infra validate" to validate terraform scripts ([#1824](https://github.com/chanzuckerberg/happy/issues/1824)) ([a57c977](https://github.com/chanzuckerberg/happy/commit/a57c9775cc436e92e3475edb6b880b49e07807b0))

## [0.81.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.80.0...shared-v0.81.0) (2023-05-30)


### Features

* Example of task usage in happy EKS ([#1776](https://github.com/chanzuckerberg/happy/issues/1776)) ([2af7c7f](https://github.com/chanzuckerberg/happy/commit/2af7c7faa87938ea859db26fe143eca429f61d86))


### Bug Fixes

* [bug] Validate credentials before stack operations and prompt user to log in and create a new token on token absence or prior invalidation ([#1806](https://github.com/chanzuckerberg/happy/issues/1806)) ([e23146a](https://github.com/chanzuckerberg/happy/commit/e23146ac94363551ff5990c533637f61344d5f94))

## [0.80.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.79.1...shared-v0.80.0) (2023-05-24)


### Features

* Collect stack configuration from existing terraform code and store in happy config ([#1761](https://github.com/chanzuckerberg/happy/issues/1761)) ([56dd781](https://github.com/chanzuckerberg/happy/commit/56dd7819d44b6464e2dd0d43ab27d77411fcf680))


### Bug Fixes

* Happy logs look weird when workspace has new format ([#1765](https://github.com/chanzuckerberg/happy/issues/1765)) ([39e8a9f](https://github.com/chanzuckerberg/happy/commit/39e8a9f5115102d21041c643f90a8157a5f5c01b))

## [0.79.1](https://github.com/chanzuckerberg/happy/compare/shared-v0.79.0...shared-v0.79.1) (2023-05-10)


### Bug Fixes

* When more than one service is specified, and the settings structure is inconsistent, happy infra generate errors out ([#1751](https://github.com/chanzuckerberg/happy/issues/1751)) ([ea166c2](https://github.com/chanzuckerberg/happy/commit/ea166c20cd6a52e0ef82a53554261f0055d680ed))

## [0.79.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.78.0...shared-v0.79.0) (2023-05-09)


### Features

* Sidecar support for services ([#1727](https://github.com/chanzuckerberg/happy/issues/1727)) ([8c5c884](https://github.com/chanzuckerberg/happy/commit/8c5c884804a4e88d1e3163f266127e6ddb336c05))


### Bug Fixes

* Refresh EKS credentials after a lengthy docker build ([#1728](https://github.com/chanzuckerberg/happy/issues/1728)) ([b9d422b](https://github.com/chanzuckerberg/happy/commit/b9d422beea1930d5806dcf6186d7fce3092c0fdd))

## [0.78.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.77.0...shared-v0.78.0) (2023-05-08)


### Features

* Implement Happy cli basic mode (terraform code is generated) ([#1684](https://github.com/chanzuckerberg/happy/issues/1684)) ([ca41c53](https://github.com/chanzuckerberg/happy/commit/ca41c538bfb99491028ab07b55308c88fc3d4a03))

## [0.77.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.76.0...shared-v0.77.0) (2023-05-02)


### Features

* add command to see the configured CI roles for env ([#1686](https://github.com/chanzuckerberg/happy/issues/1686)) ([a249cc0](https://github.com/chanzuckerberg/happy/commit/a249cc0a4fc61af413312b300f1fc4695529ee2e))

## [0.76.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.75.0...shared-v0.76.0) (2023-04-27)


### Features

* CCIE-960 do not require docker to be running for functions that don t use it ([#1659](https://github.com/chanzuckerberg/happy/issues/1659)) ([7c53ee6](https://github.com/chanzuckerberg/happy/commit/7c53ee6492300f89724182701a305d65c62b1aa1))

## [0.75.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.74.0...shared-v0.75.0) (2023-04-27)


### Features

* allow for stacks to migrate container artifacts ([#1619](https://github.com/chanzuckerberg/happy/issues/1619)) ([09cea95](https://github.com/chanzuckerberg/happy/commit/09cea95566c41b34f12a1d2f858ff3bef8d598a6))

## [0.74.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.73.3...shared-v0.74.0) (2023-04-27)


### Features

* Happy CLI integration test ([#1662](https://github.com/chanzuckerberg/happy/issues/1662)) ([a3c4d2c](https://github.com/chanzuckerberg/happy/commit/a3c4d2ce28a095f47d9c66c9ddfd24b231b864b6))

## [0.73.3](https://github.com/chanzuckerberg/happy/compare/shared-v0.73.2...shared-v0.73.3) (2023-04-27)


### Bug Fixes

* workspaces with no state will not have info ([#1663](https://github.com/chanzuckerberg/happy/issues/1663)) ([892b463](https://github.com/chanzuckerberg/happy/commit/892b4633e9bf97fb71c1e153369bd705bbc70f26))

## [0.73.2](https://github.com/chanzuckerberg/happy/compare/shared-v0.73.1...shared-v0.73.2) (2023-04-24)


### Bug Fixes

* state drilling with dryrun; use context ([#1607](https://github.com/chanzuckerberg/happy/issues/1607)) ([a75376a](https://github.com/chanzuckerberg/happy/commit/a75376a849940d9cdf45accbc1ec0357dbd0c3f8))

## [0.73.1](https://github.com/chanzuckerberg/happy/compare/shared-v0.73.0...shared-v0.73.1) (2023-04-20)


### Bug Fixes

* add more environment types ([#1588](https://github.com/chanzuckerberg/happy/issues/1588)) ([48d85fe](https://github.com/chanzuckerberg/happy/commit/48d85fe30aa2a05868fc7db075a27bc4b7e4eaa2))
* duplicate envs ([#1606](https://github.com/chanzuckerberg/happy/issues/1606)) ([077e0c0](https://github.com/chanzuckerberg/happy/commit/077e0c0943fc3f61399b527b5e7ae534b4403060))

## [0.73.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.72.0...shared-v0.73.0) (2023-04-20)


### Features

* filter the stacks by app; display git info ([#1594](https://github.com/chanzuckerberg/happy/issues/1594)) ([665f35c](https://github.com/chanzuckerberg/happy/commit/665f35c39d7eff37ff8d0bca34f37db08f0eb753))

## [0.72.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.71.0...shared-v0.72.0) (2023-04-20)


### Features

* Happy debug feature support for EKS ([#1592](https://github.com/chanzuckerberg/happy/issues/1592)) ([08eb06a](https://github.com/chanzuckerberg/happy/commit/08eb06acda5990fe5c4fd4aedc57eaf7179233d0))

## [0.71.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.70.1...shared-v0.71.0) (2023-04-18)


### Features

* Support the happy events feature to visualize events from key applicaiton levels ([#1579](https://github.com/chanzuckerberg/happy/issues/1579)) ([367d958](https://github.com/chanzuckerberg/happy/commit/367d958486536d2812940865d314bd1cd2490d23))

## [0.70.1](https://github.com/chanzuckerberg/happy/compare/shared-v0.70.0...shared-v0.70.1) (2023-04-12)


### Bug Fixes

* Breaking change in a kubernetes api ([#1565](https://github.com/chanzuckerberg/happy/issues/1565)) ([5967f4a](https://github.com/chanzuckerberg/happy/commit/5967f4a6680ed9d4495cc241b843f88a40c7f8cc))

## [0.70.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.69.2...shared-v0.70.0) (2023-04-10)


### Features

* Happy Service integrity check ([#1495](https://github.com/chanzuckerberg/happy/issues/1495)) ([29f7804](https://github.com/chanzuckerberg/happy/commit/29f780437bf28f4ae9c309ad47f1dd752b156559))

## [0.69.2](https://github.com/chanzuckerberg/happy/compare/shared-v0.69.1...shared-v0.69.2) (2023-04-07)


### Bug Fixes

* use aws credentials from request ([#1493](https://github.com/chanzuckerberg/happy/issues/1493)) ([8608647](https://github.com/chanzuckerberg/happy/commit/8608647a6e7e8ee2024f211a12fcff7fdf4fae4e))

## [0.69.1](https://github.com/chanzuckerberg/happy/compare/shared-v0.69.0...shared-v0.69.1) (2023-04-07)


### Bug Fixes

* omit empty errors in api response ([#1491](https://github.com/chanzuckerberg/happy/issues/1491)) ([fbefd5f](https://github.com/chanzuckerberg/happy/commit/fbefd5f6d523799a1b9e9d7a0c8cffa9b7398abc))

## [0.69.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.68.0...shared-v0.69.0) (2023-04-07)


### Features

* Exclude local terraform artifacts when zipping up the stack on create/update ([#1488](https://github.com/chanzuckerberg/happy/issues/1488)) ([2c28fc1](https://github.com/chanzuckerberg/happy/commit/2c28fc1ff8c13a0d9f587c713eeeb5f2027c8073))
* Expose stack TFE status, TFE Url, and Endpoints through HAPI ([#1469](https://github.com/chanzuckerberg/happy/issues/1469)) ([820396a](https://github.com/chanzuckerberg/happy/commit/820396ac31c9416ba49afe0ac73dfd816ad2e9c4))
* Happy config: Make aws region configurable ([#1487](https://github.com/chanzuckerberg/happy/issues/1487)) ([b70ad5e](https://github.com/chanzuckerberg/happy/commit/b70ad5e43e020965b7683eec82e62aa1ca02bff5))
* Remove happy config from backend ([#1472](https://github.com/chanzuckerberg/happy/issues/1472)) ([7421240](https://github.com/chanzuckerberg/happy/commit/7421240f96be6b891b43be893429b7d62e574c80))

## [0.68.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.67.0...shared-v0.68.0) (2023-04-05)


### Features

* Move backend, workspace_repo package to shared ([#1467](https://github.com/chanzuckerberg/happy/issues/1467)) ([d0b64ed](https://github.com/chanzuckerberg/happy/commit/d0b64edd690e91690438de6c35671a90d248f9ba))

## [0.67.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.66.2...shared-v0.67.0) (2023-03-29)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.66.2](https://github.com/chanzuckerberg/happy/compare/shared-v0.66.1...shared-v0.66.2) (2023-03-28)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.66.1](https://github.com/chanzuckerberg/happy/compare/shared-v0.66.0...shared-v0.66.1) (2023-03-28)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.66.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.65.0...shared-v0.66.0) (2023-03-27)


### Features

* add happy_stacklist data item to provider ([#1388](https://github.com/chanzuckerberg/happy/issues/1388)) ([9225f4d](https://github.com/chanzuckerberg/happy/commit/9225f4d6ff27d379882bf20d59b012feb3fb2023))

## [0.65.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.64.0...shared-v0.65.0) (2023-03-15)


### Features

* CCIE-900 Automatically check whether Happy is up to date ([#1355](https://github.com/chanzuckerberg/happy/issues/1355)) ([7cec2dd](https://github.com/chanzuckerberg/happy/commit/7cec2dd277b1eaf995780d9cd4ffdba3fcbb46fe))

## [0.64.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.63.1...shared-v0.64.0) (2023-03-13)


### Features

* 'happy logs' integration with Cloudwatch Insights ([#1315](https://github.com/chanzuckerberg/happy/issues/1315)) ([9ff4861](https://github.com/chanzuckerberg/happy/commit/9ff48617f79273457018d21de2a1ad78b9109a07))

## [0.63.1](https://github.com/chanzuckerberg/happy/compare/shared-v0.63.0...shared-v0.63.1) (2023-03-08)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.63.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.62.1...shared-v0.63.0) (2023-03-07)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.62.1](https://github.com/chanzuckerberg/happy/compare/shared-v0.62.0...shared-v0.62.1) (2023-03-07)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.62.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.61.0...shared-v0.62.0) (2023-03-07)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.61.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.60.0...shared-v0.61.0) (2023-03-07)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.60.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.59.0...shared-v0.60.0) (2023-03-07)


### ⚠ BREAKING CHANGES

* refactor update/clean; autocreate ECR ([#1232](https://github.com/chanzuckerberg/happy/issues/1232))

### Features

* refactor update/clean; autocreate ECR ([#1232](https://github.com/chanzuckerberg/happy/issues/1232)) ([b498074](https://github.com/chanzuckerberg/happy/commit/b4980740c3ddc716abe530fb2112dfe41bc6ab60))

## [0.59.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.58.0...shared-v0.59.0) (2023-02-28)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.58.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.57.0...shared-v0.58.0) (2023-02-24)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.57.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.56.1...shared-v0.57.0) (2023-02-23)


### Features

* List of all AWS infra associated with a k8s happy stack ([#1217](https://github.com/chanzuckerberg/happy/issues/1217)) ([83586fb](https://github.com/chanzuckerberg/happy/commit/83586fb2950a30677884245c3dc6cc8efa4968a7))
* return created_at/updated_at in api responses ([#1216](https://github.com/chanzuckerberg/happy/issues/1216)) ([0edcdab](https://github.com/chanzuckerberg/happy/commit/0edcdab9f9745baa6d630f4ac6c725b4ef80b67c))

## [0.56.1](https://github.com/chanzuckerberg/happy/compare/shared-v0.56.0...shared-v0.56.1) (2023-02-21)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.56.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.55.1...shared-v0.56.0) (2023-02-17)


### Bug Fixes

* ran 'go mod tidy' ([#1172](https://github.com/chanzuckerberg/happy/issues/1172)) ([fd0fcc7](https://github.com/chanzuckerberg/happy/commit/fd0fcc7782e18229979c7eaa622ecceeadf1b528))

## [0.55.1](https://github.com/chanzuckerberg/happy/compare/shared-v0.55.0...shared-v0.55.1) (2023-02-13)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.55.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.54.1...shared-v0.55.0) (2023-02-13)


### Features

* inject release version into api image ([#1139](https://github.com/chanzuckerberg/happy/issues/1139)) ([cf8b017](https://github.com/chanzuckerberg/happy/commit/cf8b0175d6367e05146b0ba6359655d9fdb14e5a))

## [0.54.1](https://github.com/chanzuckerberg/happy/compare/shared-v0.54.0...shared-v0.54.1) (2023-02-13)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.54.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.53.6...shared-v0.54.0) (2023-02-13)


### ⚠ BREAKING CHANGES

* inject happy config to stacks ([#1108](https://github.com/chanzuckerberg/happy/issues/1108))

### Features

* inject happy config to stacks ([#1108](https://github.com/chanzuckerberg/happy/issues/1108)) ([9cb49c7](https://github.com/chanzuckerberg/happy/commit/9cb49c7f7bd6819541510e4f31ab5fd112579457))

## [0.53.6](https://github.com/chanzuckerberg/happy/compare/shared-v0.53.5...shared-v0.53.6) (2023-02-10)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.53.5](https://github.com/chanzuckerberg/happy/compare/shared-v0.53.4...shared-v0.53.5) (2023-02-10)


### Bug Fixes

* better error handling ([#1128](https://github.com/chanzuckerberg/happy/issues/1128)) ([3ff095a](https://github.com/chanzuckerberg/happy/commit/3ff095a7ec9b5c2ddb96fdd2c3b9e62fde2dbc43))

## [0.53.4](https://github.com/chanzuckerberg/happy/compare/shared-v0.53.3...shared-v0.53.4) (2023-02-10)


### Bug Fixes

* make sure region is used to configure the provider ([#1126](https://github.com/chanzuckerberg/happy/issues/1126)) ([423a6aa](https://github.com/chanzuckerberg/happy/commit/423a6aaafb9f7dec012051fe4e22bd9afc1ba069))

## [0.53.3](https://github.com/chanzuckerberg/happy/compare/shared-v0.53.2...shared-v0.53.3) (2023-02-09)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.53.2](https://github.com/chanzuckerberg/happy/compare/shared-v0.53.1...shared-v0.53.2) (2023-02-09)


### Bug Fixes

* find git root using rev-parse ([#1113](https://github.com/chanzuckerberg/happy/issues/1113)) ([9f16ba6](https://github.com/chanzuckerberg/happy/commit/9f16ba6907b10159ec4db2c19ff28c80628e6139))

## [0.53.1](https://github.com/chanzuckerberg/happy/compare/shared-v0.53.0...shared-v0.53.1) (2023-02-09)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.53.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.52.0...shared-v0.53.0) (2023-02-09)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.52.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.51.0...shared-v0.52.0) (2023-02-08)


### Features

* use query string for GET requests to happy api ([#1101](https://github.com/chanzuckerberg/happy/issues/1101)) ([7a18eb8](https://github.com/chanzuckerberg/happy/commit/7a18eb8dd5bc2eaebdb246dbebd44f4c389b17e2))

## [0.51.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.50.2...shared-v0.51.0) (2023-02-08)


### Features

* CCIE-926 List of all happy stacks for an app env ([#1068](https://github.com/chanzuckerberg/happy/issues/1068)) ([fc8d8b1](https://github.com/chanzuckerberg/happy/commit/fc8d8b1353f822e7768d39734adc533e90c49876))

## [0.50.2](https://github.com/chanzuckerberg/happy/compare/shared-v0.50.1...shared-v0.50.2) (2023-01-30)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.50.1](https://github.com/chanzuckerberg/happy/compare/shared-v0.50.0...shared-v0.50.1) (2023-01-30)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.50.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.49.0...shared-v0.50.0) (2023-01-27)


### Features

* Abstract out kubernetes authentication ([#1024](https://github.com/chanzuckerberg/happy/issues/1024)) ([e5712ef](https://github.com/chanzuckerberg/happy/commit/e5712ef334bcb7d60c07c36ed1f6afe22566a1d9))
* Move backend interfaces to a shared module ([#1026](https://github.com/chanzuckerberg/happy/issues/1026)) ([b0921a8](https://github.com/chanzuckerberg/happy/commit/b0921a834e52895f0cd92eebf7b65fc56f7425fc))

## [0.49.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.48.0...shared-v0.49.0) (2023-01-24)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.48.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.47.1...shared-v0.48.0) (2023-01-19)


### Features

* send aws creds in requests to api ([#962](https://github.com/chanzuckerberg/happy/issues/962)) ([01c6b79](https://github.com/chanzuckerberg/happy/commit/01c6b79d1b4ea27ee54d3dc96a9a247075189aa0))

## [0.47.1](https://github.com/chanzuckerberg/happy/compare/shared-v0.47.0...shared-v0.47.1) (2023-01-17)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.47.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.46.1...shared-v0.47.0) (2023-01-17)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.46.1](https://github.com/chanzuckerberg/happy/compare/shared-v0.46.0...shared-v0.46.1) (2023-01-09)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.46.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.45.0...shared-v0.46.0) (2023-01-04)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.45.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.44.0...shared-v0.45.0) (2022-12-21)


### Features

* add api meta-command ([#903](https://github.com/chanzuckerberg/happy/issues/903)) ([b81871b](https://github.com/chanzuckerberg/happy/commit/b81871bf694063ce172267e3dcbfe08d737f4120))

## [0.44.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.43.2...shared-v0.44.0) (2022-12-20)


### Features

* send auth header in api requests ([#785](https://github.com/chanzuckerberg/happy/issues/785)) ([d83c9b3](https://github.com/chanzuckerberg/happy/commit/d83c9b3c57950b1747d8233166e276d883cda4a7))

## [0.43.2](https://github.com/chanzuckerberg/happy/compare/shared-v0.43.1...shared-v0.43.2) (2022-12-20)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.43.1](https://github.com/chanzuckerberg/happy/compare/shared-v0.43.0...shared-v0.43.1) (2022-12-19)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.43.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.42.1...shared-v0.43.0) (2022-12-16)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.42.1](https://github.com/chanzuckerberg/happy/compare/shared-v0.42.0...shared-v0.42.1) (2022-12-13)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.42.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.41.5...shared-v0.42.0) (2022-12-12)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.41.5](https://github.com/chanzuckerberg/happy/compare/shared-v0.41.4...shared-v0.41.5) (2022-12-12)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.41.4](https://github.com/chanzuckerberg/happy/compare/shared-v0.41.3...shared-v0.41.4) (2022-12-08)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.41.3](https://github.com/chanzuckerberg/happy/compare/shared-v0.41.2...shared-v0.41.3) (2022-12-07)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.41.2](https://github.com/chanzuckerberg/happy/compare/shared-v0.41.1...shared-v0.41.2) (2022-12-02)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.41.1](https://github.com/chanzuckerberg/happy/compare/shared-v0.41.0...shared-v0.41.1) (2022-11-17)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.41.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.40.1...shared-v0.41.0) (2022-11-16)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.40.1](https://github.com/chanzuckerberg/happy/compare/shared-v0.40.0...shared-v0.40.1) (2022-11-03)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.40.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.39.0...shared-v0.40.0) (2022-11-03)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.39.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.38.0...shared-v0.39.0) (2022-10-25)


### Features

* create terraform provider for happy-api ([#699](https://github.com/chanzuckerberg/happy/issues/699)) ([3325039](https://github.com/chanzuckerberg/happy/commit/3325039ae0fa433ee4d59307762869ed543b8554))

## [0.38.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.37.0...shared-v0.38.0) (2022-10-24)


### Features

* roll out config feature to cli ([#660](https://github.com/chanzuckerberg/happy/issues/660)) ([a72c965](https://github.com/chanzuckerberg/happy/commit/a72c965f6bd2c9113c8152c9155330971e808b46))

## [0.37.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.36.0...shared-v0.37.0) (2022-10-13)


### Features

* move models to shared package ([#657](https://github.com/chanzuckerberg/happy/issues/657)) ([2f42c9d](https://github.com/chanzuckerberg/happy/commit/2f42c9df6629c2adba23498b320c56cfe58335c0))

## [0.36.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.35.2...shared-v0.36.0) (2022-10-12)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.35.2](https://github.com/chanzuckerberg/happy/compare/shared-v0.35.1...shared-v0.35.2) (2022-10-12)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.35.1](https://github.com/chanzuckerberg/happy/compare/shared-v0.35.0...shared-v0.35.1) (2022-10-11)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.35.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.34.4...shared-v0.35.0) (2022-10-11)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.34.4](https://github.com/chanzuckerberg/happy/compare/shared-v0.34.3...shared-v0.34.4) (2022-10-11)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.34.3](https://github.com/chanzuckerberg/happy/compare/shared-v0.34.2...shared-v0.34.3) (2022-10-10)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.34.2](https://github.com/chanzuckerberg/happy/compare/shared-v0.34.1...shared-v0.34.2) (2022-10-10)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.34.1](https://github.com/chanzuckerberg/happy/compare/shared-v0.34.0...shared-v0.34.1) (2022-10-10)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.34.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.33.0...shared-v0.34.0) (2022-10-10)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## [0.33.0](https://github.com/chanzuckerberg/happy/compare/shared-v0.32.0...shared-v0.33.0) (2022-10-07)


### Miscellaneous Chores

* **shared:** Synchronize happy platform versions

## 0.32.0 (2022-10-07)


### Features

* add shared package ([#620](https://github.com/chanzuckerberg/happy/issues/620)) ([159bd8e](https://github.com/chanzuckerberg/happy/commit/159bd8e372cdf4c2897ca71395c1d65667b0b423))
