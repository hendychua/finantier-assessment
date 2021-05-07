# Finantier assessment

## Overview
This project consists of 2 microservices written in Golang. The first microservice is `encryption-service` and the second microservice is `data-service`.

`encryption-service` is a service that encrypts a payload with AES256 and returns the encrypted data.

`data-service` is a service that takes in a stock symbol, e.g. TSLA, calls a third-party service to get stock quotes for the symbol, and calls the `encryption-service` to encrypt the stock quote data and returns the encrypted data.

### Encryption-service
This project uses go modules to manage its dependencies. Enter the project directory and run `go mod download` to download the dependencies. After that, you can run `go run main.go` to start the service.

This service requires a 32-byte encryption key to start. Before starting the service, set the key in the environment using the variable `AES256_ENC_KEY`.

### Data-service
This project uses go modules to manage its dependencies. Enter the project directory and run `go mod download` to download the dependencies. After that, you can run `go run main.go` to start the service.

This service uses [Alpha Vantage](https://www.alphavantage.co/) to get stock quote. To get started, you need to register for an API key first.

This service requires a few settings to start. Specify these settings in the environment:
- `ENCRYPTION_SERVICE_HOST` - where the encryption service is running, e.g. `http://localhost:8080/`
- `ALPHA_VANTAGE_API_KEY` - your Alpha Vantage API key.

Optional environment variables:
- `SKIP_SSL_VERIFY` - By default this is disabled, meaning SSL certificates are checked when making network requests to Alpha Vantage. Set this to `1` to enable it. This should **only be used in local testing**. This was added only for convenience because the Ubuntu container used in this project does not recognise the CA used to generate the SSL certificate for Alpha Vantage.

## Running with docker-compose
This project uses `docker-compose` to spin up both services together and to allow them to talk to one another. Before spinning them up, you need to build the linux binary and create environment variable file for each of the services.

```
cd data-service
GOOS=linux GOARCH=arm64 go build .

cd ../encryption-service
GOOS=linux GOARCH=arm64 go build .

# go back to project root directory
cd ..
```

In the project root directory:

`encryption-service.env`
```
AES256_ENC_KEY=my-32-byte-key
```

`data-service.env`
```
# encryption-service is the name of the service defined in docker-compose.yml
ENCRYPTION_SERVICE_HOST=http://encryption-service:8080/
ALPHA_VANTAGE_API_KEY=xxxxyyyy
# This is needed for the Ubuntu & SSL cert authority issue mentioned above.
SKIP_SSL_VERIFY=1
```

After both env files are created, you can run `docker-compose up` to spin up both services.

## Usages

```
# Assuming your data-service is running at 8081 locally
curl http://localhost:8081/symbol/TSLA
# some encrypted payload
```
