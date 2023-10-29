# envdir

![coverage](https://raw.githubusercontent.com/ajgon/envdir/badges/.badges/master/coverage.svg)

Envdir is a simple command, aimed for docker/k8s environments to pass all env variables stored as files, directly to the process.

## Usage

```bash
envdir -d /secrets/dir -f -p -lf json -ll debug mycommand -to -execute -with arguments
```

| Argument | Corresponding ENV   | Default    | Description                                                                                                    |
|----------|---------------------|------------|----------------------------------------------------------------------------------------------------------------|
| `-d`     | `ENVDIR_DIRECTORY`  | `/secrets` | Directory to pick variables from                                                                               |
| `-f`     | `ENVDIR_FAIL`       | `false`    | If `true`, command will fail if directory cannot be accesed. If `false`, directory processing will be ignored. |
| `-p`     | `ENVDIR_PARANOID`   | `false`    | See [How paranoid works](#how-paranoid-works)                                                                  |
| `-lf`    | `ENVDIR_LOG_FORMAT` | `text`     | Format of log lines - either `text` or `json`                                                                  |
| `-ll`    | `ENVDIR_LOG_LEVEL`  | `warn`     | Minimal level of log files to be displayed - either `debug`, `info`, `warn` or `error`                         |

### How paranoid works

When envdir is run in "paranoid" mode (`-p`) only `HOME`, `HOSTNAME`, `PATH`, `PWD`, `TERM`, `TZ` and `UMASK` variables will be passed to subprocess.
Any other env variable needs to stored in env directory. This ensures that no unexpected env var will leak in. With this mode disabled, every exported
variable will be passed to the subcommand.

### Use as container entrypoint

Envdir can be used as a shebang in docker entrypoint file, for example:

```bash
#!/usr/bin/envdir /bin/sh

env
```

Default settings are set to most forgiving, meaning if no `/secrets` (default) directory exist, command won't fail. It will also pass all exported envs
(paranoid is also disabled). If you need to customize the command, use it in exec:

```bash
#!/usr/bin/env /bin/sh

exec /usr/bin/envdir -d /my-env-dir -p "$@"
```

### Adding to docker image

```Dockerfile
FROM alpine:3.18

RUN apk add --no-cache --virtual .build-deps curl \
 && curl -OL "https://github.com/ajgon/envdir/releases/download/v0.1.0/envdir_$(uname -s)_$(uname -m).tar.gz" \
 && tar -xzf "envdir_$(uname -s)_$(uname -m).tar.gz" -C /usr/bin/ envdir \
 && rm -rf "envdir_$(uname -s)_$(uname -m).tar.gz" \
 && apk del .build-deps
```

## Example

```bash
# create some dummy envs for testing
mkdir -p /tmp/env
echo ipsum > /tmp/env/LOREM
echo 42 > /tmp/env/THE_ANSWER
echo true > /tmp/env/DROP_DATABASE

# run in paranoid mode
/usr/bin/envdir -d /tmp/env -p env | sort
DROP_DATABASE=true
HOME=/home/ajgon
HOSTNAME=
LOREM=ipsum
PATH=/usr/local/bin:/usr/local/sbin:/usr/local/bin:/usr/bin
PWD=/tmp
TERM=xterm-256color
THE_ANSWER=42
TZ=Europe/Warsaw
UMASK=
```
