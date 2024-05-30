# corsproxy

`corsproxy` is yet another CORS proxy written in Go, designed to bypass CORS restrictions when making requests from web applications.

## Features

- **Private Network Targets Disallowed**: By default, `corsproxy` disallows private network targets to enhance security.
- **Configurable Allowed Targets**: You can configure the allowed targets to specify which domains are allowed to be accessed through the proxy.

## Installation

```shell
go install github.com/zachcheung/corsproxy/cmd/corsproxy
```

## Usage

```shell
corsproxy -allowedHeaders "Authorization" -allowedTargets "https://*.example.com,https://foo.bar"
```

Please refer to [rs/cors](https://github.com/rs/cors) for detailed information on CORS-related options.

## License

[MIT](LICENSE)
