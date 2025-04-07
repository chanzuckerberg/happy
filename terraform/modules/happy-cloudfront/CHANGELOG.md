# Changelog

## [1.5.0](https://github.com/chanzuckerberg/happy/compare/happy-cloudfront-v1.4.0...happy-cloudfront-v1.5.0) (2025-04-07)


### Features

* add ability to blacklist and remove georestriction in happy cloudfront (CCIE-4360) ([#3943](https://github.com/chanzuckerberg/happy/issues/3943)) ([69d1ad7](https://github.com/chanzuckerberg/happy/commit/69d1ad7549aec118116ab0f398c9dc99ba95856a))

## [1.4.0](https://github.com/chanzuckerberg/happy/compare/happy-cloudfront-v1.3.0...happy-cloudfront-v1.4.0) (2024-11-25)


### Features

* Adding optional origin_access_control_id for cross-account S3 origins ([#3688](https://github.com/chanzuckerberg/happy/issues/3688)) ([3e3d850](https://github.com/chanzuckerberg/happy/commit/3e3d85018fbb00681f8e8f6885609baf36d8af1b))

## [1.3.0](https://github.com/chanzuckerberg/happy/compare/happy-cloudfront-v1.2.2...happy-cloudfront-v1.3.0) (2024-11-22)


### Features

* enabling support for S3 access ([#3670](https://github.com/chanzuckerberg/happy/issues/3670)) ([f435f63](https://github.com/chanzuckerberg/happy/commit/f435f63c6c21976fd9459e4a9d150ddff1a1059a))

## [1.2.2](https://github.com/chanzuckerberg/happy/compare/happy-cloudfront-v1.2.1...happy-cloudfront-v1.2.2) (2024-07-19)


### Bug Fixes

* make sure cloudfront is created in same provider as the cert ([#3432](https://github.com/chanzuckerberg/happy/issues/3432)) ([3bdc81a](https://github.com/chanzuckerberg/happy/commit/3bdc81aa2c45c08d3c15cb1ee3b58e3364843538))

## [1.2.1](https://github.com/chanzuckerberg/happy/compare/happy-cloudfront-v1.2.0...happy-cloudfront-v1.2.1) (2024-07-10)


### Bug Fixes

* make the cloudfront comment static so we won't reach length errors ([#3417](https://github.com/chanzuckerberg/happy/issues/3417)) ([5bf832d](https://github.com/chanzuckerberg/happy/commit/5bf832df8b2c87c713266c2b1841351017469208))

## [1.2.0](https://github.com/chanzuckerberg/happy/compare/happy-cloudfront-v1.1.0...happy-cloudfront-v1.2.0) (2024-07-10)


### Features

* allow cloudfront module to use multiple origins ([#3413](https://github.com/chanzuckerberg/happy/issues/3413)) ([13c873f](https://github.com/chanzuckerberg/happy/commit/13c873fe23be3ca694094b9594f85b27d95be89a))

## [1.1.0](https://github.com/chanzuckerberg/happy/compare/happy-cloudfront-v1.0.1...happy-cloudfront-v1.1.0) (2023-11-02)


### Features

* separate configuration of allow methods and allowed cache methods ([#2667](https://github.com/chanzuckerberg/happy/issues/2667)) ([af4518a](https://github.com/chanzuckerberg/happy/commit/af4518ac1cc90096294b2bb1c629e62db2f6b700))

## [1.0.1](https://github.com/chanzuckerberg/happy/compare/happy-cloudfront-v1.0.0...happy-cloudfront-v1.0.1) (2023-10-31)


### Bug Fixes

* redirect HTTP traffic to HTTPS from viewer to CloudFront ([#2665](https://github.com/chanzuckerberg/happy/issues/2665)) ([fd95af9](https://github.com/chanzuckerberg/happy/commit/fd95af94e710d05becd2a769f6afe5c3be2ee532))

## 1.0.0 (2023-10-03)


### Features

* cloudfront added to stack module ([#2487](https://github.com/chanzuckerberg/happy/issues/2487)) ([de3d85e](https://github.com/chanzuckerberg/happy/commit/de3d85e63e5978bc349b86d93270aebe464da866))
