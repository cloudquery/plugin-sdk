# Table: relation_table

Description for relational table

The primary key for this table is **_cq_id**.

## Relations

This table depends on [test_table](test_table.md).

The following tables depend on relation_table:
  - [relation_relation_table_a](relation_relation_table_a.md)
  - [relation_relation_table_b](relation_relation_table_b.md)

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_source_name|String|
|_cq_sync_time|Timestamp|
|_cq_id (PK)|UUID|
|_cq_parent_id|UUID|
|string_col|String|
