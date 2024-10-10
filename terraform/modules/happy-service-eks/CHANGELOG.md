# Changelog

## [3.28.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.27.0...happy-service-eks-v3.28.0) (2024-10-10)


### Features

* allow k6 operator service account in rdev/staging (CCIE-3437) ([#3586](https://github.com/chanzuckerberg/happy/issues/3586)) ([afd9685](https://github.com/chanzuckerberg/happy/commit/afd9685716043b0d99b613904c034c53df700e6f))

## [3.27.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.26.0...happy-service-eks-v3.27.0) (2024-07-31)


### Features

* allow for a path to marked with a fixed-response deny message ([#3455](https://github.com/chanzuckerberg/happy/issues/3455)) ([1cb8dda](https://github.com/chanzuckerberg/happy/commit/1cb8dda981cd09a2354ecba470e397657abb6f0d))

## [3.26.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.25.0...happy-service-eks-v3.26.0) (2024-04-19)


### Features

* [CCIE-2375] ephemeral volume support ([#2977](https://github.com/chanzuckerberg/happy/issues/2977)) ([9fa7d24](https://github.com/chanzuckerberg/happy/commit/9fa7d24a50ca9b901e46cfd2bd9f47f5dc903a6a))

## [3.25.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.24.0...happy-service-eks-v3.25.0) (2024-03-12)


### Features

* command healthchecks and cli services with no network endpoints ([#3110](https://github.com/chanzuckerberg/happy/issues/3110)) ([6966e84](https://github.com/chanzuckerberg/happy/commit/6966e84d22e8192e5e486657605bfb2ef606ec15))

## [3.24.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.23.0...happy-service-eks-v3.24.0) (2024-02-26)


### Features

* Add shared cache volume support ([#3038](https://github.com/chanzuckerberg/happy/issues/3038)) ([c6d9786](https://github.com/chanzuckerberg/happy/commit/c6d9786491885fd6b4e65cc13282f12ef412b657))
* inject v2 configs into sidecars ([#3036](https://github.com/chanzuckerberg/happy/issues/3036)) ([3883ac3](https://github.com/chanzuckerberg/happy/commit/3883ac324f07b8dde3dee47b13a2fe2d69a65b07))

## [3.23.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.22.1...happy-service-eks-v3.23.0) (2024-01-20)


### Features

* allow for progress_deadline_seconds to be set ([#2963](https://github.com/chanzuckerberg/happy/issues/2963)) ([bdd581d](https://github.com/chanzuckerberg/happy/commit/bdd581dbcca11ab3e70fd6fc416346b48b3ea801))
* make the default stack behavior to use target type IP ([#2961](https://github.com/chanzuckerberg/happy/issues/2961)) ([79bca1b](https://github.com/chanzuckerberg/happy/commit/79bca1b7c143f0a1d07f71d84d03806d31bec38a))

## [3.22.1](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.22.0...happy-service-eks-v3.22.1) (2024-01-12)


### Bug Fixes

* make synthetics use additional_hostnames ([#2944](https://github.com/chanzuckerberg/happy/issues/2944)) ([12ce1b8](https://github.com/chanzuckerberg/happy/commit/12ce1b8e4a5d42a4028a5c7f46ae16a12f65ce04))

## [3.22.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.21.0...happy-service-eks-v3.22.0) (2024-01-09)


### Features

* Support args and cmd arguments for sidecar containers ([#2935](https://github.com/chanzuckerberg/happy/issues/2935)) ([ca30025](https://github.com/chanzuckerberg/happy/commit/ca300250302f8eb2ebcf4126252e563be23e419d))

## [3.21.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.20.0...happy-service-eks-v3.21.0) (2024-01-05)


### Features

* use env_from to inject happy-config secrets into k8s deployments ([#2899](https://github.com/chanzuckerberg/happy/issues/2899)) ([e73fc41](https://github.com/chanzuckerberg/happy/commit/e73fc41838855100cd49803eff57742a9ca5f1a8))

## [3.20.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.19.0...happy-service-eks-v3.20.0) (2023-11-21)


### Features

* Add support for init containers ([#2778](https://github.com/chanzuckerberg/happy/issues/2778)) ([0831554](https://github.com/chanzuckerberg/happy/commit/0831554aed28b657f68aa21a134393786b31db11))
* create a new target group before destroy ([#2617](https://github.com/chanzuckerberg/happy/issues/2617)) ([a977c0c](https://github.com/chanzuckerberg/happy/commit/a977c0cbcc2ea0878d4b7d5c8f3bca9ec0d54628))

## [3.19.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.18.1...happy-service-eks-v3.19.0) (2023-11-02)


### Features

* allow multiple hosts to be specified for a stack ([#2669](https://github.com/chanzuckerberg/happy/issues/2669)) ([f2023a3](https://github.com/chanzuckerberg/happy/commit/f2023a329322e59fd603208d8f1cb309e2b7541f))
* CCIE-2069: Add liveness and readiness timeouts ([#2664](https://github.com/chanzuckerberg/happy/issues/2664)) ([aa5734a](https://github.com/chanzuckerberg/happy/commit/aa5734afa18a40f011975f2557205fd1bea0bdd3))

## [3.18.1](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.18.0...happy-service-eks-v3.18.1) (2023-10-31)


### Bug Fixes

* Do not create a pod disruption budget for deployment with the desired_count=max_unavailable_count ([#2663](https://github.com/chanzuckerberg/happy/issues/2663)) ([6a63976](https://github.com/chanzuckerberg/happy/commit/6a639761fe383fc01b3707e82c7840a22c0a74d7))

## [3.18.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.17.2...happy-service-eks-v3.18.0) (2023-10-24)


### Features

* Enable support for pod disruption budgets and pod anti-affinity rules ([#2532](https://github.com/chanzuckerberg/happy/issues/2532)) ([71e7cd6](https://github.com/chanzuckerberg/happy/commit/71e7cd6b49aa1a3f7411fee8bf0e88c9b30df625))

## [3.17.2](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.17.1...happy-service-eks-v3.17.2) (2023-10-16)


### Bug Fixes

* trim the target group name to only 32 chars ([#2572](https://github.com/chanzuckerberg/happy/issues/2572)) ([2527f87](https://github.com/chanzuckerberg/happy/commit/2527f8761c3a4d913f76563be22291ba00af3421))

## [3.17.1](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.17.0...happy-service-eks-v3.17.1) (2023-10-12)


### Bug Fixes

* duplicate target group names ([#2566](https://github.com/chanzuckerberg/happy/issues/2566)) ([ccbcceb](https://github.com/chanzuckerberg/happy/commit/ccbccebb0b1bf3b9f042b8f4751cf623e5320624))

## [3.17.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.16.0...happy-service-eks-v3.17.0) (2023-10-03)


### Features

* cloudfront added to stack module ([#2487](https://github.com/chanzuckerberg/happy/issues/2487)) ([de3d85e](https://github.com/chanzuckerberg/happy/commit/de3d85e63e5978bc349b86d93270aebe464da866))

## [3.16.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.15.0...happy-service-eks-v3.16.0) (2023-09-22)


### Features

* adding idle timeout config for alb ([#2486](https://github.com/chanzuckerberg/happy/issues/2486)) ([5df73b7](https://github.com/chanzuckerberg/happy/commit/5df73b7af22f7bbdc19bd960ae45bf1769819961))

## [3.15.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.14.0...happy-service-eks-v3.15.0) (2023-09-14)


### Features

* add service_account_name option to allow_mesh_services ([#2443](https://github.com/chanzuckerberg/happy/issues/2443)) ([d7c76dc](https://github.com/chanzuckerberg/happy/commit/d7c76dc2e6fcbc5344af0cba3ae76353fe3d8b3b))
* Prevent non-gpu non-system workloads from being scheduled on GPU nodes ([#2442](https://github.com/chanzuckerberg/happy/issues/2442)) ([83765bf](https://github.com/chanzuckerberg/happy/commit/83765bf4d962c9922304f1e69519a5d658b8018f))

## [3.14.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.13.0...happy-service-eks-v3.14.0) (2023-08-29)


### Features

* [CCIE-1729] create internal alb for service_type = "VPC" ([#2060](https://github.com/chanzuckerberg/happy/issues/2060)) ([211b1e2](https://github.com/chanzuckerberg/happy/commit/211b1e270f0e9ad00dd9b59e0cd51ce9489064c2))
* Support for ECR tag immutability ([#2376](https://github.com/chanzuckerberg/happy/issues/2376)) ([c1d5f5b](https://github.com/chanzuckerberg/happy/commit/c1d5f5b6e6a093c19ba2a092111842cc0e4f195f))

## [3.13.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.12.2...happy-service-eks-v3.13.0) (2023-08-25)


### Features

* GPU support ([#2349](https://github.com/chanzuckerberg/happy/issues/2349)) ([d889c80](https://github.com/chanzuckerberg/happy/commit/d889c80983c24a172e0ebb051166dbd72f1a6edf))

## [3.12.2](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.12.1...happy-service-eks-v3.12.2) (2023-08-17)


### Bug Fixes

* Fail-fast is not applicable to k8s cleanup ([#2278](https://github.com/chanzuckerberg/happy/issues/2278)) ([37c3c84](https://github.com/chanzuckerberg/happy/commit/37c3c84013cae43729d923017faeb5c7a52b27be))
* Happy update reports success on failed deployment when ECS rolls back the task version ([#2268](https://github.com/chanzuckerberg/happy/issues/2268)) ([7adf8e6](https://github.com/chanzuckerberg/happy/commit/7adf8e654979bedd01c9c824ba1489901524b2d1))
* Only Always restart policy is supported ([#2277](https://github.com/chanzuckerberg/happy/issues/2277)) ([fe389c4](https://github.com/chanzuckerberg/happy/commit/fe389c436af187851dcf56978e517cdf0170fb65))

## [3.12.1](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.12.0...happy-service-eks-v3.12.1) (2023-08-14)


### Bug Fixes

* adding service type to happy routing tags ([#2253](https://github.com/chanzuckerberg/happy/issues/2253)) ([cd2cf64](https://github.com/chanzuckerberg/happy/commit/cd2cf649e16dd699de550806b2b8337604634025))

## [3.12.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.11.0...happy-service-eks-v3.12.0) (2023-08-08)


### Features

* add cmd and args for services and tasks ([#2207](https://github.com/chanzuckerberg/happy/issues/2207)) ([7e648b7](https://github.com/chanzuckerberg/happy/commit/7e648b79110bd367752d835a8178c2009902c807))
* iam service accounts for tasks ([#2214](https://github.com/chanzuckerberg/happy/issues/2214)) ([0f34396](https://github.com/chanzuckerberg/happy/commit/0f34396cc4d7915442201d0aa392a2e1dd1eb122))


### Bug Fixes

* wrong variable types for aws_iam object ([#2215](https://github.com/chanzuckerberg/happy/issues/2215)) ([cfe079f](https://github.com/chanzuckerberg/happy/commit/cfe079f41231da232a0decfe9226fbacf1cdb6ac))

## [3.11.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.10.0...happy-service-eks-v3.11.0) (2023-08-02)


### Features

* allow specifying base directory for secret mounting ([#2168](https://github.com/chanzuckerberg/happy/issues/2168)) ([c28be3a](https://github.com/chanzuckerberg/happy/commit/c28be3a8e686ae84eb2ab0dbba90b02a9161fb08))
* expose imagepullpolicy on happy stack ([#2129](https://github.com/chanzuckerberg/happy/issues/2129)) ([e2f3b0d](https://github.com/chanzuckerberg/happy/commit/e2f3b0de238f12189aae62c70b4146910e13808b))

## [3.10.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.9.0...happy-service-eks-v3.10.0) (2023-07-28)


### Features

* Linkerd Service Mesh For E2E Encryption and Access Control ([#1839](https://github.com/chanzuckerberg/happy/issues/1839)) ([e3f34da](https://github.com/chanzuckerberg/happy/commit/e3f34da289232f0ea92c0c3ef9d8d63e3c71f05c))

## [3.9.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.8.0...happy-service-eks-v3.9.0) (2023-06-22)


### Features

* reuse happy client ([#1960](https://github.com/chanzuckerberg/happy/issues/1960)) ([fc3991d](https://github.com/chanzuckerberg/happy/commit/fc3991d0670579e34013e854e6a5a4f3fc4e189e))

## [3.8.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.7.0...happy-service-eks-v3.8.0) (2023-06-20)


### Features

* pass along additional labels to pods ([#1905](https://github.com/chanzuckerberg/happy/issues/1905)) ([1f4de06](https://github.com/chanzuckerberg/happy/commit/1f4de06b1243a9e46ba2bdb6406179484204c868))

## [3.7.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.6.1...happy-service-eks-v3.7.0) (2023-06-02)


### Features

* by default make the memory and cpu limits small ([#1835](https://github.com/chanzuckerberg/happy/issues/1835)) ([d80989b](https://github.com/chanzuckerberg/happy/commit/d80989bb6840f50089e30586b3b62cf26029a3c5))

## [3.6.1](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.6.0...happy-service-eks-v3.6.1) (2023-05-30)


### Bug Fixes

* aws provider 5.0 deprecated source_json ([#1810](https://github.com/chanzuckerberg/happy/issues/1810)) ([7b69d30](https://github.com/chanzuckerberg/happy/commit/7b69d3086112972c5792edf31509dc1bde4ba23b))
* Handle empty and null ecr policies ([#1813](https://github.com/chanzuckerberg/happy/issues/1813)) ([b2e60f1](https://github.com/chanzuckerberg/happy/commit/b2e60f1dcb948a1cc3ec860c26b3ed541112b5de))

## [3.6.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.5.0...happy-service-eks-v3.6.0) (2023-05-24)


### Features

* Ingress for pods exposing HTTPS ([#1775](https://github.com/chanzuckerberg/happy/issues/1775)) ([e02675f](https://github.com/chanzuckerberg/happy/commit/e02675fbcd1c01acbc77a510c1fe385d9e42e5cb))
* new features in support of sidecar SSL termination ([#1762](https://github.com/chanzuckerberg/happy/issues/1762)) ([f78522b](https://github.com/chanzuckerberg/happy/commit/f78522b2ed847ade83d04c06d82656b4490af9bf))

## [3.5.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.4.2...happy-service-eks-v3.5.0) (2023-05-09)


### Features

* Sidecar support for services ([#1727](https://github.com/chanzuckerberg/happy/issues/1727)) ([8c5c884](https://github.com/chanzuckerberg/happy/commit/8c5c884804a4e88d1e3163f266127e6ddb336c05))

## [3.4.2](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.4.1...happy-service-eks-v3.4.2) (2023-04-27)


### Bug Fixes

* error for invalid policies on happy-iam-service-account-eks module ([#1648](https://github.com/chanzuckerberg/happy/issues/1648)) ([781d465](https://github.com/chanzuckerberg/happy/commit/781d465542984b29b62d945863205f281595440d))

## [3.4.1](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.4.0...happy-service-eks-v3.4.1) (2023-04-20)


### Bug Fixes

* ecr was saved with lockfile ([#1586](https://github.com/chanzuckerberg/happy/issues/1586)) ([12b380a](https://github.com/chanzuckerberg/happy/commit/12b380adc8dd322bffdcb141e9f20743463303e0))

## [3.4.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.3.0...happy-service-eks-v3.4.0) (2023-04-07)


### Features

* add example for target_group_only ([#1489](https://github.com/chanzuckerberg/happy/issues/1489)) ([807d4cc](https://github.com/chanzuckerberg/happy/commit/807d4ccb493dc055030a584714b737fa28580059))


### Bug Fixes

* terraform and config.json in first example project ([#1483](https://github.com/chanzuckerberg/happy/issues/1483)) ([2a90b99](https://github.com/chanzuckerberg/happy/commit/2a90b99a374beffb055886f2233a49c4246ef2ba))

## [3.3.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.2.0...happy-service-eks-v3.3.0) (2023-04-05)


### Features

* add examples folder with first happy example project ([#1470](https://github.com/chanzuckerberg/happy/issues/1470)) ([145c593](https://github.com/chanzuckerberg/happy/commit/145c593ccf42efa175622b45e19c263c884d672a))


### Bug Fixes

* empty tuple error ([#1471](https://github.com/chanzuckerberg/happy/issues/1471)) ([9b05b95](https://github.com/chanzuckerberg/happy/commit/9b05b9523d7501d4cea7a89a90f7e4eb79ea9468))

## [3.2.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.1.1...happy-service-eks-v3.2.0) (2023-03-16)


### Features

* Platform architecture affinity ([#1375](https://github.com/chanzuckerberg/happy/issues/1375)) ([e9b81be](https://github.com/chanzuckerberg/happy/commit/e9b81be2bb737d078902f4e1c75b65e8fb73db11))

## [3.1.1](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.1.0...happy-service-eks-v3.1.1) (2023-03-15)


### Bug Fixes

* Remove un-needed eks parameters from happy-service-eks ([#1372](https://github.com/chanzuckerberg/happy/issues/1372)) ([89afcdf](https://github.com/chanzuckerberg/happy/commit/89afcdf9a4a3dd232d6a2d912d8e161fed0e6e8e))

## [3.1.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v3.0.0...happy-service-eks-v3.1.0) (2023-03-07)


### Features

* give the modules an option to configure a Web ACL to protect its endpoints ([#1275](https://github.com/chanzuckerberg/happy/issues/1275)) ([90dae59](https://github.com/chanzuckerberg/happy/commit/90dae59595b041d24765123ca56c85021fe46cdb))


### Bug Fixes

* Randomize group name to prevent ALB reuse when stack is recreated ([#1287](https://github.com/chanzuckerberg/happy/issues/1287)) ([d953a78](https://github.com/chanzuckerberg/happy/commit/d953a78640ab69286f116b2ed9ab7bf418c72c20))

## [3.0.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v2.5.0...happy-service-eks-v3.0.0) (2023-03-07)


### ⚠ BREAKING CHANGES

* refactor update/clean; autocreate ECR ([#1232](https://github.com/chanzuckerberg/happy/issues/1232))

### Features

* refactor update/clean; autocreate ECR ([#1232](https://github.com/chanzuckerberg/happy/issues/1232)) ([b498074](https://github.com/chanzuckerberg/happy/commit/b4980740c3ddc716abe530fb2112dfe41bc6ab60))


### Bug Fixes

* Reduce the wait time for failing EKS deployments ([#1274](https://github.com/chanzuckerberg/happy/issues/1274)) ([21801fa](https://github.com/chanzuckerberg/happy/commit/21801fa645aa8f7377242c7545a7ee4c92e64634))

## [2.5.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v2.4.1...happy-service-eks-v2.5.0) (2023-02-24)


### Features

* Annotate k8s resources created by happy with stack ownership labels ([#1247](https://github.com/chanzuckerberg/happy/issues/1247)) ([4403cd8](https://github.com/chanzuckerberg/happy/commit/4403cd8404ccdec96936bb033a94a3d7a2f4e58b))

## [2.4.1](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v2.4.0...happy-service-eks-v2.4.1) (2023-02-21)


### Bug Fixes

* Fix HPA target ([#1213](https://github.com/chanzuckerberg/happy/issues/1213)) ([46f91dd](https://github.com/chanzuckerberg/happy/commit/46f91ddb0ad4834ecb62b22f4f673d3e73da0c07))

## [2.4.0](https://github.com/chanzuckerberg/happy/compare/happy-service-eks-v2.3.0...happy-service-eks-v2.4.0) (2023-02-21)


### Features

* Happy EKS application autoscaling support ([#1174](https://github.com/chanzuckerberg/happy/issues/1174)) ([749d23f](https://github.com/chanzuckerberg/happy/commit/749d23fec3fc46cd24ec5f387fd10abc3d1993a0))

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
