# Changelog

## [2.2.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-ecs-v2.1.0...happy-stack-ecs-v2.2.0) (2023-10-24)


### Features

* bump happy provider version ([#2508](https://github.com/chanzuckerberg/happy/issues/2508)) ([334cb3e](https://github.com/chanzuckerberg/happy/commit/334cb3e673a1e362973fabfa268649a6baa32f5d))

## [2.1.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-ecs-v2.0.1...happy-stack-ecs-v2.1.0) (2023-08-29)


### Features

* Support for ECR tag immutability ([#2376](https://github.com/chanzuckerberg/happy/issues/2376)) ([c1d5f5b](https://github.com/chanzuckerberg/happy/commit/c1d5f5b6e6a093c19ba2a092111842cc0e4f195f))

## [2.0.1](https://github.com/chanzuckerberg/happy/compare/happy-stack-ecs-v2.0.0...happy-stack-ecs-v2.0.1) (2023-08-17)


### Bug Fixes

* Happy update reports success on failed deployment when ECS rolls back the task version ([#2268](https://github.com/chanzuckerberg/happy/issues/2268)) ([7adf8e6](https://github.com/chanzuckerberg/happy/commit/7adf8e654979bedd01c9c824ba1489901524b2d1))
* Upgrade happy provider version in happy-env-ecs ([#2274](https://github.com/chanzuckerberg/happy/issues/2274)) ([8dbb6a0](https://github.com/chanzuckerberg/happy/commit/8dbb6a0d378de4eb906bfb73b953b231791a6d39))

## [2.0.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-ecs-v1.6.3...happy-stack-ecs-v2.0.0) (2023-02-13)


### ⚠ BREAKING CHANGES

* inject happy config to stacks ([#1108](https://github.com/chanzuckerberg/happy/issues/1108))

### Features

* inject happy config to stacks ([#1108](https://github.com/chanzuckerberg/happy/issues/1108)) ([9cb49c7](https://github.com/chanzuckerberg/happy/commit/9cb49c7f7bd6819541510e4f31ab5fd112579457))

## [1.6.3](https://github.com/chanzuckerberg/happy/compare/happy-stack-ecs-v1.6.2...happy-stack-ecs-v1.6.3) (2023-02-09)


### Bug Fixes

* bump happy-stack-ecs to pull in 'Dynamically allocate ECS task resources based on AWS guidelines' change ([#1124](https://github.com/chanzuckerberg/happy/issues/1124)) ([9c17e5d](https://github.com/chanzuckerberg/happy/commit/9c17e5d675cbf795d2bff559cf4d3328ed3341e0))

## [1.6.2](https://github.com/chanzuckerberg/happy/compare/happy-stack-ecs-v1.6.1...happy-stack-ecs-v1.6.2) (2023-02-09)


### Bug Fixes

* task cpu needs to be at least the sum of the containers ([#1120](https://github.com/chanzuckerberg/happy/issues/1120)) ([acaf25f](https://github.com/chanzuckerberg/happy/commit/acaf25f7f09587fb94f607bf8bd392ac7dcf6a5a))

## [1.6.1](https://github.com/chanzuckerberg/happy/compare/happy-stack-ecs-v1.6.0...happy-stack-ecs-v1.6.1) (2023-02-09)


### Bug Fixes

* remove deprecated int secret attribute ([#1112](https://github.com/chanzuckerberg/happy/issues/1112)) ([914b45c](https://github.com/chanzuckerberg/happy/commit/914b45c7ac04c6926ae04e319b37c906e7819069))

## [1.6.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-ecs-v1.5.0...happy-stack-ecs-v1.6.0) (2023-02-02)


### Features

* Add grouping tags ([#1060](https://github.com/chanzuckerberg/happy/issues/1060)) ([713015f](https://github.com/chanzuckerberg/happy/commit/713015ff7c24278c6315b9ad0ce04e98fb56bb4e))

## [1.5.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-ecs-v1.4.0...happy-stack-ecs-v1.5.0) (2023-01-31)


### Features

* Tag stack level metrics for EKS and ECS ([#1033](https://github.com/chanzuckerberg/happy/issues/1033)) ([1466430](https://github.com/chanzuckerberg/happy/commit/146643014a9c60cf2bac67fd25d6881827b9b3e9))

## [1.4.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-ecs-v1.3.0...happy-stack-ecs-v1.4.0) (2023-01-27)


### Features

* add synthetics to ecs stacks module ([#1008](https://github.com/chanzuckerberg/happy/issues/1008)) ([7ad6192](https://github.com/chanzuckerberg/happy/commit/7ad6192edf208908b50ec8ff906994fef4a15829))

## [1.3.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-ecs-v1.2.0...happy-stack-ecs-v1.3.0) (2022-12-12)


### Features

* terraform cicd ([#847](https://github.com/chanzuckerberg/happy/issues/847)) ([1be9354](https://github.com/chanzuckerberg/happy/commit/1be9354192ce8085fa967c0c9280a772a4bb6daa))


### Bug Fixes

* upgrade aws provider version ([#845](https://github.com/chanzuckerberg/happy/issues/845)) ([09e5161](https://github.com/chanzuckerberg/happy/commit/09e51613e7e5fc2a8559fd3b00dbf410fe6082f4))

## [1.2.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-ecs-v1.1.0...happy-stack-ecs-v1.2.0) (2022-12-08)


### Features

* add optional/depin internal modules ([#846](https://github.com/chanzuckerberg/happy/issues/846)) ([348fc78](https://github.com/chanzuckerberg/happy/commit/348fc7876fd7427487d7ea340171898a39d4b05b))

## [1.1.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-ecs-v1.0.0...happy-stack-ecs-v1.1.0) (2022-12-07)


### Features

* add a service account to pods ([#835](https://github.com/chanzuckerberg/happy/issues/835)) ([203c129](https://github.com/chanzuckerberg/happy/commit/203c1294602160dfc4aacc15adf8ebc91e83af5a))

## 1.0.0 (2022-11-30)


### Features

* terraform modules for happy sharing ([#800](https://github.com/chanzuckerberg/happy/issues/800)) ([d909860](https://github.com/chanzuckerberg/happy/commit/d9098607e37b29c71bdc3ddac9fabd7ba280606b))
