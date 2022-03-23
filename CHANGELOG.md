# Changelog

### [0.15.1](https://github.com/chanzuckerberg/happy/compare/v0.15.0...v0.15.1) (2022-03-23)


### Misc

* bump github.com/aws/aws-sdk-go from 1.43.21 to 1.43.23 ([#183](https://github.com/chanzuckerberg/happy/issues/183)) ([e9934ce](https://github.com/chanzuckerberg/happy/commit/e9934ce985a23a3f18232c5017be1bb5b58578de))

## [0.15.0](https://github.com/chanzuckerberg/happy/compare/v0.14.0...v0.15.0) (2022-03-21)


### Features

* ensure that stack name inputs to CLI are compatible with DNS character set ([#168](https://github.com/chanzuckerberg/happy/issues/168)) ([f9b030c](https://github.com/chanzuckerberg/happy/commit/f9b030cc7e879798478408eedd3232d9efd0e773))

## [0.14.0](https://github.com/chanzuckerberg/happy/compare/v0.13.0...v0.14.0) (2022-03-21)


### Features

* Added support for streaming of logs, 'get stack', and ability to force update when stack doesn't exist, fixes [#137](https://github.com/chanzuckerberg/happy/issues/137), [#139](https://github.com/chanzuckerberg/happy/issues/139) ([#143](https://github.com/chanzuckerberg/happy/issues/143)) ([71dbf72](https://github.com/chanzuckerberg/happy/commit/71dbf72c500ecda5d477c85692cdb411b1ffc6cc))


### Misc

* bump github.com/aws/aws-sdk-go from 1.43.14 to 1.43.17 ([#169](https://github.com/chanzuckerberg/happy/issues/169)) ([ada9ea5](https://github.com/chanzuckerberg/happy/commit/ada9ea53409771d51d64968ba2753c5256a18a4d))
* bump github.com/aws/aws-sdk-go from 1.43.17 to 1.43.18 ([#171](https://github.com/chanzuckerberg/happy/issues/171)) ([399503c](https://github.com/chanzuckerberg/happy/commit/399503c673ded958dc8a8957a451f096b5bf904b))
* bump github.com/aws/aws-sdk-go from 1.43.18 to 1.43.21 ([#177](https://github.com/chanzuckerberg/happy/issues/177)) ([d43117a](https://github.com/chanzuckerberg/happy/commit/d43117a9b40e70f8ec450d82d5b9fb1b770258e3))
* bump github.com/docker/docker ([#166](https://github.com/chanzuckerberg/happy/issues/166)) ([60eb892](https://github.com/chanzuckerberg/happy/commit/60eb8920c7cb9cdab572e07b226f8c6230cc3fe4))
* bump github.com/spf13/cobra from 1.3.0 to 1.4.0 ([#167](https://github.com/chanzuckerberg/happy/issues/167)) ([bfd1830](https://github.com/chanzuckerberg/happy/commit/bfd18306f233b78c7a1b53868c4e87daa0c4b19a))
* bump github.com/stretchr/testify from 1.7.0 to 1.7.1 ([#173](https://github.com/chanzuckerberg/happy/issues/173)) ([2e9c501](https://github.com/chanzuckerberg/happy/commit/2e9c501da9f5e93b0d362aa0e954d888b51a422d))


### BugFixes

* docker-compose version contraints from panicing ([#162](https://github.com/chanzuckerberg/happy/issues/162)) ([2c3da6c](https://github.com/chanzuckerberg/happy/commit/2c3da6ce4c0864f8b9de2758833d8c27162b552d))
* Refresh existing TFE token via browser instead of always creating new one ([#172](https://github.com/chanzuckerberg/happy/issues/172)) ([b1cd38e](https://github.com/chanzuckerberg/happy/commit/b1cd38ee54fda478cc10ce381bef4146b4710c14))

## [0.13.0](https://github.com/chanzuckerberg/happy/compare/v0.12.0...v0.13.0) (2022-03-09)


### Features

* Pre-releaser github workflow ([#148](https://github.com/chanzuckerberg/happy/issues/148)) ([fcf0ffd](https://github.com/chanzuckerberg/happy/commit/fcf0ffd0f76a143a3ab33d5d4fa11208db24674b))
* Profile the runtime of happy create ([#147](https://github.com/chanzuckerberg/happy/issues/147)) ([0c4315d](https://github.com/chanzuckerberg/happy/commit/0c4315dbb823a1f6839df162f2a78b66d099cd08))


### Misc

* bump github.com/aws/aws-sdk-go from 1.43.12 to 1.43.13 ([#151](https://github.com/chanzuckerberg/happy/issues/151)) ([28febf1](https://github.com/chanzuckerberg/happy/commit/28febf1a6b505594d779c580ea344e6b8cf473b9))
* bump github.com/aws/aws-sdk-go from 1.43.13 to 1.43.14 ([#155](https://github.com/chanzuckerberg/happy/issues/155)) ([399704b](https://github.com/chanzuckerberg/happy/commit/399704bc9b081065bb9db1b090fb41332d4444c4))
* bump github.com/aws/aws-sdk-go from 1.43.8 to 1.43.9 ([#136](https://github.com/chanzuckerberg/happy/issues/136)) ([8372812](https://github.com/chanzuckerberg/happy/commit/83728126672272cfa1a2b23b150b483975a7d184))
* bump github.com/aws/aws-sdk-go from 1.43.9 to 1.43.12 ([#145](https://github.com/chanzuckerberg/happy/issues/145)) ([3f28724](https://github.com/chanzuckerberg/happy/commit/3f28724b836595a205b15a7d287f9138f1d7418c))
* bump github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs ([#150](https://github.com/chanzuckerberg/happy/issues/150)) ([8ab7901](https://github.com/chanzuckerberg/happy/commit/8ab79010a3b5ae375ef191c12ad8137b4edd9945))
* bump github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs ([#154](https://github.com/chanzuckerberg/happy/issues/154)) ([fd48193](https://github.com/chanzuckerberg/happy/commit/fd4819338525b3d1053dee55d2b51fd00f584115))
* bump github.com/go-playground/validator/v10 ([#152](https://github.com/chanzuckerberg/happy/issues/152)) ([875cc5f](https://github.com/chanzuckerberg/happy/commit/875cc5fcddc44b09b1c420896ab4b28e48cb230c))
* bump github.com/gruntwork-io/terratest from 0.40.5 to 0.40.6 ([#146](https://github.com/chanzuckerberg/happy/issues/146)) ([fc8b44d](https://github.com/chanzuckerberg/happy/commit/fc8b44df41cac7630eb1f57d8f60b3185836cd70))
* Update coverage ([#134](https://github.com/chanzuckerberg/happy/issues/134)) ([252fad6](https://github.com/chanzuckerberg/happy/commit/252fad674c5324fb3f6501ab1117b24f8c2f008e))


### BugFixes

* Allow users to specify AWS_PROFILE; override config ([#153](https://github.com/chanzuckerberg/happy/issues/153)) ([b542199](https://github.com/chanzuckerberg/happy/commit/b5421994fe3a06d0c85ad0e730b92d1e7d07a845))
* Empty AWS profile means don't set one ([#160](https://github.com/chanzuckerberg/happy/issues/160)) ([4183357](https://github.com/chanzuckerberg/happy/commit/41833571a30078d018564ef49575931ba23cfec2))
* GitHub Release action s/jobs/needs/ ([#156](https://github.com/chanzuckerberg/happy/issues/156)) ([50a4027](https://github.com/chanzuckerberg/happy/commit/50a4027355b09104af0bdc1894efaba426da49b8))
* prerelease needs github token ([#159](https://github.com/chanzuckerberg/happy/issues/159)) ([b9b0202](https://github.com/chanzuckerberg/happy/commit/b9b02029f7642b098c5884a7ccf38f519932fd3c))
* prerelease s/tag/version/ ([#158](https://github.com/chanzuckerberg/happy/issues/158)) ([302b6f7](https://github.com/chanzuckerberg/happy/commit/302b6f707819a3d6f817961b476d4c29f6cff155))
* Prerelease wasn't using annotated tags ([#157](https://github.com/chanzuckerberg/happy/issues/157)) ([60989ae](https://github.com/chanzuckerberg/happy/commit/60989aef7baa7698a3b12561ad1c80f98752705b))

## [0.12.0](https://github.com/chanzuckerberg/happy/compare/v0.11.1...v0.12.0) (2022-03-01)


### Features

* Display links to installation instructions for missing dependencies, resolves [#128](https://github.com/chanzuckerberg/happy/issues/128) ([#129](https://github.com/chanzuckerberg/happy/issues/129)) ([172d22d](https://github.com/chanzuckerberg/happy/commit/172d22d7303b7ea790ac54a977755773e05871d3))

### [0.11.1](https://github.com/chanzuckerberg/happy/compare/v0.11.0...v0.11.1) (2022-02-28)


### BugFixes

* Fetch cloudwatch logs when running tasks ([#126](https://github.com/chanzuckerberg/happy/issues/126)) ([d5e8c6e](https://github.com/chanzuckerberg/happy/commit/d5e8c6e5c8c845e41232c0eb730af6dab4d9201d))

## [0.11.0](https://github.com/chanzuckerberg/happy/compare/v0.10.0...v0.11.0) (2022-02-28)


### Features

* Refresh TFE token if expired ([#120](https://github.com/chanzuckerberg/happy/issues/120)) ([135b70b](https://github.com/chanzuckerberg/happy/commit/135b70b243d50fbcfd04a4c118715ea5f002c513))

## [0.10.0](https://github.com/chanzuckerberg/happy/compare/v0.9.1...v0.10.0) (2022-02-24)


### ⚠ BREAKING CHANGES

* match existing config key to auto_run_migrations (#118)

### Features

* enable dependabot ([#117](https://github.com/chanzuckerberg/happy/issues/117)) ([ed4c973](https://github.com/chanzuckerberg/happy/commit/ed4c973b8752fbc71a2e2043ca7330c429e858bb))


### BugFixes

* match existing config key to auto_run_migrations ([#118](https://github.com/chanzuckerberg/happy/issues/118)) ([03ad529](https://github.com/chanzuckerberg/happy/commit/03ad529a8428ffde421b3700781d985f72f9c0e7))

### [0.9.1](https://github.com/chanzuckerberg/happy/compare/v0.9.0...v0.9.1) (2022-02-24)


### BugFixes

* typo ([#115](https://github.com/chanzuckerberg/happy/issues/115)) ([a0ef4f8](https://github.com/chanzuckerberg/happy/commit/a0ef4f87a4bc78c1bb589812318b6e9164833824))

## [0.9.0](https://github.com/chanzuckerberg/happy/compare/v0.8.0...v0.9.0) (2022-02-24)


### Features

* show usage when cli cmd validation fails ([#112](https://github.com/chanzuckerberg/happy/issues/112)) ([5da10ed](https://github.com/chanzuckerberg/happy/commit/5da10ede8ad88e4149833f5bbb1641a5c804a46f))


### Misc

* Parallelize CI jobs ([#110](https://github.com/chanzuckerberg/happy/issues/110)) ([37b16b0](https://github.com/chanzuckerberg/happy/commit/37b16b080dc77c322c06c9a219f5c84d695db5d4))

## [0.8.0](https://github.com/chanzuckerberg/happy/compare/v0.7.0...v0.8.0) (2022-02-24)


### Features

* add a version command; make sure we update stack's owner tag ([#108](https://github.com/chanzuckerberg/happy/issues/108)) ([ab17468](https://github.com/chanzuckerberg/happy/commit/ab17468256484482e4fc1ae8fdd61ca027607fcf))

## [0.7.0](https://github.com/chanzuckerberg/happy/compare/v0.6.1...v0.7.0) (2022-02-23)


### Features

* check for docker-compose v2 ([#107](https://github.com/chanzuckerberg/happy/issues/107)) ([bbbf036](https://github.com/chanzuckerberg/happy/commit/bbbf03624928924174fb0835c080b87f3fefb9c0))


### Misc

* Builder pattern and launch conditions validation ([#103](https://github.com/chanzuckerberg/happy/issues/103)) ([c5347c0](https://github.com/chanzuckerberg/happy/commit/c5347c0f303b8bb4957203562a9d2123c02d37a0))
* Improve orchestrator package code coverage ([#101](https://github.com/chanzuckerberg/happy/issues/101)) ([5a23d44](https://github.com/chanzuckerberg/happy/commit/5a23d44258990c6789615e698692dc4c1e04f626))
* Remove unused parameters ([#99](https://github.com/chanzuckerberg/happy/issues/99)) ([eb5b4dd](https://github.com/chanzuckerberg/happy/commit/eb5b4dd38465e97b76f67405a1e48b9880bb7339))
* Update coverage ([#105](https://github.com/chanzuckerberg/happy/issues/105)) ([5e58068](https://github.com/chanzuckerberg/happy/commit/5e58068eaa0ff83261d58066de4e42e3d0d09673))


### BugFixes

* Increased golangci-lint timeout ([#102](https://github.com/chanzuckerberg/happy/issues/102)) ([f7925d4](https://github.com/chanzuckerberg/happy/commit/f7925d4f187becf78ea62086f9b9e1c34b7f1173))
* more consistently use slices; imagetags set properly ([#106](https://github.com/chanzuckerberg/happy/issues/106)) ([9b8e8d2](https://github.com/chanzuckerberg/happy/commit/9b8e8d212a1bdce5130487c263161b20d54a441f))
* Updated offending dependencie (docker/distribution and opencontainers/image-spec) ([#104](https://github.com/chanzuckerberg/happy/issues/104)) ([41cb783](https://github.com/chanzuckerberg/happy/commit/41cb783517ae72c602fe87bb0ea22d60f41c1f71))

### [0.6.1](https://github.com/chanzuckerberg/happy/compare/v0.6.0...v0.6.1) (2022-02-18)


### Misc

* Style cleanup ([#96](https://github.com/chanzuckerberg/happy/issues/96)) ([6febcb4](https://github.com/chanzuckerberg/happy/commit/6febcb42d1d8fab26c4b428b0945d272bc91876b))


### BugFixes

* Create --force now works when partial stack exists ([#98](https://github.com/chanzuckerberg/happy/issues/98)) ([94b4f16](https://github.com/chanzuckerberg/happy/commit/94b4f168824f6554b04eca782bc6f9a5cac61981))
* Unused parameter in retagimages and inconsistent receiver names ([#95](https://github.com/chanzuckerberg/happy/issues/95)) ([1105f70](https://github.com/chanzuckerberg/happy/commit/1105f70d74920ec623052499c9543c130547130f))

## [0.6.0](https://github.com/chanzuckerberg/happy/compare/v0.5.0...v0.6.0) (2022-02-17)


### Features

* Audible alert upon failure ([#89](https://github.com/chanzuckerberg/happy/issues/89)) ([bd23143](https://github.com/chanzuckerberg/happy/commit/bd23143b8f7d5a75eba90b6b56029613696e5048))
* Display image tags on happy push ([#91](https://github.com/chanzuckerberg/happy/issues/91)) ([2466a5e](https://github.com/chanzuckerberg/happy/commit/2466a5ecfa0984e388599830b0cfa9a0c8042bd0))


### BugFixes

* Additional tags are ignored on happy push ([#93](https://github.com/chanzuckerberg/happy/issues/93)) ([2e0dc6e](https://github.com/chanzuckerberg/happy/commit/2e0dc6eecc6984c543314cb6e7bd61093f39d672))

## [0.5.0](https://github.com/chanzuckerberg/happy/compare/v0.4.1...v0.5.0) (2022-02-17)


### ⚠ BREAKING CHANGES

* Reinterpret slices, so they are compatible with docker-compose profiles (#77)

### Features

* colorize output and make it more human-readable ([#82](https://github.com/chanzuckerberg/happy/issues/82)) ([bfe0987](https://github.com/chanzuckerberg/happy/commit/bfe0987ba0e4b271ca740eddbf31cf8136d3af4d))
* Friendlier configuration validation messages [#90](https://github.com/chanzuckerberg/happy/issues/90) ([bfbe3a7](https://github.com/chanzuckerberg/happy/commit/bfbe3a7b220b8c302b1654e0851dab0ce0eb160f))
* Reinterpret slices, so they are compatible with docker-compose profiles ([#77](https://github.com/chanzuckerberg/happy/issues/77)) ([80fea88](https://github.com/chanzuckerberg/happy/commit/80fea88c4f25343f47194059618cc7a7b88a3cf6))


### BugFixes

* list uses tableprinter package and handles errors better ([#87](https://github.com/chanzuckerberg/happy/issues/87)) ([69abb4d](https://github.com/chanzuckerberg/happy/commit/69abb4dc4f131a448092a29b7c5fe8943d8d0153))
* Set owner tag when missing ([#88](https://github.com/chanzuckerberg/happy/issues/88)) ([8d2f13a](https://github.com/chanzuckerberg/happy/commit/8d2f13ab953945cbe788a29e6016c39435d15a51))
* TFE run loop should succeed on applied and  run_plan_and_finished ([#86](https://github.com/chanzuckerberg/happy/issues/86)) ([c11d25f](https://github.com/chanzuckerberg/happy/commit/c11d25fdeb8a01dc9a850bccfe95e4caaff88a15))
* user logger; tfe sentinel status planned_and_finished; delete succeeds/noop if stack not found ([#85](https://github.com/chanzuckerberg/happy/issues/85)) ([d684d90](https://github.com/chanzuckerberg/happy/commit/d684d902d5952a1365ab93ad6a8b257645fd226f))


### Misc

* Code coverage ([#80](https://github.com/chanzuckerberg/happy/issues/80)) ([d838a19](https://github.com/chanzuckerberg/happy/commit/d838a19c6dc8aba6401e36ff098a2d78ced3c3fb))
* Improve code coverage ([#76](https://github.com/chanzuckerberg/happy/issues/76)) ([f796107](https://github.com/chanzuckerberg/happy/commit/f796107c2d8b30867cfbc96d1655272bc27efd53))
* make apply message friendlier [#83](https://github.com/chanzuckerberg/happy/issues/83) ([6227f00](https://github.com/chanzuckerberg/happy/commit/6227f0030e4b330cc63e51b98824a7b3ad54c8d7))
* Orchestrator package coverage improvements ([#84](https://github.com/chanzuckerberg/happy/issues/84)) ([085135d](https://github.com/chanzuckerberg/happy/commit/085135d3ed4058de41d8b1505870cf1e5600452d))
* pin linter version ([#92](https://github.com/chanzuckerberg/happy/issues/92)) ([48236f6](https://github.com/chanzuckerberg/happy/commit/48236f6e0040f912bd0fb96e8509320b153b4cc9))
* Update code coverage for workspace_repo package ([#79](https://github.com/chanzuckerberg/happy/issues/79)) ([84ec051](https://github.com/chanzuckerberg/happy/commit/84ec0514479c1a3a470e9d388139a7e41c5112bd))
* Update coverage ([#81](https://github.com/chanzuckerberg/happy/issues/81)) ([ee6c151](https://github.com/chanzuckerberg/happy/commit/ee6c151c0e4002e7d4cf1f63cea0cb3b67e3457d))

### [0.4.1](https://github.com/chanzuckerberg/happy/compare/v0.4.0...v0.4.1) (2022-02-14)


### BugFixes

* Addressed addtags messaging an made flags required ([#73](https://github.com/chanzuckerberg/happy/issues/73)) ([89ce489](https://github.com/chanzuckerberg/happy/commit/89ce489e4d0fa8b01ff3855fda4f24720d77383c))
* Genepi requires imagetags tag to be a valid json even if no tags are present ([#74](https://github.com/chanzuckerberg/happy/issues/74)) ([8259d70](https://github.com/chanzuckerberg/happy/commit/8259d7025c2c9032094f6979dd03a29bd91e2b91))
* Happy migrate fails because task definition is not specified and subnets info is incorrect ([#69](https://github.com/chanzuckerberg/happy/issues/69)) ([020ff7b](https://github.com/chanzuckerberg/happy/commit/020ff7b2038f9b0b070bd72e6a893306b514166b))
* Happy shell isn't working for czgenepi ([#72](https://github.com/chanzuckerberg/happy/issues/72)) ([489a87e](https://github.com/chanzuckerberg/happy/commit/489a87ea6584297445eefa8a1279b3ff560370dd))
* Network configuration didn't have a complete list of subnets and security groups ([#71](https://github.com/chanzuckerberg/happy/issues/71)) ([8c86f75](https://github.com/chanzuckerberg/happy/commit/8c86f75d5dbcad3f68f8e4b801b183fa74af5b02))


### Misc

* Improve code coverage for orchestrator and hostname_manager ([#75](https://github.com/chanzuckerberg/happy/issues/75)) ([e75de18](https://github.com/chanzuckerberg/happy/commit/e75de18cac4f6dbe5b40795c550c75251b905eea))

## [0.4.0](https://github.com/chanzuckerberg/happy/compare/v0.3.1...v0.4.0) (2022-02-10)


### ⚠ BREAKING CHANGES

* For clarity, default_compose_env setting has been superseded by default_compose_env_file

### Features

* Discovery of docker compose env files, absolute and relative ([#56](https://github.com/chanzuckerberg/happy/issues/56)) ([7f19d69](https://github.com/chanzuckerberg/happy/commit/7f19d6927065d555084acd97550e95bfd45410c2))
* Read terraform token from env var, tfrc file, or prompt terraform login ([#58](https://github.com/chanzuckerberg/happy/issues/58)) ([e599e8e](https://github.com/chanzuckerberg/happy/commit/e599e8e9707b26e1d3cd0dc6baf08122bb1a7a5b))
* Switched to docker compose v2 ([#60](https://github.com/chanzuckerberg/happy/issues/60)) ([cf5dcad](https://github.com/chanzuckerberg/happy/commit/cf5dcad9cf8dbffd02ad536b51d7eb7b9d63b60b))


### BugFixes

* AWS Backend set default AWS profile ([#61](https://github.com/chanzuckerberg/happy/issues/61)) ([b9788d2](https://github.com/chanzuckerberg/happy/commit/b9788d27f8329a31ae710ac91816fc70f1331d20))
* Docker tag cannot have an @ sign present ([#64](https://github.com/chanzuckerberg/happy/issues/64)) ([6ff1a5a](https://github.com/chanzuckerberg/happy/commit/6ff1a5a552e932db7a1bbac212cf387839159ceb))
* happy hosts install breaking because of incorrect type casting ([#68](https://github.com/chanzuckerberg/happy/issues/68)) ([6673c0d](https://github.com/chanzuckerberg/happy/commit/6673c0daea588b1a4c1ee54799665c3b07ebdcf4))
* Implement global dockerComposeEnvFile setting with the default fallback ([#55](https://github.com/chanzuckerberg/happy/issues/55)) ([9c1bd78](https://github.com/chanzuckerberg/happy/commit/9c1bd783d279e103fe5394ecd86c9edd61156dee))
* Split env and composeEnv for clarity ([#51](https://github.com/chanzuckerberg/happy/issues/51)) ([140f643](https://github.com/chanzuckerberg/happy/commit/140f643fd138c75eae7e19a1b56ce599c9b4b498))
* testbackend package to make testing the backend easier ([#66](https://github.com/chanzuckerberg/happy/issues/66)) ([a34ccc9](https://github.com/chanzuckerberg/happy/commit/a34ccc93b840794ab734f7970f63c6815cdc383f))
* various: TFE url sanitize; docker login ecr registries; integration secret parsing ([#62](https://github.com/chanzuckerberg/happy/issues/62)) ([4f1b166](https://github.com/chanzuckerberg/happy/commit/4f1b166678ed4f19637c7f2aa326041a3e067510))
* Verify aws profile exists when creating the backend ([#63](https://github.com/chanzuckerberg/happy/issues/63)) ([7b6689e](https://github.com/chanzuckerberg/happy/commit/7b6689e96a533db553e46b8e6ec153a0980caba6))
* workspace_repo tests and coverage ([#65](https://github.com/chanzuckerberg/happy/issues/65)) ([89d85c9](https://github.com/chanzuckerberg/happy/commit/89d85c909a27f946fe4b17b87e3b8a7985ecc022))


### Misc

* Added semantic clarity to GetServiceRegistries() method ([#59](https://github.com/chanzuckerberg/happy/issues/59)) ([472572e](https://github.com/chanzuckerberg/happy/commit/472572e466fe5968f64178b2fba473722f880e4a))
* Combined docker compose invokations ([#53](https://github.com/chanzuckerberg/happy/issues/53)) ([c373542](https://github.com/chanzuckerberg/happy/commit/c3735427be4ed184af36e741bda048b463bc179d))
* refactor backend to make it easier to work with and test ([#54](https://github.com/chanzuckerberg/happy/issues/54)) ([bef351b](https://github.com/chanzuckerberg/happy/commit/bef351b6a672d706ae1e1034e01f43efb536674d))
* Silence CLI Usage on errors ([#67](https://github.com/chanzuckerberg/happy/issues/67)) ([678b448](https://github.com/chanzuckerberg/happy/commit/678b4485cc63e68586e067922fccd22acd591ba6))
* Update coverage ([#57](https://github.com/chanzuckerberg/happy/issues/57)) ([4caf025](https://github.com/chanzuckerberg/happy/commit/4caf02508d82d6cea2eccab8c426b82269312479))

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

* goreleaser GitHub action trigger ([#7](https://www.github.com/chanzuckerberg/happy-deploy/issues/7)) ([8b7984a](https://www.github.com/chanzuckerberg/happy-deploy/commit/8b7984a9ad7f2996dfba9c7534359984e26f2053))

### [0.1.1](https://www.github.com/chanzuckerberg/happy-deploy/compare/v0.1.0...v0.1.1) (2021-12-15)


### BugFixes

* update goreleaser.yml trigger [#5](https://www.github.com/chanzuckerberg/happy-deploy/issues/5) ([adf93da](https://www.github.com/chanzuckerberg/happy-deploy/commit/adf93da43ff1a833c2725a8b2b2ddf99a15285e3))

## [0.1.0](https://www.github.com/chanzuckerberg/happy-deploy/compare/v0.0.8...v0.1.0) (2021-12-15)


### Features

* Configure release process GitHub actions ([#3](https://www.github.com/chanzuckerberg/happy-deploy/issues/3)) ([988a472](https://www.github.com/chanzuckerberg/happy-deploy/commit/988a4727e6a2baeaf52a9fabbda4c8d210b90f05))
