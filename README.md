# Bulk mailer

A very primitive template renderer and sender, written in Go.

## Configuration

You can invoke this using a `config.yaml`. Either make sure this file is present in the current working directory, or
specify it using the `-c` flag.

```yaml
from: Bart Oostveen <bart@bartoostveen.nl>
reply_to: Bart Oostveen <bart@bartoostveen.nl>
subject: Reminder!
target_dir: target
smtp:
  host: mx.bartoostveen.nl
  port: 587
  username: bart@bartoostveen.nl
  password: verysecretpassword
dry: true
jobs_file: example_jobs.json
unique_file_names: false
```

The same configuration options exist as environment variables with the `BULKMAILER` prefix, e.g. the smtp password
translates to `BULKMAILER_SMTP_PASSWORD`.

## License

This project is licensed under the GPLv3. You can find a copy in [the LICENSE file](./LICENSE).
