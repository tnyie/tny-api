# TnyAPI

[![Go Reference](https://pkg.go.dev/badge/github.com/tnyie/tny-api.svg)](https://pkg.go.dev/github.com/tnyie/tny-api)


A link shortener, backend for [Tny.ie](https://tny.ie) and its associated browser extension. 

## TODO
- Add unit tests
- Switch logging library for better logging
  - Also change how much is logged & remove stale and leaky logs

### Future possible features
- Support user-owned domains
- Use subdomains for user-created websites/link trees perhaps?
  - e.g. thomas.tny.ie may have links for social medias, etc.

# Usage

## Docker

A dockerfile is incuded, an image will be available eventually

An example docker-compose file is in the \_\_example\_\_ directory

### Environment

| ENVIRONMENT_VARIABLE | Usage                                  |
|----------------------|----------------------------------------|
| TNY_AUTH_KEY         | Signing key used for cookie's sessions |
| DB_USER              | Database user                          |
| DB_PASS              | Database password                      |
| DB_NAME              | Database name                          |
| DB_HOST              | Database host                          |
| DB_PORT              | Database port                          |
