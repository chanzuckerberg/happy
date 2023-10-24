# Changelog
<!-- bump -->
## [0.114.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.113.1...cli-v0.114.0) (2023-10-24)


### Features

* fixed batch delete happy stack ([#2613](https://github.com/chanzuckerberg/happy/issues/2613)) ([21f927b](https://github.com/chanzuckerberg/happy/commit/21f927b9ac095bb2645b6e2a51d914cbfe1a265c))


### Bug Fixes

* sync go versions ([#2635](https://github.com/chanzuckerberg/happy/issues/2635)) ([e479c13](https://github.com/chanzuckerberg/happy/commit/e479c136a1f2cf83b4e6b430097e74d5512f31ee))

## [0.113.1](https://github.com/chanzuckerberg/happy/compare/cli-v0.113.0...cli-v0.113.1) (2023-10-18)


### Bug Fixes

* Allow image src aws role arn to be provided for cross-account image promotion ([#2611](https://github.com/chanzuckerberg/happy/issues/2611)) ([2c69389](https://github.com/chanzuckerberg/happy/commit/2c693897054b03d530c4d23a1969da7c8558e5d1))

## [0.113.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.112.1...cli-v0.113.0) (2023-10-17)


### Miscellaneous Chores

* **cli:** Synchronize happy platform versions

## [0.112.1](https://github.com/chanzuckerberg/happy/compare/cli-v0.112.0...cli-v0.112.1) (2023-10-17)


### Miscellaneous Chores

* **cli:** Synchronize happy platform versions

## [0.112.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.111.0...cli-v0.112.0) (2023-10-16)


### Miscellaneous Chores

* **cli:** Synchronize happy platform versions

## [0.111.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.110.1...cli-v0.111.0) (2023-10-13)


### Features

* replace Gorm with Ent ORM ([#2530](https://github.com/chanzuckerberg/happy/issues/2530)) ([fa87b1a](https://github.com/chanzuckerberg/happy/commit/fa87b1a0bbd2c6b41ac4e9f013c8c60ff5409913))

## [0.110.1](https://github.com/chanzuckerberg/happy/compare/cli-v0.110.0...cli-v0.110.1) (2023-10-05)


### Bug Fixes

* Happy push doesn't push the image if service_ecrs output not present, but gives no indication of it ([#2513](https://github.com/chanzuckerberg/happy/issues/2513)) ([1f43f6d](https://github.com/chanzuckerberg/happy/commit/1f43f6d8ba9fa7b96c709a06a7d06ee85af5eb32))
* Suppress validation errors if module cannot be downloaded ([#2528](https://github.com/chanzuckerberg/happy/issues/2528)) ([fdc8e18](https://github.com/chanzuckerberg/happy/commit/fdc8e18fbaa2556fe8b5a39520173a22473279d7))

## [0.110.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.109.0...cli-v0.110.0) (2023-10-04)


### Features

* throw error to required key in get command ([#2510](https://github.com/chanzuckerberg/happy/issues/2510)) ([f47ddd1](https://github.com/chanzuckerberg/happy/commit/f47ddd132f0a64f470a452d6a9f0e15d9cd86832))

## [0.109.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.108.0...cli-v0.109.0) (2023-10-03)


### Features

* Validate happy configuration on every happy operation ([#2511](https://github.com/chanzuckerberg/happy/issues/2511)) ([c1084f2](https://github.com/chanzuckerberg/happy/commit/c1084f2eca552f76e4010f5f1673e47f5981fa15))

## [0.108.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.107.0...cli-v0.108.0) (2023-09-25)


### Miscellaneous Chores

* **cli:** Synchronize happy platform versions

## [0.107.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.106.2...cli-v0.107.0) (2023-09-22)


### Miscellaneous Chores

* **cli:** Synchronize happy platform versions

## [0.106.2](https://github.com/chanzuckerberg/happy/compare/cli-v0.106.1...cli-v0.106.2) (2023-09-21)


### Bug Fixes

* Fix ECR scanning when scanning is not enabled ([#2483](https://github.com/chanzuckerberg/happy/issues/2483)) ([9506729](https://github.com/chanzuckerberg/happy/commit/9506729d6121989b90fe58708b8bd07530e3bc0c))

## [0.106.1](https://github.com/chanzuckerberg/happy/compare/cli-v0.106.0...cli-v0.106.1) (2023-09-20)


### Bug Fixes

* Improve ECR scanning messaging ([#2480](https://github.com/chanzuckerberg/happy/issues/2480)) ([1d58703](https://github.com/chanzuckerberg/happy/commit/1d587039606ecf36212f65d24489cff811ca3588))

## [0.106.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.105.0...cli-v0.106.0) (2023-09-19)


### Features

* Warn when ECR scans fail before deployment ([#2477](https://github.com/chanzuckerberg/happy/issues/2477)) ([772d6c1](https://github.com/chanzuckerberg/happy/commit/772d6c1fafa7fbda4f12d42ab852e043bac8eed0))

## [0.105.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.104.1...cli-v0.105.0) (2023-09-15)


### Features

* Allow execution of shell commands non-interactively ([#2457](https://github.com/chanzuckerberg/happy/issues/2457)) ([cbbc2a5](https://github.com/chanzuckerberg/happy/commit/cbbc2a5bc4fe3803901465d5da6fc29386937d04))

## [0.104.1](https://github.com/chanzuckerberg/happy/compare/cli-v0.104.0...cli-v0.104.1) (2023-09-14)


### Bug Fixes

* better error reporting for happy cli and happy tf provider ([#2445](https://github.com/chanzuckerberg/happy/issues/2445)) ([894b4bd](https://github.com/chanzuckerberg/happy/commit/894b4bd804558e956e12e51b91304bb6ff12053d))

## [0.104.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.103.0...cli-v0.104.0) (2023-09-07)


### Miscellaneous Chores

* **cli:** Synchronize happy platform versions

## [0.103.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.102.2...cli-v0.103.0) (2023-09-07)


### Miscellaneous Chores

* **cli:** Synchronize happy platform versions

## [0.102.2](https://github.com/chanzuckerberg/happy/compare/cli-v0.102.1...cli-v0.102.2) (2023-09-01)


### Miscellaneous Chores

* **cli:** Synchronize happy platform versions

## [0.102.1](https://github.com/chanzuckerberg/happy/compare/cli-v0.102.0...cli-v0.102.1) (2023-09-01)


### Bug Fixes

* On bootstrap, prompt the user if docker-compose.yml already exists ([#2392](https://github.com/chanzuckerberg/happy/issues/2392)) ([5cefe53](https://github.com/chanzuckerberg/happy/commit/5cefe53bd543eedfe886df5d33cf280682ef4717))

## [0.102.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.101.0...cli-v0.102.0) (2023-09-01)


### Bug Fixes

* pass correct launch type to api /v1/stacks ([#2385](https://github.com/chanzuckerberg/happy/issues/2385)) ([3e103df](https://github.com/chanzuckerberg/happy/commit/3e103dfc4e2d96736d96ff79c053a33cef6c236c))
* using the wrong AWS profile when promoting images ([#2393](https://github.com/chanzuckerberg/happy/issues/2393)) ([43330fc](https://github.com/chanzuckerberg/happy/commit/43330fc37dadf7458f5ba4806b2d19deff12859e))

## [0.101.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.100.0...cli-v0.101.0) (2023-08-29)


### Miscellaneous Chores

* **cli:** Synchronize happy platform versions

## [0.100.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.99.0...cli-v0.100.0) (2023-08-28)


### Bug Fixes

* Notify user to restart docker engine, allow pre-release for docker compose ([#2377](https://github.com/chanzuckerberg/happy/issues/2377)) ([48745e6](https://github.com/chanzuckerberg/happy/commit/48745e66116b0c5a6e82be71b0ec2f3653f36606))

## [0.99.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.98.0...cli-v0.99.0) (2023-08-25)


### Miscellaneous Chores

* **cli:** Synchronize happy platform versions

## [0.98.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.97.1...cli-v0.98.0) (2023-08-22)


### Features

* allow ID tokens to be set using env variable ([#2341](https://github.com/chanzuckerberg/happy/issues/2341)) ([66ebc83](https://github.com/chanzuckerberg/happy/commit/66ebc835735798640f7d1ba228a9f8d223598e9c))


### Bug Fixes

* Disallow the same stack name from being used between two applications sharing compute ([#2302](https://github.com/chanzuckerberg/happy/issues/2302)) ([f05a7da](https://github.com/chanzuckerberg/happy/commit/f05a7daf878899bf152df13948ac776519abcf4f))

## [0.97.1](https://github.com/chanzuckerberg/happy/compare/cli-v0.97.0...cli-v0.97.1) (2023-08-16)


### Bug Fixes

* only include stacks for the requested app in api response ([#2269](https://github.com/chanzuckerberg/happy/issues/2269)) ([4491496](https://github.com/chanzuckerberg/happy/commit/4491496f8d81f9e4c002aef2901fbd59bc173494))

## [0.97.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.96.1...cli-v0.97.0) (2023-08-08)


### Features

* Automatically trigger the detached behavior when output is redirected ([#2240](https://github.com/chanzuckerberg/happy/issues/2240)) ([bb93722](https://github.com/chanzuckerberg/happy/commit/bb93722038c60699a5f6cdcdc1c739309aa299b8))
* Hide happy config list values by default ([#2206](https://github.com/chanzuckerberg/happy/issues/2206)) ([5e1c347](https://github.com/chanzuckerberg/happy/commit/5e1c347d5df1889bff9aab27411b0edd392f52c5))

## [0.96.1](https://github.com/chanzuckerberg/happy/compare/cli-v0.96.0...cli-v0.96.1) (2023-08-04)


### Bug Fixes

* Bootstrap doesn't generate code for all configured environments ([#2204](https://github.com/chanzuckerberg/happy/issues/2204)) ([a54f67d](https://github.com/chanzuckerberg/happy/commit/a54f67d8448f800efa8f77f4145323e62854acf2))
* Make stack and service lookup behavior consistent across happy commands ([#2201](https://github.com/chanzuckerberg/happy/issues/2201)) ([93c6479](https://github.com/chanzuckerberg/happy/commit/93c647932c41b9a3df22c6caf8c9162c69ee8d2a))

## [0.96.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.95.0...cli-v0.96.0) (2023-08-04)


### Features

* add 'happy config exec -- &lt;run-app&gt;' ([#2077](https://github.com/chanzuckerberg/happy/issues/2077)) ([3b63c7f](https://github.com/chanzuckerberg/happy/commit/3b63c7fb1c497a08efff35f58e308174c1fcf7b0))

## [0.95.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.94.0...cli-v0.95.0) (2023-08-04)


### Miscellaneous Chores

* **cli:** Synchronize happy platform versions

## [0.94.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.93.0...cli-v0.94.0) (2023-08-02)


### Features

* CCIE-1619: happy version manager v1 ([#2066](https://github.com/chanzuckerberg/happy/issues/2066)) ([816447b](https://github.com/chanzuckerberg/happy/commit/816447b5255f22cafd3795ef244e628b1af4ea4a))
* consolidate stack service in shared pkg ([#2096](https://github.com/chanzuckerberg/happy/issues/2096)) ([24d885c](https://github.com/chanzuckerberg/happy/commit/24d885cd8a8845d1e1d1934c1c3e345cfb0e951e))
* use feature flag to determine whether to use api for stacklist retrieval ([#2167](https://github.com/chanzuckerberg/happy/issues/2167)) ([5efcc18](https://github.com/chanzuckerberg/happy/commit/5efcc18612bd0cd0e27143f8a24bd5fd0773e5e5))

## [0.93.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.92.0...cli-v0.93.0) (2023-07-28)


### Features

* happy restart stack ([#2127](https://github.com/chanzuckerberg/happy/issues/2127)) ([975ad28](https://github.com/chanzuckerberg/happy/commit/975ad28d547c2a5c8b784736af1883adfc6f0f43))

## [0.92.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.91.1...cli-v0.92.0) (2023-07-10)


### Features

* improve conditional job runs to prevent dependency merge bot from merging broken code ([#2019](https://github.com/chanzuckerberg/happy/issues/2019)) ([956b18c](https://github.com/chanzuckerberg/happy/commit/956b18c3a574301a76353cb20934c47817500440))


### Bug Fixes

* 'happy list --output json' returns an invalid value ([#2023](https://github.com/chanzuckerberg/happy/issues/2023)) ([47873e7](https://github.com/chanzuckerberg/happy/commit/47873e756a736d93a57a36add20b885fb74de301))
* Tasks when executed do not receive environment information ([#2026](https://github.com/chanzuckerberg/happy/issues/2026)) ([c281786](https://github.com/chanzuckerberg/happy/commit/c281786cdcb9537c7f57ae537fcd91c3b167d9c2))

## [0.91.1](https://github.com/chanzuckerberg/happy/compare/cli-v0.91.0...cli-v0.91.1) (2023-07-03)


### Bug Fixes

* broken filepath in shared stack package ([#1998](https://github.com/chanzuckerberg/happy/issues/1998)) ([dd7e714](https://github.com/chanzuckerberg/happy/commit/dd7e714b06247d97e4a9785f2dd238474f8cca58))

## [0.91.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.90.1...cli-v0.91.0) (2023-06-30)


### Features

* fix duplicates returning from API ([#1990](https://github.com/chanzuckerberg/happy/issues/1990)) ([58a0aa7](https://github.com/chanzuckerberg/happy/commit/58a0aa745a9646d34fc7adc418001d8f63d65047))

## [0.90.1](https://github.com/chanzuckerberg/happy/compare/cli-v0.90.0...cli-v0.90.1) (2023-06-28)


### Bug Fixes

* reuse datastructure ([#1982](https://github.com/chanzuckerberg/happy/issues/1982)) ([4bd98db](https://github.com/chanzuckerberg/happy/commit/4bd98db581e8a72a8ad9c6032126215eac220cc3))

## [0.90.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.89.1...cli-v0.90.0) (2023-06-22)


### Features

* reuse happy client ([#1960](https://github.com/chanzuckerberg/happy/issues/1960)) ([fc3991d](https://github.com/chanzuckerberg/happy/commit/fc3991d0670579e34013e854e6a5a4f3fc4e189e))

## [0.89.1](https://github.com/chanzuckerberg/happy/compare/cli-v0.89.0...cli-v0.89.1) (2023-06-21)


### Bug Fixes

* docker-compose.yml doesn't allow for a name attribute on services ([#1959](https://github.com/chanzuckerberg/happy/issues/1959)) ([4e18e5e](https://github.com/chanzuckerberg/happy/commit/4e18e5e082c9277348f3cff31ca85f8db7fdd66a))

## [0.89.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.88.0...cli-v0.89.0) (2023-06-21)


### Features

* consolidate stack operations on cli and hapi ([#1867](https://github.com/chanzuckerberg/happy/issues/1867)) ([a4a8b5d](https://github.com/chanzuckerberg/happy/commit/a4a8b5db6ce01811592278107da58cb0aba5fc5b))

## [0.88.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.87.0...cli-v0.88.0) (2023-06-20)


### Features

* Implement "happy bootstrap" to init happy configuration on the existing GitHub repo ([#1866](https://github.com/chanzuckerberg/happy/issues/1866)) ([6cd3084](https://github.com/chanzuckerberg/happy/commit/6cd3084fd720f972f4434e82db2112b225230ee3))


### Bug Fixes

* messaging around slice validation ([#1870](https://github.com/chanzuckerberg/happy/issues/1870)) ([ae47d5a](https://github.com/chanzuckerberg/happy/commit/ae47d5a4958f27348869077359de497618ebd919))

## [0.87.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.86.0...cli-v0.87.0) (2023-06-07)


### Features

* move stacks to happy eks ([#1843](https://github.com/chanzuckerberg/happy/issues/1843)) ([0e6b5f0](https://github.com/chanzuckerberg/happy/commit/0e6b5f0d28e560768c4eea17bb5f32bd699945a8))

## [0.86.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.85.1...cli-v0.86.0) (2023-06-06)


### Features

* Consume and produce docker-compose.yml via happy ingest and happy generate ([#1852](https://github.com/chanzuckerberg/happy/issues/1852)) ([addb506](https://github.com/chanzuckerberg/happy/commit/addb506505db527e6c08c71a33717cb38fd1b570))

## [0.85.1](https://github.com/chanzuckerberg/happy/compare/cli-v0.85.0...cli-v0.85.1) (2023-06-05)


### Bug Fixes

* slices failing validation of config services ([#1853](https://github.com/chanzuckerberg/happy/issues/1853)) ([18070e1](https://github.com/chanzuckerberg/happy/commit/18070e1defa6f37464787ab9593ad4aa39c5ccf1))

## [0.85.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.84.0...cli-v0.85.0) (2023-06-05)


### Features

* Add sidecar support to happy logs ([#1844](https://github.com/chanzuckerberg/happy/issues/1844)) ([12949d7](https://github.com/chanzuckerberg/happy/commit/12949d7b027721b69b0acf4e2b0f71dc5c4b1fb9))

## [0.84.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.83.0...cli-v0.84.0) (2023-06-02)


### Features

* Add sidecar support to happy shell ([#1842](https://github.com/chanzuckerberg/happy/issues/1842)) ([9c52320](https://github.com/chanzuckerberg/happy/commit/9c5232066acebd6562541de03b91028bca1fc8bb))
* CCIE-1507 do not check docker daemon when not dealing with tags in happy create/update ([#1837](https://github.com/chanzuckerberg/happy/issues/1837)) ([cb360ca](https://github.com/chanzuckerberg/happy/commit/cb360caed01df403b4575aebc96a52816a836d94))

## [0.83.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.82.1...cli-v0.83.0) (2023-06-01)


### Features

* Implement "happy infra refresh" to refresh terraform scripts ([#1832](https://github.com/chanzuckerberg/happy/issues/1832)) ([52fc23d](https://github.com/chanzuckerberg/happy/commit/52fc23dc3517c7fbe209aa82ac95ee9cf41c7e9f))
* multistack destroy; refactor destroy ([#1833](https://github.com/chanzuckerberg/happy/issues/1833)) ([7c37665](https://github.com/chanzuckerberg/happy/commit/7c3766504521025b4b8bfc8d07264b723ac5a4f6))

## [0.82.1](https://github.com/chanzuckerberg/happy/compare/cli-v0.82.0...cli-v0.82.1) (2023-05-31)


### Bug Fixes

* can't use persistenprerun ([#1829](https://github.com/chanzuckerberg/happy/issues/1829)) ([68075be](https://github.com/chanzuckerberg/happy/commit/68075bea6e23daef8834584d236bb44be5786a39))

## [0.82.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.81.0...cli-v0.82.0) (2023-05-31)


### Features

* Implement "happy infra validate" to validate terraform scripts ([#1824](https://github.com/chanzuckerberg/happy/issues/1824)) ([a57c977](https://github.com/chanzuckerberg/happy/commit/a57c9775cc436e92e3475edb6b880b49e07807b0))


### Bug Fixes

* force update; deadlock changes ([#1826](https://github.com/chanzuckerberg/happy/issues/1826)) ([bffa24a](https://github.com/chanzuckerberg/happy/commit/bffa24a267e768b9ea54278ce55576848318e18b))

## [0.81.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.80.0...cli-v0.81.0) (2023-05-30)


### Features

* Example of task usage in happy EKS ([#1776](https://github.com/chanzuckerberg/happy/issues/1776)) ([2af7c7f](https://github.com/chanzuckerberg/happy/commit/2af7c7faa87938ea859db26fe143eca429f61d86))


### Bug Fixes

* [bug] Validate credentials before stack operations and prompt user to log in and create a new token on token absence or prior invalidation ([#1806](https://github.com/chanzuckerberg/happy/issues/1806)) ([e23146a](https://github.com/chanzuckerberg/happy/commit/e23146ac94363551ff5990c533637f61344d5f94))

## [0.80.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.79.1...cli-v0.80.0) (2023-05-24)


### Features

* Collect stack configuration from existing terraform code and store in happy config ([#1761](https://github.com/chanzuckerberg/happy/issues/1761)) ([56dd781](https://github.com/chanzuckerberg/happy/commit/56dd7819d44b6464e2dd0d43ab27d77411fcf680))


### Bug Fixes

* Follow happy delete logs to prevent race conditions ([#1755](https://github.com/chanzuckerberg/happy/issues/1755)) ([d9e786e](https://github.com/chanzuckerberg/happy/commit/d9e786ed53d8ffa28d8f2e59b8aea2aea4a5aa70))
* session manager not needed unless using shell command ([#1757](https://github.com/chanzuckerberg/happy/issues/1757)) ([5a0fbb5](https://github.com/chanzuckerberg/happy/commit/5a0fbb59feae94ae10157a4de2e0ead4abebfbb0))

## [0.79.1](https://github.com/chanzuckerberg/happy/compare/cli-v0.79.0...cli-v0.79.1) (2023-05-10)


### Bug Fixes

* When more than one service is specified, and the settings structure is inconsistent, happy infra generate errors out ([#1751](https://github.com/chanzuckerberg/happy/issues/1751)) ([ea166c2](https://github.com/chanzuckerberg/happy/commit/ea166c20cd6a52e0ef82a53554261f0055d680ed))

## [0.79.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.78.0...cli-v0.79.0) (2023-05-09)


### Features

* Sidecar support for services ([#1727](https://github.com/chanzuckerberg/happy/issues/1727)) ([8c5c884](https://github.com/chanzuckerberg/happy/commit/8c5c884804a4e88d1e3163f266127e6ddb336c05))


### Bug Fixes

* add happy_config_secret back in meta ([#1743](https://github.com/chanzuckerberg/happy/issues/1743)) ([2cfac76](https://github.com/chanzuckerberg/happy/commit/2cfac76ee0cb8f28f32c94a4842818d493743d0c))
* Refresh EKS credentials after a lengthy docker build ([#1728](https://github.com/chanzuckerberg/happy/issues/1728)) ([b9d422b](https://github.com/chanzuckerberg/happy/commit/b9d422beea1930d5806dcf6186d7fce3092c0fdd))

## [0.78.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.77.0...cli-v0.78.0) (2023-05-08)


### Features

* Implement Happy cli basic mode (terraform code is generated) ([#1684](https://github.com/chanzuckerberg/happy/issues/1684)) ([ca41c53](https://github.com/chanzuckerberg/happy/commit/ca41c538bfb99491028ab07b55308c88fc3d4a03))

## [0.77.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.76.0...cli-v0.77.0) (2023-05-02)


### Features

* add command to see the configured CI roles for env ([#1686](https://github.com/chanzuckerberg/happy/issues/1686)) ([a249cc0](https://github.com/chanzuckerberg/happy/commit/a249cc0a4fc61af413312b300f1fc4695529ee2e))

## [0.76.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.75.0...cli-v0.76.0) (2023-04-27)


### Features

* CCIE-960 do not require docker to be running for functions that don t use it ([#1659](https://github.com/chanzuckerberg/happy/issues/1659)) ([7c53ee6](https://github.com/chanzuckerberg/happy/commit/7c53ee6492300f89724182701a305d65c62b1aa1))

## [0.75.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.74.0...cli-v0.75.0) (2023-04-27)


### Features

* allow for stacks to migrate container artifacts ([#1619](https://github.com/chanzuckerberg/happy/issues/1619)) ([09cea95](https://github.com/chanzuckerberg/happy/commit/09cea95566c41b34f12a1d2f858ff3bef8d598a6))


### Bug Fixes

* merged too fast ([#1669](https://github.com/chanzuckerberg/happy/issues/1669)) ([2b28c1b](https://github.com/chanzuckerberg/happy/commit/2b28c1b91eadbab50ffba5252ce95226132df8cc))

## [0.74.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.73.3...cli-v0.74.0) (2023-04-27)


### Features

* Happy CLI integration test ([#1662](https://github.com/chanzuckerberg/happy/issues/1662)) ([a3c4d2c](https://github.com/chanzuckerberg/happy/commit/a3c4d2ce28a095f47d9c66c9ddfd24b231b864b6))

## [0.73.3](https://github.com/chanzuckerberg/happy/compare/cli-v0.73.2...cli-v0.73.3) (2023-04-27)


### Miscellaneous Chores

* **cli:** Synchronize happy platform versions

## [0.73.2](https://github.com/chanzuckerberg/happy/compare/cli-v0.73.1...cli-v0.73.2) (2023-04-24)


### Bug Fixes

* state drilling with dryrun; use context ([#1607](https://github.com/chanzuckerberg/happy/issues/1607)) ([a75376a](https://github.com/chanzuckerberg/happy/commit/a75376a849940d9cdf45accbc1ec0357dbd0c3f8))

## [0.73.1](https://github.com/chanzuckerberg/happy/compare/cli-v0.73.0...cli-v0.73.1) (2023-04-20)


### Miscellaneous Chores

* **cli:** Synchronize happy platform versions

## [0.73.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.72.0...cli-v0.73.0) (2023-04-20)


### Features

* filter the stacks by app; display git info ([#1594](https://github.com/chanzuckerberg/happy/issues/1594)) ([665f35c](https://github.com/chanzuckerberg/happy/commit/665f35c39d7eff37ff8d0bca34f37db08f0eb753))

## [0.72.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.71.0...cli-v0.72.0) (2023-04-20)


### Features

* Happy debug feature support for EKS ([#1592](https://github.com/chanzuckerberg/happy/issues/1592)) ([08eb06a](https://github.com/chanzuckerberg/happy/commit/08eb06acda5990fe5c4fd4aedc57eaf7179233d0))

## [0.71.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.70.1...cli-v0.71.0) (2023-04-18)


### Features

* clean up playground happy stacks on friday night ([#1496](https://github.com/chanzuckerberg/happy/issues/1496)) ([fffc252](https://github.com/chanzuckerberg/happy/commit/fffc2523ff3eeebe7fe2878541150f740f76a477))
* Support the happy events feature to visualize events from key applicaiton levels ([#1579](https://github.com/chanzuckerberg/happy/issues/1579)) ([367d958](https://github.com/chanzuckerberg/happy/commit/367d958486536d2812940865d314bd1cd2490d23))

## [0.70.1](https://github.com/chanzuckerberg/happy/compare/cli-v0.70.0...cli-v0.70.1) (2023-04-12)


### Bug Fixes

* Breaking change in a kubernetes api ([#1565](https://github.com/chanzuckerberg/happy/issues/1565)) ([5967f4a](https://github.com/chanzuckerberg/happy/commit/5967f4a6680ed9d4495cc241b843f88a40c7f8cc))
* Happy update --dry-run deletes the stack ([#1563](https://github.com/chanzuckerberg/happy/issues/1563)) ([ff0e840](https://github.com/chanzuckerberg/happy/commit/ff0e840523a712a5f31110d6af83e94f98a21fd0))

## [0.70.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.69.2...cli-v0.70.0) (2023-04-10)


### Features

* Happy Service integrity check ([#1495](https://github.com/chanzuckerberg/happy/issues/1495)) ([29f7804](https://github.com/chanzuckerberg/happy/commit/29f780437bf28f4ae9c309ad47f1dd752b156559))


### Bug Fixes

* Dry Run is broken for "happy update" ([#1525](https://github.com/chanzuckerberg/happy/issues/1525)) ([70e75ec](https://github.com/chanzuckerberg/happy/commit/70e75ecd394e736963a6504d91371d8c976c480c))

## [0.69.2](https://github.com/chanzuckerberg/happy/compare/cli-v0.69.1...cli-v0.69.2) (2023-04-07)


### Bug Fixes

* use aws credentials from request ([#1493](https://github.com/chanzuckerberg/happy/issues/1493)) ([8608647](https://github.com/chanzuckerberg/happy/commit/8608647a6e7e8ee2024f211a12fcff7fdf4fae4e))

## [0.69.1](https://github.com/chanzuckerberg/happy/compare/cli-v0.69.0...cli-v0.69.1) (2023-04-07)


### Miscellaneous Chores

* **cli:** Synchronize happy platform versions

## [0.69.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.68.0...cli-v0.69.0) (2023-04-07)


### Features

* Expose stack TFE status, TFE Url, and Endpoints through HAPI ([#1469](https://github.com/chanzuckerberg/happy/issues/1469)) ([820396a](https://github.com/chanzuckerberg/happy/commit/820396ac31c9416ba49afe0ac73dfd816ad2e9c4))
* Happy config: Make aws region configurable ([#1487](https://github.com/chanzuckerberg/happy/issues/1487)) ([b70ad5e](https://github.com/chanzuckerberg/happy/commit/b70ad5e43e020965b7683eec82e62aa1ca02bff5))
* Remove happy config from backend ([#1472](https://github.com/chanzuckerberg/happy/issues/1472)) ([7421240](https://github.com/chanzuckerberg/happy/commit/7421240f96be6b891b43be893429b7d62e574c80))

## [0.68.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.67.0...cli-v0.68.0) (2023-04-05)


### Features

* Move backend, workspace_repo package to shared ([#1467](https://github.com/chanzuckerberg/happy/issues/1467)) ([d0b64ed](https://github.com/chanzuckerberg/happy/commit/d0b64edd690e91690438de6c35671a90d248f9ba))


### Bug Fixes

* happy list no longer displays service endpoints ([#1466](https://github.com/chanzuckerberg/happy/issues/1466)) ([d057963](https://github.com/chanzuckerberg/happy/commit/d05796312bce308ff34759f93cb871f3b10155c6))

## [0.67.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.66.2...cli-v0.67.0) (2023-03-29)


### Features

* CCIE-858 - Part 2 - Lock happy version in repo ([#1454](https://github.com/chanzuckerberg/happy/issues/1454)) ([6f49ea1](https://github.com/chanzuckerberg/happy/commit/6f49ea169e49d259a3cabe82fa94ce3655f765a8))
* CCIE-858: lock the version of happy cli in a repo like fogg ([#1371](https://github.com/chanzuckerberg/happy/issues/1371)) ([dddb799](https://github.com/chanzuckerberg/happy/commit/dddb799092fadb5a6443577d6036b547874ca442))

## [0.66.2](https://github.com/chanzuckerberg/happy/compare/cli-v0.66.1...cli-v0.66.2) (2023-03-28)


### Bug Fixes

* fix broken default tags in the default happy update/create cmds ([#1447](https://github.com/chanzuckerberg/happy/issues/1447)) ([588a804](https://github.com/chanzuckerberg/happy/commit/588a80408c843d2649d97bee76575881d72517c2))

## [0.66.1](https://github.com/chanzuckerberg/happy/compare/cli-v0.66.0...cli-v0.66.1) (2023-03-28)


### Bug Fixes

* accidently adding empty tags ([#1443](https://github.com/chanzuckerberg/happy/issues/1443)) ([ba589a0](https://github.com/chanzuckerberg/happy/commit/ba589a02ceae4cc471e392e4989734be56305ce8))
* range over wrong argument ([#1445](https://github.com/chanzuckerberg/happy/issues/1445)) ([c880fbb](https://github.com/chanzuckerberg/happy/commit/c880fbb976bdc591f8f0621a207a291a7812edae))

## [0.66.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.65.0...cli-v0.66.0) (2023-03-27)


### Bug Fixes

* tags flag for push and other create/update commands ([#1442](https://github.com/chanzuckerberg/happy/issues/1442)) ([0a1257f](https://github.com/chanzuckerberg/happy/commit/0a1257f6e75191e7f963d6d5a882b6bff29d5dd9))
* update addtags function to use latest ECR naming convention ([#1441](https://github.com/chanzuckerberg/happy/issues/1441)) ([e4b4d91](https://github.com/chanzuckerberg/happy/commit/e4b4d9166531ffc352f10ae2ac1a912d0d272652))

## [0.65.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.64.0...cli-v0.65.0) (2023-03-15)


### Features

* CCIE-900 Automatically check whether Happy is up to date ([#1355](https://github.com/chanzuckerberg/happy/issues/1355)) ([7cec2dd](https://github.com/chanzuckerberg/happy/commit/7cec2dd277b1eaf995780d9cd4ffdba3fcbb46fe))

## [0.64.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.63.1...cli-v0.64.0) (2023-03-13)


### Features

* 'happy logs' integration with Cloudwatch Insights ([#1315](https://github.com/chanzuckerberg/happy/issues/1315)) ([9ff4861](https://github.com/chanzuckerberg/happy/commit/9ff48617f79273457018d21de2a1ad78b9109a07))

## [0.63.1](https://github.com/chanzuckerberg/happy/compare/cli-v0.63.0...cli-v0.63.1) (2023-03-08)


### Miscellaneous Chores

* **cli:** Synchronize happy platform versions

## [0.63.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.62.1...cli-v0.63.0) (2023-03-07)


### Features

* Detailed stack workspace deletion message ([#1300](https://github.com/chanzuckerberg/happy/issues/1300)) ([879456e](https://github.com/chanzuckerberg/happy/commit/879456efb75efb560396655e52f0512f6d593325))

## [0.62.1](https://github.com/chanzuckerberg/happy/compare/cli-v0.62.0...cli-v0.62.1) (2023-03-07)


### Bug Fixes

* Shorten the happy version number in TFE messages ([#1297](https://github.com/chanzuckerberg/happy/issues/1297)) ([14338f2](https://github.com/chanzuckerberg/happy/commit/14338f2ce05c167fc0685848df5bcacfa8943328))

## [0.62.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.61.0...cli-v0.62.0) (2023-03-07)


### Features

* Include a specific message in creator workspace runs and add happy version number to the message ([#1296](https://github.com/chanzuckerberg/happy/issues/1296)) ([9e32d6f](https://github.com/chanzuckerberg/happy/commit/9e32d6f12e71614f5adecad0d047870a74b46d78))


### Bug Fixes

* push command to push to stack ECRs ([#1294](https://github.com/chanzuckerberg/happy/issues/1294)) ([87dc9df](https://github.com/chanzuckerberg/happy/commit/87dc9dff0ad0c9fe3cba2015e17eeb1572562855))

## [0.61.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.60.0...cli-v0.61.0) (2023-03-07)


### Features

* Replace generic message with an actual one ([#1293](https://github.com/chanzuckerberg/happy/issues/1293)) ([023ee1c](https://github.com/chanzuckerberg/happy/commit/023ee1c0d99a4fdba3f531cfb5842b038bf0c478))


### Bug Fixes

* don't throw error in create ([#1291](https://github.com/chanzuckerberg/happy/issues/1291)) ([ae727e4](https://github.com/chanzuckerberg/happy/commit/ae727e4eff3fd4789cd92e401f0e377446ad37e5))

## [0.60.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.59.0...cli-v0.60.0) (2023-03-07)


### ⚠ BREAKING CHANGES

* refactor update/clean; autocreate ECR ([#1232](https://github.com/chanzuckerberg/happy/issues/1232))

### Features

* refactor update/clean; autocreate ECR ([#1232](https://github.com/chanzuckerberg/happy/issues/1232)) ([b498074](https://github.com/chanzuckerberg/happy/commit/b4980740c3ddc716abe530fb2112dfe41bc6ab60))

## [0.59.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.58.0...cli-v0.59.0) (2023-02-28)


### Features

* CCIE-1114: Remove terraform version config ([#1256](https://github.com/chanzuckerberg/happy/issues/1256)) ([09ce85a](https://github.com/chanzuckerberg/happy/commit/09ce85a7e1e8e9aa4db0abd992908ffeecd87452))

## [0.58.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.57.0...cli-v0.58.0) (2023-02-24)


### Features

* Annotate k8s resources created by happy with stack ownership labels ([#1247](https://github.com/chanzuckerberg/happy/issues/1247)) ([4403cd8](https://github.com/chanzuckerberg/happy/commit/4403cd8404ccdec96936bb033a94a3d7a2f4e58b))

## [0.57.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.56.1...cli-v0.57.0) (2023-02-23)


### Features

* List of all AWS infra associated with a k8s happy stack ([#1217](https://github.com/chanzuckerberg/happy/issues/1217)) ([83586fb](https://github.com/chanzuckerberg/happy/commit/83586fb2950a30677884245c3dc6cc8efa4968a7))

## [0.56.1](https://github.com/chanzuckerberg/happy/compare/cli-v0.56.0...cli-v0.56.1) (2023-02-21)


### Bug Fixes

* Without go mod tidy golangci-lint breaks (claims there are no go files) ([#1210](https://github.com/chanzuckerberg/happy/issues/1210)) ([836a038](https://github.com/chanzuckerberg/happy/commit/836a038ef0913b167082f4ca95d47051063a7e18))

## [0.56.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.55.1...cli-v0.56.0) (2023-02-17)


### Features

* allow users to create bypasses for their OIDC ([#1149](https://github.com/chanzuckerberg/happy/issues/1149)) ([078ee17](https://github.com/chanzuckerberg/happy/commit/078ee17b36436ce92b5ad0efdade143d1f306879))

## [0.55.1](https://github.com/chanzuckerberg/happy/compare/cli-v0.55.0...cli-v0.55.1) (2023-02-13)


### Bug Fixes

* Happy addtags tags images outside of the application ([#1142](https://github.com/chanzuckerberg/happy/issues/1142)) ([bcd5b28](https://github.com/chanzuckerberg/happy/commit/bcd5b286948d6b4a717c448b1587bf42bdbdea0c))

## [0.55.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.54.1...cli-v0.55.0) (2023-02-13)


### Miscellaneous Chores

* **cli:** Synchronize happy platform versions

## [0.54.1](https://github.com/chanzuckerberg/happy/compare/cli-v0.54.0...cli-v0.54.1) (2023-02-13)


### Bug Fixes

* Happy version command stopped working ([#1138](https://github.com/chanzuckerberg/happy/issues/1138)) ([bfe75d4](https://github.com/chanzuckerberg/happy/commit/bfe75d4ea59d5f8bbe49561d6aa86c7c2803490d))

## [0.54.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.53.6...cli-v0.54.0) (2023-02-13)


### ⚠ BREAKING CHANGES

* inject happy config to stacks ([#1108](https://github.com/chanzuckerberg/happy/issues/1108))

### Features

* inject happy config to stacks ([#1108](https://github.com/chanzuckerberg/happy/issues/1108)) ([9cb49c7](https://github.com/chanzuckerberg/happy/commit/9cb49c7f7bd6819541510e4f31ab5fd112579457))

## [0.53.6](https://github.com/chanzuckerberg/happy/compare/cli-v0.53.5...cli-v0.53.6) (2023-02-10)


### Bug Fixes

* update happy api oidc client id ([#1133](https://github.com/chanzuckerberg/happy/issues/1133)) ([d27a82f](https://github.com/chanzuckerberg/happy/commit/d27a82f6f0bd376cd9ae81ae1b9a1e863ad8fd6f))

## [0.53.5](https://github.com/chanzuckerberg/happy/compare/cli-v0.53.4...cli-v0.53.5) (2023-02-10)


### Miscellaneous Chores

* **cli:** Synchronize happy platform versions

## [0.53.4](https://github.com/chanzuckerberg/happy/compare/cli-v0.53.3...cli-v0.53.4) (2023-02-10)


### Miscellaneous Chores

* **cli:** Synchronize happy platform versions

## [0.53.3](https://github.com/chanzuckerberg/happy/compare/cli-v0.53.2...cli-v0.53.3) (2023-02-09)


### Miscellaneous Chores

* **cli:** Synchronize happy platform versions

## [0.53.2](https://github.com/chanzuckerberg/happy/compare/cli-v0.53.1...cli-v0.53.2) (2023-02-09)


### Bug Fixes

* find git root using rev-parse ([#1113](https://github.com/chanzuckerberg/happy/issues/1113)) ([9f16ba6](https://github.com/chanzuckerberg/happy/commit/9f16ba6907b10159ec4db2c19ff28c80628e6139))

## [0.53.1](https://github.com/chanzuckerberg/happy/compare/cli-v0.53.0...cli-v0.53.1) (2023-02-09)


### Miscellaneous Chores

* **cli:** Synchronize happy platform versions

## [0.53.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.52.0...cli-v0.53.0) (2023-02-09)


### Miscellaneous Chores

* **cli:** Synchronize happy platform versions

## [0.52.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.51.0...cli-v0.52.0) (2023-02-08)


### Features

* use query string for GET requests to happy api ([#1101](https://github.com/chanzuckerberg/happy/issues/1101)) ([7a18eb8](https://github.com/chanzuckerberg/happy/commit/7a18eb8dd5bc2eaebdb246dbebd44f4c389b17e2))

## [0.51.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.50.2...cli-v0.51.0) (2023-02-08)


### Miscellaneous Chores

* **cli:** Synchronize happy platform versions

## [0.50.2](https://github.com/chanzuckerberg/happy/compare/cli-v0.50.1...cli-v0.50.2) (2023-01-30)


### Bug Fixes

* Remove duplicate oidc code and keep aws-go-sdk dependency out ([#1031](https://github.com/chanzuckerberg/happy/issues/1031)) ([12c46fa](https://github.com/chanzuckerberg/happy/commit/12c46fa8adff5f193b7064ee1673c84db16bfb8f))

## [0.50.1](https://github.com/chanzuckerberg/happy/compare/cli-v0.50.0...cli-v0.50.1) (2023-01-30)


### Bug Fixes

* Do not stop when unable to determine if git tree is dirty ([#1028](https://github.com/chanzuckerberg/happy/issues/1028)) ([81a4ca2](https://github.com/chanzuckerberg/happy/commit/81a4ca2b6df0399e978306b0590f11479a83cf99))

## [0.50.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.49.0...cli-v0.50.0) (2023-01-27)


### Features

* Abstract out kubernetes authentication ([#1024](https://github.com/chanzuckerberg/happy/issues/1024)) ([e5712ef](https://github.com/chanzuckerberg/happy/commit/e5712ef334bcb7d60c07c36ed1f6afe22566a1d9))
* Move backend interfaces to a shared module ([#1026](https://github.com/chanzuckerberg/happy/issues/1026)) ([b0921a8](https://github.com/chanzuckerberg/happy/commit/b0921a834e52895f0cd92eebf7b65fc56f7425fc))

## [0.49.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.48.0...cli-v0.49.0) (2023-01-24)


### Miscellaneous Chores

* **cli:** Synchronize happy platform versions

## [0.48.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.47.1...cli-v0.48.0) (2023-01-19)


### Features

* send aws creds in requests to api ([#962](https://github.com/chanzuckerberg/happy/issues/962)) ([01c6b79](https://github.com/chanzuckerberg/happy/commit/01c6b79d1b4ea27ee54d3dc96a9a247075189aa0))

## [0.47.1](https://github.com/chanzuckerberg/happy/compare/cli-v0.47.0...cli-v0.47.1) (2023-01-17)


### Miscellaneous Chores

* **cli:** Synchronize happy platform versions

## [0.47.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.46.1...cli-v0.47.0) (2023-01-17)


### Miscellaneous Chores

* **cli:** Synchronize happy platform versions

## [0.46.1](https://github.com/chanzuckerberg/happy/compare/cli-v0.46.0...cli-v0.46.1) (2023-01-09)


### Bug Fixes

* stack only getting added when empty ([#958](https://github.com/chanzuckerberg/happy/issues/958)) ([bcefb57](https://github.com/chanzuckerberg/happy/commit/bcefb57f5e4538e5f1c64cd177d79e1718f6bf88))

## [0.46.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.45.0...cli-v0.46.0) (2023-01-04)


### Miscellaneous Chores

* **cli:** Synchronize happy platform versions

## [0.45.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.44.0...cli-v0.45.0) (2022-12-21)


### Features

* add api meta-command ([#903](https://github.com/chanzuckerberg/happy/issues/903)) ([b81871b](https://github.com/chanzuckerberg/happy/commit/b81871bf694063ce172267e3dcbfe08d737f4120))

## [0.44.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.43.2...cli-v0.44.0) (2022-12-20)


### Features

* send auth header in api requests ([#785](https://github.com/chanzuckerberg/happy/issues/785)) ([d83c9b3](https://github.com/chanzuckerberg/happy/commit/d83c9b3c57950b1747d8233166e276d883cda4a7))

## [0.43.2](https://github.com/chanzuckerberg/happy/compare/cli-v0.43.1...cli-v0.43.2) (2022-12-20)


### Miscellaneous Chores

* **cli:** Synchronize happy platform versions

## [0.43.1](https://github.com/chanzuckerberg/happy/compare/cli-v0.43.0...cli-v0.43.1) (2022-12-19)


### Miscellaneous Chores

* **cli:** Synchronize happy platform versions

## [0.43.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.42.1...cli-v0.43.0) (2022-12-16)


### Miscellaneous Chores

* **cli:** Synchronize happy platform versions

## [0.42.1](https://github.com/chanzuckerberg/happy/compare/cli-v0.42.0...cli-v0.42.1) (2022-12-13)


### Miscellaneous Chores

* **cli:** Synchronize happy platform versions

## [0.42.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.41.5...cli-v0.42.0) (2022-12-12)


### Features

* Warning for happy create/update to encourage pushing clean commits ([#866](https://github.com/chanzuckerberg/happy/issues/866)) ([dbfafb7](https://github.com/chanzuckerberg/happy/commit/dbfafb79cd8d8face4d669ba63eb7b77d7afdf3e))


### Bug Fixes

* Fix a linter problem due to use of deprecated exec.Stream() ([#864](https://github.com/chanzuckerberg/happy/issues/864)) ([2beef7d](https://github.com/chanzuckerberg/happy/commit/2beef7d596e8a57659d47f13eddc3f51360ab8fe))

## [0.41.5](https://github.com/chanzuckerberg/happy/compare/cli-v0.41.4...cli-v0.41.5) (2022-12-12)


### Miscellaneous Chores

* **cli:** Synchronize happy platform versions

## [0.41.4](https://github.com/chanzuckerberg/happy/compare/cli-v0.41.3...cli-v0.41.4) (2022-12-08)


### Bug Fixes

* Service not found when too many services exist on an ECS cluster ([#842](https://github.com/chanzuckerberg/happy/issues/842)) ([945271d](https://github.com/chanzuckerberg/happy/commit/945271d2faa072b3e957c563a70ddacc6e8ae40f))

## [0.41.3](https://github.com/chanzuckerberg/happy/compare/cli-v0.41.2...cli-v0.41.3) (2022-12-07)


### Bug Fixes

* Pin go-slug to avoid the relative path issue (again) ([#833](https://github.com/chanzuckerberg/happy/issues/833)) ([427f67a](https://github.com/chanzuckerberg/happy/commit/427f67a7b49525f432a6107f630eea909c11d001))

## [0.41.2](https://github.com/chanzuckerberg/happy/compare/cli-v0.41.1...cli-v0.41.2) (2022-12-02)


### Bug Fixes

* Remove the misleading .env loading message ([#815](https://github.com/chanzuckerberg/happy/issues/815)) ([a029fcd](https://github.com/chanzuckerberg/happy/commit/a029fcde26affa9b3ac269d510894ad0131be415))

## [0.41.1](https://github.com/chanzuckerberg/happy/compare/cli-v0.41.0...cli-v0.41.1) (2022-11-17)


### Miscellaneous Chores

* **cli:** Synchronize happy platform versions

## [0.41.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.40.1...cli-v0.41.0) (2022-11-16)


### Bug Fixes

* Migration task fails to start ([#759](https://github.com/chanzuckerberg/happy/issues/759)) ([fc97c75](https://github.com/chanzuckerberg/happy/commit/fc97c751a2fe236b7598c6e673307001e18fd4bd))

## [0.40.1](https://github.com/chanzuckerberg/happy/compare/cli-v0.40.0...cli-v0.40.1) (2022-11-03)


### Bug Fixes

* [CCIE-714] Enforce consistency in priorities when specifying an env ([#673](https://github.com/chanzuckerberg/happy/issues/673)) ([1d77b94](https://github.com/chanzuckerberg/happy/commit/1d77b9453b94750107437f436ca899d5cfb4c11a))

## [0.40.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.39.0...cli-v0.40.0) (2022-11-03)


### Miscellaneous Chores

* **cli:** Synchronize happy platform versions

## [0.39.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.38.0...cli-v0.39.0) (2022-10-25)


### Features

* Extract common interface from ECSComputeLogPrinter and use it in PrintLogs() for k8s implementation ([#700](https://github.com/chanzuckerberg/happy/issues/700)) ([c1191b5](https://github.com/chanzuckerberg/happy/commit/c1191b5f9e9a52bfcd5c3935ec224b0690e07046))

## [0.38.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.37.0...cli-v0.38.0) (2022-10-24)


### Features

* Add support for non-string terraform outputs ([#659](https://github.com/chanzuckerberg/happy/issues/659)) ([1d04832](https://github.com/chanzuckerberg/happy/commit/1d048323650daed62330163c603c5bfdce73db48))
* Kubernetes migration and deletion task support ([#686](https://github.com/chanzuckerberg/happy/issues/686)) ([5ca47b3](https://github.com/chanzuckerberg/happy/commit/5ca47b38a5597f716a3bfd1b26f12ffcdafa549d))
* roll out config feature to cli ([#660](https://github.com/chanzuckerberg/happy/issues/660)) ([a72c965](https://github.com/chanzuckerberg/happy/commit/a72c965f6bd2c9113c8152c9155330971e808b46))


### Bug Fixes

* Add build make file changes and update root ignore file to ensure we do not add binary files ([#683](https://github.com/chanzuckerberg/happy/issues/683)) ([d1ef4e2](https://github.com/chanzuckerberg/happy/commit/d1ef4e23fe41baa03bfdee5e6931446391ac9029))

## [0.37.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.36.0...cli-v0.37.0) (2022-10-13)


### Miscellaneous Chores

* **cli:** Synchronize happy platform versions

## [0.36.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.35.2...cli-v0.36.0) (2022-10-12)


### Miscellaneous Chores

* **cli:** Synchronize happy platform versions

## [0.35.2](https://github.com/chanzuckerberg/happy/compare/cli-v0.35.1...cli-v0.35.2) (2022-10-12)


### Bug Fixes

* reorder release steps so tag is present ([#654](https://github.com/chanzuckerberg/happy/issues/654)) ([9d1e55d](https://github.com/chanzuckerberg/happy/commit/9d1e55d39a5f507d8994f22bd5e0bfb1e28e2364))

## [0.35.1](https://github.com/chanzuckerberg/happy/compare/cli-v0.35.0...cli-v0.35.1) (2022-10-11)


### Bug Fixes

* increase sleep time after tag creation ([#650](https://github.com/chanzuckerberg/happy/issues/650)) ([4899e90](https://github.com/chanzuckerberg/happy/commit/4899e9016bfc85c77a08cf64ee176b2a61f66069))

## [0.35.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.34.4...cli-v0.35.0) (2022-10-11)


### Features

* Move remaining ECS services to ECS compute backend ([#644](https://github.com/chanzuckerberg/happy/issues/644)) ([3bbc0d9](https://github.com/chanzuckerberg/happy/commit/3bbc0d9b8fa61df2629a5f1ae00dd1262c994a52))

## [0.34.4](https://github.com/chanzuckerberg/happy/compare/cli-v0.34.3...cli-v0.34.4) (2022-10-11)


### Bug Fixes

* allow goreleaser rerun without erroring on tag creation ([#646](https://github.com/chanzuckerberg/happy/issues/646)) ([ce7d3b6](https://github.com/chanzuckerberg/happy/commit/ce7d3b6c1561ea5004c6eabb35cb3327eccc6140))

## [0.34.3](https://github.com/chanzuckerberg/happy/compare/cli-v0.34.2...cli-v0.34.3) (2022-10-10)


### Bug Fixes

* add await to goreleaser tag creation ([#642](https://github.com/chanzuckerberg/happy/issues/642)) ([30feee9](https://github.com/chanzuckerberg/happy/commit/30feee94056a12a45eabdd0f31d64fd5df082afe))

## [0.34.2](https://github.com/chanzuckerberg/happy/compare/cli-v0.34.1...cli-v0.34.2) (2022-10-10)


### Bug Fixes

* create tag for goreleaser ([#640](https://github.com/chanzuckerberg/happy/issues/640)) ([6c5f60e](https://github.com/chanzuckerberg/happy/commit/6c5f60e12fb63cdf1ea61488374dbdf14ac5a0a2))
* Error syncing load balancer: failed to ensure load balancer: could not find any suitable subnets for creating the ELB ([#637](https://github.com/chanzuckerberg/happy/issues/637)) ([dc21f81](https://github.com/chanzuckerberg/happy/commit/dc21f811607bcbf2d2747069766e4f522517873d))

## [0.34.1](https://github.com/chanzuckerberg/happy/compare/cli-v0.34.0...cli-v0.34.1) (2022-10-10)


### Bug Fixes

* use latest goreleaser ([#638](https://github.com/chanzuckerberg/happy/issues/638)) ([0381b7e](https://github.com/chanzuckerberg/happy/commit/0381b7e99e379c52afd1ff1bfdc833266d22f123))

## [0.34.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.33.0...cli-v0.34.0) (2022-10-10)


### Features

* Implement an abstraction around GetEvents to support ECS and K8S ([#628](https://github.com/chanzuckerberg/happy/issues/628)) ([ce97dfd](https://github.com/chanzuckerberg/happy/commit/ce97dfdec95a629cc8917401e018038fe2824ef8))

## [0.33.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.32.0...cli-v0.33.0) (2022-10-07)


### Features

* Implement exec abstraction layer over ECS and Kubernetes in Happy CLI ([#623](https://github.com/chanzuckerberg/happy/issues/623)) ([3a89421](https://github.com/chanzuckerberg/happy/commit/3a89421878b7f4e48ef4dff04c6705ecf0899750))

## [0.32.0](https://github.com/chanzuckerberg/happy/compare/cli-v0.31.4...cli-v0.32.0) (2022-10-07)


### Features

* add shared package ([#620](https://github.com/chanzuckerberg/happy/issues/620)) ([159bd8e](https://github.com/chanzuckerberg/happy/commit/159bd8e372cdf4c2897ca71395c1d65667b0b423))
* Implement logging abstraction layer over ECS and kubernetes in happy cli ([#607](https://github.com/chanzuckerberg/happy/issues/607)) ([af144dc](https://github.com/chanzuckerberg/happy/commit/af144dc4a3ad98ae45f44be66bb2d3847dc1b7f2))

## [0.31.4](https://github.com/chanzuckerberg/happy/compare/v0.31.3...v0.31.4) (2022-10-05)


### Misc

* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#602](https://github.com/chanzuckerberg/happy/issues/602)) ([a881b64](https://github.com/chanzuckerberg/happy/commit/a881b646fab678031a678bb6f9ed79a61150a779))
* bump github.com/aws/aws-sdk-go-v2/service/ecs ([#601](https://github.com/chanzuckerberg/happy/issues/601)) ([6b24d49](https://github.com/chanzuckerberg/happy/commit/6b24d49fffec8a7189af9d07e5a0970e3eea0840))

## [0.31.3](https://github.com/chanzuckerberg/happy/compare/v0.31.2...v0.31.3) (2022-10-04)


### BugFixes

* Fix happy deploy inline documentation ([#599](https://github.com/chanzuckerberg/happy/issues/599)) ([a27e4d9](https://github.com/chanzuckerberg/happy/commit/a27e4d98571965f8a09a713871889e7fa984046d))
* Made cobra descriptions uniform and remove a defunct "test" ([#597](https://github.com/chanzuckerberg/happy/issues/597)) ([0f98bd1](https://github.com/chanzuckerberg/happy/commit/0f98bd1bbd945c6e26ed3e04703da62ed53f65e9))


### Misc

* Update coverage ([#600](https://github.com/chanzuckerberg/happy/issues/600)) ([4145c02](https://github.com/chanzuckerberg/happy/commit/4145c02ff0a226119de1f87c7be25cd432eefebc))

## [0.31.2](https://github.com/chanzuckerberg/happy/compare/v0.31.1...v0.31.2) (2022-10-04)


### Misc

* Add more descriptive error messaging to file walking ([#595](https://github.com/chanzuckerberg/happy/issues/595)) ([ed22e14](https://github.com/chanzuckerberg/happy/commit/ed22e141579a3e2d1f0d15bd7d137946f4136dd3))

## [0.31.1](https://github.com/chanzuckerberg/happy/compare/v0.31.0...v0.31.1) (2022-10-04)


### BugFixes

* (CCIE-707) De-duplicate stacks on force create ([#587](https://github.com/chanzuckerberg/happy/issues/587)) ([1b2cf1e](https://github.com/chanzuckerberg/happy/commit/1b2cf1e7f21640c0185664ee3839e4370fd02c2f))


### Misc

* bump github.com/aws/aws-sdk-go-v2/config from 1.17.7 to 1.17.8 ([#593](https://github.com/chanzuckerberg/happy/issues/593)) ([7f3c32c](https://github.com/chanzuckerberg/happy/commit/7f3c32c5dd49ef7288ca236c9d2471ddbfeff220))
* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#591](https://github.com/chanzuckerberg/happy/issues/591)) ([6c4df0c](https://github.com/chanzuckerberg/happy/commit/6c4df0c37c1d11f65a42d4f445d91cc7b3a4e1d5))
* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#594](https://github.com/chanzuckerberg/happy/issues/594)) ([2c0d027](https://github.com/chanzuckerberg/happy/commit/2c0d027a2b1f6eb8b4290cad6251aeb9fd08e3a2))
* bump github.com/aws/aws-sdk-go-v2/service/secretsmanager ([#590](https://github.com/chanzuckerberg/happy/issues/590)) ([0211967](https://github.com/chanzuckerberg/happy/commit/0211967c807cd2e7ecd7720301eb45468c4fb7b7))
* bump github.com/aws/aws-sdk-go-v2/service/ssm ([#585](https://github.com/chanzuckerberg/happy/issues/585)) ([edd61ca](https://github.com/chanzuckerberg/happy/commit/edd61ca0a5df1e214e9a743b3adfebdd9ada466a))
* Update coverage ([#588](https://github.com/chanzuckerberg/happy/issues/588)) ([36b6a3e](https://github.com/chanzuckerberg/happy/commit/36b6a3eacd6465043cbeb94e4edc2c881107f775))

## [0.31.0](https://github.com/chanzuckerberg/happy/compare/v0.30.0...v0.31.0) (2022-09-28)


### Features

* use namespaced stacklist ([#586](https://github.com/chanzuckerberg/happy/issues/586)) ([09bb2c2](https://github.com/chanzuckerberg/happy/commit/09bb2c2e4dd62ca859c53e0279b9f42c8b76d9be))


### Misc

* add stacklist override tests ([#577](https://github.com/chanzuckerberg/happy/issues/577)) ([feabb47](https://github.com/chanzuckerberg/happy/commit/feabb479afa95f6ddf4bfebe2416fe52f7a0c322))
* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#582](https://github.com/chanzuckerberg/happy/issues/582)) ([82fb88a](https://github.com/chanzuckerberg/happy/commit/82fb88a6adfc922a831c075d5f38a7878bc02291))
* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#584](https://github.com/chanzuckerberg/happy/issues/584)) ([a1fa17f](https://github.com/chanzuckerberg/happy/commit/a1fa17f93b1537d243250d79b84e3a42a6b77f5e))
* bump github.com/aws/aws-sdk-go-v2/service/ssm ([#583](https://github.com/chanzuckerberg/happy/issues/583)) ([7c8de26](https://github.com/chanzuckerberg/happy/commit/7c8de266be1bc864a31f2e8da43074f8050a630a))
* bump k8s.io/api from 0.25.1 to 0.25.2 ([#580](https://github.com/chanzuckerberg/happy/issues/580)) ([b960046](https://github.com/chanzuckerberg/happy/commit/b96004657f7fc2865b50e21cfafa759585ad4921))
* bump k8s.io/apimachinery from 0.25.1 to 0.25.2 ([#581](https://github.com/chanzuckerberg/happy/issues/581)) ([b09ca71](https://github.com/chanzuckerberg/happy/commit/b09ca71eecceb079989f014e2b7a6dea761905fc))
* bump k8s.io/client-go from 0.25.1 to 0.25.2 ([#579](https://github.com/chanzuckerberg/happy/issues/579)) ([b984d45](https://github.com/chanzuckerberg/happy/commit/b984d45b36bf00d50e393c827fd48af7ebab98eb))

## [0.30.0](https://github.com/chanzuckerberg/happy/compare/v0.29.2...v0.30.0) (2022-09-21)


### Features

* Create unit tests for kubernetes compute ([#554](https://github.com/chanzuckerberg/happy/issues/554)) ([5eaa25b](https://github.com/chanzuckerberg/happy/commit/5eaa25bf52fe95257de407fa72ab656b84af2bad))
* Implement stacklist retrieval abstraction ([#564](https://github.com/chanzuckerberg/happy/issues/564)) ([f025175](https://github.com/chanzuckerberg/happy/commit/f025175357146c87677783e48340590acc35a63a))
* Improved code coverage ([#561](https://github.com/chanzuckerberg/happy/issues/561)) ([46cdbb2](https://github.com/chanzuckerberg/happy/commit/46cdbb24ac42d501a114fdbe89412356c90e711e))
* Localstack support ([#448](https://github.com/chanzuckerberg/happy/issues/448)) ([c4b219b](https://github.com/chanzuckerberg/happy/commit/c4b219b07f547de103420d8b2098ff3760eb029f))
* make stacklist path configurable ([#576](https://github.com/chanzuckerberg/happy/issues/576)) ([33dcfa3](https://github.com/chanzuckerberg/happy/commit/33dcfa3a815b30fe87bde38b9bd0b215c9cd5e60))
* Modify .happy config structure to accommodate kubernetes config ([#540](https://github.com/chanzuckerberg/happy/issues/540)) ([5ec7023](https://github.com/chanzuckerberg/happy/commit/5ec702391d45da637bd9076b14ec7e894918f231))
* Store integration secret in kuberentes secret (self containment) ([#553](https://github.com/chanzuckerberg/happy/issues/553)) ([a6d190a](https://github.com/chanzuckerberg/happy/commit/a6d190a11ce13c67ecd0e8475c1f5849ea274042))


### Misc

* bump github.com/AlecAivazis/survey/v2 from 2.3.5 to 2.3.6 ([#539](https://github.com/chanzuckerberg/happy/issues/539)) ([048a353](https://github.com/chanzuckerberg/happy/commit/048a353962bb0c66c7d8408e356f939acf653b45))
* bump github.com/aws/aws-sdk-go-v2/config from 1.17.5 to 1.17.6 ([#547](https://github.com/chanzuckerberg/happy/issues/547)) ([2b5b464](https://github.com/chanzuckerberg/happy/commit/2b5b464b9c74be6b2db67bb13c85a4d0f3d9948e))
* bump github.com/aws/aws-sdk-go-v2/config from 1.17.6 to 1.17.7 ([#572](https://github.com/chanzuckerberg/happy/issues/572)) ([794d725](https://github.com/chanzuckerberg/happy/commit/794d72585cc82b73050bc172f27b0798f95363a3))
* bump github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs ([#545](https://github.com/chanzuckerberg/happy/issues/545)) ([d2afe43](https://github.com/chanzuckerberg/happy/commit/d2afe4333175ed6f37b81252df4247bf170f4d19))
* bump github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs ([#566](https://github.com/chanzuckerberg/happy/issues/566)) ([60f725c](https://github.com/chanzuckerberg/happy/commit/60f725cebbdcac399837dbd7c916ebbf576c63de))
* bump github.com/aws/aws-sdk-go-v2/service/dynamodb ([#550](https://github.com/chanzuckerberg/happy/issues/550)) ([0732c69](https://github.com/chanzuckerberg/happy/commit/0732c69e03efbdea3999f0e4f15eb6099286b7f3))
* bump github.com/aws/aws-sdk-go-v2/service/dynamodb ([#556](https://github.com/chanzuckerberg/happy/issues/556)) ([d6aee58](https://github.com/chanzuckerberg/happy/commit/d6aee587e476b9750bb2a2bbef3abd40e5e2f130))
* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#542](https://github.com/chanzuckerberg/happy/issues/542)) ([5410558](https://github.com/chanzuckerberg/happy/commit/5410558f6ff2a663ebd76f55737c6922d50ab287))
* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#555](https://github.com/chanzuckerberg/happy/issues/555)) ([ec2ea53](https://github.com/chanzuckerberg/happy/commit/ec2ea537e9e6ec1842f2694b2827aada387310b3))
* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#562](https://github.com/chanzuckerberg/happy/issues/562)) ([1cffdc3](https://github.com/chanzuckerberg/happy/commit/1cffdc341bab3698bd51a4fa903c969380b4e97a))
* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#569](https://github.com/chanzuckerberg/happy/issues/569)) ([a95b0c5](https://github.com/chanzuckerberg/happy/commit/a95b0c521a4339dd10ca23b879a82bab07d14ac2))
* bump github.com/aws/aws-sdk-go-v2/service/ecr ([#543](https://github.com/chanzuckerberg/happy/issues/543)) ([2865a90](https://github.com/chanzuckerberg/happy/commit/2865a90e8b9253df48f6b1eac4a10a57473c2160))
* bump github.com/aws/aws-sdk-go-v2/service/ecr ([#568](https://github.com/chanzuckerberg/happy/issues/568)) ([d79c0ba](https://github.com/chanzuckerberg/happy/commit/d79c0ba17f0100e7f45842b960dd60e064cc878b))
* bump github.com/aws/aws-sdk-go-v2/service/ecs ([#549](https://github.com/chanzuckerberg/happy/issues/549)) ([fecd148](https://github.com/chanzuckerberg/happy/commit/fecd1482c41245db4794c687c20b360a91eb5c64))
* bump github.com/aws/aws-sdk-go-v2/service/ecs ([#559](https://github.com/chanzuckerberg/happy/issues/559)) ([e9c2d4f](https://github.com/chanzuckerberg/happy/commit/e9c2d4f090b1215c852d2c1c064034bc4a6c2eaf))
* bump github.com/aws/aws-sdk-go-v2/service/ecs ([#570](https://github.com/chanzuckerberg/happy/issues/570)) ([1f2c0d8](https://github.com/chanzuckerberg/happy/commit/1f2c0d8bb12b256c84bb9785df664ca84817d9bc))
* bump github.com/aws/aws-sdk-go-v2/service/eks ([#574](https://github.com/chanzuckerberg/happy/issues/574)) ([5848e22](https://github.com/chanzuckerberg/happy/commit/5848e22869773a842429996a91cb667a2dcbc237))
* bump github.com/aws/aws-sdk-go-v2/service/secretsmanager ([#552](https://github.com/chanzuckerberg/happy/issues/552)) ([1a17071](https://github.com/chanzuckerberg/happy/commit/1a170714c05b41aef3f6b1f8c0d4bb53cf1aa30e))
* bump github.com/aws/aws-sdk-go-v2/service/secretsmanager ([#573](https://github.com/chanzuckerberg/happy/issues/573)) ([fae5889](https://github.com/chanzuckerberg/happy/commit/fae58891acb7256ef30436c8ee95cb6c47fab13e))
* bump github.com/aws/aws-sdk-go-v2/service/ssm ([#571](https://github.com/chanzuckerberg/happy/issues/571)) ([b250c0d](https://github.com/chanzuckerberg/happy/commit/b250c0d72eb87d1b2494a46396a698990a8314e9))
* bump github.com/aws/aws-sdk-go-v2/service/sts ([#541](https://github.com/chanzuckerberg/happy/issues/541)) ([aacf22b](https://github.com/chanzuckerberg/happy/commit/aacf22b6d83322907f8bcd4b508afeff7b47e3bd))
* bump github.com/docker/docker ([#538](https://github.com/chanzuckerberg/happy/issues/538)) ([79a41e8](https://github.com/chanzuckerberg/happy/commit/79a41e8146d07e5356dc7ab056ee992d882a454d))
* bump github.com/gruntwork-io/terratest from 0.40.21 to 0.40.22 ([#536](https://github.com/chanzuckerberg/happy/issues/536)) ([a8b6fb3](https://github.com/chanzuckerberg/happy/commit/a8b6fb3aa86e094e63637264fc52fd39267deeef))
* bump github.com/hashicorp/go-tfe from 1.9.0 to 1.10.0 ([#548](https://github.com/chanzuckerberg/happy/issues/548)) ([07bb6a5](https://github.com/chanzuckerberg/happy/commit/07bb6a56a911e09260f16566897203f476e0f217))
* bump k8s.io/api from 0.25.0 to 0.25.1 ([#560](https://github.com/chanzuckerberg/happy/issues/560)) ([0e1dae8](https://github.com/chanzuckerberg/happy/commit/0e1dae8febde9b0347a452f37620e29ca6bf53af))
* bump k8s.io/apimachinery from 0.25.0 to 0.25.1 ([#557](https://github.com/chanzuckerberg/happy/issues/557)) ([5b4d15a](https://github.com/chanzuckerberg/happy/commit/5b4d15ae71d9ffc4007e1dcdf76b82a75812e331))
* bump k8s.io/client-go from 0.25.0 to 0.25.1 ([#558](https://github.com/chanzuckerberg/happy/issues/558)) ([0845201](https://github.com/chanzuckerberg/happy/commit/0845201505382c1907d0e3ac3e51cf77ce576ef9))
* Update coverage ([#565](https://github.com/chanzuckerberg/happy/issues/565)) ([6e0e2fd](https://github.com/chanzuckerberg/happy/commit/6e0e2fd74b5f53f60def42b7a644d09bcf8ce707))

## [0.29.2](https://github.com/chanzuckerberg/happy/compare/v0.29.1...v0.29.2) (2022-09-06)


### Misc

* bump github.com/aws/aws-sdk-go-v2/config from 1.17.4 to 1.17.5 ([#523](https://github.com/chanzuckerberg/happy/issues/523)) ([9230b74](https://github.com/chanzuckerberg/happy/commit/9230b74123526f2c35d092fe786698b5f46bc757))
* bump github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs ([#526](https://github.com/chanzuckerberg/happy/issues/526)) ([96d2c55](https://github.com/chanzuckerberg/happy/commit/96d2c5583b09c287baf64986417f48a43643f6e7))
* bump github.com/aws/aws-sdk-go-v2/service/dynamodb ([#527](https://github.com/chanzuckerberg/happy/issues/527)) ([264bea0](https://github.com/chanzuckerberg/happy/commit/264bea06fdfbc457277bb6ee9b579937221d2e68))
* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#535](https://github.com/chanzuckerberg/happy/issues/535)) ([d104ada](https://github.com/chanzuckerberg/happy/commit/d104ada60865157b0ff022f83d03d121f80bbc35))
* bump github.com/aws/aws-sdk-go-v2/service/ecr ([#531](https://github.com/chanzuckerberg/happy/issues/531)) ([d453a1f](https://github.com/chanzuckerberg/happy/commit/d453a1f6f4b54a9b0b432196061a7cdd6ca2d0b2))
* bump github.com/aws/aws-sdk-go-v2/service/ecs ([#530](https://github.com/chanzuckerberg/happy/issues/530)) ([e64cb16](https://github.com/chanzuckerberg/happy/commit/e64cb165b8b4c91c5788d116845a4b2d57fff3f1))
* bump github.com/aws/aws-sdk-go-v2/service/secretsmanager ([#532](https://github.com/chanzuckerberg/happy/issues/532)) ([a20bdd5](https://github.com/chanzuckerberg/happy/commit/a20bdd5e698881677df75643f15e5589998c0ff0))
* bump github.com/aws/aws-sdk-go-v2/service/ssm ([#529](https://github.com/chanzuckerberg/happy/issues/529)) ([3dc6b1a](https://github.com/chanzuckerberg/happy/commit/3dc6b1ac662eba0900c9040e7ba91d4c44fa1483))
* bump github.com/aws/aws-sdk-go-v2/service/sts ([#534](https://github.com/chanzuckerberg/happy/issues/534)) ([10987cb](https://github.com/chanzuckerberg/happy/commit/10987cbc06d52dd846a6e45c8c10bbe6559faaad))
* bump github.com/gruntwork-io/terratest from 0.40.20 to 0.40.21 ([#524](https://github.com/chanzuckerberg/happy/issues/524)) ([a00b12f](https://github.com/chanzuckerberg/happy/commit/a00b12f7ef26cacb36f2fc6eb9dc7bac3e4e4b0a))

## [0.29.1](https://github.com/chanzuckerberg/happy/compare/v0.29.0...v0.29.1) (2022-09-01)


### Misc

* bump github.com/aws/aws-sdk-go-v2 from 1.16.11 to 1.16.12 ([#496](https://github.com/chanzuckerberg/happy/issues/496)) ([da34fd8](https://github.com/chanzuckerberg/happy/commit/da34fd8131a61fa22239387fc348e7c29aedaec2))
* bump github.com/aws/aws-sdk-go-v2/config from 1.17.1 to 1.17.2 ([#492](https://github.com/chanzuckerberg/happy/issues/492)) ([41b990a](https://github.com/chanzuckerberg/happy/commit/41b990a6f4555e0d5a5cfc0b7c751d4a38b49872))
* bump github.com/aws/aws-sdk-go-v2/config from 1.17.2 to 1.17.3 ([#509](https://github.com/chanzuckerberg/happy/issues/509)) ([ab98752](https://github.com/chanzuckerberg/happy/commit/ab98752b1e4f1f569bb19b5f9bfee9a8ecf4435c))
* bump github.com/aws/aws-sdk-go-v2/config from 1.17.3 to 1.17.4 ([#522](https://github.com/chanzuckerberg/happy/issues/522)) ([f713b95](https://github.com/chanzuckerberg/happy/commit/f713b95069c27c831240393d10b8f021275499e8))
* bump github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs ([#499](https://github.com/chanzuckerberg/happy/issues/499)) ([3d94b27](https://github.com/chanzuckerberg/happy/commit/3d94b27d76c8212d447166fbfdb31ec8e0493a19))
* bump github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs ([#502](https://github.com/chanzuckerberg/happy/issues/502)) ([beed629](https://github.com/chanzuckerberg/happy/commit/beed629dadf67e17e9f63876e89a2e1b0b56eb69))
* bump github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs ([#519](https://github.com/chanzuckerberg/happy/issues/519)) ([97f8a42](https://github.com/chanzuckerberg/happy/commit/97f8a420abe9020f1c303134b97f15817fcbee7c))
* bump github.com/aws/aws-sdk-go-v2/service/dynamodb ([#497](https://github.com/chanzuckerberg/happy/issues/497)) ([f1b513e](https://github.com/chanzuckerberg/happy/commit/f1b513ef2127f1caf8e95f1bd02627b311c4b285))
* bump github.com/aws/aws-sdk-go-v2/service/dynamodb ([#505](https://github.com/chanzuckerberg/happy/issues/505)) ([2ee6288](https://github.com/chanzuckerberg/happy/commit/2ee6288a300a74d6dad184f66ecad92fdad0e6cb))
* bump github.com/aws/aws-sdk-go-v2/service/dynamodb ([#518](https://github.com/chanzuckerberg/happy/issues/518)) ([aeb547e](https://github.com/chanzuckerberg/happy/commit/aeb547e33415a55972caa8178f4352d7b3f2d469))
* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#494](https://github.com/chanzuckerberg/happy/issues/494)) ([f166aba](https://github.com/chanzuckerberg/happy/commit/f166abad24416cd9f80fbdb71ca06dc70992c21e))
* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#503](https://github.com/chanzuckerberg/happy/issues/503)) ([e8aacd7](https://github.com/chanzuckerberg/happy/commit/e8aacd77a11f86906ed7d3589d18aa17f32a4fb6))
* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#515](https://github.com/chanzuckerberg/happy/issues/515)) ([40bae2a](https://github.com/chanzuckerberg/happy/commit/40bae2a72bd7e6bb0d3deed47125ca8c026d8834))
* bump github.com/aws/aws-sdk-go-v2/service/ecr ([#493](https://github.com/chanzuckerberg/happy/issues/493)) ([0835c1d](https://github.com/chanzuckerberg/happy/commit/0835c1d1c8a2121f06c19c6e02000536990af0cf))
* bump github.com/aws/aws-sdk-go-v2/service/ecr ([#506](https://github.com/chanzuckerberg/happy/issues/506)) ([04f9036](https://github.com/chanzuckerberg/happy/commit/04f9036c869b2db43b6196069a71455a09db0f8e))
* bump github.com/aws/aws-sdk-go-v2/service/ecr ([#520](https://github.com/chanzuckerberg/happy/issues/520)) ([2b6f620](https://github.com/chanzuckerberg/happy/commit/2b6f620ee4567e05f4e4c998e2b3d217f9acc948))
* bump github.com/aws/aws-sdk-go-v2/service/ecs ([#491](https://github.com/chanzuckerberg/happy/issues/491)) ([f476ab7](https://github.com/chanzuckerberg/happy/commit/f476ab7aad625adf14ef15311fec635384c87998))
* bump github.com/aws/aws-sdk-go-v2/service/ecs ([#510](https://github.com/chanzuckerberg/happy/issues/510)) ([8c58f24](https://github.com/chanzuckerberg/happy/commit/8c58f244bc69e6061c920198d87ff0c374b89acf))
* bump github.com/aws/aws-sdk-go-v2/service/ecs ([#511](https://github.com/chanzuckerberg/happy/issues/511)) ([0162a27](https://github.com/chanzuckerberg/happy/commit/0162a2773933af369fe2178b2b4e1cd2f6de7ffe))
* bump github.com/aws/aws-sdk-go-v2/service/secretsmanager ([#489](https://github.com/chanzuckerberg/happy/issues/489)) ([4909c8e](https://github.com/chanzuckerberg/happy/commit/4909c8e10e8bb5060c2121af46b417e5cef51e76))
* bump github.com/aws/aws-sdk-go-v2/service/secretsmanager ([#504](https://github.com/chanzuckerberg/happy/issues/504)) ([13a012a](https://github.com/chanzuckerberg/happy/commit/13a012aedafe7e9cdeb4c758eb56a499cf458651))
* bump github.com/aws/aws-sdk-go-v2/service/secretsmanager ([#512](https://github.com/chanzuckerberg/happy/issues/512)) ([cca9735](https://github.com/chanzuckerberg/happy/commit/cca9735b48bc6cef733abef54c2d2aeb8010882d))
* bump github.com/aws/aws-sdk-go-v2/service/ssm ([#490](https://github.com/chanzuckerberg/happy/issues/490)) ([a9808b2](https://github.com/chanzuckerberg/happy/commit/a9808b2977746df41d8ce27dafc7d6f37613bafe))
* bump github.com/aws/aws-sdk-go-v2/service/ssm ([#507](https://github.com/chanzuckerberg/happy/issues/507)) ([67bc59c](https://github.com/chanzuckerberg/happy/commit/67bc59ca7a64415a21bac85de1ffeb2f7a2254a5))
* bump github.com/aws/aws-sdk-go-v2/service/ssm ([#513](https://github.com/chanzuckerberg/happy/issues/513)) ([9714f2c](https://github.com/chanzuckerberg/happy/commit/9714f2cc22f8fca9152fda9e6fa4907933809903))
* bump github.com/aws/aws-sdk-go-v2/service/sts ([#500](https://github.com/chanzuckerberg/happy/issues/500)) ([243cb67](https://github.com/chanzuckerberg/happy/commit/243cb6797bb7a39896c3a74dcc2803d76a96f817))
* bump github.com/aws/aws-sdk-go-v2/service/sts ([#521](https://github.com/chanzuckerberg/happy/issues/521)) ([f2f4fab](https://github.com/chanzuckerberg/happy/commit/f2f4fab544a4964ae5449d9a26fe73a549c3380f))
* bump github.com/aws/smithy-go from 1.12.1 to 1.13.0 ([#488](https://github.com/chanzuckerberg/happy/issues/488)) ([3b77f2b](https://github.com/chanzuckerberg/happy/commit/3b77f2bddfd525ed74f6dbff03cd59984202d3c4))
* bump github.com/docker/go-units from 0.4.0 to 0.5.0 ([#516](https://github.com/chanzuckerberg/happy/issues/516)) ([5390a32](https://github.com/chanzuckerberg/happy/commit/5390a325809cad652dfd899166e388259fe35756))
* bump github.com/hashicorp/go-tfe from 1.8.0 to 1.9.0 ([#508](https://github.com/chanzuckerberg/happy/issues/508)) ([3e91cda](https://github.com/chanzuckerberg/happy/commit/3e91cda3905988b4acb7f5d2054f384b12d2718a))

## [0.29.0](https://github.com/chanzuckerberg/happy/compare/v0.28.2...v0.29.0) (2022-08-29)


### Features

* allow people to output logs to a file ([#481](https://github.com/chanzuckerberg/happy/issues/481)) ([f11c4c3](https://github.com/chanzuckerberg/happy/commit/f11c4c370df92e7057a447fb2cc00219482bb3b7))


### Misc

* bump github.com/aws/aws-sdk-go-v2/service/dynamodb ([#485](https://github.com/chanzuckerberg/happy/issues/485)) ([75a9716](https://github.com/chanzuckerberg/happy/commit/75a9716a979f8352c086537839d9fc0740524000))
* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#482](https://github.com/chanzuckerberg/happy/issues/482)) ([21b3d70](https://github.com/chanzuckerberg/happy/commit/21b3d70e526d5ce8c7fd177a5c9414be065205b3))
* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#487](https://github.com/chanzuckerberg/happy/issues/487)) ([1cd1ae2](https://github.com/chanzuckerberg/happy/commit/1cd1ae29447d66447c1aecc5e279f97f86854158))
* bump github.com/aws/aws-sdk-go-v2/service/secretsmanager ([#484](https://github.com/chanzuckerberg/happy/issues/484)) ([53ab690](https://github.com/chanzuckerberg/happy/commit/53ab690dd24e153bf88f5b0231fa0a33f0eff58d))
* bump github.com/hashicorp/go-tfe from 1.7.0 to 1.8.0 ([#483](https://github.com/chanzuckerberg/happy/issues/483)) ([81c7836](https://github.com/chanzuckerberg/happy/commit/81c7836d4a771fe4bf03f5c4334706d61186d0d0))

## [0.28.2](https://github.com/chanzuckerberg/happy/compare/v0.28.1...v0.28.2) (2022-08-18)


### Misc

* bump github.com/aws/aws-sdk-go-v2/config from 1.16.1 to 1.17.0 ([#474](https://github.com/chanzuckerberg/happy/issues/474)) ([e72cbc6](https://github.com/chanzuckerberg/happy/commit/e72cbc6b562537c3eaa246644a2343ed1918b47a))
* bump github.com/aws/aws-sdk-go-v2/config from 1.17.0 to 1.17.1 ([#476](https://github.com/chanzuckerberg/happy/issues/476)) ([ed09342](https://github.com/chanzuckerberg/happy/commit/ed09342ab2fe38adbdf5a6eff521386f6e4ae9a2))
* bump github.com/gruntwork-io/terratest from 0.40.19 to 0.40.20 ([#477](https://github.com/chanzuckerberg/happy/issues/477)) ([802e105](https://github.com/chanzuckerberg/happy/commit/802e105eac74cf5d913b90ed653a99b5094f617e))
* bump github.com/hashicorp/go-tfe from 1.7.0 to 1.8.0 ([#480](https://github.com/chanzuckerberg/happy/issues/480)) ([f3c1444](https://github.com/chanzuckerberg/happy/commit/f3c1444f738a201ec13ee081e476c3cb13d5da14))


### BugFixes

* refactor and combine log streams ([#479](https://github.com/chanzuckerberg/happy/issues/479)) ([16d2ce8](https://github.com/chanzuckerberg/happy/commit/16d2ce895c2e88f87c2a3bd636ee37ba25b4ed6c))

## [0.28.1](https://github.com/chanzuckerberg/happy/compare/v0.28.0...v0.28.1) (2022-08-12)


### Misc

* bump github.com/aws/aws-sdk-go-v2 from 1.16.9 to 1.16.10 ([#459](https://github.com/chanzuckerberg/happy/issues/459)) ([73b7089](https://github.com/chanzuckerberg/happy/commit/73b70897bb84e5987597b1bcd830a000aa4b8daa))
* bump github.com/aws/aws-sdk-go-v2/config from 1.15.15 to 1.15.16 ([#439](https://github.com/chanzuckerberg/happy/issues/439)) ([033ab22](https://github.com/chanzuckerberg/happy/commit/033ab22d80b64cfdee748db10dd720a6a7f0d75a))
* bump github.com/aws/aws-sdk-go-v2/config from 1.15.16 to 1.15.17 ([#451](https://github.com/chanzuckerberg/happy/issues/451)) ([e992bb4](https://github.com/chanzuckerberg/happy/commit/e992bb45aa1dfebf4db5685004af47d235d3ca1b))
* bump github.com/aws/aws-sdk-go-v2/config from 1.15.17 to 1.16.0 ([#460](https://github.com/chanzuckerberg/happy/issues/460)) ([b1d5102](https://github.com/chanzuckerberg/happy/commit/b1d51027468d3ebb42dac2dfd1bf00ce151ba493))
* bump github.com/aws/aws-sdk-go-v2/config from 1.16.0 to 1.16.1 ([#472](https://github.com/chanzuckerberg/happy/issues/472)) ([0293020](https://github.com/chanzuckerberg/happy/commit/029302074e2fd7af6df5d91defe63d663ab7f19f))
* bump github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs ([#444](https://github.com/chanzuckerberg/happy/issues/444)) ([e695492](https://github.com/chanzuckerberg/happy/commit/e6954929c7cd17b4f0b10b5744a9e9847ef680ea))
* bump github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs ([#457](https://github.com/chanzuckerberg/happy/issues/457)) ([d48dab7](https://github.com/chanzuckerberg/happy/commit/d48dab7ba1ab7d502ed2f0b891fdbfec8b30ebe8))
* bump github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs ([#464](https://github.com/chanzuckerberg/happy/issues/464)) ([07a5834](https://github.com/chanzuckerberg/happy/commit/07a5834876cf245b93f8b9815c04e45897a4b0f4))
* bump github.com/aws/aws-sdk-go-v2/service/dynamodb ([#438](https://github.com/chanzuckerberg/happy/issues/438)) ([d75d112](https://github.com/chanzuckerberg/happy/commit/d75d112c8dff36a49501d3406ad7297644bdba87))
* bump github.com/aws/aws-sdk-go-v2/service/dynamodb ([#454](https://github.com/chanzuckerberg/happy/issues/454)) ([066aefe](https://github.com/chanzuckerberg/happy/commit/066aefea0e9ad1414eb1b4c3b76778513042ae3f))
* bump github.com/aws/aws-sdk-go-v2/service/dynamodb ([#465](https://github.com/chanzuckerberg/happy/issues/465)) ([b550330](https://github.com/chanzuckerberg/happy/commit/b55033045158374717d5f9a90c9c64a6e579f32b))
* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#441](https://github.com/chanzuckerberg/happy/issues/441)) ([48f3151](https://github.com/chanzuckerberg/happy/commit/48f31514c86faf161ab7c079e5e25dae14c7daa6))
* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#456](https://github.com/chanzuckerberg/happy/issues/456)) ([cac023e](https://github.com/chanzuckerberg/happy/commit/cac023ef234c7b26d067746fadde648b16040ac6))
* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#463](https://github.com/chanzuckerberg/happy/issues/463)) ([4067daf](https://github.com/chanzuckerberg/happy/commit/4067daf97dd2a64a934f250c637af292a6f132d5))
* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#467](https://github.com/chanzuckerberg/happy/issues/467)) ([f0d3ceb](https://github.com/chanzuckerberg/happy/commit/f0d3cebba235cf784a268ee37f7a6b3bf5ef65d2))
* bump github.com/aws/aws-sdk-go-v2/service/ecr ([#440](https://github.com/chanzuckerberg/happy/issues/440)) ([7d0c5a2](https://github.com/chanzuckerberg/happy/commit/7d0c5a2296c37524af2f20dedd00ba336ab0ec22))
* bump github.com/aws/aws-sdk-go-v2/service/ecr ([#458](https://github.com/chanzuckerberg/happy/issues/458)) ([875d1af](https://github.com/chanzuckerberg/happy/commit/875d1aff6631f7542b63ac9837b6a3a4d42f67a9))
* bump github.com/aws/aws-sdk-go-v2/service/ecr ([#470](https://github.com/chanzuckerberg/happy/issues/470)) ([dd84884](https://github.com/chanzuckerberg/happy/commit/dd848843beef003b7c2bba171a37e7becff1d042))
* bump github.com/aws/aws-sdk-go-v2/service/ecs ([#447](https://github.com/chanzuckerberg/happy/issues/447)) ([65c576a](https://github.com/chanzuckerberg/happy/commit/65c576a0d87afa08fd3248ef9f742208e91fa62d))
* bump github.com/aws/aws-sdk-go-v2/service/ecs ([#452](https://github.com/chanzuckerberg/happy/issues/452)) ([e87e0e4](https://github.com/chanzuckerberg/happy/commit/e87e0e4bbe8c242216907facda718a8dfda34f24))
* bump github.com/aws/aws-sdk-go-v2/service/ecs ([#466](https://github.com/chanzuckerberg/happy/issues/466)) ([544f43b](https://github.com/chanzuckerberg/happy/commit/544f43ba1fca797af08dd43d8ffce539d014cc9a))
* bump github.com/aws/aws-sdk-go-v2/service/secretsmanager ([#445](https://github.com/chanzuckerberg/happy/issues/445)) ([0b8f9ad](https://github.com/chanzuckerberg/happy/commit/0b8f9ad146fa39ad1815429e64bad4f3a5093d90))
* bump github.com/aws/aws-sdk-go-v2/service/secretsmanager ([#453](https://github.com/chanzuckerberg/happy/issues/453)) ([ca93aaa](https://github.com/chanzuckerberg/happy/commit/ca93aaa74818d7a617887815d36337ec16b0561d))
* bump github.com/aws/aws-sdk-go-v2/service/secretsmanager ([#468](https://github.com/chanzuckerberg/happy/issues/468)) ([f411b71](https://github.com/chanzuckerberg/happy/commit/f411b713a839a6f26a99b44a690e8e4213c558ad))
* bump github.com/aws/aws-sdk-go-v2/service/ssm ([#442](https://github.com/chanzuckerberg/happy/issues/442)) ([52c2bd2](https://github.com/chanzuckerberg/happy/commit/52c2bd282b3351a28b488ce977e3e5c0cd6ceba1))
* bump github.com/aws/aws-sdk-go-v2/service/ssm ([#461](https://github.com/chanzuckerberg/happy/issues/461)) ([d01b633](https://github.com/chanzuckerberg/happy/commit/d01b6335de5a212c9ad7a3b87daa32e4e30df75a))
* bump github.com/aws/aws-sdk-go-v2/service/ssm ([#471](https://github.com/chanzuckerberg/happy/issues/471)) ([5ee695b](https://github.com/chanzuckerberg/happy/commit/5ee695b1e80bd6adc78a3a4299e9c7b2c1416e43))
* bump github.com/aws/aws-sdk-go-v2/service/sts ([#443](https://github.com/chanzuckerberg/happy/issues/443)) ([44e8aa5](https://github.com/chanzuckerberg/happy/commit/44e8aa5ed24f4def4146ee19a1dc0cf4eb9f38c4))
* bump github.com/aws/aws-sdk-go-v2/service/sts ([#455](https://github.com/chanzuckerberg/happy/issues/455)) ([de82079](https://github.com/chanzuckerberg/happy/commit/de8207905940ca03624d8dfbeb3215ad50b4b2e6))
* bump github.com/aws/aws-sdk-go-v2/service/sts ([#473](https://github.com/chanzuckerberg/happy/issues/473)) ([2c06299](https://github.com/chanzuckerberg/happy/commit/2c0629943838c36918dcc9d6bab0bc4fcd4e3909))
* bump github.com/gruntwork-io/terratest from 0.40.18 to 0.40.19 ([#436](https://github.com/chanzuckerberg/happy/issues/436)) ([9fc1902](https://github.com/chanzuckerberg/happy/commit/9fc19027f7d7de91da50eee7ef24644a2e4ad6c9))
* bump github.com/hashicorp/go-tfe from 1.6.0 to 1.7.0 ([#462](https://github.com/chanzuckerberg/happy/issues/462)) ([d761d67](https://github.com/chanzuckerberg/happy/commit/d761d67bb19d5a24fd4a144193f689ae6ce96684))
* replicate control flow changes from TFE run info in profiler ([#399](https://github.com/chanzuckerberg/happy/issues/399)) ([3c061c7](https://github.com/chanzuckerberg/happy/commit/3c061c7997f7fa9906dc31718abc8e19dcea9e60))
* Update coverage ([#449](https://github.com/chanzuckerberg/happy/issues/449)) ([1d04442](https://github.com/chanzuckerberg/happy/commit/1d04442c78b28364ef2c5d24529ba6d41fbeedb4))

## [0.28.0](https://github.com/chanzuckerberg/happy/compare/v0.27.2...v0.28.0) (2022-08-03)


### Features

* (CCIE-176) Add link for TFE run to CLI output ([#374](https://github.com/chanzuckerberg/happy/issues/374)) ([4f417f7](https://github.com/chanzuckerberg/happy/commit/4f417f7f844fd05ca2b2db7afb03326319998476))


### BugFixes

* (CCIE-452) Happy migrate in single-cell-data-portal errors out waiting for the log stream ([#435](https://github.com/chanzuckerberg/happy/issues/435)) ([4bb51c7](https://github.com/chanzuckerberg/happy/commit/4bb51c7736c45646cb4b271685f0f20de6d6871d))

## [0.27.2](https://github.com/chanzuckerberg/happy/compare/v0.27.1...v0.27.2) (2022-08-02)


### Misc

* bump github.com/aws/aws-sdk-go-v2/config from 1.15.14 to 1.15.15 ([#428](https://github.com/chanzuckerberg/happy/issues/428)) ([3632da3](https://github.com/chanzuckerberg/happy/commit/3632da376887fa25aca5d0ffe5e852e53e773731))
* bump github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs ([#429](https://github.com/chanzuckerberg/happy/issues/429)) ([1bb91c1](https://github.com/chanzuckerberg/happy/commit/1bb91c119bf96b44d77f4bdc83cd9dd5fa352f76))
* bump github.com/aws/aws-sdk-go-v2/service/dynamodb ([#433](https://github.com/chanzuckerberg/happy/issues/433)) ([404194d](https://github.com/chanzuckerberg/happy/commit/404194dedd9090a5a0a239df2a4cb678203b3318))
* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#422](https://github.com/chanzuckerberg/happy/issues/422)) ([3689541](https://github.com/chanzuckerberg/happy/commit/3689541c3c121617cd4306d99838949ed9282e6c))
* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#423](https://github.com/chanzuckerberg/happy/issues/423)) ([62d70cc](https://github.com/chanzuckerberg/happy/commit/62d70cca40d5cbdee141bf7d2a94c7dfb8160327))
* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#432](https://github.com/chanzuckerberg/happy/issues/432)) ([5c49d43](https://github.com/chanzuckerberg/happy/commit/5c49d43c9a6feaa31f56451f4c5eb020e8ae927d))
* bump github.com/aws/aws-sdk-go-v2/service/ecr ([#431](https://github.com/chanzuckerberg/happy/issues/431)) ([a483399](https://github.com/chanzuckerberg/happy/commit/a483399d0de21174ea1955f1aacb38036413c1db))
* bump github.com/aws/aws-sdk-go-v2/service/ecs ([#430](https://github.com/chanzuckerberg/happy/issues/430)) ([8b9579b](https://github.com/chanzuckerberg/happy/commit/8b9579b30ae0583c22cbcb48442eed2a74f1f123))
* bump github.com/aws/aws-sdk-go-v2/service/secretsmanager ([#426](https://github.com/chanzuckerberg/happy/issues/426)) ([7f38638](https://github.com/chanzuckerberg/happy/commit/7f3863844c18c4538eb4c60303612c9084537a1f))
* bump github.com/aws/aws-sdk-go-v2/service/ssm ([#420](https://github.com/chanzuckerberg/happy/issues/420)) ([6fc0d50](https://github.com/chanzuckerberg/happy/commit/6fc0d50178a4040604bf21006c8b2818e0c38fdb))
* bump github.com/aws/aws-sdk-go-v2/service/ssm ([#425](https://github.com/chanzuckerberg/happy/issues/425)) ([8851f04](https://github.com/chanzuckerberg/happy/commit/8851f04c517841b2355384548177303258448be4))
* bump github.com/aws/aws-sdk-go-v2/service/sts ([#427](https://github.com/chanzuckerberg/happy/issues/427)) ([dd8d5da](https://github.com/chanzuckerberg/happy/commit/dd8d5dac36c932aebc6a3246e5b56add9b880ad1))

## [0.27.1](https://github.com/chanzuckerberg/happy/compare/v0.27.0...v0.27.1) (2022-07-27)


### BugFixes

* (CCIE-443) Enforce latest log streams fetched by logs command ([#418](https://github.com/chanzuckerberg/happy/issues/418)) ([754f58d](https://github.com/chanzuckerberg/happy/commit/754f58d1e3826b88174841edbe04cd98373ead0d))

## [0.27.0](https://github.com/chanzuckerberg/happy/compare/v0.26.1...v0.27.0) (2022-07-27)


### Features

* Replace tablewriter with tableprinter and use header struct annotations ([#410](https://github.com/chanzuckerberg/happy/issues/410)) ([cd131ad](https://github.com/chanzuckerberg/happy/commit/cd131ad377e22d5a8b0f324b55214a2d0456186f))


### Misc

* bump github.com/hashicorp/go-tfe from 1.5.0 to 1.6.0 ([#415](https://github.com/chanzuckerberg/happy/issues/415)) ([514d18e](https://github.com/chanzuckerberg/happy/commit/514d18eb4570817cd1395b3846f45d2c7e4464ac))
* Update coverage ([#416](https://github.com/chanzuckerberg/happy/issues/416)) ([e04b5ea](https://github.com/chanzuckerberg/happy/commit/e04b5eab155b082b67e46cf01f2a68f255f551f8))
* Update dependencies ([#413](https://github.com/chanzuckerberg/happy/issues/413)) ([21ebe14](https://github.com/chanzuckerberg/happy/commit/21ebe145da80c565624933727bb85fde27195813))


### BugFixes

* (CCIE-436) Fix happy shell ([#417](https://github.com/chanzuckerberg/happy/issues/417)) ([ef242c5](https://github.com/chanzuckerberg/happy/commit/ef242c50cf2d439d54327aaa265a07487cd2d729))

## [0.26.1](https://github.com/chanzuckerberg/happy/compare/v0.26.0...v0.26.1) (2022-07-20)


### Misc

* (CCIE-291) update helpstring for happy push ([15807e8](https://github.com/chanzuckerberg/happy/commit/15807e8498104a5b5cb07d45a33e54ee0d45aa0d))
* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#406](https://github.com/chanzuckerberg/happy/issues/406)) ([b1fce72](https://github.com/chanzuckerberg/happy/commit/b1fce727fc7554552d1cafbddf60e2e12a5a282d))
* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#408](https://github.com/chanzuckerberg/happy/issues/408)) ([fba284c](https://github.com/chanzuckerberg/happy/commit/fba284cc3abb9109eb1d35773ea1ec8ebaea9ffd))
* bump github.com/hashicorp/go-tfe from 1.4.0 to 1.5.0 ([#405](https://github.com/chanzuckerberg/happy/issues/405)) ([19a1ce5](https://github.com/chanzuckerberg/happy/commit/19a1ce52544970ed9da6f4f6229dbf84c92c5124))
* bump github.com/sirupsen/logrus from 1.8.1 to 1.9.0 ([#411](https://github.com/chanzuckerberg/happy/issues/411)) ([50fee4b](https://github.com/chanzuckerberg/happy/commit/50fee4b67c3b8eca8dc075f53374b31ee97c9000))


### BugFixes

* Figure out why logs command picks a non-existent stream id for single-cell-data-portal devstack log group ([#412](https://github.com/chanzuckerberg/happy/issues/412)) ([c72ac77](https://github.com/chanzuckerberg/happy/commit/c72ac779236d6711d96878513fc8655b2c0fd058))

## [0.26.0](https://github.com/chanzuckerberg/happy/compare/v0.25.0...v0.26.0) (2022-07-14)


### Features

* Add json output to happy CLI for `happy list` ([#401](https://github.com/chanzuckerberg/happy/issues/401)) ([8fe80a9](https://github.com/chanzuckerberg/happy/commit/8fe80a902f26d716100559c6e424e4a97a413694))


### Misc

* bump github.com/aws/aws-sdk-go-v2/config from 1.15.12 to 1.15.13 ([#388](https://github.com/chanzuckerberg/happy/issues/388)) ([2eba8c6](https://github.com/chanzuckerberg/happy/commit/2eba8c60ccf67beb57464fdc0c0e6ccb29c08528))
* bump github.com/aws/aws-sdk-go-v2/config from 1.15.13 to 1.15.14 ([#404](https://github.com/chanzuckerberg/happy/issues/404)) ([dfad5e8](https://github.com/chanzuckerberg/happy/commit/dfad5e8e9c4c543b959e9abf4127c6d5690b696d))
* bump github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs ([#392](https://github.com/chanzuckerberg/happy/issues/392)) ([d89a195](https://github.com/chanzuckerberg/happy/commit/d89a1956713f649feec9cd9552a019859f727473))
* bump github.com/aws/aws-sdk-go-v2/service/dynamodb ([#397](https://github.com/chanzuckerberg/happy/issues/397)) ([b5bdb4f](https://github.com/chanzuckerberg/happy/commit/b5bdb4f6ed7847d1fbcc5d7ef279f9b14535feab))
* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#395](https://github.com/chanzuckerberg/happy/issues/395)) ([803a1f5](https://github.com/chanzuckerberg/happy/commit/803a1f542599a80c9009d6a6f73e31aec7cb870f))
* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#402](https://github.com/chanzuckerberg/happy/issues/402)) ([aa75783](https://github.com/chanzuckerberg/happy/commit/aa757830894285af688e2827f10ddcb4db411c61))
* bump github.com/aws/aws-sdk-go-v2/service/ecr ([#393](https://github.com/chanzuckerberg/happy/issues/393)) ([07bdcc5](https://github.com/chanzuckerberg/happy/commit/07bdcc5204ffa981a7cc38830316e17ee9afa862))
* bump github.com/aws/aws-sdk-go-v2/service/ecs ([#391](https://github.com/chanzuckerberg/happy/issues/391)) ([f1d244a](https://github.com/chanzuckerberg/happy/commit/f1d244a55bb18126283482eadb656df7effa8538))
* bump github.com/aws/aws-sdk-go-v2/service/secretsmanager ([#394](https://github.com/chanzuckerberg/happy/issues/394)) ([130e447](https://github.com/chanzuckerberg/happy/commit/130e447097946757331fc80c1bcf54cd60591e72))
* bump github.com/aws/aws-sdk-go-v2/service/ssm ([#396](https://github.com/chanzuckerberg/happy/issues/396)) ([3ea4df1](https://github.com/chanzuckerberg/happy/commit/3ea4df1cccc32d6676230230b9617015920005eb))
* bump github.com/aws/aws-sdk-go-v2/service/sts ([#389](https://github.com/chanzuckerberg/happy/issues/389)) ([ea5177a](https://github.com/chanzuckerberg/happy/commit/ea5177a647f2efc95e754af987018a300493313a))
* bump github.com/gruntwork-io/terratest from 0.40.17 to 0.40.18 ([#403](https://github.com/chanzuckerberg/happy/issues/403)) ([5c0e562](https://github.com/chanzuckerberg/happy/commit/5c0e56215a9d105d2ebb61ef2cd8d7c7bcdc941e))
* bump github.com/hashicorp/go-tfe from 1.3.0 to 1.4.0 ([#398](https://github.com/chanzuckerberg/happy/issues/398)) ([c15cb3b](https://github.com/chanzuckerberg/happy/commit/c15cb3bc3fa3258c24eb41e3f86ebfba1e15eb13))

## [0.25.0](https://github.com/chanzuckerberg/happy/compare/v0.24.0...v0.25.0) (2022-06-30)


### Features

* Add `plan` support to the happy CLI ([#364](https://github.com/chanzuckerberg/happy/issues/364)) ([3f1b200](https://github.com/chanzuckerberg/happy/commit/3f1b20032badde383d7f4d8f319ab5f622fdaaef))


### BugFixes

* Consolidated two lines of code retrieving images from ECR ([#370](https://github.com/chanzuckerberg/happy/issues/370)) ([0b4d2de](https://github.com/chanzuckerberg/happy/commit/0b4d2ded45215d414ff35474dde3875e2a9b8ba8))


### Misc

* bump github.com/aws/aws-sdk-go-v2/config from 1.15.11 to 1.15.12 ([#377](https://github.com/chanzuckerberg/happy/issues/377)) ([eacc049](https://github.com/chanzuckerberg/happy/commit/eacc0495ab3c630a2077340e15a453990330e3e7))
* bump github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs ([#382](https://github.com/chanzuckerberg/happy/issues/382)) ([c72d5f0](https://github.com/chanzuckerberg/happy/commit/c72d5f0cd2d52b3991e76eed2ef98c2cfc552330))
* bump github.com/aws/aws-sdk-go-v2/service/dynamodb ([#368](https://github.com/chanzuckerberg/happy/issues/368)) ([dce9b3d](https://github.com/chanzuckerberg/happy/commit/dce9b3d12d8e6b8c740dc5868a9915aa85369fa0))
* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#375](https://github.com/chanzuckerberg/happy/issues/375)) ([6847a3b](https://github.com/chanzuckerberg/happy/commit/6847a3b18b4a5c7fa248302c07ab17566ba2c0d4))
* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#387](https://github.com/chanzuckerberg/happy/issues/387)) ([a85b172](https://github.com/chanzuckerberg/happy/commit/a85b1726d6030c92b30ed0e4dbdb223b3b320a19))
* bump github.com/aws/aws-sdk-go-v2/service/ecr ([#379](https://github.com/chanzuckerberg/happy/issues/379)) ([c12f1d3](https://github.com/chanzuckerberg/happy/commit/c12f1d30ae046e6a9b3035d39b8eacecaa1674fd))
* bump github.com/aws/aws-sdk-go-v2/service/ecs ([#380](https://github.com/chanzuckerberg/happy/issues/380)) ([5ef1305](https://github.com/chanzuckerberg/happy/commit/5ef130589afc81169a0fc06686c7d7fde8691be1))
* bump github.com/aws/aws-sdk-go-v2/service/secretsmanager ([#384](https://github.com/chanzuckerberg/happy/issues/384)) ([860ed8f](https://github.com/chanzuckerberg/happy/commit/860ed8f72d5a8b4ca18feda0bcb31a5c73b39113))
* bump github.com/aws/aws-sdk-go-v2/service/ssm ([#383](https://github.com/chanzuckerberg/happy/issues/383)) ([a684003](https://github.com/chanzuckerberg/happy/commit/a6840035f13d740108fc092fd2ed5af84c5cedd9))
* bump github.com/aws/aws-sdk-go-v2/service/sts ([#376](https://github.com/chanzuckerberg/happy/issues/376)) ([1657955](https://github.com/chanzuckerberg/happy/commit/165795552171a5d8dd6b3fccb35083560f97aeff))
* bump github.com/stretchr/testify from 1.7.4 to 1.7.5 ([#371](https://github.com/chanzuckerberg/happy/issues/371)) ([c682069](https://github.com/chanzuckerberg/happy/commit/c682069744820d75cbc7b75b9f3aa6aa6197809e))
* bump github.com/stretchr/testify from 1.7.5 to 1.8.0 ([#386](https://github.com/chanzuckerberg/happy/issues/386)) ([9cfc9f7](https://github.com/chanzuckerberg/happy/commit/9cfc9f72c68a56defd10600698de52d259c40acd))
* Update coverage ([#373](https://github.com/chanzuckerberg/happy/issues/373)) ([c2b5650](https://github.com/chanzuckerberg/happy/commit/c2b56507bc7a88ab94fb6b38c037671bf4efa566))

## [0.24.0](https://github.com/chanzuckerberg/happy/compare/v0.23.1...v0.24.0) (2022-06-22)


### Features

* use distributed dynamo locks to prevent stacklist race condition ([#315](https://github.com/chanzuckerberg/happy/issues/315)) ([e8b8f70](https://github.com/chanzuckerberg/happy/commit/e8b8f70438de713e8da628553b91f5485c813b1d))


### Misc

* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#366](https://github.com/chanzuckerberg/happy/issues/366)) ([44c48b8](https://github.com/chanzuckerberg/happy/commit/44c48b8290643eb43e4f9c7326a5d7c4499503ad))
* bump github.com/aws/aws-sdk-go-v2/service/ecs ([#365](https://github.com/chanzuckerberg/happy/issues/365)) ([6ad91f3](https://github.com/chanzuckerberg/happy/commit/6ad91f32ef689c3bd61c0169e6c16839247ea33d))
* bump github.com/spf13/cobra from 1.4.0 to 1.5.0 ([#362](https://github.com/chanzuckerberg/happy/issues/362)) ([54958ca](https://github.com/chanzuckerberg/happy/commit/54958cab4560e8e64f1aef55b1b7a84ef4521b18))
* bump github.com/stretchr/testify from 1.7.2 to 1.7.4 ([#361](https://github.com/chanzuckerberg/happy/issues/361)) ([a7b59a0](https://github.com/chanzuckerberg/happy/commit/a7b59a0b0130a6977f5bb1cca0c9d683c8885ede))


### BugFixes

* fix image retagging (missing media formats, and registry id) ([#367](https://github.com/chanzuckerberg/happy/issues/367)) ([1f0191f](https://github.com/chanzuckerberg/happy/commit/1f0191fb419d7780718edb2a4b00933cdeefe076))

## [0.23.1](https://github.com/chanzuckerberg/happy/compare/v0.23.0...v0.23.1) (2022-06-17)


### Misc

* bump github.com/aws/aws-sdk-go-v2/config from 1.15.10 to 1.15.11 ([#358](https://github.com/chanzuckerberg/happy/issues/358)) ([66c0f9f](https://github.com/chanzuckerberg/happy/commit/66c0f9f94ab8581ad55f6550796973b627b0ef4a))
* bump github.com/aws/aws-sdk-go-v2/service/secretsmanager ([#357](https://github.com/chanzuckerberg/happy/issues/357)) ([032c446](https://github.com/chanzuckerberg/happy/commit/032c446f95288ab3c0f66116896beceff9f42bc1))


### BugFixes

* Unable to delete a stack with a non-existent workspace ([#360](https://github.com/chanzuckerberg/happy/issues/360)) ([9b873cd](https://github.com/chanzuckerberg/happy/commit/9b873cd90c06e6a9463025182d6d9a7b821edfc1))

## [0.23.0](https://github.com/chanzuckerberg/happy/compare/v0.22.1...v0.23.0) (2022-06-15)


### Features

* Notify Happy console user of TFE backlogs ([#319](https://github.com/chanzuckerberg/happy/issues/319)) ([b1abc16](https://github.com/chanzuckerberg/happy/commit/b1abc1611c85354a8b237167937d52527ef3466a))


### BugFixes

* Stacks with invalid names cannot be deleted ([#354](https://github.com/chanzuckerberg/happy/issues/354)) ([e71183b](https://github.com/chanzuckerberg/happy/commit/e71183bf403a8be85c377755b152b6d436c7bc56))

## [0.22.1](https://github.com/chanzuckerberg/happy/compare/v0.22.0...v0.22.1) (2022-06-15)


### Misc

* bump github.com/gruntwork-io/terratest from 0.40.16 to 0.40.17 ([#352](https://github.com/chanzuckerberg/happy/issues/352)) ([19d81cb](https://github.com/chanzuckerberg/happy/commit/19d81cbdf2496c637a108cba994cf5b2b0c98a2d))
* Update coverage ([#350](https://github.com/chanzuckerberg/happy/issues/350)) ([47b9f27](https://github.com/chanzuckerberg/happy/commit/47b9f27a2f1759a5bed95dd47729f3ac9ea5c4c6))


### BugFixes

* Cannot delete a stack in napari-hub due to a non-present task ([#353](https://github.com/chanzuckerberg/happy/issues/353)) ([9b2a236](https://github.com/chanzuckerberg/happy/commit/9b2a23650f4a36fccec0c326a2f7ab0ad7f1c483))

## [0.22.0](https://github.com/chanzuckerberg/happy/compare/v0.21.1...v0.22.0) (2022-06-14)


### Features

* Migrate happy-deploy.py to a GitHub action ([#345](https://github.com/chanzuckerberg/happy/issues/345)) ([e9f0f71](https://github.com/chanzuckerberg/happy/commit/e9f0f71ee1b3132ea67a91b34bbbb61015137c0b))


### Misc

* bump github.com/gruntwork-io/terratest from 0.40.15 to 0.40.16 ([#346](https://github.com/chanzuckerberg/happy/issues/346)) ([38e992b](https://github.com/chanzuckerberg/happy/commit/38e992bc88fa6e1e23edefabaa98ebe15140bdd7))
* bump github.com/hashicorp/go-tfe from 1.2.0 to 1.3.0 ([#347](https://github.com/chanzuckerberg/happy/issues/347)) ([51428bc](https://github.com/chanzuckerberg/happy/commit/51428bcef07a1f792c13badd5199a2e4682f3cad))


### BugFixes

* Cloudwatch log times out ([#349](https://github.com/chanzuckerberg/happy/issues/349)) ([0febfd3](https://github.com/chanzuckerberg/happy/commit/0febfd373441b37375bdfa4c7a70952113b1692c))

## [0.21.1](https://github.com/chanzuckerberg/happy/compare/v0.21.0...v0.21.1) (2022-06-08)


### Misc

* bump github.com/AlecAivazis/survey/v2 from 2.3.4 to 2.3.5 ([#340](https://github.com/chanzuckerberg/happy/issues/340)) ([8fc0594](https://github.com/chanzuckerberg/happy/commit/8fc0594c9c5c57653a1b6b3fd9f2867d5e7e9ef0))
* bump github.com/aws/aws-sdk-go-v2 from 1.16.4 to 1.16.5 ([#338](https://github.com/chanzuckerberg/happy/issues/338)) ([462a2f6](https://github.com/chanzuckerberg/happy/commit/462a2f6819bb57533b909373084880c988f4cec3))
* bump github.com/aws/aws-sdk-go-v2/config from 1.15.9 to 1.15.10 ([#336](https://github.com/chanzuckerberg/happy/issues/336)) ([9dead21](https://github.com/chanzuckerberg/happy/commit/9dead211613fe69272ff56237f99b53bcbd1b937))
* bump github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs ([#335](https://github.com/chanzuckerberg/happy/issues/335)) ([b490de1](https://github.com/chanzuckerberg/happy/commit/b490de18d07ca95e0b43bd45d197a53c6de95101))
* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#341](https://github.com/chanzuckerberg/happy/issues/341)) ([5a08561](https://github.com/chanzuckerberg/happy/commit/5a08561bf306a80cea7ac245eea79cf1c28b41a7))
* bump github.com/aws/aws-sdk-go-v2/service/ecr ([#337](https://github.com/chanzuckerberg/happy/issues/337)) ([5638cf7](https://github.com/chanzuckerberg/happy/commit/5638cf7166f238cb492a4591b49a5ef15fe7986b))
* bump github.com/aws/aws-sdk-go-v2/service/secretsmanager ([#342](https://github.com/chanzuckerberg/happy/issues/342)) ([20be215](https://github.com/chanzuckerberg/happy/commit/20be2158d0d5ebbf9c6c8b90bb8c972062604316))
* bump github.com/aws/aws-sdk-go-v2/service/ssm ([#339](https://github.com/chanzuckerberg/happy/issues/339)) ([8bb47ef](https://github.com/chanzuckerberg/happy/commit/8bb47ef051afaf9f2978658abfef05f23365c73e))
* bump github.com/aws/aws-sdk-go-v2/service/sts ([#334](https://github.com/chanzuckerberg/happy/issues/334)) ([f55c492](https://github.com/chanzuckerberg/happy/commit/f55c49216072dd9635e2ced325a47d15c618d401))
* Update coverage ([#331](https://github.com/chanzuckerberg/happy/issues/331)) ([bc046db](https://github.com/chanzuckerberg/happy/commit/bc046db59fe1a814d2e0bc1647364bd81b37801b))
* Update dependencies ([#344](https://github.com/chanzuckerberg/happy/issues/344)) ([c85cb58](https://github.com/chanzuckerberg/happy/commit/c85cb58d1d7876791715c828681e35392efee984))

## [0.21.0](https://github.com/chanzuckerberg/happy/compare/v0.20.0...v0.21.0) (2022-06-07)


### Features

* refactor performance profiler into a context [CCIE-4] ([#310](https://github.com/chanzuckerberg/happy/issues/310)) ([c3ae92c](https://github.com/chanzuckerberg/happy/commit/c3ae92c1ec953825ada66722bbd73a834ad36bd4))


### Misc

* bump github.com/aws/aws-sdk-go-v2/config from 1.15.7 to 1.15.8 ([#317](https://github.com/chanzuckerberg/happy/issues/317)) ([5ac29a2](https://github.com/chanzuckerberg/happy/commit/5ac29a27d3463694c85c95aa52e5aaa282c8ce38))
* bump github.com/aws/aws-sdk-go-v2/config from 1.15.8 to 1.15.9 ([#322](https://github.com/chanzuckerberg/happy/issues/322)) ([cefc121](https://github.com/chanzuckerberg/happy/commit/cefc121862e3f5c22173957045b92873fe559691))
* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#316](https://github.com/chanzuckerberg/happy/issues/316)) ([b4314e7](https://github.com/chanzuckerberg/happy/commit/b4314e78482ba5a4755db7623d9ec210af49aee6))
* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#321](https://github.com/chanzuckerberg/happy/issues/321)) ([65722c5](https://github.com/chanzuckerberg/happy/commit/65722c53c1f943f50984e946c6e1f55fee5a64f4))
* bump github.com/aws/aws-sdk-go-v2/service/secretsmanager ([#318](https://github.com/chanzuckerberg/happy/issues/318)) ([6ad8fa2](https://github.com/chanzuckerberg/happy/commit/6ad8fa2f54959a1443911cc038576466a1761d5e))
* bump github.com/docker/docker ([#330](https://github.com/chanzuckerberg/happy/issues/330)) ([eba4005](https://github.com/chanzuckerberg/happy/commit/eba4005ade06dd0f97119a7b5a3169bcd802e7c8))
* bump github.com/gruntwork-io/terratest from 0.40.10 to 0.40.15 ([#327](https://github.com/chanzuckerberg/happy/issues/327)) ([fd822b5](https://github.com/chanzuckerberg/happy/commit/fd822b50cd9168a0914335a73dcf10d17af555e2))
* bump github.com/stretchr/testify from 1.7.1 to 1.7.2 ([#329](https://github.com/chanzuckerberg/happy/issues/329)) ([c0f1a20](https://github.com/chanzuckerberg/happy/commit/c0f1a20c81a8d2f87d74a7f8ba2a1eabda3acb23))


### BugFixes

* CCIE-220 rate-limit to GetLogEvents ([#326](https://github.com/chanzuckerberg/happy/issues/326)) ([9d509c9](https://github.com/chanzuckerberg/happy/commit/9d509c9a2c09904d7138279ee45e7bffca4c9ad8))
* If the migration/deletion task succeeds too quickly, happy indicates failure status ([#328](https://github.com/chanzuckerberg/happy/issues/328)) ([29a296d](https://github.com/chanzuckerberg/happy/commit/29a296d5df098ef276e365b9f6e2cb2ee6ff6db7))

## [0.20.0](https://github.com/chanzuckerberg/happy/compare/v0.19.1...v0.20.0) (2022-05-20)


### Features

* CCIE-174 When running in GHA; set owner as GHA actor  ([#311](https://github.com/chanzuckerberg/happy/issues/311)) ([76fcb55](https://github.com/chanzuckerberg/happy/commit/76fcb55fec049b5e5a88562b1588e58cf99cab25))
* List also prints the date last updated ([#312](https://github.com/chanzuckerberg/happy/issues/312)) ([1637f65](https://github.com/chanzuckerberg/happy/commit/1637f65c445ba0327ea25a26e6f0cd7eeb7c6000))
* list prints stacks sorted by name ([#308](https://github.com/chanzuckerberg/happy/issues/308)) ([62afc12](https://github.com/chanzuckerberg/happy/commit/62afc12c7e9c0edecb3f82d5a79ff935f13616f1))


### Misc

* bump github.com/aws/aws-sdk-go-v2/config from 1.15.5 to 1.15.6 ([#294](https://github.com/chanzuckerberg/happy/issues/294)) ([8f803c1](https://github.com/chanzuckerberg/happy/commit/8f803c165fd4b2c3915c16339e8c82f4569d1d75))
* bump github.com/aws/aws-sdk-go-v2/service/secretsmanager ([#305](https://github.com/chanzuckerberg/happy/issues/305)) ([3f67460](https://github.com/chanzuckerberg/happy/commit/3f6746030463202dde143b7d3bc027fa4c4ba0cc))
* bump github.com/aws/aws-sdk-go-v2/service/sts ([#304](https://github.com/chanzuckerberg/happy/issues/304)) ([01f9113](https://github.com/chanzuckerberg/happy/commit/01f9113ada653134601b7a08365ae58616ea058f))
* simplify logic for determining list of stack names ([#296](https://github.com/chanzuckerberg/happy/issues/296)) ([49a2505](https://github.com/chanzuckerberg/happy/commit/49a2505929bfe6737bd5c45f0726efd236d79567))
* upgrade deps ([#314](https://github.com/chanzuckerberg/happy/issues/314)) ([8f15160](https://github.com/chanzuckerberg/happy/commit/8f15160c501549151b0427c949af2ec5cc8003ad))

### [0.19.1](https://github.com/chanzuckerberg/happy/compare/v0.19.0...v0.19.1) (2022-05-13)


### Misc

* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#287](https://github.com/chanzuckerberg/happy/issues/287)) ([4b69320](https://github.com/chanzuckerberg/happy/commit/4b693205b3ef016cce5e91ddff4c06c29858cfa9))
* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#290](https://github.com/chanzuckerberg/happy/issues/290)) ([740f873](https://github.com/chanzuckerberg/happy/commit/740f873b8b0adbae7783fdf527b7150727d1812d))
* bump github.com/aws/aws-sdk-go-v2/service/secretsmanager ([#286](https://github.com/chanzuckerberg/happy/issues/286)) ([5d2df11](https://github.com/chanzuckerberg/happy/commit/5d2df11a750f420d77d042651de5db6134a044b9))
* bump github.com/docker/docker ([#291](https://github.com/chanzuckerberg/happy/issues/291)) ([0b93a41](https://github.com/chanzuckerberg/happy/commit/0b93a41b51ef2859aba22509423593eeab98adf6))


### BugFixes

* add HAPPY_ENV={env} to docker build args ([#292](https://github.com/chanzuckerberg/happy/issues/292)) ([c6d122f](https://github.com/chanzuckerberg/happy/commit/c6d122f87d01e6530746f47d3bcd492e14f34467))
* fix formatting of debug log for slow applies ([#289](https://github.com/chanzuckerberg/happy/issues/289)) ([8dade44](https://github.com/chanzuckerberg/happy/commit/8dade4417c8e958c2d28ba58a7cbb35a25f800d0))

## [0.19.0](https://github.com/chanzuckerberg/happy/compare/v0.18.0...v0.19.0) (2022-05-11)


### Features

* Warn if target compute platform is not specified in docker-compose.yml ([#270](https://github.com/chanzuckerberg/happy/issues/270)) ([9d82b1e](https://github.com/chanzuckerberg/happy/commit/9d82b1e46fe86c6e4a64dc028957e987651d9e4c))


### Misc

* bump github.com/aws/aws-sdk-go-v2/config from 1.15.4 to 1.15.5 ([#283](https://github.com/chanzuckerberg/happy/issues/283)) ([afe650b](https://github.com/chanzuckerberg/happy/commit/afe650b0cec49170512101737d76d5b31ce2ee19))
* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#281](https://github.com/chanzuckerberg/happy/issues/281)) ([99955dc](https://github.com/chanzuckerberg/happy/commit/99955dcd06d7bba38944c66aafbfab757c3e4e52))
* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#285](https://github.com/chanzuckerberg/happy/issues/285)) ([fc43737](https://github.com/chanzuckerberg/happy/commit/fc4373714f24c67431a73e322ef4c5d6139915df))
* Update coverage ([#284](https://github.com/chanzuckerberg/happy/issues/284)) ([cca954a](https://github.com/chanzuckerberg/happy/commit/cca954a583d0e14ce2b0773cb095c2bb46fc9879))

## [0.18.0](https://github.com/chanzuckerberg/happy/compare/v0.17.1...v0.18.0) (2022-05-06)


### Features

* Allow overriding migration behavior via flags ([#272](https://github.com/chanzuckerberg/happy/issues/272)) ([1efa2fb](https://github.com/chanzuckerberg/happy/commit/1efa2fbb662496cf308a5879e2904fe8554c0d43))

### [0.17.1](https://github.com/chanzuckerberg/happy/compare/v0.17.0...v0.17.1) (2022-05-06)


### BugFixes

* also run migrations on Update if requested ([#269](https://github.com/chanzuckerberg/happy/issues/269)) ([d7f7f26](https://github.com/chanzuckerberg/happy/commit/d7f7f26e22806b2d41f9985753820bf7f63a95c7))
* Happy logs should support one-off tasks ([#261](https://github.com/chanzuckerberg/happy/issues/261)) ([dc86600](https://github.com/chanzuckerberg/happy/commit/dc86600804757ea83937355871ae542b7420e20a))
* Support slices in create CLI ([#242](https://github.com/chanzuckerberg/happy/issues/242)) [CCIE-3] ([9da3edc](https://github.com/chanzuckerberg/happy/commit/9da3edc8063180279c7f31b3eab08e614ef14732))


### Misc

* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#264](https://github.com/chanzuckerberg/happy/issues/264)) ([ed75f0f](https://github.com/chanzuckerberg/happy/commit/ed75f0f81a176c3d114c9af56e6c7dbe31466cc1))
* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#273](https://github.com/chanzuckerberg/happy/issues/273)) ([132bb9e](https://github.com/chanzuckerberg/happy/commit/132bb9ed0141c7c875d1fc4a2139ade67b12138e))
* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#276](https://github.com/chanzuckerberg/happy/issues/276)) ([8a53b54](https://github.com/chanzuckerberg/happy/commit/8a53b54d1792235395406c2aff7f306a45cca911))
* bump github.com/aws/aws-sdk-go-v2/service/ssm ([#266](https://github.com/chanzuckerberg/happy/issues/266)) ([d3e3ef9](https://github.com/chanzuckerberg/happy/commit/d3e3ef98e963b3515c08055934709686522a3eb7))
* bump github.com/aws/aws-sdk-go-v2/service/ssm ([#274](https://github.com/chanzuckerberg/happy/issues/274)) ([c976d4a](https://github.com/chanzuckerberg/happy/commit/c976d4a671542b59be31f96843c8f7c7c62dae88))
* bump github.com/docker/docker ([#277](https://github.com/chanzuckerberg/happy/issues/277)) ([c4b64ef](https://github.com/chanzuckerberg/happy/commit/c4b64ef8075ddc5a8e352681b23247ce32831fb7))
* bump github.com/go-playground/validator/v10 ([#267](https://github.com/chanzuckerberg/happy/issues/267)) ([f86f207](https://github.com/chanzuckerberg/happy/commit/f86f20737c4978d0e7e7ab1a3f1e3b545bcdbd28))
* bump github.com/gruntwork-io/terratest from 0.40.7 to 0.40.8 ([#278](https://github.com/chanzuckerberg/happy/issues/278)) ([cdbe6cf](https://github.com/chanzuckerberg/happy/commit/cdbe6cfe25bb27bb16797d3b0edb8c18b48edcc3))
* bump github.com/hashicorp/go-tfe from 1.1.0 to 1.2.0 ([#275](https://github.com/chanzuckerberg/happy/issues/275)) ([7cb6d94](https://github.com/chanzuckerberg/happy/commit/7cb6d9403603d9f7693e8639b3aba143ba2df4e7))
* Update coverage ([#271](https://github.com/chanzuckerberg/happy/issues/271)) ([8d861ef](https://github.com/chanzuckerberg/happy/commit/8d861efe36a3b4e4e391cf349f5afbef6e1850d8))

## [0.17.0](https://github.com/chanzuckerberg/happy/compare/v0.16.3...v0.17.0) (2022-04-27)


### Features

* Force delete support ([#247](https://github.com/chanzuckerberg/happy/issues/247)) ([6299c1c](https://github.com/chanzuckerberg/happy/commit/6299c1cedcbfe29839f923f63130a9f1d2d3f70c))


### Misc

* Add a note to document image pushing + building behavior wrt ECR registries in the integration secret ([#260](https://github.com/chanzuckerberg/happy/issues/260)) ([b48946a](https://github.com/chanzuckerberg/happy/commit/b48946ae0659504b71b7815ec16e5b0616c1fef9))
* bump github.com/aws/aws-sdk-go-v2/config from 1.15.3 to 1.15.4 ([#256](https://github.com/chanzuckerberg/happy/issues/256)) ([25d5bc9](https://github.com/chanzuckerberg/happy/commit/25d5bc9a7dbf973cc05ccb4d1e9aa7dd195b02dd))
* bump github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs ([#258](https://github.com/chanzuckerberg/happy/issues/258)) ([0b8280c](https://github.com/chanzuckerberg/happy/commit/0b8280ce92e3eecff70e47d928408724b624264b))
* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#249](https://github.com/chanzuckerberg/happy/issues/249)) ([a42a4e5](https://github.com/chanzuckerberg/happy/commit/a42a4e563325ddbebf598aab5987dff1e35287cd))
* bump github.com/aws/aws-sdk-go-v2/service/secretsmanager ([#244](https://github.com/chanzuckerberg/happy/issues/244)) ([3f23cd2](https://github.com/chanzuckerberg/happy/commit/3f23cd26285fe7b1a0b232eaa5a000ac92f6e7ae))
* bump github.com/aws/aws-sdk-go-v2/service/ssm ([#257](https://github.com/chanzuckerberg/happy/issues/257)) ([c616e59](https://github.com/chanzuckerberg/happy/commit/c616e593bfd541f20f36406c36bc2b1bcf53bc24))

### [0.16.3](https://github.com/chanzuckerberg/happy/compare/v0.16.2...v0.16.3) (2022-04-20)


### Misc

* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#232](https://github.com/chanzuckerberg/happy/issues/232)) ([986c8a3](https://github.com/chanzuckerberg/happy/commit/986c8a39ee153e1274712163f21caefc4f36c265))
* bump github.com/aws/aws-sdk-go-v2/service/ec2 ([#237](https://github.com/chanzuckerberg/happy/issues/237)) ([4f7967f](https://github.com/chanzuckerberg/happy/commit/4f7967f60ea8614fb17db97a67e8c40c1ff84ca9))
* bump github.com/aws/aws-sdk-go-v2/service/ssm ([#239](https://github.com/chanzuckerberg/happy/issues/239)) ([b8968b4](https://github.com/chanzuckerberg/happy/commit/b8968b462e1b9e23e1c56c5cf4bec364f37248b5))

### [0.16.2](https://github.com/chanzuckerberg/happy/compare/v0.16.1...v0.16.2) (2022-04-14)


### BugFixes

* Happy list should only warn when individual stack fetching fails ([#229](https://github.com/chanzuckerberg/happy/issues/229)) ([1695b43](https://github.com/chanzuckerberg/happy/commit/1695b436052b233bd4d099bcaec20fce9f56f45a))


### Misc

* bump github.com/gruntwork-io/terratest from 0.40.6 to 0.40.7 ([#233](https://github.com/chanzuckerberg/happy/issues/233)) ([c345e45](https://github.com/chanzuckerberg/happy/commit/c345e45f76f7560ff7f3cc4e99fc8a24ecba29ec))
* bump github.com/hashicorp/go-uuid from 1.0.2 to 1.0.3 ([#226](https://github.com/chanzuckerberg/happy/issues/226)) ([800fb13](https://github.com/chanzuckerberg/happy/commit/800fb13a1b6335b9a255a1c72a43c5277fab0555))
* Update coverage ([#231](https://github.com/chanzuckerberg/happy/issues/231)) ([2414b1c](https://github.com/chanzuckerberg/happy/commit/2414b1cd00436070e02ef87acd9a5b8878f860a4))

### [0.16.1](https://github.com/chanzuckerberg/happy/compare/v0.16.0...v0.16.1) (2022-04-11)


### BugFixes

* Get logs path from ECS instead of generating ([#216](https://github.com/chanzuckerberg/happy/issues/216)) ([6c01055](https://github.com/chanzuckerberg/happy/commit/6c010558b60238cec2a6925ea905d6ccbc8033f1))

## [0.16.0](https://github.com/chanzuckerberg/happy/compare/v0.15.3...v0.16.0) (2022-04-06)


### Features

* Disable colors via --no-color flag or NO_COLOR env ([#223](https://github.com/chanzuckerberg/happy/issues/223)) ([2d23ef2](https://github.com/chanzuckerberg/happy/commit/2d23ef2043bdd31a37b80a4adde7112ed666f845))


### Misc

* Upgrade aws-sdk-go to v2 ([#212](https://github.com/chanzuckerberg/happy/issues/212)) ([86203c5](https://github.com/chanzuckerberg/happy/commit/86203c5aaf08bbd27817277467d5c1102f0b4a50))


### BugFixes

* Do not show TFE output when in CI ([#222](https://github.com/chanzuckerberg/happy/issues/222)) ([8107bc0](https://github.com/chanzuckerberg/happy/commit/8107bc0b7235188698e70d076412907aba5c0506))

### [0.15.3](https://github.com/chanzuckerberg/happy/compare/v0.15.2...v0.15.3) (2022-03-30)


### Misc

* bump github.com/aws/aws-sdk-go from 1.43.27 to 1.43.28 ([#202](https://github.com/chanzuckerberg/happy/issues/202)) ([b849321](https://github.com/chanzuckerberg/happy/commit/b84932139fb8c1084965e9e5cd7520b432935bbf))


### BugFixes

* Update command suppports --create-tag flag like Create does ([#208](https://github.com/chanzuckerberg/happy/issues/208)) ([9b1b5dd](https://github.com/chanzuckerberg/happy/commit/9b1b5dd35bd5edea687b06077d45cd95f5aa4576))

### [0.15.2](https://github.com/chanzuckerberg/happy/compare/v0.15.1...v0.15.2) (2022-03-29)


### Misc

* bump github.com/AlecAivazis/survey/v2 from 2.3.2 to 2.3.4 ([#186](https://github.com/chanzuckerberg/happy/issues/186)) ([08696fe](https://github.com/chanzuckerberg/happy/commit/08696fe5c8c33d915a7cda76430d1e97b797e65c))
* bump github.com/aws/aws-sdk-go from 1.43.23 to 1.43.25 ([#193](https://github.com/chanzuckerberg/happy/issues/193)) ([0d85730](https://github.com/chanzuckerberg/happy/commit/0d857302b87cf90c05af2717c3f87c112bf21849))
* bump github.com/aws/aws-sdk-go from 1.43.25 to 1.43.26 ([#195](https://github.com/chanzuckerberg/happy/issues/195)) ([72b6337](https://github.com/chanzuckerberg/happy/commit/72b633725a3d7b7f0839a4b942e3efcc8ed2bbe5))
* bump github.com/aws/aws-sdk-go from 1.43.26 to 1.43.27 ([#198](https://github.com/chanzuckerberg/happy/issues/198)) ([44a1379](https://github.com/chanzuckerberg/happy/commit/44a13797e3d3876d774b7a42247fd7820fa10b59))
* bump github.com/aws/aws-sdk-go-v2/config from 1.14.0 to 1.15.2 ([#192](https://github.com/chanzuckerberg/happy/issues/192)) ([78bc05e](https://github.com/chanzuckerberg/happy/commit/78bc05e340f3dce2095060a76ce99ec39ccfde86))
* bump github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs ([#194](https://github.com/chanzuckerberg/happy/issues/194)) ([c85ae6b](https://github.com/chanzuckerberg/happy/commit/c85ae6b06a47d5695a7e3f2a7bfe2bfabd2321cf))
* bump github.com/docker/docker ([#187](https://github.com/chanzuckerberg/happy/issues/187)) ([f19ccd3](https://github.com/chanzuckerberg/happy/commit/f19ccd373e2b2f718f232beccb6d27fddb41a08d))
* Update coverage ([#182](https://github.com/chanzuckerberg/happy/issues/182)) ([ba0f039](https://github.com/chanzuckerberg/happy/commit/ba0f0394a7419fa3df14da1e9c821dbba1c991a3))
* Upgrade tfe client to v1.0.0 and other dependencies ([#191](https://github.com/chanzuckerberg/happy/issues/191)) ([8ecd194](https://github.com/chanzuckerberg/happy/commit/8ecd19454a70ea06fdb0eaf6311056d87b0341ea))


### BugFixes

* Docker uses buildkit [#201](https://github.com/chanzuckerberg/happy/issues/201) ([3f9521e](https://github.com/chanzuckerberg/happy/commit/3f9521ecb88c5d9ef06f4a61bfd41dee22a7740d))

### [0.15.1](https://github.com/chanzuckerberg/happy/compare/v0.15.0...v0.15.1) (2022-03-23)


### Misc

* bump github.com/aws/aws-sdk-go from 1.43.21 to 1.43.23 ([#183](https://github.com/chanzuckerberg/happy/issues/183)) ([e9934ce](https://github.com/chanzuckerberg/happy/commit/e9934ce985a23a3f18232c5017be1bb5b58578de))

## [0.15.0](https://github.com/chanzuckerberg/happy/compare/v0.14.0...v0.15.0) (2022-03-21)


### Features

* ensure that stack name inputs to CLI are compatible with DNS character set ([#168](https://github.com/chanzuckerberg/happy/issues/168)) ([f9b030c](https://github.com/chanzuckerberg/happy/commit/f9b030cc7e879798478408eedd3232d9efd0e773))

## [0.14.0](https://github.com/chanzuckerberg/happy/compare/v0.13.0...v0.14.0) (2022-03-21)


### Features

* Added support for streaming of logs, 'get stack', and ability to force update when stack doesn't exist, fixes [#137](https://github.com/chanzuckerberg/happy/issues/137), [#139](https://github.com/chanzuckerberg/happy/issues/139) ([#143](https://github.com/chanzuckerberg/happy/issues/143)) ([71dbf72](https://github.com/chanzuckerberg/happy/commit/71dbf72c500ecda5d477c85692cdb411b1ffc6cc))


### Misc

* bump github.com/aws/aws-sdk-go from 1.43.14 to 1.43.17 ([#169](https://github.com/chanzuckerberg/happy/issues/169)) ([ada9ea5](https://github.com/chanzuckerberg/happy/commit/ada9ea53409771d51d64968ba2753c5256a18a4d))
* bump github.com/aws/aws-sdk-go from 1.43.17 to 1.43.18 ([#171](https://github.com/chanzuckerberg/happy/issues/171)) ([399503c](https://github.com/chanzuckerberg/happy/commit/399503c673ded958dc8a8957a451f096b5bf904b))
* bump github.com/aws/aws-sdk-go from 1.43.18 to 1.43.21 ([#177](https://github.com/chanzuckerberg/happy/issues/177)) ([d43117a](https://github.com/chanzuckerberg/happy/commit/d43117a9b40e70f8ec450d82d5b9fb1b770258e3))
* bump github.com/docker/docker ([#166](https://github.com/chanzuckerberg/happy/issues/166)) ([60eb892](https://github.com/chanzuckerberg/happy/commit/60eb8920c7cb9cdab572e07b226f8c6230cc3fe4))
* bump github.com/spf13/cobra from 1.3.0 to 1.4.0 ([#167](https://github.com/chanzuckerberg/happy/issues/167)) ([bfd1830](https://github.com/chanzuckerberg/happy/commit/bfd18306f233b78c7a1b53868c4e87daa0c4b19a))
* bump github.com/stretchr/testify from 1.7.0 to 1.7.1 ([#173](https://github.com/chanzuckerberg/happy/issues/173)) ([2e9c501](https://github.com/chanzuckerberg/happy/commit/2e9c501da9f5e93b0d362aa0e954d888b51a422d))


### BugFixes

* docker-compose version contraints from panicing ([#162](https://github.com/chanzuckerberg/happy/issues/162)) ([2c3da6c](https://github.com/chanzuckerberg/happy/commit/2c3da6ce4c0864f8b9de2758833d8c27162b552d))
* Refresh existing TFE token via browser instead of always creating new one ([#172](https://github.com/chanzuckerberg/happy/issues/172)) ([b1cd38e](https://github.com/chanzuckerberg/happy/commit/b1cd38ee54fda478cc10ce381bef4146b4710c14))

## [0.13.0](https://github.com/chanzuckerberg/happy/compare/v0.12.0...v0.13.0) (2022-03-09)


### Features

* Pre-releaser github workflow ([#148](https://github.com/chanzuckerberg/happy/issues/148)) ([fcf0ffd](https://github.com/chanzuckerberg/happy/commit/fcf0ffd0f76a143a3ab33d5d4fa11208db24674b))
* Profile the runtime of happy create ([#147](https://github.com/chanzuckerberg/happy/issues/147)) ([0c4315d](https://github.com/chanzuckerberg/happy/commit/0c4315dbb823a1f6839df162f2a78b66d099cd08))


### Misc

* bump github.com/aws/aws-sdk-go from 1.43.12 to 1.43.13 ([#151](https://github.com/chanzuckerberg/happy/issues/151)) ([28febf1](https://github.com/chanzuckerberg/happy/commit/28febf1a6b505594d779c580ea344e6b8cf473b9))
* bump github.com/aws/aws-sdk-go from 1.43.13 to 1.43.14 ([#155](https://github.com/chanzuckerberg/happy/issues/155)) ([399704b](https://github.com/chanzuckerberg/happy/commit/399704bc9b081065bb9db1b090fb41332d4444c4))
* bump github.com/aws/aws-sdk-go from 1.43.8 to 1.43.9 ([#136](https://github.com/chanzuckerberg/happy/issues/136)) ([8372812](https://github.com/chanzuckerberg/happy/commit/83728126672272cfa1a2b23b150b483975a7d184))
* bump github.com/aws/aws-sdk-go from 1.43.9 to 1.43.12 ([#145](https://github.com/chanzuckerberg/happy/issues/145)) ([3f28724](https://github.com/chanzuckerberg/happy/commit/3f28724b836595a205b15a7d287f9138f1d7418c))
* bump github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs ([#150](https://github.com/chanzuckerberg/happy/issues/150)) ([8ab7901](https://github.com/chanzuckerberg/happy/commit/8ab79010a3b5ae375ef191c12ad8137b4edd9945))
* bump github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs ([#154](https://github.com/chanzuckerberg/happy/issues/154)) ([fd48193](https://github.com/chanzuckerberg/happy/commit/fd4819338525b3d1053dee55d2b51fd00f584115))
* bump github.com/go-playground/validator/v10 ([#152](https://github.com/chanzuckerberg/happy/issues/152)) ([875cc5f](https://github.com/chanzuckerberg/happy/commit/875cc5fcddc44b09b1c420896ab4b28e48cb230c))
* bump github.com/gruntwork-io/terratest from 0.40.5 to 0.40.6 ([#146](https://github.com/chanzuckerberg/happy/issues/146)) ([fc8b44d](https://github.com/chanzuckerberg/happy/commit/fc8b44df41cac7630eb1f57d8f60b3185836cd70))
* Update coverage ([#134](https://github.com/chanzuckerberg/happy/issues/134)) ([252fad6](https://github.com/chanzuckerberg/happy/commit/252fad674c5324fb3f6501ab1117b24f8c2f008e))


### BugFixes

* Allow users to specify AWS_PROFILE; override config ([#153](https://github.com/chanzuckerberg/happy/issues/153)) ([b542199](https://github.com/chanzuckerberg/happy/commit/b5421994fe3a06d0c85ad0e730b92d1e7d07a845))
* Empty AWS profile means don't set one ([#160](https://github.com/chanzuckerberg/happy/issues/160)) ([4183357](https://github.com/chanzuckerberg/happy/commit/41833571a30078d018564ef49575931ba23cfec2))
* GitHub Release action s/jobs/needs/ ([#156](https://github.com/chanzuckerberg/happy/issues/156)) ([50a4027](https://github.com/chanzuckerberg/happy/commit/50a4027355b09104af0bdc1894efaba426da49b8))
* prerelease needs github token ([#159](https://github.com/chanzuckerberg/happy/issues/159)) ([b9b0202](https://github.com/chanzuckerberg/happy/commit/b9b02029f7642b098c5884a7ccf38f519932fd3c))
* prerelease s/tag/version/ ([#158](https://github.com/chanzuckerberg/happy/issues/158)) ([302b6f7](https://github.com/chanzuckerberg/happy/commit/302b6f707819a3d6f817961b476d4c29f6cff155))
* Prerelease wasn't using annotated tags ([#157](https://github.com/chanzuckerberg/happy/issues/157)) ([60989ae](https://github.com/chanzuckerberg/happy/commit/60989aef7baa7698a3b12561ad1c80f98752705b))

## [0.12.0](https://github.com/chanzuckerberg/happy/compare/v0.11.1...v0.12.0) (2022-03-01)


### Features

* Display links to installation instructions for missing dependencies, resolves [#128](https://github.com/chanzuckerberg/happy/issues/128) ([#129](https://github.com/chanzuckerberg/happy/issues/129)) ([172d22d](https://github.com/chanzuckerberg/happy/commit/172d22d7303b7ea790ac54a977755773e05871d3))

### [0.11.1](https://github.com/chanzuckerberg/happy/compare/v0.11.0...v0.11.1) (2022-02-28)


### BugFixes

* Fetch cloudwatch logs when running tasks ([#126](https://github.com/chanzuckerberg/happy/issues/126)) ([d5e8c6e](https://github.com/chanzuckerberg/happy/commit/d5e8c6e5c8c845e41232c0eb730af6dab4d9201d))

## [0.11.0](https://github.com/chanzuckerberg/happy/compare/v0.10.0...v0.11.0) (2022-02-28)


### Features

* Refresh TFE token if expired ([#120](https://github.com/chanzuckerberg/happy/issues/120)) ([135b70b](https://github.com/chanzuckerberg/happy/commit/135b70b243d50fbcfd04a4c118715ea5f002c513))

## [0.10.0](https://github.com/chanzuckerberg/happy/compare/v0.9.1...v0.10.0) (2022-02-24)


### ⚠ BREAKING CHANGES

* match existing config key to auto_run_migrations (#118)

### Features

* enable dependabot ([#117](https://github.com/chanzuckerberg/happy/issues/117)) ([ed4c973](https://github.com/chanzuckerberg/happy/commit/ed4c973b8752fbc71a2e2043ca7330c429e858bb))


### BugFixes

* match existing config key to auto_run_migrations ([#118](https://github.com/chanzuckerberg/happy/issues/118)) ([03ad529](https://github.com/chanzuckerberg/happy/commit/03ad529a8428ffde421b3700781d985f72f9c0e7))

### [0.9.1](https://github.com/chanzuckerberg/happy/compare/v0.9.0...v0.9.1) (2022-02-24)


### BugFixes

* typo ([#115](https://github.com/chanzuckerberg/happy/issues/115)) ([a0ef4f8](https://github.com/chanzuckerberg/happy/commit/a0ef4f87a4bc78c1bb589812318b6e9164833824))

## [0.9.0](https://github.com/chanzuckerberg/happy/compare/v0.8.0...v0.9.0) (2022-02-24)


### Features

* show usage when cli cmd validation fails ([#112](https://github.com/chanzuckerberg/happy/issues/112)) ([5da10ed](https://github.com/chanzuckerberg/happy/commit/5da10ede8ad88e4149833f5bbb1641a5c804a46f))


### Misc

* Parallelize CI jobs ([#110](https://github.com/chanzuckerberg/happy/issues/110)) ([37b16b0](https://github.com/chanzuckerberg/happy/commit/37b16b080dc77c322c06c9a219f5c84d695db5d4))

## [0.8.0](https://github.com/chanzuckerberg/happy/compare/v0.7.0...v0.8.0) (2022-02-24)


### Features

* add a version command; make sure we update stack's owner tag ([#108](https://github.com/chanzuckerberg/happy/issues/108)) ([ab17468](https://github.com/chanzuckerberg/happy/commit/ab17468256484482e4fc1ae8fdd61ca027607fcf))

## [0.7.0](https://github.com/chanzuckerberg/happy/compare/v0.6.1...v0.7.0) (2022-02-23)


### Features

* check for docker-compose v2 ([#107](https://github.com/chanzuckerberg/happy/issues/107)) ([bbbf036](https://github.com/chanzuckerberg/happy/commit/bbbf03624928924174fb0835c080b87f3fefb9c0))


### Misc

* Builder pattern and launch conditions validation ([#103](https://github.com/chanzuckerberg/happy/issues/103)) ([c5347c0](https://github.com/chanzuckerberg/happy/commit/c5347c0f303b8bb4957203562a9d2123c02d37a0))
* Improve orchestrator package code coverage ([#101](https://github.com/chanzuckerberg/happy/issues/101)) ([5a23d44](https://github.com/chanzuckerberg/happy/commit/5a23d44258990c6789615e698692dc4c1e04f626))
* Remove unused parameters ([#99](https://github.com/chanzuckerberg/happy/issues/99)) ([eb5b4dd](https://github.com/chanzuckerberg/happy/commit/eb5b4dd38465e97b76f67405a1e48b9880bb7339))
* Update coverage ([#105](https://github.com/chanzuckerberg/happy/issues/105)) ([5e58068](https://github.com/chanzuckerberg/happy/commit/5e58068eaa0ff83261d58066de4e42e3d0d09673))


### BugFixes

* Increased golangci-lint timeout ([#102](https://github.com/chanzuckerberg/happy/issues/102)) ([f7925d4](https://github.com/chanzuckerberg/happy/commit/f7925d4f187becf78ea62086f9b9e1c34b7f1173))
* more consistently use slices; imagetags set properly ([#106](https://github.com/chanzuckerberg/happy/issues/106)) ([9b8e8d2](https://github.com/chanzuckerberg/happy/commit/9b8e8d212a1bdce5130487c263161b20d54a441f))
* Updated offending dependencie (docker/distribution and opencontainers/image-spec) ([#104](https://github.com/chanzuckerberg/happy/issues/104)) ([41cb783](https://github.com/chanzuckerberg/happy/commit/41cb783517ae72c602fe87bb0ea22d60f41c1f71))

### [0.6.1](https://github.com/chanzuckerberg/happy/compare/v0.6.0...v0.6.1) (2022-02-18)


### Misc

* Style cleanup ([#96](https://github.com/chanzuckerberg/happy/issues/96)) ([6febcb4](https://github.com/chanzuckerberg/happy/commit/6febcb42d1d8fab26c4b428b0945d272bc91876b))


### BugFixes

* Create --force now works when partial stack exists ([#98](https://github.com/chanzuckerberg/happy/issues/98)) ([94b4f16](https://github.com/chanzuckerberg/happy/commit/94b4f168824f6554b04eca782bc6f9a5cac61981))
* Unused parameter in retagimages and inconsistent receiver names ([#95](https://github.com/chanzuckerberg/happy/issues/95)) ([1105f70](https://github.com/chanzuckerberg/happy/commit/1105f70d74920ec623052499c9543c130547130f))

## [0.6.0](https://github.com/chanzuckerberg/happy/compare/v0.5.0...v0.6.0) (2022-02-17)


### Features

* Audible alert upon failure ([#89](https://github.com/chanzuckerberg/happy/issues/89)) ([bd23143](https://github.com/chanzuckerberg/happy/commit/bd23143b8f7d5a75eba90b6b56029613696e5048))
* Display image tags on happy push ([#91](https://github.com/chanzuckerberg/happy/issues/91)) ([2466a5e](https://github.com/chanzuckerberg/happy/commit/2466a5ecfa0984e388599830b0cfa9a0c8042bd0))


### BugFixes

* Additional tags are ignored on happy push ([#93](https://github.com/chanzuckerberg/happy/issues/93)) ([2e0dc6e](https://github.com/chanzuckerberg/happy/commit/2e0dc6eecc6984c543314cb6e7bd61093f39d672))

## [0.5.0](https://github.com/chanzuckerberg/happy/compare/v0.4.1...v0.5.0) (2022-02-17)


### ⚠ BREAKING CHANGES

* Reinterpret slices, so they are compatible with docker-compose profiles (#77)

### Features

* colorize output and make it more human-readable ([#82](https://github.com/chanzuckerberg/happy/issues/82)) ([bfe0987](https://github.com/chanzuckerberg/happy/commit/bfe0987ba0e4b271ca740eddbf31cf8136d3af4d))
* Friendlier configuration validation messages [#90](https://github.com/chanzuckerberg/happy/issues/90) ([bfbe3a7](https://github.com/chanzuckerberg/happy/commit/bfbe3a7b220b8c302b1654e0851dab0ce0eb160f))
* Reinterpret slices, so they are compatible with docker-compose profiles ([#77](https://github.com/chanzuckerberg/happy/issues/77)) ([80fea88](https://github.com/chanzuckerberg/happy/commit/80fea88c4f25343f47194059618cc7a7b88a3cf6))


### BugFixes

* list uses tableprinter package and handles errors better ([#87](https://github.com/chanzuckerberg/happy/issues/87)) ([69abb4d](https://github.com/chanzuckerberg/happy/commit/69abb4dc4f131a448092a29b7c5fe8943d8d0153))
* Set owner tag when missing ([#88](https://github.com/chanzuckerberg/happy/issues/88)) ([8d2f13a](https://github.com/chanzuckerberg/happy/commit/8d2f13ab953945cbe788a29e6016c39435d15a51))
* TFE run loop should succeed on applied and  run_plan_and_finished ([#86](https://github.com/chanzuckerberg/happy/issues/86)) ([c11d25f](https://github.com/chanzuckerberg/happy/commit/c11d25fdeb8a01dc9a850bccfe95e4caaff88a15))
* user logger; tfe sentinel status planned_and_finished; delete succeeds/noop if stack not found ([#85](https://github.com/chanzuckerberg/happy/issues/85)) ([d684d90](https://github.com/chanzuckerberg/happy/commit/d684d902d5952a1365ab93ad6a8b257645fd226f))


### Misc

* Code coverage ([#80](https://github.com/chanzuckerberg/happy/issues/80)) ([d838a19](https://github.com/chanzuckerberg/happy/commit/d838a19c6dc8aba6401e36ff098a2d78ced3c3fb))
* Improve code coverage ([#76](https://github.com/chanzuckerberg/happy/issues/76)) ([f796107](https://github.com/chanzuckerberg/happy/commit/f796107c2d8b30867cfbc96d1655272bc27efd53))
* make apply message friendlier [#83](https://github.com/chanzuckerberg/happy/issues/83) ([6227f00](https://github.com/chanzuckerberg/happy/commit/6227f0030e4b330cc63e51b98824a7b3ad54c8d7))
* Orchestrator package coverage improvements ([#84](https://github.com/chanzuckerberg/happy/issues/84)) ([085135d](https://github.com/chanzuckerberg/happy/commit/085135d3ed4058de41d8b1505870cf1e5600452d))
* pin linter version ([#92](https://github.com/chanzuckerberg/happy/issues/92)) ([48236f6](https://github.com/chanzuckerberg/happy/commit/48236f6e0040f912bd0fb96e8509320b153b4cc9))
* Update code coverage for workspace_repo package ([#79](https://github.com/chanzuckerberg/happy/issues/79)) ([84ec051](https://github.com/chanzuckerberg/happy/commit/84ec0514479c1a3a470e9d388139a7e41c5112bd))
* Update coverage ([#81](https://github.com/chanzuckerberg/happy/issues/81)) ([ee6c151](https://github.com/chanzuckerberg/happy/commit/ee6c151c0e4002e7d4cf1f63cea0cb3b67e3457d))

### [0.4.1](https://github.com/chanzuckerberg/happy/compare/v0.4.0...v0.4.1) (2022-02-14)


### BugFixes

* Addressed addtags messaging an made flags required ([#73](https://github.com/chanzuckerberg/happy/issues/73)) ([89ce489](https://github.com/chanzuckerberg/happy/commit/89ce489e4d0fa8b01ff3855fda4f24720d77383c))
* Genepi requires imagetags tag to be a valid json even if no tags are present ([#74](https://github.com/chanzuckerberg/happy/issues/74)) ([8259d70](https://github.com/chanzuckerberg/happy/commit/8259d7025c2c9032094f6979dd03a29bd91e2b91))
* Happy migrate fails because task definition is not specified and subnets info is incorrect ([#69](https://github.com/chanzuckerberg/happy/issues/69)) ([020ff7b](https://github.com/chanzuckerberg/happy/commit/020ff7b2038f9b0b070bd72e6a893306b514166b))
* Happy shell isn't working for czgenepi ([#72](https://github.com/chanzuckerberg/happy/issues/72)) ([489a87e](https://github.com/chanzuckerberg/happy/commit/489a87ea6584297445eefa8a1279b3ff560370dd))
* Network configuration didn't have a complete list of subnets and security groups ([#71](https://github.com/chanzuckerberg/happy/issues/71)) ([8c86f75](https://github.com/chanzuckerberg/happy/commit/8c86f75d5dbcad3f68f8e4b801b183fa74af5b02))


### Misc

* Improve code coverage for orchestrator and hostname_manager ([#75](https://github.com/chanzuckerberg/happy/issues/75)) ([e75de18](https://github.com/chanzuckerberg/happy/commit/e75de18cac4f6dbe5b40795c550c75251b905eea))

## [0.4.0](https://github.com/chanzuckerberg/happy/compare/v0.3.1...v0.4.0) (2022-02-10)


### ⚠ BREAKING CHANGES

* For clarity, default_compose_env setting has been superseded by default_compose_env_file

### Features

* Discovery of docker compose env files, absolute and relative ([#56](https://github.com/chanzuckerberg/happy/issues/56)) ([7f19d69](https://github.com/chanzuckerberg/happy/commit/7f19d6927065d555084acd97550e95bfd45410c2))
* Read terraform token from env var, tfrc file, or prompt terraform login ([#58](https://github.com/chanzuckerberg/happy/issues/58)) ([e599e8e](https://github.com/chanzuckerberg/happy/commit/e599e8e9707b26e1d3cd0dc6baf08122bb1a7a5b))
* Switched to docker compose v2 ([#60](https://github.com/chanzuckerberg/happy/issues/60)) ([cf5dcad](https://github.com/chanzuckerberg/happy/commit/cf5dcad9cf8dbffd02ad536b51d7eb7b9d63b60b))


### BugFixes

* AWS Backend set default AWS profile ([#61](https://github.com/chanzuckerberg/happy/issues/61)) ([b9788d2](https://github.com/chanzuckerberg/happy/commit/b9788d27f8329a31ae710ac91816fc70f1331d20))
* Docker tag cannot have an @ sign present ([#64](https://github.com/chanzuckerberg/happy/issues/64)) ([6ff1a5a](https://github.com/chanzuckerberg/happy/commit/6ff1a5a552e932db7a1bbac212cf387839159ceb))
* happy hosts install breaking because of incorrect type casting ([#68](https://github.com/chanzuckerberg/happy/issues/68)) ([6673c0d](https://github.com/chanzuckerberg/happy/commit/6673c0daea588b1a4c1ee54799665c3b07ebdcf4))
* Implement global dockerComposeEnvFile setting with the default fallback ([#55](https://github.com/chanzuckerberg/happy/issues/55)) ([9c1bd78](https://github.com/chanzuckerberg/happy/commit/9c1bd783d279e103fe5394ecd86c9edd61156dee))
* Split env and composeEnv for clarity ([#51](https://github.com/chanzuckerberg/happy/issues/51)) ([140f643](https://github.com/chanzuckerberg/happy/commit/140f643fd138c75eae7e19a1b56ce599c9b4b498))
* testbackend package to make testing the backend easier ([#66](https://github.com/chanzuckerberg/happy/issues/66)) ([a34ccc9](https://github.com/chanzuckerberg/happy/commit/a34ccc93b840794ab734f7970f63c6815cdc383f))
* various: TFE url sanitize; docker login ecr registries; integration secret parsing ([#62](https://github.com/chanzuckerberg/happy/issues/62)) ([4f1b166](https://github.com/chanzuckerberg/happy/commit/4f1b166678ed4f19637c7f2aa326041a3e067510))
* Verify aws profile exists when creating the backend ([#63](https://github.com/chanzuckerberg/happy/issues/63)) ([7b6689e](https://github.com/chanzuckerberg/happy/commit/7b6689e96a533db553e46b8e6ec153a0980caba6))
* workspace_repo tests and coverage ([#65](https://github.com/chanzuckerberg/happy/issues/65)) ([89d85c9](https://github.com/chanzuckerberg/happy/commit/89d85c909a27f946fe4b17b87e3b8a7985ecc022))


### Misc

* Added semantic clarity to GetServiceRegistries() method ([#59](https://github.com/chanzuckerberg/happy/issues/59)) ([472572e](https://github.com/chanzuckerberg/happy/commit/472572e466fe5968f64178b2fba473722f880e4a))
* Combined docker compose invokations ([#53](https://github.com/chanzuckerberg/happy/issues/53)) ([c373542](https://github.com/chanzuckerberg/happy/commit/c3735427be4ed184af36e741bda048b463bc179d))
* refactor backend to make it easier to work with and test ([#54](https://github.com/chanzuckerberg/happy/issues/54)) ([bef351b](https://github.com/chanzuckerberg/happy/commit/bef351b6a672d706ae1e1034e01f43efb536674d))
* Silence CLI Usage on errors ([#67](https://github.com/chanzuckerberg/happy/issues/67)) ([678b448](https://github.com/chanzuckerberg/happy/commit/678b4485cc63e68586e067922fccd22acd591ba6))
* Update coverage ([#57](https://github.com/chanzuckerberg/happy/issues/57)) ([4caf025](https://github.com/chanzuckerberg/happy/commit/4caf02508d82d6cea2eccab8c426b82269312479))

### [0.3.1](https://github.com/chanzuckerberg/happy/compare/v0.3.0...v0.3.1) (2022-02-07)


### BugFixes

* goreleaser needs full git history ([#49](https://github.com/chanzuckerberg/happy/issues/49)) ([02c706b](https://github.com/chanzuckerberg/happy/commit/02c706baf34db2039f45f9ce12dbdbd5bef31498))

## [0.3.0](https://github.com/chanzuckerberg/happy/compare/v0.2.1...v0.3.0) (2022-02-07)


### ⚠ BREAKING CHANGES

* Enforce coverage

### Features

* Happy CLI Fargate support ([#27](https://github.com/chanzuckerberg/happy/issues/27)) ([39ae2bd](https://github.com/chanzuckerberg/happy/commit/39ae2bd26bbaffd53aab5fa8cd244f07f83d0bb7))
* Happy create/update should fail if the specified tag is missing ([#32](https://github.com/chanzuckerberg/happy/issues/32)) ([0c07556](https://github.com/chanzuckerberg/happy/commit/0c0755638173fee211e454b4fddbc23ead28c7d7))
* Initial pass at consolidating configuration ([#31](https://github.com/chanzuckerberg/happy/issues/31)) ([045f768](https://github.com/chanzuckerberg/happy/commit/045f768e42e241ea3ba20ee065b95723afdf445d))
* Make env configurable ([#40](https://github.com/chanzuckerberg/happy/issues/40)) ([282c815](https://github.com/chanzuckerberg/happy/commit/282c8158719f11f9120ca2e54187201d78558e7d))
* Make happy CLI asks questions when updates takes too long ([ddb53b6](https://github.com/chanzuckerberg/happy/commit/ddb53b6431002b9bbb1e850a36066b92d4c14f0c))
* Performance improvement for happy list commands ([#22](https://github.com/chanzuckerberg/happy/issues/22)) ([cf4405c](https://github.com/chanzuckerberg/happy/commit/cf4405ca569e2583734e680ee26bf6329741bd91))
* search for happy root path in current directory tree if available + configure more lint rules  ([#46](https://github.com/chanzuckerberg/happy/issues/46)) ([e49aa96](https://github.com/chanzuckerberg/happy/commit/e49aa96895f44b08137efb1330fcdf9e24edc290))
* Skip tagging of non existing images ([#21](https://github.com/chanzuckerberg/happy/issues/21)) ([199b435](https://github.com/chanzuckerberg/happy/commit/199b435d71528aff7c2f6ebb50cb072cf55a896e))


### Misc

* fix more lint ([#29](https://github.com/chanzuckerberg/happy/issues/29)) ([12e783b](https://github.com/chanzuckerberg/happy/commit/12e783b2c2ead229bd9b87cc00f2d73a5481c96e))
* reduce number of ignored error cases ([#28](https://github.com/chanzuckerberg/happy/issues/28)) ([75132e6](https://github.com/chanzuckerberg/happy/commit/75132e6d233c7d4e34bddbb7e3ba18d42f7cb734))
* removing circular dependencies in config ([#47](https://github.com/chanzuckerberg/happy/issues/47)) ([c23487c](https://github.com/chanzuckerberg/happy/commit/c23487cd6c825f6fae7b0a5c954d8cbd8a9b5e53))
* Update coverage ([#35](https://github.com/chanzuckerberg/happy/issues/35)) ([c4088cb](https://github.com/chanzuckerberg/happy/commit/c4088cb5c7eef7d46928acd4690d89152a0d6260))
* Write silly tests to make coverage kick in ([#24](https://github.com/chanzuckerberg/happy/issues/24)) ([83b0138](https://github.com/chanzuckerberg/happy/commit/83b0138e9db755e05c4b2a6ede4c4eca5f5c5732))


### BugFixes

* Coverage action name should be upgrade-coverage ([768ae24](https://github.com/chanzuckerberg/happy/commit/768ae24137482969581f566ce31296b11e9ad33a))
* Enforce coverage ([cff2fb4](https://github.com/chanzuckerberg/happy/commit/cff2fb4810c18874cee154cf29648d513dcbfc62))
* Image existence check needs to be skipped when the tag is not specified ([#41](https://github.com/chanzuckerberg/happy/issues/41)) ([359792e](https://github.com/chanzuckerberg/happy/commit/359792e65d5ad338b11e76f985d118ca7b941174))
* Non-numeric tag values break happy list ([#42](https://github.com/chanzuckerberg/happy/issues/42)) ([227978d](https://github.com/chanzuckerberg/happy/commit/227978d5089bb1dadd487db4c4ee025dab363307))
* RegistryId is just the AWS account id, not a registry host name ([#37](https://github.com/chanzuckerberg/happy/issues/37)) ([ebc9144](https://github.com/chanzuckerberg/happy/commit/ebc9144a93701d59b1a2dd4055bcc271a9b949b4))
* Usability improvements; primarily logging and console interaction ([#44](https://github.com/chanzuckerberg/happy/issues/44)) ([6a63982](https://github.com/chanzuckerberg/happy/commit/6a63982f103b1ef85314a03fea0f6a151a6f4bf1))
* Use a more friendly docker tag format ([#48](https://github.com/chanzuckerberg/happy/issues/48)) ([84f5952](https://github.com/chanzuckerberg/happy/commit/84f595230a4f09ff2c17e30da8c0b9f735b87a94))
* Use the initialized tag ([#36](https://github.com/chanzuckerberg/happy/issues/36)) ([c14b3e7](https://github.com/chanzuckerberg/happy/commit/c14b3e7402d0ca92198ed70c5644405e5e172063))

### [0.2.1](https://github.com/chanzuckerberg/happy/compare/v0.2.0...v0.2.1) (2022-01-25)


### BugFixes

* release action needs code and Go installed ([#15](https://github.com/chanzuckerberg/happy/issues/15)) ([f81984b](https://github.com/chanzuckerberg/happy/commit/f81984b6b74cda602791964905e6cdbcbe3d66cf))

## [0.2.0](https://github.com/chanzuckerberg/happy/compare/v0.1.2...v0.2.0) (2022-01-25)


### ⚠ BREAKING CHANGES

* Rename happy-deploy to happy everywhere (#14)
* We call the binary "happy" rather than happy deploy (#9)

### Features

* Get CI to a good place ([#12](https://github.com/chanzuckerberg/happy/issues/12)) ([01b5e73](https://github.com/chanzuckerberg/happy/commit/01b5e739daaf7e79e0bf9a970f3b3268f1f4587c))


### Misc

* Add CODEOWNERS ([#10](https://github.com/chanzuckerberg/happy/issues/10)) ([b555bf9](https://github.com/chanzuckerberg/happy/commit/b555bf9f92f0433569eff14db5c1e0b9728e43a4))
* Rename happy-deploy to happy everywhere ([#14](https://github.com/chanzuckerberg/happy/issues/14)) ([d7794a2](https://github.com/chanzuckerberg/happy/commit/d7794a2fc40d0f83c4324be8cb5c989536e3aa67))
* We call the binary "happy" rather than happy deploy ([#9](https://github.com/chanzuckerberg/happy/issues/9)) ([1355910](https://github.com/chanzuckerberg/happy/commit/13559103b1c3151ac9baf942963af034e11df408))

### [0.1.2](https://www.github.com/chanzuckerberg/happy-deploy/compare/v0.1.1...v0.1.2) (2021-12-15)


### BugFixes

* goreleaser GitHub action trigger ([#7](https://www.github.com/chanzuckerberg/happy-deploy/issues/7)) ([8b7984a](https://www.github.com/chanzuckerberg/happy-deploy/commit/8b7984a9ad7f2996dfba9c7534359984e26f2053))

### [0.1.1](https://www.github.com/chanzuckerberg/happy-deploy/compare/v0.1.0...v0.1.1) (2021-12-15)


### BugFixes

* update goreleaser.yml trigger [#5](https://www.github.com/chanzuckerberg/happy-deploy/issues/5) ([adf93da](https://www.github.com/chanzuckerberg/happy-deploy/commit/adf93da43ff1a833c2725a8b2b2ddf99a15285e3))

## [0.1.0](https://www.github.com/chanzuckerberg/happy-deploy/compare/v0.0.8...v0.1.0) (2021-12-15)


### Features

* Configure release process GitHub actions ([#3](https://www.github.com/chanzuckerberg/happy-deploy/issues/3)) ([988a472](https://www.github.com/chanzuckerberg/happy-deploy/commit/988a4727e6a2baeaf52a9fabbda4c8d210b90f05))
