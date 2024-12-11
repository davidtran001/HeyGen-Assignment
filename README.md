# HeyGen Assignment - `hgclient`

## Overview

The `hgclient` package provides a client library to query a video translation status endpoint. Instead of making repeated calls manually and handling logic for timing and retries, you can rely on the `hgclient.Client` to handle the polling. The client polls the server until one of the following states is reached `completed`, `error`, or `timeout`.

It uses an exponential backoff strategy for polling, starting with a short initial interval and doubling it until a maximum wait interval is reached. This helps balance server load without causing significant delays for the client.

## Prerequisites
- Go 1.23.4 or above.

## Installation
To install and add the package to your Go project, run the following command:

`
go get github.com/davidtran001/HeyGen-Assignment/hgclient
`

## How To Use
Once the `hgclient` package is installed into your Go project, you will need to import `hgclient`. For example:

```go
import (
	"github.com/davidtran001/HeyGen-Assignment/hgclient"
)
```

### Creating a Client

Use the `hgclient.NewClient` function to create a new `hgclient.Client`. The parameters of the function allow you to configure the polling strategy and total timeout:

|Parameter 			|Description|
|:--- 				|:---- |
|`baseURL`			| The base URL of the server |
|`initialWait`		| Initial polling interval in seconds |
|`maxWait`			| The maximum polling interval in seconds |
|`maxTotalTime`		| The maximum total time in seconds to wait for a final result (`completed` or `error`)|
|`backoffFactor`	| The factor by which the wait interval is multiplied each time until it reaches `maxWait`|

A full example of how the `hgclient` package can be used is included in this repository as `client_example/client_example.go`.

## Running Tests
This project includes several tests.

### Integration Tests
To run the integration test, run the following command:

`go test -v -run TestIntegrationServer`

### Unit & Integration Tests
To run all unit and integration tests, run the following command:

`go test ./... -v`



