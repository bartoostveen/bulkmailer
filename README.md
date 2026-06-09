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

## Usage

You may create a jobs.json that contains data of the following format:

```json
[
  {
    "recipient": "johndoe@example.com",
    "template": "example_template.txt",
    "extraAttrs": {
      "Foo": "bar"
    }
  },
  {
    "recipient": "anotheruser@example.com",
    "template": "example_template.txt",
    "extraAttrs": {
      "Foo": "baz"
    }
  }
]
```

Then, for all recipients, the template `example_template.txt` gets rendered. Attributes such as `Recipient` are
available, as well as the `extraAttrs` in the `Data` object. Here is an example template that makes use of these fields.

```
Hello, {{.Recipient}}!

Sample data: "{{.Data.Foo}}"
```

You may then invoke the main binary using:

```shell
go run bartoostveen.nl/bulkmailer/cmd/bulkmailer
```

## License

This project is licensed under the GPLv3. You can find a copy in [the LICENSE file](./LICENSE).
