# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.6.0-rc1] - 2021-12-29
### :gear: Changed
* **Breaking Change**: changed column attribute `IgnoreInTests` to `IgnoreInTests` API [#138](https://github.com/cloudquery/cq-provider-sdk/pull/137)

### :rocket: Added
* Added `SkipEmptyColumn` and `SkipEmptyRows` to `ResourceTestCase`
* If test fail it will print what are the missing columns as well.

## [v0.6.0-beta] - 2021-12-29
### :gear: Changed
* **Breaking Change**: changed `TestResource` API [#137](https://github.com/cloudquery/cq-provider-sdk/pull/137)

## [v0.5.7]- 2021-12-20

### :gear: Changed
* Fix table and column name limit tests [#134](https://github.com/cloudquery/cq-provider-sdk/pull/134).

## [v0.5.6] - 2021-12-18

### :gear: Changed
* SDK e2e testing terraform apply now also logs [#130](https://github.com/cloudquery/cq-provider-sdk/pull/130).

### :rocket: Added
* Added new test for table and column name limits [#133](https://github.com/cloudquery/cq-provider-sdk/pull/133).

## [v0.5.5] - 2021-12-15

### :gear: Changed
* Added support for error interface for diagnostics [#128](https://github.com/cloudquery/cq-provider-sdk/pull/128).
* Improved doc generation to remove unused files [#127](https://github.com/cloudquery/cq-provider-sdk/pull/127) fixes [#116](https://github.com/cloudquery/cq-provider-sdk/issues/116).
* Added warning about file descriptor usage [#126](https://github.com/cloudquery/cq-provider-sdk/pull/126) fixes [cloudquery/cloudquery#285](https://github.com/cloudquery/cloudquery/issues/285).

## [v0.5.4] - 2021-12-09

### :spider: Fixed
* fixed bad execution error validation [#125](https://github.com/cloudquery/cq-provider-sdk/pull/125)

## [v0.5.3] - 2021-12-06

### :gear: Changed
* Updated SDK dependencies [#121](https://github.com/cloudquery/cq-provider-sdk/pull/121)
* Add column name to resolver errors [#114](https://github.com/cloudquery/cq-provider-sdk/issues/114)
* Improve plugin serve execution message [#117](https://github.com/cloudquery/cq-provider-sdk/issues/117)


## [v0.5.2] - 2021-11-23

### :rocket: Added
* Support IPAddressesResolver for TypeInetArray [#112](https://github.com/cloudquery/cq-provider-sdk/pull/112)
* []struct now can be parsed automatically to jsons [#109](https://github.com/cloudquery/cq-provider-sdk/issues/109)


## [v0.5.1] - 2021-11-01

### :rocket: Added
 * feat/implementation of parallel clients limit by @fdistorted in [#103](https://github.com/cloudquery/cq-provider-sdk/pull/103)
 * Support passing table meta information over cqproto [#107](https://github.com/cloudquery/cq-provider-sdk/pull/107)

## [v0.5.0] - 2021-10-21

### :rocket: Added
* Support diagnostics from fetch executions, allow providers to define custom diagnostic classifiers for errosr received from the fetch execution [#104](https://github.com/cloudquery/cq-provider-sdk/pull/104)
* Add more metadata sent on fetched resources completion, status, resource count and diagnostics [#104](https://github.com/cloudquery/cq-provider-sdk/pull/104)

## [v0.4.10] - 2021-10-18

Fix drop provider tables due to out of shared memory, a large number of tables exceeds the transaction memory limit of
usual database configurations [#105](https://github.com/cloudquery/cq-provider-sdk/pull/105)
    
## [v0.4.9] - 2021-10-07

### :spider: Fixed
* fixed missing stale filter `--disable-delete` in cloudquery [#102](https://github.com/cloudquery/cq-provider-sdk/pull/102)

## [v0.4.8] - 2021-10-05

### :spider: Fixed
* updated integration test validation, allowing at least 1 results [#101](https://github.com/cloudquery/cq-provider-sdk/pull/101)


## [v0.4.7] - 2021-09-23

### :rocket: Added
* Added support to remove stale data based on `last_updated` column that wasn't fetched in latest refresh, activate with `--disable-delete` in cloudquery [#95](https://github.com/cloudquery/cq-provider-sdk/pull/95)

### :gear: Changed
* Integration tesing should fail if provider has internal error [#98](https://github.com/cloudquery/cq-provider-sdk/pull/98)

### :spider: Fixed
* fixed default resolver for resource valus to be json for unknown types [#99](https://github.com/cloudquery/cq-provider-sdk/pull/99)

## [v0.4.6] - 2021-09-14

### :gear: Changed
* Debugging providers will print debug level by default. Trace enabled via env variable `CQ_PROVIDER_DEBUG_TRACE_LOG` [#93](https://github.com/cloudquery/cq-provider-sdk/pull/93)

## [v0.4.5] - 2021-09-14

### :rocket: Added
* Added support to close migrator connection [#92](https://github.com/cloudquery/cq-provider-sdk/pull/92)


## [v0.4.4] - 2021-09-13

### :spider: Fixed
* fix resource insert logging error, print syntax error SQL on failure [#89](https://github.com/cloudquery/cq-provider-sdk/pull/89)


## [v0.4.3] - 2021-09-06

### :rocket: Added
* Added support to fetch all resources with "*" [#87](https://github.com/cloudquery/cq-provider-sdk/pull/87)

### :gear: Changed
* Partial fetch flag enabled by default on configuration (cq init [provider]) creation for new providers [#87](https://github.com/cloudquery/cq-provider-sdk/pull/87)


## [v0.4.2] - 2021-09-04

### :spider: Fixed
* fix missing table name from partial fetch error [#85](https://github.com/cloudquery/cq-provider-sdk/issues/85)


## [v0.4.1] - 2021-09-02 

### :spider: Fixed
* fix missing database connection url set [#84](https://github.com/cloudquery/cq-provider-sdk/issues/84)


## [v0.4.0] - 2021-09-02

### :rocket: Added
* Added support for partial fetching [#60](https://github.com/cloudquery/cq-provider-sdk/pull/76)


## [v0.3.4] - 2021-08-25

### :spider: Fixed
* fix edge case of migration jumps [#78](https://github.com/cloudquery/cq-provider-sdk/issues/78)


### :rocket: Added
* Added support for provider migrations [#71](https://github.com/cloudquery/cq-provider-sdk/issues/71)

## [v0.3.2] - 2021-08-11

### :spider: Fixed
* Generate random cq_id if some primary keys are null [#63](https://github.com/cloudquery/cq-provider-sdk/issues/63) fixed in [#68](https://github.com/cloudquery/cq-provider-sdk/issues/63) 

### :rocket: Added
* Added support for common resolvers [#61](https://github.com/cloudquery/cq-provider-sdk/issues/61)
    *  IP Resolver
    * INET resolver
    * Mac resolver
    * UUID resolver
    * Datetime Resolver
    * Date Resolver
    * String Transformer
    * Int Transformer

## [v0.3.1] - 2021-07-30

### :spider: Fixed
* Return error on duplicate resources request fixes [#58](https://github.com/cloudquery/cq-provider-sdk/issues/58)
* Add better recovery from panic in resolvers, printing stack and errors in log [#55](https://github.com/cloudquery/cq-provider-sdk/issues/55)

## [v0.3.0] - 2021-07-28

### :rocket: Added

* Added a changelog :)
* Added support for user defined Primary Keys in [#41](https://github.com/cloudquery/cq-provider-sdk/pull/41)
* Added support to disable delete of data [#41](https://github.com/cloudquery/cq-provider-sdk/pull/41)
* Added meta field, meta information on the resource, for example: when resource updated last. [#41](https://github.com/cloudquery/cq-provider-sdk/pull/41)

### :gear: Changed
* Changed default insert in provider from Insert to Copy-From, this method improves insert performance [#48](https://github.com/cloudquery/cq-provider-sdk/pull/48)
* **Breaking Change**: default CloudQuery "id" from `id` to `cq_id` [#41](https://github.com/cloudquery/cq-provider-sdk/pull/41)

## [0.2.8] - 2021-07-15

Base version at which changelog was introduced.
