# API Rate Limit
API Rate limit is a controller to limit the number of each IP or token to use the API or any server. In this project use the many famous rate limit algorithm with redis to implement the rate limt controller also provide the easy usage to let's  the developer choose different rate limit algorthm base on he's/she's object.

## Current Provide Algorithm
* Token Bucket Algorithm  
**Pros**  
  Good for memory usage because limit the bucket number.  
**Cons**  
  Hard to adjust the proper bucket size and refill speed.
* Fix Window Counter Algorithm  
**Pros**  
Good for memory usage
Easy to understand  
**Cons**        
if the request increasing when the interval change could couse the more request be allowed.
* Slide Window Log Algorithm
**Pros**  
Very accurate solve the Fix Window Counter Algorithm problem
**Cons**        
Need more memory to record the request time stamp


## install
```shell
go get github.com/Noahnut/rate_limit
```

## Usage