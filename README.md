# Tinfoil Proxy

A tiny local reverse proxy that provides an OpenAI-compatible API endpoint,
forwarding requests to a [Tinfoil](https://tinfoil.sh/) secure enclave with
verified attestation. Tinfoil Proxy is a thin wrapper over the
[Tinfoil Go SDK](https://github.com/tinfoilsh/tinfoil-go).

As described in the
[Tinfoil docs](https://docs.tinfoil.sh/sdk/overview#direct-api-access-not-recommended),
using the Tinfoil SDK is more secure than calling Tinfoil's HTTPS endpoint
directly because the SDK verifies the attestation locally, ensuring a trusted
end-to-end encrypted channel to a secure enclave.

Tinfoil proxy authenticates requests with your own `TINFOIL_PROXY_API_KEY`. It
then proxies them securely to a Tinfoil enclave using the `TINFOIL_API_KEY`.

You can also provide these two keys as paths `TINFOIL_PROXY_API_KEY_PATH` and
`TINFOIL_API_KEY_PATH` pointing to files containing the keys. This is more
secure than direct env vars and the preferred way to supply credentials when
running Tinfoil Proxy as a systemd unit.

## Endpoints

| Endpoint               | Method | Description                                |
| ---------------------- | ------ | ------------------------------------------ |
| `/v1/chat/completions` | POST   | Chat completions (proxied to enclave)      |
| `/v1/models`           | GET    | List available models (proxied to enclave) |

## Requirements

- Go 1.25+

## Build

```
make build
```

## Run

```
export TINFOIL_API_KEY=<your-tinfoil-api-key>
export TINFOIL_PROXY_API_KEY=<your-proxy-api-key>
./tinfoil-proxy
```

The proxy listens on `127.0.0.1:17349` by default. Override with the `-listen`
flag:

```
./tinfoil-proxy -listen 0.0.0.0:8080
```

## Usage

Point any OpenAI-compatible client at the proxy using your
`TINFOIL_PROXY_API_KEY` as the bearer token:

```
curl http://127.0.0.1:17349/v1/chat/completions \
  -H "Authorization: Bearer $TINFOIL_PROXY_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"model":"...","messages":[{"role":"user","content":"Hello"}]}'
```

## License

AGPL-3.0. See LICENSE file for details.
