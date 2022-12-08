# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
