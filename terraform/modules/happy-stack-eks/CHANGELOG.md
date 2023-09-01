# Changelog

## [4.16.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v4.15.0...happy-stack-eks-v4.16.0) (2023-09-01)


### Features

* add ability to skip injecting configs ([#2399](https://github.com/chanzuckerberg/happy/issues/2399)) ([0a858ec](https://github.com/chanzuckerberg/happy/commit/0a858ec98eaf416a22653151eea750929d55927c))

## [4.15.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v4.14.0...happy-stack-eks-v4.15.0) (2023-08-29)


### Features

* [CCIE-1729] create internal alb for service_type = "VPC" ([#2060](https://github.com/chanzuckerberg/happy/issues/2060)) ([211b1e2](https://github.com/chanzuckerberg/happy/commit/211b1e270f0e9ad00dd9b59e0cd51ce9489064c2))
* Support for ECR tag immutability ([#2376](https://github.com/chanzuckerberg/happy/issues/2376)) ([c1d5f5b](https://github.com/chanzuckerberg/happy/commit/c1d5f5b6e6a093c19ba2a092111842cc0e4f195f))

## [4.14.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v4.13.1...happy-stack-eks-v4.14.0) (2023-08-25)


### Features

* GPU support ([#2349](https://github.com/chanzuckerberg/happy/issues/2349)) ([d889c80](https://github.com/chanzuckerberg/happy/commit/d889c80983c24a172e0ebb051166dbd72f1a6edf))

## [4.13.1](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v4.13.0...happy-stack-eks-v4.13.1) (2023-08-17)


### Bug Fixes

* Fail-fast is not applicable to k8s cleanup ([#2278](https://github.com/chanzuckerberg/happy/issues/2278)) ([37c3c84](https://github.com/chanzuckerberg/happy/commit/37c3c84013cae43729d923017faeb5c7a52b27be))
* Happy update reports success on failed deployment when ECS rolls back the task version ([#2268](https://github.com/chanzuckerberg/happy/issues/2268)) ([7adf8e6](https://github.com/chanzuckerberg/happy/commit/7adf8e654979bedd01c9c824ba1489901524b2d1))

## [4.13.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v4.12.0...happy-stack-eks-v4.13.0) (2023-08-08)


### Features

* add cmd and args for services and tasks ([#2207](https://github.com/chanzuckerberg/happy/issues/2207)) ([7e648b7](https://github.com/chanzuckerberg/happy/commit/7e648b79110bd367752d835a8178c2009902c807))
* additional envs for tasks ([#2213](https://github.com/chanzuckerberg/happy/issues/2213)) ([4707d22](https://github.com/chanzuckerberg/happy/commit/4707d22c064fad221fb385e2c2e572c8dcd90736))
* iam service accounts for tasks ([#2214](https://github.com/chanzuckerberg/happy/issues/2214)) ([0f34396](https://github.com/chanzuckerberg/happy/commit/0f34396cc4d7915442201d0aa392a2e1dd1eb122))


### Bug Fixes

* wrong variable types for aws_iam object ([#2215](https://github.com/chanzuckerberg/happy/issues/2215)) ([cfe079f](https://github.com/chanzuckerberg/happy/commit/cfe079f41231da232a0decfe9226fbacf1cdb6ac))

## [4.12.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v4.11.0...happy-stack-eks-v4.12.0) (2023-08-02)


### Features

* allow specifying base directory for secret mounting ([#2168](https://github.com/chanzuckerberg/happy/issues/2168)) ([c28be3a](https://github.com/chanzuckerberg/happy/commit/c28be3a8e686ae84eb2ab0dbba90b02a9161fb08))
* expose env vars and cron vars to stack ([#2098](https://github.com/chanzuckerberg/happy/issues/2098)) ([7d370c8](https://github.com/chanzuckerberg/happy/commit/7d370c8018af7f6744ddac5fa7a492e9c2fb9515))
* expose imagepullpolicy on happy stack ([#2129](https://github.com/chanzuckerberg/happy/issues/2129)) ([e2f3b0d](https://github.com/chanzuckerberg/happy/commit/e2f3b0de238f12189aae62c70b4146910e13808b))

## [4.11.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v4.10.0...happy-stack-eks-v4.11.0) (2023-07-28)


### Features

* Linkerd Service Mesh For E2E Encryption and Access Control ([#1839](https://github.com/chanzuckerberg/happy/issues/1839)) ([e3f34da](https://github.com/chanzuckerberg/happy/commit/e3f34da289232f0ea92c0c3ef9d8d63e3c71f05c))

## [4.10.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v4.9.0...happy-stack-eks-v4.10.0) (2023-06-22)


### Features

* reuse happy client ([#1960](https://github.com/chanzuckerberg/happy/issues/1960)) ([fc3991d](https://github.com/chanzuckerberg/happy/commit/fc3991d0670579e34013e854e6a5a4f3fc4e189e))

## [4.9.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v4.8.0...happy-stack-eks-v4.9.0) (2023-06-20)


### Features

* pass along additional labels to pods ([#1905](https://github.com/chanzuckerberg/happy/issues/1905)) ([1f4de06](https://github.com/chanzuckerberg/happy/commit/1f4de06b1243a9e46ba2bdb6406179484204c868))

## [4.8.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v4.7.0...happy-stack-eks-v4.8.0) (2023-05-30)


### Features

* Example of task usage in happy EKS ([#1776](https://github.com/chanzuckerberg/happy/issues/1776)) ([2af7c7f](https://github.com/chanzuckerberg/happy/commit/2af7c7faa87938ea859db26fe143eca429f61d86))

## [4.7.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v4.6.0...happy-stack-eks-v4.7.0) (2023-05-24)


### Features

* Ingress for pods exposing HTTPS ([#1775](https://github.com/chanzuckerberg/happy/issues/1775)) ([e02675f](https://github.com/chanzuckerberg/happy/commit/e02675fbcd1c01acbc77a510c1fe385d9e42e5cb))
* new features in support of sidecar SSL termination ([#1762](https://github.com/chanzuckerberg/happy/issues/1762)) ([f78522b](https://github.com/chanzuckerberg/happy/commit/f78522b2ed847ade83d04c06d82656b4490af9bf))


### Bug Fixes

* Service port is not populated ([#1801](https://github.com/chanzuckerberg/happy/issues/1801)) ([bda5172](https://github.com/chanzuckerberg/happy/commit/bda5172115ab192fb6c1197ab987f6cee823fef0))

## [4.6.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v4.5.2...happy-stack-eks-v4.6.0) (2023-05-09)


### Features

* Sidecar support for services ([#1727](https://github.com/chanzuckerberg/happy/issues/1727)) ([8c5c884](https://github.com/chanzuckerberg/happy/commit/8c5c884804a4e88d1e3163f266127e6ddb336c05))


### Bug Fixes

* Fix sidecar validation rules ([#1744](https://github.com/chanzuckerberg/happy/issues/1744)) ([f61534e](https://github.com/chanzuckerberg/happy/commit/f61534eeb699b79d70d63e8c00571c96cfd581e8))

## [4.5.2](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v4.5.1...happy-stack-eks-v4.5.2) (2023-04-24)


### Bug Fixes

* state drilling with dryrun; use context ([#1607](https://github.com/chanzuckerberg/happy/issues/1607)) ([a75376a](https://github.com/chanzuckerberg/happy/commit/a75376a849940d9cdf45accbc1ec0357dbd0c3f8))

## [4.5.1](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v4.5.0...happy-stack-eks-v4.5.1) (2023-04-20)


### Bug Fixes

* store waf arn instead of name ([#1604](https://github.com/chanzuckerberg/happy/issues/1604)) ([e0f58ba](https://github.com/chanzuckerberg/happy/commit/e0f58ba94e59c79840fe8fb61df877e6dd0bb233))

## [4.5.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v4.4.0...happy-stack-eks-v4.5.0) (2023-04-20)


### Features

* filter the stacks by app; display git info ([#1594](https://github.com/chanzuckerberg/happy/issues/1594)) ([665f35c](https://github.com/chanzuckerberg/happy/commit/665f35c39d7eff37ff8d0bca34f37db08f0eb753))

## [4.4.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v4.3.0...happy-stack-eks-v4.4.0) (2023-04-07)


### Features

* add example for target_group_only ([#1489](https://github.com/chanzuckerberg/happy/issues/1489)) ([807d4cc](https://github.com/chanzuckerberg/happy/commit/807d4ccb493dc055030a584714b737fa28580059))


### Bug Fixes

* terraform and config.json in first example project ([#1483](https://github.com/chanzuckerberg/happy/issues/1483)) ([2a90b99](https://github.com/chanzuckerberg/happy/commit/2a90b99a374beffb055886f2233a49c4246ef2ba))

## [4.3.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v4.2.1...happy-stack-eks-v4.3.0) (2023-04-05)


### Features

* add examples folder with first happy example project ([#1470](https://github.com/chanzuckerberg/happy/issues/1470)) ([145c593](https://github.com/chanzuckerberg/happy/commit/145c593ccf42efa175622b45e19c263c884d672a))

## [4.2.1](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v4.2.0...happy-stack-eks-v4.2.1) (2023-03-17)


### Bug Fixes

* front stack with WAF if it exists, ignore otherwise ([#1385](https://github.com/chanzuckerberg/happy/issues/1385)) ([faf9f0d](https://github.com/chanzuckerberg/happy/commit/faf9f0d710878f072bc4913e282b20ebb5391167))

## [4.2.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v4.1.1...happy-stack-eks-v4.2.0) (2023-03-16)


### Features

* Platform architecture affinity ([#1375](https://github.com/chanzuckerberg/happy/issues/1375)) ([e9b81be](https://github.com/chanzuckerberg/happy/commit/e9b81be2bb737d078902f4e1c75b65e8fb73db11))

## [4.1.1](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v4.1.0...happy-stack-eks-v4.1.1) (2023-03-15)


### Bug Fixes

* Remove un-needed eks parameters from happy-service-eks ([#1372](https://github.com/chanzuckerberg/happy/issues/1372)) ([89afcdf](https://github.com/chanzuckerberg/happy/commit/89afcdf9a4a3dd232d6a2d912d8e161fed0e6e8e))

## [4.1.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v4.0.0...happy-stack-eks-v4.1.0) (2023-03-07)


### Features

* give the modules an option to configure a Web ACL to protect its endpoints ([#1275](https://github.com/chanzuckerberg/happy/issues/1275)) ([90dae59](https://github.com/chanzuckerberg/happy/commit/90dae59595b041d24765123ca56c85021fe46cdb))


### Bug Fixes

* Randomize group name to prevent ALB reuse when stack is recreated ([#1287](https://github.com/chanzuckerberg/happy/issues/1287)) ([d953a78](https://github.com/chanzuckerberg/happy/commit/d953a78640ab69286f116b2ed9ab7bf418c72c20))

## [4.0.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v3.1.0...happy-stack-eks-v4.0.0) (2023-03-07)


### ⚠ BREAKING CHANGES

* refactor update/clean; autocreate ECR ([#1232](https://github.com/chanzuckerberg/happy/issues/1232))

### Features

* refactor update/clean; autocreate ECR ([#1232](https://github.com/chanzuckerberg/happy/issues/1232)) ([b498074](https://github.com/chanzuckerberg/happy/commit/b4980740c3ddc716abe530fb2112dfe41bc6ab60))

## [3.1.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v3.0.0...happy-stack-eks-v3.1.0) (2023-02-17)


### Features

* allow users to create bypasses for their OIDC ([#1149](https://github.com/chanzuckerberg/happy/issues/1149)) ([078ee17](https://github.com/chanzuckerberg/happy/commit/078ee17b36436ce92b5ad0efdade143d1f306879))

## [3.0.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v2.4.0...happy-stack-eks-v3.0.0) (2023-02-13)


### ⚠ BREAKING CHANGES

* inject happy config to stacks ([#1108](https://github.com/chanzuckerberg/happy/issues/1108))

### Features

* inject happy config to stacks ([#1108](https://github.com/chanzuckerberg/happy/issues/1108)) ([9cb49c7](https://github.com/chanzuckerberg/happy/commit/9cb49c7f7bd6819541510e4f31ab5fd112579457))

## [2.4.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v2.3.0...happy-stack-eks-v2.4.0) (2023-02-03)


### Features

* Sample Happy Environment EKS Datadog dashboard ([#1066](https://github.com/chanzuckerberg/happy/issues/1066)) ([b4c9f3f](https://github.com/chanzuckerberg/happy/commit/b4c9f3fb7df7d131093a282cb2b54fe83f1e5143))

## [2.3.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v2.2.0...happy-stack-eks-v2.3.0) (2023-02-02)


### Features

* Add grouping tags ([#1060](https://github.com/chanzuckerberg/happy/issues/1060)) ([713015f](https://github.com/chanzuckerberg/happy/commit/713015ff7c24278c6315b9ad0ce04e98fb56bb4e))
* Create stack level dashboard ([#1062](https://github.com/chanzuckerberg/happy/issues/1062)) ([b346da9](https://github.com/chanzuckerberg/happy/commit/b346da951e14f4c10ab7c9a936990d47913b4c92))

## [2.2.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v2.1.1...happy-stack-eks-v2.2.0) (2023-02-01)


### Features

* Host level routing assigns an incorrect group name to an ingress, causing ALB not to be provisioned ([#1054](https://github.com/chanzuckerberg/happy/issues/1054)) ([16432eb](https://github.com/chanzuckerberg/happy/commit/16432eb0f5170f4ae601f460e8e73abdc0a747b0))

## [2.1.1](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v2.1.0...happy-stack-eks-v2.1.1) (2023-02-01)


### Bug Fixes

* Force-release happy-stack-eks ([#1052](https://github.com/chanzuckerberg/happy/issues/1052)) ([6e41fe9](https://github.com/chanzuckerberg/happy/commit/6e41fe95a43aa19e3867127d5a4596b3ca62c2ab))

## [2.1.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v2.0.0...happy-stack-eks-v2.1.0) (2023-02-01)


### Features

* Enable mapping of additional environment variables from secrets ([#1046](https://github.com/chanzuckerberg/happy/issues/1046)) ([6ef2fea](https://github.com/chanzuckerberg/happy/commit/6ef2feaf13d07a7848f8498ed14653610d1edc94))

## [2.0.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v1.8.0...happy-stack-eks-v2.0.0) (2023-01-31)


### ⚠ BREAKING CHANGES

* authenticate ALBs for ingresses ([#1021](https://github.com/chanzuckerberg/happy/issues/1021))

### Features

* authenticate ALBs for ingresses ([#1021](https://github.com/chanzuckerberg/happy/issues/1021)) ([7cd9375](https://github.com/chanzuckerberg/happy/commit/7cd937576a11b16cbf07e3babf268649c48c0976))

## [1.8.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v1.7.0...happy-stack-eks-v1.8.0) (2023-01-25)


### Features

* Allow configuration for the liveness and health probe timings to propagate to the stack ([#1010](https://github.com/chanzuckerberg/happy/issues/1010)) ([00976e8](https://github.com/chanzuckerberg/happy/commit/00976e84a8810bb79612da2b8d714392c4316702))

## [1.7.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v1.6.0...happy-stack-eks-v1.7.0) (2023-01-24)


### Features

* (CCIE-1004) Enable creation of stack-level ingress resources with a context based routing support ([#986](https://github.com/chanzuckerberg/happy/issues/986)) ([f258387](https://github.com/chanzuckerberg/happy/commit/f258387b72c1a0753c2779a79b0de8da56df71f1))

## [1.6.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v1.5.0...happy-stack-eks-v1.6.0) (2023-01-24)


### Features

* add synthetics to happy stack services ([#988](https://github.com/chanzuckerberg/happy/issues/988)) ([0f8eb5c](https://github.com/chanzuckerberg/happy/commit/0f8eb5c908b5133fecd35fc3d39fe7e441abd091))
* add tags from integration secret ([#990](https://github.com/chanzuckerberg/happy/issues/990)) ([46fcd8a](https://github.com/chanzuckerberg/happy/commit/46fcd8a99118b70add0feaecc0d9dd4358100bf0))

## [1.5.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v1.4.2...happy-stack-eks-v1.5.0) (2023-01-17)


### Features

* allow users to specify env vars in stacks ([#972](https://github.com/chanzuckerberg/happy/issues/972)) ([f53858e](https://github.com/chanzuckerberg/happy/commit/f53858e512d3588e44c651cbce0e2dc12fe69edf))

## [1.4.2](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v1.4.1...happy-stack-eks-v1.4.2) (2023-01-04)


### Bug Fixes

* make env vars nonsensitive to allow for_each ([#927](https://github.com/chanzuckerberg/happy/issues/927)) ([0aaf238](https://github.com/chanzuckerberg/happy/commit/0aaf23826c54d1980f6947c20a7623076a5954e1))

## [1.4.1](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v1.4.0...happy-stack-eks-v1.4.1) (2022-12-21)


### Bug Fixes

* update additional_env_vars for eks stack ([#910](https://github.com/chanzuckerberg/happy/issues/910)) ([3e0cea1](https://github.com/chanzuckerberg/happy/commit/3e0cea11efc9770626e7f1af66080e1d9fcc8be1))

## [1.4.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v1.3.0...happy-stack-eks-v1.4.0) (2022-12-21)


### Features

* automatically inject db env vars into eks env ([#908](https://github.com/chanzuckerberg/happy/issues/908)) ([99123b2](https://github.com/chanzuckerberg/happy/commit/99123b2b1648b1b7c6ce756942c9fb925b31e07a))

## [1.3.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v1.2.0...happy-stack-eks-v1.3.0) (2022-12-12)


### Features

* terraform cicd ([#847](https://github.com/chanzuckerberg/happy/issues/847)) ([1be9354](https://github.com/chanzuckerberg/happy/commit/1be9354192ce8085fa967c0c9280a772a4bb6daa))

## [1.2.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v1.1.0...happy-stack-eks-v1.2.0) (2022-12-08)


### Features

* add optional/depin internal modules ([#846](https://github.com/chanzuckerberg/happy/issues/846)) ([348fc78](https://github.com/chanzuckerberg/happy/commit/348fc7876fd7427487d7ea340171898a39d4b05b))


### Bug Fixes

* remove old files ([#839](https://github.com/chanzuckerberg/happy/issues/839)) ([8659e46](https://github.com/chanzuckerberg/happy/commit/8659e463f73e4ce16f9a43a49e4134f66c6ba518))

## [1.1.0](https://github.com/chanzuckerberg/happy/compare/happy-stack-eks-v1.0.0...happy-stack-eks-v1.1.0) (2022-12-07)


### Features

* add a service account to pods ([#835](https://github.com/chanzuckerberg/happy/issues/835)) ([203c129](https://github.com/chanzuckerberg/happy/commit/203c1294602160dfc4aacc15adf8ebc91e83af5a))


### Bug Fixes

* input variables and tags in happy eks modules ([#838](https://github.com/chanzuckerberg/happy/issues/838)) ([175dc76](https://github.com/chanzuckerberg/happy/commit/175dc7652735e5683dced24d8cdfa48101355c72))

## 1.0.0 (2022-11-30)


### Features

* terraform modules for happy sharing ([#800](https://github.com/chanzuckerberg/happy/issues/800)) ([d909860](https://github.com/chanzuckerberg/happy/commit/d9098607e37b29c71bdc3ddac9fabd7ba280606b))


### Bug Fixes

* happy module bugs ([#806](https://github.com/chanzuckerberg/happy/issues/806)) ([7a87501](https://github.com/chanzuckerberg/happy/commit/7a875019afda4bc016558ee06c846c940a71a6dd))
