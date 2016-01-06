# Caffeine proxy

[![Build Status](https://travis-ci.org/tsuru/caffeine.png?branch=master)](https://travis-ci.org/tsuru/caffeine)

Proxy to start asleep apps

## Configuration

Some environment variables should be configured:

- `HIPACHE_REDIS_HOST`: Hipache Redis hostname (default: `localhost`)
- `HIPACHE_REDIS_PORT`: Hipache Redis port (default: `6379`)
- `HIPACHE_REDIS_MAXCONN`: Hipache Redis maxconn value (default: `10`)
- `TSURU_HOST`: hostname for Tsuru API (default: `http://localhost`)
- `TSURU_APP_PROXY`: the name of this app in Tsuru. Used to ignore requests that shouldn't be proxied
- `TSURU_TOKEN`: user token, used to make requests to Tsuru API
- `WAIT_BEFORE_PROXY`: time to wait, in seconds, between starting the app and proxying the request (default: `0`)
