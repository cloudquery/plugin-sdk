# Table: test_table

Description for test table

The composite primary key for this table is (**id_col**, **id_col2**).

## Relations

The following tables depend on test_table:
  - [relation_table](relation_table.md)
  - [relation_table2](relation_table2.md)

## Columns
| Name          | Type          |
| ------------- | ------------- |
|_cq_source_name|String|
|_cq_sync_time|Timestamp|
|_cq_id|UUID|
|_cq_parent_id|UUID|
|int_col|Int|
|id_col (PK)|Int|
|id_col2 (PK)|Int|
