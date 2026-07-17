# Tinfoil Proxy

A tiny proxy that provides an OpenAI-compatible API endpoint, forwarding
requests to a [Tinfoil](https://tinfoil.sh/) secure enclave.

This gives you the best of both worlds:

- A connection to a secure inference enclave with locally verified attestation
- An OpenAI-compatible endpoint so your client doesn't need any changes

Tinfoil Proxy is a thin wrapper over the
[Tinfoil Go SDK](https://github.com/tinfoilsh/tinfoil-go). As described in the
[Tinfoil docs](https://docs.tinfoil.sh/sdk/overview#direct-api-access-not-recommended),
using the Tinfoil SDK is more secure than calling Tinfoil's OpenAI-compatible
endpoint directly because the SDK verifies the attestation locally, ensuring a
trusted end-to-end encrypted channel to a secure enclave.

Tinfoil Proxy authenticates incoming requests with `TINFOIL_PROXY_API_KEY`, a
key of your choosing. It then proxies them securely to a Tinfoil enclave using
`TINFOIL_API_KEY`, which you need to
[obtain from Tinfoil](https://dash.tinfoil.sh/?tab=api-keys).

You can also set `TINFOIL_PROXY_API_KEY_PATH` and `TINFOIL_API_KEY_PATH` to file
paths containing the respective keys. This is more secure than using env vars
directly, and is the preferred way to supply credentials when running Tinfoil
Proxy as a systemd unit.

## Endpoints

| Endpoint               | Method | Description                                |
| ---------------------- | ------ | ------------------------------------------ |
| `/v1/chat/completions` | POST   | Chat completions (proxied to enclave)      |
| `/v1/models`           | GET    | List available models (proxied to enclave) |

## Requirements

- Go 1.26+

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
