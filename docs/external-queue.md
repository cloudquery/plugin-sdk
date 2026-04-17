# External Queue Storage (Preview)

The SDK supports offloading the scheduler's work state to an external queue backend. This is useful for large syncs that would otherwise exhaust RAM — queued work units and parent resources can be spilled to local disk.

## When to use it

Use an external queue when:
- Your sync OOMs or runs near memory limits on its host.
- You're using the `shuffle-queue` scheduler strategy.

Do NOT use it when:
- You're using `dfs`, `round-robin`, or `shuffle` strategies — those cannot use an external queue (they'd need execution-model rewrites).
- Your sync completes well within memory budgets — the backend adds disk I/O overhead.

## Configuration

Add a `queue` block alongside `scheduler` in your source plugin spec:

```yaml
spec:
  scheduler: shuffle-queue
  queue:
    type: badger
    path: /var/lib/cq/queue
```

### Backends

| `type` | Description |
|---|---|
| `in-memory` (default) | Current behavior. Everything held in process memory. |
| `badger` | Embedded BadgerDB on local disk. Required: `path`. |

## Requirements for plugin authors

Plugins must use `transformers.TransformWithStruct` for tables that have relations. This is already universal — no changes needed. If a plugin has a custom `Transform`, call `table.SetItemSample(yourItemType{})` explicitly; otherwise the sync will fail fast at startup with a clear error.

## Caveats

- **No crash recovery.** A crashed sync's queue state is discarded; restart from scratch.
- **No encryption at rest.** Items are stored as plain JSON. Do not configure a disk backend on shared filesystems without filesystem-level encryption if you sync sensitive data.
- **Per-invocation path isolation.** The actual Badger directory is `{queue.path}/{invocation_id}`, so multiple concurrent syncs of the same plugin don't collide. Stale directories from crashed syncs are left for the user to clean up.
- **No remote backends in v1.** Redis/Postgres/SQS are deferred. The contract-test architecture makes adding them a small project.

## Troubleshooting

- `queue backend already locked by another process` → a prior sync or an orphaned Badger instance holds the directory lock. Remove the stale directory or pass a fresh `queue.path`.
- `queue: table "X" has relations but no itemSample` → the plugin has a table with relations that doesn't use `TransformWithStruct`. Add `table.SetItemSample(XItemType{})` in the table definition.
- Disk fills up mid-sync → v1 doesn't bound disk usage. Point `queue.path` at a volume with generous headroom; capacity ≈ number-of-resources × relation-depth × average-resource-size.
