# exchain
A simple project to provide a consistent API for getting data from crypto exchanges.

## Build and deploy

The default Dockerfile can be used to deploy to Heroku.

You can also deploy to a local docker instance:
```
docker build -f Dockerfile.local -t <<image>> .

docker run --name <<container>> -p <<port>>:80 -d <<image>>
```

## Usage

### Tickers

Endpoint: /ex/v1/tickers/{base}

Description: Get ticker information for a currency

```
 curl -s http://<<hostname:port>>/ex/v1/tickers/ETH
```