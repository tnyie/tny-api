# [Tny](https://tny.ie)

[![Go Reference](https://pkg.go.dev/badge/github.com/tnyie/tny-api.svg)](https://pkg.go.dev/github.com/tnyie/tny-api)

A link shortener with some extra kick

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
