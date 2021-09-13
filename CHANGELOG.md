# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).


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
