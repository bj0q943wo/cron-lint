# cron-lint

Static analyzer for cron expressions that validates schedules and warns about overlapping jobs.

## Installation

```bash
go install github.com/cron-lint/cron-lint@latest
```

## Usage

Run `cron-lint` against a file containing cron expressions (one per line):

```bash
cron-lint schedule.cron
```

**Example input** (`schedule.cron`):
```
*/5 * * * *   backup-job
0 * * * *     hourly-report
*/10 * * * *  sync-data
0 0 * * *     daily-cleanup
```

**Example output**:
```
schedule.cron:1: WARNING: "*/5 * * * *" overlaps with "*/10 * * * *" (every 10 minutes)
schedule.cron:3: OK
schedule.cron:4: OK
```

You can also pipe expressions directly:

```bash
echo "*/5 * * * * my-job" | cron-lint -
```

### Flags

| Flag | Description |
|------|-------------|
| `-strict` | Treat warnings as errors |
| `-format json` | Output results as JSON |
| `-quiet` | Suppress passing entries |

## What It Checks

- Invalid cron expression syntax
- Overlapping job schedules
- Unreachable expressions (e.g., `0 0 31 2 *`)
- Duplicate schedules across entries

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

## License

[MIT](LICENSE)