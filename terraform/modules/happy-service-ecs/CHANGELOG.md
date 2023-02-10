# Changelog

## [1.7.0](https://github.com/chanzuckerberg/happy/compare/happy-service-ecs-v1.6.0...happy-service-ecs-v1.7.0) (2023-02-10)


### Features

* (CCIE-1123) Make DataDog sidecar optional for ECS tasks ([#1130](https://github.com/chanzuckerberg/happy/issues/1130)) ([8921c53](https://github.com/chanzuckerberg/happy/commit/8921c53369c044d356f7f98009dfcef88469a4c1))

## [1.6.0](https://github.com/chanzuckerberg/happy/compare/happy-service-ecs-v1.5.1...happy-service-ecs-v1.6.0) (2023-02-09)


### Features

* Dynamically allocate ECS task resources based on AWS guidelines ([#1122](https://github.com/chanzuckerberg/happy/issues/1122)) ([cf7bca0](https://github.com/chanzuckerberg/happy/commit/cf7bca04b33c65a439535d7fcb7ba6aee48f7b48))

## [1.5.1](https://github.com/chanzuckerberg/happy/compare/happy-service-ecs-v1.5.0...happy-service-ecs-v1.5.1) (2023-02-09)


### Bug Fixes

* task cpu needs to be at least the sum of the containers ([#1120](https://github.com/chanzuckerberg/happy/issues/1120)) ([acaf25f](https://github.com/chanzuckerberg/happy/commit/acaf25f7f09587fb94f607bf8bd392ac7dcf6a5a))

## [1.5.0](https://github.com/chanzuckerberg/happy/compare/happy-service-ecs-v1.4.0...happy-service-ecs-v1.5.0) (2023-02-02)


### Features

* Add grouping tags ([#1060](https://github.com/chanzuckerberg/happy/issues/1060)) ([713015f](https://github.com/chanzuckerberg/happy/commit/713015ff7c24278c6315b9ad0ce04e98fb56bb4e))

## [1.4.0](https://github.com/chanzuckerberg/happy/compare/happy-service-ecs-v1.3.0...happy-service-ecs-v1.4.0) (2023-01-31)


### Features

* Tag stack level metrics for EKS and ECS ([#1033](https://github.com/chanzuckerberg/happy/issues/1033)) ([1466430](https://github.com/chanzuckerberg/happy/commit/146643014a9c60cf2bac67fd25d6881827b9b3e9))

## [1.3.0](https://github.com/chanzuckerberg/happy/compare/happy-service-ecs-v1.2.0...happy-service-ecs-v1.3.0) (2022-12-12)


### Features

* terraform cicd ([#847](https://github.com/chanzuckerberg/happy/issues/847)) ([1be9354](https://github.com/chanzuckerberg/happy/commit/1be9354192ce8085fa967c0c9280a772a4bb6daa))


### Bug Fixes

* duplicate providers ([#849](https://github.com/chanzuckerberg/happy/issues/849)) ([59c45f8](https://github.com/chanzuckerberg/happy/commit/59c45f8b6fbf9b877a8de60662793ccc45292f09))

## [1.2.0](https://github.com/chanzuckerberg/happy/compare/happy-service-ecs-v1.1.0...happy-service-ecs-v1.2.0) (2022-12-08)


### Features

* add optional/depin internal modules ([#846](https://github.com/chanzuckerberg/happy/issues/846)) ([348fc78](https://github.com/chanzuckerberg/happy/commit/348fc7876fd7427487d7ea340171898a39d4b05b))

## [1.1.0](https://github.com/chanzuckerberg/happy/compare/happy-service-ecs-v1.0.0...happy-service-ecs-v1.1.0) (2022-12-07)


### Features

* add a service account to pods ([#835](https://github.com/chanzuckerberg/happy/issues/835)) ([203c129](https://github.com/chanzuckerberg/happy/commit/203c1294602160dfc4aacc15adf8ebc91e83af5a))

## 1.0.0 (2022-11-30)


### Features

* terraform modules for happy sharing ([#800](https://github.com/chanzuckerberg/happy/issues/800)) ([d909860](https://github.com/chanzuckerberg/happy/commit/d9098607e37b29c71bdc3ddac9fabd7ba280606b))


### Bug Fixes

* happy module bugs ([#806](https://github.com/chanzuckerberg/happy/issues/806)) ([7a87501](https://github.com/chanzuckerberg/happy/commit/7a875019afda4bc016558ee06c846c940a71a6dd))
