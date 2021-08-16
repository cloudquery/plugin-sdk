# Integration Tests

## Description

Integration tests use terraform to deploy resources. Every resource terraform file should be described in
./resources/testData and have the same name as tested resource(<provider>_<domain>_<resource_name>.tf).
Example: `aws_iam_users.tf`
Testing routine copies this file among with default *.tf files to a separate folder, sets default variables prefix=<
resource_name> suffix=<machine_hostname>, deploys resources, fetches data from provider to database, queries data
using `Filter` field or using default filter based on tags, compares the received data with expected values.

## Run

To run integration tests you need:

- see deployment examples in the cloudquery docs
- terraform executable in $PATH
- provider credentials configured via config files or env variables
- sql database deployed  
  
to run the tests use command below in PROVIDER root dir:
```shell
 go test -v -p 20 ./resources --tags=integration
```
Tests can be marked with `integration_skip` tag. This means are not ready yet for testing.

## Debugging

For debugging, you can set env variable `TF_NO_DESTROY=true` to leave the directory and resources after the test.
Resources should be destroyed manually by running

```
tf destroy -var test_prefix=<value>  -var test_suffix=<value>
```

values for `-var` arguments can be found in test execution output near `<resource name> tf apply` log entry

To avoid long terraform deploys you can set `TF_EXEC_TIMEOUT=<minutes>`