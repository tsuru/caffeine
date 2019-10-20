# Caffeine proxy

[![Build Status](https://travis-ci.org/tsuru/caffeine.png?branch=master)](https://travis-ci.org/tsuru/caffeine)

Proxy to start asleep apps



## Documentation

See [documentation](https://godoc.org/github.com/tsuru/caffeine
)
for instructions.


## Configuration

Some environment variables should be configured:

- `CUSTOM_HEADER_VALUE`: if defined, this value is used in a custom header `X-Caffeine` added to the proxied request
- `TSURU_TOKEN`: user token, used to make requests to Tsuru API
- `TSURU_HOST`: hostname for Tsuru API
- `WAIT_BEFORE_PROXY`: time to wait, in seconds, between starting the app and proxying the request (default: `0`)

## Testing

- `GO15VENDOREXPERIMENT`: If using go1.5, ensure you set this to 1
