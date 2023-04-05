# Changelog

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
