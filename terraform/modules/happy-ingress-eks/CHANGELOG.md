# Changelog

## [2.12.0](https://github.com/chanzuckerberg/happy/compare/happy-ingress-eks-v2.11.0...happy-ingress-eks-v2.12.0) (2024-07-31)


### Features

* allow for a path to marked with a fixed-response deny message ([#3455](https://github.com/chanzuckerberg/happy/issues/3455)) ([1cb8dda](https://github.com/chanzuckerberg/happy/commit/1cb8dda981cd09a2354ecba470e397657abb6f0d))

## [2.11.0](https://github.com/chanzuckerberg/happy/compare/happy-ingress-eks-v2.10.0...happy-ingress-eks-v2.11.0) (2024-01-20)


### Features

* make the default stack behavior to use target type IP ([#2961](https://github.com/chanzuckerberg/happy/issues/2961)) ([79bca1b](https://github.com/chanzuckerberg/happy/commit/79bca1b7c143f0a1d07f71d84d03806d31bec38a))

## [2.10.0](https://github.com/chanzuckerberg/happy/compare/happy-ingress-eks-v2.9.0...happy-ingress-eks-v2.10.0) (2023-11-02)


### Features

* allow multiple hosts to be specified for a stack ([#2669](https://github.com/chanzuckerberg/happy/issues/2669)) ([f2023a3](https://github.com/chanzuckerberg/happy/commit/f2023a329322e59fd603208d8f1cb309e2b7541f))

## [2.9.0](https://github.com/chanzuckerberg/happy/compare/happy-ingress-eks-v2.8.0...happy-ingress-eks-v2.9.0) (2023-09-22)


### Features

* adding idle timeout config for alb ([#2486](https://github.com/chanzuckerberg/happy/issues/2486)) ([5df73b7](https://github.com/chanzuckerberg/happy/commit/5df73b7af22f7bbdc19bd960ae45bf1769819961))

## [2.8.0](https://github.com/chanzuckerberg/happy/compare/happy-ingress-eks-v2.7.0...happy-ingress-eks-v2.8.0) (2023-08-29)


### Features

* [CCIE-1729] create internal alb for service_type = "VPC" ([#2060](https://github.com/chanzuckerberg/happy/issues/2060)) ([211b1e2](https://github.com/chanzuckerberg/happy/commit/211b1e270f0e9ad00dd9b59e0cd51ce9489064c2))

## [2.7.0](https://github.com/chanzuckerberg/happy/compare/happy-ingress-eks-v2.6.0...happy-ingress-eks-v2.7.0) (2023-08-28)


### Features

* increase the health check interval ([#2378](https://github.com/chanzuckerberg/happy/issues/2378)) ([fcb0fad](https://github.com/chanzuckerberg/happy/commit/fcb0fad658ee0cecd01921dd0cb3f45901cfaf68))

## [2.6.0](https://github.com/chanzuckerberg/happy/compare/happy-ingress-eks-v2.5.0...happy-ingress-eks-v2.6.0) (2023-07-28)


### Features

* Linkerd Service Mesh For E2E Encryption and Access Control ([#1839](https://github.com/chanzuckerberg/happy/issues/1839)) ([e3f34da](https://github.com/chanzuckerberg/happy/commit/e3f34da289232f0ea92c0c3ef9d8d63e3c71f05c))

## [2.5.0](https://github.com/chanzuckerberg/happy/compare/happy-ingress-eks-v2.4.0...happy-ingress-eks-v2.5.0) (2023-05-24)


### Features

* Ingress for pods exposing HTTPS ([#1775](https://github.com/chanzuckerberg/happy/issues/1775)) ([e02675f](https://github.com/chanzuckerberg/happy/commit/e02675fbcd1c01acbc77a510c1fe385d9e42e5cb))

## [2.4.0](https://github.com/chanzuckerberg/happy/compare/happy-ingress-eks-v2.3.0...happy-ingress-eks-v2.4.0) (2023-03-07)


### Features

* give the modules an option to configure a Web ACL to protect its endpoints ([#1275](https://github.com/chanzuckerberg/happy/issues/1275)) ([90dae59](https://github.com/chanzuckerberg/happy/commit/90dae59595b041d24765123ca56c85021fe46cdb))


### Bug Fixes

* WAF assignment null condition ([#1301](https://github.com/chanzuckerberg/happy/issues/1301)) ([7ce142e](https://github.com/chanzuckerberg/happy/commit/7ce142ead96e012a192901fa5529ed6a0c2cb7bc))

## [2.3.0](https://github.com/chanzuckerberg/happy/compare/happy-ingress-eks-v2.2.0...happy-ingress-eks-v2.3.0) (2023-02-24)


### Features

* Annotate k8s resources created by happy with stack ownership labels ([#1247](https://github.com/chanzuckerberg/happy/issues/1247)) ([4403cd8](https://github.com/chanzuckerberg/happy/commit/4403cd8404ccdec96936bb033a94a3d7a2f4e58b))

## [2.2.0](https://github.com/chanzuckerberg/happy/compare/happy-ingress-eks-v2.1.0...happy-ingress-eks-v2.2.0) (2023-02-17)


### Features

* allow users to create bypasses for their OIDC ([#1149](https://github.com/chanzuckerberg/happy/issues/1149)) ([078ee17](https://github.com/chanzuckerberg/happy/commit/078ee17b36436ce92b5ad0efdade143d1f306879))

## [2.1.0](https://github.com/chanzuckerberg/happy/compare/happy-ingress-eks-v2.0.0...happy-ingress-eks-v2.1.0) (2023-02-03)


### Features

* Sample Happy Environment EKS Datadog dashboard ([#1066](https://github.com/chanzuckerberg/happy/issues/1066)) ([b4c9f3f](https://github.com/chanzuckerberg/happy/commit/b4c9f3fb7df7d131093a282cb2b54fe83f1e5143))

## [2.0.0](https://github.com/chanzuckerberg/happy/compare/happy-ingress-eks-v1.0.0...happy-ingress-eks-v2.0.0) (2023-01-31)


### ⚠ BREAKING CHANGES

* authenticate ALBs for ingresses ([#1021](https://github.com/chanzuckerberg/happy/issues/1021))

### Features

* authenticate ALBs for ingresses ([#1021](https://github.com/chanzuckerberg/happy/issues/1021)) ([7cd9375](https://github.com/chanzuckerberg/happy/commit/7cd937576a11b16cbf07e3babf268649c48c0976))

## 1.0.0 (2023-01-24)


### Features

* (CCIE-1004) Enable creation of stack-level ingress resources with a context based routing support ([#986](https://github.com/chanzuckerberg/happy/issues/986)) ([f258387](https://github.com/chanzuckerberg/happy/commit/f258387b72c1a0753c2779a79b0de8da56df71f1))
