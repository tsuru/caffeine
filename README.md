# Caffeine proxy

[![Build Status](https://travis-ci.org/tsuru/caffeine.png?branch=master)](https://travis-ci.org/tsuru/caffeine)

Proxy to start asleep apps

## Configuration

Some environment variables should be configured:

- `TSURU_APP_PROXY`: the name of this app in Tsuru. Used to ignore requests that shouldn't be proxied
- `TSURU_TOKEN`: user token, used to make requests to Tsuru API
- `WAIT_BEFORE_PROXY`: time to wait, in seconds, between starting the app and proxying the request (default: `0`)

## Testing

- `GO15VENDOREXPERIMENT`: If using go1.5, ensure you set this to 1
