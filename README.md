# nomad-holepunch

Proxy the Nomad API via Workload Identity.

[![MPL License](https://img.shields.io/github/license/shoenig/nomad-holepunch?color=g&style=flat-square)](https://github.com/shoenig/nomad-holepunch/blob/main/LICENSE)
[![Run CI Tests](https://github.com/shoenig/nomad-holepunch/actions/workflows/ci.yaml/badge.svg)](https://github.com/shoenig/nomad-holepunch/actions/workflows/ci.yaml)

## Overview

This `nomad-holepunch` program can be run as a Nomad task that proxies the Nomad API
by making use of Nomad's Workload Identity authorization token and Unix domain socket.

## Configuration

`nomad-holepunch` is configured via environment variables.

| Environment Variable | Description | Default |
| ---------------------|-------------|---------|
| `HOLEPUNCH_BIND` | The TCP address to bind to | `0.0.0.0` |
| `HOLEPUNCH_PORT` | The TCP port to listen on | `6120` |
| `NOMAD_SOCKET_PATH` | The filepath to find the Nomad API `api.sock` file | `$NOMAD_SECRETS_DIR/api.sock` |
| `HOLEPUNCH_ALLOW_ALL` | Allow access to all Nomad endpoints | `false` |
| `HOLEPUNCH_ALLOW_METRICS` | Allow access to Nomad /metrics API endpoint | `true` |
| `HOLEPUNCH_ALLOW_NODES` | Allow access to Nomad /node API endpoint | `false` |
| `HOLEPUNCH_ALLOW_AGENT_HEALTH` | Allow access to Nomad /agent/health endpoint | `true` |
| `HOLEPUNCH_ALLOW_AGENT_SELF` | Allow access to Nomad /agent/self endpoint | `false` |
| `HOLEPUNCH_ALLOW_AGENT_MEMBERS` | Allow access to Nomad /agent/members endpoint | `false` |
| `HOLEPUNCH_ALLOW_AGENT_SERVERS` | Allow access to Nomad /agent/servers endpoint | `false` |
| `HOLEPUNCH_ALLOW_AGENT_HOST` | Allow access to Nomad /agent/host endpoint | `false` |
| `HOLEPUNCH_ALLOW_AGENT_SCHEDULERS` | Allow access to Nomad /agent/schedulers(/config) endpoints | `false` |
| `HOLEPUNCH_ALLOW_PLUGINS` | Allow access to Nomad /plugins endpoint | `false` |
| `HOLEPUNCH_ALLOW_SERVICES` | Allow access to Nomad /service(s/*) endpoints | `false` |
| `HOLEPUNCH_ALLOW_STATUS` | Allow access to Nomad /status/(leader,peers) endpoints | `false` |

## Local Development

For local development, `hack/localdev.hcl` provides a convenient way to run the
`nomad-holepunch` program as a `raw_exec` nomad job. It is assumed that `nomad-holepunch`
is present somewhere on the user's `$PATH` (see Compile), and that Nomad agent
is running.

##### running a nomad agent

```shell-session
sudo nomad agent -dev
```

##### compile and install

into `$GOPATH/bin`, assumed to be on `$PATH`

```shell-session
go install
```

###### run the localdev nomad job

```shell-session
nomad job run -var=user=$USER ./hack/localdev.hcl
```

### Compile

The `Makefile` provides targets for building `nomad-holepunch` locally.

- `make build` (default) - compile and place binary into `output/`

- `make test` - run `go test` to run test cases

### Container

A container is built for every version. They live in the [GitHub Container Registry](https://github.com/shoenig/nomad-holepunch/pkgs/container/nomad-holepunch).

Although it isn't useful to run the container outside of Nomad, it is still possible, e.g.

```shell-session
➜ podman run --rm ghcr.io/shoenig/nomad-holepunch:v0.1.1
2023/06/25 19:23:52 INFO  [main] ^^ startup nomad-holepunch ^^
2023/06/25 19:23:52 TRACE [main] HOLEPUNCH_BIND = 0.0.0.0
2023/06/25 19:23:52 TRACE [main] HOLEPUNCH_PORT = 6120
2023/06/25 19:23:52 TRACE [main] HOLEPUNCH_TOKEN = <redacted>
2023/06/25 19:23:52 TRACE [main] HOLEPUNCH_SOCKET_PATH = /secrets/api.sock
2023/06/25 19:23:52 TRACE [main] HOLEPUNCH_ALLOW_ALL = false
2023/06/25 19:23:52 TRACE [main] HOLEPUNCH_ALLOW_METRICS = true
```
