# API Rate Limit
API Rate limit is a controller to limit the number of each IP or token to use the API or any server. In this project use the many famous rate limit algorithm with redis to implement the rate limt controller also provide the easy usage to let's  the developer choose different rate limit algorthm base on he's/she's object.

## Current Provide Algorithm
* Token Bucket Algorithm
* Fix Window Counter Algorithm


## install
```shell
go get github.com/Noahnut/rate_limit
```

## Usage