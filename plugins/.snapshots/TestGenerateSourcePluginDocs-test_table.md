# Table: test_table
Description for test table

The composite primary key for this table is (**id_col**, **id_col2**).

## Relations
The following tables depend on `test_table`:
  - [`relation_table`](relation_table.md)

## Columns
| Name          | Type          |
| ------------- | ------------- |
|int_col|Int|
|id_col (PK)|Int|
|id_col2 (PK)|Int|
|_cq_id|UUID|
|_cq_fetch_time|Timestamp|
