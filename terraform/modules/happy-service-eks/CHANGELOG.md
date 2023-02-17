# Changelog

## [2.3.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v2.2.0...happy-service-eks-v2.3.0) (2023-02-17)


### Features

* allow users to create bypasses for their OIDC ([#1149](https://github.com/chanzuckerberg/happy/issues/1149)) ([078ee17](https://github.com/chanzuckerberg/happy/commit/078ee17b36436ce92b5ad0efdade143d1f306879))

## [2.2.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v2.1.1...happy-service-eks-v2.2.0) (2023-02-02)


### Features

* Add grouping tags ([#1060](https://github.com/chanzuckerberg/happy/issues/1060)) ([713015f](https://github.com/chanzuckerberg/happy/commit/713015ff7c24278c6315b9ad0ce04e98fb56bb4e))

## [2.1.1](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v2.1.0...happy-service-eks-v2.1.1) (2023-02-01)


### Bug Fixes

* Datadog stack level annotations ([#1048](https://github.com/chanzuckerberg/happy/issues/1048)) ([69d025c](https://github.com/chanzuckerberg/happy/commit/69d025ccad8ad7b39175489b5f3d8a7392863500))

## [2.1.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v2.0.0...happy-service-eks-v2.1.0) (2023-02-01)


### Features

* Enable mapping of additional environment variables from secrets ([#1046](https://github.com/chanzuckerberg/happy/issues/1046)) ([6ef2fea](https://github.com/chanzuckerberg/happy/commit/6ef2feaf13d07a7848f8498ed14653610d1edc94))

## [2.0.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v1.7.0...happy-service-eks-v2.0.0) (2023-01-31)


### ⚠ BREAKING CHANGES

* authenticate ALBs for ingresses ([#1021](https://github.com/chanzuckerberg/happy/issues/1021))

### Features

* authenticate ALBs for ingresses ([#1021](https://github.com/chanzuckerberg/happy/issues/1021)) ([7cd9375](https://github.com/chanzuckerberg/happy/commit/7cd937576a11b16cbf07e3babf268649c48c0976))

## [1.7.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v1.6.0...happy-service-eks-v1.7.0) (2023-01-31)


### Features

* Tag stack level metrics for EKS and ECS ([#1033](https://github.com/chanzuckerberg/happy/issues/1033)) ([1466430](https://github.com/chanzuckerberg/happy/commit/146643014a9c60cf2bac67fd25d6881827b9b3e9))

## [1.6.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v1.5.0...happy-service-eks-v1.6.0) (2023-01-24)


### Features

* (CCIE-1004) Enable creation of stack-level ingress resources with a context based routing support ([#986](https://github.com/chanzuckerberg/happy/issues/986)) ([f258387](https://github.com/chanzuckerberg/happy/commit/f258387b72c1a0753c2779a79b0de8da56df71f1))

## [1.5.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v1.4.2...happy-service-eks-v1.5.0) (2023-01-24)


### Features

* add tags from integration secret ([#990](https://github.com/chanzuckerberg/happy/issues/990)) ([46fcd8a](https://github.com/chanzuckerberg/happy/commit/46fcd8a99118b70add0feaecc0d9dd4358100bf0))

## [1.4.2](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v1.4.1...happy-service-eks-v1.4.2) (2023-01-04)


### Bug Fixes

* make env vars nonsensitive to allow for_each ([#927](https://github.com/chanzuckerberg/happy/issues/927)) ([0aaf238](https://github.com/chanzuckerberg/happy/commit/0aaf23826c54d1980f6947c20a7623076a5954e1))

## [1.4.1](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v1.4.0...happy-service-eks-v1.4.1) (2022-12-21)


### Bug Fixes

* update additional_env_vars for eks stack ([#910](https://github.com/chanzuckerberg/happy/issues/910)) ([3e0cea1](https://github.com/chanzuckerberg/happy/commit/3e0cea11efc9770626e7f1af66080e1d9fcc8be1))

## [1.4.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v1.3.0...happy-service-eks-v1.4.0) (2022-12-21)


### Features

* automatically inject db env vars into eks env ([#908](https://github.com/chanzuckerberg/happy/issues/908)) ([99123b2](https://github.com/chanzuckerberg/happy/commit/99123b2b1648b1b7c6ce756942c9fb925b31e07a))

## [1.3.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v1.2.0...happy-service-eks-v1.3.0) (2022-12-12)


### Features

* terraform cicd ([#847](https://github.com/chanzuckerberg/happy/issues/847)) ([1be9354](https://github.com/chanzuckerberg/happy/commit/1be9354192ce8085fa967c0c9280a772a4bb6daa))

## [1.2.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v1.1.0...happy-service-eks-v1.2.0) (2022-12-08)


### Features

* add optional/depin internal modules ([#846](https://github.com/chanzuckerberg/happy/issues/846)) ([348fc78](https://github.com/chanzuckerberg/happy/commit/348fc7876fd7427487d7ea340171898a39d4b05b))


### Bug Fixes

* remove old files ([#839](https://github.com/chanzuckerberg/happy/issues/839)) ([8659e46](https://github.com/chanzuckerberg/happy/commit/8659e463f73e4ce16f9a43a49e4134f66c6ba518))

## [1.1.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v1.0.1...happy-service-eks-v1.1.0) (2022-12-07)


### Features

* add a service account to pods ([#835](https://github.com/chanzuckerberg/happy/issues/835)) ([203c129](https://github.com/chanzuckerberg/happy/commit/203c1294602160dfc4aacc15adf8ebc91e83af5a))


### Bug Fixes

* bugs in modules ([#837](https://github.com/chanzuckerberg/happy/issues/837)) ([c911306](https://github.com/chanzuckerberg/happy/commit/c91130667c04b449deb4dd82049baf29f17acc01))
* input variables and tags in happy eks modules ([#838](https://github.com/chanzuckerberg/happy/issues/838)) ([175dc76](https://github.com/chanzuckerberg/happy/commit/175dc7652735e5683dced24d8cdfa48101355c72))

## [1.0.1](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v1.0.0...happy-service-eks-v1.0.1) (2022-12-01)


### Bug Fixes

* secret block instead of = ([#813](https://github.com/chanzuckerberg/happy/issues/813)) ([4f70fcd](https://github.com/chanzuckerberg/happy/commit/4f70fcd199d149937f09a9b2c363d0db0790e5ca))

## 1.0.0 (2022-11-30)


### Features

* add integration secret to volume mounted ([#812](https://github.com/chanzuckerberg/happy/issues/812)) ([1a2ae56](https://github.com/chanzuckerberg/happy/commit/1a2ae56d3bb3a4a0eaef6bfc50d18a0aa99e2f1a))
* terraform modules for happy sharing ([#800](https://github.com/chanzuckerberg/happy/issues/800)) ([d909860](https://github.com/chanzuckerberg/happy/commit/d9098607e37b29c71bdc3ddac9fabd7ba280606b))


### Bug Fixes

* happy module bugs ([#806](https://github.com/chanzuckerberg/happy/issues/806)) ([7a87501](https://github.com/chanzuckerberg/happy/commit/7a875019afda4bc016558ee06c846c940a71a6dd))
* replace dashes with underscores ([#811](https://github.com/chanzuckerberg/happy/issues/811)) ([2af12f4](https://github.com/chanzuckerberg/happy/commit/2af12f42e9cceb89985e94f8e08e8d4a19e88915))
