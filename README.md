# corsproxy

`corsproxy` is yet another CORS proxy written in Go, designed to bypass CORS restrictions when making requests from web applications.

## Features

- **Private Network Targets Disallowed**: By default, `corsproxy` disallows private network targets to enhance security.
- **Configurable Allowed Targets**: You can configure the allowed targets to specify which domains are allowed to be accessed through the proxy.

## Installation

* [releases](https://github.com/zachcheung/corsproxy/releases)

* docker

```shell
docker pull ghcr.io/zachcheung/corsproxy
```

* go install

```shell
go install github.com/zachcheung/corsproxy/cmd/corsproxy@latest
```

## Usage

```shell
corsproxy -allowedTargets "https://*.example.com,http://ip.me"
```

Please refer to [rs/cors](https://github.com/rs/cors) for detailed information on CORS-related options.

Request examples:

* `http://localhost:8000/http://ip.me`

## License

[MIT](LICENSE)
