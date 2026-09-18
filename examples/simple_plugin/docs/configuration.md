---
title: Configuration
slug: configuration
description: How to configure the simple_plugin example.
---

# Configuration

Every field is optional.

| Field | Type | Description |
|---|---|---|
| `num_clients` | int | How many multiplexed clients to fetch with. Defaults to `1`. |
| `changes` | list | Canned events replayed over the `Watch` stream, in order. |
| `refresh_data` | map | The rows a refresh sync can find, per table. |

## Event-based sync

`changes` and `refresh_data` drive `sync_mode: events` without a cloud behind
the plugin.

```yaml
spec:
  num_clients: 1
  changes:
    - table: test_some_table
      keys: { column1: "test_project_id", client_id: 0 }
      op: upsert          # upsert (default) or delete
      event_id: "1"
    - table: test_some_table
      keys: { column1: "gone", client_id: 0 }
      op: delete
      event_id: "2"
  refresh_data:
    test_some_table:
      - { column1: "test_project_id", column2: 123, client_id: 0 }
```

- `changes[].keys` must carry every primary key column of the table. For
  `test_some_table` those are `column1` and `client_id`.
- `changes[].event_id` must be unique. Acking one checkpoints every change up to
  and including it, in the state backend under `events:checkpoint` (with
  `events:checkpoint_time` alongside it). A later `Watch` with no explicit
  position resumes after that event.
- `refresh_data` rows must also carry every primary key column. A refresh for
  keys that match a row emits an insert for it; keys that match nothing emit a
  delete record.
