# Changelog

## [4.12.2](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.12.1...happy-env-eks-v4.12.2) (2023-09-22)


### Bug Fixes

* Add permissions to github ci role to allow retrieval of ecr image scanning ([#2492](https://github.com/chanzuckerberg/happy/issues/2492)) ([a6b7116](https://github.com/chanzuckerberg/happy/commit/a6b71169993165a55b8f139225a53e5a229367e9))

## [4.12.1](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.12.0...happy-env-eks-v4.12.1) (2023-09-22)


### Bug Fixes

* Fix scan on push typo ([#2489](https://github.com/chanzuckerberg/happy/issues/2489)) ([140fad2](https://github.com/chanzuckerberg/happy/commit/140fad21b7ecbb06d65615db2ccde223c3ce7821))

## [4.12.0](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.11.1...happy-env-eks-v4.12.0) (2023-09-21)


### Features

* Bump happy-env-eks ([#2485](https://github.com/chanzuckerberg/happy/issues/2485)) ([99825aa](https://github.com/chanzuckerberg/happy/commit/99825aa063f0c60de9f0fdab82476d8d4fb2ce81))

## [4.11.1](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.11.0...happy-env-eks-v4.11.1) (2023-09-15)


### Bug Fixes

* allow for changing the default token endpoint auth ([#2458](https://github.com/chanzuckerberg/happy/issues/2458)) ([1e5b66e](https://github.com/chanzuckerberg/happy/commit/1e5b66e528ddc29643e0e00beea7be6d2ae77ce9))
* thread through the token auth config ([#2460](https://github.com/chanzuckerberg/happy/issues/2460)) ([5a5db25](https://github.com/chanzuckerberg/happy/commit/5a5db25099fd29e9f96d3c9485069f8e5206f9a1))

## [4.11.0](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.10.0...happy-env-eks-v4.11.0) (2023-09-14)


### Features

* allow happy envs to configure oidc ([#2446](https://github.com/chanzuckerberg/happy/issues/2446)) ([0041153](https://github.com/chanzuckerberg/happy/commit/0041153f07a124c1101dc685379ef1c0249eafd3))

## [4.10.0](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.9.0...happy-env-eks-v4.10.0) (2023-08-29)


### Features

* Support for ECR tag immutability ([#2376](https://github.com/chanzuckerberg/happy/issues/2376)) ([c1d5f5b](https://github.com/chanzuckerberg/happy/commit/c1d5f5b6e6a093c19ba2a092111842cc0e4f195f))

## [4.9.0](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.8.2...happy-env-eks-v4.9.0) (2023-08-16)


### Features

* add policy variable for s3 buckets ([#2270](https://github.com/chanzuckerberg/happy/issues/2270)) ([37640c4](https://github.com/chanzuckerberg/happy/commit/37640c4037a14de96128285e6d5c0eb8d63fd344))

## [4.8.2](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.8.1...happy-env-eks-v4.8.2) (2023-08-15)


### Bug Fixes

* allow for a bucket policy to be specified ([#2257](https://github.com/chanzuckerberg/happy/issues/2257)) ([47c3ddc](https://github.com/chanzuckerberg/happy/commit/47c3ddce7fbd57d31c07161908e011c93c8900e5))

## [4.8.1](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.8.0...happy-env-eks-v4.8.1) (2023-07-19)


### Bug Fixes

* underscores in db parameter groups ([#2063](https://github.com/chanzuckerberg/happy/issues/2063)) ([67edf7c](https://github.com/chanzuckerberg/happy/commit/67edf7caff1354e48ae9f52ce032915c89548275))

## [4.8.0](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.7.5...happy-env-eks-v4.8.0) (2023-07-07)


### Features

* support EKS paths for new service accounts ([#2012](https://github.com/chanzuckerberg/happy/issues/2012)) ([e1f407b](https://github.com/chanzuckerberg/happy/commit/e1f407b95baa0daee5069ec3d67eef64263c4383))

## [4.7.5](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.7.4...happy-env-eks-v4.7.5) (2023-06-27)


### Bug Fixes

* leave out WAF configuration in Happy EKS setup ([#1970](https://github.com/chanzuckerberg/happy/issues/1970)) ([e906ada](https://github.com/chanzuckerberg/happy/commit/e906ada431c012639cf6e4a9a1e183daada016e1))

## [4.7.4](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.7.3...happy-env-eks-v4.7.4) (2023-05-30)


### Bug Fixes

* aws provider 5.0 deprecated source_json ([#1810](https://github.com/chanzuckerberg/happy/issues/1810)) ([7b69d30](https://github.com/chanzuckerberg/happy/commit/7b69d3086112972c5792edf31509dc1bde4ba23b))
* Handle empty and null ecr policies ([#1813](https://github.com/chanzuckerberg/happy/issues/1813)) ([b2e60f1](https://github.com/chanzuckerberg/happy/commit/b2e60f1dcb948a1cc3ec860c26b3ed541112b5de))

## [4.7.3](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.7.2...happy-env-eks-v4.7.3) (2023-05-09)


### Bug Fixes

* always create dynamo rules for ci role ([#1731](https://github.com/chanzuckerberg/happy/issues/1731)) ([c8de064](https://github.com/chanzuckerberg/happy/commit/c8de0644aa28a014e1275100d6d1e1e80197107d))

## [4.7.2](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.7.1...happy-env-eks-v4.7.2) (2023-05-05)


### Bug Fixes

* wrong count argument ([#1689](https://github.com/chanzuckerberg/happy/issues/1689)) ([45b2709](https://github.com/chanzuckerberg/happy/commit/45b27099d6504d16f58789e81bca7d7ef1c7e2b0))

## [4.7.1](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.7.0...happy-env-eks-v4.7.1) (2023-05-03)


### Bug Fixes

* dynamo alphanumeric sid ([#1687](https://github.com/chanzuckerberg/happy/issues/1687)) ([61afbda](https://github.com/chanzuckerberg/happy/commit/61afbdac2796d213b2b722c5e4f42044e00cfe48))

## [4.7.0](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.6.7...happy-env-eks-v4.7.0) (2023-05-02)


### Features

* add command to see the configured CI roles for env ([#1686](https://github.com/chanzuckerberg/happy/issues/1686)) ([a249cc0](https://github.com/chanzuckerberg/happy/commit/a249cc0a4fc61af413312b300f1fc4695529ee2e))


### Bug Fixes

* thread CI role through happy env ([#1683](https://github.com/chanzuckerberg/happy/issues/1683)) ([78b4be9](https://github.com/chanzuckerberg/happy/commit/78b4be95b7f4f4be95cf18a3d3b9920a28f409da))

## [4.6.7](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.6.6...happy-env-eks-v4.6.7) (2023-04-27)


### Bug Fixes

* Trigger the release of happy-env-eks ([#1672](https://github.com/chanzuckerberg/happy/issues/1672)) ([021ca49](https://github.com/chanzuckerberg/happy/commit/021ca49728c015a9ffdd71f922086e61a4a5d20c))

## [4.6.6](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.6.5...happy-env-eks-v4.6.6) (2023-04-27)


### Bug Fixes

* count being used on resources known after apply ([#1665](https://github.com/chanzuckerberg/happy/issues/1665)) ([00eec6d](https://github.com/chanzuckerberg/happy/commit/00eec6d86b489408c2347ff57179d5ad9de43414))

## [4.6.5](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.6.4...happy-env-eks-v4.6.5) (2023-04-25)


### Bug Fixes

* Trigger a happy-env-eks release ([#1654](https://github.com/chanzuckerberg/happy/issues/1654)) ([05b018a](https://github.com/chanzuckerberg/happy/commit/05b018a23979d2045fcd297e78b27ea5ef7015e3))

## [4.6.4](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.6.3...happy-env-eks-v4.6.4) (2023-04-24)


### Bug Fixes

* waf is a count ([#1616](https://github.com/chanzuckerberg/happy/issues/1616)) ([3696d34](https://github.com/chanzuckerberg/happy/commit/3696d34d1c70c60f84c252fffddf32b0dd0a545a))

## [4.6.3](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.6.2...happy-env-eks-v4.6.3) (2023-04-20)


### Bug Fixes

* store waf arn instead of name ([#1604](https://github.com/chanzuckerberg/happy/issues/1604)) ([e0f58ba](https://github.com/chanzuckerberg/happy/commit/e0f58ba94e59c79840fe8fb61df877e6dd0bb233))

## [4.6.2](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.6.1...happy-env-eks-v4.6.2) (2023-04-20)


### Bug Fixes

* remove skip; actually validate child modules ([#1602](https://github.com/chanzuckerberg/happy/issues/1602)) ([79c6719](https://github.com/chanzuckerberg/happy/commit/79c671919e4fa897c93d441fa60825694f65b1ce))

## [4.6.1](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.6.0...happy-env-eks-v4.6.1) (2023-04-20)


### Bug Fixes

* Disable terraform validate on select modules ([#1600](https://github.com/chanzuckerberg/happy/issues/1600)) ([0294798](https://github.com/chanzuckerberg/happy/commit/0294798010874c57e601c4f78f0a4efd899796a8))

## [4.6.0](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.5.7...happy-env-eks-v4.6.0) (2023-04-20)


### Features

* CDI-1356 - add database creds to output ([#1595](https://github.com/chanzuckerberg/happy/issues/1595)) ([652d6b6](https://github.com/chanzuckerberg/happy/commit/652d6b6473df879b39a7723327628476441a6eb6))
* enable EKS cluster permissions for CI role ([#1589](https://github.com/chanzuckerberg/happy/issues/1589)) ([c0c6452](https://github.com/chanzuckerberg/happy/commit/c0c6452d977ae0ef0f7ba1f82ddf859ec0fdea65))

## [4.5.7](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.5.6...happy-env-eks-v4.5.7) (2023-04-18)


### Bug Fixes

* CDI-1222 Swap `rds_cluster_parameter` data type ([#1567](https://github.com/chanzuckerberg/happy/issues/1567)) ([67a6a0c](https://github.com/chanzuckerberg/happy/commit/67a6a0cca9f6c8a39922ec83a1fe94986c9965cf))

## [4.5.6](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.5.5...happy-env-eks-v4.5.6) (2023-04-14)


### Bug Fixes

* upgrade WAF version in happy-env-eks ([#1568](https://github.com/chanzuckerberg/happy/issues/1568)) ([0b5f0aa](https://github.com/chanzuckerberg/happy/commit/0b5f0aa25b65400af97954a6fd646d6c9113edb0))

## [4.5.5](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.5.4...happy-env-eks-v4.5.5) (2023-04-05)


### Bug Fixes

* version bump to get role fixed ([#1481](https://github.com/chanzuckerberg/happy/issues/1481)) ([fa42318](https://github.com/chanzuckerberg/happy/commit/fa4231834cf47ac32d3be47c7c6bf1baf4a40fc9))

## [4.5.4](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.5.3...happy-env-eks-v4.5.4) (2023-03-29)


### Bug Fixes

* bump happy-env-eks version ([#1458](https://github.com/chanzuckerberg/happy/issues/1458)) ([46a9ca7](https://github.com/chanzuckerberg/happy/commit/46a9ca7be7e1c66a63d386b5a006d7081044c551))

## [4.5.3](https://github.com/chanzuckerberg/happy/compare/happy-env-eks-v4.5.2...happy-env-eks-v4.5.3) (2023-03-29)


### Bug Fixes

* typo in provider aliases ([#1455](https://github.com/chanzuckerberg/happy/issues/1455)) ([66e3346](https://github.com/chanzuckerberg/happy/commit/66e33460f87c49eb96505217798c524c4d02a921))

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
