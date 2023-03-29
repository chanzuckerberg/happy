# Changelog

## [4.5.2](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.5.1...happy-env-eks-v4.5.2) (2023-03-29)


### Bug Fixes

* write the SSM parameters into czi-si ([#1453](https://github.com/chanzuckerberg/happy/issues/1453)) ([edc4430](https://github.com/chanzuckerberg/happy/commit/edc4430d8ca7039d141625304712b95f85b4fe6c))

## [4.5.1](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.5.0...happy-env-eks-v4.5.1) (2023-03-28)


### Bug Fixes

* remove optional arguments for ppr ([#1449](https://github.com/chanzuckerberg/happy/issues/1449)) ([f1f95df](https://github.com/chanzuckerberg/happy/commit/f1f95df37dd69bc9c25336d91a060148224cb2f7))

## [4.5.0](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.4.1...happy-env-eks-v4.5.0) (2023-03-27)


### Features

* version bump on happy-env-eks ([#1389](https://github.com/chanzuckerberg/happy/issues/1389)) ([afa081e](https://github.com/chanzuckerberg/happy/commit/afa081e55647fa026fa4dfcd6cdd83d6cddfb95d))

## [4.4.1](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.4.0...happy-env-eks-v4.4.1) (2023-03-17)


### Bug Fixes

* reorganize WAF integration secrets to be independent of WAF apply ([#1383](https://github.com/chanzuckerberg/happy/issues/1383)) ([33c775e](https://github.com/chanzuckerberg/happy/commit/33c775ee46476870c44efc1e3324ba62a2b521e8))

## [4.4.0](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.3.0...happy-env-eks-v4.4.0) (2023-03-17)


### Features

* configure WAF per happy-env-eks ([#1381](https://github.com/chanzuckerberg/happy/issues/1381)) ([77a8c18](https://github.com/chanzuckerberg/happy/commit/77a8c18327afa030e7875ab70ba8cc317a4e4840))

## [4.3.0](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.2.1...happy-env-eks-v4.3.0) (2023-03-15)


### Features

* automate adding OIDC providers for new happy apps to happy api ([#1353](https://github.com/chanzuckerberg/happy/issues/1353)) ([782a650](https://github.com/chanzuckerberg/happy/commit/782a650aa6366d7b8f27d94642c0bb21fd99c10c))


### Bug Fixes

* Remove unused EKS vars from happy-env-eks eks input ([#1370](https://github.com/chanzuckerberg/happy/issues/1370)) ([b0de9f1](https://github.com/chanzuckerberg/happy/commit/b0de9f1ac2cfdd33ce937e6d194df7a9d07173ad))

## [4.2.1](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.2.0...happy-env-eks-v4.2.1) (2023-03-08)


### Bug Fixes

* Adjust rds_dbs typing of rds_cluster_parameters ([#1307](https://github.com/chanzuckerberg/happy/issues/1307)) ([6929c4c](https://github.com/chanzuckerberg/happy/commit/6929c4c7cadc164a4ee0ed72a70c87f6981965e7))

## [4.2.0](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.1.0...happy-env-eks-v4.2.0) (2023-03-08)


### Features

* Add rds_cluster_parameters to happy-env-eks rds variables ([#1303](https://github.com/chanzuckerberg/happy/issues/1303)) ([2ef51a3](https://github.com/chanzuckerberg/happy/commit/2ef51a306d4e6bc5bc5b22b6ae1abaced184bcee))

## [4.1.0](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.0.0...happy-env-eks-v4.1.0) (2023-02-28)


### Features

* prevent changes in dynamic tags ([#1233](https://github.com/chanzuckerberg/happy/issues/1233)) ([5ca2403](https://github.com/chanzuckerberg/happy/commit/5ca2403bf2f52797ed92525f13e700866b91ac77))

## [4.0.0](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v3.2.1...happy-env-eks-v4.0.0) (2023-02-13)


### ⚠ BREAKING CHANGES

* inject happy config to stacks ([#1108](https://github.com/chanzuckerberg/happy/issues/1108))

### Features

* inject happy config to stacks ([#1108](https://github.com/chanzuckerberg/happy/issues/1108)) ([9cb49c7](https://github.com/chanzuckerberg/happy/commit/9cb49c7f7bd6819541510e4f31ab5fd112579457))

## [3.2.1](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v3.2.0...happy-env-eks-v3.2.1) (2023-02-06)


### Bug Fixes

* Mark integration secret as sensitive ([#1096](https://github.com/chanzuckerberg/happy/issues/1096)) ([f5fefc1](https://github.com/chanzuckerberg/happy/commit/f5fefc12f55c04f5e2a8d8138eec12718d6cc958))

## [3.2.0](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v3.1.0...happy-env-eks-v3.2.0) (2023-02-06)


### Features

* Add integration secret output to happy-env-eks ([#1094](https://github.com/chanzuckerberg/happy/issues/1094)) ([3ea1a33](https://github.com/chanzuckerberg/happy/commit/3ea1a33d906394a283294522cdbe82852d8bde3b))

## [3.1.0](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v3.0.0...happy-env-eks-v3.1.0) (2023-02-03)


### Features

* Sample Happy Environment EKS Datadog dashboard ([#1066](https://github.com/chanzuckerberg/happy/issues/1066)) ([b4c9f3f](https://github.com/chanzuckerberg/happy/commit/b4c9f3fb7df7d131093a282cb2b54fe83f1e5143))

## [3.0.0](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v2.1.0...happy-env-eks-v3.0.0) (2023-01-31)


### ⚠ BREAKING CHANGES

* authenticate ALBs for ingresses ([#1021](https://github.com/chanzuckerberg/happy/issues/1021))

### Features

* authenticate ALBs for ingresses ([#1021](https://github.com/chanzuckerberg/happy/issues/1021)) ([7cd9375](https://github.com/chanzuckerberg/happy/commit/7cd937576a11b16cbf07e3babf268649c48c0976))

## [2.1.0](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v2.0.0...happy-env-eks-v2.1.0) (2023-01-31)


### Features

* Add namespace output to happy-env-eks ([#1039](https://github.com/chanzuckerberg/happy/issues/1039)) ([b500c16](https://github.com/chanzuckerberg/happy/commit/b500c1657d360364912410c14a9e717b08cc8ce7))

## [2.0.0](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v1.6.0...happy-env-eks-v2.0.0) (2023-01-27)


### ⚠ BREAKING CHANGES

* add oidc configuration to happy env,int secret ([#1020](https://github.com/chanzuckerberg/happy/issues/1020))

### Features

* add oidc configuration to happy env,int secret ([#1020](https://github.com/chanzuckerberg/happy/issues/1020)) ([d887dff](https://github.com/chanzuckerberg/happy/commit/d887dff7755a6899e2cf09e592a70b906ae53671))

## [1.6.0](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v1.5.0...happy-env-eks-v1.6.0) (2023-01-24)


### Features

* (CCIE-1004) Enable creation of stack-level ingress resources with a context based routing support ([#986](https://github.com/chanzuckerberg/happy/issues/986)) ([f258387](https://github.com/chanzuckerberg/happy/commit/f258387b72c1a0753c2779a79b0de8da56df71f1))

## [1.5.0](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v1.4.0...happy-env-eks-v1.5.0) (2023-01-24)


### Features

* add tags from integration secret ([#990](https://github.com/chanzuckerberg/happy/issues/990)) ([46fcd8a](https://github.com/chanzuckerberg/happy/commit/46fcd8a99118b70add0feaecc0d9dd4358100bf0))

## [1.4.0](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v1.3.0...happy-env-eks-v1.4.0) (2023-01-09)


### Features

* Add oauth proxy bypass paths to happy-eks ([#952](https://github.com/chanzuckerberg/happy/issues/952)) ([b363d8f](https://github.com/chanzuckerberg/happy/commit/b363d8f3ee91c6e0d9cb14c21d895980995d7da8))

## [1.3.0](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v1.2.4...happy-env-eks-v1.3.0) (2023-01-06)


### Features

* release please bump ([#949](https://github.com/chanzuckerberg/happy/issues/949)) ([ac65376](https://github.com/chanzuckerberg/happy/commit/ac6537687b3d2182291c05edf962b79094234a6c))

## [1.2.4](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v1.2.3...happy-env-eks-v1.2.4) (2023-01-05)


### Bug Fixes

* revert main changes from proxy ([#931](https://github.com/chanzuckerberg/happy/issues/931)) ([4b6873c](https://github.com/chanzuckerberg/happy/commit/4b6873cd5b7d6e1efa9c6dbaa960ff7d20c67a27))

## [1.2.3](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v1.2.2...happy-env-eks-v1.2.3) (2022-12-21)


### Bug Fixes

* make security group unique across multi-dbs ([#912](https://github.com/chanzuckerberg/happy/issues/912)) ([6cd790f](https://github.com/chanzuckerberg/happy/commit/6cd790fec881324daa3440dd1a692462653857d3))

## [1.2.2](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v1.2.1...happy-env-eks-v1.2.2) (2022-12-13)


### Bug Fixes

* make var iterable ([#868](https://github.com/chanzuckerberg/happy/issues/868)) ([8d08ee7](https://github.com/chanzuckerberg/happy/commit/8d08ee7c37899931633d51fa317637094bac766e))

## [1.2.1](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v1.2.0...happy-env-eks-v1.2.1) (2022-12-12)


### Bug Fixes

* move required_version ([#855](https://github.com/chanzuckerberg/happy/issues/855)) ([b13832c](https://github.com/chanzuckerberg/happy/commit/b13832ca61af7ed8ca0caec643e24bd8633ea4c1))

## [1.2.0](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v1.1.0...happy-env-eks-v1.2.0) (2022-12-12)


### Features

* terraform cicd ([#847](https://github.com/chanzuckerberg/happy/issues/847)) ([1be9354](https://github.com/chanzuckerberg/happy/commit/1be9354192ce8085fa967c0c9280a772a4bb6daa))

## [1.1.0](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v1.0.0...happy-env-eks-v1.1.0) (2022-12-07)


### Features

* add a service account to pods ([#835](https://github.com/chanzuckerberg/happy/issues/835)) ([203c129](https://github.com/chanzuckerberg/happy/commit/203c1294602160dfc4aacc15adf8ebc91e83af5a))

## 1.0.0 (2022-11-30)


### Features

* terraform modules for happy sharing ([#800](https://github.com/chanzuckerberg/happy/issues/800)) ([d909860](https://github.com/chanzuckerberg/happy/commit/d9098607e37b29c71bdc3ddac9fabd7ba280606b))


### Bug Fixes

* happy module bugs ([#806](https://github.com/chanzuckerberg/happy/issues/806)) ([7a87501](https://github.com/chanzuckerberg/happy/commit/7a875019afda4bc016558ee06c846c940a71a6dd))
