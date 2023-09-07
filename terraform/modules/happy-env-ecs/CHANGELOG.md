# Changelog

## [3.2.1](https://github.com/chanzuckerberg/happy/compare/happy-env-ecs-v3.2.0...happy-env-ecs-v3.2.1) (2023-09-07)


### Bug Fixes

* wrong outputs and ecr types ([#2420](https://github.com/chanzuckerberg/happy/issues/2420)) ([b1cd390](https://github.com/chanzuckerberg/happy/commit/b1cd39024d1b70cf987378768c02b55d07569cf1))

## [3.2.0](https://github.com/chanzuckerberg/happy/compare/happy-env-ecs-v3.1.0...happy-env-ecs-v3.2.0) (2023-09-07)


### Features

* bump the ecs multidomain proxy module ([#2418](https://github.com/chanzuckerberg/happy/issues/2418)) ([77f4dc1](https://github.com/chanzuckerberg/happy/commit/77f4dc12ce44f88998949807b0ea3730699c77f8))


### Bug Fixes

* WAF attachment errors in happy-env-ecs (CCIE-1824) ([#2417](https://github.com/chanzuckerberg/happy/issues/2417)) ([f258119](https://github.com/chanzuckerberg/happy/commit/f258119d8c63cb0a3666ea295847161d742e760a))

## [3.1.0](https://github.com/chanzuckerberg/happy/compare/happy-env-ecs-v3.0.0...happy-env-ecs-v3.1.0) (2023-08-29)


### Features

* Support for ECR tag immutability ([#2376](https://github.com/chanzuckerberg/happy/issues/2376)) ([c1d5f5b](https://github.com/chanzuckerberg/happy/commit/c1d5f5b6e6a093c19ba2a092111842cc0e4f195f))

## [3.0.0](https://github.com/chanzuckerberg/happy/compare/happy-env-ecs-v2.2.8...happy-env-ecs-v3.0.0) (2023-08-16)


### ⚠ BREAKING CHANGES

* Arguments missing from happy_github_ci_role ([#2273](https://github.com/chanzuckerberg/happy/issues/2273))

### Bug Fixes

* Arguments missing from happy_github_ci_role ([#2273](https://github.com/chanzuckerberg/happy/issues/2273)) ([e73b096](https://github.com/chanzuckerberg/happy/commit/e73b0964bc03ac208e026f97b1f6f0536b1a5d49))

## [2.2.8](https://github.com/chanzuckerberg/happy/compare/happy-env-ecs-v2.2.7...happy-env-ecs-v2.2.8) (2023-05-30)


### Bug Fixes

* aws provider 5.0 deprecated source_json ([#1810](https://github.com/chanzuckerberg/happy/issues/1810)) ([7b69d30](https://github.com/chanzuckerberg/happy/commit/7b69d3086112972c5792edf31509dc1bde4ba23b))
* Handle empty and null ecr policies ([#1813](https://github.com/chanzuckerberg/happy/issues/1813)) ([b2e60f1](https://github.com/chanzuckerberg/happy/commit/b2e60f1dcb948a1cc3ec860c26b3ed541112b5de))

## [2.2.7](https://github.com/chanzuckerberg/happy/compare/happy-env-ecs-v2.2.6...happy-env-ecs-v2.2.7) (2023-04-25)


### Bug Fixes

* broken integration secret ([#1655](https://github.com/chanzuckerberg/happy/issues/1655)) ([64ae962](https://github.com/chanzuckerberg/happy/commit/64ae962f99a3f69288fcd75d4ab501afab04c494))

## [2.2.6](https://github.com/chanzuckerberg/happy/compare/happy-env-ecs-v2.2.5...happy-env-ecs-v2.2.6) (2023-04-25)


### Bug Fixes

* type mismatch for additional variables ([#1646](https://github.com/chanzuckerberg/happy/issues/1646)) ([3aca7d0](https://github.com/chanzuckerberg/happy/commit/3aca7d07bb35a9db89a78dc664a9833399aee43b))

## [2.2.5](https://github.com/chanzuckerberg/happy/compare/happy-env-ecs-v2.2.4...happy-env-ecs-v2.2.5) (2023-04-24)


### Bug Fixes

* needs_private_waf_attachment and needs_public_waf_attachment ([#1608](https://github.com/chanzuckerberg/happy/issues/1608)) ([50a8b98](https://github.com/chanzuckerberg/happy/commit/50a8b9866c7d1bfcedbb0705d3f9d9d465129ed1))

## [2.2.4](https://github.com/chanzuckerberg/happy/compare/happy-env-ecs-v2.2.3...happy-env-ecs-v2.2.4) (2023-04-20)


### Bug Fixes

* remove skip; actually validate child modules ([#1602](https://github.com/chanzuckerberg/happy/issues/1602)) ([79c6719](https://github.com/chanzuckerberg/happy/commit/79c671919e4fa897c93d441fa60825694f65b1ce))

## [2.2.3](https://github.com/chanzuckerberg/happy/compare/happy-env-ecs-v2.2.2...happy-env-ecs-v2.2.3) (2023-04-20)


### Bug Fixes

* Disable terraform validate on select modules ([#1600](https://github.com/chanzuckerberg/happy/issues/1600)) ([0294798](https://github.com/chanzuckerberg/happy/commit/0294798010874c57e601c4f78f0a4efd899796a8))

## [2.2.2](https://github.com/chanzuckerberg/happy/compare/happy-env-ecs-v2.2.1...happy-env-ecs-v2.2.2) (2023-04-20)


### Bug Fixes

* support init_script in terraform/modules/happy-env-ecs/batch.tf ([#1596](https://github.com/chanzuckerberg/happy/issues/1596)) ([202d3b9](https://github.com/chanzuckerberg/happy/commit/202d3b9b835f8178ef02eab866644f06c9c4d4a9))

## [2.2.1](https://github.com/chanzuckerberg/happy/compare/happy-env-ecs-v2.2.0...happy-env-ecs-v2.2.1) (2023-04-14)


### Bug Fixes

* When no WAF ARN is present, change the iterator variable type ([#1562](https://github.com/chanzuckerberg/happy/issues/1562)) ([15a982a](https://github.com/chanzuckerberg/happy/commit/15a982aee1d828a2761edd44c9aa5ba0e59d6ac9))

## [2.2.0](https://github.com/chanzuckerberg/happy/compare/happy-env-ecs-v2.1.0...happy-env-ecs-v2.2.0) (2023-03-15)


### Features

* automate adding OIDC providers for new happy apps to happy api ([#1353](https://github.com/chanzuckerberg/happy/issues/1353)) ([782a650](https://github.com/chanzuckerberg/happy/commit/782a650aa6366d7b8f27d94642c0bb21fd99c10c))

## [2.1.0](https://github.com/chanzuckerberg/happy/compare/happy-env-ecs-v2.0.0...happy-env-ecs-v2.1.0) (2023-03-07)


### Features

* give the modules an option to configure a Web ACL to protect its endpoints ([#1275](https://github.com/chanzuckerberg/happy/issues/1275)) ([90dae59](https://github.com/chanzuckerberg/happy/commit/90dae59595b041d24765123ca56c85021fe46cdb))


### Bug Fixes

* WAF assignment null condition ([#1301](https://github.com/chanzuckerberg/happy/issues/1301)) ([7ce142e](https://github.com/chanzuckerberg/happy/commit/7ce142ead96e012a192901fa5529ed6a0c2cb7bc))

## [2.0.0](https://github.com/chanzuckerberg/happy/compare/happy-env-ecs-v1.5.1...happy-env-ecs-v2.0.0) (2023-02-13)


### ⚠ BREAKING CHANGES

* inject happy config to stacks ([#1108](https://github.com/chanzuckerberg/happy/issues/1108))

### Features

* inject happy config to stacks ([#1108](https://github.com/chanzuckerberg/happy/issues/1108)) ([9cb49c7](https://github.com/chanzuckerberg/happy/commit/9cb49c7f7bd6819541510e4f31ab5fd112579457))

## [1.5.1](https://github.com/chanzuckerberg/happy/compare/happy-env-ecs-v1.5.0...happy-env-ecs-v1.5.1) (2023-01-31)


### Bug Fixes

* bump version of proxy ([#979](https://github.com/chanzuckerberg/happy/issues/979)) ([2af63ce](https://github.com/chanzuckerberg/happy/commit/2af63ced8c26eb2b74da8eb421e8d8af76194d95))

## [1.5.0](https://github.com/chanzuckerberg/happy/compare/happy-env-ecs-v1.4.0...happy-env-ecs-v1.5.0) (2023-01-31)


### Features

* Tag stack level metrics for EKS and ECS ([#1033](https://github.com/chanzuckerberg/happy/issues/1033)) ([1466430](https://github.com/chanzuckerberg/happy/commit/146643014a9c60cf2bac67fd25d6881827b9b3e9))

## [1.4.0](https://github.com/chanzuckerberg/happy/compare/happy-env-ecs-v1.3.3...happy-env-ecs-v1.4.0) (2023-01-27)


### Features

* add synthetics to ecs stacks module ([#1008](https://github.com/chanzuckerberg/happy/issues/1008)) ([7ad6192](https://github.com/chanzuckerberg/happy/commit/7ad6192edf208908b50ec8ff906994fef4a15829))

## [1.3.3](https://github.com/chanzuckerberg/happy/compare/happy-env-ecs-v1.3.2...happy-env-ecs-v1.3.3) (2023-01-19)


### Bug Fixes

* attach dynamo locktable policy to github ci role ([#978](https://github.com/chanzuckerberg/happy/issues/978)) ([f9fe4d6](https://github.com/chanzuckerberg/happy/commit/f9fe4d6b40d5fd0e7e2ce11384819f704b8ad2af))

## [1.3.2](https://github.com/chanzuckerberg/happy/compare/happy-env-ecs-v1.3.1...happy-env-ecs-v1.3.2) (2023-01-09)


### Bug Fixes

* happy env ecs ([#960](https://github.com/chanzuckerberg/happy/issues/960)) ([323a6cc](https://github.com/chanzuckerberg/happy/commit/323a6cc0796056076f0c3c4ba75e3bd055232a5f))

## [1.3.1](https://github.com/chanzuckerberg/happy/compare/happy-env-ecs-v1.3.0...happy-env-ecs-v1.3.1) (2023-01-09)


### Bug Fixes

* ecs-multi-fix same port ([#956](https://github.com/chanzuckerberg/happy/issues/956)) ([36e7697](https://github.com/chanzuckerberg/happy/commit/36e7697e1d15f5a306ac9e0c7259117ad8fdb727))

## [1.3.0](https://github.com/chanzuckerberg/happy/compare/happy-env-ecs-v1.2.2...happy-env-ecs-v1.3.0) (2023-01-09)


### Features

* metrics server ([#954](https://github.com/chanzuckerberg/happy/issues/954)) ([3e4011d](https://github.com/chanzuckerberg/happy/commit/3e4011d8db8700650d49a24cc255734ee1c6c46c))

## [1.2.2](https://github.com/chanzuckerberg/happy/compare/happy-env-ecs-v1.2.1...happy-env-ecs-v1.2.2) (2022-12-13)


### Bug Fixes

* make var iterable ([#868](https://github.com/chanzuckerberg/happy/issues/868)) ([8d08ee7](https://github.com/chanzuckerberg/happy/commit/8d08ee7c37899931633d51fa317637094bac766e))

## [1.2.1](https://github.com/chanzuckerberg/happy/compare/happy-env-ecs-v1.2.0...happy-env-ecs-v1.2.1) (2022-12-12)


### Bug Fixes

* move required_version ([#855](https://github.com/chanzuckerberg/happy/issues/855)) ([b13832c](https://github.com/chanzuckerberg/happy/commit/b13832ca61af7ed8ca0caec643e24bd8633ea4c1))

## [1.2.0](https://github.com/chanzuckerberg/happy/compare/happy-env-ecs-v1.1.0...happy-env-ecs-v1.2.0) (2022-12-12)


### Features

* terraform cicd ([#847](https://github.com/chanzuckerberg/happy/issues/847)) ([1be9354](https://github.com/chanzuckerberg/happy/commit/1be9354192ce8085fa967c0c9280a772a4bb6daa))

## [1.1.0](https://github.com/chanzuckerberg/happy/compare/happy-env-ecs-v1.0.0...happy-env-ecs-v1.1.0) (2022-12-07)


### Features

* add a service account to pods ([#835](https://github.com/chanzuckerberg/happy/issues/835)) ([203c129](https://github.com/chanzuckerberg/happy/commit/203c1294602160dfc4aacc15adf8ebc91e83af5a))

## 1.0.0 (2022-11-30)


### Features

* terraform modules for happy sharing ([#800](https://github.com/chanzuckerberg/happy/issues/800)) ([d909860](https://github.com/chanzuckerberg/happy/commit/d9098607e37b29c71bdc3ddac9fabd7ba280606b))
