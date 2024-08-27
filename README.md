# Launchdarkly conding test
I have passed this coding interview with this solution.. enjoy.

# Server-Sent-Events listener

This daemon connects to an URL endpoint streaming Server-Sent-Events in the following format:
```
event: score
data: {"exam": 3, "studentId": "foo", score: .991}
```

## Build
The build process is managed by make
```make```

the resulting binary is under `server`

## Run the server
The server is started with ` $ ./server` and commandline options can be accessed with

```$ ./server -help``` 

## Run test suite
test-script.sh will execute the server, assuming it was build with `make`, and run curl against it.
This scripts assume localhost:8080 is available.
```$ ./test-script.sh```


## Observability
No observability has been implemented, for production readiness i'd add whatever metrics endpoint is supported by the company eg.  
`/metrics` for prometheus exporters  
`/liveness` for k8s readiness/liveness
