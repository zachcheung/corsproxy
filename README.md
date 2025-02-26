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
corsproxy -allowedTargets "https://*.example.com,https://ipinfo.io"
```

Please refer to [rs/cors](https://github.com/rs/cors) for detailed information on CORS-related options.

Request examples:

* `http://localhost:8000/https://ipinfo.io/json`

Live examples:

* https://corsproxy-c6ln.onrender.com/ - render.com free web service tier

```console
$ curl https://corsproxy-c6ln.onrender.com/https://ipinfo.io/json
{
  "ip": "44.227.217.144",
  "hostname": "ec2-44-227-217-144.us-west-2.compute.amazonaws.com",
  "city": "Boardman",
  "region": "Oregon",
  "country": "US",
  "loc": "45.8399,-119.7006",
  "org": "AS16509 Amazon.com, Inc.",
  "postal": "97818",
  "timezone": "America/Los_Angeles",
  "readme": "https://ipinfo.io/missingauth"
}
```

## License

[MIT](LICENSE)
