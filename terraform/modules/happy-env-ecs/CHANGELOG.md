# Changelog

## [6.4.1](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v6.4.0...happy-env-ecs-v6.4.1) (2022-11-23)


### Bug Fixes

* move ecs stuff out to eks can use module ([#6768](https://github.com/chanzuckerberg/shared-infra/issues/6768)) ([53053fa](https://github.com/chanzuckerberg/shared-infra/commit/53053fa9d0c1cf333838d8c49625309a5f43e4b3))

## [6.4.0](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v6.3.2...happy-env-ecs-v6.4.0) (2022-11-17)


### Features

* pin the new module version ci role + release ([#6735](https://github.com/chanzuckerberg/shared-infra/issues/6735)) ([78f0e78](https://github.com/chanzuckerberg/shared-infra/commit/78f0e782203de750a62375ae826b3e4ef514bb27))

## [6.3.2](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v6.3.1...happy-env-ecs-v6.3.2) (2022-11-16)


### Bug Fixes

* de-conflict ecs reader policy name ([#6726](https://github.com/chanzuckerberg/shared-infra/issues/6726)) ([3bfae67](https://github.com/chanzuckerberg/shared-infra/commit/3bfae67bdd31349fdd52ec20e3d64e30f66a9b48))

## [6.3.1](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v6.3.0...happy-env-ecs-v6.3.1) (2022-11-15)


### Bug Fixes

* can't use count in a tf plan when resource does not exist yet ([#6709](https://github.com/chanzuckerberg/shared-infra/issues/6709)) ([0348160](https://github.com/chanzuckerberg/shared-infra/commit/0348160ed0c6dd711d6d2eca7d1bdfbd3be4ef08))

## [6.3.0](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v6.2.0...happy-env-ecs-v6.3.0) (2022-11-04)


### Features

* add database configuration to integration secret ([#6644](https://github.com/chanzuckerberg/shared-infra/issues/6644)) ([0c49b90](https://github.com/chanzuckerberg/shared-infra/commit/0c49b90f27757e03d11374ede2d22eae2fd5defc))


### Bug Fixes

* Tweak ssm_reader_writer policy name in happy-github-ci-role ([#6641](https://github.com/chanzuckerberg/shared-infra/issues/6641)) ([e132837](https://github.com/chanzuckerberg/shared-infra/commit/e132837c2b2b2ab205b22b341b256c33c27720fd))

## [6.2.0](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v6.1.2...happy-env-ecs-v6.2.0) (2022-11-03)


### Features

* add CI role creation to happy-env-* modules ([#6622](https://github.com/chanzuckerberg/shared-infra/issues/6622)) ([694f6f7](https://github.com/chanzuckerberg/shared-infra/commit/694f6f751f1709346848c779ad450cc6a4d3fba7))

## [6.1.2](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v6.1.1...happy-env-ecs-v6.1.2) (2022-10-03)


### Bug Fixes

* **CCIE-719:** czi-lp-poc happy shell errors out ([#6365](https://github.com/chanzuckerberg/shared-infra/issues/6365)) ([d6c390c](https://github.com/chanzuckerberg/shared-infra/commit/d6c390c3d5871363460447cb068613cb834aa327))

## [6.1.1](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v6.1.0...happy-env-ecs-v6.1.1) (2022-09-27)


### Bug Fixes

* Add vpc-id to the happy-env-eks integration secret ([#6289](https://github.com/chanzuckerberg/shared-infra/issues/6289)) ([586fbd7](https://github.com/chanzuckerberg/shared-infra/commit/586fbd7dea6eeda0dc341223d47e790ac1882f9e))

## [6.1.0](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v6.0.3...happy-env-ecs-v6.1.0) (2022-09-19)


### Features

* update proxy to trunk lengths of services ([#6278](https://github.com/chanzuckerberg/shared-infra/issues/6278)) ([ee85afc](https://github.com/chanzuckerberg/shared-infra/commit/ee85afce76743006213f62b2064d4563e4498916))

## [6.0.3](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v6.0.2...happy-env-ecs-v6.0.3) (2022-09-07)


### Bug Fixes

* Upgrade TLS policies 1.1 to lowest current TLS 1.2 policy ([#6208](https://github.com/chanzuckerberg/shared-infra/issues/6208)) ([89812af](https://github.com/chanzuckerberg/shared-infra/commit/89812af9d3c44c91d2e5942fd54eac5be3a35366))

## [6.0.2](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v6.0.1...happy-env-ecs-v6.0.2) (2022-09-07)


### Bug Fixes

* Upgrade TLS policies 1.1 to lowest current TLS 1.2 policy ([#6208](https://github.com/chanzuckerberg/shared-infra/issues/6208)) ([89812af](https://github.com/chanzuckerberg/shared-infra/commit/89812af9d3c44c91d2e5942fd54eac5be3a35366))

## [6.0.1](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v6.0.0...happy-env-ecs-v6.0.1) (2022-09-01)


### Bug Fixes

* namespacing in tfe-agents ([#6157](https://github.com/chanzuckerberg/shared-infra/issues/6157)) ([7e59401](https://github.com/chanzuckerberg/shared-infra/commit/7e59401116599fa04a5503d151e13128dabee025))

## [6.0.0](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v5.0.9...happy-env-ecs-v6.0.0) (2022-08-31)


### ⚠ BREAKING CHANGES

* [CDI-223] Allow each happy service to configure own LB priv/pub (#5678)
* (happy-env-ecs) Set Default Aurora RDS engine_version to latest (14.x); also parameterize it (#5632)

### Features

* (happy-env-ecs) Set Default Aurora RDS engine_version to latest (14.x); also parameterize it ([#5632](https://github.com/chanzuckerberg/shared-infra/issues/5632)) ([13905df](https://github.com/chanzuckerberg/shared-infra/commit/13905df0003b03fba541578c09309adc758cc93c))
* [CDI-223] Allow each happy service to configure own LB priv/pub ([#5678](https://github.com/chanzuckerberg/shared-infra/issues/5678)) ([027dbf6](https://github.com/chanzuckerberg/shared-infra/commit/027dbf6e08f30083670a15e2f8a6726e8b66c58d))


### Bug Fixes

* add databricks to dev-cutter providers ([#5900](https://github.com/chanzuckerberg/shared-infra/issues/5900)) ([fe7d0cc](https://github.com/chanzuckerberg/shared-infra/commit/fe7d0ccea6a694728e7d7aab4a35a08f630d02d9))
* all for_each in happy use the local services ([#5757](https://github.com/chanzuckerberg/shared-infra/issues/5757)) ([d0506e5](https://github.com/chanzuckerberg/shared-infra/commit/d0506e54968d4197cd3d7491235aa82037eef9b3))
* happy public sg for_each needs map not list ([#5765](https://github.com/chanzuckerberg/shared-infra/issues/5765)) ([043d450](https://github.com/chanzuckerberg/shared-infra/commit/043d45048bcc6e87bd21ed777046041c5ceb56ca))
* private_lb_services is a var in happy not local ([#5754](https://github.com/chanzuckerberg/shared-infra/issues/5754)) ([0797575](https://github.com/chanzuckerberg/shared-infra/commit/0797575a0796af84d102edfb62da24c078f454b2))
* public_lb_services is a var ([#5708](https://github.com/chanzuckerberg/shared-infra/issues/5708)) ([5873989](https://github.com/chanzuckerberg/shared-infra/commit/587398917a7728480aa7aae06d255b7f8235b585))
* remove deprecated grpc alb from happy secret ([#5707](https://github.com/chanzuckerberg/shared-infra/issues/5707)) ([23f0969](https://github.com/chanzuckerberg/shared-infra/commit/23f096952a8a3377473be640a8c543ef5e533783))
* remove public from pubic services ([#5932](https://github.com/chanzuckerberg/shared-infra/issues/5932)) ([221f8da](https://github.com/chanzuckerberg/shared-infra/commit/221f8da80a752d9eb7c1e68fa315aef29347f543))
* remove unused ssh key name variable ([#6145](https://github.com/chanzuckerberg/shared-infra/issues/6145)) ([75983f4](https://github.com/chanzuckerberg/shared-infra/commit/75983f4c0efd36e6c61cb31d93e91788e36eb187))
* straggling oauth proxy references ([#5687](https://github.com/chanzuckerberg/shared-infra/issues/5687)) ([ce32909](https://github.com/chanzuckerberg/shared-infra/commit/ce32909a457b40f00ca4f39678b7541f022a995e))
* use string list not map in services' for each ([#5710](https://github.com/chanzuckerberg/shared-infra/issues/5710)) ([0636172](https://github.com/chanzuckerberg/shared-infra/commit/063617280f7085c83434e5adbda96c3691034e6d))

## [5.0.9](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v5.0.8...happy-env-ecs-v5.0.9) (2022-08-31)


### Bug Fixes

* remove unused ssh key name variable ([#6145](https://github.com/chanzuckerberg/shared-infra/issues/6145)) ([75983f4](https://github.com/chanzuckerberg/shared-infra/commit/75983f4c0efd36e6c61cb31d93e91788e36eb187))

## [5.0.8](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v5.0.7...happy-env-ecs-v5.0.8) (2022-08-09)


### Bug Fixes

* add databricks to dev-cutter providers ([#5900](https://github.com/chanzuckerberg/shared-infra/issues/5900)) ([fe7d0cc](https://github.com/chanzuckerberg/shared-infra/commit/fe7d0ccea6a694728e7d7aab4a35a08f630d02d9))
* remove public from pubic services ([#5932](https://github.com/chanzuckerberg/shared-infra/issues/5932)) ([221f8da](https://github.com/chanzuckerberg/shared-infra/commit/221f8da80a752d9eb7c1e68fa315aef29347f543))

## [5.0.7](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v5.0.6...happy-env-ecs-v5.0.7) (2022-07-25)


### Bug Fixes

* happy public sg for_each needs map not list ([#5765](https://github.com/chanzuckerberg/shared-infra/issues/5765)) ([043d450](https://github.com/chanzuckerberg/shared-infra/commit/043d45048bcc6e87bd21ed777046041c5ceb56ca))

## [5.0.6](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v5.0.5...happy-env-ecs-v5.0.6) (2022-07-25)


### Bug Fixes

* all for_each in happy use the local services ([#5757](https://github.com/chanzuckerberg/shared-infra/issues/5757)) ([d0506e5](https://github.com/chanzuckerberg/shared-infra/commit/d0506e54968d4197cd3d7491235aa82037eef9b3))

## [5.0.5](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v5.0.4...happy-env-ecs-v5.0.5) (2022-07-25)


### Bug Fixes

* private_lb_services is a var in happy not local ([#5754](https://github.com/chanzuckerberg/shared-infra/issues/5754)) ([0797575](https://github.com/chanzuckerberg/shared-infra/commit/0797575a0796af84d102edfb62da24c078f454b2))

## [5.0.4](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v5.0.3...happy-env-ecs-v5.0.4) (2022-07-21)


### Bug Fixes

* use string list not map in services' for each ([#5710](https://github.com/chanzuckerberg/shared-infra/issues/5710)) ([0636172](https://github.com/chanzuckerberg/shared-infra/commit/063617280f7085c83434e5adbda96c3691034e6d))

## [5.0.3](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v5.0.2...happy-env-ecs-v5.0.3) (2022-07-16)


### Bug Fixes

* public_lb_services is a var ([#5708](https://github.com/chanzuckerberg/shared-infra/issues/5708)) ([5873989](https://github.com/chanzuckerberg/shared-infra/commit/587398917a7728480aa7aae06d255b7f8235b585))

## [5.0.2](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v5.0.1...happy-env-ecs-v5.0.2) (2022-07-15)


### Bug Fixes

* remove deprecated grpc alb from happy secret ([#5707](https://github.com/chanzuckerberg/shared-infra/issues/5707)) ([23f0969](https://github.com/chanzuckerberg/shared-infra/commit/23f096952a8a3377473be640a8c543ef5e533783))

## [5.0.1](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v5.0.0...happy-env-ecs-v5.0.1) (2022-07-15)


### Bug Fixes

* straggling oauth proxy references ([#5687](https://github.com/chanzuckerberg/shared-infra/issues/5687)) ([ce32909](https://github.com/chanzuckerberg/shared-infra/commit/ce32909a457b40f00ca4f39678b7541f022a995e))

## [5.0.0](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v4.0.0...happy-env-ecs-v5.0.0) (2022-07-14)


### ⚠ BREAKING CHANGES

* [CDI-223] Allow each happy service to configure own LB priv/pub (#5678)

### Features

* [CDI-223] Allow each happy service to configure own LB priv/pub ([#5678](https://github.com/chanzuckerberg/shared-infra/issues/5678)) ([027dbf6](https://github.com/chanzuckerberg/shared-infra/commit/027dbf6e08f30083670a15e2f8a6726e8b66c58d))

## [4.0.0](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v3.5.3...happy-env-ecs-v4.0.0) (2022-07-12)


### ⚠ BREAKING CHANGES

* (happy-env-ecs) Set Default Aurora RDS engine_version to latest (14.x); also parameterize it (#5632)

### Features

* (happy-env-ecs) Set Default Aurora RDS engine_version to latest (14.x); also parameterize it ([#5632](https://github.com/chanzuckerberg/shared-infra/issues/5632)) ([13905df](https://github.com/chanzuckerberg/shared-infra/commit/13905df0003b03fba541578c09309adc758cc93c))

## [3.5.3](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v3.5.2...happy-env-ecs-v3.5.3) (2022-06-23)


### Bug Fixes

* bug in happy-env-ecs for gzip ([#5550](https://github.com/chanzuckerberg/shared-infra/issues/5550)) ([b234911](https://github.com/chanzuckerberg/shared-infra/commit/b23491118df5b31a6d4b61d4b91505c8ba79941f))

## [3.5.2](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v3.5.1...happy-env-ecs-v3.5.2) (2022-06-23)


### Bug Fixes

* broken happy-env-ecs module ([#5546](https://github.com/chanzuckerberg/shared-infra/issues/5546)) ([f04074d](https://github.com/chanzuckerberg/shared-infra/commit/f04074d944f187bccc2bda8c1dc42b78c6229ae3))

## [3.5.1](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v3.5.0...happy-env-ecs-v3.5.1) (2022-06-23)


### Bug Fixes

* swipe arguments in new moduel ([#5545](https://github.com/chanzuckerberg/shared-infra/issues/5545)) ([5653f7b](https://github.com/chanzuckerberg/shared-infra/commit/5653f7b8934b130e10da2a2b87d558021f161c15))

## [3.5.0](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v3.4.0...happy-env-ecs-v3.5.0) (2022-06-23)


### Features

* allow happy-env-ecs to pass oauth2-proxy args ([#5543](https://github.com/chanzuckerberg/shared-infra/issues/5543)) ([07a3b03](https://github.com/chanzuckerberg/shared-infra/commit/07a3b03f2581db138c3d60f6b91851796ae49e18))

## [3.4.0](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v3.3.0...happy-env-ecs-v3.4.0) (2022-06-06)


### Features

* Enable ssh for swipe [CCIE-5] ([#5153](https://github.com/chanzuckerberg/shared-infra/issues/5153)) ([3d57002](https://github.com/chanzuckerberg/shared-infra/commit/3d570027c795d79d53e6734d279a286bf976d13b))

## [3.3.0](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v3.2.7...happy-env-ecs-v3.3.0) (2022-05-23)


### Features

* add dynamo locktable to happy-env-ecs module ([#5318](https://github.com/chanzuckerberg/shared-infra/issues/5318)) ([dc4fec1](https://github.com/chanzuckerberg/shared-infra/commit/dc4fec189cda1cb94fba511bd87292f94b722ac2))

### [3.2.7](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v3.2.6...happy-env-ecs-v3.2.7) (2022-05-20)


### Bug Fixes

* add local prefix ([#5333](https://github.com/chanzuckerberg/shared-infra/issues/5333)) ([72225a0](https://github.com/chanzuckerberg/shared-infra/commit/72225a0e1c34550bee6460a19a49796c5707daab))

### [3.2.6](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v3.2.5...happy-env-ecs-v3.2.6) (2022-05-20)


### Bug Fixes

* generate map of happy private services from set ([#5331](https://github.com/chanzuckerberg/shared-infra/issues/5331)) ([412887c](https://github.com/chanzuckerberg/shared-infra/commit/412887cfd896c4a69394217068b19c4aa877bc5a))

### [3.2.5](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v3.2.4...happy-env-ecs-v3.2.5) (2022-05-20)


### Bug Fixes

* happy conditional outputs were inconsistent ([#5329](https://github.com/chanzuckerberg/shared-infra/issues/5329)) ([d6266a4](https://github.com/chanzuckerberg/shared-infra/commit/d6266a43d4f5a1c1c1766e46424f27a08120658e))

### [3.2.4](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v3.2.3...happy-env-ecs-v3.2.4) (2022-05-20)


### Bug Fixes

* pull keys from services ([#5327](https://github.com/chanzuckerberg/shared-infra/issues/5327)) ([927aef2](https://github.com/chanzuckerberg/shared-infra/commit/927aef2b2ae2993d1ae712b05a119ba144efeb37))

### [3.2.3](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v3.2.2...happy-env-ecs-v3.2.3) (2022-05-20)


### Bug Fixes

* toset not set in happy private lb ([#5325](https://github.com/chanzuckerberg/shared-infra/issues/5325)) ([20330ca](https://github.com/chanzuckerberg/shared-infra/commit/20330cac4fbd5678db1eb9bc602fd05a2b4df4f8))

### [3.2.2](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v3.2.1...happy-env-ecs-v3.2.2) (2022-05-19)


### Bug Fixes

* grpc alb listener variable name changed ([#5316](https://github.com/chanzuckerberg/shared-infra/issues/5316)) ([67f4a31](https://github.com/chanzuckerberg/shared-infra/commit/67f4a31273088a360be12a895ee2583e95aeba7c))

### [3.2.1](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v3.2.0...happy-env-ecs-v3.2.1) (2022-05-18)


### Bug Fixes

* use setsubtract in happy albs ([#5313](https://github.com/chanzuckerberg/shared-infra/issues/5313)) ([43bd713](https://github.com/chanzuckerberg/shared-infra/commit/43bd71359d90658c8dc3ad6872aad58556c0e584))

## [3.2.0](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v3.1.5...happy-env-ecs-v3.2.0) (2022-05-18)


### Features

* add http redirect for happy grpc; remove standard alb if grpc enabled ([#5311](https://github.com/chanzuckerberg/shared-infra/issues/5311)) ([988be7a](https://github.com/chanzuckerberg/shared-infra/commit/988be7a68b65473db13d1adfd71b6699e782dc32))
* Enable deletion protection for aurora postgres dbs ([#5287](https://github.com/chanzuckerberg/shared-infra/issues/5287)) ([0709c78](https://github.com/chanzuckerberg/shared-infra/commit/0709c78ac64af6d7d0bdf2e52af786f7013cc6ac))

### [3.1.5](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v3.1.4...happy-env-ecs-v3.1.5) (2022-05-09)


### Bug Fixes

* CCIE-40 happy Remove unused ALB default listeners ([#5136](https://github.com/chanzuckerberg/shared-infra/issues/5136)) ([f8ff6e8](https://github.com/chanzuckerberg/shared-infra/commit/f8ff6e828d28969ad98fb6b11ac226a5939160c7))

### [3.1.4](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v3.1.3...happy-env-ecs-v3.1.4) (2022-05-02)


### Bug Fixes

* Add cert for happy grpc lb ([#5139](https://github.com/chanzuckerberg/shared-infra/issues/5139)) ([fbeb4bc](https://github.com/chanzuckerberg/shared-infra/commit/fbeb4bc7f539c27ef9aaac3edd68bd52282829e3))

### [3.1.3](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v3.1.2...happy-env-ecs-v3.1.3) (2022-05-02)


### Bug Fixes

* [CDI-75] No idle_timeout in new gRPC set ([#5135](https://github.com/chanzuckerberg/shared-infra/issues/5135)) ([1422d20](https://github.com/chanzuckerberg/shared-infra/commit/1422d205f724325b9f227819fc8a0a41279346b0))

### [3.1.2](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v3.1.1...happy-env-ecs-v3.1.2) (2022-04-29)


### Bug Fixes

* grpc list must be a set ([#5129](https://github.com/chanzuckerberg/shared-infra/issues/5129)) ([f15e299](https://github.com/chanzuckerberg/shared-infra/commit/f15e299840171be3a8482a7cb4f610b31829e7e4))

### [3.1.1](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v3.1.0...happy-env-ecs-v3.1.1) (2022-04-29)


### Bug Fixes

* missing lb references in happy ([#5125](https://github.com/chanzuckerberg/shared-infra/issues/5125)) ([b42ff44](https://github.com/chanzuckerberg/shared-infra/commit/b42ff446a791955539125f84325a384926a78f9c))

## [3.1.0](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v3.0.0...happy-env-ecs-v3.1.0) (2022-04-29)


### Features

* [CDI-75] Add Private gRPC load balancer to happy module ([#5122](https://github.com/chanzuckerberg/shared-infra/issues/5122)) ([b74975c](https://github.com/chanzuckerberg/shared-infra/commit/b74975c30c5d8829b37c1d3f7ff841bfa3e0ae53))

## [3.0.0](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v2.0.0...happy-env-ecs-v3.0.0) (2022-04-22)


### ⚠ BREAKING CHANGES

* happy-env-ecs/multidomai-proxy fixes namespace inconsistency; allow more than 1 per "project" (#5074)

### Bug Fixes

* happy-env-ecs/multidomai-proxy fixes namespace inconsistency; allow more than 1 per "project" ([#5074](https://github.com/chanzuckerberg/shared-infra/issues/5074)) ([2b6036c](https://github.com/chanzuckerberg/shared-infra/commit/2b6036ca2dc5df389ac44121c37028241b946e0d))

## [2.0.0](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v1.1.0...happy-env-ecs-v2.0.0) (2022-03-31)


### ⚠ BREAKING CHANGES

* AWS Batch Env has SSH+Bless enabled; make Happy modules consistent with cloud-env outputs (#4896)

### Features

* AWS Batch Env has SSH+Bless enabled; make Happy modules consistent with cloud-env outputs ([#4896](https://github.com/chanzuckerberg/shared-infra/issues/4896)) ([8fa76af](https://github.com/chanzuckerberg/shared-infra/commit/8fa76af0cead00c01ecde34055065a193dd562af))

## [1.1.0](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v1.0.1...happy-env-ecs-v1.1.0) (2022-02-16)


### Features

* Update swipe for happy integration ([#4812](https://github.com/chanzuckerberg/shared-infra/issues/4812)) ([65c1137](https://github.com/chanzuckerberg/shared-infra/commit/65c11371ee2e8a4889b68564f9e87ff08dace65a))

### [1.0.1](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-ecs-v1.0.0...happy-env-ecs-v1.0.1) (2022-01-24)


### Bug Fixes

* lifecycle rule for happy env ecs ([#4705](https://github.com/chanzuckerberg/shared-infra/issues/4705)) ([bc37ce9](https://github.com/chanzuckerberg/shared-infra/commit/bc37ce983156275c8249ad9c269ddd2969b1315e))

## 1.0.0 (2022-01-21)


### Features

* Add Module Swipe to happy-env-ecs ([#4675](https://github.com/chanzuckerberg/shared-infra/issues/4675)) ([4e2dad5](https://github.com/chanzuckerberg/shared-infra/commit/4e2dad5eaa2036cdacadda9735614009ed441dd9))
