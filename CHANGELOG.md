# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
