# APKCTL Integration testing

## Pre-requisites for running integration tests
- In order to run the integration tests, a K8s cluster needs to be started. By default integration tests are configured to run against local K8s cluster.


## Executing command

| All commands must be run from *integration* directory

### Flags ###

- Required:

   `-archive` :  apkctl archive file that is to be tested

- Optional:

   `-run` : Run specific test fucntion only

   `-v` : Print verbose test output, useful for debugging

   `-logtransport` : Print http transport level request/responses



### Command ###

Before start running the integration tests, navigate to `CTL` directory and build the executable file by running the following command.
`go build apkctl.go`

- Basic command

```
go test -p 1 -timeout 0 -archive <apkctl archive name>

example: go test -p 1 -timeout 0 -archive apkctl-4.1.0-linux-x64.tar.gz

```

- Run a specific test function only

```
go test -p 1 -timeout 0 -archive <apkctl archive name> -run <Test function name or partial name regex>

example: go test -p 1 -timeout 0 -archive apkctl-4.1.0-linux-x64.tar.gz -run TestVersion
```

- Print verbose output

```
go test -p 1 -timeout 0 -archive <apkctl archive name> -v

example: go test -p 1 -timeout 0 -archive apkctl-4.1.0-linux-x64.tar.gz -v
```

- Print http transport request/responses

```
go test -p 1 -timeout 0 -archive <apkctl archive name> -logtransport

example: go test -p 1 -timeout 0 -archive apkctl-4.1.0-linux-x64.tar.gz -logtransport
```

---
- [1] https://github.com/golang/go/issues/3575
- [2] https://wilsonmar.github.io/maximum-limits/
