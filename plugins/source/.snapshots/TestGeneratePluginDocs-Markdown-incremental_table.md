# Table: incremental_table

Description for incremental table

The primary key for this table is **id_col**.
It supports incremental syncs based on the (**id_col**, **id_col2**) columns.

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_source_name|String|
|_cq_sync_time|Timestamp|
|_cq_id|UUID|
|_cq_parent_id|UUID|
|int_col|Int|
|id_col (PK) (Incremental Key)|Int|
|id_col2 (Incremental Key)|Int|
