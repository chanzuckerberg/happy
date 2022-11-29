# Changelog

## [0.2.1](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-eks-v0.2.0...happy-env-eks-v0.2.1) (2022-11-23)


### Bug Fixes

* move ecs stuff out to eks can use module ([#6768](https://github.com/chanzuckerberg/shared-infra/issues/6768)) ([53053fa](https://github.com/chanzuckerberg/shared-infra/commit/53053fa9d0c1cf333838d8c49625309a5f43e4b3))

## [0.2.0](https://github.com/chanzuckerberg/shared-infra/compare/happy-env-eks-v0.1.0...happy-env-eks-v0.2.0) (2022-11-04)


### Features

* add database configuration to integration secret ([#6644](https://github.com/chanzuckerberg/shared-infra/issues/6644)) ([0c49b90](https://github.com/chanzuckerberg/shared-infra/commit/0c49b90f27757e03d11374ede2d22eae2fd5defc))

## 0.1.0 (2022-11-03)


### Features

* add CI role creation to happy-env-* modules ([#6622](https://github.com/chanzuckerberg/shared-infra/issues/6622)) ([694f6f7](https://github.com/chanzuckerberg/shared-infra/commit/694f6f751f1709346848c779ad450cc6a4d3fba7))
* **CCIE-599:** Build out EKS, K8s-core in the happy-env-eks ([#6121](https://github.com/chanzuckerberg/shared-infra/issues/6121)) ([9791c34](https://github.com/chanzuckerberg/shared-infra/commit/9791c3419f35215e88a0e5dd41e48d24e53ddd20))
* Create an wildcard ACM cert and store the arn in integration secret ([#6300](https://github.com/chanzuckerberg/shared-infra/issues/6300)) ([a2d6573](https://github.com/chanzuckerberg/shared-infra/commit/a2d657305c9c8391a205ef6a782f246f8404f951))
* Expose oauth proxy service name through the integration secret ([#6562](https://github.com/chanzuckerberg/shared-infra/issues/6562)) ([64ed2b1](https://github.com/chanzuckerberg/shared-infra/commit/64ed2b1546c65700b7cef8d528a8d6334e923abc))
* Implement a non-externally accessible service in k8s-happy-app ([#6486](https://github.com/chanzuckerberg/shared-infra/issues/6486)) ([60578ea](https://github.com/chanzuckerberg/shared-infra/commit/60578ea20aa5d874d7cbd1d05bba983341dd39be))
* Implement oauth proxy (eks-multi-domain-oauth-proxy) in happy-env-eks ([#6321](https://github.com/chanzuckerberg/shared-infra/issues/6321)) ([3709aa2](https://github.com/chanzuckerberg/shared-infra/commit/3709aa28ef660f3849cafb33e7e61ee0c3caeb7b))
* refactor eks happy mult domain ([#6504](https://github.com/chanzuckerberg/shared-infra/issues/6504)) ([312ac56](https://github.com/chanzuckerberg/shared-infra/commit/312ac56f77b7f2a23496ead231bd4171b59af4f0))


### Bug Fixes

* Add tags to happy-eks integration secret ([#6320](https://github.com/chanzuckerberg/shared-infra/issues/6320)) ([0212345](https://github.com/chanzuckerberg/shared-infra/commit/02123450bc161c702b7281033bf2e8ea3cf97554))
* Add value for internal zone in happy-env-eks ([#6551](https://github.com/chanzuckerberg/shared-infra/issues/6551)) ([871ebdd](https://github.com/chanzuckerberg/shared-infra/commit/871ebdd1c12f7ffe273408b5cd7388bc84370dc5))
* Add vpc-id to the happy-env-eks integration secret ([#6289](https://github.com/chanzuckerberg/shared-infra/issues/6289)) ([586fbd7](https://github.com/chanzuckerberg/shared-infra/commit/586fbd7dea6eeda0dc341223d47e790ac1882f9e))
* Fix happy-eks module syntax ([#6566](https://github.com/chanzuckerberg/shared-infra/issues/6566)) ([b36ca85](https://github.com/chanzuckerberg/shared-infra/commit/b36ca85823f8deb7b228ccc67c10d320bf668807))
* namespacing in tfe-agents ([#6157](https://github.com/chanzuckerberg/shared-infra/issues/6157)) ([7e59401](https://github.com/chanzuckerberg/shared-infra/commit/7e59401116599fa04a5503d151e13128dabee025))
* Only one level of wildcard is supported ([#6309](https://github.com/chanzuckerberg/shared-infra/issues/6309)) ([fdee273](https://github.com/chanzuckerberg/shared-infra/commit/fdee273ce9920f2dcf96fdc44133cdeaa097b82d))
* Remove secret manager references from happy-env-eks ([#6281](https://github.com/chanzuckerberg/shared-infra/issues/6281)) ([4a6d482](https://github.com/chanzuckerberg/shared-infra/commit/4a6d482dbf143a326690b23f7ff6850d708af7e5))
