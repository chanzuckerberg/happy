# Changelog

## [0.56.1](https://github.com/chanzuckerberg/happy/compare/api-v0.56.0...api-v0.56.1) (2023-02-21)


### Miscellaneous Chores

* **api:** Synchronize happy platform versions

## [0.56.0](https://github.com/chanzuckerberg/happy/compare/api-v0.55.1...api-v0.56.0) (2023-02-17)


### Features

* allow oidc providers to be configured by env variable ([#1144](https://github.com/chanzuckerberg/happy/issues/1144)) ([a5766bd](https://github.com/chanzuckerberg/happy/commit/a5766bd41ae10100f66d6a72b1418c9b5169f123))


### Bug Fixes

* prevent error when ssm stacklist param doesn't exist ([#1173](https://github.com/chanzuckerberg/happy/issues/1173)) ([21b06f4](https://github.com/chanzuckerberg/happy/commit/21b06f4af78bd7bd48a0cf90b638f5d62a53897c))
* ran 'go mod tidy' ([#1172](https://github.com/chanzuckerberg/happy/issues/1172)) ([fd0fcc7](https://github.com/chanzuckerberg/happy/commit/fd0fcc7782e18229979c7eaa622ecceeadf1b528))

## [0.55.1](https://github.com/chanzuckerberg/happy/compare/api-v0.55.0...api-v0.55.1) (2023-02-13)


### Miscellaneous Chores

* **api:** Synchronize happy platform versions

## [0.55.0](https://github.com/chanzuckerberg/happy/compare/api-v0.54.1...api-v0.55.0) (2023-02-13)


### Features

* inject release version into api image ([#1139](https://github.com/chanzuckerberg/happy/issues/1139)) ([cf8b017](https://github.com/chanzuckerberg/happy/commit/cf8b0175d6367e05146b0ba6359655d9fdb14e5a))

## [0.54.1](https://github.com/chanzuckerberg/happy/compare/api-v0.54.0...api-v0.54.1) (2023-02-13)


### Miscellaneous Chores

* **api:** Synchronize happy platform versions

## [0.54.0](https://github.com/chanzuckerberg/happy/compare/api-v0.53.6...api-v0.54.0) (2023-02-13)


### ⚠ BREAKING CHANGES

* inject happy config to stacks ([#1108](https://github.com/chanzuckerberg/happy/issues/1108))

### Features

* inject happy config to stacks ([#1108](https://github.com/chanzuckerberg/happy/issues/1108)) ([9cb49c7](https://github.com/chanzuckerberg/happy/commit/9cb49c7f7bd6819541510e4f31ab5fd112579457))

## [0.53.6](https://github.com/chanzuckerberg/happy/compare/api-v0.53.5...api-v0.53.6) (2023-02-10)


### Bug Fixes

* update happy api oidc client id ([#1133](https://github.com/chanzuckerberg/happy/issues/1133)) ([d27a82f](https://github.com/chanzuckerberg/happy/commit/d27a82f6f0bd376cd9ae81ae1b9a1e863ad8fd6f))

## [0.53.5](https://github.com/chanzuckerberg/happy/compare/api-v0.53.4...api-v0.53.5) (2023-02-10)


### Miscellaneous Chores

* **api:** Synchronize happy platform versions

## [0.53.4](https://github.com/chanzuckerberg/happy/compare/api-v0.53.3...api-v0.53.4) (2023-02-10)


### Miscellaneous Chores

* **api:** Synchronize happy platform versions

## [0.53.3](https://github.com/chanzuckerberg/happy/compare/api-v0.53.2...api-v0.53.3) (2023-02-09)


### Bug Fixes

* use patched version of happy-stack-ecs modules ([#1121](https://github.com/chanzuckerberg/happy/issues/1121)) ([9807059](https://github.com/chanzuckerberg/happy/commit/98070599562ac303c4e7f34ebba1199bc12e56f7))

## [0.53.2](https://github.com/chanzuckerberg/happy/compare/api-v0.53.1...api-v0.53.2) (2023-02-09)


### Miscellaneous Chores

* **api:** Synchronize happy platform versions

## [0.53.1](https://github.com/chanzuckerberg/happy/compare/api-v0.53.0...api-v0.53.1) (2023-02-09)


### Bug Fixes

* use patched happy-stack-ecs modules ([#1118](https://github.com/chanzuckerberg/happy/issues/1118)) ([217a5e0](https://github.com/chanzuckerberg/happy/commit/217a5e03c9f377c176aec66508bb289813dd9657))

## [0.53.0](https://github.com/chanzuckerberg/happy/compare/api-v0.52.0...api-v0.53.0) (2023-02-09)


### Features

* reject requests from old tf provider versions ([#1106](https://github.com/chanzuckerberg/happy/issues/1106)) ([22461e5](https://github.com/chanzuckerberg/happy/commit/22461e51a253c054306a20aa0369f776e77a5d05))
* use new happy-stack-ecs modules in api ([#1109](https://github.com/chanzuckerberg/happy/issues/1109)) ([992c1f6](https://github.com/chanzuckerberg/happy/commit/992c1f6c727f7da567a8af221e8238e1dd7abe96))


### Bug Fixes

* remove deprecated int secret attribute ([#1112](https://github.com/chanzuckerberg/happy/issues/1112)) ([914b45c](https://github.com/chanzuckerberg/happy/commit/914b45c7ac04c6926ae04e319b37c906e7819069))

## [0.52.0](https://github.com/chanzuckerberg/happy/compare/api-v0.51.0...api-v0.52.0) (2023-02-08)


### Features

* use query string for GET requests to happy api ([#1101](https://github.com/chanzuckerberg/happy/issues/1101)) ([7a18eb8](https://github.com/chanzuckerberg/happy/commit/7a18eb8dd5bc2eaebdb246dbebd44f4c389b17e2))

## [0.51.0](https://github.com/chanzuckerberg/happy/compare/api-v0.50.2...api-v0.51.0) (2023-02-08)


### Features

* CCIE-926 List of all happy stacks for an app env ([#1068](https://github.com/chanzuckerberg/happy/issues/1068)) ([fc8d8b1](https://github.com/chanzuckerberg/happy/commit/fc8d8b1353f822e7768d39734adc533e90c49876))

## [0.50.2](https://github.com/chanzuckerberg/happy/compare/api-v0.50.1...api-v0.50.2) (2023-01-30)


### Miscellaneous Chores

* **api:** Synchronize happy platform versions

## [0.50.1](https://github.com/chanzuckerberg/happy/compare/api-v0.50.0...api-v0.50.1) (2023-01-30)


### Miscellaneous Chores

* **api:** Synchronize happy platform versions

## [0.50.0](https://github.com/chanzuckerberg/happy/compare/api-v0.49.0...api-v0.50.0) (2023-01-27)


### Miscellaneous Chores

* **api:** Synchronize happy platform versions

## [0.49.0](https://github.com/chanzuckerberg/happy/compare/api-v0.48.0...api-v0.49.0) (2023-01-24)


### Features

* allow cors ([#1005](https://github.com/chanzuckerberg/happy/issues/1005)) ([87a5cfe](https://github.com/chanzuckerberg/happy/commit/87a5cfe1a56ff6e272ef5893142ad993fb08ef91))

## [0.48.0](https://github.com/chanzuckerberg/happy/compare/api-v0.47.1...api-v0.48.0) (2023-01-19)


### Miscellaneous Chores

* **api:** Synchronize happy platform versions

## [0.47.1](https://github.com/chanzuckerberg/happy/compare/api-v0.47.0...api-v0.47.1) (2023-01-17)


### Miscellaneous Chores

* **api:** Synchronize happy platform versions

## [0.47.0](https://github.com/chanzuckerberg/happy/compare/api-v0.46.1...api-v0.47.0) (2023-01-17)


### Miscellaneous Chores

* **api:** Synchronize happy platform versions

## [0.46.1](https://github.com/chanzuckerberg/happy/compare/api-v0.46.0...api-v0.46.1) (2023-01-09)


### Miscellaneous Chores

* **api:** Synchronize happy platform versions

## [0.46.0](https://github.com/chanzuckerberg/happy/compare/api-v0.45.0...api-v0.46.0) (2023-01-04)


### Miscellaneous Chores

* **api:** Synchronize happy platform versions

## [0.45.0](https://github.com/chanzuckerberg/happy/compare/api-v0.44.0...api-v0.45.0) (2022-12-21)


### Features

* add api meta-command ([#903](https://github.com/chanzuckerberg/happy/issues/903)) ([b81871b](https://github.com/chanzuckerberg/happy/commit/b81871bf694063ce172267e3dcbfe08d737f4120))

## [0.44.0](https://github.com/chanzuckerberg/happy/compare/api-v0.43.2...api-v0.44.0) (2022-12-20)


### Features

* send auth header in api requests ([#785](https://github.com/chanzuckerberg/happy/issues/785)) ([d83c9b3](https://github.com/chanzuckerberg/happy/commit/d83c9b3c57950b1747d8233166e276d883cda4a7))

## [0.43.2](https://github.com/chanzuckerberg/happy/compare/api-v0.43.1...api-v0.43.2) (2022-12-20)


### Bug Fixes

* re-enable okta for prod env ([#898](https://github.com/chanzuckerberg/happy/issues/898)) ([78a0632](https://github.com/chanzuckerberg/happy/commit/78a06321e7659a85ecbbe338bbcc57dfc79338a6))

## [0.43.1](https://github.com/chanzuckerberg/happy/compare/api-v0.43.0...api-v0.43.1) (2022-12-19)


### Bug Fixes

* prevent verifier error ([#892](https://github.com/chanzuckerberg/happy/issues/892)) ([0d94a85](https://github.com/chanzuckerberg/happy/commit/0d94a8540eb212b7a8dd54790b9d9ba199c1f5ca))

## [0.43.0](https://github.com/chanzuckerberg/happy/compare/api-v0.42.1...api-v0.43.0) (2022-12-16)


### Features

* allow for configuration to present n oidc ([#871](https://github.com/chanzuckerberg/happy/issues/871)) ([09fc0be](https://github.com/chanzuckerberg/happy/commit/09fc0be272c2c51d5749b69d3b61681ff178414a))

## [0.42.1](https://github.com/chanzuckerberg/happy/compare/api-v0.42.0...api-v0.42.1) (2022-12-13)


### Bug Fixes

* upgrade to patched module ([#863](https://github.com/chanzuckerberg/happy/issues/863)) ([853ff25](https://github.com/chanzuckerberg/happy/commit/853ff25536019a6eb087bca31c91fe6d26cde32e))

## [0.42.0](https://github.com/chanzuckerberg/happy/compare/api-v0.41.5...api-v0.42.0) (2022-12-12)


### Miscellaneous Chores

* **api:** Synchronize happy platform versions

## [0.41.5](https://github.com/chanzuckerberg/happy/compare/api-v0.41.4...api-v0.41.5) (2022-12-12)


### Bug Fixes

* upgrade aws provider version ([#845](https://github.com/chanzuckerberg/happy/issues/845)) ([09e5161](https://github.com/chanzuckerberg/happy/commit/09e51613e7e5fc2a8559fd3b00dbf410fe6082f4))

## [0.41.4](https://github.com/chanzuckerberg/happy/compare/api-v0.41.3...api-v0.41.4) (2022-12-08)


### Miscellaneous Chores

* **api:** Synchronize happy platform versions

## [0.41.3](https://github.com/chanzuckerberg/happy/compare/api-v0.41.2...api-v0.41.3) (2022-12-07)


### Miscellaneous Chores

* **api:** Synchronize happy platform versions

## [0.41.2](https://github.com/chanzuckerberg/happy/compare/api-v0.41.1...api-v0.41.2) (2022-12-02)


### Miscellaneous Chores

* **api:** Synchronize happy platform versions

## [0.41.1](https://github.com/chanzuckerberg/happy/compare/api-v0.41.0...api-v0.41.1) (2022-11-17)


### Bug Fixes

* simplify staging/prod dns routing ([#774](https://github.com/chanzuckerberg/happy/issues/774)) ([f4f277e](https://github.com/chanzuckerberg/happy/commit/f4f277ead206a05b0e146a416cd664a5abd2e0cc))

## [0.41.0](https://github.com/chanzuckerberg/happy/compare/api-v0.40.1...api-v0.41.0) (2022-11-16)


### Features

* upgrade happy-ecs-stack to get database env vars in service ([#735](https://github.com/chanzuckerberg/happy/issues/735)) ([4791bc4](https://github.com/chanzuckerberg/happy/commit/4791bc4b038c7971837694469c1b21925e46deb1))


### Bug Fixes

* increase ReadBufferSize and ReadTimeout to allow oidc tokens to be passed ([#725](https://github.com/chanzuckerberg/happy/issues/725)) ([4f4b9ee](https://github.com/chanzuckerberg/happy/commit/4f4b9ee303217781a938923e4aa3cb75245c613f))

## [0.40.1](https://github.com/chanzuckerberg/happy/compare/api-v0.40.0...api-v0.40.1) (2022-11-03)


### Miscellaneous Chores

* **api:** Synchronize happy platform versions

## [0.40.0](https://github.com/chanzuckerberg/happy/compare/api-v0.39.0...api-v0.40.0) (2022-11-03)


### Features

* Add scratch container for Happy-API ([#722](https://github.com/chanzuckerberg/happy/issues/722)) ([d03ac67](https://github.com/chanzuckerberg/happy/commit/d03ac67fb1f71dafe17fd555f2e1792f86ad342c))


### Bug Fixes

* use correct secret_arn and log_group_prefix in happy-api config ([#723](https://github.com/chanzuckerberg/happy/issues/723)) ([e2d9d54](https://github.com/chanzuckerberg/happy/commit/e2d9d54db5d845df8724b1e8811732ec23592bb8))

## [0.39.0](https://github.com/chanzuckerberg/happy/compare/api-v0.38.0...api-v0.39.0) (2022-10-25)


### Features

* create terraform provider for happy-api ([#699](https://github.com/chanzuckerberg/happy/issues/699)) ([3325039](https://github.com/chanzuckerberg/happy/commit/3325039ae0fa433ee4d59307762869ed543b8554))

## [0.38.0](https://github.com/chanzuckerberg/happy/compare/api-v0.37.0...api-v0.38.0) (2022-10-24)


### Features

* disable auth for test and dev environments ([#672](https://github.com/chanzuckerberg/happy/issues/672)) ([6c4000f](https://github.com/chanzuckerberg/happy/commit/6c4000fcbfff2169741d3bbe6c1d3366c71b204d))
* roll out config feature to cli ([#660](https://github.com/chanzuckerberg/happy/issues/660)) ([a72c965](https://github.com/chanzuckerberg/happy/commit/a72c965f6bd2c9113c8152c9155330971e808b46))

## [0.37.0](https://github.com/chanzuckerberg/happy/compare/api-v0.36.0...api-v0.37.0) (2022-10-13)


### Features

* move models to shared package ([#657](https://github.com/chanzuckerberg/happy/issues/657)) ([2f42c9d](https://github.com/chanzuckerberg/happy/commit/2f42c9df6629c2adba23498b320c56cfe58335c0))

## [0.36.0](https://github.com/chanzuckerberg/happy/compare/api-v0.35.2...api-v0.36.0) (2022-10-12)


### Features

* auth in api ([#629](https://github.com/chanzuckerberg/happy/issues/629)) ([25aaa25](https://github.com/chanzuckerberg/happy/commit/25aaa2558f0228bc1f063ff6667160c954313d3e))

## [0.35.2](https://github.com/chanzuckerberg/happy/compare/api-v0.35.1...api-v0.35.2) (2022-10-12)


### Miscellaneous Chores

* **api:** Synchronize happy platform versions

## [0.35.1](https://github.com/chanzuckerberg/happy/compare/api-v0.35.0...api-v0.35.1) (2022-10-11)


### Miscellaneous Chores

* **api:** Synchronize happy platform versions

## [0.35.0](https://github.com/chanzuckerberg/happy/compare/api-v0.34.4...api-v0.35.0) (2022-10-11)


### Miscellaneous Chores

* **api:** Synchronize happy platform versions

## [0.34.4](https://github.com/chanzuckerberg/happy/compare/api-v0.34.3...api-v0.34.4) (2022-10-11)


### Miscellaneous Chores

* **api:** Synchronize happy platform versions

## [0.34.3](https://github.com/chanzuckerberg/happy/compare/api-v0.34.2...api-v0.34.3) (2022-10-10)


### Miscellaneous Chores

* **api:** Synchronize happy platform versions

## [0.34.2](https://github.com/chanzuckerberg/happy/compare/api-v0.34.1...api-v0.34.2) (2022-10-10)


### Miscellaneous Chores

* **api:** Synchronize happy platform versions

## [0.34.1](https://github.com/chanzuckerberg/happy/compare/api-v0.34.0...api-v0.34.1) (2022-10-10)


### Miscellaneous Chores

* **api:** Synchronize happy platform versions

## [0.34.0](https://github.com/chanzuckerberg/happy/compare/api-v0.33.0...api-v0.34.0) (2022-10-10)


### Miscellaneous Chores

* **api:** Synchronize happy platform versions

## [0.33.0](https://github.com/chanzuckerberg/happy/compare/api-v0.32.0...api-v0.33.0) (2022-10-07)


### Miscellaneous Chores

* **api:** Synchronize happy platform versions

## 0.32.0 (2022-10-07)


### Features

* add shared package ([#620](https://github.com/chanzuckerberg/happy/issues/620)) ([159bd8e](https://github.com/chanzuckerberg/happy/commit/159bd8e372cdf4c2897ca71395c1d65667b0b423))
