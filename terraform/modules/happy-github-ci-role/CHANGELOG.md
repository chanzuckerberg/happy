# Changelog

## [1.4.2](https://github.com/chanzuckerberg/happy/compare/happy-github-ci-role-v1.4.1...happy-github-ci-role-v1.4.2) (2023-09-22)


### Bug Fixes

* Add permissions to github ci role to allow retrieval of ecr image scanning ([#2492](https://github.com/chanzuckerberg/happy/issues/2492)) ([a6b7116](https://github.com/chanzuckerberg/happy/commit/a6b71169993165a55b8f139225a53e5a229367e9))

## [1.4.1](https://github.com/chanzuckerberg/happy/compare/happy-github-ci-role-v1.4.0...happy-github-ci-role-v1.4.1) (2023-09-21)


### Bug Fixes

* Fix ECR scanning when scanning is not enabled ([#2483](https://github.com/chanzuckerberg/happy/issues/2483)) ([9506729](https://github.com/chanzuckerberg/happy/commit/9506729d6121989b90fe58708b8bd07530e3bc0c))

## [1.4.0](https://github.com/chanzuckerberg/happy/compare/happy-github-ci-role-v1.3.3...happy-github-ci-role-v1.4.0) (2023-08-29)


### Features

* Support for ECR tag immutability ([#2376](https://github.com/chanzuckerberg/happy/issues/2376)) ([c1d5f5b](https://github.com/chanzuckerberg/happy/commit/c1d5f5b6e6a093c19ba2a092111842cc0e4f195f))

## [1.3.3](https://github.com/chanzuckerberg/happy/compare/happy-github-ci-role-v1.3.2...happy-github-ci-role-v1.3.3) (2023-05-09)


### Bug Fixes

* always make dynamo policy ([#1723](https://github.com/chanzuckerberg/happy/issues/1723)) ([2a3b43f](https://github.com/chanzuckerberg/happy/commit/2a3b43f9e5de9f93be400d67d966a356df50f7f3))

## [1.3.2](https://github.com/chanzuckerberg/happy/compare/happy-github-ci-role-v1.3.1...happy-github-ci-role-v1.3.2) (2023-05-05)


### Bug Fixes

* wrong count argument ([#1689](https://github.com/chanzuckerberg/happy/issues/1689)) ([45b2709](https://github.com/chanzuckerberg/happy/commit/45b27099d6504d16f58789e81bca7d7ef1c7e2b0))

## [1.3.1](https://github.com/chanzuckerberg/happy/compare/happy-github-ci-role-v1.3.0...happy-github-ci-role-v1.3.1) (2023-05-03)


### Bug Fixes

* dynamo alphanumeric sid ([#1687](https://github.com/chanzuckerberg/happy/issues/1687)) ([61afbda](https://github.com/chanzuckerberg/happy/commit/61afbdac2796d213b2b722c5e4f42044e00cfe48))

## [1.3.0](https://github.com/chanzuckerberg/happy/compare/happy-github-ci-role-v1.2.4...happy-github-ci-role-v1.3.0) (2023-05-02)


### Features

* add command to see the configured CI roles for env ([#1686](https://github.com/chanzuckerberg/happy/issues/1686)) ([a249cc0](https://github.com/chanzuckerberg/happy/commit/a249cc0a4fc61af413312b300f1fc4695529ee2e))


### Bug Fixes

* thread CI role through happy env ([#1683](https://github.com/chanzuckerberg/happy/issues/1683)) ([78b4be9](https://github.com/chanzuckerberg/happy/commit/78b4be95b7f4f4be95cf18a3d3b9920a28f409da))

## [1.2.4](https://github.com/chanzuckerberg/happy/compare/happy-github-ci-role-v1.2.3...happy-github-ci-role-v1.2.4) (2023-04-27)


### Bug Fixes

* count being used on resources known after apply ([#1665](https://github.com/chanzuckerberg/happy/issues/1665)) ([00eec6d](https://github.com/chanzuckerberg/happy/commit/00eec6d86b489408c2347ff57179d5ad9de43414))

## [1.2.3](https://github.com/chanzuckerberg/happy/compare/happy-github-ci-role-v1.2.2...happy-github-ci-role-v1.2.3) (2023-04-25)


### Bug Fixes

* move the iam role to opensource cztack ([#1651](https://github.com/chanzuckerberg/happy/issues/1651)) ([a490871](https://github.com/chanzuckerberg/happy/commit/a490871da60a4c2c672f02a78278298bef53fc06))

## [1.2.2](https://github.com/chanzuckerberg/happy/compare/happy-github-ci-role-v1.2.1...happy-github-ci-role-v1.2.2) (2023-04-18)


### Bug Fixes

* the CI role with latest permissions ([#1497](https://github.com/chanzuckerberg/happy/issues/1497)) ([a856f6c](https://github.com/chanzuckerberg/happy/commit/a856f6ce50b661e227db7d26e4943f82da37bab0))

## [1.2.1](https://github.com/chanzuckerberg/happy/compare/happy-github-ci-role-v1.2.0...happy-github-ci-role-v1.2.1) (2023-04-05)


### Bug Fixes

* don't make policy if no ecrs ([#1473](https://github.com/chanzuckerberg/happy/issues/1473)) ([1317dd1](https://github.com/chanzuckerberg/happy/commit/1317dd167d5ef5c28fce0f0fd2721951a7e1ed5b))

## [1.2.0](https://github.com/chanzuckerberg/happy/compare/happy-github-ci-role-v1.1.0...happy-github-ci-role-v1.2.0) (2022-12-12)


### Features

* terraform cicd ([#847](https://github.com/chanzuckerberg/happy/issues/847)) ([1be9354](https://github.com/chanzuckerberg/happy/commit/1be9354192ce8085fa967c0c9280a772a4bb6daa))

## [1.1.0](https://github.com/chanzuckerberg/happy/compare/happy-github-ci-role-v1.0.0...happy-github-ci-role-v1.1.0) (2022-12-07)


### Features

* add a service account to pods ([#835](https://github.com/chanzuckerberg/happy/issues/835)) ([203c129](https://github.com/chanzuckerberg/happy/commit/203c1294602160dfc4aacc15adf8ebc91e83af5a))

## 1.0.0 (2022-11-30)


### Features

* terraform modules for happy sharing ([#800](https://github.com/chanzuckerberg/happy/issues/800)) ([d909860](https://github.com/chanzuckerberg/happy/commit/d9098607e37b29c71bdc3ddac9fabd7ba280606b))
