# Table: incremental_table

This table shows data for Incremental Table.

Description for incremental table

The primary key for this table is **id_col**.
It supports incremental syncs based on the (**id_col**, **id_col2**) columns.

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_source_name|utf8|
|_cq_sync_time|timestamp[us, tz=UTC]|
|_cq_id|uuid|
|_cq_parent_id|uuid|
|int_col|int64|
|id_col (PK) (Incremental Key)|int64|
|id_col2 (Incremental Key)|int64|
