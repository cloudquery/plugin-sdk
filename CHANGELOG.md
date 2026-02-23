# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [4.94.3](https://github.com/cloudquery/plugin-sdk/compare/v4.94.2...v4.94.3) (2026-02-23)


### Bug Fixes

* Better Go build flags for plugins packaging ([#2415](https://github.com/cloudquery/plugin-sdk/issues/2415)) ([934226b](https://github.com/cloudquery/plugin-sdk/commit/934226b27bb83c4d7f411678c3887ffafc621329))
* **deps:** Update module github.com/cloudquery/codegen to v0.3.37 ([#2410](https://github.com/cloudquery/plugin-sdk/issues/2410)) ([c812b39](https://github.com/cloudquery/plugin-sdk/commit/c812b39c2e3760d4eb702be3edc89dc870b696f9))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.27.7 ([#2414](https://github.com/cloudquery/plugin-sdk/issues/2414)) ([51f4d60](https://github.com/cloudquery/plugin-sdk/commit/51f4d60ec1dc0d8f67c042ceef67a3247ae7b6e3))
* **deps:** Update module google.golang.org/grpc to v1.79.1 ([#2412](https://github.com/cloudquery/plugin-sdk/issues/2412)) ([f617806](https://github.com/cloudquery/plugin-sdk/commit/f6178065380e07ccf62e3fbca4c35f4e86af33fc))

## [4.94.2](https://github.com/cloudquery/plugin-sdk/compare/v4.94.1...v4.94.2) (2026-02-06)


### Bug Fixes

* **deps:** Update module github.com/cloudquery/codegen to v0.3.36 ([#2403](https://github.com/cloudquery/plugin-sdk/issues/2403)) ([b7188f1](https://github.com/cloudquery/plugin-sdk/commit/b7188f174bcb1286ccd1194d85d7268fb1ecb4ce))
* Fix race condition on batchsender when sending resources from multiple goroutines ([#2405](https://github.com/cloudquery/plugin-sdk/issues/2405)) ([a0e2801](https://github.com/cloudquery/plugin-sdk/commit/a0e28013ec8178bad7e123cd86c1f7701bc27f86))

## [4.94.1](https://github.com/cloudquery/plugin-sdk/compare/v4.94.0...v4.94.1) (2026-02-02)


### Bug Fixes

* **deps:** Update aws-sdk-go-v2 monorepo ([#2390](https://github.com/cloudquery/plugin-sdk/issues/2390)) ([f50c14c](https://github.com/cloudquery/plugin-sdk/commit/f50c14c3b6bb1b9b1314dc013cb4e78ad5327707))
* **deps:** Update golang.org/x/exp digest to 716be56 ([#2395](https://github.com/cloudquery/plugin-sdk/issues/2395)) ([3a7a913](https://github.com/cloudquery/plugin-sdk/commit/3a7a913693bb9ec81ca30135cc4c5f5d03e25e2e))
* **deps:** Update module github.com/apache/arrow-go/v18 to v18.5.1 ([#2400](https://github.com/cloudquery/plugin-sdk/issues/2400)) ([78985aa](https://github.com/cloudquery/plugin-sdk/commit/78985aa956ac1c0e524b897cf4881d28c26560d9))
* **deps:** Update module github.com/aws/aws-sdk-go-v2/service/licensemanager to v1.37.6 ([#2393](https://github.com/cloudquery/plugin-sdk/issues/2393)) ([5b129f8](https://github.com/cloudquery/plugin-sdk/commit/5b129f8df656fdb4f5e78d6e93136cb279c8be3f))
* **deps:** Update module github.com/aws/aws-sdk-go-v2/service/marketplacemetering to v1.35.6 ([#2394](https://github.com/cloudquery/plugin-sdk/issues/2394)) ([cf0438b](https://github.com/cloudquery/plugin-sdk/commit/cf0438b0840804c47430dec43c1ed2ef9fa2f57a))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.27.6 ([#2401](https://github.com/cloudquery/plugin-sdk/issues/2401)) ([1cc9ada](https://github.com/cloudquery/plugin-sdk/commit/1cc9adad97bff7ccffaea7dcc6015aaa73f91418))
* **deps:** Update module github.com/getsentry/sentry-go to v0.41.0 ([#2397](https://github.com/cloudquery/plugin-sdk/issues/2397)) ([907b607](https://github.com/cloudquery/plugin-sdk/commit/907b607408f1f2279bd0957ed467a327d612c6ff))

## [4.94.0](https://github.com/cloudquery/plugin-sdk/compare/v4.93.1...v4.94.0) (2026-01-12)


### Features

* Add optional Sentry support ([#2386](https://github.com/cloudquery/plugin-sdk/issues/2386)) ([dab2ce8](https://github.com/cloudquery/plugin-sdk/commit/dab2ce8c2e7545fa3da3fce1b9d9c1c9ca057c7f))


### Bug Fixes

* **deps:** Update aws-sdk-go-v2 monorepo ([#2371](https://github.com/cloudquery/plugin-sdk/issues/2371)) ([2d1165a](https://github.com/cloudquery/plugin-sdk/commit/2d1165a6f19a799890a8d7ebf67b142b512754d2))
* **deps:** Update golang.org/x/exp digest to 944ab1f ([#2376](https://github.com/cloudquery/plugin-sdk/issues/2376)) ([4da0fd0](https://github.com/cloudquery/plugin-sdk/commit/4da0fd0c31e864161e042d0443c94d9e75113bfb))
* **deps:** Update module github.com/aws/aws-sdk-go-v2/config to v1.32.6 ([#2374](https://github.com/cloudquery/plugin-sdk/issues/2374)) ([0f0d985](https://github.com/cloudquery/plugin-sdk/commit/0f0d985b076731a74a4d82503108d6ecce61fc8a))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.14.8 ([#2378](https://github.com/cloudquery/plugin-sdk/issues/2378)) ([b3e0ae7](https://github.com/cloudquery/plugin-sdk/commit/b3e0ae7b3b26ac2df4eb57417a9dbb020485a641))
* **deps:** Update module github.com/cloudquery/codegen to v0.3.35 ([#2381](https://github.com/cloudquery/plugin-sdk/issues/2381)) ([2daf99c](https://github.com/cloudquery/plugin-sdk/commit/2daf99c00ba2b9cf359a1e08cf715c4d7fa6f859))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.27.3 ([#2380](https://github.com/cloudquery/plugin-sdk/issues/2380)) ([fb04081](https://github.com/cloudquery/plugin-sdk/commit/fb040815a46fa35b52f47570223d000cab82d63e))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.27.4 ([#2382](https://github.com/cloudquery/plugin-sdk/issues/2382)) ([9ba1f1f](https://github.com/cloudquery/plugin-sdk/commit/9ba1f1fba501d3e10a0adb9fc77d91785c66824c))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.27.5 ([#2385](https://github.com/cloudquery/plugin-sdk/issues/2385)) ([0eee0ed](https://github.com/cloudquery/plugin-sdk/commit/0eee0ed757b25aea0eac949a86f0cfda34cba1be))
* **deps:** Update module google.golang.org/grpc to v1.78.0 ([#2383](https://github.com/cloudquery/plugin-sdk/issues/2383)) ([815d332](https://github.com/cloudquery/plugin-sdk/commit/815d332ae113cc493b29f7b38f82dbf0e947d67e))

## [4.93.1](https://github.com/cloudquery/plugin-sdk/compare/v4.93.0...v4.93.1) (2025-12-16)


### Bug Fixes

* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.14.7 ([#2365](https://github.com/cloudquery/plugin-sdk/issues/2365)) ([1953885](https://github.com/cloudquery/plugin-sdk/commit/19538857776889bba19f91f5acbaa103703768b4))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.27.1 ([#2367](https://github.com/cloudquery/plugin-sdk/issues/2367)) ([985d54a](https://github.com/cloudquery/plugin-sdk/commit/985d54ace1f535e9c8afb7f6f9e87e6921fe9ef2))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.27.2 ([#2368](https://github.com/cloudquery/plugin-sdk/issues/2368)) ([d54cb52](https://github.com/cloudquery/plugin-sdk/commit/d54cb521324ded7593cf9add421bb38c280166b0))
* **deps:** Update module github.com/spf13/cobra to v1.10.2 ([#2353](https://github.com/cloudquery/plugin-sdk/issues/2353)) ([8c37fb4](https://github.com/cloudquery/plugin-sdk/commit/8c37fb478994612422dc4b83ef7eec10a249e019))
* **deps:** Update module golang.org/x/oauth2 to v0.34.0 ([#2364](https://github.com/cloudquery/plugin-sdk/issues/2364)) ([d8a6cff](https://github.com/cloudquery/plugin-sdk/commit/d8a6cffc93620ce623e809451b9a659ccf649f80))
* **deps:** Update module golang.org/x/sync to v0.19.0 ([#2355](https://github.com/cloudquery/plugin-sdk/issues/2355)) ([72f27ac](https://github.com/cloudquery/plugin-sdk/commit/72f27ac9573dcf3e00ecfe557b8519bac6d70610))

## [4.93.0](https://github.com/cloudquery/plugin-sdk/compare/v4.92.1...v4.93.0) (2025-12-12)


### Features

* Update OTEL to latest version ([#2349](https://github.com/cloudquery/plugin-sdk/issues/2349)) ([d4ee1ab](https://github.com/cloudquery/plugin-sdk/commit/d4ee1ab7d7e5f837068f9ed3b077c694856d7dff))


### Bug Fixes

* **deps:** Update aws-sdk-go-v2 monorepo ([#2343](https://github.com/cloudquery/plugin-sdk/issues/2343)) ([e1aaee1](https://github.com/cloudquery/plugin-sdk/commit/e1aaee11d33a6c868fcf1e04497690819a3608e5))
* **deps:** Update aws-sdk-go-v2 monorepo ([#2347](https://github.com/cloudquery/plugin-sdk/issues/2347)) ([81598ca](https://github.com/cloudquery/plugin-sdk/commit/81598cad9c5b011dd76c7e20e0ad2ea2b79a69cc))
* **deps:** Update opentelemetry-go monorepo ([#2348](https://github.com/cloudquery/plugin-sdk/issues/2348)) ([3a7f6d5](https://github.com/cloudquery/plugin-sdk/commit/3a7f6d59a0fa544a495179e41f454700d3aa6994))

## [4.92.1](https://github.com/cloudquery/plugin-sdk/compare/v4.92.0...v4.92.1) (2025-12-05)


### Bug Fixes

* **deps:** Update aws-sdk-go-v2 monorepo ([#2324](https://github.com/cloudquery/plugin-sdk/issues/2324)) ([681c70f](https://github.com/cloudquery/plugin-sdk/commit/681c70f615659baa4e3a6ca6d96ad97d98b7d926))
* **deps:** Update aws-sdk-go-v2 monorepo ([#2327](https://github.com/cloudquery/plugin-sdk/issues/2327)) ([36789b9](https://github.com/cloudquery/plugin-sdk/commit/36789b9829a53f7a04294b6f2f21b00d0a4be3ef))
* **deps:** Update aws-sdk-go-v2 monorepo ([#2330](https://github.com/cloudquery/plugin-sdk/issues/2330)) ([9316597](https://github.com/cloudquery/plugin-sdk/commit/9316597f8140c976a8f2a238350f78be63f4ec74))
* **deps:** Update aws-sdk-go-v2 monorepo ([#2332](https://github.com/cloudquery/plugin-sdk/issues/2332)) ([5f05015](https://github.com/cloudquery/plugin-sdk/commit/5f050154219dd7409e6f6b5ac0e08ef0365853d6))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.14.6 ([#2339](https://github.com/cloudquery/plugin-sdk/issues/2339)) ([5d84688](https://github.com/cloudquery/plugin-sdk/commit/5d846884c345d5acc5efd304fe18876b020dada8))
* **deps:** Update module github.com/grpc-ecosystem/go-grpc-middleware/v2 to v2.3.3 ([#2334](https://github.com/cloudquery/plugin-sdk/issues/2334)) ([10231d3](https://github.com/cloudquery/plugin-sdk/commit/10231d34221340423edb18ab771c699a1266690c))
* **deps:** Update module golang.org/x/oauth2 to v0.33.0 ([#2336](https://github.com/cloudquery/plugin-sdk/issues/2336)) ([1a46188](https://github.com/cloudquery/plugin-sdk/commit/1a46188f0fb38516258c40ff9b9832e73b0b1c79))
* **deps:** Update module golang.org/x/sync to v0.18.0 ([#2337](https://github.com/cloudquery/plugin-sdk/issues/2337)) ([4418937](https://github.com/cloudquery/plugin-sdk/commit/441893708695aae2448dda45deac1ab9eba5f150))
* **deps:** Update module google.golang.org/grpc to v1.77.0 ([#2333](https://github.com/cloudquery/plugin-sdk/issues/2333)) ([b6fc293](https://github.com/cloudquery/plugin-sdk/commit/b6fc293082e9b964e85e7841bf9f24f495ce86e0))
* Validate all rows in a record ([#2341](https://github.com/cloudquery/plugin-sdk/issues/2341)) ([17285b9](https://github.com/cloudquery/plugin-sdk/commit/17285b9195ace5a00083210c5bf9bb6301568cfd))

## [4.92.0](https://github.com/cloudquery/plugin-sdk/compare/v4.91.0...v4.92.0) (2025-11-06)


### Features

* Support chunks in resource resolvers ([#2287](https://github.com/cloudquery/plugin-sdk/issues/2287)) ([087ef9a](https://github.com/cloudquery/plugin-sdk/commit/087ef9a8cc52b1e671eea5169acb273a685e609f))


### Bug Fixes

* **deps:** Update aws-sdk-go-v2 monorepo ([#2316](https://github.com/cloudquery/plugin-sdk/issues/2316)) ([828c4c2](https://github.com/cloudquery/plugin-sdk/commit/828c4c2d40d1c70b595e33a00c2d9e8a8d3c3fbc))
* **deps:** Update golang.org/x/exp digest to a4bb9ff ([#2315](https://github.com/cloudquery/plugin-sdk/issues/2315)) ([4bdaf94](https://github.com/cloudquery/plugin-sdk/commit/4bdaf94045c4f90732fa2b8e29d377d40e5510cc))
* **deps:** Update module github.com/cloudquery/codegen to v0.3.33 ([#2321](https://github.com/cloudquery/plugin-sdk/issues/2321)) ([8080660](https://github.com/cloudquery/plugin-sdk/commit/8080660ab4845c1f6a2a330491c724884925b6e2))
* **deps:** Update module github.com/samber/lo to v1.52.0 ([#2318](https://github.com/cloudquery/plugin-sdk/issues/2318)) ([6e4c424](https://github.com/cloudquery/plugin-sdk/commit/6e4c424b2cf6611e32d91cdb256c97410d49f4a1))
* **deps:** Update module golang.org/x/oauth2 to v0.32.0 ([#2319](https://github.com/cloudquery/plugin-sdk/issues/2319)) ([48e0e01](https://github.com/cloudquery/plugin-sdk/commit/48e0e01dce96a4662957d654905ae79549d3ae95))

## [4.91.0](https://github.com/cloudquery/plugin-sdk/compare/v4.90.0...v4.91.0) (2025-10-31)


### Features

* Expose transformer option for Skipping PK Validation ([#2312](https://github.com/cloudquery/plugin-sdk/issues/2312)) ([d43e7cc](https://github.com/cloudquery/plugin-sdk/commit/d43e7cc2f23a7e4b86319209b798bb717542c168))

## [4.90.0](https://github.com/cloudquery/plugin-sdk/compare/v4.89.1...v4.90.0) (2025-10-31)


### Features

* Support skipping PK validation for columns ([#2310](https://github.com/cloudquery/plugin-sdk/issues/2310)) ([5ea21e3](https://github.com/cloudquery/plugin-sdk/commit/5ea21e3486065c2ffe3010bda883c3c432c3e40d))


### Bug Fixes

* **deps:** Update aws-sdk-go-v2 monorepo ([#2277](https://github.com/cloudquery/plugin-sdk/issues/2277)) ([ff1b11a](https://github.com/cloudquery/plugin-sdk/commit/ff1b11a2aa13739c12cd0f57b217f9de0b80022f))
* **deps:** Update aws-sdk-go-v2 monorepo ([#2286](https://github.com/cloudquery/plugin-sdk/issues/2286)) ([f5de4f4](https://github.com/cloudquery/plugin-sdk/commit/f5de4f4758d902e29bdddb80caf86e31a962e14c))
* **deps:** Update aws-sdk-go-v2 monorepo ([#2298](https://github.com/cloudquery/plugin-sdk/issues/2298)) ([c1edab7](https://github.com/cloudquery/plugin-sdk/commit/c1edab76e81e05324647573f6a8428fb519f9a7a))
* **deps:** Update aws-sdk-go-v2 monorepo ([#2308](https://github.com/cloudquery/plugin-sdk/issues/2308)) ([bc9b670](https://github.com/cloudquery/plugin-sdk/commit/bc9b670bd24ee9f2f1c7b5fa9fdbc26ba1af2ebf))
* **deps:** Update golang.org/x/exp digest to 8b4c13b ([#2271](https://github.com/cloudquery/plugin-sdk/issues/2271)) ([79eb8e1](https://github.com/cloudquery/plugin-sdk/commit/79eb8e1c824bdd20ae4ffd6d744a02c8aa1f6481))
* **deps:** Update golang.org/x/exp digest to df92998 ([#2291](https://github.com/cloudquery/plugin-sdk/issues/2291)) ([0708ef5](https://github.com/cloudquery/plugin-sdk/commit/0708ef5eba4c783ad2e1c43bbe0683d94b93fe74))
* **deps:** Update Google Golang modules ([#2278](https://github.com/cloudquery/plugin-sdk/issues/2278)) ([3a714e8](https://github.com/cloudquery/plugin-sdk/commit/3a714e852a179cf25a97c8867c4609f747de0faf))
* **deps:** Update Google Golang modules ([#2285](https://github.com/cloudquery/plugin-sdk/issues/2285)) ([162e001](https://github.com/cloudquery/plugin-sdk/commit/162e001c8d4a40452d1ea011a39417b58c32acc3))
* **deps:** Update module github.com/aws/aws-sdk-go-v2/config to v1.31.12 ([#2301](https://github.com/cloudquery/plugin-sdk/issues/2301)) ([3db0e41](https://github.com/cloudquery/plugin-sdk/commit/3db0e4153de5de59eadcb45e06651033b96bebc0))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.14.2 ([#2275](https://github.com/cloudquery/plugin-sdk/issues/2275)) ([6be5fe6](https://github.com/cloudquery/plugin-sdk/commit/6be5fe6411812cb59cbdd0d3b8c03361a34f3e71))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.14.3 ([#2283](https://github.com/cloudquery/plugin-sdk/issues/2283)) ([1c7c102](https://github.com/cloudquery/plugin-sdk/commit/1c7c10259d224426895209321141a7a2383b33eb))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.14.4 ([#2284](https://github.com/cloudquery/plugin-sdk/issues/2284)) ([2c86149](https://github.com/cloudquery/plugin-sdk/commit/2c86149e86ae4e3bb82ab9d106701cde988fd8ec))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.14.5 ([#2297](https://github.com/cloudquery/plugin-sdk/issues/2297)) ([b45ad80](https://github.com/cloudquery/plugin-sdk/commit/b45ad80b825fc8f1de8bc92402588ce19e1690d5))
* **deps:** Update module github.com/cloudquery/codegen to v0.3.32 ([#2288](https://github.com/cloudquery/plugin-sdk/issues/2288)) ([7a64b2d](https://github.com/cloudquery/plugin-sdk/commit/7a64b2d7807965254577b3fc8a1650f6cd47c519))
* **deps:** Update module github.com/spf13/cobra to v1.10.1 ([#2295](https://github.com/cloudquery/plugin-sdk/issues/2295)) ([9f9bb1b](https://github.com/cloudquery/plugin-sdk/commit/9f9bb1b223e5c8f941f625e7bd70941ca911385a))
* **deps:** Update module github.com/stretchr/testify to v1.11.1 ([#2292](https://github.com/cloudquery/plugin-sdk/issues/2292)) ([767ae65](https://github.com/cloudquery/plugin-sdk/commit/767ae6569d8b976814e80bee202dba5a26048a65))
* **deps:** Update module golang.org/x/oauth2 to v0.31.0 ([#2296](https://github.com/cloudquery/plugin-sdk/issues/2296)) ([c07ee17](https://github.com/cloudquery/plugin-sdk/commit/c07ee17da339b50708e039d64c6a38faac699c2d))
* **deps:** Update module google.golang.org/grpc to v1.76.0 ([#2304](https://github.com/cloudquery/plugin-sdk/issues/2304)) ([7abe806](https://github.com/cloudquery/plugin-sdk/commit/7abe806a13ee701664cfe9031e9b21ca5e61a706))
* **deps:** Update module google.golang.org/protobuf to v1.36.10 ([#2302](https://github.com/cloudquery/plugin-sdk/issues/2302)) ([21637c4](https://github.com/cloudquery/plugin-sdk/commit/21637c44acfd0948ac38293c3c9758b6fc512395))

## [4.89.1](https://github.com/cloudquery/plugin-sdk/compare/v4.89.0...v4.89.1) (2025-09-01)


### Bug Fixes

* **deps:** Update aws-sdk-go-v2 monorepo ([#2262](https://github.com/cloudquery/plugin-sdk/issues/2262)) ([b8cf390](https://github.com/cloudquery/plugin-sdk/commit/b8cf3902daef8d1bfbde5d7c7908c70f8085c0a6))
* **deps:** Update module github.com/samber/lo to v1.51.0 ([#2265](https://github.com/cloudquery/plugin-sdk/issues/2265)) ([ea8ca00](https://github.com/cloudquery/plugin-sdk/commit/ea8ca0036e5b447ad5c9af677a318dc52ca69296))
* **deps:** Update module github.com/stretchr/testify to v1.11.0 ([#2266](https://github.com/cloudquery/plugin-sdk/issues/2266)) ([691cf32](https://github.com/cloudquery/plugin-sdk/commit/691cf32340d5aa4bfa0b3c939957e0c8af7110ab))
* **deps:** Update module golang.org/x/text to v0.28.0 ([#2267](https://github.com/cloudquery/plugin-sdk/issues/2267)) ([2e436ed](https://github.com/cloudquery/plugin-sdk/commit/2e436edccb2adbb8ce5c828b4bfed0519dc831d8))

## [4.89.0](https://github.com/cloudquery/plugin-sdk/compare/v4.88.1...v4.89.0) (2025-08-19)


### Features

* Make getColumnChangeSummary public for use in plugins ([#2260](https://github.com/cloudquery/plugin-sdk/issues/2260)) ([5c0a06e](https://github.com/cloudquery/plugin-sdk/commit/5c0a06ef60a1ee576b4222e3eb5cac8aad338c36))


### Bug Fixes

* **deps:** Update aws-sdk-go-v2 monorepo ([#2256](https://github.com/cloudquery/plugin-sdk/issues/2256)) ([0b64895](https://github.com/cloudquery/plugin-sdk/commit/0b64895bac95452bc6f9ae00618050b9201eef65))
* **deps:** Update aws-sdk-go-v2 monorepo ([#2257](https://github.com/cloudquery/plugin-sdk/issues/2257)) ([b546297](https://github.com/cloudquery/plugin-sdk/commit/b5462974de129450f50b0333830810203c831620))
* **deps:** Update golang.org/x/exp digest to 645b1fa ([#2250](https://github.com/cloudquery/plugin-sdk/issues/2250)) ([ab469c1](https://github.com/cloudquery/plugin-sdk/commit/ab469c12bd91a322f189426353da8449282c7ac1))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.14.0 ([#2246](https://github.com/cloudquery/plugin-sdk/issues/2246)) ([291b0d9](https://github.com/cloudquery/plugin-sdk/commit/291b0d9aff6ad1a3b9aab8b4764463937bc4f59b))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.14.1 ([#2254](https://github.com/cloudquery/plugin-sdk/issues/2254)) ([fb148a1](https://github.com/cloudquery/plugin-sdk/commit/fb148a1e180aa94ad0fa816e0a88d2fef4e27623))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.26.16 ([#2245](https://github.com/cloudquery/plugin-sdk/issues/2245)) ([5223700](https://github.com/cloudquery/plugin-sdk/commit/5223700c37b646825c7f2e41c46f07c61b3b4a31))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.26.17 ([#2248](https://github.com/cloudquery/plugin-sdk/issues/2248)) ([3b8a166](https://github.com/cloudquery/plugin-sdk/commit/3b8a1668687bb1b471c5f15a1f6e23ebfb78aa05))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.26.18 ([#2249](https://github.com/cloudquery/plugin-sdk/issues/2249)) ([8e35433](https://github.com/cloudquery/plugin-sdk/commit/8e354333c60518df08bfd17d8c0961d55acd170e))
* **deps:** Update module github.com/grpc-ecosystem/go-grpc-middleware/v2 to v2.3.2 ([#2251](https://github.com/cloudquery/plugin-sdk/issues/2251)) ([a76d1c2](https://github.com/cloudquery/plugin-sdk/commit/a76d1c2c34e4c48c4fc38dcef3960365c9639c1f))
* **deps:** Update module github.com/spf13/cobra to v1.9.1 ([#2252](https://github.com/cloudquery/plugin-sdk/issues/2252)) ([3db1576](https://github.com/cloudquery/plugin-sdk/commit/3db1576eda360fcf9eeac7674295f54c0822cf92))
* **deps:** Update module google.golang.org/grpc to v1.74.2 ([#2255](https://github.com/cloudquery/plugin-sdk/issues/2255)) ([5d8368f](https://github.com/cloudquery/plugin-sdk/commit/5d8368fa82ea32bfc46d361dfed4b58f4cc9701b))
* **deps:** Update module google.golang.org/protobuf to v1.36.7 ([#2258](https://github.com/cloudquery/plugin-sdk/issues/2258)) ([cd611d3](https://github.com/cloudquery/plugin-sdk/commit/cd611d350e9ce1b81461142f8262fa5e4f68dae6))
* Require Row count to be greater than 0 ([#2259](https://github.com/cloudquery/plugin-sdk/issues/2259)) ([8721bdc](https://github.com/cloudquery/plugin-sdk/commit/8721bdc9d861ccd7f21e7e1a67658a3c37acce24))

## [4.88.1](https://github.com/cloudquery/plugin-sdk/compare/v4.88.0...v4.88.1) (2025-07-29)


### Bug Fixes

* **deps:** Update opentelemetry-go monorepo ([#2242](https://github.com/cloudquery/plugin-sdk/issues/2242)) ([71468d0](https://github.com/cloudquery/plugin-sdk/commit/71468d054d7d3a59c5e20f3f92c1b5f3cce78592))

## [4.88.0](https://github.com/cloudquery/plugin-sdk/compare/v4.87.4...v4.88.0) (2025-07-28)


### Features

* Add better changes summary helper ([#2240](https://github.com/cloudquery/plugin-sdk/issues/2240)) ([529532d](https://github.com/cloudquery/plugin-sdk/commit/529532d0555243f88727850997dc86849fb6f993))


### Bug Fixes

* **deps:** Update module google.golang.org/grpc to v1.74.1 ([#2238](https://github.com/cloudquery/plugin-sdk/issues/2238)) ([eeed6a6](https://github.com/cloudquery/plugin-sdk/commit/eeed6a681af717dc5c96e5564c4fb1c4b2a4d2a8))

## [4.87.4](https://github.com/cloudquery/plugin-sdk/compare/v4.87.3...v4.87.4) (2025-07-23)


### Bug Fixes

* **deps:** Update module github.com/apache/arrow-go/v18 to v18.4.0 ([#2234](https://github.com/cloudquery/plugin-sdk/issues/2234)) ([3955c1d](https://github.com/cloudquery/plugin-sdk/commit/3955c1d4d2b3c8827a75ba99b713f40948b3adb2))
* Don't lose IP data in `AppendValueFromString` ([#2236](https://github.com/cloudquery/plugin-sdk/issues/2236)) ([6f1db88](https://github.com/cloudquery/plugin-sdk/commit/6f1db88d07830f45b0140be41a13a36363ca7d81))

## [4.87.3](https://github.com/cloudquery/plugin-sdk/compare/v4.87.2...v4.87.3) (2025-07-17)


### Bug Fixes

* **deps:** Update module github.com/aws/aws-sdk-go-v2/service/licensemanager to v1.32.0 ([#2230](https://github.com/cloudquery/plugin-sdk/issues/2230)) ([28a1479](https://github.com/cloudquery/plugin-sdk/commit/28a147900315aa9756fc66e7329e3cfbe111e4b5))
* Improve telemetry allocations ([#2185](https://github.com/cloudquery/plugin-sdk/issues/2185)) ([b07ce76](https://github.com/cloudquery/plugin-sdk/commit/b07ce76c3d3bd1e00a13b30f413d981c733d595e))
* Upgrade golangci-lint to v2 ([#2228](https://github.com/cloudquery/plugin-sdk/issues/2228)) ([7fc238c](https://github.com/cloudquery/plugin-sdk/commit/7fc238c8e7aa2c044bb1f62901c18d44c5f36d7c))

## [4.87.2](https://github.com/cloudquery/plugin-sdk/compare/v4.87.1...v4.87.2) (2025-07-09)


### Bug Fixes

* **deps:** Update module github.com/cloudquery/codegen to v0.3.31 ([#2224](https://github.com/cloudquery/plugin-sdk/issues/2224)) ([c4b5329](https://github.com/cloudquery/plugin-sdk/commit/c4b5329614ad573a1e28332eadd60bf56b9daa3f))
* Handle nil internal buffer ([#2226](https://github.com/cloudquery/plugin-sdk/issues/2226)) ([0d79bb6](https://github.com/cloudquery/plugin-sdk/commit/0d79bb6ae149498e4aaa45f4cb81561141dccfd5))

## [4.87.1](https://github.com/cloudquery/plugin-sdk/compare/v4.87.0...v4.87.1) (2025-07-09)


### Bug Fixes

* **deps:** Update module github.com/cloudquery/codegen to v0.3.30 ([#2221](https://github.com/cloudquery/plugin-sdk/issues/2221)) ([0453cbe](https://github.com/cloudquery/plugin-sdk/commit/0453cbe3391d72d5ccd490dc7761cfec93c3ea2b))
* Don't use `ValueStr`, get raw bytes instead ([#2220](https://github.com/cloudquery/plugin-sdk/issues/2220)) ([6d71d18](https://github.com/cloudquery/plugin-sdk/commit/6d71d1807b452c6d1a10decd443fb81f09501ae5))

## [4.87.0](https://github.com/cloudquery/plugin-sdk/compare/v4.86.2...v4.87.0) (2025-07-08)


### Features

* Better error when panic happens in `UnifiedDiff` ([#2217](https://github.com/cloudquery/plugin-sdk/issues/2217)) ([a282ce6](https://github.com/cloudquery/plugin-sdk/commit/a282ce6433d320e49ae0a8cd59704149cc1d3ce4))

## [4.86.2](https://github.com/cloudquery/plugin-sdk/compare/v4.86.1...v4.86.2) (2025-07-07)


### Bug Fixes

* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.13.11 ([#2211](https://github.com/cloudquery/plugin-sdk/issues/2211)) ([5d97720](https://github.com/cloudquery/plugin-sdk/commit/5d9772020c8e1e6c02ddebfb267ce2c888075ecf))
* **deps:** Update module github.com/cloudquery/codegen to v0.3.29 ([#2214](https://github.com/cloudquery/plugin-sdk/issues/2214)) ([c7534ba](https://github.com/cloudquery/plugin-sdk/commit/c7534ba2ed9b30127a8a2ffade77e90f8663ec07))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.26.15 ([#2212](https://github.com/cloudquery/plugin-sdk/issues/2212)) ([afcdfd4](https://github.com/cloudquery/plugin-sdk/commit/afcdfd4b90fbdce8ac47779976a6be4c1615c61a))
* Improve usage error wording ([#2215](https://github.com/cloudquery/plugin-sdk/issues/2215)) ([1e2c257](https://github.com/cloudquery/plugin-sdk/commit/1e2c257d343d7a2a1ac39ee06a6fc6d6659da010))

## [4.86.1](https://github.com/cloudquery/plugin-sdk/compare/v4.86.0...v4.86.1) (2025-07-01)


### Bug Fixes

* **deps:** Update dependency go to v1.24.4 ([#2209](https://github.com/cloudquery/plugin-sdk/issues/2209)) ([6b91b19](https://github.com/cloudquery/plugin-sdk/commit/6b91b194f8a9280c026670332e40225311f535f2))
* **deps:** Update golang.org/x/exp digest to b7579e2 ([#2207](https://github.com/cloudquery/plugin-sdk/issues/2207)) ([5970950](https://github.com/cloudquery/plugin-sdk/commit/5970950cd733d7031cfa1449eec91b8e51ad45ad))
* Validate and normalize inet test values ([#2205](https://github.com/cloudquery/plugin-sdk/issues/2205)) ([c9f45a2](https://github.com/cloudquery/plugin-sdk/commit/c9f45a24ead0c9fa22f0530daeef5ad40d57cd56))

## [4.86.0](https://github.com/cloudquery/plugin-sdk/compare/v4.85.0...v4.86.0) (2025-06-30)


### Features

* Make more data dynamic ([#2197](https://github.com/cloudquery/plugin-sdk/issues/2197)) ([57d0285](https://github.com/cloudquery/plugin-sdk/commit/57d0285bac8b19e7568d6d42600bda0e1dcbf5be))


### Bug Fixes

* **deps:** Update aws-sdk-go-v2 monorepo ([#2201](https://github.com/cloudquery/plugin-sdk/issues/2201)) ([aad5137](https://github.com/cloudquery/plugin-sdk/commit/aad513750f94d0839434d4c4ba965c9accf93628))
* **deps:** Update module github.com/apache/arrow-go/v18 to v18.3.1 ([#2199](https://github.com/cloudquery/plugin-sdk/issues/2199)) ([7f27c56](https://github.com/cloudquery/plugin-sdk/commit/7f27c56ff8fc5a80dd0718fc906044318a6f3e2d))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.13.10 ([#2200](https://github.com/cloudquery/plugin-sdk/issues/2200)) ([a9f5dc1](https://github.com/cloudquery/plugin-sdk/commit/a9f5dc1655b25c401d39c7d810f4d9bf24095871))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.26.14 ([#2203](https://github.com/cloudquery/plugin-sdk/issues/2203)) ([29d53f3](https://github.com/cloudquery/plugin-sdk/commit/29d53f37071513c01b60801b9640e8a2b5c1a3b7))

## [4.85.0](https://github.com/cloudquery/plugin-sdk/compare/v4.84.2...v4.85.0) (2025-06-27)


### Features

* Add handling of error messages to sdk ([#2195](https://github.com/cloudquery/plugin-sdk/issues/2195)) ([c5273da](https://github.com/cloudquery/plugin-sdk/commit/c5273da9f82e289452a9bfa54bb02b9f5c615a01))


### Bug Fixes

* **deps:** Update aws-sdk-go-v2 monorepo ([#2193](https://github.com/cloudquery/plugin-sdk/issues/2193)) ([d220f63](https://github.com/cloudquery/plugin-sdk/commit/d220f6354c0d2a901eba541126580137254b82b0))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.26.13 ([#2196](https://github.com/cloudquery/plugin-sdk/issues/2196)) ([140b6f3](https://github.com/cloudquery/plugin-sdk/commit/140b6f3d90e2b0d37f74308a50ff5f61cfeaf20c))

## [4.84.2](https://github.com/cloudquery/plugin-sdk/compare/v4.84.1...v4.84.2) (2025-06-18)


### Bug Fixes

* Add time delay in DeleteStaleAll test for destinations ([#2191](https://github.com/cloudquery/plugin-sdk/issues/2191)) ([d98a293](https://github.com/cloudquery/plugin-sdk/commit/d98a29334ac8d34fb2e85cea5600bf326c7c90ea))
* **deps:** Update module github.com/aws/aws-sdk-go-v2/config to v1.29.15 ([#2189](https://github.com/cloudquery/plugin-sdk/issues/2189)) ([9860e20](https://github.com/cloudquery/plugin-sdk/commit/9860e20cb37ca3be41b1b7723cfe2e980469b88c))
* **deps:** Update module github.com/aws/aws-sdk-go-v2/service/licensemanager to v1.31.1 ([#2186](https://github.com/cloudquery/plugin-sdk/issues/2186)) ([7647d77](https://github.com/cloudquery/plugin-sdk/commit/7647d778e746ab2b016a5912b311f053746cd85a))
* **deps:** Update module google.golang.org/grpc to v1.72.2 ([#2187](https://github.com/cloudquery/plugin-sdk/issues/2187)) ([a999c81](https://github.com/cloudquery/plugin-sdk/commit/a999c818c057c34468b0ad2c865c31dbf57fb7ad))
* **deps:** Update module google.golang.org/grpc to v1.73.0 ([#2190](https://github.com/cloudquery/plugin-sdk/issues/2190)) ([2e3c192](https://github.com/cloudquery/plugin-sdk/commit/2e3c192ff2f2455de31f51f0b6d4c16577d15c41))
* Error handling in StreamingBatchWriter ([#1921](https://github.com/cloudquery/plugin-sdk/issues/1921)) ([6d71fb1](https://github.com/cloudquery/plugin-sdk/commit/6d71fb1099792438f6527f5854a6dc37eaf298ec))

## [4.84.1](https://github.com/cloudquery/plugin-sdk/compare/v4.84.0...v4.84.1) (2025-05-30)


### Bug Fixes

* Correctly validate backend_options configuration ([#2182](https://github.com/cloudquery/plugin-sdk/issues/2182)) ([50dd38f](https://github.com/cloudquery/plugin-sdk/commit/50dd38f74af0f33eec311b0d48e14b976a3ff131))

## [4.84.0](https://github.com/cloudquery/plugin-sdk/compare/v4.83.0...v4.84.0) (2025-05-30)


### Features

* Make SDK FIPS-compliant by using internal SHA1 module ([#2179](https://github.com/cloudquery/plugin-sdk/issues/2179)) ([5a34e35](https://github.com/cloudquery/plugin-sdk/commit/5a34e3522179831f991dd8b4b59844bc1c918c1b))

## [4.83.0](https://github.com/cloudquery/plugin-sdk/compare/v4.82.2...v4.83.0) (2025-05-28)


### Features

* Switch state grpc client to NewClient rather than DialContext ([#2176](https://github.com/cloudquery/plugin-sdk/issues/2176)) ([9356d9d](https://github.com/cloudquery/plugin-sdk/commit/9356d9d14f89d3c1ea58848ae3e53d671f5b4c8f))


### Bug Fixes

* **deps:** Update dependency go to v1.24.3 ([#2041](https://github.com/cloudquery/plugin-sdk/issues/2041)) ([c438d69](https://github.com/cloudquery/plugin-sdk/commit/c438d690057cb2b8fb4944a5108b0c9bd5bfe294))

## [4.82.2](https://github.com/cloudquery/plugin-sdk/compare/v4.82.1...v4.82.2) (2025-05-26)


### Bug Fixes

* **deps:** Update module github.com/apache/arrow-go/v18 to v18.3.0 ([#2173](https://github.com/cloudquery/plugin-sdk/issues/2173)) ([f9f136d](https://github.com/cloudquery/plugin-sdk/commit/f9f136d48c5687ecd288bdfcd722d1554990827f))
* **deps:** Update module github.com/aws/aws-sdk-go-v2/service/licensemanager to v1.31.0 ([#2171](https://github.com/cloudquery/plugin-sdk/issues/2171)) ([bf74fd2](https://github.com/cloudquery/plugin-sdk/commit/bf74fd20bf667f742a657239feb7ba4c563aaaf0))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.26.12 ([#2174](https://github.com/cloudquery/plugin-sdk/issues/2174)) ([34a2d67](https://github.com/cloudquery/plugin-sdk/commit/34a2d67bd7ff96b5565e96fb053dc5d6c360b11a))

## [4.82.1](https://github.com/cloudquery/plugin-sdk/compare/v4.82.0...v4.82.1) (2025-05-23)


### Bug Fixes

* **deps:** Update module github.com/santhosh-tekuri/jsonschema/v6 to v6.0.2 ([#2168](https://github.com/cloudquery/plugin-sdk/issues/2168)) ([e8b3ecd](https://github.com/cloudquery/plugin-sdk/commit/e8b3ecdcadea5dde7c9ee805967ed665f1954c56))
* **deps:** Update opentelemetry-go monorepo ([#2166](https://github.com/cloudquery/plugin-sdk/issues/2166)) ([dfa09cd](https://github.com/cloudquery/plugin-sdk/commit/dfa09cd0dcf2126d0102e942c4f6b64eb2040a26))
* **deps:** Update opentelemetry-go monorepo ([#2169](https://github.com/cloudquery/plugin-sdk/issues/2169)) ([76ddcf6](https://github.com/cloudquery/plugin-sdk/commit/76ddcf672373e7090db9b0bfffc3f71fdd1833ef))

## [4.82.0](https://github.com/cloudquery/plugin-sdk/compare/v4.81.0...v4.82.0) (2025-05-20)


### Features

* Add `cloudquery.plugin.path` attribute to traces ([#1978](https://github.com/cloudquery/plugin-sdk/issues/1978)) ([889102e](https://github.com/cloudquery/plugin-sdk/commit/889102efe055db2f5bb9f8d4caa3dc772b4a358d))
* Improve shard chunking ([#2163](https://github.com/cloudquery/plugin-sdk/issues/2163)) ([f1c0106](https://github.com/cloudquery/plugin-sdk/commit/f1c0106f03e17cee0a8da620c7fa4638488c57b6))

## [4.81.0](https://github.com/cloudquery/plugin-sdk/compare/v4.80.3...v4.81.0) (2025-05-19)


### Features

* Remove reflect methods ([#2129](https://github.com/cloudquery/plugin-sdk/issues/2129)) ([bd277cc](https://github.com/cloudquery/plugin-sdk/commit/bd277cca0af25def7dda7d5cc8d6f8ec762ecb76))

## [4.80.3](https://github.com/cloudquery/plugin-sdk/compare/v4.80.2...v4.80.3) (2025-05-19)


### Bug Fixes

* Pass correct value for plugin version ([#2156](https://github.com/cloudquery/plugin-sdk/issues/2156)) ([37b4157](https://github.com/cloudquery/plugin-sdk/commit/37b41572165f96c8d0c67390bded2b20815a506d))

## [4.80.2](https://github.com/cloudquery/plugin-sdk/compare/v4.80.1...v4.80.2) (2025-05-19)


### Bug Fixes

* Change logic for batch writing to write when batch size is reached, not exceeded ([#2153](https://github.com/cloudquery/plugin-sdk/issues/2153)) ([58c8a1e](https://github.com/cloudquery/plugin-sdk/commit/58c8a1e35d8d77f7cb1ae1c73e70e4a21b23e0a7))
* Flush DeleteRecord messages when batch writer is flushed ([#2154](https://github.com/cloudquery/plugin-sdk/issues/2154)) ([791c865](https://github.com/cloudquery/plugin-sdk/commit/791c8658a0b0224f080dce8ea0cf734dbc9ce911))

## [4.80.1](https://github.com/cloudquery/plugin-sdk/compare/v4.80.0...v4.80.1) (2025-05-12)


### Bug Fixes

* **deps:** Update module github.com/cloudquery/codegen to v0.3.27 ([#2148](https://github.com/cloudquery/plugin-sdk/issues/2148)) ([1fd7b1e](https://github.com/cloudquery/plugin-sdk/commit/1fd7b1e95a262017152616290da791e08fc497b8))
* **deps:** Update module github.com/cloudquery/codegen to v0.3.28 ([#2150](https://github.com/cloudquery/plugin-sdk/issues/2150)) ([4d1409c](https://github.com/cloudquery/plugin-sdk/commit/4d1409c402fd5e8b793ee0586e9d2ec8c3812cd5))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.26.11 ([#2149](https://github.com/cloudquery/plugin-sdk/issues/2149)) ([a904b46](https://github.com/cloudquery/plugin-sdk/commit/a904b461a5ca7eebb0c6f4c807db9aa23011a2b5))
* **deps:** Update module github.com/samber/lo to v1.49.1 ([#2139](https://github.com/cloudquery/plugin-sdk/issues/2139)) ([f11b5e6](https://github.com/cloudquery/plugin-sdk/commit/f11b5e67d95af1f54410236bbefb071c01c5df82))

## [4.80.0](https://github.com/cloudquery/plugin-sdk/compare/v4.79.1...v4.80.0) (2025-05-09)


### Features

* Add SensitiveColumns to tables schema ([#2134](https://github.com/cloudquery/plugin-sdk/issues/2134)) ([e95674f](https://github.com/cloudquery/plugin-sdk/commit/e95674f255c7225a9b5d593daf54c0a373c2ef50))


### Bug Fixes

* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.13.9 ([#2143](https://github.com/cloudquery/plugin-sdk/issues/2143)) ([77d4b9b](https://github.com/cloudquery/plugin-sdk/commit/77d4b9b317dfa06ccbf4c0696adcdaa6de724173))

## [4.79.1](https://github.com/cloudquery/plugin-sdk/compare/v4.79.0...v4.79.1) (2025-05-03)


### Bug Fixes

* **deps:** Update golang.org/x/exp digest to 7e4ce0a ([#2135](https://github.com/cloudquery/plugin-sdk/issues/2135)) ([efb8813](https://github.com/cloudquery/plugin-sdk/commit/efb8813182f813a25ec9bf0ba465f5a1419c937d))
* **deps:** Update module github.com/grpc-ecosystem/go-grpc-middleware/v2 to v2.3.1 ([#2136](https://github.com/cloudquery/plugin-sdk/issues/2136)) ([5534187](https://github.com/cloudquery/plugin-sdk/commit/553418708ebc1fcbc1eb13ea29907346258aa36b))
* **deps:** Update module github.com/rs/zerolog to v1.34.0 ([#2138](https://github.com/cloudquery/plugin-sdk/issues/2138)) ([e1bcb05](https://github.com/cloudquery/plugin-sdk/commit/e1bcb0532b7b1b42e38571254c655769296240e4))
* **deps:** Update module google.golang.org/grpc to v1.72.0 ([#2141](https://github.com/cloudquery/plugin-sdk/issues/2141)) ([a0f27a3](https://github.com/cloudquery/plugin-sdk/commit/a0f27a3ca01a870912ad7afb9ec92fb8e52f78b7))
* Pass installation ID from env for usage report ([#2140](https://github.com/cloudquery/plugin-sdk/issues/2140)) ([4d36bfb](https://github.com/cloudquery/plugin-sdk/commit/4d36bfba0a572995c11e10fc2e448e36beb67208))

## [4.79.0](https://github.com/cloudquery/plugin-sdk/compare/v4.78.0...v4.79.0) (2025-04-28)


### Features

* Add transformer to update table description with its table options ([#2128](https://github.com/cloudquery/plugin-sdk/issues/2128)) ([2387b57](https://github.com/cloudquery/plugin-sdk/commit/2387b5765133a03e202f332290bfd94c2ac50eab))
* Show plugin version in plugin server logs ([#2124](https://github.com/cloudquery/plugin-sdk/issues/2124)) ([be08606](https://github.com/cloudquery/plugin-sdk/commit/be08606413a4392d04d3e388ccc1edbe64439c14))


### Bug Fixes

* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.26.10 ([#2132](https://github.com/cloudquery/plugin-sdk/issues/2132)) ([775d537](https://github.com/cloudquery/plugin-sdk/commit/775d537cc9ee69ee548d69be67e7172baa03fcab))
* Prevent deadlock in transformer ([#2130](https://github.com/cloudquery/plugin-sdk/issues/2130)) ([a65b101](https://github.com/cloudquery/plugin-sdk/commit/a65b101d05ee43bc1f1f0033736ede756ee55604))

## [4.78.0](https://github.com/cloudquery/plugin-sdk/compare/v4.77.0...v4.78.0) (2025-04-22)


### Features

* Add logger to context ([#2125](https://github.com/cloudquery/plugin-sdk/issues/2125)) ([718e8ed](https://github.com/cloudquery/plugin-sdk/commit/718e8ed781fb27130636a87c76bfeb6c00348383))


### Bug Fixes

* **deps:** Update aws-sdk-go-v2 monorepo ([#2119](https://github.com/cloudquery/plugin-sdk/issues/2119)) ([5554039](https://github.com/cloudquery/plugin-sdk/commit/5554039d4358a66f21e765b8dc7c3203b7437f04))
* **deps:** Update aws-sdk-go-v2 monorepo ([#2121](https://github.com/cloudquery/plugin-sdk/issues/2121)) ([7b54577](https://github.com/cloudquery/plugin-sdk/commit/7b54577964b523aba6ca93497d65c8bad6132149))
* **deps:** Update aws-sdk-go-v2 monorepo ([#2123](https://github.com/cloudquery/plugin-sdk/issues/2123)) ([8f370f8](https://github.com/cloudquery/plugin-sdk/commit/8f370f80da7ba9f8c48896aa31414bf16c57fbf1))
* **deps:** Update Google Golang modules ([#2118](https://github.com/cloudquery/plugin-sdk/issues/2118)) ([93d9203](https://github.com/cloudquery/plugin-sdk/commit/93d9203936fb499ab516fe9e847e078e758afb36))
* **deps:** Update module golang.org/x/net to v0.38.0 [SECURITY] ([#2122](https://github.com/cloudquery/plugin-sdk/issues/2122)) ([0b0e187](https://github.com/cloudquery/plugin-sdk/commit/0b0e18763cccad1d01cf61c7a8b4c6c5e5ec343c))

## [4.77.0](https://github.com/cloudquery/plugin-sdk/compare/v4.76.0...v4.77.0) (2025-04-03)


### Features

* Allow skipping PK components mismatch validation ([#2115](https://github.com/cloudquery/plugin-sdk/issues/2115)) ([caf7f92](https://github.com/cloudquery/plugin-sdk/commit/caf7f92168d69d2cd95350162fa9ea8aef33cd5d))

## [4.76.0](https://github.com/cloudquery/plugin-sdk/compare/v4.75.0...v4.76.0) (2025-04-02)


### Features

* Pass installation ID from env to usage report ([#2106](https://github.com/cloudquery/plugin-sdk/issues/2106)) ([0bea6e7](https://github.com/cloudquery/plugin-sdk/commit/0bea6e792c2024231e6cf20804e0e493cf63f012))


### Bug Fixes

* **deps:** Update golang.org/x/exp digest to 054e65f ([#2110](https://github.com/cloudquery/plugin-sdk/issues/2110)) ([f9875f8](https://github.com/cloudquery/plugin-sdk/commit/f9875f8706268a9c0a3d3d95eb73da9d6d258f9f))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.26.9 ([#2112](https://github.com/cloudquery/plugin-sdk/issues/2112)) ([abd2117](https://github.com/cloudquery/plugin-sdk/commit/abd2117c47837393fb06eecf270bf4d25ce511a7))
* Error if both PKs and PK components are set ([#2113](https://github.com/cloudquery/plugin-sdk/issues/2113)) ([4f0b312](https://github.com/cloudquery/plugin-sdk/commit/4f0b3129391339df8b53358433ad17555149399e))

## [4.75.0](https://github.com/cloudquery/plugin-sdk/compare/v4.74.2...v4.75.0) (2025-03-31)


### Features

* Add internal columns helpers ([#2105](https://github.com/cloudquery/plugin-sdk/issues/2105)) ([1dea99c](https://github.com/cloudquery/plugin-sdk/commit/1dea99c53b1e959470a5532c0d76f6fec2238876))

## [4.74.2](https://github.com/cloudquery/plugin-sdk/compare/v4.74.1...v4.74.2) (2025-03-24)


### Bug Fixes

* **deps:** Update aws-sdk-go-v2 monorepo ([#2096](https://github.com/cloudquery/plugin-sdk/issues/2096)) ([f49534a](https://github.com/cloudquery/plugin-sdk/commit/f49534ada5fe0f189d534deaa0fd5d2990122ddb))
* **deps:** Update aws-sdk-go-v2 monorepo ([#2100](https://github.com/cloudquery/plugin-sdk/issues/2100)) ([07a3ed8](https://github.com/cloudquery/plugin-sdk/commit/07a3ed85d1196818c743234d56a4e76490ac3213))
* **deps:** Update module github.com/apache/arrow-go/v18 to v18.2.0 ([#2103](https://github.com/cloudquery/plugin-sdk/issues/2103)) ([f6b7143](https://github.com/cloudquery/plugin-sdk/commit/f6b7143c1168bfdc17410a10e98b6459e76b480d))
* **deps:** Update module github.com/aws/aws-sdk-go-v2/service/marketplacemetering to v1.26.2 ([#2102](https://github.com/cloudquery/plugin-sdk/issues/2102)) ([ddae6e0](https://github.com/cloudquery/plugin-sdk/commit/ddae6e070cc34b8e9ffdb2170e60bfcf6e38591f))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.13.7 ([#2099](https://github.com/cloudquery/plugin-sdk/issues/2099)) ([316ff40](https://github.com/cloudquery/plugin-sdk/commit/316ff406cda86f3376a6d464e77553d05c78db60))
* **deps:** Update module golang.org/x/net to v0.36.0 [SECURITY] ([#2098](https://github.com/cloudquery/plugin-sdk/issues/2098)) ([b41044d](https://github.com/cloudquery/plugin-sdk/commit/b41044d3d9cd1de9c20012a38808e8498bd51f2e))
* **deps:** Update module google.golang.org/grpc to v1.71.0 ([#2101](https://github.com/cloudquery/plugin-sdk/issues/2101)) ([7086507](https://github.com/cloudquery/plugin-sdk/commit/7086507c5536ec1930a9936799146a947ac10e95))

## [4.74.1](https://github.com/cloudquery/plugin-sdk/compare/v4.74.0...v4.74.1) (2025-03-07)


### Bug Fixes

* Revert faker pointer change ([#2093](https://github.com/cloudquery/plugin-sdk/issues/2093)) ([4157755](https://github.com/cloudquery/plugin-sdk/commit/4157755d75399f3e831a7d6fed719cf456c15189))

## [4.74.0](https://github.com/cloudquery/plugin-sdk/compare/v4.73.4...v4.74.0) (2025-03-07)


### Features

* Add description to time JSONschema ([#2083](https://github.com/cloudquery/plugin-sdk/issues/2083)) ([fc27b14](https://github.com/cloudquery/plugin-sdk/commit/fc27b14c801d44ffba458fc676a6e3f0d577dfec))
* Add way to get Configtype time value for hashing ([#2077](https://github.com/cloudquery/plugin-sdk/issues/2077)) ([ed27292](https://github.com/cloudquery/plugin-sdk/commit/ed27292685e032027802410265e92e96fbfc375e))
* Add WithNullableFieldTransformer to transformers.TransformWithStruct options ([#2084](https://github.com/cloudquery/plugin-sdk/issues/2084)) ([2175946](https://github.com/cloudquery/plugin-sdk/commit/217594694dfa253f1534aead73c3481978210a9d))


### Bug Fixes

* **deps:** Update aws-sdk-go-v2 monorepo ([#2074](https://github.com/cloudquery/plugin-sdk/issues/2074)) ([091e2c7](https://github.com/cloudquery/plugin-sdk/commit/091e2c7c028717a8dc91cf6440bd6307468f2fdc))
* **deps:** Update aws-sdk-go-v2 monorepo ([#2080](https://github.com/cloudquery/plugin-sdk/issues/2080)) ([14ad9db](https://github.com/cloudquery/plugin-sdk/commit/14ad9db4509fdf38516b7140b519af040ea8f218))
* **deps:** Update aws-sdk-go-v2 monorepo ([#2085](https://github.com/cloudquery/plugin-sdk/issues/2085)) ([3da5572](https://github.com/cloudquery/plugin-sdk/commit/3da55726a9264d8d74883d40886de122f0377e1c))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.13.6 ([#2089](https://github.com/cloudquery/plugin-sdk/issues/2089)) ([4faff6a](https://github.com/cloudquery/plugin-sdk/commit/4faff6a844959875d43a5ceb8e1d18ca9d35e3ce))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.26.7 ([#2076](https://github.com/cloudquery/plugin-sdk/issues/2076)) ([dbdae2f](https://github.com/cloudquery/plugin-sdk/commit/dbdae2f4829cfd440e395dfaa9047e12f87d1fb2))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.26.8 ([#2090](https://github.com/cloudquery/plugin-sdk/issues/2090)) ([c4c9cd6](https://github.com/cloudquery/plugin-sdk/commit/c4c9cd60e35cd50e6567afdfc9901381bb293dd7))
* **deps:** Update module github.com/goccy/go-json to v0.10.5 ([#2086](https://github.com/cloudquery/plugin-sdk/issues/2086)) ([b238237](https://github.com/cloudquery/plugin-sdk/commit/b23823758f6a0d74a7ab9a96c3dfd5bdc09dce40))
* **deps:** Update module github.com/google/go-cmp to v0.7.0 ([#2088](https://github.com/cloudquery/plugin-sdk/issues/2088)) ([1895ddf](https://github.com/cloudquery/plugin-sdk/commit/1895ddf6d5438df771a16f2ad33e14c9ee77cbd4))
* **deps:** Update module google.golang.org/protobuf to v1.36.5 ([#2081](https://github.com/cloudquery/plugin-sdk/issues/2081)) ([833b19c](https://github.com/cloudquery/plugin-sdk/commit/833b19c406de76be59c19bb9c53695c675468325))
* Ignore context cancelled errors. ([#2091](https://github.com/cloudquery/plugin-sdk/issues/2091)) ([bc50fd3](https://github.com/cloudquery/plugin-sdk/commit/bc50fd3d2be414edba8f8ad5bb7739a012840bf1))

## [4.73.4](https://github.com/cloudquery/plugin-sdk/compare/v4.73.3...v4.73.4) (2025-02-03)


### Bug Fixes

* **deps:** Update aws-sdk-go-v2 monorepo ([#2061](https://github.com/cloudquery/plugin-sdk/issues/2061)) ([7346223](https://github.com/cloudquery/plugin-sdk/commit/7346223d04dd709ab5e89b35f3914971ba8677d9))
* **deps:** Update aws-sdk-go-v2 monorepo ([#2067](https://github.com/cloudquery/plugin-sdk/issues/2067)) ([21125d0](https://github.com/cloudquery/plugin-sdk/commit/21125d0c4f39553312fe0c82578d41bde6c707ca))
* **deps:** Update Google Golang modules ([#2060](https://github.com/cloudquery/plugin-sdk/issues/2060)) ([d3a180d](https://github.com/cloudquery/plugin-sdk/commit/d3a180d8968b6fc30f645d7f34cf6a63d5632497))
* **deps:** Update Google Golang modules ([#2066](https://github.com/cloudquery/plugin-sdk/issues/2066)) ([6c32c4a](https://github.com/cloudquery/plugin-sdk/commit/6c32c4a75a5fda58efc40b0da149bf339ccd54ff))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.26.6 ([#2072](https://github.com/cloudquery/plugin-sdk/issues/2072)) ([00ce2d7](https://github.com/cloudquery/plugin-sdk/commit/00ce2d772ee4e27629252c0b222843cc71036d2c))
* **deps:** Update module github.com/invopop/jsonschema to v0.13.0 ([#2068](https://github.com/cloudquery/plugin-sdk/issues/2068)) ([c8122a2](https://github.com/cloudquery/plugin-sdk/commit/c8122a2685b057f1f522991d55a78859e7ed67e2))
* **deps:** Update module golang.org/x/oauth2 to v0.25.0 ([#2069](https://github.com/cloudquery/plugin-sdk/issues/2069)) ([9448009](https://github.com/cloudquery/plugin-sdk/commit/944800907fa7b065a0bb2ee81cae12598fae07e6))
* **deps:** Update opentelemetry-go monorepo ([#2070](https://github.com/cloudquery/plugin-sdk/issues/2070)) ([66793b9](https://github.com/cloudquery/plugin-sdk/commit/66793b9e2ec5694c5aa186ff951aee8aeee1530b))

## [4.73.3](https://github.com/cloudquery/plugin-sdk/compare/v4.73.2...v4.73.3) (2025-01-20)


### Bug Fixes

* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.26.5 ([#2057](https://github.com/cloudquery/plugin-sdk/issues/2057)) ([c91a230](https://github.com/cloudquery/plugin-sdk/commit/c91a23052cba9fcc58cd55b420dabdee027afe3b))

## [4.73.2](https://github.com/cloudquery/plugin-sdk/compare/v4.73.1...v4.73.2) (2025-01-20)


### Bug Fixes

* **deps:** Update aws-sdk-go-v2 monorepo ([#2052](https://github.com/cloudquery/plugin-sdk/issues/2052)) ([ea0d787](https://github.com/cloudquery/plugin-sdk/commit/ea0d787ef1229b2b9d89be81842e85da78504b1e))
* **deps:** Update module github.com/apache/arrow-go/v18 to v18.1.0 ([#2055](https://github.com/cloudquery/plugin-sdk/issues/2055)) ([a0f0dc6](https://github.com/cloudquery/plugin-sdk/commit/a0f0dc6a80826a7a1e79e4f0f596b5f6313d7cd1))
* **deps:** Update module google.golang.org/protobuf to v1.36.2 ([#2053](https://github.com/cloudquery/plugin-sdk/issues/2053)) ([78a26e4](https://github.com/cloudquery/plugin-sdk/commit/78a26e46b60891d9e9ad1cd6ad6ec3272db7f5dc))

## [4.73.1](https://github.com/cloudquery/plugin-sdk/compare/v4.73.0...v4.73.1) (2025-01-15)


### Bug Fixes

* AWS Marketplace Integration ([#2049](https://github.com/cloudquery/plugin-sdk/issues/2049)) ([97a3706](https://github.com/cloudquery/plugin-sdk/commit/97a3706efa87e74924b6769775926d66d602484e))

## [4.73.0](https://github.com/cloudquery/plugin-sdk/compare/v4.72.6...v4.73.0) (2025-01-08)


### Features

* Enable storing _cq_client_id. ([#2046](https://github.com/cloudquery/plugin-sdk/issues/2046)) ([3b28991](https://github.com/cloudquery/plugin-sdk/commit/3b2899111934e3d4f925ec75f2ffec61fbc038eb))

## [4.72.6](https://github.com/cloudquery/plugin-sdk/compare/v4.72.5...v4.72.6) (2025-01-07)


### Bug Fixes

* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.26.4 ([#2044](https://github.com/cloudquery/plugin-sdk/issues/2044)) ([c7bd2d2](https://github.com/cloudquery/plugin-sdk/commit/c7bd2d26b6ad3d298f77f46812b365ff40e6cf25))
* **deps:** Update module github.com/goccy/go-json to v0.10.4 ([#2040](https://github.com/cloudquery/plugin-sdk/issues/2040)) ([f6e0201](https://github.com/cloudquery/plugin-sdk/commit/f6e0201f18745a8d96bb2fbae92b9cfcffdb646a))
* **deps:** Update module google.golang.org/protobuf to v1.36.1 ([#2043](https://github.com/cloudquery/plugin-sdk/issues/2043)) ([13437c2](https://github.com/cloudquery/plugin-sdk/commit/13437c25d40dd7b2d0b32efa560d65b6f7a5af1d))
* **deps:** Update opentelemetry-go monorepo ([#2042](https://github.com/cloudquery/plugin-sdk/issues/2042)) ([e6123c3](https://github.com/cloudquery/plugin-sdk/commit/e6123c3c14920bd60c723e37bf600e15821c87f1))
* Log warning instead of erroring out of PK component validation failure ([#2039](https://github.com/cloudquery/plugin-sdk/issues/2039)) ([c98b5c5](https://github.com/cloudquery/plugin-sdk/commit/c98b5c5e23af1246f6ed8a47b57a79300d403f10))
* Validate missing PK components ([#2037](https://github.com/cloudquery/plugin-sdk/issues/2037)) ([d2cff6b](https://github.com/cloudquery/plugin-sdk/commit/d2cff6b5fb90ebdd7063384166115a50ce2132f4))

## [4.72.5](https://github.com/cloudquery/plugin-sdk/compare/v4.72.4...v4.72.5) (2024-12-23)


### Bug Fixes

* **deps:** Update aws-sdk-go-v2 monorepo ([#2032](https://github.com/cloudquery/plugin-sdk/issues/2032)) ([7f6fb0a](https://github.com/cloudquery/plugin-sdk/commit/7f6fb0a5fe726a94cad0c970c94234611efc1453))
* **deps:** Update golang.org/x/exp digest to b2144cd ([#2029](https://github.com/cloudquery/plugin-sdk/issues/2029)) ([f955d43](https://github.com/cloudquery/plugin-sdk/commit/f955d4384eb6f89af3582315824061756a5b2251))
* **deps:** Update Google Golang modules ([#2031](https://github.com/cloudquery/plugin-sdk/issues/2031)) ([1da77a5](https://github.com/cloudquery/plugin-sdk/commit/1da77a55cc74d433bed8d3e815b80e529d66a5b5))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.26.3 ([#2035](https://github.com/cloudquery/plugin-sdk/issues/2035)) ([310aea1](https://github.com/cloudquery/plugin-sdk/commit/310aea168cf7ce58cf1fa459381a32529f813458))
* **deps:** Update module github.com/grpc-ecosystem/go-grpc-middleware/v2 to v2.2.0 ([#2033](https://github.com/cloudquery/plugin-sdk/issues/2033)) ([58d29b3](https://github.com/cloudquery/plugin-sdk/commit/58d29b3bc5a47cabd0add3ee72b44d92ee5223e6))
* **deps:** Update module google.golang.org/grpc to v1.69.0 ([#2027](https://github.com/cloudquery/plugin-sdk/issues/2027)) ([8542575](https://github.com/cloudquery/plugin-sdk/commit/85425759ef318b73112d8e4951ce3efb4dd39831))

## [4.72.4](https://github.com/cloudquery/plugin-sdk/compare/v4.72.3...v4.72.4) (2024-12-20)


### Bug Fixes

* Revert "fix(deps): Update module google.golang.org/grpc to v1.69.0" ([#2023](https://github.com/cloudquery/plugin-sdk/issues/2023)) ([78a6371](https://github.com/cloudquery/plugin-sdk/commit/78a6371a757ffdb53a3aa0158d57a2d4ad313f1b))

## [4.72.3](https://github.com/cloudquery/plugin-sdk/compare/v4.72.2...v4.72.3) (2024-12-20)


### Bug Fixes

* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.26.1 ([#2021](https://github.com/cloudquery/plugin-sdk/issues/2021)) ([d4e9d15](https://github.com/cloudquery/plugin-sdk/commit/d4e9d1553c71a063df0c5b998fd493d59f17af5f))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.26.2 ([#2022](https://github.com/cloudquery/plugin-sdk/issues/2022)) ([b0a7640](https://github.com/cloudquery/plugin-sdk/commit/b0a76401de2c9ade87a01f465efba815798a6c2c))
* Revert "fix(deps): Update module google.golang.org/grpc to v1.69.0" ([#2018](https://github.com/cloudquery/plugin-sdk/issues/2018)) ([6d72c67](https://github.com/cloudquery/plugin-sdk/commit/6d72c67b06ee3770bb396fad58c5f0eacc911ec2))

## [4.72.2](https://github.com/cloudquery/plugin-sdk/compare/v4.72.1...v4.72.2) (2024-12-19)


### Bug Fixes

* Use field name in json type schema if json tag is missing ([#2011](https://github.com/cloudquery/plugin-sdk/issues/2011)) ([7ca8009](https://github.com/cloudquery/plugin-sdk/commit/7ca8009bec4214928fdeb2473b7c04294ae7952e))

## [4.72.1](https://github.com/cloudquery/plugin-sdk/compare/v4.72.0...v4.72.1) (2024-12-19)


### Bug Fixes

* **deps:** Update aws-sdk-go-v2 monorepo ([#2007](https://github.com/cloudquery/plugin-sdk/issues/2007)) ([7f3818d](https://github.com/cloudquery/plugin-sdk/commit/7f3818d51a2d60bc7dc2a3846ef038a783d984bc))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.13.5 ([#2015](https://github.com/cloudquery/plugin-sdk/issues/2015)) ([9b6e9f2](https://github.com/cloudquery/plugin-sdk/commit/9b6e9f29ac3d165bf5470e933b8638a961b4bd64))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.26.1 ([#2010](https://github.com/cloudquery/plugin-sdk/issues/2010)) ([b12dc10](https://github.com/cloudquery/plugin-sdk/commit/b12dc1033a5130629c4ff3eb76c233704df81747))
* **deps:** Update module golang.org/x/net to v0.33.0 [SECURITY] ([#2014](https://github.com/cloudquery/plugin-sdk/issues/2014)) ([7360bd2](https://github.com/cloudquery/plugin-sdk/commit/7360bd26d49e76f48182efdad8d75a07a95e0263))
* **deps:** Update module google.golang.org/grpc to v1.69.0 ([#2008](https://github.com/cloudquery/plugin-sdk/issues/2008)) ([aae018f](https://github.com/cloudquery/plugin-sdk/commit/aae018f9838c80c3ff5f10ec6b47c41d809b4694))
* OpenTelemetry schema URL panic ([#2012](https://github.com/cloudquery/plugin-sdk/issues/2012)) ([b616279](https://github.com/cloudquery/plugin-sdk/commit/b6162796a417cea0b8bb0efee8074917dce63415))

## [4.72.0](https://github.com/cloudquery/plugin-sdk/compare/v4.71.1...v4.72.0) (2024-12-13)


### Features

* Update to Arrow v18 ([#1997](https://github.com/cloudquery/plugin-sdk/issues/1997)) ([5b84d3d](https://github.com/cloudquery/plugin-sdk/commit/5b84d3deb1d9bad854fecd183871a58213ea4773))

## [4.71.1](https://github.com/cloudquery/plugin-sdk/compare/v4.71.0...v4.71.1) (2024-12-12)


### Bug Fixes

* Fix Transform hang in CLI sync ([#2001](https://github.com/cloudquery/plugin-sdk/issues/2001)) ([0474b9f](https://github.com/cloudquery/plugin-sdk/commit/0474b9fa43f4fd8a99e7fc55f6e651d1c7963213))

## [4.71.0](https://github.com/cloudquery/plugin-sdk/compare/v4.70.2...v4.71.0) (2024-12-09)


### Features

* Implement batch sender. ([#1995](https://github.com/cloudquery/plugin-sdk/issues/1995)) ([371b20f](https://github.com/cloudquery/plugin-sdk/commit/371b20fd192e69681e07c79302e7a06fc89b4a71))

## [4.70.2](https://github.com/cloudquery/plugin-sdk/compare/v4.70.1...v4.70.2) (2024-12-05)


### Bug Fixes

* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.13.4 ([#1992](https://github.com/cloudquery/plugin-sdk/issues/1992)) ([cd4dc4b](https://github.com/cloudquery/plugin-sdk/commit/cd4dc4bcfb9eb42227bce7cf77899a5a31635a20))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.25.5 ([#1991](https://github.com/cloudquery/plugin-sdk/issues/1991)) ([037a6d9](https://github.com/cloudquery/plugin-sdk/commit/037a6d97ccf8dbcf6f5a9a28fa6f945f8892af25))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.25.6 ([#1994](https://github.com/cloudquery/plugin-sdk/issues/1994)) ([32855ea](https://github.com/cloudquery/plugin-sdk/commit/32855ea19975675bba2ed4cfaf1a00013f49a7b5))
* Handle integer overflow ([#1996](https://github.com/cloudquery/plugin-sdk/issues/1996)) ([6af9c22](https://github.com/cloudquery/plugin-sdk/commit/6af9c22a82b3872a3fbffdc7f70c61d63450be6e))

## [4.70.1](https://github.com/cloudquery/plugin-sdk/compare/v4.70.0...v4.70.1) (2024-12-02)


### Bug Fixes

* **deps:** Update golang.org/x/exp digest to 2d47ceb ([#1979](https://github.com/cloudquery/plugin-sdk/issues/1979)) ([3785e86](https://github.com/cloudquery/plugin-sdk/commit/3785e8624e3a6248b136b8dabfbe40c550a1e4ee))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.13.3 ([#1985](https://github.com/cloudquery/plugin-sdk/issues/1985)) ([03944f2](https://github.com/cloudquery/plugin-sdk/commit/03944f20f8d1b9b0018ab466200a2b938c32a203))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.25.4 ([#1987](https://github.com/cloudquery/plugin-sdk/issues/1987)) ([1de3891](https://github.com/cloudquery/plugin-sdk/commit/1de3891a63df6e0cce1ad606a4292f8277831d44))
* **deps:** Update module github.com/stretchr/testify to v1.10.0 ([#1982](https://github.com/cloudquery/plugin-sdk/issues/1982)) ([c3ac76c](https://github.com/cloudquery/plugin-sdk/commit/c3ac76cc0655e1222ec1478c2ca6b60f561a425c))
* **deps:** Update module golang.org/x/oauth2 to v0.24.0 ([#1983](https://github.com/cloudquery/plugin-sdk/issues/1983)) ([90b9728](https://github.com/cloudquery/plugin-sdk/commit/90b972850bdc4994425992ea1ed595d9ab1d8ad0))
* **deps:** Update opentelemetry-go monorepo ([#1984](https://github.com/cloudquery/plugin-sdk/issues/1984)) ([3894891](https://github.com/cloudquery/plugin-sdk/commit/38948915ed0fc7dbd59158b80755cf3b6f4497e6))

## [4.70.0](https://github.com/cloudquery/plugin-sdk/compare/v4.69.0...v4.70.0) (2024-11-26)


### Features

* Expose resource stats in traces as separate event ([#1973](https://github.com/cloudquery/plugin-sdk/issues/1973)) ([e74bb27](https://github.com/cloudquery/plugin-sdk/commit/e74bb27ae66713f5734d8cd5fbe1905226c3e696))


### Bug Fixes

* **deps:** Update aws-sdk-go-v2 monorepo ([#1976](https://github.com/cloudquery/plugin-sdk/issues/1976)) ([566409d](https://github.com/cloudquery/plugin-sdk/commit/566409d09d422fb26453879e3e44a00f0608ab01))

## [4.69.0](https://github.com/cloudquery/plugin-sdk/compare/v4.68.3...v4.69.0) (2024-11-22)


### Features

* Allow to include symbols in binaries during `package` ([#1974](https://github.com/cloudquery/plugin-sdk/issues/1974)) ([aa3b3e4](https://github.com/cloudquery/plugin-sdk/commit/aa3b3e45e206d57cf812f131020cfed7ad44e9f7))


### Bug Fixes

* **deps:** Update module github.com/aws/aws-sdk-go-v2/config to v1.28.4 ([#1968](https://github.com/cloudquery/plugin-sdk/issues/1968)) ([4e35df7](https://github.com/cloudquery/plugin-sdk/commit/4e35df709bb2190794b6aa9713c00188200b07f1))
* **deps:** Update module google.golang.org/protobuf to v1.35.2 ([#1971](https://github.com/cloudquery/plugin-sdk/issues/1971)) ([7076899](https://github.com/cloudquery/plugin-sdk/commit/7076899d9933b53dbb03533bbf28b83dc1b37a27))

## [4.68.3](https://github.com/cloudquery/plugin-sdk/compare/v4.68.2...v4.68.3) (2024-11-14)


### Bug Fixes

* Correctly handle success in DryRun invocation ([#1966](https://github.com/cloudquery/plugin-sdk/issues/1966)) ([9b3f292](https://github.com/cloudquery/plugin-sdk/commit/9b3f2924703532b907255807abaac380a42c6aeb))
* **deps:** Update aws-sdk-go-v2 monorepo ([#1963](https://github.com/cloudquery/plugin-sdk/issues/1963)) ([41f717e](https://github.com/cloudquery/plugin-sdk/commit/41f717eb97e1cb56dbd6636620e9b0f41567d942))
* **deps:** Update module google.golang.org/grpc to v1.68.0 ([#1964](https://github.com/cloudquery/plugin-sdk/issues/1964)) ([763d55f](https://github.com/cloudquery/plugin-sdk/commit/763d55fa8d9ba4a7543f9debad52e312c5fab4c8))

## [4.68.2](https://github.com/cloudquery/plugin-sdk/compare/v4.68.1...v4.68.2) (2024-11-04)


### Bug Fixes

* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.25.2 ([#1960](https://github.com/cloudquery/plugin-sdk/issues/1960)) ([68fbf32](https://github.com/cloudquery/plugin-sdk/commit/68fbf32417d6799a8589b7fb78fa909edbdbf169))

## [4.68.1](https://github.com/cloudquery/plugin-sdk/compare/v4.68.0...v4.68.1) (2024-11-04)


### Bug Fixes

* **deps:** Update aws-sdk-go-v2 monorepo ([#1957](https://github.com/cloudquery/plugin-sdk/issues/1957)) ([360cc57](https://github.com/cloudquery/plugin-sdk/commit/360cc579c7d6dffb74872183b9d5308e73b41f15))
* **deps:** Update golang.org/x/exp digest to f66d83c ([#1954](https://github.com/cloudquery/plugin-sdk/issues/1954)) ([18cb1b2](https://github.com/cloudquery/plugin-sdk/commit/18cb1b2428480aa3143d2d4a0bf40aaf3c4b803b))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.25.1 ([#1958](https://github.com/cloudquery/plugin-sdk/issues/1958)) ([f537b56](https://github.com/cloudquery/plugin-sdk/commit/f537b56c0527d40ec89e571c1510adf6669c113c))
* **deps:** Update opentelemetry-go monorepo ([#1956](https://github.com/cloudquery/plugin-sdk/issues/1956)) ([ea171a4](https://github.com/cloudquery/plugin-sdk/commit/ea171a44fef368ea3e60a6ae5b5018d349c0b989))

## [4.68.0](https://github.com/cloudquery/plugin-sdk/compare/v4.67.1...v4.68.0) (2024-10-31)


### Features

* Add Time configtype ([#1905](https://github.com/cloudquery/plugin-sdk/issues/1905)) ([f57c3eb](https://github.com/cloudquery/plugin-sdk/commit/f57c3ebecf99f0d7fe546c058d4086e2454075ba))
* Support for quota query interval header ([#1948](https://github.com/cloudquery/plugin-sdk/issues/1948)) ([bfce6fe](https://github.com/cloudquery/plugin-sdk/commit/bfce6fee435085af67163f4fed6168d4459aa87b))
* Test `MeterUsage` API call on initial setup of client ([#1906](https://github.com/cloudquery/plugin-sdk/issues/1906)) ([78df77d](https://github.com/cloudquery/plugin-sdk/commit/78df77d3c20a5f0a4ccc037fc82c6f626a6d5e1c))


### Bug Fixes

* Clean up usage retry logic ([#1950](https://github.com/cloudquery/plugin-sdk/issues/1950)) ([ca982f9](https://github.com/cloudquery/plugin-sdk/commit/ca982f92d65dbf55bd849fbe7688200f1d03c66a))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.25.0 ([#1946](https://github.com/cloudquery/plugin-sdk/issues/1946)) ([b8e3e10](https://github.com/cloudquery/plugin-sdk/commit/b8e3e104071fa3454d74762fc4c45d0cc98f31ab))

## [4.67.1](https://github.com/cloudquery/plugin-sdk/compare/v4.67.0...v4.67.1) (2024-10-22)


### Bug Fixes

* **deps:** Update module github.com/aws/aws-sdk-go-v2/config to v1.28.0 ([#1940](https://github.com/cloudquery/plugin-sdk/issues/1940)) ([35cf587](https://github.com/cloudquery/plugin-sdk/commit/35cf587f2c96d8bbadbd0b4cdb0484039a77f089))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.24.1 ([#1943](https://github.com/cloudquery/plugin-sdk/issues/1943)) ([14f44ad](https://github.com/cloudquery/plugin-sdk/commit/14f44adf41ba797e156378a208ec1528070d4fcd))
* Ensure module field exists in all log messages ([#1941](https://github.com/cloudquery/plugin-sdk/issues/1941)) ([b1ca41c](https://github.com/cloudquery/plugin-sdk/commit/b1ca41c632069900225b556339e74fb6d2136c6c))

## [4.67.0](https://github.com/cloudquery/plugin-sdk/compare/v4.66.1...v4.67.0) (2024-10-18)


### Features

* Make state client versioning the default, remove option ([#1938](https://github.com/cloudquery/plugin-sdk/issues/1938)) ([f105651](https://github.com/cloudquery/plugin-sdk/commit/f1056512e7ec6675808f1ed5ea398ba3c3da82ad))


### Bug Fixes

* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.23.0 ([#1934](https://github.com/cloudquery/plugin-sdk/issues/1934)) ([ea8b17a](https://github.com/cloudquery/plugin-sdk/commit/ea8b17a368ddf4762e49397558d32844d96f53dd))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.23.1 ([#1936](https://github.com/cloudquery/plugin-sdk/issues/1936)) ([0b152ba](https://github.com/cloudquery/plugin-sdk/commit/0b152ba7d5d073933e47ecc0913a6a525a282a6a))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.24.0 ([#1937](https://github.com/cloudquery/plugin-sdk/issues/1937)) ([d9e6f47](https://github.com/cloudquery/plugin-sdk/commit/d9e6f478a8d9ebfdda9cc308f1c0facf0ca72c25))

## [4.66.1](https://github.com/cloudquery/plugin-sdk/compare/v4.66.0...v4.66.1) (2024-10-14)


### Bug Fixes

* **deps:** Update aws-sdk-go-v2 monorepo ([#1928](https://github.com/cloudquery/plugin-sdk/issues/1928)) ([75cabcd](https://github.com/cloudquery/plugin-sdk/commit/75cabcd798e5a2fb073b36a93306f353e5b4f447))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.13.1 ([#1931](https://github.com/cloudquery/plugin-sdk/issues/1931)) ([b8a88d0](https://github.com/cloudquery/plugin-sdk/commit/b8a88d079f2e713b0d93c0fd348845f6defe4301))
* **deps:** Update module google.golang.org/protobuf to v1.35.1 ([#1929](https://github.com/cloudquery/plugin-sdk/issues/1929)) ([94a8638](https://github.com/cloudquery/plugin-sdk/commit/94a86387d10e4837aef79b88ef1ad84eb71533a7))

## [4.66.0](https://github.com/cloudquery/plugin-sdk/compare/v4.65.0...v4.66.0) (2024-10-07)


### Features

* Add time.Sleep to mitigate race condition. ([#1923](https://github.com/cloudquery/plugin-sdk/issues/1923)) ([83dfcad](https://github.com/cloudquery/plugin-sdk/commit/83dfcad9fcfa802b38bc4e97587e25218822814b))


### Bug Fixes

* **deps:** Update aws-sdk-go-v2 monorepo ([#1926](https://github.com/cloudquery/plugin-sdk/issues/1926)) ([4fc8896](https://github.com/cloudquery/plugin-sdk/commit/4fc8896e6c72f3fc8fbea2bb569d31cf8b34c961))
* **deps:** Update module google.golang.org/grpc to v1.67.1 ([#1925](https://github.com/cloudquery/plugin-sdk/issues/1925)) ([5e0305d](https://github.com/cloudquery/plugin-sdk/commit/5e0305dd47297e6d0499fbd2b70589b57e17c625))

## [4.65.0](https://github.com/cloudquery/plugin-sdk/compare/v4.64.1...v4.65.0) (2024-10-04)


### Features

* Implement RandomQueue scheduler strategy ([#1914](https://github.com/cloudquery/plugin-sdk/issues/1914)) ([af8ac87](https://github.com/cloudquery/plugin-sdk/commit/af8ac87178cc318d2f31cd17efc7c921d6d52e6b))


### Bug Fixes

* Revert "fix: Error handling in StreamingBatchWriter" ([#1918](https://github.com/cloudquery/plugin-sdk/issues/1918)) ([38b4bfd](https://github.com/cloudquery/plugin-sdk/commit/38b4bfd20e17a00d5a2c83e1d48b8b16270592ba))
* **tests:** WriterTestSuite.handleNulls should not overwrite columns ([#1920](https://github.com/cloudquery/plugin-sdk/issues/1920)) ([08e18e2](https://github.com/cloudquery/plugin-sdk/commit/08e18e265dfb7e6e77c32244f56acd0f63bf4ead))

## [4.64.1](https://github.com/cloudquery/plugin-sdk/compare/v4.64.0...v4.64.1) (2024-10-02)


### Bug Fixes

* Error handling in StreamingBatchWriter ([#1913](https://github.com/cloudquery/plugin-sdk/issues/1913)) ([d852119](https://github.com/cloudquery/plugin-sdk/commit/d8521194dee50d93d74a7156ed607d442ab1db45))

## [4.64.0](https://github.com/cloudquery/plugin-sdk/compare/v4.63.0...v4.64.0) (2024-10-01)


### Features

* Add `opts.SchedulerOpts()` helper to convert `plugin.SyncOptions` for scheduler ([#1900](https://github.com/cloudquery/plugin-sdk/issues/1900)) ([242fb55](https://github.com/cloudquery/plugin-sdk/commit/242fb55088032f65e1e743dcd861b8f05d8d60ce))
* **remoteoauth:** Add `WithToken` option ([#1898](https://github.com/cloudquery/plugin-sdk/issues/1898)) ([ff7a485](https://github.com/cloudquery/plugin-sdk/commit/ff7a485df334cdaa00f8b1a4671595d4fa3fbcdf))
* Update concurrency formula. ([#1907](https://github.com/cloudquery/plugin-sdk/issues/1907)) ([adce99c](https://github.com/cloudquery/plugin-sdk/commit/adce99c9613131a3ef160c9127a5c0d33d33e8af))


### Bug Fixes

* **deps:** Update aws-sdk-go-v2 monorepo ([#1903](https://github.com/cloudquery/plugin-sdk/issues/1903)) ([ce2a0ef](https://github.com/cloudquery/plugin-sdk/commit/ce2a0efa3da3d388be954030153919a5577b586f))
* **deps:** Update aws-sdk-go-v2 monorepo ([#1908](https://github.com/cloudquery/plugin-sdk/issues/1908)) ([bea3b00](https://github.com/cloudquery/plugin-sdk/commit/bea3b00a52b65f65e564e679a202d8fbd8108712))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.22.4 ([#1912](https://github.com/cloudquery/plugin-sdk/issues/1912)) ([c28aabe](https://github.com/cloudquery/plugin-sdk/commit/c28aabeb93fb23432069956d3e3b302bae8b6ed9))
* **deps:** Update module golang.org/x/oauth2 to v0.23.0 ([#1910](https://github.com/cloudquery/plugin-sdk/issues/1910)) ([6fe6414](https://github.com/cloudquery/plugin-sdk/commit/6fe64140337ba8d5c1af795abf64318e6138bdf3))
* **deps:** Update module google.golang.org/grpc to v1.67.0 ([#1904](https://github.com/cloudquery/plugin-sdk/issues/1904)) ([a349812](https://github.com/cloudquery/plugin-sdk/commit/a3498124b325616d085d302fc0faaffb11c77856))
* **deps:** Update opentelemetry-go monorepo ([#1911](https://github.com/cloudquery/plugin-sdk/issues/1911)) ([78e05e1](https://github.com/cloudquery/plugin-sdk/commit/78e05e12bfcb38f675dd83dab0b8b442b6227944))

## [4.63.0](https://github.com/cloudquery/plugin-sdk/compare/v4.62.0...v4.63.0) (2024-09-18)


### Features

* Add sync sharding ([#1891](https://github.com/cloudquery/plugin-sdk/issues/1891)) ([e1823f8](https://github.com/cloudquery/plugin-sdk/commit/e1823f82fd3c457f1f58c266bfd9519b547f31c9))


### Bug Fixes

* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.22.3 ([#1895](https://github.com/cloudquery/plugin-sdk/issues/1895)) ([b05d24b](https://github.com/cloudquery/plugin-sdk/commit/b05d24b345ec519deef70156377338e6b41d8108))
* **deps:** Update module google.golang.org/grpc to v1.66.2 ([#1893](https://github.com/cloudquery/plugin-sdk/issues/1893)) ([6d70b88](https://github.com/cloudquery/plugin-sdk/commit/6d70b88808aa144c4c05e007b291bd8d958858e4))

## [4.62.0](https://github.com/cloudquery/plugin-sdk/compare/v4.61.0...v4.62.0) (2024-09-07)


### Features

* Support Contract Listing For AWS Marketplace ([#1889](https://github.com/cloudquery/plugin-sdk/issues/1889)) ([4654866](https://github.com/cloudquery/plugin-sdk/commit/4654866cb423d237cddb696384e910f59539e1d9))


### Bug Fixes

* **deps:** Update aws-sdk-go-v2 monorepo ([#1890](https://github.com/cloudquery/plugin-sdk/issues/1890)) ([b185e11](https://github.com/cloudquery/plugin-sdk/commit/b185e11bad937fbbeb9178f88f0ede749088efc7))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.22.2 ([#1887](https://github.com/cloudquery/plugin-sdk/issues/1887)) ([a881fac](https://github.com/cloudquery/plugin-sdk/commit/a881fac8976ecfb83101c6268d114e28e19bd2f2))

## [4.61.0](https://github.com/cloudquery/plugin-sdk/compare/v4.60.0...v4.61.0) (2024-09-02)


### Features

* Add remoteoauth helpers (`TokenAuthTransport` and `TokenAuthEditor`) ([#1875](https://github.com/cloudquery/plugin-sdk/issues/1875)) ([bb1be84](https://github.com/cloudquery/plugin-sdk/commit/bb1be8421bbe8086c71c3c02cc4ab281e0eceb5b))
* Add warning on duplicate clients for `round-robin` and `shuffle` schedulers ([#1878](https://github.com/cloudquery/plugin-sdk/issues/1878)) ([d148b94](https://github.com/cloudquery/plugin-sdk/commit/d148b940b09dd832f771a7bf229e4900659d7846))


### Bug Fixes

* **deps:** Update aws-sdk-go-v2 monorepo ([#1872](https://github.com/cloudquery/plugin-sdk/issues/1872)) ([80eb38a](https://github.com/cloudquery/plugin-sdk/commit/80eb38a318bbfd14db2d6a0031e0a2ef467e8a29))
* **deps:** Update golang.org/x/exp digest to 9b4947d ([#1881](https://github.com/cloudquery/plugin-sdk/issues/1881)) ([bbeb846](https://github.com/cloudquery/plugin-sdk/commit/bbeb846aadac0c6f4c8592003a3b4aac2e60b024))
* **deps:** Update module github.com/aws/aws-sdk-go-v2/config to v1.27.30 ([#1876](https://github.com/cloudquery/plugin-sdk/issues/1876)) ([0319ff3](https://github.com/cloudquery/plugin-sdk/commit/0319ff3023b3c79f3463e28f0dfc9a19441d5063))
* **deps:** Update module github.com/aws/aws-sdk-go-v2/config to v1.27.31 ([#1879](https://github.com/cloudquery/plugin-sdk/issues/1879)) ([4dc8f41](https://github.com/cloudquery/plugin-sdk/commit/4dc8f417d986749565e67f9bce0cb172e789d74f))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.12.8 ([#1870](https://github.com/cloudquery/plugin-sdk/issues/1870)) ([96a5194](https://github.com/cloudquery/plugin-sdk/commit/96a51947cd67a22545fb863c4437fe21de170dfb))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.12.9 ([#1873](https://github.com/cloudquery/plugin-sdk/issues/1873)) ([76d4f9f](https://github.com/cloudquery/plugin-sdk/commit/76d4f9f11b8a4f10327d02894ef109e282f1f58b))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.13.0 ([#1874](https://github.com/cloudquery/plugin-sdk/issues/1874)) ([e091d8a](https://github.com/cloudquery/plugin-sdk/commit/e091d8a7091f9d52da068813efacdaa37b7ae0b5))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.22.1 ([#1877](https://github.com/cloudquery/plugin-sdk/issues/1877)) ([11aaab4](https://github.com/cloudquery/plugin-sdk/commit/11aaab425f9182af49bf0d92d5829a70d624b538))
* **deps:** Update module golang.org/x/oauth2 to v0.22.0 ([#1883](https://github.com/cloudquery/plugin-sdk/issues/1883)) ([2a40306](https://github.com/cloudquery/plugin-sdk/commit/2a40306b74e7926078b4576d9f1940e772f0ee1b))
* **deps:** Update module google.golang.org/grpc to v1.66.0 ([#1880](https://github.com/cloudquery/plugin-sdk/issues/1880)) ([a907ea6](https://github.com/cloudquery/plugin-sdk/commit/a907ea632a7e5e0803202a1930222eeeaca50d8e))
* **deps:** Update opentelemetry-go monorepo ([#1884](https://github.com/cloudquery/plugin-sdk/issues/1884)) ([9be63fe](https://github.com/cloudquery/plugin-sdk/commit/9be63feb754ad6503dedb45d2e921aee2c804ade))
* Fix panic when converting schema changes to string ([#1885](https://github.com/cloudquery/plugin-sdk/issues/1885)) ([8274f17](https://github.com/cloudquery/plugin-sdk/commit/8274f172ebf65c085a8d004808404564f7903ffa))

## [4.60.0](https://github.com/cloudquery/plugin-sdk/compare/v4.59.0...v4.60.0) (2024-08-12)


### Features

* Add RemoteOAuth Token helper to refresh `access_token` from cloud environment ([#1866](https://github.com/cloudquery/plugin-sdk/issues/1866)) ([bcd9081](https://github.com/cloudquery/plugin-sdk/commit/bcd9081baf6b1e7311237a8b46e0a13c109ac4ba))


### Bug Fixes

* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.22.0 ([#1864](https://github.com/cloudquery/plugin-sdk/issues/1864)) ([382f980](https://github.com/cloudquery/plugin-sdk/commit/382f98014ae8b72a5493bd06e72d4e1de8398e88))

## [4.59.0](https://github.com/cloudquery/plugin-sdk/compare/v4.58.1...v4.59.0) (2024-08-08)


### Features

* Add basic testing large syncs support ([#1862](https://github.com/cloudquery/plugin-sdk/issues/1862)) ([40a0095](https://github.com/cloudquery/plugin-sdk/commit/40a009574bb3392865a3da7217385c8e389b7a55))


### Bug Fixes

* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.12.7 ([#1860](https://github.com/cloudquery/plugin-sdk/issues/1860)) ([25ed3d2](https://github.com/cloudquery/plugin-sdk/commit/25ed3d25a529a22f351ab92e22fb03a19c9557d4))

## [4.58.1](https://github.com/cloudquery/plugin-sdk/compare/v4.58.0...v4.58.1) (2024-08-03)


### Bug Fixes

* **deps:** Update aws-sdk-go-v2 monorepo ([#1857](https://github.com/cloudquery/plugin-sdk/issues/1857)) ([45a74e8](https://github.com/cloudquery/plugin-sdk/commit/45a74e83f22e2564cb43f121295738d7e10cec6a))

## [4.58.0](https://github.com/cloudquery/plugin-sdk/compare/v4.57.1...v4.58.0) (2024-08-02)


### Features

* Support AWS usage marketplace ([#1770](https://github.com/cloudquery/plugin-sdk/issues/1770)) ([1eb6d1a](https://github.com/cloudquery/plugin-sdk/commit/1eb6d1aabab7db02458f56ea448af750ffb082ae))

## [4.57.1](https://github.com/cloudquery/plugin-sdk/compare/v4.57.0...v4.57.1) (2024-08-02)


### Bug Fixes

* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.12.6 ([#1850](https://github.com/cloudquery/plugin-sdk/issues/1850)) ([4ef35bf](https://github.com/cloudquery/plugin-sdk/commit/4ef35bf659f2a028cf9c46c42c7c9abb496d772b))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.21.5 ([#1852](https://github.com/cloudquery/plugin-sdk/issues/1852)) ([a13ee97](https://github.com/cloudquery/plugin-sdk/commit/a13ee97503ca1312e5547c9d383397c088e6c5d7))

## [4.57.0](https://github.com/cloudquery/plugin-sdk/compare/v4.56.0...v4.57.0) (2024-08-01)


### Features

* Allow setting JSON type schema max depth ([#1844](https://github.com/cloudquery/plugin-sdk/issues/1844)) ([0b28389](https://github.com/cloudquery/plugin-sdk/commit/0b28389bb53cd2c076cca3ddaa93ca4d24e40b7b))


### Bug Fixes

* **deps:** Update golang.org/x/exp digest to 8a7402a ([#1845](https://github.com/cloudquery/plugin-sdk/issues/1845)) ([5f7eb25](https://github.com/cloudquery/plugin-sdk/commit/5f7eb25df3208ed738c1f0e6f17c5366b89fcc30))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.12.5 ([#1841](https://github.com/cloudquery/plugin-sdk/issues/1841)) ([4361e84](https://github.com/cloudquery/plugin-sdk/commit/4361e8442c05b77a0be772cab52e1d217810bf47))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.21.4 ([#1847](https://github.com/cloudquery/plugin-sdk/issues/1847)) ([281b945](https://github.com/cloudquery/plugin-sdk/commit/281b94510552962af13d5c3c2e735669d4fa4bd4))
* **deps:** Update opentelemetry-go monorepo to v1.28.0 ([#1846](https://github.com/cloudquery/plugin-sdk/issues/1846)) ([3a5c90c](https://github.com/cloudquery/plugin-sdk/commit/3a5c90c0045aa6f2df57c02ed2afc7e5596c4bb7))

## [4.56.0](https://github.com/cloudquery/plugin-sdk/compare/v4.55.0...v4.56.0) (2024-07-31)


### Features

* Implement TransformSchema support. ([#1838](https://github.com/cloudquery/plugin-sdk/issues/1838)) ([30875d6](https://github.com/cloudquery/plugin-sdk/commit/30875d6f134f399f5c2ea16dad49b0b5aa4dd3e9))


### Bug Fixes

* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.21.2 ([#1837](https://github.com/cloudquery/plugin-sdk/issues/1837)) ([47bb424](https://github.com/cloudquery/plugin-sdk/commit/47bb424c2151363cc312d155ac5823abfc7d23c5))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.21.3 ([#1840](https://github.com/cloudquery/plugin-sdk/issues/1840)) ([d2c5c7b](https://github.com/cloudquery/plugin-sdk/commit/d2c5c7b54a933a268fe5090a0ca83f1995be9082))

## [4.55.0](https://github.com/cloudquery/plugin-sdk/compare/v4.54.0...v4.55.0) (2024-07-30)


### Features

* Add `PermissionsNeeded` to tables schema ([#1827](https://github.com/cloudquery/plugin-sdk/issues/1827)) ([863b906](https://github.com/cloudquery/plugin-sdk/commit/863b9068bd296dac7c879ae3980a2f2f3ec4c359))


### Bug Fixes

* Handle commas in permissions array ([#1835](https://github.com/cloudquery/plugin-sdk/issues/1835)) ([b633aed](https://github.com/cloudquery/plugin-sdk/commit/b633aed0dc0e6fa8f8af58c8f84e5309375f4608))

## [4.54.0](https://github.com/cloudquery/plugin-sdk/compare/v4.53.1...v4.54.0) (2024-07-30)


### Features

* Add PluginKind "transformer". ([#1828](https://github.com/cloudquery/plugin-sdk/issues/1828)) ([2c78878](https://github.com/cloudquery/plugin-sdk/commit/2c788784e0cb7d2d70ce4389113cd9e758b8e146))


### Bug Fixes

* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.12.4 ([#1830](https://github.com/cloudquery/plugin-sdk/issues/1830)) ([57a606f](https://github.com/cloudquery/plugin-sdk/commit/57a606f774210314f5f19e096be6a18caa5f343a))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.21.0 ([#1829](https://github.com/cloudquery/plugin-sdk/issues/1829)) ([1a614cb](https://github.com/cloudquery/plugin-sdk/commit/1a614cb0dc469bd1c2620c22a91d0c2c4f7779b0))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.21.1 ([#1832](https://github.com/cloudquery/plugin-sdk/issues/1832)) ([33bd55f](https://github.com/cloudquery/plugin-sdk/commit/33bd55f19de89cd877018e7d10ff19c8f5e2a2e2))

## [4.53.1](https://github.com/cloudquery/plugin-sdk/compare/v4.53.0...v4.53.1) (2024-07-26)


### Bug Fixes

* Don't include non exported fields or ignored types in JSON schema ([#1824](https://github.com/cloudquery/plugin-sdk/issues/1824)) ([e97f243](https://github.com/cloudquery/plugin-sdk/commit/e97f2439145961bbae86b09d1e2f1c4ba28af5c4))

## [4.53.0](https://github.com/cloudquery/plugin-sdk/compare/v4.52.1...v4.53.0) (2024-07-25)


### Features

* Add `zerolog.Logger` to `retryablehttp.LeveledLogger` adapter struct ([#1821](https://github.com/cloudquery/plugin-sdk/issues/1821)) ([5c77cee](https://github.com/cloudquery/plugin-sdk/commit/5c77cee87d9fca292e9e81663c9ce3775962a623))

## [4.52.1](https://github.com/cloudquery/plugin-sdk/compare/v4.52.0...v4.52.1) (2024-07-24)


### Bug Fixes

* Don't panic when trying `ToSnake` on `s` character ([#1816](https://github.com/cloudquery/plugin-sdk/issues/1816)) ([30e02da](https://github.com/cloudquery/plugin-sdk/commit/30e02da227bac041cf4ea0c918ad81d360c05084))
* Properly handle map and slice pointers ([#1817](https://github.com/cloudquery/plugin-sdk/issues/1817)) ([8fe9081](https://github.com/cloudquery/plugin-sdk/commit/8fe9081b133892c95fc6dd23223e2a64572a164e))
* Reduce JSON column schema nesting ([#1819](https://github.com/cloudquery/plugin-sdk/issues/1819)) ([2e1112f](https://github.com/cloudquery/plugin-sdk/commit/2e1112fd9fc8d7442f1de21c883287e5d314bb32))

## [4.52.0](https://github.com/cloudquery/plugin-sdk/compare/v4.51.0...v4.52.0) (2024-07-24)


### Features

* Add JSON type schema ([#1796](https://github.com/cloudquery/plugin-sdk/issues/1796)) ([dbc534b](https://github.com/cloudquery/plugin-sdk/commit/dbc534bc54a3f9f02fd3468bfe256b6c46971614))

## [4.51.0](https://github.com/cloudquery/plugin-sdk/compare/v4.50.1...v4.51.0) (2024-07-22)


### Features

* Send plugin logs via OTEL ([#1807](https://github.com/cloudquery/plugin-sdk/issues/1807)) ([9897b83](https://github.com/cloudquery/plugin-sdk/commit/9897b837f25e8d0338ce19fb795a8c96d4bc6223))

## [4.50.1](https://github.com/cloudquery/plugin-sdk/compare/v4.50.0...v4.50.1) (2024-07-22)


### Bug Fixes

* **deps:** Update module github.com/apache/arrow/go/v16 to v17 ([#1809](https://github.com/cloudquery/plugin-sdk/issues/1809)) ([0d6e62d](https://github.com/cloudquery/plugin-sdk/commit/0d6e62d1c25839b2e0472a4e64f0b408a8f51042))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.20.7 ([#1806](https://github.com/cloudquery/plugin-sdk/issues/1806)) ([eca6061](https://github.com/cloudquery/plugin-sdk/commit/eca606147aa94363d8da26dcb8d925d3fb0bdc0d))

## [4.50.0](https://github.com/cloudquery/plugin-sdk/compare/v4.49.4...v4.50.0) (2024-07-19)


### Features

* Implement transformations logic. ([#1800](https://github.com/cloudquery/plugin-sdk/issues/1800)) ([377194e](https://github.com/cloudquery/plugin-sdk/commit/377194e65da224dc11e7e540e4b4e12519de1a95))


### Bug Fixes

* Use trace level for source batcher log message ([#1803](https://github.com/cloudquery/plugin-sdk/issues/1803)) ([6ccf7e6](https://github.com/cloudquery/plugin-sdk/commit/6ccf7e67ae075035f4d2c3b201001a866e6c7624))

## [4.49.4](https://github.com/cloudquery/plugin-sdk/compare/v4.49.3...v4.49.4) (2024-07-19)


### Bug Fixes

* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.12.2 ([#1797](https://github.com/cloudquery/plugin-sdk/issues/1797)) ([98d187b](https://github.com/cloudquery/plugin-sdk/commit/98d187b10c0bda4e54b3f4d09393f35a13df15ce))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.12.3 ([#1801](https://github.com/cloudquery/plugin-sdk/issues/1801)) ([470fbcf](https://github.com/cloudquery/plugin-sdk/commit/470fbcf03991334fe94d20a2e97fe0795f580b94))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.20.5 ([#1798](https://github.com/cloudquery/plugin-sdk/issues/1798)) ([bd584ce](https://github.com/cloudquery/plugin-sdk/commit/bd584ce3998ffb5cb87d66f05886abebee45d7e8))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.20.6 ([#1799](https://github.com/cloudquery/plugin-sdk/issues/1799)) ([d647712](https://github.com/cloudquery/plugin-sdk/commit/d647712917f94020fc5c4b73600b01367baf6d9c))
* **deps:** Update module github.com/santhosh-tekuri/jsonschema/v5 to v6 ([#1782](https://github.com/cloudquery/plugin-sdk/issues/1782)) ([413453c](https://github.com/cloudquery/plugin-sdk/commit/413453c5db66ed27a86acb1340a068e2c4231c78))

## [4.49.3](https://github.com/cloudquery/plugin-sdk/compare/v4.49.2...v4.49.3) (2024-07-09)


### Bug Fixes

* Log OTEL errors as warning level instead of debug ([#1791](https://github.com/cloudquery/plugin-sdk/issues/1791)) ([c7a6179](https://github.com/cloudquery/plugin-sdk/commit/c7a6179fd07cda66fade13f83ee9d9f04094e74b))

## [4.49.2](https://github.com/cloudquery/plugin-sdk/compare/v4.49.1...v4.49.2) (2024-07-09)


### Bug Fixes

* Properly handle relational tables metrics ([#1788](https://github.com/cloudquery/plugin-sdk/issues/1788)) ([ee16898](https://github.com/cloudquery/plugin-sdk/commit/ee168981e13ea8479d95cc0257cf582d0c275183))

## [4.49.1](https://github.com/cloudquery/plugin-sdk/compare/v4.49.0...v4.49.1) (2024-07-08)


### Bug Fixes

* **deps:** Update golang.org/x/exp digest to 7f521ea ([#1778](https://github.com/cloudquery/plugin-sdk/issues/1778)) ([2a1f2d6](https://github.com/cloudquery/plugin-sdk/commit/2a1f2d6b5559403e178038520eea16dd45b9849e))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.12.0 ([#1783](https://github.com/cloudquery/plugin-sdk/issues/1783)) ([812115d](https://github.com/cloudquery/plugin-sdk/commit/812115d04ce38c03b68274a3f53453b3132442e6))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.12.1 ([#1784](https://github.com/cloudquery/plugin-sdk/issues/1784)) ([9cf0394](https://github.com/cloudquery/plugin-sdk/commit/9cf0394018088e382429473d7b4e312770c339c7))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.20.3 ([#1786](https://github.com/cloudquery/plugin-sdk/issues/1786)) ([7b1fc58](https://github.com/cloudquery/plugin-sdk/commit/7b1fc584bd19b1c639bbc32a600b7e21bf438fd1))
* **deps:** Update module github.com/spf13/cobra to v1.8.1 ([#1779](https://github.com/cloudquery/plugin-sdk/issues/1779)) ([e3566a3](https://github.com/cloudquery/plugin-sdk/commit/e3566a3ed202e4d613eb71c2f948aba34dda735b))
* **deps:** Update module google.golang.org/grpc to v1.65.0 ([#1785](https://github.com/cloudquery/plugin-sdk/issues/1785)) ([fcad58f](https://github.com/cloudquery/plugin-sdk/commit/fcad58f3f55c70a735aa3272cb315d2a232ece58))
* Reuse builders on state client flush ([#1777](https://github.com/cloudquery/plugin-sdk/issues/1777)) ([49d43b6](https://github.com/cloudquery/plugin-sdk/commit/49d43b6b5f478697dd88ad83bcf736e7ab8ae7c1))

## [4.49.0](https://github.com/cloudquery/plugin-sdk/compare/v4.48.0...v4.49.0) (2024-06-27)


### Features

* Better OTEL traces, add metrics ([#1751](https://github.com/cloudquery/plugin-sdk/issues/1751)) ([874c33a](https://github.com/cloudquery/plugin-sdk/commit/874c33a9ddc1fcc7aefb86cf1a0076f701f07735))


### Bug Fixes

* **deps:** Update module github.com/hashicorp/go-retryablehttp to v0.7.7 [SECURITY] ([#1774](https://github.com/cloudquery/plugin-sdk/issues/1774)) ([e5e8e7e](https://github.com/cloudquery/plugin-sdk/commit/e5e8e7ea650862e53214f381db6c22173fb04edb))

## [4.48.0](https://github.com/cloudquery/plugin-sdk/compare/v4.47.1...v4.48.0) (2024-06-24)


### Features

* Enable batching resources on source side by default ([#1771](https://github.com/cloudquery/plugin-sdk/issues/1771)) ([1a99a66](https://github.com/cloudquery/plugin-sdk/commit/1a99a66c23cf039e74dac089c8b67d9953653c51))

## [4.47.1](https://github.com/cloudquery/plugin-sdk/compare/v4.47.0...v4.47.1) (2024-06-21)


### Bug Fixes

* Use Atomic Pointer for updating duration metric ([#1766](https://github.com/cloudquery/plugin-sdk/issues/1766)) ([61e698e](https://github.com/cloudquery/plugin-sdk/commit/61e698ed0c094d97411cc14a83ec6b7544c3e83f))

## [4.47.0](https://github.com/cloudquery/plugin-sdk/compare/v4.46.1...v4.47.0) (2024-06-21)


### Features

* Add `duration_ms` to `table sync finished` log message ([#1757](https://github.com/cloudquery/plugin-sdk/issues/1757)) ([9ea034d](https://github.com/cloudquery/plugin-sdk/commit/9ea034daa3b093975ea787cac90785a66953f66c))

## [4.46.1](https://github.com/cloudquery/plugin-sdk/compare/v4.46.0...v4.46.1) (2024-06-21)


### Bug Fixes

* Don't allocate many loggers for source batching ([#1759](https://github.com/cloudquery/plugin-sdk/issues/1759)) ([f29f046](https://github.com/cloudquery/plugin-sdk/commit/f29f0461ec394e2917654f4e8750d659d1c5cbb8))

## [4.46.0](https://github.com/cloudquery/plugin-sdk/compare/v4.45.6...v4.46.0) (2024-06-20)


### Features

* Batch resources into a single record on source side ([#1642](https://github.com/cloudquery/plugin-sdk/issues/1642)) ([f86dcb5](https://github.com/cloudquery/plugin-sdk/commit/f86dcb5513c9ece6a79c24025a07570cddbd5247))

## [4.45.6](https://github.com/cloudquery/plugin-sdk/compare/v4.45.5...v4.45.6) (2024-06-20)


### Bug Fixes

* Account for bytes limit properly when batching records for writing ([#1719](https://github.com/cloudquery/plugin-sdk/issues/1719)) ([25e554e](https://github.com/cloudquery/plugin-sdk/commit/25e554e622e001b7c4c81e4111d70c9143ab29f1))
* **deps:** Update dependency go to v1.21.11 ([#1752](https://github.com/cloudquery/plugin-sdk/issues/1752)) ([abcb2d4](https://github.com/cloudquery/plugin-sdk/commit/abcb2d40cdd6191c900dc9bc50074694471356e7))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.20.2 ([#1754](https://github.com/cloudquery/plugin-sdk/issues/1754)) ([6288710](https://github.com/cloudquery/plugin-sdk/commit/6288710b43da942854766a3e58b492f2f26e5d72))

## [4.45.5](https://github.com/cloudquery/plugin-sdk/compare/v4.45.4...v4.45.5) (2024-06-19)


### Bug Fixes

* Send parent in migrate message ([#1742](https://github.com/cloudquery/plugin-sdk/issues/1742)) ([a862f4a](https://github.com/cloudquery/plugin-sdk/commit/a862f4a2488b794cf5346977716262fa2da6ca9e))

## [4.45.4](https://github.com/cloudquery/plugin-sdk/compare/v4.45.3...v4.45.4) (2024-06-19)


### Bug Fixes

* Revert "fix: Allow marshaling plain string values into JSON scalars" ([#1746](https://github.com/cloudquery/plugin-sdk/issues/1746)) ([096bd88](https://github.com/cloudquery/plugin-sdk/commit/096bd88863f1356e8bb01e873cb08b82ae4d4363))

## [4.45.3](https://github.com/cloudquery/plugin-sdk/compare/v4.45.2...v4.45.3) (2024-06-19)


### Bug Fixes

* Allow marshaling plain string values into JSON scalars ([#1743](https://github.com/cloudquery/plugin-sdk/issues/1743)) ([87e90b8](https://github.com/cloudquery/plugin-sdk/commit/87e90b843bab5bfa281c568b2dd94513d195d11b))

## [4.45.2](https://github.com/cloudquery/plugin-sdk/compare/v4.45.1...v4.45.2) (2024-06-17)


### Bug Fixes

* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.20.1 ([#1739](https://github.com/cloudquery/plugin-sdk/issues/1739)) ([cdb1b6b](https://github.com/cloudquery/plugin-sdk/commit/cdb1b6b9688c10822a184787a107bb02a4e2ebaf))
* **deps:** Update module google.golang.org/protobuf to v1.34.2 ([#1737](https://github.com/cloudquery/plugin-sdk/issues/1737)) ([ced9333](https://github.com/cloudquery/plugin-sdk/commit/ced933389a2ce661a85bf7d73b4fffd215fd2a74))
* Remove no sentry deprecation warning ([#1740](https://github.com/cloudquery/plugin-sdk/issues/1740)) ([1fb06bc](https://github.com/cloudquery/plugin-sdk/commit/1fb06bc46db3de657682f474a5a00e3d160b9823))

## [4.45.1](https://github.com/cloudquery/plugin-sdk/compare/v4.45.0...v4.45.1) (2024-06-14)


### Bug Fixes

* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.20.0 ([#1734](https://github.com/cloudquery/plugin-sdk/issues/1734)) ([a308b19](https://github.com/cloudquery/plugin-sdk/commit/a308b19a460ecd0dc96bb3d2cda80542ab502292))

## [4.45.0](https://github.com/cloudquery/plugin-sdk/compare/v4.44.2...v4.45.0) (2024-06-14)


### Features

* Remove plugin option for logging error events to Sentry ([#1724](https://github.com/cloudquery/plugin-sdk/issues/1724)) ([7732fe8](https://github.com/cloudquery/plugin-sdk/commit/7732fe898d2ce2d6579ff9fc8165551e042c3d33))


### Bug Fixes

* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.19.19 ([#1726](https://github.com/cloudquery/plugin-sdk/issues/1726)) ([a1dd044](https://github.com/cloudquery/plugin-sdk/commit/a1dd04437f51a0167b430482aa1d152118672320))
* Don't include other relation siblings if not specified in config ([#1720](https://github.com/cloudquery/plugin-sdk/issues/1720)) ([f730ec5](https://github.com/cloudquery/plugin-sdk/commit/f730ec52f565a89736364dff78b9d78f6ed02507))

## [4.44.2](https://github.com/cloudquery/plugin-sdk/compare/v4.44.1...v4.44.2) (2024-06-03)


### Bug Fixes

* **deps:** Update golang.org/x/exp digest to fd00a4e ([#1708](https://github.com/cloudquery/plugin-sdk/issues/1708)) ([93866a9](https://github.com/cloudquery/plugin-sdk/commit/93866a9b94f8d93ff56fd14e5f44f54b18b6f531))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.11.3 ([#1716](https://github.com/cloudquery/plugin-sdk/issues/1716)) ([36c97c8](https://github.com/cloudquery/plugin-sdk/commit/36c97c819d45cc0d41abe1c9e4afdd4ec6004c2a))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.19.18 ([#1717](https://github.com/cloudquery/plugin-sdk/issues/1717)) ([f36d5d2](https://github.com/cloudquery/plugin-sdk/commit/f36d5d2947de236f5a601be99c35977af7e4e6c7))
* **deps:** Update module github.com/getsentry/sentry-go to v0.28.0 ([#1712](https://github.com/cloudquery/plugin-sdk/issues/1712)) ([82d78cb](https://github.com/cloudquery/plugin-sdk/commit/82d78cbf2d552907171a192e9b1490c41907ff25))
* **deps:** Update module github.com/goccy/go-json to v0.10.3 ([#1709](https://github.com/cloudquery/plugin-sdk/issues/1709)) ([32a2dca](https://github.com/cloudquery/plugin-sdk/commit/32a2dca6919bd5ebd2fd18405bd5f3ed22d4fc47))
* **deps:** Update module github.com/rs/zerolog to v1.33.0 ([#1713](https://github.com/cloudquery/plugin-sdk/issues/1713)) ([a09376d](https://github.com/cloudquery/plugin-sdk/commit/a09376d7b8a4904d5cb90e3b7613e38b208ed217))
* **deps:** Update opentelemetry-go monorepo to v1.27.0 ([#1714](https://github.com/cloudquery/plugin-sdk/issues/1714)) ([4f29cf1](https://github.com/cloudquery/plugin-sdk/commit/4f29cf1bf50f3ccbea1d8039cd4addc0fb758d77))

## [4.44.1](https://github.com/cloudquery/plugin-sdk/compare/v4.44.0...v4.44.1) (2024-05-31)


### Bug Fixes

* Added support for list pointers ([#1705](https://github.com/cloudquery/plugin-sdk/issues/1705)) ([0368a01](https://github.com/cloudquery/plugin-sdk/commit/0368a0117e31f8baff83c408eabda93a874edf9e))

## [4.44.0](https://github.com/cloudquery/plugin-sdk/compare/v4.43.1...v4.44.0) (2024-05-24)


### Features

* Enable `NewConnectedClientWithOptions` to set `ClientOptions` ([#1700](https://github.com/cloudquery/plugin-sdk/issues/1700)) ([8797a18](https://github.com/cloudquery/plugin-sdk/commit/8797a182be8f6d667f6c567c4a6f9132402ebf00))


### Bug Fixes

* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.19.16 ([#1699](https://github.com/cloudquery/plugin-sdk/issues/1699)) ([3b15ac6](https://github.com/cloudquery/plugin-sdk/commit/3b15ac6730b2d36898b8ea418bf8ee15414356f6))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.19.17 ([#1703](https://github.com/cloudquery/plugin-sdk/issues/1703)) ([7501fdd](https://github.com/cloudquery/plugin-sdk/commit/7501fdd5a6c7c484eac4e6ade3e90f762ea7855f))

## [4.43.1](https://github.com/cloudquery/plugin-sdk/compare/v4.43.0...v4.43.1) (2024-05-20)


### Bug Fixes

* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.19.14 ([#1696](https://github.com/cloudquery/plugin-sdk/issues/1696)) ([4f1f3f8](https://github.com/cloudquery/plugin-sdk/commit/4f1f3f8fa56eafd20c9df08ef587fe2a60d80daa))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.19.15 ([#1697](https://github.com/cloudquery/plugin-sdk/issues/1697)) ([0135160](https://github.com/cloudquery/plugin-sdk/commit/0135160f1f2bb5805edef62707104e1757138a95))
* **deps:** Update module google.golang.org/grpc to v1.64.0 ([#1692](https://github.com/cloudquery/plugin-sdk/issues/1692)) ([f9e2053](https://github.com/cloudquery/plugin-sdk/commit/f9e20536d0abd4f5ae8cac67b17af04c4ae6faa9))

## [4.43.0](https://github.com/cloudquery/plugin-sdk/compare/v4.42.2...v4.43.0) (2024-05-20)


### Features

* Add test connection ([#1682](https://github.com/cloudquery/plugin-sdk/issues/1682)) ([03493f5](https://github.com/cloudquery/plugin-sdk/commit/03493f5c9e63c8d50ff9790c1bd77ac71dfd9139))


### Bug Fixes

* **deps:** Update module github.com/apache/arrow/go/v16 to v16.1.0 ([#1693](https://github.com/cloudquery/plugin-sdk/issues/1693)) ([461e352](https://github.com/cloudquery/plugin-sdk/commit/461e3520304d41a19e0152bf8c58f35546223022))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.11.2 ([#1687](https://github.com/cloudquery/plugin-sdk/issues/1687)) ([30b52f7](https://github.com/cloudquery/plugin-sdk/commit/30b52f7561b635a75f8e789ab721ec69ed57922c))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.19.13 ([#1690](https://github.com/cloudquery/plugin-sdk/issues/1690)) ([c35be5d](https://github.com/cloudquery/plugin-sdk/commit/c35be5d074b2e13a6b0371a12b05cb7538a15065))

## [4.42.2](https://github.com/cloudquery/plugin-sdk/compare/v4.42.1...v4.42.2) (2024-05-17)


### Bug Fixes

* Remove JSON schema validation warnings ([#1685](https://github.com/cloudquery/plugin-sdk/issues/1685)) ([c6f39f4](https://github.com/cloudquery/plugin-sdk/commit/c6f39f4da587fd5b8397994cad35c8e81a80de4c))

## [4.42.1](https://github.com/cloudquery/plugin-sdk/compare/v4.42.0...v4.42.1) (2024-05-15)


### Bug Fixes

* Correct error message on Read failure ([#1680](https://github.com/cloudquery/plugin-sdk/issues/1680)) ([dc31c3a](https://github.com/cloudquery/plugin-sdk/commit/dc31c3aa8639250df277de8b857727b88e503d3b))
* Properly handle records with multiple rows in batching ([#1647](https://github.com/cloudquery/plugin-sdk/issues/1647)) ([926a7fc](https://github.com/cloudquery/plugin-sdk/commit/926a7fc97ae08adf6be30662f3ece58157294f22))

## [4.42.0](https://github.com/cloudquery/plugin-sdk/compare/v4.41.1...v4.42.0) (2024-05-11)


### Features

* Re-configure batch updater using response headers ([#1677](https://github.com/cloudquery/plugin-sdk/issues/1677)) ([e6313f9](https://github.com/cloudquery/plugin-sdk/commit/e6313f9bf20677cdbcce75875946c33c05dc8fc7))


### Bug Fixes

* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.11.1 ([#1675](https://github.com/cloudquery/plugin-sdk/issues/1675)) ([6cd035d](https://github.com/cloudquery/plugin-sdk/commit/6cd035dfc9ad9a3ae6c35a5006d5b67c97685786))
* **deps:** Update module google.golang.org/protobuf to v1.34.1 ([#1678](https://github.com/cloudquery/plugin-sdk/issues/1678)) ([49cadc8](https://github.com/cloudquery/plugin-sdk/commit/49cadc89eb01e46e13520c613b58fe356a2f12c0))

## [4.41.1](https://github.com/cloudquery/plugin-sdk/compare/v4.41.0...v4.41.1) (2024-05-09)


### Bug Fixes

* Expose IncreaseForTable on UsageClient interface ([#1672](https://github.com/cloudquery/plugin-sdk/issues/1672)) ([52c145c](https://github.com/cloudquery/plugin-sdk/commit/52c145c0be55023eb044ac32d1c099f7b1f6bd25))

## [4.41.0](https://github.com/cloudquery/plugin-sdk/compare/v4.40.2...v4.41.0) (2024-05-09)


### Features

* Allow reporting usage as breakdown by table ([#1668](https://github.com/cloudquery/plugin-sdk/issues/1668)) ([0a93aec](https://github.com/cloudquery/plugin-sdk/commit/0a93aecd8ad365b54d03034f6984d6228f134b53))


### Bug Fixes

* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.10.0 ([#1669](https://github.com/cloudquery/plugin-sdk/issues/1669)) ([7068bcb](https://github.com/cloudquery/plugin-sdk/commit/7068bcb7b2b6f7c310e5c029c27868dc2646c798))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.11.0 ([#1670](https://github.com/cloudquery/plugin-sdk/issues/1670)) ([32a78c9](https://github.com/cloudquery/plugin-sdk/commit/32a78c975f51e4dcf1c4a950dd7b51ff2a3175cf))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.9.2 ([#1665](https://github.com/cloudquery/plugin-sdk/issues/1665)) ([bdbc8ca](https://github.com/cloudquery/plugin-sdk/commit/bdbc8ca6488a541b19256d91a4f3fd30ffd0f035))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.19.12 ([#1667](https://github.com/cloudquery/plugin-sdk/issues/1667)) ([36024dc](https://github.com/cloudquery/plugin-sdk/commit/36024dc799e734cff0adb587532010f69cec5c87))

## [4.40.2](https://github.com/cloudquery/plugin-sdk/compare/v4.40.1...v4.40.2) (2024-05-06)


### Bug Fixes

* **deps:** Upgrade `github.com/apache/arrow/go` to `v16` ([#1661](https://github.com/cloudquery/plugin-sdk/issues/1661)) ([04d9585](https://github.com/cloudquery/plugin-sdk/commit/04d95859f000e2fc2823bb5b470ccd6a7d117bb7))

## [4.40.1](https://github.com/cloudquery/plugin-sdk/compare/v4.40.0...v4.40.1) (2024-05-06)


### Bug Fixes

* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.19.10 ([#1658](https://github.com/cloudquery/plugin-sdk/issues/1658)) ([cf1c5a0](https://github.com/cloudquery/plugin-sdk/commit/cf1c5a08dfa93a4c7a1a50e1a24ad3d25c4ac85d))
* **deps:** Update module golang.org/x/exp to v0.0.0-20240416160154-fe59bbe5cc7f ([#1653](https://github.com/cloudquery/plugin-sdk/issues/1653)) ([e759eac](https://github.com/cloudquery/plugin-sdk/commit/e759eacb2a1a97bdca4fa7c92d81a22a540199ea))
* **deps:** Update module google.golang.org/protobuf to v1.34.0 ([#1657](https://github.com/cloudquery/plugin-sdk/issues/1657)) ([b492bdc](https://github.com/cloudquery/plugin-sdk/commit/b492bdc86d913a264353e06391c01fbba6e6d3aa))
* **deps:** Update opentelemetry-go monorepo to v1.26.0 ([#1654](https://github.com/cloudquery/plugin-sdk/issues/1654)) ([4ea5b0d](https://github.com/cloudquery/plugin-sdk/commit/4ea5b0d71a4f2e280fdc0486f66f5c59e23461ea))

## [4.40.0](https://github.com/cloudquery/plugin-sdk/compare/v4.39.1...v4.40.0) (2024-04-29)


### Features

* Add function to create state client and gRPC backend connection ([#1650](https://github.com/cloudquery/plugin-sdk/issues/1650)) ([e150c58](https://github.com/cloudquery/plugin-sdk/commit/e150c5813de51175d10a4d79ec331f724e2593fc))

## [4.39.1](https://github.com/cloudquery/plugin-sdk/compare/v4.39.0...v4.39.1) (2024-04-22)


### Bug Fixes

* Fix link to billing page ([#1643](https://github.com/cloudquery/plugin-sdk/issues/1643)) ([ca216b6](https://github.com/cloudquery/plugin-sdk/commit/ca216b6dbacfb16afc587d03bdf84fef796eff33))
* Use `clear` for mixed batch writer ([#1645](https://github.com/cloudquery/plugin-sdk/issues/1645)) ([07945ac](https://github.com/cloudquery/plugin-sdk/commit/07945ac593e2965b8f80ed3b4ed8724ac610f2bd))

## [4.39.0](https://github.com/cloudquery/plugin-sdk/compare/v4.38.2...v4.39.0) (2024-04-19)


### Features

* Fill in `CC` and `CXX` environment variables for `package` command ([#1637](https://github.com/cloudquery/plugin-sdk/issues/1637)) ([c0282e1](https://github.com/cloudquery/plugin-sdk/commit/c0282e11e873112a563406678fb2da5b85ee006c))

## [4.38.2](https://github.com/cloudquery/plugin-sdk/compare/v4.38.1...v4.38.2) (2024-04-15)


### Bug Fixes

* **deps:** Update opentelemetry-go monorepo to v1.25.0 ([#1634](https://github.com/cloudquery/plugin-sdk/issues/1634)) ([c1477d5](https://github.com/cloudquery/plugin-sdk/commit/c1477d544a6000f45ae30f019c17e22fc679ff13))

## [4.38.1](https://github.com/cloudquery/plugin-sdk/compare/v4.38.0...v4.38.1) (2024-04-12)


### Bug Fixes

* **deps:** Update golang.org/x/exp digest to c0f41cb ([#1615](https://github.com/cloudquery/plugin-sdk/issues/1615)) ([0c21bfb](https://github.com/cloudquery/plugin-sdk/commit/0c21bfbc0faed5c12e415820dc95c1f1cb2c8e7d))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.19.9 ([#1631](https://github.com/cloudquery/plugin-sdk/issues/1631)) ([5d45003](https://github.com/cloudquery/plugin-sdk/commit/5d450034a835d2e9620310394bc0cf776d44c337))
* **deps:** Update module google.golang.org/grpc to v1.63.2 ([#1617](https://github.com/cloudquery/plugin-sdk/issues/1617)) ([02461b1](https://github.com/cloudquery/plugin-sdk/commit/02461b1ba112a87dc97575320fe5327383e8bae9))
* **test:** Slice rows properly when reading in tests ([#1632](https://github.com/cloudquery/plugin-sdk/issues/1632)) ([537b64c](https://github.com/cloudquery/plugin-sdk/commit/537b64c6ca2fbe2c7208e152fb4f77015fc875f1))

## [4.38.0](https://github.com/cloudquery/plugin-sdk/compare/v4.37.0...v4.38.0) (2024-04-08)


### Features

* Support arbitrary map values for structs (scalar) ([#1611](https://github.com/cloudquery/plugin-sdk/issues/1611)) ([d8fde8c](https://github.com/cloudquery/plugin-sdk/commit/d8fde8c87a82617dc180c273cfd9cb70dcddfe13))
* Test duplicated primary key insertion ([#1584](https://github.com/cloudquery/plugin-sdk/issues/1584)) ([6c57402](https://github.com/cloudquery/plugin-sdk/commit/6c57402388df153d752d85a4f7793499a36a78bd))

## [4.37.0](https://github.com/cloudquery/plugin-sdk/compare/v4.36.5...v4.37.0) (2024-04-05)


### Features

* Add versioning to state client ([#1604](https://github.com/cloudquery/plugin-sdk/issues/1604)) ([8957223](https://github.com/cloudquery/plugin-sdk/commit/89572235c39042d1dc03469c2819909e26b3ca17))

## [4.36.5](https://github.com/cloudquery/plugin-sdk/compare/v4.36.4...v4.36.5) (2024-04-04)


### Bug Fixes

* Update Otel Schema version ([#1605](https://github.com/cloudquery/plugin-sdk/issues/1605)) ([601ef35](https://github.com/cloudquery/plugin-sdk/commit/601ef352b4871c0396ad98b3b41aa30636f6924f))

## [4.36.4](https://github.com/cloudquery/plugin-sdk/compare/v4.36.3...v4.36.4) (2024-04-04)


### Bug Fixes

* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.9.1 ([#1600](https://github.com/cloudquery/plugin-sdk/issues/1600)) ([34d4501](https://github.com/cloudquery/plugin-sdk/commit/34d45014e7838e029bc88978fb783c07b1a7228e))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.19.8 ([#1602](https://github.com/cloudquery/plugin-sdk/issues/1602)) ([477be38](https://github.com/cloudquery/plugin-sdk/commit/477be386acab9ecf2af3aec0aa09f9db0b1e3674))

## [4.36.3](https://github.com/cloudquery/plugin-sdk/compare/v4.36.2...v4.36.3) (2024-04-01)


### Bug Fixes

* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.19.7 ([#1596](https://github.com/cloudquery/plugin-sdk/issues/1596)) ([2c7f34c](https://github.com/cloudquery/plugin-sdk/commit/2c7f34c7bd62eee5e2d8259c5c8f9cc38fea8334))

## [4.36.2](https://github.com/cloudquery/plugin-sdk/compare/v4.36.1...v4.36.2) (2024-04-01)


### Bug Fixes

* **deps:** Update golang.org/x/exp digest to a685a6e ([#1585](https://github.com/cloudquery/plugin-sdk/issues/1585)) ([824e745](https://github.com/cloudquery/plugin-sdk/commit/824e7455a6be58b59cc6d322216e3bea00738269))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.19.6 ([#1594](https://github.com/cloudquery/plugin-sdk/issues/1594)) ([dd25ea9](https://github.com/cloudquery/plugin-sdk/commit/dd25ea9b4d3f0552491c83d6034834499cac6f49))
* **deps:** Update module github.com/getsentry/sentry-go to v0.27.0 ([#1588](https://github.com/cloudquery/plugin-sdk/issues/1588)) ([88ec704](https://github.com/cloudquery/plugin-sdk/commit/88ec704fda1dc4fa599e7bcfdc3cfb5a27bf13e4))
* **deps:** Update module github.com/grpc-ecosystem/go-grpc-middleware/v2 to v2.1.0 ([#1589](https://github.com/cloudquery/plugin-sdk/issues/1589)) ([5dfa082](https://github.com/cloudquery/plugin-sdk/commit/5dfa0829d476c0a3f958a37fbafbef39b704127e))
* **deps:** Update module github.com/invopop/jsonschema to v0.12.0 ([#1590](https://github.com/cloudquery/plugin-sdk/issues/1590)) ([3e71418](https://github.com/cloudquery/plugin-sdk/commit/3e7141855dde7a746dd0111d58f2af4b015d0feb))
* **deps:** Update module github.com/rs/zerolog to v1.32.0 ([#1591](https://github.com/cloudquery/plugin-sdk/issues/1591)) ([5331564](https://github.com/cloudquery/plugin-sdk/commit/5331564babe505c2145329e47a67adf126b25f0c))
* **deps:** Update module github.com/spf13/cobra to v1.8.0 ([#1592](https://github.com/cloudquery/plugin-sdk/issues/1592)) ([fc8558b](https://github.com/cloudquery/plugin-sdk/commit/fc8558b20e90a8c7eaeff86c177dcf09cf81f63a))
* **deps:** Update module github.com/stretchr/testify to v1.9.0 ([#1593](https://github.com/cloudquery/plugin-sdk/issues/1593)) ([59cc967](https://github.com/cloudquery/plugin-sdk/commit/59cc9677f363a92f635852b0e711c9136315d30d))

## [4.36.1](https://github.com/cloudquery/plugin-sdk/compare/v4.36.0...v4.36.1) (2024-03-28)


### Bug Fixes

* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.9.0 ([#1578](https://github.com/cloudquery/plugin-sdk/issues/1578)) ([f8d350a](https://github.com/cloudquery/plugin-sdk/commit/f8d350a50d9b01b88321bcefabde7795fdcf00b6))

## [4.36.0](https://github.com/cloudquery/plugin-sdk/compare/v4.35.0...v4.36.0) (2024-03-25)


### Features

* Expose InvocationID to Plugin Client ([#1571](https://github.com/cloudquery/plugin-sdk/issues/1571)) ([038e401](https://github.com/cloudquery/plugin-sdk/commit/038e401e37062ef82d7c3439dbbdd998ab520fab))

## [4.35.0](https://github.com/cloudquery/plugin-sdk/compare/v4.34.2...v4.35.0) (2024-03-22)


### Features

* Handle unknown token types when getting team name in usage client ([#1572](https://github.com/cloudquery/plugin-sdk/issues/1572)) ([b6cb796](https://github.com/cloudquery/plugin-sdk/commit/b6cb79643a10bd79016478ee74629e8db6d16031))

## [4.34.2](https://github.com/cloudquery/plugin-sdk/compare/v4.34.1...v4.34.2) (2024-03-18)


### Bug Fixes

* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.8.1 ([#1567](https://github.com/cloudquery/plugin-sdk/issues/1567)) ([d6f5c18](https://github.com/cloudquery/plugin-sdk/commit/d6f5c18aad252a1a82451362d6834238360baa0f))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.19.4 ([#1569](https://github.com/cloudquery/plugin-sdk/issues/1569)) ([e4895d3](https://github.com/cloudquery/plugin-sdk/commit/e4895d3c17bc4c01e81b620d19d88527f3fb1bf3))

## [4.34.1](https://github.com/cloudquery/plugin-sdk/compare/v4.34.0...v4.34.1) (2024-03-15)


### Bug Fixes

* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.8.0 ([#1563](https://github.com/cloudquery/plugin-sdk/issues/1563)) ([abf3794](https://github.com/cloudquery/plugin-sdk/commit/abf37940ef2b413774d452dce907003f9deb7ff6))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.19.3 ([#1565](https://github.com/cloudquery/plugin-sdk/issues/1565)) ([5107ad0](https://github.com/cloudquery/plugin-sdk/commit/5107ad08419f4c3f5d46e0c68da942d0ebba41a5))

## [4.34.0](https://github.com/cloudquery/plugin-sdk/compare/v4.33.0...v4.34.0) (2024-03-15)


### Features

* Enable destinations to completely skip migration that are not supported ([#1560](https://github.com/cloudquery/plugin-sdk/issues/1560)) ([3d3479b](https://github.com/cloudquery/plugin-sdk/commit/3d3479bcf7f4b10a50f19bf32af4aa71361b526f))

## [4.33.0](https://github.com/cloudquery/plugin-sdk/compare/v4.32.1...v4.33.0) (2024-03-13)


### Features

* Add destination tests for removing a unique constraint ([#1558](https://github.com/cloudquery/plugin-sdk/issues/1558)) ([8add2b3](https://github.com/cloudquery/plugin-sdk/commit/8add2b36b8a9bf37bd24bc6cb03597c9843a592e))


### Bug Fixes

* **deps:** Update Google Golang modules ([#1556](https://github.com/cloudquery/plugin-sdk/issues/1556)) ([e89d4ce](https://github.com/cloudquery/plugin-sdk/commit/e89d4cea569abf81973d14db600a00db1d5133f3))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.19.2 ([#1554](https://github.com/cloudquery/plugin-sdk/issues/1554)) ([09c24f1](https://github.com/cloudquery/plugin-sdk/commit/09c24f1e6b2f53c1a40df65bb7e0cea7b66c3722))

## [4.32.1](https://github.com/cloudquery/plugin-sdk/compare/v4.32.0...v4.32.1) (2024-03-06)


### Bug Fixes

* **deps:** Update golang.org/x/exp digest to 814bf88 ([#1540](https://github.com/cloudquery/plugin-sdk/issues/1540)) ([e80fb24](https://github.com/cloudquery/plugin-sdk/commit/e80fb24ad916e84e391595ed482b4285ea5e1a9c))
* **deps:** Update google.golang.org/genproto/googleapis/api digest to df926f6 ([#1541](https://github.com/cloudquery/plugin-sdk/issues/1541)) ([9d8a3ec](https://github.com/cloudquery/plugin-sdk/commit/9d8a3ec5c7a4bffe3e625f148de43d71c836794d))
* **deps:** Update google.golang.org/genproto/googleapis/rpc digest to df926f6 ([#1543](https://github.com/cloudquery/plugin-sdk/issues/1543)) ([9315c16](https://github.com/cloudquery/plugin-sdk/commit/9315c1639e02474e97670ffb6c9b198b63aec5ef))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.19.1 ([#1549](https://github.com/cloudquery/plugin-sdk/issues/1549)) ([3112739](https://github.com/cloudquery/plugin-sdk/commit/3112739d2a897b784ed85f27ee7632f5fbcb5091))
* **deps:** Update module github.com/klauspost/compress to v1.17.7 ([#1544](https://github.com/cloudquery/plugin-sdk/issues/1544)) ([4e04027](https://github.com/cloudquery/plugin-sdk/commit/4e04027488cb1c32830d5fd14440beabf4a07500))
* **deps:** Update module github.com/klauspost/cpuid/v2 to v2.2.7 ([#1545](https://github.com/cloudquery/plugin-sdk/issues/1545)) ([0fff7ed](https://github.com/cloudquery/plugin-sdk/commit/0fff7ed4464ac572e00eb5e0dc289e467b8e7afb))
* **deps:** Update module github.com/tdewolff/minify/v2 to v2.20.18 ([#1546](https://github.com/cloudquery/plugin-sdk/issues/1546)) ([45fa641](https://github.com/cloudquery/plugin-sdk/commit/45fa641b50f177d2ab01298b0c14fc764464fcd7))
* **deps:** Update module github.com/ugorji/go/codec to v1.2.12 ([#1547](https://github.com/cloudquery/plugin-sdk/issues/1547)) ([cd3488a](https://github.com/cloudquery/plugin-sdk/commit/cd3488ab730499dd513996d73987c9b86fca34c0))
* **deps:** Update module google.golang.org/grpc to v1.62.0 ([#1550](https://github.com/cloudquery/plugin-sdk/issues/1550)) ([9ccec98](https://github.com/cloudquery/plugin-sdk/commit/9ccec989cd143e685fd7d3f66d840c2e2cb8d74b))
* **deps:** Update module google.golang.org/grpc to v1.62.0 ([#1551](https://github.com/cloudquery/plugin-sdk/issues/1551)) ([d907120](https://github.com/cloudquery/plugin-sdk/commit/d907120661cb2ebead90c68b0f1a42767112bba3))
* MixedBatchWriter should nil the slice instead of zeroing ([#1553](https://github.com/cloudquery/plugin-sdk/issues/1553)) ([f565da8](https://github.com/cloudquery/plugin-sdk/commit/f565da8961db0b9f88efcdaa6f083faa789de324))

## [4.32.0](https://github.com/cloudquery/plugin-sdk/compare/v4.31.0...v4.32.0) (2024-02-28)


### Features

* Skip table validation during init ([#1536](https://github.com/cloudquery/plugin-sdk/issues/1536)) ([fb09f20](https://github.com/cloudquery/plugin-sdk/commit/fb09f20cdc3603bcbfdff6fc060c49992ee4e881))


### Bug Fixes

* Remove unreachable code ([#1537](https://github.com/cloudquery/plugin-sdk/issues/1537)) ([6cae5a4](https://github.com/cloudquery/plugin-sdk/commit/6cae5a44323169d568c0c7a8e236bbcac83d7498))

## [4.31.0](https://github.com/cloudquery/plugin-sdk/compare/v4.30.0...v4.31.0) (2024-02-27)


### Features

* Allow homogeneous data types to be configured ([#1533](https://github.com/cloudquery/plugin-sdk/issues/1533)) ([ca7cdb8](https://github.com/cloudquery/plugin-sdk/commit/ca7cdb8b150900a315a694626d394775bcfc6b90))


### Bug Fixes

* Default Plugin logger assumes plugin is a `source` ([#1531](https://github.com/cloudquery/plugin-sdk/issues/1531)) ([b7dcd56](https://github.com/cloudquery/plugin-sdk/commit/b7dcd56e25abfea5992f4746910d5c39ce93e121))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.18.0 ([#1528](https://github.com/cloudquery/plugin-sdk/issues/1528)) ([4cc6ade](https://github.com/cloudquery/plugin-sdk/commit/4cc6adeb4edfb9bf8b8b51716ceefec284d43548))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.19.0 ([#1532](https://github.com/cloudquery/plugin-sdk/issues/1532)) ([4b475bb](https://github.com/cloudquery/plugin-sdk/commit/4b475bbd5d61e533fd0aad1d80f7dfe5e58b039d))
* Support list scalars from JSON ([#1530](https://github.com/cloudquery/plugin-sdk/issues/1530)) ([cf13dd5](https://github.com/cloudquery/plugin-sdk/commit/cf13dd56e0d54a1c6d5f72d7590991d8c676a233))

## [4.30.0](https://github.com/cloudquery/plugin-sdk/compare/v4.29.1...v4.30.0) (2024-02-16)


### Features

* Enhance test suite ([#1523](https://github.com/cloudquery/plugin-sdk/issues/1523)) ([668a297](https://github.com/cloudquery/plugin-sdk/commit/668a29752331c54208bad5e4e5ddfeb90c15f52f))
* Implement `GetSpecSchema` call ([#1521](https://github.com/cloudquery/plugin-sdk/issues/1521)) ([87bea95](https://github.com/cloudquery/plugin-sdk/commit/87bea95367b6e70335e788c410dd982c70c04dd4))
* Support offline licensing for all plugins from a specific team ([#1517](https://github.com/cloudquery/plugin-sdk/issues/1517)) ([d3755dd](https://github.com/cloudquery/plugin-sdk/commit/d3755dd40df0a0addb52ba30e4a0793848416d6d))


### Bug Fixes

* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.7.3 ([#1516](https://github.com/cloudquery/plugin-sdk/issues/1516)) ([54baf21](https://github.com/cloudquery/plugin-sdk/commit/54baf21490d3843931ffb5b39f8caf79cb069db0))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.7.4 ([#1524](https://github.com/cloudquery/plugin-sdk/issues/1524)) ([e1a3f77](https://github.com/cloudquery/plugin-sdk/commit/e1a3f779776fe87606a975bebda3a19f1ddd0a3e))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.7.5 ([#1525](https://github.com/cloudquery/plugin-sdk/issues/1525)) ([c1fae76](https://github.com/cloudquery/plugin-sdk/commit/c1fae76f2694e07964fa1562d08c93641b46a940))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.16.8 ([#1514](https://github.com/cloudquery/plugin-sdk/issues/1514)) ([5b43629](https://github.com/cloudquery/plugin-sdk/commit/5b43629100296bb6a5e687b4a0fb15491a1b0e35))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.17.0 ([#1519](https://github.com/cloudquery/plugin-sdk/issues/1519)) ([209b081](https://github.com/cloudquery/plugin-sdk/commit/209b081e11ac25e36a496afd1b054a3f8a45a290))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.17.1 ([#1520](https://github.com/cloudquery/plugin-sdk/issues/1520)) ([b858608](https://github.com/cloudquery/plugin-sdk/commit/b858608a89f2d509287b3ffd76dd6ad0b63c3c0f))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.17.2 ([#1526](https://github.com/cloudquery/plugin-sdk/issues/1526)) ([84a22a9](https://github.com/cloudquery/plugin-sdk/commit/84a22a97dba72365900b3e29becb87919449404f))

## [4.29.1](https://github.com/cloudquery/plugin-sdk/compare/v4.29.0...v4.29.1) (2024-02-01)


### Bug Fixes

* **deps:** Update golang.org/x/exp digest to 1b97071 ([#1505](https://github.com/cloudquery/plugin-sdk/issues/1505)) ([14d8545](https://github.com/cloudquery/plugin-sdk/commit/14d8545ac6c39d64f893c60d97dc19d2e144bdbc))
* **deps:** Update google.golang.org/genproto/googleapis/api digest to 1f4bbc5 ([#1506](https://github.com/cloudquery/plugin-sdk/issues/1506)) ([4021d65](https://github.com/cloudquery/plugin-sdk/commit/4021d65d966363f5efc37c16626c81f1e4b2f435))
* **deps:** Update google.golang.org/genproto/googleapis/rpc digest to 1f4bbc5 ([#1507](https://github.com/cloudquery/plugin-sdk/issues/1507)) ([b1316a8](https://github.com/cloudquery/plugin-sdk/commit/b1316a8423902b454505bc67f0582df9282ae0c1))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.7.2 ([#1501](https://github.com/cloudquery/plugin-sdk/issues/1501)) ([f5ecd8e](https://github.com/cloudquery/plugin-sdk/commit/f5ecd8e65e00a44b85cad183277d6cf448b721d2))
* **deps:** Update module github.com/kataras/iris/v12 to v12.2.10 ([#1508](https://github.com/cloudquery/plugin-sdk/issues/1508)) ([611982b](https://github.com/cloudquery/plugin-sdk/commit/611982b154ddd56c4722c809422dc394b1be2bef))
* **deps:** Update module github.com/klauspost/compress to v1.17.5 ([#1509](https://github.com/cloudquery/plugin-sdk/issues/1509)) ([e8d3c6b](https://github.com/cloudquery/plugin-sdk/commit/e8d3c6b2f4b518d05d5bf2f5b7a8415a064e79e0))
* **deps:** Update module github.com/pierrec/lz4/v4 to v4.1.21 ([#1510](https://github.com/cloudquery/plugin-sdk/issues/1510)) ([8af0e4e](https://github.com/cloudquery/plugin-sdk/commit/8af0e4e47fcebb0ef888ecdc364a1df1467418d0))
* **deps:** Update module github.com/tdewolff/minify/v2 to v2.20.16 ([#1511](https://github.com/cloudquery/plugin-sdk/issues/1511)) ([b1433cc](https://github.com/cloudquery/plugin-sdk/commit/b1433cc85889209d18c4c264a78b15d7bfd5c1dc))
* **deps:** Update module github.com/tdewolff/parse/v2 to v2.7.11 ([#1512](https://github.com/cloudquery/plugin-sdk/issues/1512)) ([401fa4a](https://github.com/cloudquery/plugin-sdk/commit/401fa4a27048f61cb2cb659e8340866466f9acf3))
* Handle PrimaryKeyComponents in packaging ([#1503](https://github.com/cloudquery/plugin-sdk/issues/1503)) ([8c8fdc9](https://github.com/cloudquery/plugin-sdk/commit/8c8fdc918569a04dbfb779f1134d273ffc1d9b1e))

## [4.29.0](https://github.com/cloudquery/plugin-sdk/compare/v4.28.0...v4.29.0) (2024-01-31)


### Features

* Introduce `PrimaryKeyComponent` ([#1491](https://github.com/cloudquery/plugin-sdk/issues/1491)) ([ae4a26e](https://github.com/cloudquery/plugin-sdk/commit/ae4a26e627f0d9d4df86eb93fee031753044f682))


### Bug Fixes

* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.7.1 ([#1499](https://github.com/cloudquery/plugin-sdk/issues/1499)) ([165be4d](https://github.com/cloudquery/plugin-sdk/commit/165be4dd7d22019c41546940f0b4913a2536f834))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.16.7 ([#1500](https://github.com/cloudquery/plugin-sdk/issues/1500)) ([2b98dab](https://github.com/cloudquery/plugin-sdk/commit/2b98daba1af1b26bd917f419a788c3a54113dd41))
* Remove access to parent tests in test suite ([#1497](https://github.com/cloudquery/plugin-sdk/issues/1497)) ([63e95e7](https://github.com/cloudquery/plugin-sdk/commit/63e95e7b36cfc9c277e03c4fc939868e7a377da6))

## [4.28.0](https://github.com/cloudquery/plugin-sdk/compare/v4.27.2...v4.28.0) (2024-01-30)


### Features

* Package JSON Schema with plugins ([#1494](https://github.com/cloudquery/plugin-sdk/issues/1494)) ([790e240](https://github.com/cloudquery/plugin-sdk/commit/790e240c3ff90e756881114037f1857b392934f0))

## [4.27.2](https://github.com/cloudquery/plugin-sdk/compare/v4.27.1...v4.27.2) (2024-01-29)


### Bug Fixes

* Better build overrides ([#1492](https://github.com/cloudquery/plugin-sdk/issues/1492)) ([ca5afc1](https://github.com/cloudquery/plugin-sdk/commit/ca5afc1aca7c2fa994ddb0e371d8b507435501e1))
* When `_cq_id` SyncMigrateMessage not sent ([#1489](https://github.com/cloudquery/plugin-sdk/issues/1489)) ([d177320](https://github.com/cloudquery/plugin-sdk/commit/d177320d62fb104a24381b1ba6cb0bfd9864c723))

## [4.27.1](https://github.com/cloudquery/plugin-sdk/compare/v4.27.0...v4.27.1) (2024-01-23)


### Bug Fixes

* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.16.6 ([#1485](https://github.com/cloudquery/plugin-sdk/issues/1485)) ([6de5f88](https://github.com/cloudquery/plugin-sdk/commit/6de5f886e052c592a7f81ec9a952df6f8d5ef641))

## [4.27.0](https://github.com/cloudquery/plugin-sdk/compare/v4.26.0...v4.27.0) (2024-01-23)


### Features

* Add Sync Run API Token Type ([#1473](https://github.com/cloudquery/plugin-sdk/issues/1473)) ([c776750](https://github.com/cloudquery/plugin-sdk/commit/c7767505318c98c7c4b11dea8796df52537c53df))


### Bug Fixes

* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.7.0 ([#1483](https://github.com/cloudquery/plugin-sdk/issues/1483)) ([01961cb](https://github.com/cloudquery/plugin-sdk/commit/01961cb11ef6e590d14dcc145ec94c9a1767d76d))

## [4.26.0](https://github.com/cloudquery/plugin-sdk/compare/v4.25.2...v4.26.0) (2024-01-23)


### Features

* Expose Migration test for new special case for moving to `_cq_id` as only primary key ([#1480](https://github.com/cloudquery/plugin-sdk/issues/1480)) ([321e355](https://github.com/cloudquery/plugin-sdk/commit/321e35557d656c92b88cd5ca1a3ae0cfda8c3bfb))
* Expose special migration paths ([#1470](https://github.com/cloudquery/plugin-sdk/issues/1470)) ([d70eaff](https://github.com/cloudquery/plugin-sdk/commit/d70eafff3324f81be67e1cbe360814d68892886c))
* Make UUID in `testdata` always deterministic like all other columns ([#1479](https://github.com/cloudquery/plugin-sdk/issues/1479)) ([78027f0](https://github.com/cloudquery/plugin-sdk/commit/78027f0eae3a66b6d13fd3f86af4c28f142329fb))
* Support CQ ID on the source only ([#1461](https://github.com/cloudquery/plugin-sdk/issues/1461)) ([f583cea](https://github.com/cloudquery/plugin-sdk/commit/f583ceabc11cab9be2371b90ce1c7f44e17f8ca4))


### Bug Fixes

* **deps:** Update github.com/apache/arrow/go/v15 digest to 7e703aa ([#1467](https://github.com/cloudquery/plugin-sdk/issues/1467)) ([7645b7a](https://github.com/cloudquery/plugin-sdk/commit/7645b7a3c9d12544a3609c181e43353fbffe777d))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.6.5 ([#1471](https://github.com/cloudquery/plugin-sdk/issues/1471)) ([acb1ac7](https://github.com/cloudquery/plugin-sdk/commit/acb1ac7dcda5b99e7f8697fc12f4b6239dac7d83))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.16.5 ([#1474](https://github.com/cloudquery/plugin-sdk/issues/1474)) ([aa35ce7](https://github.com/cloudquery/plugin-sdk/commit/aa35ce79e4821aad32893d36da95a6bc9b229f79))
* Eternal recursion in `scalar.MonthInterval` ([#1477](https://github.com/cloudquery/plugin-sdk/issues/1477)) ([78219a6](https://github.com/cloudquery/plugin-sdk/commit/78219a68e536c4cab7dcf1ab93c40f4691101487))
* Handle unrelated licenses ([#1472](https://github.com/cloudquery/plugin-sdk/issues/1472)) ([4936425](https://github.com/cloudquery/plugin-sdk/commit/49364255b2fe4c2a2982d0e616cd60c16f9f54a8))
* Verify `nil` value consistently ([#1478](https://github.com/cloudquery/plugin-sdk/issues/1478)) ([31085d2](https://github.com/cloudquery/plugin-sdk/commit/31085d23fc3940ba8ee422e25b33978695df5272))

## [4.25.2](https://github.com/cloudquery/plugin-sdk/compare/v4.25.1...v4.25.2) (2024-01-12)


### Bug Fixes

* **deps:** Update github.com/apache/arrow/go/v15 digest to 6d44906 ([#1462](https://github.com/cloudquery/plugin-sdk/issues/1462)) ([45533de](https://github.com/cloudquery/plugin-sdk/commit/45533de27eddbb0207f2d9ac23cfbb088827cb36))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.16.4 ([#1464](https://github.com/cloudquery/plugin-sdk/issues/1464)) ([098749c](https://github.com/cloudquery/plugin-sdk/commit/098749c2413cd7f6f72efa9caef423e19e3f1189))

## [4.25.1](https://github.com/cloudquery/plugin-sdk/compare/v4.25.0...v4.25.1) (2024-01-05)


### Bug Fixes

* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.6.4 ([#1459](https://github.com/cloudquery/plugin-sdk/issues/1459)) ([5ec8f8d](https://github.com/cloudquery/plugin-sdk/commit/5ec8f8d9c2f35f937ebe03007bf321a51a368ab1))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.16.2 ([#1456](https://github.com/cloudquery/plugin-sdk/issues/1456)) ([341d770](https://github.com/cloudquery/plugin-sdk/commit/341d770669f8cc4db30edb9c40e44af49eb0ecfe))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.16.3 ([#1458](https://github.com/cloudquery/plugin-sdk/issues/1458)) ([4dd2130](https://github.com/cloudquery/plugin-sdk/commit/4dd2130e8129ea15a5e06eb5a619bcebd6770c44))

## [4.25.0](https://github.com/cloudquery/plugin-sdk/compare/v4.24.1...v4.25.0) (2024-01-02)


### Features

* Support multiple and/or specific plugin licenses ([#1451](https://github.com/cloudquery/plugin-sdk/issues/1451)) ([993e352](https://github.com/cloudquery/plugin-sdk/commit/993e352dd2abbdfaa1ff5d6a3cc48c38457fa7f8))


### Bug Fixes

* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.16.1 ([#1454](https://github.com/cloudquery/plugin-sdk/issues/1454)) ([dc4afb6](https://github.com/cloudquery/plugin-sdk/commit/dc4afb6994c673623ff10539ba04fca34b9a02d8))

## [4.24.1](https://github.com/cloudquery/plugin-sdk/compare/v4.24.0...v4.24.1) (2024-01-01)


### Bug Fixes

* **deps:** Update github.com/apache/arrow/go/v15 digest to 7c3480e ([#1443](https://github.com/cloudquery/plugin-sdk/issues/1443)) ([bc8644f](https://github.com/cloudquery/plugin-sdk/commit/bc8644f40c11ab9d39d14e90d2cdb07d7b89898d))
* **deps:** Update github.com/gomarkdown/markdown digest to 1d6d208 ([#1445](https://github.com/cloudquery/plugin-sdk/issues/1445)) ([9a29286](https://github.com/cloudquery/plugin-sdk/commit/9a2928606c7f627ab7b5c74efdfc4b2d2484d98f))
* **deps:** Update golang.org/x/exp digest to 02704c9 ([#1446](https://github.com/cloudquery/plugin-sdk/issues/1446)) ([496d59d](https://github.com/cloudquery/plugin-sdk/commit/496d59d34ef540d8a2b4a683f838938e3de3b239))
* **deps:** Update google.golang.org/genproto/googleapis/api digest to 995d672 ([#1447](https://github.com/cloudquery/plugin-sdk/issues/1447)) ([21771e7](https://github.com/cloudquery/plugin-sdk/commit/21771e759b0dba180c679e7221bad62a26466ce1))
* **deps:** Update google.golang.org/genproto/googleapis/rpc digest to 995d672 ([#1448](https://github.com/cloudquery/plugin-sdk/issues/1448)) ([2135e11](https://github.com/cloudquery/plugin-sdk/commit/2135e1105800bd65a57806cf8ed6c1a0283e0188))
* **deps:** Update module github.com/klauspost/compress to v1.17.4 ([#1450](https://github.com/cloudquery/plugin-sdk/issues/1450)) ([04323d7](https://github.com/cloudquery/plugin-sdk/commit/04323d7f599f10693b072322eb6e6ec1714fa835))

## [4.24.0](https://github.com/cloudquery/plugin-sdk/compare/v4.23.0...v4.24.0) (2023-12-29)


### Features

* Offline licensing support ([1fdf892](https://github.com/cloudquery/plugin-sdk/commit/1fdf892b8b4e4a90da4e69a463af0dd7d8b6a420))


### Bug Fixes

* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.15.0 ([#1438](https://github.com/cloudquery/plugin-sdk/issues/1438)) ([e0c2a4b](https://github.com/cloudquery/plugin-sdk/commit/e0c2a4bbf6248294ae62e47e129a65ed8dc01277))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.16.0 ([#1440](https://github.com/cloudquery/plugin-sdk/issues/1440)) ([d2a5850](https://github.com/cloudquery/plugin-sdk/commit/d2a5850e126368fd3e03f0d993383ac0e355c8bc))

## [4.23.0](https://github.com/cloudquery/plugin-sdk/compare/v4.22.0...v4.23.0) (2023-12-27)


### Features

* Introduce a per resource rate limit in addition to a global resource rate limit ([2918402](https://github.com/cloudquery/plugin-sdk/commit/29184024a39264669b1f2e70daf2149361ef9c7f))
* Set default rate limit of `5` for `SingleResourceMaxConcurrency` and `SingleNestedTableMaxConcurrency` ([2918402](https://github.com/cloudquery/plugin-sdk/commit/29184024a39264669b1f2e70daf2149361ef9c7f))

## [4.22.0](https://github.com/cloudquery/plugin-sdk/compare/v4.21.3...v4.22.0) (2023-12-26)


### Features

* Expose otel headers and url_path as flags ([#1430](https://github.com/cloudquery/plugin-sdk/issues/1430)) ([3541726](https://github.com/cloudquery/plugin-sdk/commit/3541726fb27d437d9b059fba40396690d758d60a))
* Faker should preserve previous values, if set ([#1429](https://github.com/cloudquery/plugin-sdk/issues/1429)) ([e44f185](https://github.com/cloudquery/plugin-sdk/commit/e44f1857856c5dafa2e7cb369cb9365d08697cb7))


### Bug Fixes

* **deps:** Update github.com/apache/arrow/go/v15 digest to ec41209 ([#1431](https://github.com/cloudquery/plugin-sdk/issues/1431)) ([b50e9ac](https://github.com/cloudquery/plugin-sdk/commit/b50e9ac396de183d6fd7b062b27aedaa047fed04))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.6.3 ([#1427](https://github.com/cloudquery/plugin-sdk/issues/1427)) ([7d8a9d9](https://github.com/cloudquery/plugin-sdk/commit/7d8a9d9d3c3cb28e71ed3c0680f180c8162fa355))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.14.6 ([#1425](https://github.com/cloudquery/plugin-sdk/issues/1425)) ([870936f](https://github.com/cloudquery/plugin-sdk/commit/870936f65c9b497b29125a3e1dec9863936864fb))

## [4.21.3](https://github.com/cloudquery/plugin-sdk/compare/v4.21.2...v4.21.3) (2023-12-19)


### Bug Fixes

* **deps:** Update `github.com/apache/arrow/go` to `v15` ([#1424](https://github.com/cloudquery/plugin-sdk/issues/1424)) ([64db12d](https://github.com/cloudquery/plugin-sdk/commit/64db12d30dfb09b434c1399e81a8943c3fa2e046))
* **deps:** Update module golang.org/x/crypto to v0.17.0 [SECURITY] ([#1422](https://github.com/cloudquery/plugin-sdk/issues/1422)) ([975adba](https://github.com/cloudquery/plugin-sdk/commit/975adba3b3ea5b3249b14a1e671631ce2604563e))

## [4.21.2](https://github.com/cloudquery/plugin-sdk/compare/v4.21.1...v4.21.2) (2023-12-18)


### Bug Fixes

* **tests:** Find empty columns for JSON types ([#1418](https://github.com/cloudquery/plugin-sdk/issues/1418)) ([027273c](https://github.com/cloudquery/plugin-sdk/commit/027273c1be6b4163406a4127139a9870c1eafec8))

## [4.21.1](https://github.com/cloudquery/plugin-sdk/compare/v4.21.0...v4.21.1) (2023-12-14)


### Bug Fixes

* Update usage limit message ([#1415](https://github.com/cloudquery/plugin-sdk/issues/1415)) ([98438ff](https://github.com/cloudquery/plugin-sdk/commit/98438ffc22ebdec58a6637b02c0fef2a45e19dd5))

## [4.21.0](https://github.com/cloudquery/plugin-sdk/compare/v4.20.0...v4.21.0) (2023-12-13)


### Features

* Individual Table and Client rate limit ([#1411](https://github.com/cloudquery/plugin-sdk/issues/1411)) ([4d13b18](https://github.com/cloudquery/plugin-sdk/commit/4d13b18b5ef33d3159155289703dce67e1ad750c))


### Bug Fixes

* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.6.2 ([#1413](https://github.com/cloudquery/plugin-sdk/issues/1413)) ([f5a0d47](https://github.com/cloudquery/plugin-sdk/commit/f5a0d47b0ea5628166eb9138a1b9f67241598344))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.14.4 ([#1408](https://github.com/cloudquery/plugin-sdk/issues/1408)) ([7544967](https://github.com/cloudquery/plugin-sdk/commit/754496784a2c182e1765aa7a5ef832a337e6a7f8))

## [4.20.0](https://github.com/cloudquery/plugin-sdk/compare/v4.19.1...v4.20.0) (2023-12-07)


### Features

* Add `GetPaidTables()` and `HasPaidTables()` to `schema.Tables` ([#1403](https://github.com/cloudquery/plugin-sdk/issues/1403)) ([b355fa0](https://github.com/cloudquery/plugin-sdk/commit/b355fa07dd8a1265b93c6f3b4f6d17f663a93912))
* Include `is_paid` field when creating tables json during package ([#1405](https://github.com/cloudquery/plugin-sdk/issues/1405)) ([455a1e3](https://github.com/cloudquery/plugin-sdk/commit/455a1e3ebf0eea79bbd11c0f31315775d8609b2b))


### Bug Fixes

* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.14.3 ([#1406](https://github.com/cloudquery/plugin-sdk/issues/1406)) ([7833342](https://github.com/cloudquery/plugin-sdk/commit/783334242e12d4d3fe78ddfd5acd11ecd8220fff))

## [4.19.1](https://github.com/cloudquery/plugin-sdk/compare/v4.19.0...v4.19.1) (2023-12-04)


### Bug Fixes

* **deps:** Update github.com/gomarkdown/markdown digest to a660076 ([#1392](https://github.com/cloudquery/plugin-sdk/issues/1392)) ([8a1c31a](https://github.com/cloudquery/plugin-sdk/commit/8a1c31a609d98319c6cef0a01c37f208968f3bba))
* **deps:** Update golang.org/x/exp digest to 6522937 ([#1394](https://github.com/cloudquery/plugin-sdk/issues/1394)) ([5b4f9ac](https://github.com/cloudquery/plugin-sdk/commit/5b4f9acb3de89cef6e0dd999c411a60eae8b68fe))
* **deps:** Update google.golang.org/genproto/googleapis/api digest to 3a041ad ([#1396](https://github.com/cloudquery/plugin-sdk/issues/1396)) ([403be86](https://github.com/cloudquery/plugin-sdk/commit/403be86e1b76ec887957fd1be791c5fa6b3074e7))
* **deps:** Update google.golang.org/genproto/googleapis/rpc digest to 3a041ad ([#1397](https://github.com/cloudquery/plugin-sdk/issues/1397)) ([89a063f](https://github.com/cloudquery/plugin-sdk/commit/89a063f63a3d4915ba50f986b8a80b94645ca26b))
* **deps:** Update module github.com/chenzhuoyu/iasm to v0.9.1 ([#1398](https://github.com/cloudquery/plugin-sdk/issues/1398)) ([a0e516a](https://github.com/cloudquery/plugin-sdk/commit/a0e516a563d663b95f88d3e726e6728c8e17b45b))
* **deps:** Update module github.com/gorilla/css to v1.0.1 ([#1399](https://github.com/cloudquery/plugin-sdk/issues/1399)) ([8bbeafa](https://github.com/cloudquery/plugin-sdk/commit/8bbeafab587b61553100741d1547f17c994f37e4))
* Fail early on usage client init if token is not set ([#1401](https://github.com/cloudquery/plugin-sdk/issues/1401)) ([dce2b0d](https://github.com/cloudquery/plugin-sdk/commit/dce2b0db513aa8ea1755b6846022004c082db49d))

## [4.19.0](https://github.com/cloudquery/plugin-sdk/compare/v4.18.3...v4.19.0) (2023-11-30)


### Features

* Improved tracing ([#1387](https://github.com/cloudquery/plugin-sdk/issues/1387)) ([68cfc32](https://github.com/cloudquery/plugin-sdk/commit/68cfc322c6e35525833bc79cfe2c1c6c8ef2fe71))


### Bug Fixes

* Cleanup batch writers ([#1386](https://github.com/cloudquery/plugin-sdk/issues/1386)) ([cde7462](https://github.com/cloudquery/plugin-sdk/commit/cde7462e3f33897136b9120e73351dc253449e8a))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.14.1 ([#1380](https://github.com/cloudquery/plugin-sdk/issues/1380)) ([e5451c6](https://github.com/cloudquery/plugin-sdk/commit/e5451c636c54b987daac2cdeaef97a826155671f))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.14.2 ([#1382](https://github.com/cloudquery/plugin-sdk/issues/1382)) ([8acdb72](https://github.com/cloudquery/plugin-sdk/commit/8acdb720ce4d16e132868ff3851b10395d596ba2))

## [4.18.3](https://github.com/cloudquery/plugin-sdk/compare/v4.18.2...v4.18.3) (2023-11-17)


### Bug Fixes

* Retrieve team for api key ([#1372](https://github.com/cloudquery/plugin-sdk/issues/1372)) ([940d87f](https://github.com/cloudquery/plugin-sdk/commit/940d87f7cc71d8e2367c3f751dcf1d081e2b8126))

## [4.18.2](https://github.com/cloudquery/plugin-sdk/compare/v4.18.1...v4.18.2) (2023-11-16)


### Bug Fixes

* Batching for mixedbatchwriter ([#1374](https://github.com/cloudquery/plugin-sdk/issues/1374)) ([ca435cf](https://github.com/cloudquery/plugin-sdk/commit/ca435cfe4a42271dadc9ea0a119a4515804efebb))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.4.6 ([#1368](https://github.com/cloudquery/plugin-sdk/issues/1368)) ([ea05199](https://github.com/cloudquery/plugin-sdk/commit/ea0519920ab1fadced3a27320a7f50a20e0bf080))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.5.1 ([#1370](https://github.com/cloudquery/plugin-sdk/issues/1370)) ([309b1cb](https://github.com/cloudquery/plugin-sdk/commit/309b1cb8267c867d6be827f36dd63fdb138485ae))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.6.0 ([#1373](https://github.com/cloudquery/plugin-sdk/issues/1373)) ([63fc4bb](https://github.com/cloudquery/plugin-sdk/commit/63fc4bbb605bf92a79def791a2b7e5d3fd09f42a))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.14.0 ([#1371](https://github.com/cloudquery/plugin-sdk/issues/1371)) ([8ec6a34](https://github.com/cloudquery/plugin-sdk/commit/8ec6a3422dc387662a5028b81a483bf8b2e8d1dc))

## [4.18.1](https://github.com/cloudquery/plugin-sdk/compare/v4.18.0...v4.18.1) (2023-11-14)


### Bug Fixes

* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.4.4 ([#1364](https://github.com/cloudquery/plugin-sdk/issues/1364)) ([d5a5760](https://github.com/cloudquery/plugin-sdk/commit/d5a5760c7f876fbb50db5fe09cfcd03bb42fdb04))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.4.5 ([#1365](https://github.com/cloudquery/plugin-sdk/issues/1365)) ([2ec138f](https://github.com/cloudquery/plugin-sdk/commit/2ec138f178100f96c36cc0a07c223a676a423a58))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.13.5 ([#1362](https://github.com/cloudquery/plugin-sdk/issues/1362)) ([6663a64](https://github.com/cloudquery/plugin-sdk/commit/6663a64ec9b0acbb3d8fea4f2585d780e8af651d))
* Mark relations as paid as well ([#1366](https://github.com/cloudquery/plugin-sdk/issues/1366)) ([ca833eb](https://github.com/cloudquery/plugin-sdk/commit/ca833eb5c83aa580d4fe2568a3dfa079b3a3614e))

## [4.18.0](https://github.com/cloudquery/plugin-sdk/compare/v4.17.2...v4.18.0) (2023-11-09)


### Features

* **package:** Check for Version variable ([#1359](https://github.com/cloudquery/plugin-sdk/issues/1359)) ([2f1aff8](https://github.com/cloudquery/plugin-sdk/commit/2f1aff831be92e20dba91a73b17e8ed4a224effb))


### Bug Fixes

* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.4.3 ([#1357](https://github.com/cloudquery/plugin-sdk/issues/1357)) ([f5cd387](https://github.com/cloudquery/plugin-sdk/commit/f5cd3870271da3593ec82ffdfba5ad835bf15d65))

## [4.17.2](https://github.com/cloudquery/plugin-sdk/compare/v4.17.1...v4.17.2) (2023-11-02)


### Bug Fixes

* **deps:** Update github.com/apache/arrow/go/v14 digest to c49e242 ([#1343](https://github.com/cloudquery/plugin-sdk/issues/1343)) ([8f6362e](https://github.com/cloudquery/plugin-sdk/commit/8f6362e8f2153c597bed2577729efa8cd7924d1b))
* **deps:** Update golang.org/x/xerrors digest to 104605a ([#1345](https://github.com/cloudquery/plugin-sdk/issues/1345)) ([5b3e9c6](https://github.com/cloudquery/plugin-sdk/commit/5b3e9c61634e9169895facb37deb2a403f833792))
* **deps:** Update google.golang.org/genproto digest to d783a09 ([#1346](https://github.com/cloudquery/plugin-sdk/issues/1346)) ([2af9c70](https://github.com/cloudquery/plugin-sdk/commit/2af9c70fe1bf54f3654d06b5028520e5ade9b2df))
* **deps:** Update google.golang.org/genproto/googleapis/api digest to d783a09 ([#1347](https://github.com/cloudquery/plugin-sdk/issues/1347)) ([6f43900](https://github.com/cloudquery/plugin-sdk/commit/6f43900227fe95b58c278cc2b86ca2bf909fcf33))
* **deps:** Update google.golang.org/genproto/googleapis/rpc digest to d783a09 ([#1348](https://github.com/cloudquery/plugin-sdk/issues/1348)) ([bdf7a32](https://github.com/cloudquery/plugin-sdk/commit/bdf7a321af9d748bb19ea08182c4333d43ed6deb))
* **deps:** Update module github.com/andybalholm/brotli to v1.0.6 ([#1349](https://github.com/cloudquery/plugin-sdk/issues/1349)) ([2e79c6f](https://github.com/cloudquery/plugin-sdk/commit/2e79c6f6d37d3f6c8496b4de35232f34639151f5))
* **deps:** Update module github.com/bytedance/sonic to v1.10.2 ([#1350](https://github.com/cloudquery/plugin-sdk/issues/1350)) ([147b381](https://github.com/cloudquery/plugin-sdk/commit/147b381f2d4a2d48d3799530463c2b41ed79e5f3))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.13.2 ([#1351](https://github.com/cloudquery/plugin-sdk/issues/1351)) ([d3d34e5](https://github.com/cloudquery/plugin-sdk/commit/d3d34e55c95d95ab95753abf3a4a9704de349f8c))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.13.3 ([#1352](https://github.com/cloudquery/plugin-sdk/issues/1352)) ([31137ad](https://github.com/cloudquery/plugin-sdk/commit/31137ad67036202d901fc1e84994e8ed050bd458))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.13.4 ([#1353](https://github.com/cloudquery/plugin-sdk/issues/1353)) ([f5c1bbe](https://github.com/cloudquery/plugin-sdk/commit/f5c1bbe4ae30029accd5698be1384d414baff4c8))
* Return clearer error when team name empty or not set ([#1354](https://github.com/cloudquery/plugin-sdk/issues/1354)) ([e82c69e](https://github.com/cloudquery/plugin-sdk/commit/e82c69ec37050432bc551b93c9526eae9716e0b4))

## [4.17.1](https://github.com/cloudquery/plugin-sdk/compare/v4.17.0...v4.17.1) (2023-10-31)


### Bug Fixes

* Fix nil pointer dereference when remaining rows not set ([#1339](https://github.com/cloudquery/plugin-sdk/issues/1339)) ([36a9d35](https://github.com/cloudquery/plugin-sdk/commit/36a9d3534c2613df926c0ddd0460f3b548336b5c))

## [4.17.0](https://github.com/cloudquery/plugin-sdk/compare/v4.16.1...v4.17.0) (2023-10-30)


### Features

* Add IsPaid flag to table definition ([#1327](https://github.com/cloudquery/plugin-sdk/issues/1327)) ([ffd14bf](https://github.com/cloudquery/plugin-sdk/commit/ffd14bf398fb8fd6831da34e3b99be0eb1a618ab))
* Add OnBeforeSend hook ([#1325](https://github.com/cloudquery/plugin-sdk/issues/1325)) ([023ebbc](https://github.com/cloudquery/plugin-sdk/commit/023ebbc522959e1826a6bf2480395cb38c27aea0))
* Adding a batch updater to allow usage updates to be batched ([#1326](https://github.com/cloudquery/plugin-sdk/issues/1326)) ([0301ed7](https://github.com/cloudquery/plugin-sdk/commit/0301ed75928a6e8bc50984cb5ec29880396cbc4f))
* Adding quota monitoring for premium plugins ([#1333](https://github.com/cloudquery/plugin-sdk/issues/1333)) ([b7a2ca5](https://github.com/cloudquery/plugin-sdk/commit/b7a2ca547a3d819eff7283d8a3afa312335617a9))
* Allow sync to be cancelled when in progress ([#1334](https://github.com/cloudquery/plugin-sdk/issues/1334)) ([6d7be0b](https://github.com/cloudquery/plugin-sdk/commit/6d7be0bd9e25710d0e92407f34fe38a11c3f8dad))


### Bug Fixes

* **deps:** Update github.com/cloudquery/arrow/go/v14 digest to 50d3871 ([#1337](https://github.com/cloudquery/plugin-sdk/issues/1337)) ([f15a89d](https://github.com/cloudquery/plugin-sdk/commit/f15a89d64e604642455951895bf3db3e04ae4afe))
* **deps:** Update github.com/cloudquery/arrow/go/v14 digest to f46436f ([#1329](https://github.com/cloudquery/plugin-sdk/issues/1329)) ([ee24384](https://github.com/cloudquery/plugin-sdk/commit/ee243848baa2e6c6e5737233c926c44897de0ef0))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.4.2 ([#1335](https://github.com/cloudquery/plugin-sdk/issues/1335)) ([2ecd2a1](https://github.com/cloudquery/plugin-sdk/commit/2ecd2a1f47ac6ad3d529da0aded01fcdd8f8cb18))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.13.0 ([#1332](https://github.com/cloudquery/plugin-sdk/issues/1332)) ([5553f85](https://github.com/cloudquery/plugin-sdk/commit/5553f8556a7dda0be9425c70f9694140c7afb103))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.13.1 ([#1336](https://github.com/cloudquery/plugin-sdk/issues/1336)) ([b782ee7](https://github.com/cloudquery/plugin-sdk/commit/b782ee714e87cac8901eac4f032e51fd4362a997))
* **deps:** Update module google.golang.org/grpc to v1.58.3 [SECURITY] ([#1331](https://github.com/cloudquery/plugin-sdk/issues/1331)) ([43f60c2](https://github.com/cloudquery/plugin-sdk/commit/43f60c2d229dc4947cb4a020bd6a54b9b4d8325e))

## [4.16.1](https://github.com/cloudquery/plugin-sdk/compare/v4.16.0...v4.16.1) (2023-10-19)


### Bug Fixes

* **package:** Only return one level down of relations when writing `tables.json` ([#1321](https://github.com/cloudquery/plugin-sdk/issues/1321)) ([3d4ebe0](https://github.com/cloudquery/plugin-sdk/commit/3d4ebe0098ba4e458d88e092e6240ee848c38c0a))

## [4.16.0](https://github.com/cloudquery/plugin-sdk/compare/v4.15.3...v4.16.0) (2023-10-19)


### Features

* Support publishing plugins with team and kind metadata set ([#1313](https://github.com/cloudquery/plugin-sdk/issues/1313)) ([933698d](https://github.com/cloudquery/plugin-sdk/commit/933698dca6da13c2a8e428f7758e8a9911326095))

## [4.15.3](https://github.com/cloudquery/plugin-sdk/compare/v4.15.2...v4.15.3) (2023-10-18)


### Bug Fixes

* Set all fields in `DeleteRecord` message ([#1316](https://github.com/cloudquery/plugin-sdk/issues/1316)) ([ad9d109](https://github.com/cloudquery/plugin-sdk/commit/ad9d10936f0362542af280fd517377d30010033b))

## [4.15.2](https://github.com/cloudquery/plugin-sdk/compare/v4.15.1...v4.15.2) (2023-10-18)


### Bug Fixes

* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.4.1 ([#1312](https://github.com/cloudquery/plugin-sdk/issues/1312)) ([0c75527](https://github.com/cloudquery/plugin-sdk/commit/0c7552704d5ca751638ad3119fc51dc882a0caf5))

## [4.15.1](https://github.com/cloudquery/plugin-sdk/compare/v4.15.0...v4.15.1) (2023-10-18)


### Bug Fixes

* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.2.9 ([#1306](https://github.com/cloudquery/plugin-sdk/issues/1306)) ([e8ebf8d](https://github.com/cloudquery/plugin-sdk/commit/e8ebf8d6037a29f6506f80db46678690c8718e7e))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.3.0 ([#1308](https://github.com/cloudquery/plugin-sdk/issues/1308)) ([15d7129](https://github.com/cloudquery/plugin-sdk/commit/15d7129baa31d6fe36d7bef6f0cb6467b7016dae))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.4.0 ([#1309](https://github.com/cloudquery/plugin-sdk/issues/1309)) ([4e90152](https://github.com/cloudquery/plugin-sdk/commit/4e9015201685061dcea2530703dd9bb757ee7763))
* Make static linking conditional only for Linux ([#1310](https://github.com/cloudquery/plugin-sdk/issues/1310)) ([35fa449](https://github.com/cloudquery/plugin-sdk/commit/35fa449c8877395cb5d12d63fbe505c983df78c3))

## [4.15.0](https://github.com/cloudquery/plugin-sdk/compare/v4.14.1...v4.15.0) (2023-10-17)


### Features

* Add JSON schema for `configtype.Duration` ([#1303](https://github.com/cloudquery/plugin-sdk/issues/1303)) ([5e1598b](https://github.com/cloudquery/plugin-sdk/commit/5e1598b48967d5a36c1bde74f4c811504a1009e1))

## [4.14.1](https://github.com/cloudquery/plugin-sdk/compare/v4.14.0...v4.14.1) (2023-10-16)


### Bug Fixes

* Enable Skipping of DeleteRecord tests ([#1299](https://github.com/cloudquery/plugin-sdk/issues/1299)) ([5dd5739](https://github.com/cloudquery/plugin-sdk/commit/5dd573908f69e6d35b3e19c2ed7a5b60be583807))

## [4.14.0](https://github.com/cloudquery/plugin-sdk/compare/v4.13.0...v4.14.0) (2023-10-16)


### Features

* Support DeleteRecord in all writers ([#1295](https://github.com/cloudquery/plugin-sdk/issues/1295)) ([5a02e27](https://github.com/cloudquery/plugin-sdk/commit/5a02e27525a2c225b55bd0e668be6038035630d5))

## [4.13.0](https://github.com/cloudquery/plugin-sdk/compare/v4.12.5...v4.13.0) (2023-10-16)


### Features

* Add support for conditional static linking of C lib to builds ([#1292](https://github.com/cloudquery/plugin-sdk/issues/1292)) ([7c27065](https://github.com/cloudquery/plugin-sdk/commit/7c27065c6ac9a4f84b8ea7dc7024f01677cc6357))
* Support Delete Record ([#1282](https://github.com/cloudquery/plugin-sdk/issues/1282)) ([1f0a603](https://github.com/cloudquery/plugin-sdk/commit/1f0a6039e61d64ee0530c6a195ee38ba183dad7f))


### Bug Fixes

* **deps:** Update github.com/cloudquery/arrow/go/v14 digest to dbcb149 ([#1291](https://github.com/cloudquery/plugin-sdk/issues/1291)) ([7c634dc](https://github.com/cloudquery/plugin-sdk/commit/7c634dc1e8e0ef6959a73922938ff8280d326682))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.12.3 ([#1289](https://github.com/cloudquery/plugin-sdk/issues/1289)) ([3e063bc](https://github.com/cloudquery/plugin-sdk/commit/3e063bc7eda88938d96ee94bc7ebdc062d4822f2))

## [4.12.5](https://github.com/cloudquery/plugin-sdk/compare/v4.12.4...v4.12.5) (2023-10-12)


### Bug Fixes

* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.2.7 ([#1285](https://github.com/cloudquery/plugin-sdk/issues/1285)) ([e27875e](https://github.com/cloudquery/plugin-sdk/commit/e27875ea0e9bc1bee07214f87cd689c67da2b04e))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.2.8 ([#1286](https://github.com/cloudquery/plugin-sdk/issues/1286)) ([9d9eb10](https://github.com/cloudquery/plugin-sdk/commit/9d9eb1007e43928de7994772c58e352acf43f7dd))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.12.2 ([#1287](https://github.com/cloudquery/plugin-sdk/issues/1287)) ([57e4795](https://github.com/cloudquery/plugin-sdk/commit/57e479507a9d4244d8a2f82731c192570ae4c6b7))
* **deps:** Update module golang.org/x/net to v0.17.0 [SECURITY] ([#1283](https://github.com/cloudquery/plugin-sdk/issues/1283)) ([4e5f9de](https://github.com/cloudquery/plugin-sdk/commit/4e5f9de50a76a29b44164a9072f179c3915b9fbb))

## [4.12.4](https://github.com/cloudquery/plugin-sdk/compare/v4.12.3...v4.12.4) (2023-10-10)


### Bug Fixes

* **deps:** Update github.com/cloudquery/arrow/go/v14 digest to d401686 ([#1277](https://github.com/cloudquery/plugin-sdk/issues/1277)) ([c94273b](https://github.com/cloudquery/plugin-sdk/commit/c94273b03bde133a1c684256ecbedc01dd730e38))
* **deps:** Update module github.com/cloudquery/cloudquery-api-go to v1.2.6 ([#1279](https://github.com/cloudquery/plugin-sdk/issues/1279)) ([d49f8dc](https://github.com/cloudquery/plugin-sdk/commit/d49f8dca4f61b4fd9e07cf970e97eb029d05282a))
* Only warn on validation err ([#1280](https://github.com/cloudquery/plugin-sdk/issues/1280)) ([299c1d3](https://github.com/cloudquery/plugin-sdk/commit/299c1d3c9a25497c724e24f7831c838b8951bb3e))

## [4.12.3](https://github.com/cloudquery/plugin-sdk/compare/v4.12.2...v4.12.3) (2023-10-05)


### Bug Fixes

* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.12.1 ([#1272](https://github.com/cloudquery/plugin-sdk/issues/1272)) ([7d7e15b](https://github.com/cloudquery/plugin-sdk/commit/7d7e15b3b712908ab0e56e9c4138154463cfe03e))

## [4.12.2](https://github.com/cloudquery/plugin-sdk/compare/v4.12.1...v4.12.2) (2023-10-05)


### Bug Fixes

* Serialize columns during package ([#1270](https://github.com/cloudquery/plugin-sdk/issues/1270)) ([cd5f79d](https://github.com/cloudquery/plugin-sdk/commit/cd5f79d15570415b49bd0eff00e1a46227ffa7f9))

## [4.12.1](https://github.com/cloudquery/plugin-sdk/compare/v4.12.0...v4.12.1) (2023-10-05)


### Bug Fixes

* Add `linux_arm64` to default build targets ([#1267](https://github.com/cloudquery/plugin-sdk/issues/1267)) ([a5f46d1](https://github.com/cloudquery/plugin-sdk/commit/a5f46d18672a434fffe94320751a28c90e7c7ec2))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.11.1 ([#1264](https://github.com/cloudquery/plugin-sdk/issues/1264)) ([7a390f0](https://github.com/cloudquery/plugin-sdk/commit/7a390f06842b0354d9359839b4129bc8efd4141d))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.12.0 ([#1268](https://github.com/cloudquery/plugin-sdk/issues/1268)) ([16669fe](https://github.com/cloudquery/plugin-sdk/commit/16669fe393926566acdca4404e0fdca089a9fe88))

## [4.12.0](https://github.com/cloudquery/plugin-sdk/compare/v4.11.1...v4.12.0) (2023-10-02)


### Features

* Add JSON schema to scheduler strategy ([#1254](https://github.com/cloudquery/plugin-sdk/issues/1254)) ([1cec01d](https://github.com/cloudquery/plugin-sdk/commit/1cec01de43faa4f6f44af58428cb95b269f97990))


### Bug Fixes

* **deps:** Update github.com/apache/arrow/go/v14 digest to 00efb06 ([#1257](https://github.com/cloudquery/plugin-sdk/issues/1257)) ([e56f6f8](https://github.com/cloudquery/plugin-sdk/commit/e56f6f82f34795f21aa1bad5fc3a62b85417fbf5))
* **deps:** Update github.com/cloudquery/arrow/go/v14 digest to 7ded38b ([#1263](https://github.com/cloudquery/plugin-sdk/issues/1263)) ([332c255](https://github.com/cloudquery/plugin-sdk/commit/332c2555cc7e13f05612a274e63fe59af4c5ba98))
* **deps:** Update google.golang.org/genproto digest to e6e6cda ([#1258](https://github.com/cloudquery/plugin-sdk/issues/1258)) ([1b75050](https://github.com/cloudquery/plugin-sdk/commit/1b75050c5fafa8ea27a3e4841dbd2ce9001d801e))
* **deps:** Update google.golang.org/genproto/googleapis/api digest to e6e6cda ([#1259](https://github.com/cloudquery/plugin-sdk/issues/1259)) ([eb6a97d](https://github.com/cloudquery/plugin-sdk/commit/eb6a97dfc702b4cc779aff42152d21de8270de7b))
* **deps:** Update google.golang.org/genproto/googleapis/rpc digest to e6e6cda ([#1260](https://github.com/cloudquery/plugin-sdk/issues/1260)) ([49940fd](https://github.com/cloudquery/plugin-sdk/commit/49940fd94bb4ab605ea511e957e02316e31e046c))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.11.0 ([#1252](https://github.com/cloudquery/plugin-sdk/issues/1252)) ([41a6561](https://github.com/cloudquery/plugin-sdk/commit/41a6561f2ab0f048c1f333d5a3de558014f58f5f))
* **deps:** Update module github.com/getsentry/sentry-go to v0.24.1 ([#1262](https://github.com/cloudquery/plugin-sdk/issues/1262)) ([be03068](https://github.com/cloudquery/plugin-sdk/commit/be030689c413afa341a4b7e0644c4d28be6c9640))
* **deps:** Update module github.com/grpc-ecosystem/go-grpc-middleware/v2 to v2.0.1 ([#1261](https://github.com/cloudquery/plugin-sdk/issues/1261)) ([cf57d20](https://github.com/cloudquery/plugin-sdk/commit/cf57d20a17de07a21a5cc364cefc9f4057cb05df))

## [4.11.1](https://github.com/cloudquery/plugin-sdk/compare/v4.11.0...v4.11.1) (2023-09-27)


### Bug Fixes

* **package:** Don't init destinations during package ([#1249](https://github.com/cloudquery/plugin-sdk/issues/1249)) ([f21e963](https://github.com/cloudquery/plugin-sdk/commit/f21e963d4b4c864102ba5afdcd03892e2b0cc969))

## [4.11.0](https://github.com/cloudquery/plugin-sdk/compare/v4.10.2...v4.11.0) (2023-09-25)


### Features

* Provide User with actionable error message when no tables are configured for syncing ([#1243](https://github.com/cloudquery/plugin-sdk/issues/1243)) ([e53d952](https://github.com/cloudquery/plugin-sdk/commit/e53d952fc7347f0c3428a588839f69c2c585a390))


### Bug Fixes

* **deps:** Update github.com/cloudquery/arrow/go/v14 digest to 64e27ff ([#1245](https://github.com/cloudquery/plugin-sdk/issues/1245)) ([ff074f4](https://github.com/cloudquery/plugin-sdk/commit/ff074f4393e15494373578293e7649b6030da803))
* Set GOOS and GOARCH in package command ([#1246](https://github.com/cloudquery/plugin-sdk/issues/1246)) ([119f962](https://github.com/cloudquery/plugin-sdk/commit/119f9628773bf7dcd946fe17571cc523968a36f6))

## [4.10.2](https://github.com/cloudquery/plugin-sdk/compare/v4.10.1...v4.10.2) (2023-09-21)


### Bug Fixes

* Rename plugin type -&gt; kind for consistency with existing configs ([#1240](https://github.com/cloudquery/plugin-sdk/issues/1240)) ([a00b8d0](https://github.com/cloudquery/plugin-sdk/commit/a00b8d0d7161e7c1675cc9d075f967a0c397bee9))

## [4.10.1](https://github.com/cloudquery/plugin-sdk/compare/v4.10.0...v4.10.1) (2023-09-21)


### Bug Fixes

* **scalar:** Don't pass typed nils in list values ([#1226](https://github.com/cloudquery/plugin-sdk/issues/1226)) ([7175e5a](https://github.com/cloudquery/plugin-sdk/commit/7175e5a478524ac99032be6f474f2130ed46985f))
* Skip tables.json when packaging destinations ([#1238](https://github.com/cloudquery/plugin-sdk/issues/1238)) ([f6471e3](https://github.com/cloudquery/plugin-sdk/commit/f6471e3ada92871e951026db9c1bf748a2e0b154))

## [4.10.0](https://github.com/cloudquery/plugin-sdk/compare/v4.9.3...v4.10.0) (2023-09-20)


### Features

* Expose `plugin.JSONSchemaValidator` to be used in schema tests ([#1233](https://github.com/cloudquery/plugin-sdk/issues/1233)) ([ef71086](https://github.com/cloudquery/plugin-sdk/commit/ef71086967c0852438631f9af17fffec304a1ba7))

## [4.9.3](https://github.com/cloudquery/plugin-sdk/compare/v4.9.2...v4.9.3) (2023-09-20)


### Bug Fixes

* Enable format assertion for JSON schema ([#1231](https://github.com/cloudquery/plugin-sdk/issues/1231)) ([b53c5ab](https://github.com/cloudquery/plugin-sdk/commit/b53c5ab519c634c39089232aebe42c0a1f939927))

## [4.9.2](https://github.com/cloudquery/plugin-sdk/compare/v4.9.1...v4.9.2) (2023-09-20)


### Bug Fixes

* **package:** Normalize tables when writing tables.json ([#1227](https://github.com/cloudquery/plugin-sdk/issues/1227)) ([06c84c0](https://github.com/cloudquery/plugin-sdk/commit/06c84c09c731817346644a3d6e337f3732aff023))

## [4.9.1](https://github.com/cloudquery/plugin-sdk/compare/v4.9.0...v4.9.1) (2023-09-20)


### Bug Fixes

* Validate spec only when connection is established ([#1223](https://github.com/cloudquery/plugin-sdk/issues/1223)) ([59aef16](https://github.com/cloudquery/plugin-sdk/commit/59aef16ebe7553faba0dc87b3d81b567acbe77b4))

## [4.9.0](https://github.com/cloudquery/plugin-sdk/compare/v4.8.0...v4.9.0) (2023-09-20)


### Features

* Add support for jsonschema ([#1214](https://github.com/cloudquery/plugin-sdk/issues/1214)) ([2d766dc](https://github.com/cloudquery/plugin-sdk/commit/2d766dc013b80ff62768b1629e69e670df25f4fa))


### Bug Fixes

* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.10.0 ([#1220](https://github.com/cloudquery/plugin-sdk/issues/1220)) ([aa01b1f](https://github.com/cloudquery/plugin-sdk/commit/aa01b1ffbdfb326e9522fd18d73ccf2b653b03df))

## [4.8.0](https://github.com/cloudquery/plugin-sdk/compare/v4.7.1...v4.8.0) (2023-09-19)


### Features

* Add Checksums to package.json format ([#1217](https://github.com/cloudquery/plugin-sdk/issues/1217)) ([720baae](https://github.com/cloudquery/plugin-sdk/commit/720baaec5191706bc52a63478d7b98cdfee6fa47))
* Add message to package command ([#1216](https://github.com/cloudquery/plugin-sdk/issues/1216)) ([44956d9](https://github.com/cloudquery/plugin-sdk/commit/44956d9e5f067909a5126c44e0420c6abf386fce))
* Add shuffle scheduler ([#1218](https://github.com/cloudquery/plugin-sdk/issues/1218)) ([2b1ba30](https://github.com/cloudquery/plugin-sdk/commit/2b1ba309828cfcda3667121557ac30b437a822ce))
* Update package command ([#1211](https://github.com/cloudquery/plugin-sdk/issues/1211)) ([39fc65e](https://github.com/cloudquery/plugin-sdk/commit/39fc65ec5261ab1a070694bed3615613fc3c4d17))


### Bug Fixes

* Add schema version to package.json ([#1212](https://github.com/cloudquery/plugin-sdk/issues/1212)) ([393c94d](https://github.com/cloudquery/plugin-sdk/commit/393c94d3a4b70242aeafe4257cb67cea0ff63236))
* **deps:** Update github.com/cloudquery/arrow/go/v14 digest to 483f6b2 ([#1209](https://github.com/cloudquery/plugin-sdk/issues/1209)) ([179769a](https://github.com/cloudquery/plugin-sdk/commit/179769a2b6dc5900c3078a235c2d19d4091a21ae))
* **deps:** Update github.com/cloudquery/arrow/go/v14 digest to ffb7089 ([#1215](https://github.com/cloudquery/plugin-sdk/issues/1215)) ([70f20bb](https://github.com/cloudquery/plugin-sdk/commit/70f20bb3244cd52d71cf09666bd10b15e1b67d41))
* Use -dir suffix for plugin package arguments ([#1213](https://github.com/cloudquery/plugin-sdk/issues/1213)) ([93f9398](https://github.com/cloudquery/plugin-sdk/commit/93f93988d0334bf2ea101fcc375bad878b396343))

## [4.7.1](https://github.com/cloudquery/plugin-sdk/compare/v4.7.0...v4.7.1) (2023-09-05)


### Bug Fixes

* Relax plugin tables and columns validation ([#1203](https://github.com/cloudquery/plugin-sdk/issues/1203)) ([59c3715](https://github.com/cloudquery/plugin-sdk/commit/59c371528a7f8dcf3618fc768e36cdaacedc55cc))

## [4.7.0](https://github.com/cloudquery/plugin-sdk/compare/v4.6.4...v4.7.0) (2023-09-05)


### Features

* Export `grpczerolog` for reuse ([#1200](https://github.com/cloudquery/plugin-sdk/issues/1200)) ([e2c8fe5](https://github.com/cloudquery/plugin-sdk/commit/e2c8fe5b5b6cae88d04acbb518b05f98554e02dc))

## [4.6.4](https://github.com/cloudquery/plugin-sdk/compare/v4.6.3...v4.6.4) (2023-09-04)


### Bug Fixes

* **caser:** ToSnake does not replace spaces with _ ([#1148](https://github.com/cloudquery/plugin-sdk/issues/1148)) ([329b601](https://github.com/cloudquery/plugin-sdk/commit/329b60164148af2a40fd1d10ef7a607ea1fbb6bc))
* **deps:** Update `github.com/grpc-ecosystem/go-grpc-middleware/v2` to `v2.0.0` ([#1197](https://github.com/cloudquery/plugin-sdk/issues/1197)) ([6d3f752](https://github.com/cloudquery/plugin-sdk/commit/6d3f752bcfaada6a35aeced2503cab7b81362283))

## [4.6.3](https://github.com/cloudquery/plugin-sdk/compare/v4.6.2...v4.6.3) (2023-09-04)


### Bug Fixes

* **deps:** Update github.com/cloudquery/arrow/go/v14 digest to cd3d411 ([#1193](https://github.com/cloudquery/plugin-sdk/issues/1193)) ([3c5e6dd](https://github.com/cloudquery/plugin-sdk/commit/3c5e6ddd8ecb990aa29791af660e7429580f574f))
* Use tables with primary key in `delete-stale` tests ([#1195](https://github.com/cloudquery/plugin-sdk/issues/1195)) ([6dd1730](https://github.com/cloudquery/plugin-sdk/commit/6dd1730b25df3d8153943e1edc05a7afe832edfe))

## [4.6.2](https://github.com/cloudquery/plugin-sdk/compare/v4.6.1...v4.6.2) (2023-09-01)


### Bug Fixes

* Basic delete stale test ([#1189](https://github.com/cloudquery/plugin-sdk/issues/1189)) ([af4aa2e](https://github.com/cloudquery/plugin-sdk/commit/af4aa2e2c896860df16a5a63af2281310d4da268))

## [4.6.1](https://github.com/cloudquery/plugin-sdk/compare/v4.6.0...v4.6.1) (2023-09-01)


### Bug Fixes

* **deps:** Update github.com/apache/arrow/go/v14 digest to 84583d6 ([#1179](https://github.com/cloudquery/plugin-sdk/issues/1179)) ([167fded](https://github.com/cloudquery/plugin-sdk/commit/167fded1e19b2e99ecf90c1eb5514c4dd5613a44))
* **deps:** Update github.com/apache/arrow/go/v14 digest to b6c0ea4 ([#1185](https://github.com/cloudquery/plugin-sdk/issues/1185)) ([7e6bad6](https://github.com/cloudquery/plugin-sdk/commit/7e6bad67ea149a6006a0d2f9049a38d65c516809))
* **deps:** Update golang.org/x/exp digest to d852ddb ([#1181](https://github.com/cloudquery/plugin-sdk/issues/1181)) ([1c8ec87](https://github.com/cloudquery/plugin-sdk/commit/1c8ec87dce3b1a972de07fc4de71dcaa7251be97))
* **deps:** Update golang.org/x/tools digest to 914b218 ([#1183](https://github.com/cloudquery/plugin-sdk/issues/1183)) ([9b9a392](https://github.com/cloudquery/plugin-sdk/commit/9b9a39217e69cfb99cd9c84f0e116f508bd41ba7))
* **deps:** Update google.golang.org/genproto digest to b8732ec ([#1182](https://github.com/cloudquery/plugin-sdk/issues/1182)) ([8d98808](https://github.com/cloudquery/plugin-sdk/commit/8d988082ad9bea01220d21c04dcb447da7456e86))
* **deps:** Update google.golang.org/genproto/googleapis/api digest to b8732ec ([#1184](https://github.com/cloudquery/plugin-sdk/issues/1184)) ([c74fb1d](https://github.com/cloudquery/plugin-sdk/commit/c74fb1dbf60daffff68a81550e014d13cef098fa))
* **deps:** Update google.golang.org/genproto/googleapis/rpc digest to b8732ec ([#1186](https://github.com/cloudquery/plugin-sdk/issues/1186)) ([15cea46](https://github.com/cloudquery/plugin-sdk/commit/15cea46d59bc6bb2e1c82022497ae3322e9190ff))
* **test:** Truncate sync time based on test options in `testDeleteStaleBasic` ([#1187](https://github.com/cloudquery/plugin-sdk/issues/1187)) ([faa64b0](https://github.com/cloudquery/plugin-sdk/commit/faa64b08ea80a173ebd38c3b8799576716f4bacd))

## [4.6.0](https://github.com/cloudquery/plugin-sdk/compare/v4.5.7...v4.6.0) (2023-08-31)


### Features

* Extensive testing for `delete-stale` ([#1175](https://github.com/cloudquery/plugin-sdk/issues/1175)) ([304e4eb](https://github.com/cloudquery/plugin-sdk/commit/304e4eba408a0782f6b47e1c47a7f86f81588ac1))

## [4.5.7](https://github.com/cloudquery/plugin-sdk/compare/v4.5.6...v4.5.7) (2023-08-28)


### Bug Fixes

* **deps:** Update `github.com/cloudquery/arrow/go/v13` to `github.com/cloudquery/arrow/go/v14` ([#1169](https://github.com/cloudquery/plugin-sdk/issues/1169)) ([6be8194](https://github.com/cloudquery/plugin-sdk/commit/6be8194a27a2d562479e8980c213e8ab152fc972))

## [4.5.6](https://github.com/cloudquery/plugin-sdk/compare/v4.5.5...v4.5.6) (2023-08-28)


### Bug Fixes

* Don't send migrate messages in destination v1 write ([#1167](https://github.com/cloudquery/plugin-sdk/issues/1167)) ([9ed543c](https://github.com/cloudquery/plugin-sdk/commit/9ed543c5e10a46fa0cb9c0ff8e942e12d2c48f37))

## [4.5.5](https://github.com/cloudquery/plugin-sdk/compare/v4.5.4...v4.5.5) (2023-08-22)


### Bug Fixes

* Skip double migration test in forced mode ([#1163](https://github.com/cloudquery/plugin-sdk/issues/1163)) ([e7b5ed1](https://github.com/cloudquery/plugin-sdk/commit/e7b5ed18868f38ae09f8a392c19566f40d0e5a83))

## [4.5.4](https://github.com/cloudquery/plugin-sdk/compare/v4.5.3...v4.5.4) (2023-08-22)


### Bug Fixes

* Fix testdata generation ([#1160](https://github.com/cloudquery/plugin-sdk/issues/1160)) ([f07869a](https://github.com/cloudquery/plugin-sdk/commit/f07869aa82f92d745a30aaa35f33ae3bf31a7f50))

## [4.5.3](https://github.com/cloudquery/plugin-sdk/compare/v4.5.2...v4.5.3) (2023-08-21)


### Bug Fixes

* Ease diff code reading ([#1157](https://github.com/cloudquery/plugin-sdk/issues/1157)) ([72fc538](https://github.com/cloudquery/plugin-sdk/commit/72fc538af0eec502bc0287dc3ab4b3f989adb448))

## [4.5.2](https://github.com/cloudquery/plugin-sdk/compare/v4.5.1...v4.5.2) (2023-08-21)


### Bug Fixes

* Change `testdata.Generate` signature ([#1153](https://github.com/cloudquery/plugin-sdk/issues/1153)) ([86e717a](https://github.com/cloudquery/plugin-sdk/commit/86e717a442c43c945239cbdcbc79ac4ece97c7c2))
* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to 5b83d4f ([#1154](https://github.com/cloudquery/plugin-sdk/issues/1154)) ([8558dd1](https://github.com/cloudquery/plugin-sdk/commit/8558dd102d359159dec64ad099bc417c97cc1477))
* **deps:** Update module github.com/cloudquery/plugin-sdk/v4 to v4.5.1 ([#1150](https://github.com/cloudquery/plugin-sdk/issues/1150)) ([b3f41b1](https://github.com/cloudquery/plugin-sdk/commit/b3f41b1620c912383e5ef83c0765af03d3224fc7))

## [4.5.1](https://github.com/cloudquery/plugin-sdk/compare/v4.5.0...v4.5.1) (2023-08-18)


### Bug Fixes

* Bring back plugin validation ([#1108](https://github.com/cloudquery/plugin-sdk/issues/1108)) ([61765a7](https://github.com/cloudquery/plugin-sdk/commit/61765a7ce6a2ec1b88ab97fd2f53514b88df4d36))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.9.3 ([#1149](https://github.com/cloudquery/plugin-sdk/issues/1149)) ([e1ea578](https://github.com/cloudquery/plugin-sdk/commit/e1ea57877f82cafce7c42a826dddc0fe22c9ff51))
* **deps:** Update module github.com/cloudquery/plugin-sdk/v4 to v4.5.0 ([#1145](https://github.com/cloudquery/plugin-sdk/issues/1145)) ([70d12e4](https://github.com/cloudquery/plugin-sdk/commit/70d12e476581c6388d08b056afd955a25dcaf888))

## [4.5.0](https://github.com/cloudquery/plugin-sdk/compare/v4.4.0...v4.5.0) (2023-08-14)


### Features

* Add publish command ([#1143](https://github.com/cloudquery/plugin-sdk/issues/1143)) ([fdd44d5](https://github.com/cloudquery/plugin-sdk/commit/fdd44d5d3a9ce12d59e168ea691a343f6f219694))


### Bug Fixes

* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to e9683e1 ([#1144](https://github.com/cloudquery/plugin-sdk/issues/1144)) ([763c549](https://github.com/cloudquery/plugin-sdk/commit/763c549a783f69d7adfb7291534d3d2b25d697e5))
* Scalar timestamp parsing ([#1109](https://github.com/cloudquery/plugin-sdk/issues/1109)) ([c15b214](https://github.com/cloudquery/plugin-sdk/commit/c15b214a346fa8a89c929858c2623317e7048211))

## [4.4.0](https://github.com/cloudquery/plugin-sdk/compare/v4.3.1...v4.4.0) (2023-08-08)


### Features

* Add Unflatten tables method ([#1138](https://github.com/cloudquery/plugin-sdk/issues/1138)) ([848e505](https://github.com/cloudquery/plugin-sdk/commit/848e505ba49bdb4fb45cfa8bb7b9b7538afc785e))

## [4.3.1](https://github.com/cloudquery/plugin-sdk/compare/v4.3.0...v4.3.1) (2023-08-08)


### Bug Fixes

* **plugin-tables:** Add missing `skip_dependent_tables` ([#1136](https://github.com/cloudquery/plugin-sdk/issues/1136)) ([65e9f1a](https://github.com/cloudquery/plugin-sdk/commit/65e9f1a9d81d4534e8a637ed5db57071fe91d831))

## [4.3.0](https://github.com/cloudquery/plugin-sdk/compare/v4.2.6...v4.3.0) (2023-08-08)


### Features

* Add more metadata to tables needed for docs generation ([#1129](https://github.com/cloudquery/plugin-sdk/issues/1129)) ([3dbd7f3](https://github.com/cloudquery/plugin-sdk/commit/3dbd7f32cdcb87dd0b7cd4dd9b71c2552b25b30e))

## [4.2.6](https://github.com/cloudquery/plugin-sdk/compare/v4.2.5...v4.2.6) (2023-08-08)


### Bug Fixes

* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to f53878d ([#1132](https://github.com/cloudquery/plugin-sdk/issues/1132)) ([0c47570](https://github.com/cloudquery/plugin-sdk/commit/0c475702592506e2fce708384dd2bd5c8b9da827))
* **writers:** StreamingBatchWriter hangs with non-append mode ([#1131](https://github.com/cloudquery/plugin-sdk/issues/1131)) ([806c85d](https://github.com/cloudquery/plugin-sdk/commit/806c85d92bb9152b0469a1e30e167a662ebd8015))

## [4.2.5](https://github.com/cloudquery/plugin-sdk/compare/v4.2.4...v4.2.5) (2023-08-02)


### Bug Fixes

* Nulls in lists ([#1127](https://github.com/cloudquery/plugin-sdk/issues/1127)) ([dc1e6be](https://github.com/cloudquery/plugin-sdk/commit/dc1e6bee22dbbbeb15b3586a8815598d50a6b434))

## [4.2.4](https://github.com/cloudquery/plugin-sdk/compare/v4.2.3...v4.2.4) (2023-08-02)


### Bug Fixes

* Check record equality before generating diff ([#1123](https://github.com/cloudquery/plugin-sdk/issues/1123)) ([b2e6331](https://github.com/cloudquery/plugin-sdk/commit/b2e63318befaf3cf4f633a95f08178ef7dbbed18))
* **deps:** Update github.com/apache/arrow/go/v13 digest to 112f949 ([#1115](https://github.com/cloudquery/plugin-sdk/issues/1115)) ([ed0e4e0](https://github.com/cloudquery/plugin-sdk/commit/ed0e4e03c271d7232258c4efaec3708f645e7d5e))
* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to 10df4b9 ([#1110](https://github.com/cloudquery/plugin-sdk/issues/1110)) ([636084c](https://github.com/cloudquery/plugin-sdk/commit/636084cb28281e4cccad76b8aff5a18306855eb1))
* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to 3452eb0 ([#1114](https://github.com/cloudquery/plugin-sdk/issues/1114)) ([af83988](https://github.com/cloudquery/plugin-sdk/commit/af839886025f534bf28484b49345faca9dcd1735))
* **deps:** Update golang.org/x/exp digest to b0cb94b ([#1116](https://github.com/cloudquery/plugin-sdk/issues/1116)) ([4a6dc5b](https://github.com/cloudquery/plugin-sdk/commit/4a6dc5b8a657ad09a4476305fb64629fbec6463f))
* **deps:** Update google.golang.org/genproto digest to e0aa005 ([#1117](https://github.com/cloudquery/plugin-sdk/issues/1117)) ([5fa4d51](https://github.com/cloudquery/plugin-sdk/commit/5fa4d5184b333fb7d7a4a2c5bed2ca695eba80fe))
* **deps:** Update google.golang.org/genproto/googleapis/api digest to e0aa005 ([#1118](https://github.com/cloudquery/plugin-sdk/issues/1118)) ([939060f](https://github.com/cloudquery/plugin-sdk/commit/939060fbbca30e17de0537d5eec42ff15beaceab))
* **deps:** Update google.golang.org/genproto/googleapis/rpc digest to e0aa005 ([#1119](https://github.com/cloudquery/plugin-sdk/issues/1119)) ([0a9f8ea](https://github.com/cloudquery/plugin-sdk/commit/0a9f8eaa4777764c654460bc7328281df9bf0ac8))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.9.0 ([#1112](https://github.com/cloudquery/plugin-sdk/issues/1112)) ([3831a88](https://github.com/cloudquery/plugin-sdk/commit/3831a88c3a4afa5f3764c908a2ae098c4f3cba5f))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.9.1 ([#1113](https://github.com/cloudquery/plugin-sdk/issues/1113)) ([67bc46e](https://github.com/cloudquery/plugin-sdk/commit/67bc46e957d6ec6e21f018823700eccb3af96027))
* **deps:** Update module github.com/klauspost/compress to v1.16.7 ([#1120](https://github.com/cloudquery/plugin-sdk/issues/1120)) ([e41a303](https://github.com/cloudquery/plugin-sdk/commit/e41a303142475b9b796214ba8909962a7a43e6a2))
* **deps:** Update module github.com/pierrec/lz4/v4 to v4.1.18 ([#1121](https://github.com/cloudquery/plugin-sdk/issues/1121)) ([6829b63](https://github.com/cloudquery/plugin-sdk/commit/6829b6356ba7b543f35c0c22d2f22a6789c59e9b))
* Process nulls for tested types, too (maps, lists, structs) ([#1125](https://github.com/cloudquery/plugin-sdk/issues/1125)) ([4a1f315](https://github.com/cloudquery/plugin-sdk/commit/4a1f31514aee9021a4c667f559eefe08b42e5c14))

## [4.2.3](https://github.com/cloudquery/plugin-sdk/compare/v4.2.2...v4.2.3) (2023-07-18)


### Bug Fixes

* **streamingbatchwriter:** Missing tickerFn on DeleteWorker ([#1103](https://github.com/cloudquery/plugin-sdk/issues/1103)) ([91eae56](https://github.com/cloudquery/plugin-sdk/commit/91eae56526588f944bdfaceb5c89de8473d84779))

## [4.2.2](https://github.com/cloudquery/plugin-sdk/compare/v4.2.1...v4.2.2) (2023-07-18)


### Bug Fixes

* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to 8e2219b ([#1095](https://github.com/cloudquery/plugin-sdk/issues/1095)) ([2f6bd18](https://github.com/cloudquery/plugin-sdk/commit/2f6bd18db9aac05ade8c21260c9f4c6fca8555ea))
* **testing:** Force migrations should allow table drops ([#1101](https://github.com/cloudquery/plugin-sdk/issues/1101)) ([5dbb23e](https://github.com/cloudquery/plugin-sdk/commit/5dbb23eb9ceab7e43a672fbc60060934b490b47c))

## [4.2.1](https://github.com/cloudquery/plugin-sdk/compare/v4.2.0...v4.2.1) (2023-07-17)


### Bug Fixes

* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.8.0 ([#1098](https://github.com/cloudquery/plugin-sdk/issues/1098)) ([cbbecb8](https://github.com/cloudquery/plugin-sdk/commit/cbbecb8ebe32b21b59d79ec5548347d86b7a370a))

## [4.2.0](https://github.com/cloudquery/plugin-sdk/compare/v4.1.1...v4.2.0) (2023-07-17)


### Features

* Add initial version of open-telemetry ([#1097](https://github.com/cloudquery/plugin-sdk/issues/1097)) ([09a880c](https://github.com/cloudquery/plugin-sdk/commit/09a880c3ad420b991f0bc21b3cb9fba3226a6d91))


### Bug Fixes

* Differentiate between errgroup context and global context being cance… ([#1082](https://github.com/cloudquery/plugin-sdk/issues/1082)) ([0532f88](https://github.com/cloudquery/plugin-sdk/commit/0532f881067c142fd7799037990963b3ceee61fa))

## [4.1.1](https://github.com/cloudquery/plugin-sdk/compare/v4.1.0...v4.1.1) (2023-07-14)


### Bug Fixes

* Add `NoConnection` to init request ([#1092](https://github.com/cloudquery/plugin-sdk/issues/1092)) ([ba16cfd](https://github.com/cloudquery/plugin-sdk/commit/ba16cfd902fa0ba86ca826fa761d1d0e72688bc0))

## [4.1.0](https://github.com/cloudquery/plugin-sdk/compare/v4.0.0...v4.1.0) (2023-07-14)


### Features

* Add `plugin.ValidateNoEmptyColumns` ([#1085](https://github.com/cloudquery/plugin-sdk/issues/1085)) ([32e1215](https://github.com/cloudquery/plugin-sdk/commit/32e1215ef3d59a1e56d14bbb342f1f33dd76146b))


### Bug Fixes

* Add random suffix to test table names ([#1086](https://github.com/cloudquery/plugin-sdk/issues/1086)) ([ad16b20](https://github.com/cloudquery/plugin-sdk/commit/ad16b20eded7ac587e41d154a7ba4e3f801e2c99))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.7.0 ([#1091](https://github.com/cloudquery/plugin-sdk/issues/1091)) ([fb124a2](https://github.com/cloudquery/plugin-sdk/commit/fb124a207d05c00c1c974efd900d06d8eb9374db))
* **testing:** Comply with given TimePrecision ([#1089](https://github.com/cloudquery/plugin-sdk/issues/1089)) ([d16ed0f](https://github.com/cloudquery/plugin-sdk/commit/d16ed0f823ee0dcad8f8fa21df64be2aa5b9bd04))

## [4.0.0](https://github.com/cloudquery/plugin-sdk/compare/v4.8.1-rc1...v4.0.0) (2023-07-12)


### Bug Fixes

* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to 0a52533 ([#1083](https://github.com/cloudquery/plugin-sdk/issues/1083)) ([0370294](https://github.com/cloudquery/plugin-sdk/commit/0370294523989c73afd808ac9678bc9018210c41))
* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to a2a76eb ([#1084](https://github.com/cloudquery/plugin-sdk/issues/1084)) ([26df75f](https://github.com/cloudquery/plugin-sdk/commit/26df75f3fc38ee8cd5c644cb62cd4ce5c720df25))
* **types-inet:** Align logic with scalar package, set `net.IPNet` `IP` field after parsing `ParseCIDR` ([#982](https://github.com/cloudquery/plugin-sdk/issues/982)) ([fa07032](https://github.com/cloudquery/plugin-sdk/commit/fa0703271ea05e46cfe171ad1f488ddbefdd96d2))
* Use background ctx in batchwriter worker ([#1079](https://github.com/cloudquery/plugin-sdk/issues/1079)) ([dea8168](https://github.com/cloudquery/plugin-sdk/commit/dea8168c37da58f0aaf6273446a68f8d752c9cef))


### Miscellaneous Chores

* release 4.0.0 ([a80ee69](https://github.com/cloudquery/plugin-sdk/commit/a80ee69c795819dfaff2512fee8a66135bf7aca8))

## [4.8.1-rc1](https://github.com/cloudquery/plugin-sdk/compare/v4.8.0-rc1...v4.8.1-rc1) (2023-07-05)


### Bug Fixes

* **scheduler:** Concurrency as `int` ([#1077](https://github.com/cloudquery/plugin-sdk/issues/1077)) ([30ba6d7](https://github.com/cloudquery/plugin-sdk/commit/30ba6d758cedea74928be4901a6f78696c0c7247))

## [4.8.0-rc1](https://github.com/cloudquery/plugin-sdk/compare/v4.7.1-rc1...v4.8.0-rc1) (2023-07-05)


### Features

* **transformers:** Add `Apply` to apply extra transformations ([#1069](https://github.com/cloudquery/plugin-sdk/issues/1069)) ([a40598e](https://github.com/cloudquery/plugin-sdk/commit/a40598e6c6fe409e7170d2c1553c85050c196562))


### Bug Fixes

* Deterministic ordering for records returned by readAll in tests ([#1072](https://github.com/cloudquery/plugin-sdk/issues/1072)) ([cf7510f](https://github.com/cloudquery/plugin-sdk/commit/cf7510fdb594f7772c8507b0f9d394c862172a9f))
* Handle null-related test options ([#1074](https://github.com/cloudquery/plugin-sdk/issues/1074)) ([88f08ee](https://github.com/cloudquery/plugin-sdk/commit/88f08ee35601d98385f3f6da4c2a57cc3ce81bd5))
* **naming:** Rename `SyncMessages.InsertMessage()` to `SyncMessages.GetInserts()` ([#1070](https://github.com/cloudquery/plugin-sdk/issues/1070)) ([ab9e768](https://github.com/cloudquery/plugin-sdk/commit/ab9e768f8e11d008236a0ff861734841524a9aea))
* Reset timers on flush ([#1076](https://github.com/cloudquery/plugin-sdk/issues/1076)) ([767327f](https://github.com/cloudquery/plugin-sdk/commit/767327fd5decbbbbd9e3a5c9664c73425b7b6dbe))
* Reverse order of records in memdb ([#1075](https://github.com/cloudquery/plugin-sdk/issues/1075)) ([8356590](https://github.com/cloudquery/plugin-sdk/commit/8356590c03f84b7ba69e7f661aba2b2a889fb2dd))
* **scalar:** Test `AppendTime` on TimestampBuilder ([#1068](https://github.com/cloudquery/plugin-sdk/issues/1068)) ([888c9ee](https://github.com/cloudquery/plugin-sdk/commit/888c9ee7e88f145b1baa2758f71bee1a24e5f60e))
* **testdata:** Exclude only the correct type ([#1067](https://github.com/cloudquery/plugin-sdk/issues/1067)) ([1c72fb2](https://github.com/cloudquery/plugin-sdk/commit/1c72fb2fc532afee425ded6f324aa7e6cd9875b1))

## [4.7.1-rc1](https://github.com/cloudquery/plugin-sdk/compare/v4.7.0-rc1...v4.7.1-rc1) (2023-07-04)


### Bug Fixes

* Add AddCqIDs helper function ([#1065](https://github.com/cloudquery/plugin-sdk/issues/1065)) ([911762d](https://github.com/cloudquery/plugin-sdk/commit/911762d2f790c9ed9facbea567dc6ff2100a6adf))
* Check record data in tests ([#1062](https://github.com/cloudquery/plugin-sdk/issues/1062)) ([f13e4cc](https://github.com/cloudquery/plugin-sdk/commit/f13e4cc4a8d401fca314c5b266b75700bdc47088))
* **configtype:** Add `Equal()` method to `Duration` ([#1059](https://github.com/cloudquery/plugin-sdk/issues/1059)) ([57c7bc2](https://github.com/cloudquery/plugin-sdk/commit/57c7bc230c3ad3150f37b4f36b8e479b1c45c64f))
* Conversion and test fixes ([#1064](https://github.com/cloudquery/plugin-sdk/issues/1064)) ([36b65cb](https://github.com/cloudquery/plugin-sdk/commit/36b65cb9132470a835aac3e1f02c5c49c3fb70f6))
* Fix test assertions for records ([#1066](https://github.com/cloudquery/plugin-sdk/issues/1066)) ([a9bd88f](https://github.com/cloudquery/plugin-sdk/commit/a9bd88f8db0a71dc4f8ea713ff35b206a0485d9a))
* **testdata:** Add missing column types ([#1061](https://github.com/cloudquery/plugin-sdk/issues/1061)) ([f5d01c9](https://github.com/cloudquery/plugin-sdk/commit/f5d01c9adf8b532e97372245827334ec6d5c6e64))

## [4.7.0-rc1](https://github.com/cloudquery/plugin-sdk/compare/v4.6.1-rc1...v4.7.0-rc1) (2023-07-04)


### Features

* Add `WriteInserts.GetRecords()` ([#1053](https://github.com/cloudquery/plugin-sdk/issues/1053)) ([05e1edd](https://github.com/cloudquery/plugin-sdk/commit/05e1eddff293504be015b4ee76c911f35b91bfba))
* Add batch timeout support to mixed batch writer ([#1055](https://github.com/cloudquery/plugin-sdk/issues/1055)) ([7fe7c64](https://github.com/cloudquery/plugin-sdk/commit/7fe7c642287609f9ea9a65e604741f2164f8f8ce))
* Add Duration configtype ([#1014](https://github.com/cloudquery/plugin-sdk/issues/1014)) ([fbde15a](https://github.com/cloudquery/plugin-sdk/commit/fbde15a62c055270ce03dc9bbbced9400c53e943))


### Bug Fixes

* Fix timer logic in batch writers ([#1056](https://github.com/cloudquery/plugin-sdk/issues/1056)) ([9179e7f](https://github.com/cloudquery/plugin-sdk/commit/9179e7f9184260e36018f83d13a5229f47dafdac))

## [4.6.1-rc1](https://github.com/cloudquery/plugin-sdk/compare/v4.6.0-rc1...v4.6.1-rc1) (2023-07-03)


### Bug Fixes

* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.6.0 ([#1050](https://github.com/cloudquery/plugin-sdk/issues/1050)) ([ba632d1](https://github.com/cloudquery/plugin-sdk/commit/ba632d1dd5feb98d64cf762ef55d7d6dd03fc2e2))
* Make gen docs work without auth ([#1052](https://github.com/cloudquery/plugin-sdk/issues/1052)) ([504f849](https://github.com/cloudquery/plugin-sdk/commit/504f8498f9f1317be849c1090728c81a76daa5ca))

## [4.6.0-rc1](https://github.com/cloudquery/plugin-sdk/compare/v4.5.0-rc1...v4.6.0-rc1) (2023-07-03)


### Features

* Add state.NoOpClient ([#1047](https://github.com/cloudquery/plugin-sdk/issues/1047)) ([ee1ee5f](https://github.com/cloudquery/plugin-sdk/commit/ee1ee5fdb455ef6216b5a591c55b94b27cd96277))

## [4.5.0-rc1](https://github.com/cloudquery/plugin-sdk/compare/v4.4.2-rc1...v4.5.0-rc1) (2023-07-03)


### Features

* **writers:** More unimplemented writer helpers ([#1038](https://github.com/cloudquery/plugin-sdk/issues/1038)) ([b1ad878](https://github.com/cloudquery/plugin-sdk/commit/b1ad878f7bb99403c4516f134a76fc165758ec0f))


### Bug Fixes

* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to df3b664 ([#1043](https://github.com/cloudquery/plugin-sdk/issues/1043)) ([5b95fce](https://github.com/cloudquery/plugin-sdk/commit/5b95fceffb74515d7141a2d56c6be1a78f0e562c))
* Make scheduler stateful to support sync option ([#1046](https://github.com/cloudquery/plugin-sdk/issues/1046)) ([d683eff](https://github.com/cloudquery/plugin-sdk/commit/d683eff3a0bc31d9dde97c50a48bfa94e5ff2895))
* **writers:** Require `Close()` for `StreamingBatchWriter` ([#1045](https://github.com/cloudquery/plugin-sdk/issues/1045)) ([2078e84](https://github.com/cloudquery/plugin-sdk/commit/2078e842f55b37966214f281ca7ac230cec8dc73))

## [4.4.2-rc1](https://github.com/cloudquery/plugin-sdk/compare/v4.4.1-rc1...v4.4.2-rc1) (2023-07-02)


### Bug Fixes

* Add GetMessageByTable to WriteMigrateTables ([#1041](https://github.com/cloudquery/plugin-sdk/issues/1041)) ([8a23f68](https://github.com/cloudquery/plugin-sdk/commit/8a23f6801a66ad35b755081d6669458acd4bc186))

## [4.4.1-rc1](https://github.com/cloudquery/plugin-sdk/compare/v4.4.0-rc1...v4.4.1-rc1) (2023-07-01)


### Bug Fixes

* **deps:** Update github.com/apache/arrow/go/v13 digest to 5a06b2e ([#1032](https://github.com/cloudquery/plugin-sdk/issues/1032)) ([d369262](https://github.com/cloudquery/plugin-sdk/commit/d36926212e837eb833e49efd20755adfb886804d))
* **deps:** Update golang.org/x/exp digest to 97b1e66 ([#1033](https://github.com/cloudquery/plugin-sdk/issues/1033)) ([791e60a](https://github.com/cloudquery/plugin-sdk/commit/791e60aa6113e2d70d245e28d4b3c4f910c32a25))
* **deps:** Update google.golang.org/genproto/googleapis/rpc digest to 9506855 ([#1034](https://github.com/cloudquery/plugin-sdk/issues/1034)) ([6999d11](https://github.com/cloudquery/plugin-sdk/commit/6999d11b674235875b22f5a86d766206f6f0b56c))
* **deps:** Update module github.com/goccy/go-json to v0.10.2 ([#1035](https://github.com/cloudquery/plugin-sdk/issues/1035)) ([521eb13](https://github.com/cloudquery/plugin-sdk/commit/521eb13730e761c0bc9f12bc1769d61cc24fec48))
* **deps:** Update module github.com/klauspost/compress to v1.16.6 ([#1036](https://github.com/cloudquery/plugin-sdk/issues/1036)) ([76bfc85](https://github.com/cloudquery/plugin-sdk/commit/76bfc8544fea19e8cb4dc3999fb0c3956f1f4e36))
* **serve:** Confusing message ([#1031](https://github.com/cloudquery/plugin-sdk/issues/1031)) ([ee873c9](https://github.com/cloudquery/plugin-sdk/commit/ee873c96f83ab05a4aa67fdae14eb5aa9d32471c))
* State add flush and fix migration bug ([#1039](https://github.com/cloudquery/plugin-sdk/issues/1039)) ([8c10291](https://github.com/cloudquery/plugin-sdk/commit/8c1029124c73e32f8581951caa8ad737ac0c2fba))

## [4.4.0-rc1](https://github.com/cloudquery/plugin-sdk/compare/v4.3.1-rc1...v4.4.0-rc1) (2023-06-30)


### Features

* Implement plugin Read ([#1027](https://github.com/cloudquery/plugin-sdk/issues/1027)) ([09fb4ce](https://github.com/cloudquery/plugin-sdk/commit/09fb4cede4159e23120726bac3d674e53e89f614))


### Bug Fixes

* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.4.0 ([#1005](https://github.com/cloudquery/plugin-sdk/issues/1005)) ([40f1c77](https://github.com/cloudquery/plugin-sdk/commit/40f1c77193e6ec380ad417ad84cea3b7fb25f810))
* Update to plugin-pb v1.5.0 ([#1026](https://github.com/cloudquery/plugin-sdk/issues/1026)) ([abe2557](https://github.com/cloudquery/plugin-sdk/commit/abe25573411e0ce1b75f76fdcd949ef497674e9d))

## [4.3.1-rc1](https://github.com/cloudquery/plugin-sdk/compare/v4.3.0-rc1...v4.3.1-rc1) (2023-06-29)


### Bug Fixes

* Enable double migration test ([#1023](https://github.com/cloudquery/plugin-sdk/issues/1023)) ([466796b](https://github.com/cloudquery/plugin-sdk/commit/466796bd312b92c9646a2ef1a170bfc4e4b27419))
* Put null helpers back ([#1002](https://github.com/cloudquery/plugin-sdk/issues/1002)) ([95ed5df](https://github.com/cloudquery/plugin-sdk/commit/95ed5dfaf505a3ecdca6be03e8cd46a5cc5a3f23))

## [4.3.0-rc1](https://github.com/cloudquery/plugin-sdk/compare/v4.2.2-rc1...v4.3.0-rc1) (2023-06-29)


### Features

* Use named message slice types in writers ([#1017](https://github.com/cloudquery/plugin-sdk/issues/1017)) ([e290234](https://github.com/cloudquery/plugin-sdk/commit/e29023429d699f095c4240d8097dda850d1933f2))
* **writers:** Add `streamingbatchwriter.Unimplemented*` handlers ([#1022](https://github.com/cloudquery/plugin-sdk/issues/1022)) ([88f4909](https://github.com/cloudquery/plugin-sdk/commit/88f4909e07c0042be20c41288eedfaa729559b5a))


### Bug Fixes

* **writers:** Allow zero timeout, remove unused timeout options from mixedbatchwriter ([#1020](https://github.com/cloudquery/plugin-sdk/issues/1020)) ([282ee45](https://github.com/cloudquery/plugin-sdk/commit/282ee45e552bde91dea36e3c9d1410e5066365ba))
* **writers:** Don't export defaults ([#1013](https://github.com/cloudquery/plugin-sdk/issues/1013)) ([d11dd56](https://github.com/cloudquery/plugin-sdk/commit/d11dd56a0bda79865be505e14159e807a6033431))

## [4.2.2-rc1](https://github.com/cloudquery/plugin-sdk/compare/v4.2.1-rc1...v4.2.2-rc1) (2023-06-29)


### Bug Fixes

* Add backward compatibility for batch_size ([#1018](https://github.com/cloudquery/plugin-sdk/issues/1018)) ([3a72b2f](https://github.com/cloudquery/plugin-sdk/commit/3a72b2f3b9570ba901267337871575f0ed4301a7))

## [4.2.1-rc1](https://github.com/cloudquery/plugin-sdk/compare/v4.2.0-rc1...v4.2.1-rc1) (2023-06-29)


### Bug Fixes

* Add back testing for all types ([#1015](https://github.com/cloudquery/plugin-sdk/issues/1015)) ([8525cc9](https://github.com/cloudquery/plugin-sdk/commit/8525cc966d6da274b2e06da0a529c80a8650fa60))

## [4.2.0-rc1](https://github.com/cloudquery/plugin-sdk/compare/v4.1.0-rc1...v4.2.0-rc1) (2023-06-28)


### Features

* Add StreamingBatchWriter ([#1004](https://github.com/cloudquery/plugin-sdk/issues/1004)) ([986340f](https://github.com/cloudquery/plugin-sdk/commit/986340fc624a0370a726150fb733b723ec96fe74))


### Bug Fixes

* **batchwriter:** Allow zero batch size, flush before exceeding batch size instead of after ([#1008](https://github.com/cloudquery/plugin-sdk/issues/1008)) ([c7ea17b](https://github.com/cloudquery/plugin-sdk/commit/c7ea17b08f340bb5eced749ec76eeb7a480cfa61))
* Naming fix for `messages.InsertMessage` (now `messages.GetInserts`) ([#1000](https://github.com/cloudquery/plugin-sdk/issues/1000)) ([b1e2bd4](https://github.com/cloudquery/plugin-sdk/commit/b1e2bd4d1a5c5f904b82d6f41ed6eb0f26ef91cc))
* Update scheduler for JSON marshal / unmarshal ([#1006](https://github.com/cloudquery/plugin-sdk/issues/1006)) ([970bad1](https://github.com/cloudquery/plugin-sdk/commit/970bad1fcf1f0e6aadcf1d0f381c8cd8177e8513))
* **writers:** Move to sub packages ([#1011](https://github.com/cloudquery/plugin-sdk/issues/1011)) ([826e816](https://github.com/cloudquery/plugin-sdk/commit/826e816d8b8ec9668806b9d34ee5f92bd4a7ff56))

## [4.1.0-rc1](https://github.com/cloudquery/plugin-sdk/compare/v4.0.2-rc1...v4.1.0-rc1) (2023-06-28)


### Features

* Split sync and write messages ([#1009](https://github.com/cloudquery/plugin-sdk/issues/1009)) ([6e35a5f](https://github.com/cloudquery/plugin-sdk/commit/6e35a5f270d5938d6f108a37efce9a51afb35119))


### Bug Fixes

* **testing:** Grammar ([#1003](https://github.com/cloudquery/plugin-sdk/issues/1003)) ([c79cde4](https://github.com/cloudquery/plugin-sdk/commit/c79cde4dfced980dbdcc721f048adfb0686174e2))

## [4.0.2-rc1](https://github.com/cloudquery/plugin-sdk/compare/v4.0.1-rc1...v4.0.2-rc1) (2023-06-26)


### Bug Fixes

* Set Sync option fields ([#997](https://github.com/cloudquery/plugin-sdk/issues/997)) ([29223ba](https://github.com/cloudquery/plugin-sdk/commit/29223baba5e1fd59087dc480bcf47066f1bda91c))

## [4.0.1-rc1](https://github.com/cloudquery/plugin-sdk/compare/v4.0.0-rc1...v4.0.1-rc1) (2023-06-26)


### Bug Fixes

* Close files in docs ([#995](https://github.com/cloudquery/plugin-sdk/issues/995)) ([152f1e1](https://github.com/cloudquery/plugin-sdk/commit/152f1e12df87b31d15ccf5d6c147dc5aef5e5181))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.3.4 ([#994](https://github.com/cloudquery/plugin-sdk/issues/994)) ([24ad6fe](https://github.com/cloudquery/plugin-sdk/commit/24ad6fefd6bbf32cf95d8713ba0e1dfc1413367a))

## [4.0.0-rc1](https://github.com/cloudquery/plugin-sdk/compare/v3.10.6...v4.0.0-rc1) (2023-06-26)


### ⚠ BREAKING CHANGES

* Update to SDK V4 ([#984](https://github.com/cloudquery/plugin-sdk/issues/984))

### Features

* Update to SDK V4 ([#984](https://github.com/cloudquery/plugin-sdk/issues/984)) ([24b19c9](https://github.com/cloudquery/plugin-sdk/commit/24b19c92db5792a0d6d531c7af60d13c20049140))


### Bug Fixes

* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to 0656028 ([#991](https://github.com/cloudquery/plugin-sdk/issues/991)) ([bc9e6e1](https://github.com/cloudquery/plugin-sdk/commit/bc9e6e1ae1bff7a6f2c2699a8a113323d689f1a5))
* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to 1e68c51 ([#973](https://github.com/cloudquery/plugin-sdk/issues/973)) ([f5cdc95](https://github.com/cloudquery/plugin-sdk/commit/f5cdc95ab76bcae67c2f9001da5000bb726feb5e))
* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to 43638cb ([#978](https://github.com/cloudquery/plugin-sdk/issues/978)) ([fb76304](https://github.com/cloudquery/plugin-sdk/commit/fb76304c503b42957b5909cf33d219c9aa4c2934))
* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to 4d76231 ([#970](https://github.com/cloudquery/plugin-sdk/issues/970)) ([646cbc0](https://github.com/cloudquery/plugin-sdk/commit/646cbc0d4471093de309a510b1f16a875a8aa484))
* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to 8366a22 ([#981](https://github.com/cloudquery/plugin-sdk/issues/981)) ([097621f](https://github.com/cloudquery/plugin-sdk/commit/097621f02e4fe1258290c3bfdd744a8fc3ab1c15))
* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to 95d3199 ([#980](https://github.com/cloudquery/plugin-sdk/issues/980)) ([b7bcd93](https://github.com/cloudquery/plugin-sdk/commit/b7bcd9326eb87f43e77e509ec5151c7f4a6f01b4))
* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to 9a09f72 ([#974](https://github.com/cloudquery/plugin-sdk/issues/974)) ([5acec96](https://github.com/cloudquery/plugin-sdk/commit/5acec961f8a86068b7662410349eeac6ba98a399))
* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to b0832be ([#976](https://github.com/cloudquery/plugin-sdk/issues/976)) ([3d95166](https://github.com/cloudquery/plugin-sdk/commit/3d951665407e3287cd236b267c1e8e018e7992ef))
* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to d01ed41 ([#975](https://github.com/cloudquery/plugin-sdk/issues/975)) ([19dae31](https://github.com/cloudquery/plugin-sdk/commit/19dae31535cc80cae0800874a4f35c8dd83d7fb8))
* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to d864719 ([#972](https://github.com/cloudquery/plugin-sdk/issues/972)) ([9e25cc4](https://github.com/cloudquery/plugin-sdk/commit/9e25cc480060ed2cdfa48f174273d45c31c1ce07))
* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to f060192 ([#989](https://github.com/cloudquery/plugin-sdk/issues/989)) ([47c4fce](https://github.com/cloudquery/plugin-sdk/commit/47c4fceb39605c7f56bcb00e3fb2fd45da69304b))
* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to f0dffc6 ([#979](https://github.com/cloudquery/plugin-sdk/issues/979)) ([3579590](https://github.com/cloudquery/plugin-sdk/commit/357959068463cbc0bd33ce9ceefb1e0d59149a51))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.1.0 ([#977](https://github.com/cloudquery/plugin-sdk/issues/977)) ([e0f8009](https://github.com/cloudquery/plugin-sdk/commit/e0f8009c5de9d88b27ec2c15c5c84d3236fa924d))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.2.0 ([#983](https://github.com/cloudquery/plugin-sdk/issues/983)) ([8ce6e06](https://github.com/cloudquery/plugin-sdk/commit/8ce6e062824625fecf6b149f66c05eb2a86566de))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.2.1 ([#985](https://github.com/cloudquery/plugin-sdk/issues/985)) ([ade3b63](https://github.com/cloudquery/plugin-sdk/commit/ade3b63469c66e4d8fb5ce4f4719435da00a5789))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.3.0 ([#987](https://github.com/cloudquery/plugin-sdk/issues/987)) ([e1a2aec](https://github.com/cloudquery/plugin-sdk/commit/e1a2aec521483b3f890e28be651281ee4512e92e))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.3.2 ([#988](https://github.com/cloudquery/plugin-sdk/issues/988)) ([28076a7](https://github.com/cloudquery/plugin-sdk/commit/28076a7f168ca85f7c8fd39f347bf67a13c423c5))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.3.3 ([#990](https://github.com/cloudquery/plugin-sdk/issues/990)) ([1f5e87c](https://github.com/cloudquery/plugin-sdk/commit/1f5e87c3a81a3900c82348edfb776bf743ca1773))


### Miscellaneous Chores

* release 4.0.0-rc1 ([21e11bf](https://github.com/cloudquery/plugin-sdk/commit/21e11bf785fb904be6b0bf5bab480ea970b8c20d))

## [3.10.6](https://github.com/cloudquery/plugin-sdk/compare/v3.10.5...v3.10.6) (2023-06-13)


### Bug Fixes

* Don't write last batch in managed writer if the context was canceled ([#964](https://github.com/cloudquery/plugin-sdk/issues/964)) ([8027e62](https://github.com/cloudquery/plugin-sdk/commit/8027e62b66c0acf799d795a479e072f86a6dc205))

## [3.10.5](https://github.com/cloudquery/plugin-sdk/compare/v3.10.4...v3.10.5) (2023-06-13)


### Bug Fixes

* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to 0f7bd3b ([#961](https://github.com/cloudquery/plugin-sdk/issues/961)) ([21f3b68](https://github.com/cloudquery/plugin-sdk/commit/21f3b68d45d79e9e726cbe395044a3560145003d))
* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to 6b7fa9c ([#962](https://github.com/cloudquery/plugin-sdk/issues/962)) ([78eecf2](https://github.com/cloudquery/plugin-sdk/commit/78eecf2f4cf4027c7c37c8297ff7845debf49fd8))
* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to 71dfe94 ([#953](https://github.com/cloudquery/plugin-sdk/issues/953)) ([b48ae1a](https://github.com/cloudquery/plugin-sdk/commit/b48ae1a546a5f4e73e88793da9afda3b91f3ba08))
* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to 7f6aaff ([#963](https://github.com/cloudquery/plugin-sdk/issues/963)) ([8c7acdd](https://github.com/cloudquery/plugin-sdk/commit/8c7acdd63318bdc270b6fd1141db14148f7ba68c))
* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to 8f72077 ([#958](https://github.com/cloudquery/plugin-sdk/issues/958)) ([6f6c993](https://github.com/cloudquery/plugin-sdk/commit/6f6c9936e24f9460c253a297df44415dd4eef64f))
* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to 90670b8 ([#955](https://github.com/cloudquery/plugin-sdk/issues/955)) ([047ab30](https://github.com/cloudquery/plugin-sdk/commit/047ab3066139a00ad665b8ff766bd7f57e70803f))
* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to b359e74 ([#960](https://github.com/cloudquery/plugin-sdk/issues/960)) ([7e95e7d](https://github.com/cloudquery/plugin-sdk/commit/7e95e7dbadd1b12772eadc8607f267310ca5583e))
* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to d8eacf8 ([#966](https://github.com/cloudquery/plugin-sdk/issues/966)) ([2d32679](https://github.com/cloudquery/plugin-sdk/commit/2d3267979a6c66c7fd89a284e2d369b10338a7af))
* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to e258cfb ([#957](https://github.com/cloudquery/plugin-sdk/issues/957)) ([df842e0](https://github.com/cloudquery/plugin-sdk/commit/df842e01437f51b835f83683001c4fb15fc36b7a))
* **transformers:** Ability to transform `any` with TypeTransformer ([#956](https://github.com/cloudquery/plugin-sdk/issues/956)) ([c989c28](https://github.com/cloudquery/plugin-sdk/commit/c989c288ab2fa6f34f6dd71ed7a8fc4597db085e))

## [3.10.4](https://github.com/cloudquery/plugin-sdk/compare/v3.10.3...v3.10.4) (2023-06-06)


### Bug Fixes

* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to 88d5dc2 ([#950](https://github.com/cloudquery/plugin-sdk/issues/950)) ([58bfa32](https://github.com/cloudquery/plugin-sdk/commit/58bfa32767d8fa690b61a263091c802e4bd246a8))
* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.0.9 ([#952](https://github.com/cloudquery/plugin-sdk/issues/952)) ([3266266](https://github.com/cloudquery/plugin-sdk/commit/3266266fb011c04933bd2f08075458e9f3f23ccf))

## [3.10.3](https://github.com/cloudquery/plugin-sdk/compare/v3.10.2...v3.10.3) (2023-06-05)


### Bug Fixes

* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to 20b0de9 ([#947](https://github.com/cloudquery/plugin-sdk/issues/947)) ([32a0c05](https://github.com/cloudquery/plugin-sdk/commit/32a0c053deae5dad2ffc5a6c61932d573cf2b5a6))
* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to 6d34568 ([#944](https://github.com/cloudquery/plugin-sdk/issues/944)) ([f92fd66](https://github.com/cloudquery/plugin-sdk/commit/f92fd66282f258888787aa0924263052fb6315d3))
* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to c655015 ([#946](https://github.com/cloudquery/plugin-sdk/issues/946)) ([4b6e3a3](https://github.com/cloudquery/plugin-sdk/commit/4b6e3a33bad11fb97422887c7f0c0e50f9e00e41))
* **types:** Extensions conversion with storage ([#948](https://github.com/cloudquery/plugin-sdk/issues/948)) ([1132c02](https://github.com/cloudquery/plugin-sdk/commit/1132c0227cc8db45d3eefa147aebbe9e32941ec4))

## [3.10.2](https://github.com/cloudquery/plugin-sdk/compare/v3.10.1...v3.10.2) (2023-06-02)


### Bug Fixes

* Remove uint validation ([#942](https://github.com/cloudquery/plugin-sdk/issues/942)) ([4df3b46](https://github.com/cloudquery/plugin-sdk/commit/4df3b46b9180a415bc42b40648653e5dd8ba84fd))

## [3.10.1](https://github.com/cloudquery/plugin-sdk/compare/v3.10.0...v3.10.1) (2023-06-02)


### Bug Fixes

* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to a7aad4c ([#941](https://github.com/cloudquery/plugin-sdk/issues/941)) ([a39f6e8](https://github.com/cloudquery/plugin-sdk/commit/a39f6e871bcb038c2cd90a8f01ebcc0cdf02b1e8))
* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to ac40107 ([#939](https://github.com/cloudquery/plugin-sdk/issues/939)) ([ef9e774](https://github.com/cloudquery/plugin-sdk/commit/ef9e7747e360eee2e61abd42c02c0d668d896e1e))

## [3.10.0](https://github.com/cloudquery/plugin-sdk/compare/v3.9.0...v3.10.0) (2023-06-01)


### Features

* **scalar:** Support all int variations in decimal scalar ([#937](https://github.com/cloudquery/plugin-sdk/issues/937)) ([159e975](https://github.com/cloudquery/plugin-sdk/commit/159e975b3bd4f74925760507f7115f1880d19f21))
* **scalar:** Support pointer dereferencing in decimal ([#938](https://github.com/cloudquery/plugin-sdk/issues/938)) ([181e676](https://github.com/cloudquery/plugin-sdk/commit/181e6765bf8d8c22b1ccc276dc3258d1f25eeec3))


### Bug Fixes

* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to 7f8dd24 ([#936](https://github.com/cloudquery/plugin-sdk/issues/936)) ([8cfc215](https://github.com/cloudquery/plugin-sdk/commit/8cfc2151893bf6c175ea517e1879e33063a261dc))
* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to c1359c5 ([#933](https://github.com/cloudquery/plugin-sdk/issues/933)) ([dab8d86](https://github.com/cloudquery/plugin-sdk/commit/dab8d86804cb47dc0a8ef6244b91763306f456bc))
* **deps:** Update github.com/cloudquery/arrow/go/v13 digest to c67fb39 ([#935](https://github.com/cloudquery/plugin-sdk/issues/935)) ([82f5f60](https://github.com/cloudquery/plugin-sdk/commit/82f5f60bb010e1430d3e7f0303c398c69a3ce301))

## [3.9.0](https://github.com/cloudquery/plugin-sdk/compare/v3.8.1...v3.9.0) (2023-06-01)


### Features

* More scalars ([#914](https://github.com/cloudquery/plugin-sdk/issues/914)) ([f8625e2](https://github.com/cloudquery/plugin-sdk/commit/f8625e25ed202711c16343799ad72a48232f1e5c))


### Bug Fixes

* **scalar:** Handle nil pointer to []byte in uuid and binary ([#922](https://github.com/cloudquery/plugin-sdk/issues/922)) ([dac967a](https://github.com/cloudquery/plugin-sdk/commit/dac967a57b36856d51ddaa6c2c71744cbf43e18a))
* **testdata:** Match map field names with type ([#930](https://github.com/cloudquery/plugin-sdk/issues/930)) ([cec067d](https://github.com/cloudquery/plugin-sdk/commit/cec067d4902c8590f8295d5b97b0683a73d28d3c))

## [3.8.1](https://github.com/cloudquery/plugin-sdk/compare/v3.8.0...v3.8.1) (2023-06-01)


### Bug Fixes

* **deps:** Update github.com/apache/arrow/go/v13 digest to cbc17a9 ([#924](https://github.com/cloudquery/plugin-sdk/issues/924)) ([dd0789e](https://github.com/cloudquery/plugin-sdk/commit/dd0789e1ca0cfd8dc8d458e234cecc210c29929f))
* **deps:** Update golang.org/x/exp digest to 2e198f4 ([#926](https://github.com/cloudquery/plugin-sdk/issues/926)) ([97440df](https://github.com/cloudquery/plugin-sdk/commit/97440df046469c800a8cc1c5e49956484fb809ac))
* **deps:** Update google.golang.org/genproto digest to e85fd2c ([#927](https://github.com/cloudquery/plugin-sdk/issues/927)) ([b185a17](https://github.com/cloudquery/plugin-sdk/commit/b185a17ede3866754782bda305ef7102abf1b565))
* **deps:** Update google.golang.org/genproto/googleapis/rpc digest to e85fd2c ([#928](https://github.com/cloudquery/plugin-sdk/issues/928)) ([c23f09d](https://github.com/cloudquery/plugin-sdk/commit/c23f09dc406eb71df10c8f6d04c1e518cb85311c))
* **test:** Use `array.WithUnorderedMapKeys` ([#921](https://github.com/cloudquery/plugin-sdk/issues/921)) ([ac2cfbd](https://github.com/cloudquery/plugin-sdk/commit/ac2cfbdc09521ae78d648fe841354351348496cb))

## [3.8.0](https://github.com/cloudquery/plugin-sdk/compare/v3.7.0...v3.8.0) (2023-05-31)


### Features

* Add the names of tables to the periodic logger ([#738](https://github.com/cloudquery/plugin-sdk/issues/738)) ([72e1d49](https://github.com/cloudquery/plugin-sdk/commit/72e1d496cbed1e76c273ac5592419ac136c6ab2a))
* Separate Queued Tables from In Progress Tables ([#920](https://github.com/cloudquery/plugin-sdk/issues/920)) ([dcb5d26](https://github.com/cloudquery/plugin-sdk/commit/dcb5d26b3ee22de436327b9d9c7f0c514abf1ada))

## [3.7.0](https://github.com/cloudquery/plugin-sdk/compare/v3.6.7...v3.7.0) (2023-05-30)


### Features

* **test:** Add `AllowNull` option for test data ([#913](https://github.com/cloudquery/plugin-sdk/issues/913)) ([9b911eb](https://github.com/cloudquery/plugin-sdk/commit/9b911eb7ea5566a8a5979443bea21a45779b4691))


### Bug Fixes

* Test Decimal type, map type and larger number ranges ([#905](https://github.com/cloudquery/plugin-sdk/issues/905)) ([9a3b4ad](https://github.com/cloudquery/plugin-sdk/commit/9a3b4ad3380f95ae6eabb59203d2a608e80ef59e))

## [3.6.7](https://github.com/cloudquery/plugin-sdk/compare/v3.6.6...v3.6.7) (2023-05-26)


### Bug Fixes

* Update Arrow to latest cqmain branch ([#910](https://github.com/cloudquery/plugin-sdk/issues/910)) ([1295559](https://github.com/cloudquery/plugin-sdk/commit/12955593507984fa51c1130732a34df1b256d800))

## [3.6.6](https://github.com/cloudquery/plugin-sdk/compare/v3.6.5...v3.6.6) (2023-05-26)


### Bug Fixes

* Use backtick around types ([#908](https://github.com/cloudquery/plugin-sdk/issues/908)) ([858fe54](https://github.com/cloudquery/plugin-sdk/commit/858fe5429bf17ab32a07957a1a60433a8780ace5))

## [3.6.5](https://github.com/cloudquery/plugin-sdk/compare/v3.6.4...v3.6.5) (2023-05-26)


### Bug Fixes

* Transform `[]any` as `JSON` ([#906](https://github.com/cloudquery/plugin-sdk/issues/906)) ([7719677](https://github.com/cloudquery/plugin-sdk/commit/771967717617e40ef809882dbdaed83d6bfad116))

## [3.6.4](https://github.com/cloudquery/plugin-sdk/compare/v3.6.3...v3.6.4) (2023-05-25)


### Bug Fixes

* Scalar set now accepts scalar type ([#902](https://github.com/cloudquery/plugin-sdk/issues/902)) ([1ff2229](https://github.com/cloudquery/plugin-sdk/commit/1ff222910356762ea2c7f48c4bc2ee3c19769e26))

## [3.6.3](https://github.com/cloudquery/plugin-sdk/compare/v3.6.2...v3.6.3) (2023-05-24)


### Bug Fixes

* Better handling for Arrow type strings in docs ([#896](https://github.com/cloudquery/plugin-sdk/issues/896)) ([78699f4](https://github.com/cloudquery/plugin-sdk/commit/78699f416c67fb701eb7f7d56a5beba37b3fc150))

## [3.6.2](https://github.com/cloudquery/plugin-sdk/compare/v3.6.1...v3.6.2) (2023-05-22)


### Bug Fixes

* **testdata:** Don't use escaping in JSON testdata (as array.Approx will check the underlying data) ([#898](https://github.com/cloudquery/plugin-sdk/issues/898)) ([f7e0ae7](https://github.com/cloudquery/plugin-sdk/commit/f7e0ae7bbf520a77d3a900fc9b0068a18fcdfab3))

## [3.6.1](https://github.com/cloudquery/plugin-sdk/compare/v3.6.0...v3.6.1) (2023-05-21)


### Bug Fixes

* Inet extension MarshalJSON ([#894](https://github.com/cloudquery/plugin-sdk/issues/894)) ([f483c57](https://github.com/cloudquery/plugin-sdk/commit/f483c572ac2b77f42a8f3a6cf8a0327fae3fce4c))

## [3.6.0](https://github.com/cloudquery/plugin-sdk/compare/v3.5.2...v3.6.0) (2023-05-21)


### Features

* Add precision options for dest testing ([#893](https://github.com/cloudquery/plugin-sdk/issues/893)) ([faacca6](https://github.com/cloudquery/plugin-sdk/commit/faacca6b52347b9cf61b0acbcb4096f535817087))
* Refactor test options and allow skipping of nulls in lists ([#892](https://github.com/cloudquery/plugin-sdk/issues/892)) ([bc3c251](https://github.com/cloudquery/plugin-sdk/commit/bc3c25193c6675317835a9642758c350260486e9))


### Bug Fixes

* Add null-row case for append-only tests ([#889](https://github.com/cloudquery/plugin-sdk/issues/889)) ([6967929](https://github.com/cloudquery/plugin-sdk/commit/6967929bc598ddc2bf6120a9a905ccbf92b97773))
* Tighter Arrow test cases ([#891](https://github.com/cloudquery/plugin-sdk/issues/891)) ([c7f2546](https://github.com/cloudquery/plugin-sdk/commit/c7f25468f5fff7176cc71301d337598837ef7d61))

## [3.5.2](https://github.com/cloudquery/plugin-sdk/compare/v3.5.1...v3.5.2) (2023-05-18)


### Bug Fixes

* **arrow:** `schema.Table` &lt;-&gt; `arrow.Schema` conversion ([#886](https://github.com/cloudquery/plugin-sdk/issues/886)) ([61d98c9](https://github.com/cloudquery/plugin-sdk/commit/61d98c9558287879137e10da0687bbf307d0d0ac))
* **destination:** Don't duplicate tables to be removed ([#886](https://github.com/cloudquery/plugin-sdk/issues/886)) ([61d98c9](https://github.com/cloudquery/plugin-sdk/commit/61d98c9558287879137e10da0687bbf307d0d0ac))
* **tables:** Flatten stripping relations ([#884](https://github.com/cloudquery/plugin-sdk/issues/884)) ([e890385](https://github.com/cloudquery/plugin-sdk/commit/e890385102e2668a16e35cff75fe2ffea32f2937))
* **testing:** CQ Parent ID column should not be NotNull ([#887](https://github.com/cloudquery/plugin-sdk/issues/887)) ([f4aa5bc](https://github.com/cloudquery/plugin-sdk/commit/f4aa5bcebc88ae1a9a5bd90937dcd5868dc0dff1))

## [3.5.1](https://github.com/cloudquery/plugin-sdk/compare/v3.5.0...v3.5.1) (2023-05-16)


### Bug Fixes

* Flatten V2 tables ([#882](https://github.com/cloudquery/plugin-sdk/issues/882)) ([28706f1](https://github.com/cloudquery/plugin-sdk/commit/28706f17eb3cc9d0766ecd9c3554eb7505d69c85))

## [3.5.0](https://github.com/cloudquery/plugin-sdk/compare/v3.4.0...v3.5.0) (2023-05-16)


### Features

* Revert "feat(test): Test writing to a child table" ([#880](https://github.com/cloudquery/plugin-sdk/issues/880)) ([9d61013](https://github.com/cloudquery/plugin-sdk/commit/9d610131faf4597fe191caac08d40a93efd8aafe))

## [3.4.0](https://github.com/cloudquery/plugin-sdk/compare/v3.3.0...v3.4.0) (2023-05-16)


### Features

* **test:** Test writing to a child table ([#878](https://github.com/cloudquery/plugin-sdk/issues/878)) ([d4154fb](https://github.com/cloudquery/plugin-sdk/commit/d4154fb4e2bc703d2974afa4e7dd9c2c774940f9)), closes [#877](https://github.com/cloudquery/plugin-sdk/issues/877)


### Bug Fixes

* **test:** Remove extra `v2/schema` import ([#876](https://github.com/cloudquery/plugin-sdk/issues/876)) ([da9ed4d](https://github.com/cloudquery/plugin-sdk/commit/da9ed4d79223ab2c21b48e816ebd194b9b42b262))

## [3.3.0](https://github.com/cloudquery/plugin-sdk/compare/v3.2.1...v3.3.0) (2023-05-15)


### Features

* Support sources in SDK V3 ([#864](https://github.com/cloudquery/plugin-sdk/issues/864)) ([a49abcb](https://github.com/cloudquery/plugin-sdk/commit/a49abcbc67e695d804b72baee1bb8813d3830a4a))


### Bug Fixes

* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.0.8 ([#874](https://github.com/cloudquery/plugin-sdk/issues/874)) ([56c0e84](https://github.com/cloudquery/plugin-sdk/commit/56c0e8451606aa2ee9e8773e640bbf339037629d))

## [3.2.1](https://github.com/cloudquery/plugin-sdk/compare/v3.2.0...v3.2.1) (2023-05-15)


### Bug Fixes

* Fix test column generation ([#872](https://github.com/cloudquery/plugin-sdk/issues/872)) ([99fb000](https://github.com/cloudquery/plugin-sdk/commit/99fb0008d216c7b63ccf91db90e99da996185c46))

## [3.2.0](https://github.com/cloudquery/plugin-sdk/compare/v3.1.0...v3.2.0) (2023-05-15)


### Features

* Allow testing of more Arrow types ([#863](https://github.com/cloudquery/plugin-sdk/issues/863)) ([28642ec](https://github.com/cloudquery/plugin-sdk/commit/28642ec7537ac9f1b97401a66e1982591b62b6d9))

## [3.1.0](https://github.com/cloudquery/plugin-sdk/compare/v3.0.1...v3.1.0) (2023-05-15)


### Features

* **schema:** Embed column creation options ([#869](https://github.com/cloudquery/plugin-sdk/issues/869)) ([7512e29](https://github.com/cloudquery/plugin-sdk/commit/7512e299168e43fb1d8b9d184d71a2b23f1d9892))
* **types:** Rename Mac -&gt; MAC ([#868](https://github.com/cloudquery/plugin-sdk/issues/868)) ([b5c76bb](https://github.com/cloudquery/plugin-sdk/commit/b5c76bb36b52c01bd27ec8529529dc69ecf0f116))


### Bug Fixes

* **deps:** Update module github.com/cloudquery/plugin-pb-go to v1.0.6 ([#865](https://github.com/cloudquery/plugin-sdk/issues/865)) ([1fb4eaf](https://github.com/cloudquery/plugin-sdk/commit/1fb4eafe3e3b0842b071948c3f2e3dd5d29dad22))

## [3.0.1](https://github.com/cloudquery/plugin-sdk/compare/v3.0.0...v3.0.1) (2023-05-11)


### Bug Fixes

* **testing:** Wrong types in v3 testdata ([#859](https://github.com/cloudquery/plugin-sdk/issues/859)) ([e494fb5](https://github.com/cloudquery/plugin-sdk/commit/e494fb51f177ea0ae9af735e9fb8f320c3a72b94))

## [3.0.0](https://github.com/cloudquery/plugin-sdk/compare/v2.7.0...v3.0.0) (2023-05-09)


### ⚠ BREAKING CHANGES

* Upgrade to SDK V3 make Column.Type an arrow.DataType ([#854](https://github.com/cloudquery/plugin-sdk/issues/854))

### Features

* Upgrade to SDK V3 make Column.Type an arrow.DataType ([#854](https://github.com/cloudquery/plugin-sdk/issues/854)) ([1265554](https://github.com/cloudquery/plugin-sdk/commit/12655541d1b7e4a1c5ab69e3c9e16f3978d2d44e))

## [2.7.0](https://github.com/cloudquery/plugin-sdk/compare/v2.6.0...v2.7.0) (2023-05-09)


### Features

* **deps:** Upgrade to Apache Arrow v13 (latest `cqmain`) ([#852](https://github.com/cloudquery/plugin-sdk/issues/852)) ([5ae502f](https://github.com/cloudquery/plugin-sdk/commit/5ae502f7fe6c41043f1a5e1392c69657d8d9062e))

## [2.6.0](https://github.com/cloudquery/plugin-sdk/compare/v2.5.4...v2.6.0) (2023-05-08)


### Features

* **arrow:** Add `types.XBuilder.NewXArray` helpers ([2df4413](https://github.com/cloudquery/plugin-sdk/commit/2df4413bed3df91ec596e2540584debab1974f4a))
* Move proto to external repository ([#844](https://github.com/cloudquery/plugin-sdk/issues/844)) ([3cd3ba7](https://github.com/cloudquery/plugin-sdk/commit/3cd3ba7d910141ba89265767d968d24516809332))

## [2.5.4](https://github.com/cloudquery/plugin-sdk/compare/v2.5.3...v2.5.4) (2023-05-05)


### Bug Fixes

* **arrow:** Allow empty and `nil` valid param in `AppendValues` ([#847](https://github.com/cloudquery/plugin-sdk/issues/847)) ([dafd05b](https://github.com/cloudquery/plugin-sdk/commit/dafd05b3e2b8dc406d4b6a4bdaf6d1143e569f1d))

## [2.5.3](https://github.com/cloudquery/plugin-sdk/compare/v2.5.2...v2.5.3) (2023-05-04)


### Bug Fixes

* **arrow:** Add missing table options ([#833](https://github.com/cloudquery/plugin-sdk/issues/833)) ([95a9f0c](https://github.com/cloudquery/plugin-sdk/commit/95a9f0c29c6c2b85fded012341bf00cff0225605))

## [2.5.2](https://github.com/cloudquery/plugin-sdk/compare/v2.5.1...v2.5.2) (2023-05-02)


### Bug Fixes

* **deps:** Update github.com/apache/arrow/go/v12 digest to 0ea1a10 ([#836](https://github.com/cloudquery/plugin-sdk/issues/836)) ([5561fa1](https://github.com/cloudquery/plugin-sdk/commit/5561fa1a59ee498d5ecb0acbde79971e82fe4fda))
* **deps:** Update golang.org/x/exp digest to 47ecfdc ([#837](https://github.com/cloudquery/plugin-sdk/issues/837)) ([bb56f9c](https://github.com/cloudquery/plugin-sdk/commit/bb56f9c67d1ce5936c32c093911b915680707954))
* **deps:** Update golang.org/x/xerrors digest to 04be3eb ([#838](https://github.com/cloudquery/plugin-sdk/issues/838)) ([42d4517](https://github.com/cloudquery/plugin-sdk/commit/42d4517d223791f75881ad301d6df90664d4e232))
* **deps:** Update google.golang.org/genproto digest to daa745c ([#839](https://github.com/cloudquery/plugin-sdk/issues/839)) ([1285222](https://github.com/cloudquery/plugin-sdk/commit/128522279101eb316f3b29665a1f3c7c65da1e3e))
* **deps:** Update module github.com/avast/retry-go/v4 to v4.3.4 ([#840](https://github.com/cloudquery/plugin-sdk/issues/840)) ([47da73d](https://github.com/cloudquery/plugin-sdk/commit/47da73dac6c2af71e13d65e9b872fd0657cb0a2a))
* Destination migration testing using incorrect mode ([#822](https://github.com/cloudquery/plugin-sdk/issues/822)) ([fa51c80](https://github.com/cloudquery/plugin-sdk/commit/fa51c80522b2bf573414eae81f12cd21b1cf549f))
* **json:** Use `GetOneForMarshal` instead of deserialization-serialization cycle ([#834](https://github.com/cloudquery/plugin-sdk/issues/834)) ([6fb7c1c](https://github.com/cloudquery/plugin-sdk/commit/6fb7c1c761a0ed49f84f61afaadcc958966e58fa))

## [2.5.1](https://github.com/cloudquery/plugin-sdk/compare/v2.5.0...v2.5.1) (2023-04-28)


### Bug Fixes

* **transformer:** Allow camel-cased json tags ([#828](https://github.com/cloudquery/plugin-sdk/issues/828)) ([653a50d](https://github.com/cloudquery/plugin-sdk/commit/653a50dccd9456f5e676a1fb63b8ff37fd5cc4e8))

## [2.5.0](https://github.com/cloudquery/plugin-sdk/compare/v2.4.0...v2.5.0) (2023-04-28)


### Features

* Add table description to Arrow schema metadata ([#824](https://github.com/cloudquery/plugin-sdk/issues/824)) ([1a8072f](https://github.com/cloudquery/plugin-sdk/commit/1a8072ff7eff1c411569a538958069ad0744a0ce))
* **arrow:** Streamline Apache Arrow extension types ([#823](https://github.com/cloudquery/plugin-sdk/issues/823)) ([f32fac3](https://github.com/cloudquery/plugin-sdk/commit/f32fac3b04c769bb86774c3d1b89991d5d2f51b3))
* **test:** Add double migration test ([#827](https://github.com/cloudquery/plugin-sdk/issues/827)) ([4cd3872](https://github.com/cloudquery/plugin-sdk/commit/4cd3872f2a281c6b7e685d13061d6b7849fff3f4))
* Time values are truncated uniformly ([#825](https://github.com/cloudquery/plugin-sdk/issues/825)) ([ffb97b0](https://github.com/cloudquery/plugin-sdk/commit/ffb97b0ddc949edccb2f05a4b67f3bc6b3ca2401))


### Bug Fixes

* TransformWithStruct/DefaultNameTransformer change for invalid column names ([#820](https://github.com/cloudquery/plugin-sdk/issues/820)) ([01e6649](https://github.com/cloudquery/plugin-sdk/commit/01e66491f6a21b1ed8fe1837ac86c0cccafd0cab))

## [2.4.0](https://github.com/cloudquery/plugin-sdk/compare/v2.3.8...v2.4.0) (2023-04-24)


### Features

* **arrow:** Pretty-print field changes ([#817](https://github.com/cloudquery/plugin-sdk/issues/817)) ([6c0d0b3](https://github.com/cloudquery/plugin-sdk/commit/6c0d0b346a2748dbac2464b81dfab86d307e6090))

## [2.3.8](https://github.com/cloudquery/plugin-sdk/compare/v2.3.7...v2.3.8) (2023-04-20)


### Bug Fixes

* Fail on empty tables ([#796](https://github.com/cloudquery/plugin-sdk/issues/796)) ([1320d32](https://github.com/cloudquery/plugin-sdk/commit/1320d32b5a2e6ea7b6bacb0b597caf45c3f26b1e))
* **testing:** Add sorting for testing dest migrations ([#814](https://github.com/cloudquery/plugin-sdk/issues/814)) ([b1437f1](https://github.com/cloudquery/plugin-sdk/commit/b1437f1fd7a67253f6d1fc68bbb713fedbbb91c2))

## [2.3.7](https://github.com/cloudquery/plugin-sdk/compare/v2.3.6...v2.3.7) (2023-04-20)


### Bug Fixes

* Use Go memory allocator for arrow ([#810](https://github.com/cloudquery/plugin-sdk/issues/810)) ([b54e5e1](https://github.com/cloudquery/plugin-sdk/commit/b54e5e16378de6dc08d6782769f1779acb92804e))

## [2.3.6](https://github.com/cloudquery/plugin-sdk/compare/v2.3.5...v2.3.6) (2023-04-19)


### Bug Fixes

* Release resource on SkipSecondAppend ([#808](https://github.com/cloudquery/plugin-sdk/issues/808)) ([6f19c2d](https://github.com/cloudquery/plugin-sdk/commit/6f19c2d69f33b9983ffe4c201058db33e97a4e13))
* **testdata:** Add old style gen testdata ([#811](https://github.com/cloudquery/plugin-sdk/issues/811)) ([494992b](https://github.com/cloudquery/plugin-sdk/commit/494992b267b3c145e63e1c97912d56bcc50da13f))

## [2.3.5](https://github.com/cloudquery/plugin-sdk/compare/v2.3.4...v2.3.5) (2023-04-19)


### Bug Fixes

* Truncate timestamp to millisecond in dest testing ([#806](https://github.com/cloudquery/plugin-sdk/issues/806)) ([eb8b7c4](https://github.com/cloudquery/plugin-sdk/commit/eb8b7c49cf788ebb8702d48cf22e75c6b56b8856))

## [2.3.4](https://github.com/cloudquery/plugin-sdk/compare/v2.3.3...v2.3.4) (2023-04-19)


### Bug Fixes

* Undo release of all resources in managed writer ([#801](https://github.com/cloudquery/plugin-sdk/issues/801)) ([d586be0](https://github.com/cloudquery/plugin-sdk/commit/d586be077b099fa6d00e405a3b6c0bd655c1b40c))

## [2.3.3](https://github.com/cloudquery/plugin-sdk/compare/v2.3.2...v2.3.3) (2023-04-19)


### Bug Fixes

* Make cq_id non required on destination ([#799](https://github.com/cloudquery/plugin-sdk/issues/799)) ([7f33b8d](https://github.com/cloudquery/plugin-sdk/commit/7f33b8df0e283fb8db5e70744a9964671f6b53d4))

## [2.3.2](https://github.com/cloudquery/plugin-sdk/compare/v2.3.1...v2.3.2) (2023-04-19)


### Bug Fixes

* Arrow Retain and Release fixes ([#795](https://github.com/cloudquery/plugin-sdk/issues/795)) ([a893db6](https://github.com/cloudquery/plugin-sdk/commit/a893db675c5f4bb8cab71a854014c65caa43d3e3))
* Disallow null character in strings per utf8 spec ([#797](https://github.com/cloudquery/plugin-sdk/issues/797)) ([591502f](https://github.com/cloudquery/plugin-sdk/commit/591502f51ea99ca852b307616e60ab665b231440))

## [2.3.1](https://github.com/cloudquery/plugin-sdk/compare/v2.3.0...v2.3.1) (2023-04-18)


### Bug Fixes

* Set _cq_id to NotNull in destinations for backward compat ([#793](https://github.com/cloudquery/plugin-sdk/issues/793)) ([1ab4350](https://github.com/cloudquery/plugin-sdk/commit/1ab4350c7b26993f71cb39adc0d9e6d3caeddb7a))

## [2.3.0](https://github.com/cloudquery/plugin-sdk/compare/v2.2.2...v2.3.0) (2023-04-18)


### Features

* Change default source tables to none ([#790](https://github.com/cloudquery/plugin-sdk/issues/790)) ([b33c777](https://github.com/cloudquery/plugin-sdk/commit/b33c77752a0b155c12ca46985410f56700a16589))


### Bug Fixes

* Update to latest Arrow (cqmain branch) ([#792](https://github.com/cloudquery/plugin-sdk/issues/792)) ([a6fdaca](https://github.com/cloudquery/plugin-sdk/commit/a6fdaca6656b79a6b420217abe8583be832ab70b))

## [2.2.2](https://github.com/cloudquery/plugin-sdk/compare/v2.2.1...v2.2.2) (2023-04-17)


### Bug Fixes

* Destination testing memory leak ([#788](https://github.com/cloudquery/plugin-sdk/issues/788)) ([c17b64d](https://github.com/cloudquery/plugin-sdk/commit/c17b64dade247d794bd191075518eeba30d03a96))

## [2.2.1](https://github.com/cloudquery/plugin-sdk/compare/v2.2.0...v2.2.1) (2023-04-17)


### Bug Fixes

* Make cq_id unique for backward compat ([#786](https://github.com/cloudquery/plugin-sdk/issues/786)) ([ad25ded](https://github.com/cloudquery/plugin-sdk/commit/ad25dedf81d0fb8538cd34dd0998627887ad5300))

## [2.2.0](https://github.com/cloudquery/plugin-sdk/compare/v2.1.0...v2.2.0) (2023-04-17)


### Features

* Use ApproxEqual in dest tests ([#784](https://github.com/cloudquery/plugin-sdk/issues/784)) ([88a677a](https://github.com/cloudquery/plugin-sdk/commit/88a677a059f24575a0019552da92827a440b6b47))


### Bug Fixes

* Add composite PK to test table ([#768](https://github.com/cloudquery/plugin-sdk/issues/768)) ([57b8edd](https://github.com/cloudquery/plugin-sdk/commit/57b8edd823df9f2f2b603f42f3a298edf2a22bef))
* Add StableTime to GenTestDataOptions and make panic message more verbose ([#783](https://github.com/cloudquery/plugin-sdk/issues/783)) ([be7a9a7](https://github.com/cloudquery/plugin-sdk/commit/be7a9a72b1317bb69c6e902d50f24705890a78c4))
* Handle When `_cq_id` only PK ([#774](https://github.com/cloudquery/plugin-sdk/issues/774)) ([06fde4b](https://github.com/cloudquery/plugin-sdk/commit/06fde4b0f4f4bf4bf07878f30d0cf6222e295642))

## [2.1.0](https://github.com/cloudquery/plugin-sdk/compare/v2.0.1...v2.1.0) (2023-04-12)


### Features

* **destination:** Remove redundant `ReverseTransformValues` method ([#778](https://github.com/cloudquery/plugin-sdk/issues/778)) ([bea4d00](https://github.com/cloudquery/plugin-sdk/commit/bea4d00d6502a0a131abb2321685733af8de62c1))


### Bug Fixes

* **unimplemented:** Conform to the interface ([#777](https://github.com/cloudquery/plugin-sdk/issues/777)) ([3a155d4](https://github.com/cloudquery/plugin-sdk/commit/3a155d4997cd76fe4459c779eaaad0d9dc47f8c6))

## [2.0.1](https://github.com/cloudquery/plugin-sdk/compare/v2.0.0...v2.0.1) (2023-04-11)


### Bug Fixes

* Update custom types with ValueStr and AppendFromValueString ([#772](https://github.com/cloudquery/plugin-sdk/issues/772)) ([166198e](https://github.com/cloudquery/plugin-sdk/commit/166198e8af595307adaa2ffe8577da5bde4fb1fa))

## [2.0.0](https://github.com/cloudquery/plugin-sdk/compare/v1.44.2...v2.0.0) (2023-04-11)


### ⚠ BREAKING CHANGES

* Arrow migration for destination

### Features

* Arrow migration for destination ([b39da64](https://github.com/cloudquery/plugin-sdk/commit/b39da6418115d7cf07902f7391de3565fcbbda0d))


### Bug Fixes

* **deps:** Update module golang.org/x/net to v0.9.0 ([#752](https://github.com/cloudquery/plugin-sdk/issues/752)) ([336a957](https://github.com/cloudquery/plugin-sdk/commit/336a957984ea12088a1783c2ce030dc148473287))
* **deps:** Update module golang.org/x/sys to v0.7.0 ([#753](https://github.com/cloudquery/plugin-sdk/issues/753)) ([8d88a50](https://github.com/cloudquery/plugin-sdk/commit/8d88a50d6a47eafeeb35610b75e61b39110de42f))
* **deps:** Update module golang.org/x/term to v0.7.0 ([#754](https://github.com/cloudquery/plugin-sdk/issues/754)) ([643d5e0](https://github.com/cloudquery/plugin-sdk/commit/643d5e0287ac62e497acfda46629a0bbfb03f5bf))
* **deps:** Update module golang.org/x/text to v0.9.0 ([#755](https://github.com/cloudquery/plugin-sdk/issues/755)) ([92d3748](https://github.com/cloudquery/plugin-sdk/commit/92d3748d239829ac843df379d5ee903865fc0543))

## [1.44.2](https://github.com/cloudquery/plugin-sdk/compare/v1.44.1...v1.44.2) (2023-04-04)


### Bug Fixes

* **deps:** Update golang.org/x/exp digest to 10a5072 ([#745](https://github.com/cloudquery/plugin-sdk/issues/745)) ([d52241c](https://github.com/cloudquery/plugin-sdk/commit/d52241c3120edf6b10cb9aacb6cee6ecad1b1764))
* **deps:** Update google.golang.org/genproto digest to dcfb400 ([#746](https://github.com/cloudquery/plugin-sdk/issues/746)) ([b510219](https://github.com/cloudquery/plugin-sdk/commit/b51021934cd6355e9040d00504404f909490278b))
* **deps:** Update module github.com/getsentry/sentry-go to v0.20.0 ([#751](https://github.com/cloudquery/plugin-sdk/issues/751)) ([47b2fbc](https://github.com/cloudquery/plugin-sdk/commit/47b2fbcaab138f1d725a697f48a7c74db609bc62))
* **deps:** Update module github.com/mattn/go-isatty to v0.0.18 ([#749](https://github.com/cloudquery/plugin-sdk/issues/749)) ([2d39af0](https://github.com/cloudquery/plugin-sdk/commit/2d39af0a6d4e71ae227c010f223889bde6157cf0))
* **deps:** Update module github.com/schollz/progressbar/v3 to v3.13.1 ([#750](https://github.com/cloudquery/plugin-sdk/issues/750)) ([ee3f17f](https://github.com/cloudquery/plugin-sdk/commit/ee3f17fc56e3ee05ce3389a38415bccb10b4d420))

## [1.44.1](https://github.com/cloudquery/plugin-sdk/compare/v1.44.0...v1.44.1) (2023-03-31)


### Bug Fixes

* **transform:** Use path instead of field name for PK options ([#739](https://github.com/cloudquery/plugin-sdk/issues/739)) ([d7649d8](https://github.com/cloudquery/plugin-sdk/commit/d7649d80f1a15cac6b7a29b6d0458a83db68cc76))

## [1.44.0](https://github.com/cloudquery/plugin-sdk/compare/v1.43.0...v1.44.0) (2023-03-17)


### Features

* Support for User Specifying Primary Key Scheme (default or cq-ids) ([#732](https://github.com/cloudquery/plugin-sdk/issues/732)) ([a41af50](https://github.com/cloudquery/plugin-sdk/commit/a41af50cca9529be77ca9c94114377dd0af006d6))

## [1.43.0](https://github.com/cloudquery/plugin-sdk/compare/v1.42.0...v1.43.0) (2023-03-14)


### Features

* Add ability to store table titles and render them in documentation ([#729](https://github.com/cloudquery/plugin-sdk/issues/729)) ([a0a58c4](https://github.com/cloudquery/plugin-sdk/commit/a0a58c4d523eee6d48e3500f3f8d1b571eef2a43))
* **source:** Expose docs generation ([#726](https://github.com/cloudquery/plugin-sdk/issues/726)) ([3360aa6](https://github.com/cloudquery/plugin-sdk/commit/3360aa6cbb9e7d383debc257a937fde0a58b4fa3))

## [1.42.0](https://github.com/cloudquery/plugin-sdk/compare/v1.41.0...v1.42.0) (2023-03-06)


### Features

* Add arrow support for timestamp and bytea ([#724](https://github.com/cloudquery/plugin-sdk/issues/724)) ([c2e84c3](https://github.com/cloudquery/plugin-sdk/commit/c2e84c369d3d7eb63fcf27de494078ee09125998))

## [1.41.0](https://github.com/cloudquery/plugin-sdk/compare/v1.40.0...v1.41.0) (2023-03-02)


### Features

* Deterministic _cq_id ([#712](https://github.com/cloudquery/plugin-sdk/issues/712)) ([2e7ad2c](https://github.com/cloudquery/plugin-sdk/commit/2e7ad2c03e9817ea00de31774a8869ef77b60325))
* **multiplex:** Detect duplicated clients ([#723](https://github.com/cloudquery/plugin-sdk/issues/723)) ([dfb039d](https://github.com/cloudquery/plugin-sdk/commit/dfb039d76c6976749c001bd7f12fcb32fa052e9d))


### Bug Fixes

* Cleanup code ([#710](https://github.com/cloudquery/plugin-sdk/issues/710)) ([963f03c](https://github.com/cloudquery/plugin-sdk/commit/963f03cd3d12a6ebdc091a5a555472abec858c00))
* **deps:** Update golang.org/x/exp digest to c95f2b4 ([#718](https://github.com/cloudquery/plugin-sdk/issues/718)) ([de52c10](https://github.com/cloudquery/plugin-sdk/commit/de52c10aa43132b2ceb08486722bb5fdd2acf8a1))
* **deps:** Update google.golang.org/genproto digest to 9b19f0b ([#719](https://github.com/cloudquery/plugin-sdk/issues/719)) ([ecfddea](https://github.com/cloudquery/plugin-sdk/commit/ecfddeaff6a6ffcc4cc9c454ae3906bd7e9e01f7))
* **deps:** Update module github.com/rivo/uniseg to v0.4.4 ([#720](https://github.com/cloudquery/plugin-sdk/issues/720)) ([0da69b6](https://github.com/cloudquery/plugin-sdk/commit/0da69b6a488fbbc6010cecea26522836a2ddba65))
* **deps:** Update module github.com/stretchr/testify to v1.8.2 ([#721](https://github.com/cloudquery/plugin-sdk/issues/721)) ([19c0742](https://github.com/cloudquery/plugin-sdk/commit/19c07425eb1c82a2ef962ed411742291557db2b8))
* **pk:** Skip filter for no PK ([#709](https://github.com/cloudquery/plugin-sdk/issues/709)) ([d0c2e26](https://github.com/cloudquery/plugin-sdk/commit/d0c2e2682b164707a0c15bfc5173ca7461cbf175))
* **types-json:** Disable HTML escaping during JSON marshalling ([#714](https://github.com/cloudquery/plugin-sdk/issues/714)) ([2f6f1d8](https://github.com/cloudquery/plugin-sdk/commit/2f6f1d8c65653d2816c263851b07fa455c3cb5d1))
* **types-timestamp:** Ensure timestamp is UTC ([#716](https://github.com/cloudquery/plugin-sdk/issues/716)) ([bb33629](https://github.com/cloudquery/plugin-sdk/commit/bb33629678bb01ee74da49296b0b14e024ce94af))

## [1.40.0](https://github.com/cloudquery/plugin-sdk/compare/v1.39.1...v1.40.0) (2023-02-23)


### Features

* **spec:** Return sources, destinations in order ([#624](https://github.com/cloudquery/plugin-sdk/issues/624)) ([4602071](https://github.com/cloudquery/plugin-sdk/commit/4602071ad83c16473a4afe899f384f5c94010252))

## [1.39.1](https://github.com/cloudquery/plugin-sdk/compare/v1.39.0...v1.39.1) (2023-02-22)


### Bug Fixes

* **destination:** Set CqID to unique at the destination level ([#704](https://github.com/cloudquery/plugin-sdk/issues/704)) ([1a97cb8](https://github.com/cloudquery/plugin-sdk/commit/1a97cb8d39c7236c72842f61f95ff514bc01cf11))

## [1.39.0](https://github.com/cloudquery/plugin-sdk/compare/v1.38.2...v1.39.0) (2023-02-21)


### Features

* **schema:** Add Unique column option, set it for CqID ([#702](https://github.com/cloudquery/plugin-sdk/issues/702)) ([d5c7636](https://github.com/cloudquery/plugin-sdk/commit/d5c763666c6e758fa39c26a362952a96de5105fa))

## [1.38.2](https://github.com/cloudquery/plugin-sdk/compare/v1.38.1...v1.38.2) (2023-02-20)


### Bug Fixes

* **test-migrate:** Add CqId to migrate tests tables ([#695](https://github.com/cloudquery/plugin-sdk/issues/695)) ([e996a11](https://github.com/cloudquery/plugin-sdk/commit/e996a11571d7039343a74e780024f60b79ca965c))
* **test-migrate:** Ignore order when comparing resources read ([#696](https://github.com/cloudquery/plugin-sdk/issues/696)) ([aea1b82](https://github.com/cloudquery/plugin-sdk/commit/aea1b82cf269a88b371bd81ee56523b79fbb5cdf))

## [1.38.1](https://github.com/cloudquery/plugin-sdk/compare/v1.38.0...v1.38.1) (2023-02-18)


### Bug Fixes

* **deps:** Update module golang.org/x/net to v0.7.0 [SECURITY] ([#692](https://github.com/cloudquery/plugin-sdk/issues/692)) ([47566c9](https://github.com/cloudquery/plugin-sdk/commit/47566c93f0ce88f6e76f1fcbe261ac14a56f77d3))

## [1.38.0](https://github.com/cloudquery/plugin-sdk/compare/v1.37.1...v1.38.0) (2023-02-16)


### Features

* Improve migration detection APIs ([#688](https://github.com/cloudquery/plugin-sdk/issues/688)) ([dc3bedf](https://github.com/cloudquery/plugin-sdk/commit/dc3bedf7af75c834882753a10499162da626a876))


### Bug Fixes

* Better string methods for TableColumnChange ([#690](https://github.com/cloudquery/plugin-sdk/issues/690)) ([a0ec52c](https://github.com/cloudquery/plugin-sdk/commit/a0ec52ca2c161cd6f77bca1285d47ae2d7616e30))

## [1.37.1](https://github.com/cloudquery/plugin-sdk/compare/v1.37.0...v1.37.1) (2023-02-14)


### Bug Fixes

* Set _cq_id not null for all tables ([#686](https://github.com/cloudquery/plugin-sdk/issues/686)) ([ff5f1d4](https://github.com/cloudquery/plugin-sdk/commit/ff5f1d423299a5bc44da635d26210ef088722234))

## [1.37.0](https://github.com/cloudquery/plugin-sdk/compare/v1.36.3...v1.37.0) (2023-02-13)


### Features

* Add unmanaged sources ([#677](https://github.com/cloudquery/plugin-sdk/issues/677)) ([f3e2b1d](https://github.com/cloudquery/plugin-sdk/commit/f3e2b1d982268ce9fa3c23a5cad5b853119c49e6))


### Bug Fixes

* Fix race in dest testing try 3 ([#683](https://github.com/cloudquery/plugin-sdk/issues/683)) ([8e8f5fe](https://github.com/cloudquery/plugin-sdk/commit/8e8f5fe75892a3c154e4ad9a809e6132f0674b8f))
* Make sure _cq_id unique across all dest plugins ([#685](https://github.com/cloudquery/plugin-sdk/issues/685)) ([a9a1173](https://github.com/cloudquery/plugin-sdk/commit/a9a1173335273858aa7baed566ec8644a059dbbf))

## [1.36.3](https://github.com/cloudquery/plugin-sdk/compare/v1.36.2...v1.36.3) (2023-02-12)


### Bug Fixes

* Take2 of fixing race in destination testing ([#680](https://github.com/cloudquery/plugin-sdk/issues/680)) ([77b74b2](https://github.com/cloudquery/plugin-sdk/commit/77b74b2cd4c28ee5f570b008105a41d0b7e8afc8))

## [1.36.2](https://github.com/cloudquery/plugin-sdk/compare/v1.36.1...v1.36.2) (2023-02-12)


### Bug Fixes

* Potential database lock/race in destination testing ([#678](https://github.com/cloudquery/plugin-sdk/issues/678)) ([50e683e](https://github.com/cloudquery/plugin-sdk/commit/50e683e7f6dfd38a25eb512c1e2417798fa832f7))

## [1.36.1](https://github.com/cloudquery/plugin-sdk/compare/v1.36.0...v1.36.1) (2023-02-12)


### Bug Fixes

* Destination testing add force tests ([#671](https://github.com/cloudquery/plugin-sdk/issues/671)) ([879f843](https://github.com/cloudquery/plugin-sdk/commit/879f843662914dc84e85c775fd62fce783c34a44))
* Fix source log message, and some debug messages ([#673](https://github.com/cloudquery/plugin-sdk/issues/673)) ([e49f593](https://github.com/cloudquery/plugin-sdk/commit/e49f5938cb9b77964ffbd4af628a27172d506baf))

## [1.36.0](https://github.com/cloudquery/plugin-sdk/compare/v1.35.0...v1.36.0) (2023-02-08)


### Features

* Add table diff methods ([#668](https://github.com/cloudquery/plugin-sdk/issues/668)) ([f6baa82](https://github.com/cloudquery/plugin-sdk/commit/f6baa82d7d1db6d28a47bf3d206306f98aa84bd4))
* Use Setpgid=true on Unix systems so that signals are not sent to the child process ([#664](https://github.com/cloudquery/plugin-sdk/issues/664)) ([2883487](https://github.com/cloudquery/plugin-sdk/commit/28834871facfe5004618362840a910c2120b11d1))


### Bug Fixes

* Remove duplicate force implementation ([#670](https://github.com/cloudquery/plugin-sdk/issues/670)) ([fe34554](https://github.com/cloudquery/plugin-sdk/commit/fe345545de1613da904f0d52c97376ab70151df4))

## [1.35.0](https://github.com/cloudquery/plugin-sdk/compare/v1.34.0...v1.35.0) (2023-02-08)


### Features

* Enable Custom Validators ([#654](https://github.com/cloudquery/plugin-sdk/issues/654)) ([6b7b5de](https://github.com/cloudquery/plugin-sdk/commit/6b7b5de46f62b3a2dbbd98fba31b790c6a170dbe))


### Bug Fixes

* **deps:** Update module golang.org/x/term to v0.5.0 ([#648](https://github.com/cloudquery/plugin-sdk/issues/648)) ([3a02bed](https://github.com/cloudquery/plugin-sdk/commit/3a02bedccb902e03ff6101ff4913dfd631977280))
* Handle null bytes in text fields ([8597f08](https://github.com/cloudquery/plugin-sdk/commit/8597f088d35b1ceb8f2f48888cb6edcbbe58a2e3))

## [1.34.0](https://github.com/cloudquery/plugin-sdk/compare/v1.33.1...v1.34.0) (2023-02-07)


### Features

* Add skip_dependent_tables option ([#662](https://github.com/cloudquery/plugin-sdk/issues/662)) ([bf34943](https://github.com/cloudquery/plugin-sdk/commit/bf349439b419b79d833a85d52073592d6ef3ba3a))


### Bug Fixes

* **logging:** Log more explicit message when OOM and other status codes occur ([#659](https://github.com/cloudquery/plugin-sdk/issues/659)) ([45c637b](https://github.com/cloudquery/plugin-sdk/commit/45c637b1127ed16ce39a471142427d19fb28fe0c))
* **logging:** Send more info logs when plugins are being terminated ([#657](https://github.com/cloudquery/plugin-sdk/issues/657)) ([6f44e1c](https://github.com/cloudquery/plugin-sdk/commit/6f44e1c597b5ca2a31e21bd099c1f556d21bf2cf))
* Remove unused `OnlyIncrementalTables` spec property ([#661](https://github.com/cloudquery/plugin-sdk/issues/661)) ([f88ba7d](https://github.com/cloudquery/plugin-sdk/commit/f88ba7d55d7644cd37a44a64719ba705e9878456))
* Trap terminate signal, log which signal we received ([#658](https://github.com/cloudquery/plugin-sdk/issues/658)) ([bb39830](https://github.com/cloudquery/plugin-sdk/commit/bb39830ff9bde75e409967f85f85f95d8919672a))

## [1.33.1](https://github.com/cloudquery/plugin-sdk/compare/v1.33.0...v1.33.1) (2023-02-01)


### Bug Fixes

* Handle numbers in env variables ([#651](https://github.com/cloudquery/plugin-sdk/issues/651)) ([0aa8f68](https://github.com/cloudquery/plugin-sdk/commit/0aa8f685e5f4c6796ef20941ed9fe7185bc44340))

## [1.33.0](https://github.com/cloudquery/plugin-sdk/compare/v1.32.0...v1.33.0) (2023-02-01)


### Features

* Support downloading plugins from other cloudquery repos ([#632](https://github.com/cloudquery/plugin-sdk/issues/632)) ([9e1501e](https://github.com/cloudquery/plugin-sdk/commit/9e1501e3db928fc283b9be43fe4b115adb6aa140))


### Bug Fixes

* **deps:** Update golang.org/x/exp digest to f062dba ([#641](https://github.com/cloudquery/plugin-sdk/issues/641)) ([c6ec154](https://github.com/cloudquery/plugin-sdk/commit/c6ec154ab4ba263b6a103f31e6e425307a6fa104))
* **deps:** Update google.golang.org/genproto digest to 1c01626 ([#642](https://github.com/cloudquery/plugin-sdk/issues/642)) ([fc9f338](https://github.com/cloudquery/plugin-sdk/commit/fc9f338804a071478ed253541cd4aff6aefd822a))
* **deps:** Update module github.com/avast/retry-go/v4 to v4.3.2 ([#643](https://github.com/cloudquery/plugin-sdk/issues/643)) ([2f6a2e8](https://github.com/cloudquery/plugin-sdk/commit/2f6a2e81cc9d687b05af21cdf96d3a29b8dfb2b4))
* **deps:** Update module github.com/getsentry/sentry-go to v0.17.0 ([#644](https://github.com/cloudquery/plugin-sdk/issues/644)) ([fb33f8c](https://github.com/cloudquery/plugin-sdk/commit/fb33f8cd3eaf426f2194c94145bd7646b355b1af))
* **deps:** Update module github.com/rs/zerolog to v1.29.0 ([#645](https://github.com/cloudquery/plugin-sdk/issues/645)) ([e864963](https://github.com/cloudquery/plugin-sdk/commit/e86496367046990d3eaf67e211225b7d3c6a9226))
* **deps:** Update module github.com/schollz/progressbar/v3 to v3.13.0 ([#646](https://github.com/cloudquery/plugin-sdk/issues/646)) ([c2146d3](https://github.com/cloudquery/plugin-sdk/commit/c2146d3cc5fba5a24041393fe5653e740e7423f2))
* **deps:** Update module golang.org/x/net to v0.5.0 ([#647](https://github.com/cloudquery/plugin-sdk/issues/647)) ([417c99d](https://github.com/cloudquery/plugin-sdk/commit/417c99d6657133312a3accd1a0e994fdab18af0a))
* **deps:** Update module golang.org/x/text to v0.6.0 ([#649](https://github.com/cloudquery/plugin-sdk/issues/649)) ([a91c7dc](https://github.com/cloudquery/plugin-sdk/commit/a91c7dc20e56c8e5858a04a74c765cc4acc2c1eb))
* **deps:** Update module google.golang.org/grpc to v1.52.3 ([#650](https://github.com/cloudquery/plugin-sdk/issues/650)) ([48d96ee](https://github.com/cloudquery/plugin-sdk/commit/48d96ee530166ae732ee34a50929eb73b8b16f2b))

## [1.32.0](https://github.com/cloudquery/plugin-sdk/compare/v1.31.0...v1.32.0) (2023-01-30)


### Features

* Return error message when download fails ([#636](https://github.com/cloudquery/plugin-sdk/issues/636)) ([0eb39af](https://github.com/cloudquery/plugin-sdk/commit/0eb39af7a294a2a9de4c81ee9950d4443e168224))


### Bug Fixes

* Add cq-dir param to discovery ([#633](https://github.com/cloudquery/plugin-sdk/issues/633)) ([13d633a](https://github.com/cloudquery/plugin-sdk/commit/13d633a6b2f1e1633325c94f7a965835e8604e88))

## [1.31.0](https://github.com/cloudquery/plugin-sdk/compare/v1.30.0...v1.31.0) (2023-01-26)


### Features

* Validate PK Creation ([#626](https://github.com/cloudquery/plugin-sdk/issues/626)) ([9ab4b46](https://github.com/cloudquery/plugin-sdk/commit/9ab4b46dfbef1872a9a16e13b0c4ab0d4e984ab3))

## [1.30.0](https://github.com/cloudquery/plugin-sdk/compare/v1.29.0...v1.30.0) (2023-01-26)


### Features

* **destination:** Filter the duplicate primary keys prior to writing batch ([#629](https://github.com/cloudquery/plugin-sdk/issues/629)) ([505709e](https://github.com/cloudquery/plugin-sdk/commit/505709eb25cee540a67bf4c55925a4ff5466a4b9)), closes [#627](https://github.com/cloudquery/plugin-sdk/issues/627)


### Bug Fixes

* Ignore env variables in comments ([#625](https://github.com/cloudquery/plugin-sdk/issues/625)) ([08bace8](https://github.com/cloudquery/plugin-sdk/commit/08bace89c708ca7f20490ce9756f8276b7e5d6f2))
* Only call `newExecutionClient` if needed ([#630](https://github.com/cloudquery/plugin-sdk/issues/630)) ([ece947f](https://github.com/cloudquery/plugin-sdk/commit/ece947f82c62be7c6bfb2f241b4644f0e2a8ae82))

## [1.29.0](https://github.com/cloudquery/plugin-sdk/compare/v1.28.0...v1.29.0) (2023-01-24)


### Features

* Add NopBackend ([#616](https://github.com/cloudquery/plugin-sdk/issues/616)) ([79f5395](https://github.com/cloudquery/plugin-sdk/commit/79f5395c5ba489564239ace9e29157d851c63158))

## [1.28.0](https://github.com/cloudquery/plugin-sdk/compare/v1.27.0...v1.28.0) (2023-01-23)


### Features

* Add version discovery service ([#619](https://github.com/cloudquery/plugin-sdk/issues/619)) ([33ab32a](https://github.com/cloudquery/plugin-sdk/commit/33ab32a690e99c00cf412097960a1d14efcff281))
* Dynamic tables and introduce proto versioning ([#610](https://github.com/cloudquery/plugin-sdk/issues/610)) ([448232c](https://github.com/cloudquery/plugin-sdk/commit/448232c8789350c8fb071902d33a5c5f07d2b82c))


### Bug Fixes

* **clients:** Update `log line too long` message ([#611](https://github.com/cloudquery/plugin-sdk/issues/611)) ([0d3ff48](https://github.com/cloudquery/plugin-sdk/commit/0d3ff48d4a8ce324b5685c3df9196943d09b2eba))
* Simplify client naming conventions ([#617](https://github.com/cloudquery/plugin-sdk/issues/617)) ([38b136b](https://github.com/cloudquery/plugin-sdk/commit/38b136b9aa15dc049f9b66dcd4ceca60fa7bdca6))

## [1.27.0](https://github.com/cloudquery/plugin-sdk/compare/v1.26.0...v1.27.0) (2023-01-17)


### Features

* **spec:** Add source, destination String methods ([#609](https://github.com/cloudquery/plugin-sdk/issues/609)) ([604b9ef](https://github.com/cloudquery/plugin-sdk/commit/604b9efe5608e60936e87114e0fbf776ea6253ea))

## [1.26.0](https://github.com/cloudquery/plugin-sdk/compare/v1.25.1...v1.26.0) (2023-01-16)


### Features

* **destinations:** Add `migrate_mode` ([#604](https://github.com/cloudquery/plugin-sdk/issues/604)) ([78b9acb](https://github.com/cloudquery/plugin-sdk/commit/78b9acbfad4183506c39ea24a4634eb1ba70c04e))


### Bug Fixes

* **destination:** Pass proper spec to client constructor ([#606](https://github.com/cloudquery/plugin-sdk/issues/606)) ([8370882](https://github.com/cloudquery/plugin-sdk/commit/837088220447a0c305888e25807163dd08042a48))

## [1.25.1](https://github.com/cloudquery/plugin-sdk/compare/v1.25.0...v1.25.1) (2023-01-14)


### Bug Fixes

* Change options for new client ([#603](https://github.com/cloudquery/plugin-sdk/issues/603)) ([f548a54](https://github.com/cloudquery/plugin-sdk/commit/f548a544f1143f60efeee3401a41f726cd707243))
* PK Addition Order ([#607](https://github.com/cloudquery/plugin-sdk/issues/607)) ([eff40e7](https://github.com/cloudquery/plugin-sdk/commit/eff40e76ae656e782a0e9745bcf34c2e5b2cd7e5))

## [1.25.0](https://github.com/cloudquery/plugin-sdk/compare/v1.24.2...v1.25.0) (2023-01-11)


### Features

* **docs:** Sort tables ([#599](https://github.com/cloudquery/plugin-sdk/issues/599)) ([8a3bfad](https://github.com/cloudquery/plugin-sdk/commit/8a3bfaddabec395cc4105ae7d2f2e99c5d31eab6))
* **transformers:** Add support for `net.IP` ([#595](https://github.com/cloudquery/plugin-sdk/issues/595)) ([a420645](https://github.com/cloudquery/plugin-sdk/commit/a420645377943939278e5d8b4a7969db957d08bf))
* **transformers:** Add WithPrimaryKeys option ([#598](https://github.com/cloudquery/plugin-sdk/issues/598)) ([107006c](https://github.com/cloudquery/plugin-sdk/commit/107006cac82e3635470bec93b086b68d0f92edf1))


### Bug Fixes

* Send resource validation errors to Sentry ([#601](https://github.com/cloudquery/plugin-sdk/issues/601)) ([5916516](https://github.com/cloudquery/plugin-sdk/commit/5916516fa9d112ba5ac146c54d02a4a1fd8850b3))

## [1.24.2](https://github.com/cloudquery/plugin-sdk/compare/v1.24.1...v1.24.2) (2023-01-11)


### Bug Fixes

* Incremental tables should not delete stale ([#594](https://github.com/cloudquery/plugin-sdk/issues/594)) ([d45e230](https://github.com/cloudquery/plugin-sdk/commit/d45e230632c2fb8035b7942dac2bb74e26d4fcb1))

## [1.24.1](https://github.com/cloudquery/plugin-sdk/compare/v1.24.0...v1.24.1) (2023-01-09)


### Bug Fixes

* Array types ([#587](https://github.com/cloudquery/plugin-sdk/issues/587)) ([73ea82c](https://github.com/cloudquery/plugin-sdk/commit/73ea82cc4abd697d428df0072f6b2ecf7002b4d1))
* Sentry errors not sent ([#592](https://github.com/cloudquery/plugin-sdk/issues/592)) ([9f1e373](https://github.com/cloudquery/plugin-sdk/commit/9f1e373b516be958f0594e84fbbbcd43951f14ad))

## [1.24.0](https://github.com/cloudquery/plugin-sdk/compare/v1.23.0...v1.24.0) (2023-01-09)


### Features

* Add local backend for storing cursor state ([#569](https://github.com/cloudquery/plugin-sdk/issues/569)) ([3b07885](https://github.com/cloudquery/plugin-sdk/commit/3b07885a57595b96dc1db5b786a6f1c22f0a5149))
* Remove codegen ([#589](https://github.com/cloudquery/plugin-sdk/issues/589)) ([1c5943a](https://github.com/cloudquery/plugin-sdk/commit/1c5943a3f1fcdd77eac89763ef3650f20f75df03))


### Bug Fixes

* **destinations:** Log correct size of batch ([#588](https://github.com/cloudquery/plugin-sdk/issues/588)) ([9cebafe](https://github.com/cloudquery/plugin-sdk/commit/9cebafef0b46c674df3027886649676cbf6c933f))

## [1.23.0](https://github.com/cloudquery/plugin-sdk/compare/v1.22.0...v1.23.0) (2023-01-09)


### Features

* Add batch size bytes as additional option ([#582](https://github.com/cloudquery/plugin-sdk/issues/582)) ([bdd76e0](https://github.com/cloudquery/plugin-sdk/commit/bdd76e04402d6da551c964a47a2bcbecd634be24))

## [1.22.0](https://github.com/cloudquery/plugin-sdk/compare/v1.21.0...v1.22.0) (2023-01-06)


### Features

* Add size in bytes to CQ types ([#510](https://github.com/cloudquery/plugin-sdk/issues/510)) ([7c15d9a](https://github.com/cloudquery/plugin-sdk/commit/7c15d9a157ef895077ac749acf4adb57deb43fd8))
* Add WithIgnoreInTestsTransformer ([#579](https://github.com/cloudquery/plugin-sdk/issues/579)) ([f836abd](https://github.com/cloudquery/plugin-sdk/commit/f836abd5addad71f3a4fa389730c4a9cdba1c219))
* Add WithResolverTransformer ([#578](https://github.com/cloudquery/plugin-sdk/issues/578)) ([5aeba0e](https://github.com/cloudquery/plugin-sdk/commit/5aeba0e1bec90a28190fae38ebc6194fa27f7653))

## [1.21.0](https://github.com/cloudquery/plugin-sdk/compare/v1.20.0...v1.21.0) (2023-01-05)


### Features

* **testing:** Add test for migrations ([#574](https://github.com/cloudquery/plugin-sdk/issues/574)) ([071a4e5](https://github.com/cloudquery/plugin-sdk/commit/071a4e5d4f91110345c69a1b787c4712ee2e7009))

## [1.20.0](https://github.com/cloudquery/plugin-sdk/compare/v1.19.0...v1.20.0) (2023-01-05)


### Features

* **transformers:** Add WithTypeTransformer ([#575](https://github.com/cloudquery/plugin-sdk/issues/575)) ([387694d](https://github.com/cloudquery/plugin-sdk/commit/387694dcbaefbbbc8154d6d237593821f64dd646))

## [1.19.0](https://github.com/cloudquery/plugin-sdk/compare/v1.18.0...v1.19.0) (2023-01-05)


### Features

* Add scheduler option and introduce Round Robin scheduler ([#545](https://github.com/cloudquery/plugin-sdk/issues/545)) ([d89a911](https://github.com/cloudquery/plugin-sdk/commit/d89a91139bf0d76833d0c756101fac58c1c15823))
* Add unwrap option to transformations ([#573](https://github.com/cloudquery/plugin-sdk/issues/573)) ([a17ee4b](https://github.com/cloudquery/plugin-sdk/commit/a17ee4bf7fb017018566ddea5d783891c7cb82d3))

## [1.18.0](https://github.com/cloudquery/plugin-sdk/compare/v1.17.2...v1.18.0) (2023-01-04)


### Features

* Add Transformer for tables (codegen replacement) ([#564](https://github.com/cloudquery/plugin-sdk/issues/564)) ([a643ddf](https://github.com/cloudquery/plugin-sdk/commit/a643ddf237fa7f40a20e525b78932d6b241b6c26))
* Support conversion of Unix timestamps in timestamptz ([#570](https://github.com/cloudquery/plugin-sdk/issues/570)) ([6b948ab](https://github.com/cloudquery/plugin-sdk/commit/6b948ab392c59c936d49182eb8b70444d81d38b5))

## [1.17.2](https://github.com/cloudquery/plugin-sdk/compare/v1.17.1...v1.17.2) (2023-01-03)


### Bug Fixes

* **testing:** Fix bug in testing missed due to reference to resource being re-used in memdb ([#567](https://github.com/cloudquery/plugin-sdk/issues/567)) ([95ab353](https://github.com/cloudquery/plugin-sdk/commit/95ab3538e29d9f253173c3a0dffd92e185cdc53c))

## [1.17.1](https://github.com/cloudquery/plugin-sdk/compare/v1.17.0...v1.17.1) (2023-01-03)


### Bug Fixes

* **testing:** Some fixes to the ordering for plugin tests ([#565](https://github.com/cloudquery/plugin-sdk/issues/565)) ([79c2b85](https://github.com/cloudquery/plugin-sdk/commit/79c2b85c38d9f42b3240559ae9b4a0d057a50607))

## [1.17.0](https://github.com/cloudquery/plugin-sdk/compare/v1.16.1...v1.17.0) (2023-01-02)


### Features

* Add primary key validation ([#563](https://github.com/cloudquery/plugin-sdk/issues/563)) ([09f891a](https://github.com/cloudquery/plugin-sdk/commit/09f891a0b34f1ec76b8143df6d7942afae506015))


### Bug Fixes

* **testing:** Sort results before comparison in tests ([#561](https://github.com/cloudquery/plugin-sdk/issues/561)) ([587715d](https://github.com/cloudquery/plugin-sdk/commit/587715de5a6fb06d861c12ece10d8b1fdf1d7ecb))

## [1.16.1](https://github.com/cloudquery/plugin-sdk/compare/v1.16.0...v1.16.1) (2023-01-01)


### Bug Fixes

* **deps:** Update golang.org/x/exp digest to 738e83a ([#546](https://github.com/cloudquery/plugin-sdk/issues/546)) ([bdf3ff1](https://github.com/cloudquery/plugin-sdk/commit/bdf3ff1e9e93164e20e73046534fa1a8dd208576))
* **deps:** Update google.golang.org/genproto digest to f9683d7 ([#552](https://github.com/cloudquery/plugin-sdk/issues/552)) ([763d22b](https://github.com/cloudquery/plugin-sdk/commit/763d22b5f209ae26d54937c84f07f8895062ebc5))
* **deps:** Update module github.com/getsentry/sentry-go to v0.16.0 ([#549](https://github.com/cloudquery/plugin-sdk/issues/549)) ([b4a0efc](https://github.com/cloudquery/plugin-sdk/commit/b4a0efc392a9323011d217c40ca3661d38351c37))
* **deps:** Update module github.com/inconshreveable/mousetrap to v1.1.0 ([#555](https://github.com/cloudquery/plugin-sdk/issues/555)) ([f449234](https://github.com/cloudquery/plugin-sdk/commit/f4492343b52a8edf3864b4d77c4e2f40d0d3e308))
* **deps:** Update module github.com/mattn/go-isatty to v0.0.17 ([#553](https://github.com/cloudquery/plugin-sdk/issues/553)) ([826006f](https://github.com/cloudquery/plugin-sdk/commit/826006f6d70e9cb1c4a062d5691be05b41514926))
* **deps:** Update module github.com/schollz/progressbar/v3 to v3.12.2 ([#547](https://github.com/cloudquery/plugin-sdk/issues/547)) ([b6640b8](https://github.com/cloudquery/plugin-sdk/commit/b6640b8134aff9d9c12c211d0994eda657a966d0))
* **deps:** Update module github.com/thoas/go-funk to v0.9.3 ([#548](https://github.com/cloudquery/plugin-sdk/issues/548)) ([6e5469a](https://github.com/cloudquery/plugin-sdk/commit/6e5469a32ec688b94070f300633050fbe2e53018))
* **deps:** Update module golang.org/x/net to v0.4.0 ([#550](https://github.com/cloudquery/plugin-sdk/issues/550)) ([9ced5ec](https://github.com/cloudquery/plugin-sdk/commit/9ced5ec92f60be484d470550781110b1a3b6a2d0))
* **deps:** Update module golang.org/x/text to v0.5.0 ([#551](https://github.com/cloudquery/plugin-sdk/issues/551)) ([1353026](https://github.com/cloudquery/plugin-sdk/commit/1353026325232a7de6c0ea403cdcbe5e821abe53))
* Managed writer log message, timeout-&gt;flush ([#536](https://github.com/cloudquery/plugin-sdk/issues/536)) ([6b0c711](https://github.com/cloudquery/plugin-sdk/commit/6b0c71174d71c4fc5f5a55f9317caa1037f75d15))

## [1.16.0](https://github.com/cloudquery/plugin-sdk/compare/v1.15.1...v1.16.0) (2022-12-28)


### Features

* **destinations:** Allow plugins to set default batch size ([#540](https://github.com/cloudquery/plugin-sdk/issues/540)) ([bc1476b](https://github.com/cloudquery/plugin-sdk/commit/bc1476b0d6a7f9b3014c2d78108fc5a499399893))

## [1.15.1](https://github.com/cloudquery/plugin-sdk/compare/v1.15.0...v1.15.1) (2022-12-28)


### Bug Fixes

* **destinations:** Set done even if no resources to flush ([#537](https://github.com/cloudquery/plugin-sdk/issues/537)) ([02eca6d](https://github.com/cloudquery/plugin-sdk/commit/02eca6d1962d306f7571cdfc4f4255ef93a98c02))

## [1.15.0](https://github.com/cloudquery/plugin-sdk/compare/v1.14.0...v1.15.0) (2022-12-28)


### Features

* Make TestData public ([#534](https://github.com/cloudquery/plugin-sdk/issues/534)) ([a476052](https://github.com/cloudquery/plugin-sdk/commit/a4760521cff17b251a0c90b4cb45eaa8257d6fe2))

## [1.14.0](https://github.com/cloudquery/plugin-sdk/compare/v1.13.1...v1.14.0) (2022-12-27)


### Features

* Add basic periodic metric INFO logger ([#496](https://github.com/cloudquery/plugin-sdk/issues/496)) ([8d1d32e](https://github.com/cloudquery/plugin-sdk/commit/8d1d32eacf34a7835cb9e712cc448c66d7894b55))


### Bug Fixes

* **destinations:** Stop writing resources when channel is closed ([#460](https://github.com/cloudquery/plugin-sdk/issues/460)) ([5590845](https://github.com/cloudquery/plugin-sdk/commit/5590845d5ce9f3395a57e6c1997c2e4071b41952))
* Don't hide errors in destination server ([#529](https://github.com/cloudquery/plugin-sdk/issues/529)) ([d91f94f](https://github.com/cloudquery/plugin-sdk/commit/d91f94fc8bd74830c88c42d4e8a1bee16bcbd2a7))

## [1.13.1](https://github.com/cloudquery/plugin-sdk/compare/v1.13.0...v1.13.1) (2022-12-22)


### Bug Fixes

* Typo manager-&gt;managed ([#526](https://github.com/cloudquery/plugin-sdk/issues/526)) ([7503b1f](https://github.com/cloudquery/plugin-sdk/commit/7503b1fba9fbd42e423207195ae8af93c988ea99))

## [1.13.0](https://github.com/cloudquery/plugin-sdk/compare/v1.12.7...v1.13.0) (2022-12-21)


### Features

* Add managed API for destination plugins ([#521](https://github.com/cloudquery/plugin-sdk/issues/521)) ([3df6129](https://github.com/cloudquery/plugin-sdk/commit/3df6129255784dc54707755da9ddd81b848b4a2d))

## [1.12.7](https://github.com/cloudquery/plugin-sdk/compare/v1.12.6...v1.12.7) (2022-12-19)


### Bug Fixes

* **destination:** Rename `NewDestinationPlugin` to `NewPlugin` ([#519](https://github.com/cloudquery/plugin-sdk/issues/519)) ([3934775](https://github.com/cloudquery/plugin-sdk/commit/39347757ba443e93ab36de86c8672223f9554145))

## [1.12.6](https://github.com/cloudquery/plugin-sdk/compare/v1.12.5...v1.12.6) (2022-12-18)


### Bug Fixes

* Add better logging/metric per table ([#513](https://github.com/cloudquery/plugin-sdk/issues/513)) ([da36396](https://github.com/cloudquery/plugin-sdk/commit/da363966a7f74adb85280cc6688e0c573112e506))
* Improve formatting of newlines in markdown files ([#492](https://github.com/cloudquery/plugin-sdk/issues/492)) ([e48ff90](https://github.com/cloudquery/plugin-sdk/commit/e48ff90e0b38ea67efc5648e0bff4895938545ce))
* Include table name in logs on panic ([#505](https://github.com/cloudquery/plugin-sdk/issues/505)) ([a0b8a46](https://github.com/cloudquery/plugin-sdk/commit/a0b8a46c05b86ce3276d7f5455ca0762579db532))
* Move source & destination plugin code to separate packages ([#516](https://github.com/cloudquery/plugin-sdk/issues/516)) ([6733785](https://github.com/cloudquery/plugin-sdk/commit/67337856a8c973ecb5fb4749078f63e9b9909129))
* Use correct error codes ([#514](https://github.com/cloudquery/plugin-sdk/issues/514)) ([8b53d76](https://github.com/cloudquery/plugin-sdk/commit/8b53d76ca155eb95526698d16a2233faf4fd4a1e))

## [1.12.5](https://github.com/cloudquery/plugin-sdk/compare/v1.12.4...v1.12.5) (2022-12-14)


### Bug Fixes

* Don't print value with error on invalid JSON ([#503](https://github.com/cloudquery/plugin-sdk/issues/503)) ([4b36824](https://github.com/cloudquery/plugin-sdk/commit/4b368246dcb470f87933bd7e7f575e201befa7c1))

## [1.12.4](https://github.com/cloudquery/plugin-sdk/compare/v1.12.3...v1.12.4) (2022-12-14)


### Bug Fixes

* Use json.Valid ([#500](https://github.com/cloudquery/plugin-sdk/issues/500)) ([4242e5e](https://github.com/cloudquery/plugin-sdk/commit/4242e5ec3ad674cccb7d8597d3c016b68ab563bd))

## [1.12.3](https://github.com/cloudquery/plugin-sdk/compare/v1.12.2...v1.12.3) (2022-12-14)


### Bug Fixes

* Throw error on empty env variable ([#499](https://github.com/cloudquery/plugin-sdk/issues/499)) ([4b77cf5](https://github.com/cloudquery/plugin-sdk/commit/4b77cf511f7c6a05fdeb96941da2eaf0c3a80fa0))
* Validate json strings and handle empty strings ([#497](https://github.com/cloudquery/plugin-sdk/issues/497)) ([dd5f008](https://github.com/cloudquery/plugin-sdk/commit/dd5f008ee46561663555fc419d0246bfc3bc8be0))

## [1.12.2](https://github.com/cloudquery/plugin-sdk/compare/v1.12.1...v1.12.2) (2022-12-13)


### Bug Fixes

* Glob table filtering ([#494](https://github.com/cloudquery/plugin-sdk/issues/494)) ([d6c126b](https://github.com/cloudquery/plugin-sdk/commit/d6c126bfa59321f8cf3f521c800a496f386ae961))

## [1.12.1](https://github.com/cloudquery/plugin-sdk/compare/v1.12.0...v1.12.1) (2022-12-13)


### Bug Fixes

* Don't panic on empty-string for timestamp ([#489](https://github.com/cloudquery/plugin-sdk/issues/489)) ([83813de](https://github.com/cloudquery/plugin-sdk/commit/83813de73b4d907bd6bdd93b47e53bf5800f0805))
* Fix deadlock off-by-one ([#493](https://github.com/cloudquery/plugin-sdk/issues/493)) ([4ea9ed8](https://github.com/cloudquery/plugin-sdk/commit/4ea9ed82eed9528a2cb2f74ffe80d8e5e75a83d6))
* Reduce default concurrency ([#491](https://github.com/cloudquery/plugin-sdk/issues/491)) ([f995da9](https://github.com/cloudquery/plugin-sdk/commit/f995da9d2f4c2dfe7d0a09107a610a7cd700ce5a))
* Refactor glob filters ([#488](https://github.com/cloudquery/plugin-sdk/issues/488)) ([cb5f6bb](https://github.com/cloudquery/plugin-sdk/commit/cb5f6bbd111a3532fa0ad37039894c60fda52ef4))

## [1.12.0](https://github.com/cloudquery/plugin-sdk/compare/v1.11.2...v1.12.0) (2022-12-11)


### Features

* Add handling for json.Number in faker ([#481](https://github.com/cloudquery/plugin-sdk/issues/481)) ([ad20787](https://github.com/cloudquery/plugin-sdk/commit/ad2078708d66b3667ba7718e24b43f95db6eba02))


### Bug Fixes

* Allow both 'yml' and 'yaml' extensions ([#476](https://github.com/cloudquery/plugin-sdk/issues/476)) ([52c4c56](https://github.com/cloudquery/plugin-sdk/commit/52c4c566b7b06498562a48f8591d24fe49c37bc7))
* **errors:** Remove usage of `codes.Internal` ([#485](https://github.com/cloudquery/plugin-sdk/issues/485)) ([62692b9](https://github.com/cloudquery/plugin-sdk/commit/62692b9cb8a3ff3465d9d14a1ec7cc801d3490af))

## [1.11.2](https://github.com/cloudquery/plugin-sdk/compare/v1.11.1...v1.11.2) (2022-12-08)


### Bug Fixes

* Initialise clients only once ([#473](https://github.com/cloudquery/plugin-sdk/issues/473)) ([c88a521](https://github.com/cloudquery/plugin-sdk/commit/c88a521dbb9793cc8acc08c11826f1b158f2669b))

## [1.11.1](https://github.com/cloudquery/plugin-sdk/compare/v1.11.0...v1.11.1) (2022-12-07)


### Bug Fixes

* **codegen:** Column type for slices ([7474c90](https://github.com/cloudquery/plugin-sdk/commit/7474c90415119082bdb1cdb145bd16d1ef51a3b2))
* Concurrent read,write to a map ([#467](https://github.com/cloudquery/plugin-sdk/issues/467)) ([ebef24a](https://github.com/cloudquery/plugin-sdk/commit/ebef24a00e667aab04c1e5258c7f9c70757894d6))
* **sentry:** Use HTTPSyncTransport, remove flush ([#465](https://github.com/cloudquery/plugin-sdk/issues/465)) ([4d48306](https://github.com/cloudquery/plugin-sdk/commit/4d483064218fbabea350297260dec59bc547bc6e))
* Skip relations when initializing metrics ([#469](https://github.com/cloudquery/plugin-sdk/issues/469)) ([5efe564](https://github.com/cloudquery/plugin-sdk/commit/5efe56493a21264172655bfc0b769be72d135c91))

## [1.11.0](https://github.com/cloudquery/plugin-sdk/compare/v1.10.0...v1.11.0) (2022-12-05)


### Features

* Add Support for net.IP in Faker ([#445](https://github.com/cloudquery/plugin-sdk/issues/445)) ([2deced1](https://github.com/cloudquery/plugin-sdk/commit/2deced12ec51d504840d064be367d70e855697f5))


### Bug Fixes

* **deps:** Update golang.org/x/exp digest to 6ab00d0 ([#449](https://github.com/cloudquery/plugin-sdk/issues/449)) ([b981e33](https://github.com/cloudquery/plugin-sdk/commit/b981e3301e53fa2f8d0b7a854b50fec84ad28a3a))
* **deps:** Update module github.com/avast/retry-go/v4 to v4.3.1 ([#450](https://github.com/cloudquery/plugin-sdk/issues/450)) ([e4116f1](https://github.com/cloudquery/plugin-sdk/commit/e4116f1982dbe6fb4bd5222dbc8d58af551b56b1))
* **deps:** Update module google.golang.org/grpc to v1.51.0 ([#451](https://github.com/cloudquery/plugin-sdk/issues/451)) ([538211c](https://github.com/cloudquery/plugin-sdk/commit/538211c863ec2d4b719b83086a842e89ecc396d3))
* Don't log start-and-finish of relational tables ([#459](https://github.com/cloudquery/plugin-sdk/issues/459)) ([4d6eeca](https://github.com/cloudquery/plugin-sdk/commit/4d6eecac9d9ed76caf064287b1f15fd321c7a651))
* Sync Metrics for Long running plugins ([#455](https://github.com/cloudquery/plugin-sdk/issues/455)) ([3fecc61](https://github.com/cloudquery/plugin-sdk/commit/3fecc612db841db289796f0dd77dfe9efa10847f))
* TablesForSpec should only return top-level tables ([#456](https://github.com/cloudquery/plugin-sdk/issues/456)) ([ab7ca97](https://github.com/cloudquery/plugin-sdk/commit/ab7ca972e0b187a7dfb66132a03f07479cd29bb7))

## [1.10.0](https://github.com/cloudquery/plugin-sdk/compare/v1.9.0...v1.10.0) (2022-11-29)


### Features

* Add function to list tables that match a source spec ([#440](https://github.com/cloudquery/plugin-sdk/issues/440)) ([a8f3690](https://github.com/cloudquery/plugin-sdk/commit/a8f369029dd90a1530112fcf2b675fc9e4f2e0d8))

## [1.9.0](https://github.com/cloudquery/plugin-sdk/compare/v1.8.2...v1.9.0) (2022-11-25)


### Features

* Handle resolving of empty maps and slices ([#430](https://github.com/cloudquery/plugin-sdk/issues/430)) ([a5672b5](https://github.com/cloudquery/plugin-sdk/commit/a5672b5faa9f41f2179650f989761217575b3934))


### Bug Fixes

* Fix docs for deeply nested tables ([#434](https://github.com/cloudquery/plugin-sdk/issues/434)) ([48e0466](https://github.com/cloudquery/plugin-sdk/commit/48e04662a6afc82dba084efa5f91bbe1470b2d43))

## [1.8.2](https://github.com/cloudquery/plugin-sdk/compare/v1.8.1...v1.8.2) (2022-11-25)


### Bug Fixes

* **test:** Values check test should account for `IgnoreInTests` column option ([#431](https://github.com/cloudquery/plugin-sdk/issues/431)) ([ffffcd5](https://github.com/cloudquery/plugin-sdk/commit/ffffcd54ff2036b2af5a3539a2b10f4b2a65abb5))

## [1.8.1](https://github.com/cloudquery/plugin-sdk/compare/v1.8.0...v1.8.1) (2022-11-24)


### Bug Fixes

* Small improvement to PK checking in codegen ([#432](https://github.com/cloudquery/plugin-sdk/issues/432)) ([15f7d1b](https://github.com/cloudquery/plugin-sdk/commit/15f7d1b4dfbdf1966650be8f93d85cb4492e0767))

## [1.8.0](https://github.com/cloudquery/plugin-sdk/compare/v1.7.0...v1.8.0) (2022-11-23)


### Features

* Resolve table relations in parallel ([#423](https://github.com/cloudquery/plugin-sdk/issues/423)) ([ede04b7](https://github.com/cloudquery/plugin-sdk/commit/ede04b7c01d11a833a2c894e229f41656f85b036))

## [1.7.0](https://github.com/cloudquery/plugin-sdk/compare/v1.6.0...v1.7.0) (2022-11-22)


### Features

* Resolve table relations in parallel ([#416](https://github.com/cloudquery/plugin-sdk/issues/416)) ([aadbde9](https://github.com/cloudquery/plugin-sdk/commit/aadbde9064eb30c2412c13d9e770e216e8c57ec9))


### Bug Fixes

* Revert "feat: Resolve table relations in parallel" ([#422](https://github.com/cloudquery/plugin-sdk/issues/422)) ([655a04b](https://github.com/cloudquery/plugin-sdk/commit/655a04b8f9d8c7857a800e0666392a02a4c805ba))
* Skip very large gRPC messages, log when it happens ([#421](https://github.com/cloudquery/plugin-sdk/issues/421)) ([0874d58](https://github.com/cloudquery/plugin-sdk/commit/0874d585d2fc4cddc890fbb9d92423ad7c1029fe))

## [1.6.0](https://github.com/cloudquery/plugin-sdk/compare/v1.5.3...v1.6.0) (2022-11-21)


### Features

* Add option to plugin doc command to output tables as JSON ([#347](https://github.com/cloudquery/plugin-sdk/issues/347)) ([c1b4240](https://github.com/cloudquery/plugin-sdk/commit/c1b424044d2e8aa33d833222b5d7d09a7b606ae7))
* Support ${file:./path} expansion in spec ([#418](https://github.com/cloudquery/plugin-sdk/issues/418)) ([58d7c44](https://github.com/cloudquery/plugin-sdk/commit/58d7c4420431142ac95fa2eb2cb16ce64d6ba179))


### Bug Fixes

* Fix Destination testing suite ([#417](https://github.com/cloudquery/plugin-sdk/issues/417)) ([4771efa](https://github.com/cloudquery/plugin-sdk/commit/4771efadf9c5a0ba8ace33af89614557a721072e))
* Increase GRPC message size limit to 50MiB ([#419](https://github.com/cloudquery/plugin-sdk/issues/419)) ([a54c6ea](https://github.com/cloudquery/plugin-sdk/commit/a54c6ea15d0af87b3c314f21f62e7ec9071e372f))

## [1.5.3](https://github.com/cloudquery/plugin-sdk/compare/v1.5.2...v1.5.3) (2022-11-15)


### Bug Fixes

* Workaround Go Inet marshal bug ([#410](https://github.com/cloudquery/plugin-sdk/issues/410)) ([bd7718c](https://github.com/cloudquery/plugin-sdk/commit/bd7718c3a5a76d8c0c70db66d5a6231450ad9e78))

## [1.5.2](https://github.com/cloudquery/plugin-sdk/compare/v1.5.1...v1.5.2) (2022-11-14)


### Bug Fixes

* Update libs ([#406](https://github.com/cloudquery/plugin-sdk/issues/406)) ([04d6ca8](https://github.com/cloudquery/plugin-sdk/commit/04d6ca88783817a51157b99f52005bf86d395d50))

## [1.5.1](https://github.com/cloudquery/plugin-sdk/compare/v1.5.0...v1.5.1) (2022-11-14)


### Bug Fixes

* Allow searching relations by name ([#404](https://github.com/cloudquery/plugin-sdk/issues/404)) ([45da719](https://github.com/cloudquery/plugin-sdk/commit/45da719a8368de20d80b6837a916fada9443d130))

## [1.5.0](https://github.com/cloudquery/plugin-sdk/compare/v1.4.1...v1.5.0) (2022-11-11)


### Features

* Add support for glob matching in config ([#398](https://github.com/cloudquery/plugin-sdk/issues/398)) ([c866573](https://github.com/cloudquery/plugin-sdk/commit/c866573ba656e4a23ed0c0bc9576c1beb708a4c6))


### Bug Fixes

* Change globbing behavior to include descendants by default ([#403](https://github.com/cloudquery/plugin-sdk/issues/403)) ([de15d26](https://github.com/cloudquery/plugin-sdk/commit/de15d2610388eb8572baa23cb0fc5df86aea1950))
* Exit early if all Write workers have stopped ([#395](https://github.com/cloudquery/plugin-sdk/issues/395)) ([5707e7a](https://github.com/cloudquery/plugin-sdk/commit/5707e7a132d44cea712753590081724edf26725d))

## [1.4.1](https://github.com/cloudquery/plugin-sdk/compare/v1.4.0...v1.4.1) (2022-11-10)


### Bug Fixes

* Pre-aggregate metrics before sending ([#396](https://github.com/cloudquery/plugin-sdk/issues/396)) ([b6b5f7f](https://github.com/cloudquery/plugin-sdk/commit/b6b5f7fb57d89e0d50deaf27467f75bb014d3616))

## [1.4.0](https://github.com/cloudquery/plugin-sdk/compare/v1.3.2...v1.4.0) (2022-11-10)


### Features

* **codegen:** Allow passing slices ([#386](https://github.com/cloudquery/plugin-sdk/issues/386)) ([dbc28d8](https://github.com/cloudquery/plugin-sdk/commit/dbc28d8419e3e3fa5682a537d11b80787ad2d036))


### Bug Fixes

* Clear skip tables error on invalid or child table skippage ([#349](https://github.com/cloudquery/plugin-sdk/issues/349)) ([bb0c60b](https://github.com/cloudquery/plugin-sdk/commit/bb0c60bd9d86f2dab5853ff6377bfb789a0dbf7d))

## [1.3.2](https://github.com/cloudquery/plugin-sdk/compare/v1.3.1...v1.3.2) (2022-11-10)


### Bug Fixes

* Add -race when running tests ([#388](https://github.com/cloudquery/plugin-sdk/issues/388)) ([3da08bb](https://github.com/cloudquery/plugin-sdk/commit/3da08bb89c3c381cbc87b5dc8b53408bef5b4a9d))
* Close zip archive when we're done with it ([#391](https://github.com/cloudquery/plugin-sdk/issues/391)) ([1c4a877](https://github.com/cloudquery/plugin-sdk/commit/1c4a877662b3a84f99b8a942d918f0d39d90e869))

## [1.3.1](https://github.com/cloudquery/plugin-sdk/compare/v1.3.0...v1.3.1) (2022-11-10)


### Bug Fixes

* **deps:** Revert dependencies updates ([#389](https://github.com/cloudquery/plugin-sdk/issues/389)) ([3bc5314](https://github.com/cloudquery/plugin-sdk/commit/3bc5314907de511ad15eeea2257588eecf68a35a))

## [1.3.0](https://github.com/cloudquery/plugin-sdk/compare/v1.2.0...v1.3.0) (2022-11-09)


### Features

* **codegen:** Add `WithPKColumns` option ([#379](https://github.com/cloudquery/plugin-sdk/issues/379)) ([0e3457d](https://github.com/cloudquery/plugin-sdk/commit/0e3457de7b3c8de1e1f21330d98a1a7a1806ccc3))

## [1.2.0](https://github.com/cloudquery/plugin-sdk/compare/v1.1.2...v1.2.0) (2022-11-09)


### Features

* **codegen:** Add sanity check to `TableDefinition` ([#376](https://github.com/cloudquery/plugin-sdk/issues/376)) ([49c27b5](https://github.com/cloudquery/plugin-sdk/commit/49c27b515d1e0318c986d9c0bd58ce7a17c0a0d7))


### Bug Fixes

* Revert "fix(faker): Use `MarshalText` for faker timestamps ([#373](https://github.com/cloudquery/plugin-sdk/issues/373))" ([#381](https://github.com/cloudquery/plugin-sdk/issues/381)) ([a01ec51](https://github.com/cloudquery/plugin-sdk/commit/a01ec517c63d18e103aaa7c09e49c620f87a8c76))
* Update `resolveResource` timeout to 10 minutes ([#384](https://github.com/cloudquery/plugin-sdk/issues/384)) ([456ef2f](https://github.com/cloudquery/plugin-sdk/commit/456ef2fd19fb1e15ccf9929bc0b092580d040011))
* Use MarshalText when serializing timestamps when applicable ([#382](https://github.com/cloudquery/plugin-sdk/issues/382)) ([b110a90](https://github.com/cloudquery/plugin-sdk/commit/b110a9095ffb705289eb8a250eeb390ba5450a50))

## [1.1.2](https://github.com/cloudquery/plugin-sdk/compare/v1.1.1...v1.1.2) (2022-11-09)


### Bug Fixes

* **faker:** Use `MarshalText` for faker timestamps ([#373](https://github.com/cloudquery/plugin-sdk/issues/373)) ([a291438](https://github.com/cloudquery/plugin-sdk/commit/a29143861b22432c81cdc8b04650d9d8d0ac9671))

## [1.1.1](https://github.com/cloudquery/plugin-sdk/compare/v1.1.0...v1.1.1) (2022-11-09)


### Bug Fixes

* Context cancelled too early for delete stale mode ([#377](https://github.com/cloudquery/plugin-sdk/issues/377)) ([cd7bf6d](https://github.com/cloudquery/plugin-sdk/commit/cd7bf6d90b8b4942919165f0b5cda7ac33b238e3))

## [1.1.0](https://github.com/cloudquery/plugin-sdk/compare/v1.0.4...v1.1.0) (2022-11-08)


### Features

* Add Testing suite for destination plugins ([#369](https://github.com/cloudquery/plugin-sdk/issues/369)) ([1a542b9](https://github.com/cloudquery/plugin-sdk/commit/1a542b9bf23219373d0b683030770c2f15502016))

## [1.0.4](https://github.com/cloudquery/plugin-sdk/compare/v1.0.3...v1.0.4) (2022-11-08)


### Bug Fixes

* Make path a required config parameter ([#368](https://github.com/cloudquery/plugin-sdk/issues/368)) ([77fdaf8](https://github.com/cloudquery/plugin-sdk/commit/77fdaf85c1f580b760694ed7fb0563be71d06726))

## [1.0.3](https://github.com/cloudquery/plugin-sdk/compare/v1.0.2...v1.0.3) (2022-11-07)


### Bug Fixes

* Allow managed clients to disable sentry logging ([#363](https://github.com/cloudquery/plugin-sdk/issues/363)) ([dc20388](https://github.com/cloudquery/plugin-sdk/commit/dc203886a6b077afa4e1b1138c3c1c60b0fcd2f2))
* Normalize Windows line breaks before parsing configuration files ([#352](https://github.com/cloudquery/plugin-sdk/issues/352)) ([979e207](https://github.com/cloudquery/plugin-sdk/commit/979e207831a2835943a420791fc9598ada2efbf7))

## [1.0.2](https://github.com/cloudquery/plugin-sdk/compare/v1.0.1...v1.0.2) (2022-11-07)


### Bug Fixes

* Revert "chore: Start SDK semantic versioning from v1" ([#366](https://github.com/cloudquery/plugin-sdk/issues/366)) ([c66be4b](https://github.com/cloudquery/plugin-sdk/commit/c66be4bb440990327c6e1aa82e5ffdd76659bd07))

## [1.0.1](https://github.com/cloudquery/plugin-sdk/compare/v1.0.0...v1.0.1) (2022-11-07)


### Bug Fixes

* Module parameter in logs of source-plugins ([#364](https://github.com/cloudquery/plugin-sdk/issues/364)) ([379d3e6](https://github.com/cloudquery/plugin-sdk/commit/379d3e639599e14fe112ef301e59f22a27923f00))

## [1.0.0](https://github.com/cloudquery/plugin-sdk/compare/v0.13.23...v1.0.0) (2022-11-07)


### Bug Fixes

* Dont use reflection in reverse transformer ([#360](https://github.com/cloudquery/plugin-sdk/issues/360)) ([9c85c1a](https://github.com/cloudquery/plugin-sdk/commit/9c85c1a14e6740af8adecf6c9580c924fd0dcd9c))


### Miscellaneous Chores

* Start SDK semantic versioning from v1 ([#362](https://github.com/cloudquery/plugin-sdk/issues/362)) ([40041c8](https://github.com/cloudquery/plugin-sdk/commit/40041c8c3544c6189a4b3975c72637abd5c52bc0))

## [0.13.23](https://github.com/cloudquery/plugin-sdk/compare/v0.13.22...v0.13.23) (2022-11-07)


### Bug Fixes

* Move cqtypes ([#357](https://github.com/cloudquery/plugin-sdk/issues/357)) ([9064bc0](https://github.com/cloudquery/plugin-sdk/commit/9064bc0bdf4da2d6dcdd038378a67dba3bd73422))

## [0.13.22](https://github.com/cloudquery/plugin-sdk/compare/v0.13.21...v0.13.22) (2022-11-06)


### Bug Fixes

* Include source path in dest to source map key ([#353](https://github.com/cloudquery/plugin-sdk/issues/353)) ([ac727f6](https://github.com/cloudquery/plugin-sdk/commit/ac727f66bb03f3a1e9cbec79cb37073819bb6981))

## [0.13.21](https://github.com/cloudquery/plugin-sdk/compare/v0.13.20...v0.13.21) (2022-11-04)


### Bug Fixes

* Disallow child tables ([#342](https://github.com/cloudquery/plugin-sdk/issues/342)) ([24922a7](https://github.com/cloudquery/plugin-sdk/commit/24922a70794ec6c6f7b134580995c608f2672cc2))

## [0.13.20](https://github.com/cloudquery/plugin-sdk/compare/v0.13.19...v0.13.20) (2022-11-04)


### Features

* Add retry logic when downloading plugins from GitHub ([#310](https://github.com/cloudquery/plugin-sdk/issues/310)) ([914d252](https://github.com/cloudquery/plugin-sdk/commit/914d252d74dd39d15402898a398673bb3553252e))
* Enable Multiline table description ([#345](https://github.com/cloudquery/plugin-sdk/issues/345)) ([d83c60a](https://github.com/cloudquery/plugin-sdk/commit/d83c60a2ce7bba0b190d3d5ae64400a2a6161195))

## [0.13.19](https://github.com/cloudquery/plugin-sdk/compare/v0.13.18...v0.13.19) (2022-11-03)


### Bug Fixes

* Check for missing parent tables ([#339](https://github.com/cloudquery/plugin-sdk/issues/339)) ([49fabc7](https://github.com/cloudquery/plugin-sdk/commit/49fabc7abf155c03d72b5417196a63c33b29495e))

## [0.13.18](https://github.com/cloudquery/plugin-sdk/compare/v0.13.17...v0.13.18) (2022-11-01)


### Bug Fixes

* Parsing timestamptz default string ([#336](https://github.com/cloudquery/plugin-sdk/issues/336)) ([acdcb02](https://github.com/cloudquery/plugin-sdk/commit/acdcb02b48ca2e0009a998d710b88f60830295d0))

## [0.13.17](https://github.com/cloudquery/plugin-sdk/compare/v0.13.16...v0.13.17) (2022-11-01)


### Bug Fixes

* **deps:** Update golang.org/x/exp digest to c99f073 ([#324](https://github.com/cloudquery/plugin-sdk/issues/324)) ([c33c33d](https://github.com/cloudquery/plugin-sdk/commit/c33c33d4a8e6ec6b7dcc32fea2358d694c6e8161))
* **deps:** Update module github.com/getsentry/sentry-go to v0.14.0 ([#328](https://github.com/cloudquery/plugin-sdk/issues/328)) ([446447a](https://github.com/cloudquery/plugin-sdk/commit/446447adcdb8dab3f7064c58280ad32438b68c3b))
* **deps:** Update module github.com/grpc-ecosystem/go-grpc-middleware/providers/zerolog/v2 to v2.0.0-rc.3 ([#325](https://github.com/cloudquery/plugin-sdk/issues/325)) ([da6e340](https://github.com/cloudquery/plugin-sdk/commit/da6e340cd3a31d049805467d39ced69a6a06dc1a))
* **deps:** Update module github.com/spf13/cobra to v1.6.1 ([#329](https://github.com/cloudquery/plugin-sdk/issues/329)) ([ec583d2](https://github.com/cloudquery/plugin-sdk/commit/ec583d2fca8e54edbabf8719f21d20a133b66331))
* **deps:** Update module github.com/stretchr/testify to v1.8.1 ([#327](https://github.com/cloudquery/plugin-sdk/issues/327)) ([f9904de](https://github.com/cloudquery/plugin-sdk/commit/f9904dee98d4411a3a0b6e62bfd7478ed4d2e81f))
* **deps:** Update module golang.org/x/net to v0.1.0 ([#330](https://github.com/cloudquery/plugin-sdk/issues/330)) ([06e8426](https://github.com/cloudquery/plugin-sdk/commit/06e84261e7fc5e9c0a146afea09e52b61a9549b9))
* **deps:** Update module golang.org/x/sync to v0.1.0 ([#331](https://github.com/cloudquery/plugin-sdk/issues/331)) ([489d6b7](https://github.com/cloudquery/plugin-sdk/commit/489d6b752cb8f88a9a6d2f89af6ad6faa2d0bb5e))
* **deps:** Update module golang.org/x/text to v0.4.0 ([#332](https://github.com/cloudquery/plugin-sdk/issues/332)) ([314a172](https://github.com/cloudquery/plugin-sdk/commit/314a1723cbce3a2020cb119899875642cb9739c1))
* **deps:** Update module google.golang.org/grpc to v1.50.1 ([#334](https://github.com/cloudquery/plugin-sdk/issues/334)) ([a24ce80](https://github.com/cloudquery/plugin-sdk/commit/a24ce8066eecd51d1c28cb30405b263a734ecb34))
* Try formatting timestamptz in a few formats ([#322](https://github.com/cloudquery/plugin-sdk/issues/322)) ([543638c](https://github.com/cloudquery/plugin-sdk/commit/543638c1fd3e975ffaaff0209c0393edffca11ec))

## [0.13.16](https://github.com/cloudquery/plugin-sdk/compare/v0.13.15...v0.13.16) (2022-10-31)


### Features

* Add CQ type system to support multiple destinations ([#320](https://github.com/cloudquery/plugin-sdk/issues/320)) ([d3b24a0](https://github.com/cloudquery/plugin-sdk/commit/d3b24a006d2f0d906076ed77b6cf427045d15fa1))

## [0.13.15](https://github.com/cloudquery/plugin-sdk/compare/v0.13.14...v0.13.15) (2022-10-30)


### Features

* Add Metrics and improve scheduler with DFS ([#318](https://github.com/cloudquery/plugin-sdk/issues/318)) ([2d7a83b](https://github.com/cloudquery/plugin-sdk/commit/2d7a83beae21e1e7ad8ff8b7aec0f5954475f476))

## [0.13.14](https://github.com/cloudquery/plugin-sdk/compare/v0.13.13...v0.13.14) (2022-10-27)


### Bug Fixes

* Revert "fix(deps): Update go-funk ([#312](https://github.com/cloudquery/plugin-sdk/issues/312))" ([#314](https://github.com/cloudquery/plugin-sdk/issues/314)) ([06a33ab](https://github.com/cloudquery/plugin-sdk/commit/06a33ab12b52c1e5b576f280a2bec03d396db063))

## [0.13.13](https://github.com/cloudquery/plugin-sdk/compare/v0.13.12...v0.13.13) (2022-10-27)


### Bug Fixes

* **deps:** Update go-funk ([#312](https://github.com/cloudquery/plugin-sdk/issues/312)) ([fea5c28](https://github.com/cloudquery/plugin-sdk/commit/fea5c2855d46d1cefacb9ed826dc78dfad45a6f7))

## [0.13.12](https://github.com/cloudquery/plugin-sdk/compare/v0.13.11...v0.13.12) (2022-10-20)


### Bug Fixes

* Set Sentry server name to empty to avoid sending it ([#305](https://github.com/cloudquery/plugin-sdk/issues/305)) ([4b0bfd4](https://github.com/cloudquery/plugin-sdk/commit/4b0bfd425e23859c19891311857dd6e1d065fa6f))

## [0.13.11](https://github.com/cloudquery/plugin-sdk/compare/v0.13.10...v0.13.11) (2022-10-19)


### Features

* Validate source plugin table and column names ([#302](https://github.com/cloudquery/plugin-sdk/issues/302)) ([718314e](https://github.com/cloudquery/plugin-sdk/commit/718314efccaa5ffb23175eced2396387dcb7195f))

## [0.13.10](https://github.com/cloudquery/plugin-sdk/compare/v0.13.9...v0.13.10) (2022-10-19)


### Bug Fixes

* Remove descriptions from table docs ([#300](https://github.com/cloudquery/plugin-sdk/issues/300)) ([6dd529e](https://github.com/cloudquery/plugin-sdk/commit/6dd529ef177d91a6ba0f6a54dcc2c701d7612be6))

## [0.13.9](https://github.com/cloudquery/plugin-sdk/compare/v0.13.8...v0.13.9) (2022-10-16)


### Bug Fixes

* Use 'source' in error message instead of 'destination' ([#295](https://github.com/cloudquery/plugin-sdk/issues/295)) ([7abc547](https://github.com/cloudquery/plugin-sdk/commit/7abc5470247554db9c2e19fc012657e421f7de44))

## [0.13.8](https://github.com/cloudquery/plugin-sdk/compare/v0.13.7...v0.13.8) (2022-10-14)


### Features

* Support application level protocol message. ([#294](https://github.com/cloudquery/plugin-sdk/issues/294)) ([3e1492b](https://github.com/cloudquery/plugin-sdk/commit/3e1492b7ff8855d983262ecbb00eb38a78f3ab69))


### Bug Fixes

* **tests:** Parallel plugin testing, remove old faker ([#292](https://github.com/cloudquery/plugin-sdk/issues/292)) ([48f953a](https://github.com/cloudquery/plugin-sdk/commit/48f953ae0f60a460ea64c4ec35051c48de66faa6))

## [0.13.7](https://github.com/cloudquery/plugin-sdk/compare/v0.13.6...v0.13.7) (2022-10-13)


### Features

* **tests:** More faker options ([#287](https://github.com/cloudquery/plugin-sdk/issues/287)) ([7219478](https://github.com/cloudquery/plugin-sdk/commit/7219478ee1223b1f55eb2f59963d0c48558fe1ae))

## [0.13.6](https://github.com/cloudquery/plugin-sdk/compare/v0.13.5...v0.13.6) (2022-10-12)


### Bug Fixes

* Fix sentry check for development environment ([#285](https://github.com/cloudquery/plugin-sdk/issues/285)) ([151a536](https://github.com/cloudquery/plugin-sdk/commit/151a536196542c60d951597c8aedd18a6d47c545))

## [0.13.5](https://github.com/cloudquery/plugin-sdk/compare/v0.13.4...v0.13.5) (2022-10-12)


### Features

* Add links to tables in table README.md, and list of relations ([#283](https://github.com/cloudquery/plugin-sdk/issues/283)) ([fcfaa42](https://github.com/cloudquery/plugin-sdk/commit/fcfaa422917be8ae4544802558ae799f5a5573c2))

## [0.13.4](https://github.com/cloudquery/plugin-sdk/compare/v0.13.3...v0.13.4) (2022-10-11)


### Bug Fixes

* Tests ([#281](https://github.com/cloudquery/plugin-sdk/issues/281)) ([983e57b](https://github.com/cloudquery/plugin-sdk/commit/983e57b8bf2979be45889ff483510754481ae7fe))

## [0.13.3](https://github.com/cloudquery/plugin-sdk/compare/v0.13.2...v0.13.3) (2022-10-11)


### Bug Fixes

* Call Release on resource semaphore ([#279](https://github.com/cloudquery/plugin-sdk/issues/279)) ([051e247](https://github.com/cloudquery/plugin-sdk/commit/051e24710b64672b4fa4eda1261e2558859cbc75))

## [0.13.2](https://github.com/cloudquery/plugin-sdk/compare/v0.13.1...v0.13.2) (2022-10-11)


### Bug Fixes

* Remove DisallowUnknownFields from Source plugin server ([#277](https://github.com/cloudquery/plugin-sdk/issues/277)) ([0fcf813](https://github.com/cloudquery/plugin-sdk/commit/0fcf813141c82049bd09414fd005d0ff6bbd0b54))

## [0.13.1](https://github.com/cloudquery/plugin-sdk/compare/v0.13.0...v0.13.1) (2022-10-10)


### Bug Fixes

* Ignore Sentry errors in dev (make comparison case insensitive) ([#273](https://github.com/cloudquery/plugin-sdk/issues/273)) ([87ca430](https://github.com/cloudquery/plugin-sdk/commit/87ca430b5855efd3a0f2ad42088aba6ad0e6ae79))
* Ignore sentry in development, case-insensitive for source plugins ([#275](https://github.com/cloudquery/plugin-sdk/issues/275)) ([e2acf4c](https://github.com/cloudquery/plugin-sdk/commit/e2acf4c7200f7f883283c7bb0bd5b88f9382088c))
* Make concurrency change backwards-compatible ([#271](https://github.com/cloudquery/plugin-sdk/issues/271)) ([59ac17a](https://github.com/cloudquery/plugin-sdk/commit/59ac17a4e4cbd3c2a069130fc38eadc29507aafb))

## [0.13.0](https://github.com/cloudquery/plugin-sdk/compare/v0.12.10...v0.13.0) (2022-10-10)


### ⚠ BREAKING CHANGES

* Support table_concurrency and resource_concurrency (#268)

### Features

* Support table_concurrency and resource_concurrency ([#268](https://github.com/cloudquery/plugin-sdk/issues/268)) ([7717d6f](https://github.com/cloudquery/plugin-sdk/commit/7717d6fff5b77f26e2b9ad23859ae03e73e93815))


### Bug Fixes

* Add custom log reader implementation to fix hang on long log lines ([#263](https://github.com/cloudquery/plugin-sdk/issues/263)) ([f8ca238](https://github.com/cloudquery/plugin-sdk/commit/f8ca23838459a67ebb98a6e6f24f954121069f32))
* DeleteStale feature ([#269](https://github.com/cloudquery/plugin-sdk/issues/269)) ([837c5f3](https://github.com/cloudquery/plugin-sdk/commit/837c5f3a56d640dd2ab626ff83d6a540dee4ba08))

## [0.12.10](https://github.com/cloudquery/plugin-sdk/compare/v0.12.9...v0.12.10) (2022-10-09)


### Bug Fixes

* Add missing defer ([#260](https://github.com/cloudquery/plugin-sdk/issues/260)) ([1ee7829](https://github.com/cloudquery/plugin-sdk/commit/1ee782901f6e8499b852af1bd6057aacd1ca7429))

## [0.12.9](https://github.com/cloudquery/plugin-sdk/compare/v0.12.8...v0.12.9) (2022-10-07)


### Bug Fixes

* Bug where first resource wasn't insert into DB ([#258](https://github.com/cloudquery/plugin-sdk/issues/258)) ([2f5b78d](https://github.com/cloudquery/plugin-sdk/commit/2f5b78d8354f11839ac6117d80e29f98562b0b74))

## [0.12.8](https://github.com/cloudquery/plugin-sdk/compare/v0.12.7...v0.12.8) (2022-10-07)


### Bug Fixes

* Error on incorrect table configuration ([#237](https://github.com/cloudquery/plugin-sdk/issues/237)) ([6ad75f5](https://github.com/cloudquery/plugin-sdk/commit/6ad75f53c8f9a632d8f68f04bf4ebc3d4e72f795))
* Exit gracefully on context cancelled ([#252](https://github.com/cloudquery/plugin-sdk/issues/252)) ([b4df92e](https://github.com/cloudquery/plugin-sdk/commit/b4df92e837dd9d892c43b0b02b7c37ed25d573c8))
* Progressbar should go into stdout ([#250](https://github.com/cloudquery/plugin-sdk/issues/250)) ([b8bcdad](https://github.com/cloudquery/plugin-sdk/commit/b8bcdadca19fc1e71ece541d3dbcb3011d8372c7))
* Recover panic in table resolver and object resolver flows ([#257](https://github.com/cloudquery/plugin-sdk/issues/257)) ([04dba02](https://github.com/cloudquery/plugin-sdk/commit/04dba024dd242c169920a15805bc5217e9e446fb))
* Stop if PreResourceResolver fails ([#251](https://github.com/cloudquery/plugin-sdk/issues/251)) ([ee83f8f](https://github.com/cloudquery/plugin-sdk/commit/ee83f8f5e4c03ac421e0a5f0d07a21a5cfd63deb))

## [0.12.7](https://github.com/cloudquery/plugin-sdk/compare/v0.12.6...v0.12.7) (2022-10-05)


### Bug Fixes

* Make progressbar work on small screens ([#248](https://github.com/cloudquery/plugin-sdk/issues/248)) ([7395250](https://github.com/cloudquery/plugin-sdk/commit/73952506b6f7666be44390b6040e8b194ae73214))

## [0.12.6](https://github.com/cloudquery/plugin-sdk/compare/v0.12.5...v0.12.6) (2022-10-05)


### Bug Fixes

* Plugin connection using Unix Domain Socket fixed for windows ([#246](https://github.com/cloudquery/plugin-sdk/issues/246)) ([9e30c60](https://github.com/cloudquery/plugin-sdk/commit/9e30c60cbf0f4136354382fed1f4252c39f52349))

## [0.12.5](https://github.com/cloudquery/plugin-sdk/compare/v0.12.4...v0.12.5) (2022-10-04)


### Bug Fixes

* Logging level ([#243](https://github.com/cloudquery/plugin-sdk/issues/243)) ([d49c44e](https://github.com/cloudquery/plugin-sdk/commit/d49c44e13deba0cf3be27f3b0d64038453ed9ef8))

## [0.12.4](https://github.com/cloudquery/plugin-sdk/compare/v0.12.3...v0.12.4) (2022-10-04)


### Bug Fixes

* Improve download message ([#240](https://github.com/cloudquery/plugin-sdk/issues/240)) ([7929bbb](https://github.com/cloudquery/plugin-sdk/commit/7929bbb7b4492305d420b75265d0721c19546a2d))
* Race condition in log streaming ([#242](https://github.com/cloudquery/plugin-sdk/issues/242)) ([3c8242a](https://github.com/cloudquery/plugin-sdk/commit/3c8242a72e0ee8ffb7fe882c3e8d383bbee6932c))

## [0.12.3](https://github.com/cloudquery/plugin-sdk/compare/v0.12.2...v0.12.3) (2022-10-04)


### Bug Fixes

* Add progressbar instead of writers for Downloads ([#238](https://github.com/cloudquery/plugin-sdk/issues/238)) ([8666d06](https://github.com/cloudquery/plugin-sdk/commit/8666d060785915387bf0a7253e5934f3f2277bce))

## [0.12.2](https://github.com/cloudquery/plugin-sdk/compare/v0.12.1...v0.12.2) (2022-10-04)


### Bug Fixes

* **deps:** Update module github.com/bradleyjkemp/cupaloy/v2 to v2.8.0 ([#215](https://github.com/cloudquery/plugin-sdk/issues/215)) ([a1e444c](https://github.com/cloudquery/plugin-sdk/commit/a1e444c0939616d88fe7507394a8864a03c90ed7))

## [0.12.1](https://github.com/cloudquery/plugin-sdk/compare/v0.12.0...v0.12.1) (2022-10-03)


### Bug Fixes

* SDK compile error, and add workflow ([#234](https://github.com/cloudquery/plugin-sdk/issues/234)) ([6ab1dc2](https://github.com/cloudquery/plugin-sdk/commit/6ab1dc24c683bdfc438e541e285567ae6201df68))

## [0.12.0](https://github.com/cloudquery/plugin-sdk/compare/v0.11.7...v0.12.0) (2022-10-03)


### ⚠ BREAKING CHANGES

* Add overwrite-delete-stale mode for destination plugins (#224)

### Features

* Add overwrite-delete-stale mode for destination plugins ([#224](https://github.com/cloudquery/plugin-sdk/issues/224)) ([567121d](https://github.com/cloudquery/plugin-sdk/commit/567121d680643024bab07988926b46dfbdfbfba6))

## [0.11.7](https://github.com/cloudquery/plugin-sdk/compare/v0.11.6...v0.11.7) (2022-10-03)


### Bug Fixes

* Set default download directory to `.cq` ([#230](https://github.com/cloudquery/plugin-sdk/issues/230)) ([689f5ed](https://github.com/cloudquery/plugin-sdk/commit/689f5ed0299d69498829fbe96c409f7ef86c8757))
* Use correct binary path on Windows ([#231](https://github.com/cloudquery/plugin-sdk/issues/231)) ([0a5dc26](https://github.com/cloudquery/plugin-sdk/commit/0a5dc262c5665fe2253cc5eb26c1b05d250e6b06))

## [0.11.6](https://github.com/cloudquery/plugin-sdk/compare/v0.11.5...v0.11.6) (2022-10-03)


### Bug Fixes

* Download destinations to 'destination' directory ([#228](https://github.com/cloudquery/plugin-sdk/issues/228)) ([d6ebfc3](https://github.com/cloudquery/plugin-sdk/commit/d6ebfc3207c6d0139d5889247754a1a6a4381391))

## [0.11.5](https://github.com/cloudquery/plugin-sdk/compare/v0.11.4...v0.11.5) (2022-10-03)


### Bug Fixes

* Create doc directory if doesn't exist ([#220](https://github.com/cloudquery/plugin-sdk/issues/220)) ([067534d](https://github.com/cloudquery/plugin-sdk/commit/067534d11afee1b39c4a54578b666cc487e12148))
* **deps:** Update golang.org/x/exp digest to 540bb73 ([#212](https://github.com/cloudquery/plugin-sdk/issues/212)) ([2e3dae3](https://github.com/cloudquery/plugin-sdk/commit/2e3dae3490eb89b5be6e6d8733edd2d269960aee))
* **deps:** Update golang.org/x/sync digest to 8fcdb60 ([#213](https://github.com/cloudquery/plugin-sdk/issues/213)) ([7d7d85f](https://github.com/cloudquery/plugin-sdk/commit/7d7d85fc1cede872eba31e02e2c0009d9f903d00))
* Remove redundant error print ([#226](https://github.com/cloudquery/plugin-sdk/issues/226)) ([9927ede](https://github.com/cloudquery/plugin-sdk/commit/9927ede89787f99f517d8883bed1e1383ee32a76))
* Remove unused docs template function ([#221](https://github.com/cloudquery/plugin-sdk/issues/221)) ([f65f023](https://github.com/cloudquery/plugin-sdk/commit/f65f02386a71529cad7a6dcd004d809c2a54ccb9))
* Use correct path for Windows zip ([#223](https://github.com/cloudquery/plugin-sdk/issues/223)) ([960f650](https://github.com/cloudquery/plugin-sdk/commit/960f650cf724175c7014bbef54f9b00c99f8a62d))

## [0.11.4](https://github.com/cloudquery/plugin-sdk/compare/v0.11.3...v0.11.4) (2022-10-01)


### Features

* Add Close() to destination interface and new writemode ([#211](https://github.com/cloudquery/plugin-sdk/issues/211)) ([8af6fcb](https://github.com/cloudquery/plugin-sdk/commit/8af6fcb3dda8e3b17626eb8783bd45dd4ca3fc68))

## [0.11.3](https://github.com/cloudquery/plugin-sdk/compare/v0.11.2...v0.11.3) (2022-09-30)


### Features

* **schema:** Add schema.TypeDuration ([#205](https://github.com/cloudquery/plugin-sdk/issues/205)) ([02fce2c](https://github.com/cloudquery/plugin-sdk/commit/02fce2c8dbdd66ba4e1ee38bf4c7ac61461a8bf8))

## [0.11.2](https://github.com/cloudquery/plugin-sdk/compare/v0.11.1...v0.11.2) (2022-09-30)


### Features

* Make NewSourceClient (and dest) as one interface ([#208](https://github.com/cloudquery/plugin-sdk/issues/208)) ([841a81b](https://github.com/cloudquery/plugin-sdk/commit/841a81b70395bb339a0b30460925b8b35119370b))

## [0.11.1](https://github.com/cloudquery/plugin-sdk/compare/v0.11.0...v0.11.1) (2022-09-30)


### Features

* Move NewSourceClientSpawn (and dest) to sdk ([#206](https://github.com/cloudquery/plugin-sdk/issues/206)) ([15754f9](https://github.com/cloudquery/plugin-sdk/commit/15754f99b4163eed3663c45daa22483075b87828))

## [0.11.0](https://github.com/cloudquery/plugin-sdk/compare/v0.10.2...v0.11.0) (2022-09-29)


### ⚠ BREAKING CHANGES

* Avoid using global variables in caser (#196)
* Remove ParentIDResolver (#202)
* Rename ParentResourceFieldResolver to ParentColumnResolver (#203)
* Make CQUUIDResolver private (#201)
* Remove ParentPathResolver (#200)

### Features

* Make CQUUIDResolver private ([#201](https://github.com/cloudquery/plugin-sdk/issues/201)) ([d879dca](https://github.com/cloudquery/plugin-sdk/commit/d879dca35b6a279f39938500f906687e5b552dbd))
* Remove ParentIDResolver ([#202](https://github.com/cloudquery/plugin-sdk/issues/202)) ([5ae38d0](https://github.com/cloudquery/plugin-sdk/commit/5ae38d0156bedade74c9aada9001664487ef290c))
* Remove ParentPathResolver ([#200](https://github.com/cloudquery/plugin-sdk/issues/200)) ([d839b2f](https://github.com/cloudquery/plugin-sdk/commit/d839b2f4d79f67e1969ea390b43423987a2ecd4d))
* Rename ParentResourceFieldResolver to ParentColumnResolver ([#203](https://github.com/cloudquery/plugin-sdk/issues/203)) ([77d515b](https://github.com/cloudquery/plugin-sdk/commit/77d515bbc369883e09ce441afa9f81f5e5155ad9))


### Bug Fixes

* Add initialisms for k8s ([#191](https://github.com/cloudquery/plugin-sdk/issues/191)) ([5c52157](https://github.com/cloudquery/plugin-sdk/commit/5c521571ed0c136e1f3ad197a10fbb9bd2462428))
* Avoid using global variables in caser ([#196](https://github.com/cloudquery/plugin-sdk/issues/196)) ([85fd56a](https://github.com/cloudquery/plugin-sdk/commit/85fd56a484bd96e2d52730fa0acc61340db6569e))

## [0.10.2](https://github.com/cloudquery/plugin-sdk/compare/v0.10.1...v0.10.2) (2022-09-28)


### Bug Fixes

* Streaming to destination plugin wasn't implemented correctly ([#187](https://github.com/cloudquery/plugin-sdk/issues/187)) ([8e28bd1](https://github.com/cloudquery/plugin-sdk/commit/8e28bd17283a4f039bf34031cf0e01e5c94ac18f))
* Validate versions only for github registry ([#188](https://github.com/cloudquery/plugin-sdk/issues/188)) ([7f9a3ba](https://github.com/cloudquery/plugin-sdk/commit/7f9a3ba8dc31eb6a65ae75bbb5cdc6f563e77aea))

## [0.10.1](https://github.com/cloudquery/plugin-sdk/compare/v0.10.0...v0.10.1) (2022-09-27)


### Bug Fixes

* Add SetDefault and Validate to SpecReader ([#185](https://github.com/cloudquery/plugin-sdk/issues/185)) ([d90acaf](https://github.com/cloudquery/plugin-sdk/commit/d90acaf59f7b803cb814214ac90301e2fd77b4c6))

## [0.10.0](https://github.com/cloudquery/plugin-sdk/compare/v0.9.2...v0.10.0) (2022-09-27)


### ⚠ BREAKING CHANGES

* SpecReader to support multiple files, dirs and yaml (#183)

### Features

* SpecReader to support multiple files, dirs and yaml ([#183](https://github.com/cloudquery/plugin-sdk/issues/183)) ([2531708](https://github.com/cloudquery/plugin-sdk/commit/2531708540298570a9d9711f05abc2d73cc34ddb))

## [0.9.2](https://github.com/cloudquery/plugin-sdk/compare/v0.9.1...v0.9.2) (2022-09-27)


### Bug Fixes

* Spec unmarshalling now supports defaults and validation ([#181](https://github.com/cloudquery/plugin-sdk/issues/181)) ([ba9128a](https://github.com/cloudquery/plugin-sdk/commit/ba9128abea487e67793fa7115376833a938e084e))

## [0.9.1](https://github.com/cloudquery/plugin-sdk/compare/v0.9.0...v0.9.1) (2022-09-27)


### Bug Fixes

* Added custom toCamel, toSnake implementation ([#171](https://github.com/cloudquery/plugin-sdk/issues/171)) ([f28e208](https://github.com/cloudquery/plugin-sdk/commit/f28e20811989abcbe567c9b9ee4420b15667a316))
* **cli:** Added more informative error when there is no config files ([#179](https://github.com/cloudquery/plugin-sdk/issues/179)) ([a7ab327](https://github.com/cloudquery/plugin-sdk/commit/a7ab3276f0890424352360db88b7e571c08fa252))

## [0.9.0](https://github.com/cloudquery/plugin-sdk/compare/v0.8.2...v0.9.0) (2022-09-25)


### ⚠ BREAKING CHANGES

* Make destinations work over gRPC only (#174)

### Bug Fixes

* Make destinations work over gRPC only ([#174](https://github.com/cloudquery/plugin-sdk/issues/174)) ([26237c3](https://github.com/cloudquery/plugin-sdk/commit/26237c357d416f3dda61e51f970660a73c05b0a6))

## [0.8.2](https://github.com/cloudquery/plugin-sdk/compare/v0.8.1...v0.8.2) (2022-09-23)


### Bug Fixes

* Fix typo in ValueTypeFromString ([#172](https://github.com/cloudquery/plugin-sdk/issues/172)) ([12cb9c9](https://github.com/cloudquery/plugin-sdk/commit/12cb9c9b9ee24dd0282da35926229d1256f11696))

## [0.8.1](https://github.com/cloudquery/plugin-sdk/compare/v0.8.0...v0.8.1) (2022-09-22)


### Features

* **codegen:** Add `WithResolverTransformer` option ([#164](https://github.com/cloudquery/plugin-sdk/issues/164)) ([9529956](https://github.com/cloudquery/plugin-sdk/commit/95299560af85a687d1e7274ab80541e02948980a))

## [0.8.0](https://github.com/cloudquery/plugin-sdk/compare/v0.7.13...v0.8.0) (2022-09-22)


### ⚠ BREAKING CHANGES

* Remove ExampleConfig from client,servers and protobuf (#167)

### Features

* Remove ExampleConfig from client,servers and protobuf ([#167](https://github.com/cloudquery/plugin-sdk/issues/167)) ([23b1575](https://github.com/cloudquery/plugin-sdk/commit/23b15758158318b0bfbad78a344a5d4e2369cf98))

## [0.7.13](https://github.com/cloudquery/plugin-sdk/compare/v0.7.12...v0.7.13) (2022-09-21)


### Features

* Ignore specific field types ([#163](https://github.com/cloudquery/plugin-sdk/issues/163)) ([792c88a](https://github.com/cloudquery/plugin-sdk/commit/792c88ab76bac2ce1495931bb4980271a7210051))

## [0.7.12](https://github.com/cloudquery/plugin-sdk/compare/v0.7.11...v0.7.12) (2022-09-21)


### Features

* **codegen:** Add WithTypeTransformer ([#157](https://github.com/cloudquery/plugin-sdk/issues/157)) ([714e5c8](https://github.com/cloudquery/plugin-sdk/commit/714e5c8103c1f771ef95cecfb2cdb5306736f94f))

## [0.7.11](https://github.com/cloudquery/plugin-sdk/compare/v0.7.10...v0.7.11) (2022-09-21)


### Features

* Test that JSON columns don't receive string values ([#156](https://github.com/cloudquery/plugin-sdk/issues/156)) ([d730fdb](https://github.com/cloudquery/plugin-sdk/commit/d730fdb0969912c096fbbd23691eca0bac5121bd))

## [0.7.10](https://github.com/cloudquery/plugin-sdk/compare/v0.7.9...v0.7.10) (2022-09-21)


### Features

* Add support for IgnoreInTests for columns during codegen ([#153](https://github.com/cloudquery/plugin-sdk/issues/153)) ([ec84ddf](https://github.com/cloudquery/plugin-sdk/commit/ec84ddf0d4697d2748eaab4e8197891daa637b4c))

## [0.7.9](https://github.com/cloudquery/plugin-sdk/compare/v0.7.8...v0.7.9) (2022-09-20)


### Features

* Make default transformer exported to use in custom transformers ([#151](https://github.com/cloudquery/plugin-sdk/issues/151)) ([bc93c52](https://github.com/cloudquery/plugin-sdk/commit/bc93c52c0f97584b17398a478206c02f4425c56c))
* make default transformer visible to use in custom transformers ([bc93c52](https://github.com/cloudquery/plugin-sdk/commit/bc93c52c0f97584b17398a478206c02f4425c56c))

## [0.7.8](https://github.com/cloudquery/plugin-sdk/compare/v0.7.7...v0.7.8) (2022-09-20)


### Bug Fixes

* Print correct number of table resources ([#143](https://github.com/cloudquery/plugin-sdk/issues/143)) ([bcbd2a2](https://github.com/cloudquery/plugin-sdk/commit/bcbd2a29ac3e8bd4573042ee526e1292289dd525))

## [0.7.7](https://github.com/cloudquery/plugin-sdk/compare/v0.7.6...v0.7.7) (2022-09-20)


### Features

* Add information about relations to generated docs ([#142](https://github.com/cloudquery/plugin-sdk/issues/142)) ([af77dd7](https://github.com/cloudquery/plugin-sdk/commit/af77dd78b71d1d59f5a9c363b65165439a841e8a))

## [0.7.6](https://github.com/cloudquery/plugin-sdk/compare/v0.7.5...v0.7.6) (2022-09-20)


### Bug Fixes

* Use plugin name to print usage ([#146](https://github.com/cloudquery/plugin-sdk/issues/146)) ([775358c](https://github.com/cloudquery/plugin-sdk/commit/775358ca468537d440b32a58941711563d6649e2))

## [0.7.5](https://github.com/cloudquery/plugin-sdk/compare/v0.7.4...v0.7.5) (2022-09-20)


### Features

* Validate undefined column in TestResource ([#144](https://github.com/cloudquery/plugin-sdk/issues/144)) ([98e8999](https://github.com/cloudquery/plugin-sdk/commit/98e8999fb5923da440d9b6622aed57d3dc9f783b))

## [0.7.4](https://github.com/cloudquery/plugin-sdk/compare/v0.7.3...v0.7.4) (2022-09-20)


### Bug Fixes

* Skip fields that have "-" json tag ([#137](https://github.com/cloudquery/plugin-sdk/issues/137)) ([de4ad3f](https://github.com/cloudquery/plugin-sdk/commit/de4ad3f8df2b64ddd3dba6a5f62df2c7f447a04b))

## [0.7.3](https://github.com/cloudquery/plugin-sdk/compare/v0.7.2...v0.7.3) (2022-09-20)


### Features

* Add PK information to generated docs ([#136](https://github.com/cloudquery/plugin-sdk/issues/136)) ([379d38c](https://github.com/cloudquery/plugin-sdk/commit/379d38c339cb9dc035928211b14c815b3c80a2ef))

## [0.7.2](https://github.com/cloudquery/plugin-sdk/compare/v0.7.1...v0.7.2) (2022-09-19)


### Bug Fixes

* Bring concurrency back ([#129](https://github.com/cloudquery/plugin-sdk/issues/129)) ([04c7f49](https://github.com/cloudquery/plugin-sdk/commit/04c7f4968884cd9430df89815348e63731f91826))

## [0.7.1](https://github.com/cloudquery/plugin-sdk/compare/v0.7.0...v0.7.1) (2022-09-19)


### Features

* Add Multiplexer function type ([#131](https://github.com/cloudquery/plugin-sdk/issues/131)) ([0c72e0c](https://github.com/cloudquery/plugin-sdk/commit/0c72e0ccf4938d492c5478db60b29266dfda5879))

## [0.7.0](https://github.com/cloudquery/plugin-sdk/compare/v0.6.4...v0.7.0) (2022-09-19)


### ⚠ BREAKING CHANGES

* Idiomatic serve interface (#126)

### Features

* Add version flag ([#127](https://github.com/cloudquery/plugin-sdk/issues/127)) ([7e7f1ba](https://github.com/cloudquery/plugin-sdk/commit/7e7f1baaa944ef1d25314b1271a8683b7ae1bd3e))
* Idiomatic serve interface ([#126](https://github.com/cloudquery/plugin-sdk/issues/126)) ([5f848de](https://github.com/cloudquery/plugin-sdk/commit/5f848de294c23dff0890dc1897d55e2e479983cd))
* Use JSON tag for column name when applicable ([#112](https://github.com/cloudquery/plugin-sdk/issues/112)) ([3aa795b](https://github.com/cloudquery/plugin-sdk/commit/3aa795be2852e025866a96acf3a4c1643c6e2022))

## [0.6.4](https://github.com/cloudquery/plugin-sdk/compare/v0.6.3...v0.6.4) (2022-09-18)


### Features

* Make GenerateSourcePluginDocs struct method ([#124](https://github.com/cloudquery/plugin-sdk/issues/124)) ([6597df7](https://github.com/cloudquery/plugin-sdk/commit/6597df73d1297974759209904ef56e1daa793e1d))

## [0.6.3](https://github.com/cloudquery/plugin-sdk/compare/v0.6.2...v0.6.3) (2022-09-18)


### Features

* Add GetDestinations function to list all destinations ([#120](https://github.com/cloudquery/plugin-sdk/issues/120)) ([c4b33fe](https://github.com/cloudquery/plugin-sdk/commit/c4b33fe80c4259cdd84b282ad664a95dce9f14bf))

## [0.6.2](https://github.com/cloudquery/plugin-sdk/compare/v0.6.1...v0.6.2) (2022-09-16)


### Bug Fixes

* Improve error message on codegen field error ([#115](https://github.com/cloudquery/plugin-sdk/issues/115)) ([f31bcec](https://github.com/cloudquery/plugin-sdk/commit/f31bcec69cf750db67b110cacc9213dea4ae3197))

## [0.6.1](https://github.com/cloudquery/plugin-sdk/compare/v0.6.0...v0.6.1) (2022-09-15)


### Features

* Add option to unwrap embedded structs 1 level down ([#111](https://github.com/cloudquery/plugin-sdk/issues/111)) ([a10efb5](https://github.com/cloudquery/plugin-sdk/commit/a10efb53a39c4688754a925173229594dbc764e7))

## [0.6.0](https://github.com/cloudquery/plugin-sdk/compare/v0.5.2...v0.6.0) (2022-09-15)


### ⚠ BREAKING CHANGES

* Remove withComment for codegen (#108)

### Features

* Remove withComment for codegen ([#108](https://github.com/cloudquery/plugin-sdk/issues/108)) ([d8a5711](https://github.com/cloudquery/plugin-sdk/commit/d8a5711ee7434b8bc887d38782094082af3ebe88))

## [0.5.2](https://github.com/cloudquery/plugin-sdk/compare/v0.5.1...v0.5.2) (2022-09-13)


### Bug Fixes

* Remove old entries from changelog ([#100](https://github.com/cloudquery/plugin-sdk/issues/100)) ([6d9290a](https://github.com/cloudquery/plugin-sdk/commit/6d9290a137e103c4448f01488786963519d9557b))

## [0.5.1](https://github.com/cloudquery/plugin-sdk/compare/v0.5.0...v0.5.1) (2022-09-13)


### Bug Fixes

* **testing:** Validate all tables and relations ([#85](https://github.com/cloudquery/plugin-sdk/issues/85)) ([d979863](https://github.com/cloudquery/plugin-sdk/commit/d9798631d9b5a6d93912bda89b7c3e123aeec251))

## [0.5.0](https://github.com/cloudquery/plugin-sdk/compare/v0.4.2...v0.5.0) (2022-09-13)


### ⚠ BREAKING CHANGES

* Enable var names lint rule and fix issues (#88)
* Disable default completion command (#96)

### Features

* Disable default completion command ([#96](https://github.com/cloudquery/plugin-sdk/issues/96)) ([67fca4b](https://github.com/cloudquery/plugin-sdk/commit/67fca4be000c6e4acee76ee95618bc323558c7c1))


### Bug Fixes

* Enable var names lint rule and fix issues ([#88](https://github.com/cloudquery/plugin-sdk/issues/88)) ([4a752b5](https://github.com/cloudquery/plugin-sdk/commit/4a752b548692659bcf203a5ea8a9d11ab3100d2a))
* Remove description from docs ([#92](https://github.com/cloudquery/plugin-sdk/issues/92)) ([7df58df](https://github.com/cloudquery/plugin-sdk/commit/7df58df426baf11a953ee541c67b00ffb15b6fff))
* Remove empty test and enable some lint rules ([#90](https://github.com/cloudquery/plugin-sdk/issues/90)) ([514fba4](https://github.com/cloudquery/plugin-sdk/commit/514fba49e1e817ec505e66b680aa6d3deb0efe07))

## [0.4.2](https://github.com/cloudquery/plugin-sdk/compare/v0.4.1...v0.4.2) (2022-09-12)


### Features

* Add PostResourceResolver to template ([#95](https://github.com/cloudquery/plugin-sdk/issues/95)) ([1f75b05](https://github.com/cloudquery/plugin-sdk/commit/1f75b052f715ce61f78819074c6d774d1301d919))

## [0.4.1](https://github.com/cloudquery/plugin-sdk/compare/v0.4.0...v0.4.1) (2022-09-12)


### Bug Fixes

* **deps:** Update module github.com/bradleyjkemp/cupaloy to v2.7.0 ([#93](https://github.com/cloudquery/plugin-sdk/issues/93)) ([070b9f1](https://github.com/cloudquery/plugin-sdk/commit/070b9f1d694dd67a15e88087ace92187cf8bd3af))

## [0.4.0](https://github.com/cloudquery/plugin-sdk/compare/v0.3.0...v0.4.0) (2022-09-12)


### ⚠ BREAKING CHANGES

* Enable export lin rule and fix option export (#89)

### Bug Fixes

* Enable export lin rule and fix option export ([#89](https://github.com/cloudquery/plugin-sdk/issues/89)) ([478682a](https://github.com/cloudquery/plugin-sdk/commit/478682a99a108f407da096c8114088a531585584))

## [0.3.0](https://github.com/cloudquery/plugin-sdk/compare/v0.2.9...v0.3.0) (2022-09-11)


### ⚠ BREAKING CHANGES

* Depracate override columns (#86)

### Features

* Depracate override columns ([62e1b16](https://github.com/cloudquery/plugin-sdk/commit/62e1b16c2cbba504144bb7496de4bfe408af12ae))
* Depracate override columns ([#86](https://github.com/cloudquery/plugin-sdk/issues/86)) ([62e1b16](https://github.com/cloudquery/plugin-sdk/commit/62e1b16c2cbba504144bb7496de4bfe408af12ae))

## [0.2.9](https://github.com/cloudquery/plugin-sdk/compare/v0.2.8...v0.2.9) (2022-09-11)


### Bug Fixes

* **deps:** Update module github.com/gofrs/uuid to v4.3.0 ([#82](https://github.com/cloudquery/plugin-sdk/issues/82)) ([dbc0c1a](https://github.com/cloudquery/plugin-sdk/commit/dbc0c1ad852520b196bc8beea57c044deac79f9f))

## [0.2.8](https://github.com/cloudquery/plugin-sdk/compare/v0.2.7...v0.2.8) (2022-09-11)


### Bug Fixes

* **deps:** Update module github.com/google/go-cmp to v0.5.9 ([#81](https://github.com/cloudquery/plugin-sdk/issues/81)) ([478f3ad](https://github.com/cloudquery/plugin-sdk/commit/478f3ad7288cba9ce4dd448a4404b407604465f1))

## [0.2.7](https://github.com/cloudquery/plugin-sdk/compare/v0.2.6...v0.2.7) (2022-09-11)


### Bug Fixes

* Add missing comma when generating relations ([#78](https://github.com/cloudquery/plugin-sdk/issues/78)) ([41172d4](https://github.com/cloudquery/plugin-sdk/commit/41172d421bcda25b52ef5747bfe5b92a89667eba))

## [0.2.6](https://github.com/cloudquery/plugin-sdk/compare/v0.2.5...v0.2.6) (2022-09-07)


### Bug Fixes

* **deps:** Update golang.org/x/sync digest to f12130a ([#76](https://github.com/cloudquery/plugin-sdk/issues/76)) ([fe8aa05](https://github.com/cloudquery/plugin-sdk/commit/fe8aa05664d21bb38c57628abf0eafdec4b1662b))

## [0.2.5](https://github.com/cloudquery/plugin-sdk/compare/v0.2.4...v0.2.5) (2022-09-07)


### Bug Fixes

* Fix typo in GetDestinationByName ([#72](https://github.com/cloudquery/plugin-sdk/issues/72)) ([3671366](https://github.com/cloudquery/plugin-sdk/commit/3671366905a2d0f222291436847bab51225e628a))

## [0.2.4](https://github.com/cloudquery/plugin-sdk/compare/v0.2.3...v0.2.4) (2022-09-07)


### Features

* Revert "feat: Generate full example configs from within SDK" ([#70](https://github.com/cloudquery/plugin-sdk/issues/70)) ([06275b6](https://github.com/cloudquery/plugin-sdk/commit/06275b68b1dd4a6fce22889dc3e0bee8d4ad035b))

## [0.2.3](https://github.com/cloudquery/plugin-sdk/compare/v0.2.2...v0.2.3) (2022-09-07)


### Features

* Generate full example configs from within SDK ([#61](https://github.com/cloudquery/plugin-sdk/issues/61)) ([e4f49e9](https://github.com/cloudquery/plugin-sdk/commit/e4f49e956cccabb2cff768b20cf5a4c8c75d052e))

## [0.2.2](https://github.com/cloudquery/plugin-sdk/compare/v0.2.1...v0.2.2) (2022-09-07)


### Features

* Make logs consistent ([#59](https://github.com/cloudquery/plugin-sdk/issues/59)) ([73fcd58](https://github.com/cloudquery/plugin-sdk/commit/73fcd58bd9e5e37f8b4ff3652d61a5a9b8f5a9c9))

## [0.2.1](https://github.com/cloudquery/plugin-sdk/compare/v0.2.0...v0.2.1) (2022-09-07)


### Features

* Remove IgnoreError and send sentry only on panics ([#60](https://github.com/cloudquery/plugin-sdk/issues/60)) ([7139e55](https://github.com/cloudquery/plugin-sdk/commit/7139e553c9e24b95329643c699ec20541206e8a8))

## [0.2.0](https://github.com/cloudquery/plugin-sdk/compare/v0.1.2...v0.2.0) (2022-09-07)


### ⚠ BREAKING CHANGES

* Remove unused table create options (#57)

### Features

* Remove unused table create options ([#57](https://github.com/cloudquery/plugin-sdk/issues/57)) ([6723465](https://github.com/cloudquery/plugin-sdk/commit/67234651a29d75746800c0730d8e0a3a2d90f0ee))


### Bug Fixes

* Ignore hidden files ([#56](https://github.com/cloudquery/plugin-sdk/issues/56)) ([1732ca1](https://github.com/cloudquery/plugin-sdk/commit/1732ca163b5f06ef890bbeae57320e716d5c3ca4))

## [0.1.2](https://github.com/cloudquery/plugin-sdk/compare/v0.1.1...v0.1.2) (2022-09-06)


### Features

* Generate source plugin docs ([#47](https://github.com/cloudquery/plugin-sdk/issues/47)) ([e00d970](https://github.com/cloudquery/plugin-sdk/commit/e00d9707873d1a42b4eeb3ffcbc4b2ee9544f087))

## [0.1.1](https://github.com/cloudquery/plugin-sdk/compare/v0.1.0...v0.1.1) (2022-09-06)


### Features

* Add custom faker ([#52](https://github.com/cloudquery/plugin-sdk/issues/52)) ([34bef4b](https://github.com/cloudquery/plugin-sdk/commit/34bef4b4ce97b3e40bfcc9116a9382df3d3b0551))
* Add sentry for serve.Serve function ([#54](https://github.com/cloudquery/plugin-sdk/issues/54)) ([c1b508f](https://github.com/cloudquery/plugin-sdk/commit/c1b508f09477b881e8862091254f86bb77c110be))

## [0.1.0](https://github.com/cloudquery/plugin-sdk/compare/v0.0.11...v0.1.0) (2022-09-04)


### ⚠ BREAKING CHANGES

* Logger wasnt passed to source plugin resulting no errors (#49)

### Bug Fixes

* Logger wasnt passed to source plugin resulting no errors ([#49](https://github.com/cloudquery/plugin-sdk/issues/49)) ([b0930e4](https://github.com/cloudquery/plugin-sdk/commit/b0930e4e98e98e634314392b0565cfe26a46ea09))

## [0.0.11](https://github.com/cloudquery/plugin-sdk/compare/v0.0.10...v0.0.11) (2022-09-03)


### Features

* Add PreResourceResolver to accommodate list/detail pattern ([#46](https://github.com/cloudquery/plugin-sdk/issues/46)) ([7afadcc](https://github.com/cloudquery/plugin-sdk/commit/7afadccfb82010675ac2cad955d8b70492669e12))

## [0.0.10](https://github.com/cloudquery/plugin-sdk/compare/v0.0.9...v0.0.10) (2022-09-01)


### Bug Fixes

* Pointers to slice are handled correctly ([#11](https://github.com/cloudquery/plugin-sdk/issues/11)) ([70e59fb](https://github.com/cloudquery/plugin-sdk/commit/70e59fb79d9211cdc60446a5d4f8710a49385354))

## [0.0.9](https://github.com/cloudquery/plugin-sdk/compare/v0.0.8...v0.0.9) (2022-09-01)


### Bug Fixes

* **deps:** Update golang.org/x/sync digest to 7fc1605 ([#33](https://github.com/cloudquery/plugin-sdk/issues/33)) ([b594dd0](https://github.com/cloudquery/plugin-sdk/commit/b594dd09cad9e4f5c208f0c76f15341e651116ae))
* **deps:** Update module github.com/rs/zerolog to v1.28.0 ([#38](https://github.com/cloudquery/plugin-sdk/issues/38)) ([17753ea](https://github.com/cloudquery/plugin-sdk/commit/17753ea5c09151bd24d4d8ca9f1241aaecc14872))
* **deps:** Update module google.golang.org/grpc to v1.49.0 ([#39](https://github.com/cloudquery/plugin-sdk/issues/39)) ([d1e0538](https://github.com/cloudquery/plugin-sdk/commit/d1e0538abcb023cfb0a4dc155f68da8c74c06a0c))

## [0.0.8](https://github.com/cloudquery/plugin-sdk/compare/v0.0.7...v0.0.8) (2022-09-01)


### Features

* Remove Unique constraint (support pks only) ([#41](https://github.com/cloudquery/plugin-sdk/issues/41)) ([7e15a30](https://github.com/cloudquery/plugin-sdk/commit/7e15a302d76d903f43560856b173bd4ef06f1b8e))


### Bug Fixes

* **deps:** Update module golang.org/x/tools to v0.1.12 ([#34](https://github.com/cloudquery/plugin-sdk/issues/34)) ([a7bacfd](https://github.com/cloudquery/plugin-sdk/commit/a7bacfda8543c010033f405e9fa5e7803247c8f3))
* **deps:** Update module google.golang.org/protobuf to v1.28.1 ([#35](https://github.com/cloudquery/plugin-sdk/issues/35)) ([b1b25a1](https://github.com/cloudquery/plugin-sdk/commit/b1b25a13a43c6dc5759cfb423ea33a727c5c0894))

## [0.0.7](https://github.com/cloudquery/plugin-sdk/compare/v0.0.6...v0.0.7) (2022-08-31)


### Features

* Improve gRPC status codes and remove .cq file suffix ([#30](https://github.com/cloudquery/plugin-sdk/issues/30)) ([4d4d987](https://github.com/cloudquery/plugin-sdk/commit/4d4d987ead9d05bb0103cf372d990d4dba11a973))

## [0.0.6](https://github.com/cloudquery/plugin-sdk/compare/v0.0.5...v0.0.6) (2022-08-31)


### Features

* Remove pgx dependency ([#26](https://github.com/cloudquery/plugin-sdk/issues/26)) ([be1f37a](https://github.com/cloudquery/plugin-sdk/commit/be1f37a12d6d034058cb8e10d77319e85b290190))

## [0.0.5](https://github.com/cloudquery/plugin-sdk/compare/v0.0.4...v0.0.5) (2022-08-30)


### Bug Fixes

* When cq_id is pkey remove unique ([#20](https://github.com/cloudquery/plugin-sdk/issues/20)) ([0cf4ff8](https://github.com/cloudquery/plugin-sdk/commit/0cf4ff84a5c55c5f0705f2ada50e18af5e3d8d0a))

## [0.0.4](https://github.com/cloudquery/plugin-sdk/compare/v0.0.3...v0.0.4) (2022-08-30)


### Features

* CloudQuery v2 ([#4](https://github.com/cloudquery/plugin-sdk/issues/4)) ([5ceaad4](https://github.com/cloudquery/plugin-sdk/commit/5ceaad4e1c955c90a767205a6ffa7ab3cbf76508))

## [0.0.3](https://github.com/cloudquery/plugin-sdk/compare/v0.0.2...v0.0.3) (2022-08-11)


### Bug Fixes

* Tests and spec unmarshalling ([#3](https://github.com/cloudquery/plugin-sdk/issues/3)) ([6638a8b](https://github.com/cloudquery/plugin-sdk/commit/6638a8ba421cb430891d572314bd5af25d2c8583))
