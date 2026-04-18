# Tinfoil Proxy

A tiny local reverse proxy that provides an OpenAI-compatible API endpoint,
forwarding requests to a [Tinfoil](https://tinfoil.sh/) secure enclave with
verified attestation.

This is more secure than calling Tinfoil's OpenAI compatible API directly
because the proxy verifies the attestation locally, ensuring a trusted
end-to-end encrypted channel to the secure enclave.

Requests are authenticated with your own `TINFOIL_PROXY_API_KEY`, then proxied
securely to the Tinfoil enclave using the `TINFOIL_API_KEY`.

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
