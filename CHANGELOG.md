# Changelog

### [0.3.1](https://github.com/chanzuckerberg/happy/compare/v0.3.0...v0.3.1) (2022-02-07)


### BugFixes

* goreleaser needs full git history ([#49](https://github.com/chanzuckerberg/happy/issues/49)) ([02c706b](https://github.com/chanzuckerberg/happy/commit/02c706baf34db2039f45f9ce12dbdbd5bef31498))

## [0.3.0](https://github.com/chanzuckerberg/happy/compare/v0.2.1...v0.3.0) (2022-02-07)


### ⚠ BREAKING CHANGES

* Enforce coverage

### Features

* Happy CLI Fargate support ([#27](https://github.com/chanzuckerberg/happy/issues/27)) ([39ae2bd](https://github.com/chanzuckerberg/happy/commit/39ae2bd26bbaffd53aab5fa8cd244f07f83d0bb7))
* Happy create/update should fail if the specified tag is missing ([#32](https://github.com/chanzuckerberg/happy/issues/32)) ([0c07556](https://github.com/chanzuckerberg/happy/commit/0c0755638173fee211e454b4fddbc23ead28c7d7))
* Initial pass at consolidating configuration ([#31](https://github.com/chanzuckerberg/happy/issues/31)) ([045f768](https://github.com/chanzuckerberg/happy/commit/045f768e42e241ea3ba20ee065b95723afdf445d))
* Make env configurable ([#40](https://github.com/chanzuckerberg/happy/issues/40)) ([282c815](https://github.com/chanzuckerberg/happy/commit/282c8158719f11f9120ca2e54187201d78558e7d))
* Make happy CLI asks questions when updates takes too long ([ddb53b6](https://github.com/chanzuckerberg/happy/commit/ddb53b6431002b9bbb1e850a36066b92d4c14f0c))
* Performance improvement for happy list commands ([#22](https://github.com/chanzuckerberg/happy/issues/22)) ([cf4405c](https://github.com/chanzuckerberg/happy/commit/cf4405ca569e2583734e680ee26bf6329741bd91))
* search for happy root path in current directory tree if available + configure more lint rules  ([#46](https://github.com/chanzuckerberg/happy/issues/46)) ([e49aa96](https://github.com/chanzuckerberg/happy/commit/e49aa96895f44b08137efb1330fcdf9e24edc290))
* Skip tagging of non existing images ([#21](https://github.com/chanzuckerberg/happy/issues/21)) ([199b435](https://github.com/chanzuckerberg/happy/commit/199b435d71528aff7c2f6ebb50cb072cf55a896e))


### Misc

* fix more lint ([#29](https://github.com/chanzuckerberg/happy/issues/29)) ([12e783b](https://github.com/chanzuckerberg/happy/commit/12e783b2c2ead229bd9b87cc00f2d73a5481c96e))
* reduce number of ignored error cases ([#28](https://github.com/chanzuckerberg/happy/issues/28)) ([75132e6](https://github.com/chanzuckerberg/happy/commit/75132e6d233c7d4e34bddbb7e3ba18d42f7cb734))
* removing circular dependencies in config ([#47](https://github.com/chanzuckerberg/happy/issues/47)) ([c23487c](https://github.com/chanzuckerberg/happy/commit/c23487cd6c825f6fae7b0a5c954d8cbd8a9b5e53))
* Update coverage ([#35](https://github.com/chanzuckerberg/happy/issues/35)) ([c4088cb](https://github.com/chanzuckerberg/happy/commit/c4088cb5c7eef7d46928acd4690d89152a0d6260))
* Write silly tests to make coverage kick in ([#24](https://github.com/chanzuckerberg/happy/issues/24)) ([83b0138](https://github.com/chanzuckerberg/happy/commit/83b0138e9db755e05c4b2a6ede4c4eca5f5c5732))


### BugFixes

* Coverage action name should be upgrade-coverage ([768ae24](https://github.com/chanzuckerberg/happy/commit/768ae24137482969581f566ce31296b11e9ad33a))
* Enforce coverage ([cff2fb4](https://github.com/chanzuckerberg/happy/commit/cff2fb4810c18874cee154cf29648d513dcbfc62))
* Image existence check needs to be skipped when the tag is not specified ([#41](https://github.com/chanzuckerberg/happy/issues/41)) ([359792e](https://github.com/chanzuckerberg/happy/commit/359792e65d5ad338b11e76f985d118ca7b941174))
* Non-numeric tag values break happy list ([#42](https://github.com/chanzuckerberg/happy/issues/42)) ([227978d](https://github.com/chanzuckerberg/happy/commit/227978d5089bb1dadd487db4c4ee025dab363307))
* RegistryId is just the AWS account id, not a registry host name ([#37](https://github.com/chanzuckerberg/happy/issues/37)) ([ebc9144](https://github.com/chanzuckerberg/happy/commit/ebc9144a93701d59b1a2dd4055bcc271a9b949b4))
* Usability improvements; primarily logging and console interaction ([#44](https://github.com/chanzuckerberg/happy/issues/44)) ([6a63982](https://github.com/chanzuckerberg/happy/commit/6a63982f103b1ef85314a03fea0f6a151a6f4bf1))
* Use a more friendly docker tag format ([#48](https://github.com/chanzuckerberg/happy/issues/48)) ([84f5952](https://github.com/chanzuckerberg/happy/commit/84f595230a4f09ff2c17e30da8c0b9f735b87a94))
* Use the initialized tag ([#36](https://github.com/chanzuckerberg/happy/issues/36)) ([c14b3e7](https://github.com/chanzuckerberg/happy/commit/c14b3e7402d0ca92198ed70c5644405e5e172063))

### [0.2.1](https://github.com/chanzuckerberg/happy/compare/v0.2.0...v0.2.1) (2022-01-25)


### BugFixes

* release action needs code and Go installed ([#15](https://github.com/chanzuckerberg/happy/issues/15)) ([f81984b](https://github.com/chanzuckerberg/happy/commit/f81984b6b74cda602791964905e6cdbcbe3d66cf))

## [0.2.0](https://github.com/chanzuckerberg/happy/compare/v0.1.2...v0.2.0) (2022-01-25)


### ⚠ BREAKING CHANGES

* Rename happy-deploy to happy everywhere (#14)
* We call the binary "happy" rather than happy deploy (#9)

### Features

* Get CI to a good place ([#12](https://github.com/chanzuckerberg/happy/issues/12)) ([01b5e73](https://github.com/chanzuckerberg/happy/commit/01b5e739daaf7e79e0bf9a970f3b3268f1f4587c))


### Misc

* Add CODEOWNERS ([#10](https://github.com/chanzuckerberg/happy/issues/10)) ([b555bf9](https://github.com/chanzuckerberg/happy/commit/b555bf9f92f0433569eff14db5c1e0b9728e43a4))
* Rename happy-deploy to happy everywhere ([#14](https://github.com/chanzuckerberg/happy/issues/14)) ([d7794a2](https://github.com/chanzuckerberg/happy/commit/d7794a2fc40d0f83c4324be8cb5c989536e3aa67))
* We call the binary "happy" rather than happy deploy ([#9](https://github.com/chanzuckerberg/happy/issues/9)) ([1355910](https://github.com/chanzuckerberg/happy/commit/13559103b1c3151ac9baf942963af034e11df408))

### [0.1.2](https://www.github.com/chanzuckerberg/happy-deploy/compare/v0.1.1...v0.1.2) (2021-12-15)


### BugFixes

* goreleaser github action trigger ([#7](https://www.github.com/chanzuckerberg/happy-deploy/issues/7)) ([8b7984a](https://www.github.com/chanzuckerberg/happy-deploy/commit/8b7984a9ad7f2996dfba9c7534359984e26f2053))

### [0.1.1](https://www.github.com/chanzuckerberg/happy-deploy/compare/v0.1.0...v0.1.1) (2021-12-15)


### BugFixes

* update goreleaser.yml trigger [#5](https://www.github.com/chanzuckerberg/happy-deploy/issues/5) ([adf93da](https://www.github.com/chanzuckerberg/happy-deploy/commit/adf93da43ff1a833c2725a8b2b2ddf99a15285e3))

## [0.1.0](https://www.github.com/chanzuckerberg/happy-deploy/compare/v0.0.8...v0.1.0) (2021-12-15)


### Features

* Configure release process GitHub actions ([#3](https://www.github.com/chanzuckerberg/happy-deploy/issues/3)) ([988a472](https://www.github.com/chanzuckerberg/happy-deploy/commit/988a4727e6a2baeaf52a9fabbda4c8d210b90f05))
