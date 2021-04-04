module github.com/cloudquery/cq-provider-sdk

go 1.15

require (
	github.com/Masterminds/squirrel v1.5.0
	github.com/creasty/defaults v1.5.1
	github.com/georgysavva/scany v0.2.7
	github.com/google/uuid v1.2.0
	github.com/hashicorp/go-hclog v0.15.0
	github.com/huandu/go-sqlbuilder v1.12.0
	github.com/iancoleman/strcase v0.1.2
	github.com/jackc/pgproto3/v2 v2.0.7 // indirect
	github.com/jackc/pgx/v4 v4.10.1
	github.com/lib/pq v1.8.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1
	github.com/stretchr/testify v1.7.0
	github.com/thoas/go-funk v0.8.0
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2 // indirect
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/thoas/go-funk => github.com/cloudquery/go-funk v0.8.1-0.20210404121448-4d824a7058bc