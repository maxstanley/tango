# tango
Tango is a go api framework wrapper to enable cross-platform deployment of endpoint functions.

## examples

### simple

Simple defines a single endpoint handler, which is then used across multiple cmd's.
Tango enables this single function to be called by each of the api frameworks, by using a wrapper to create a single event for the function.

Simple currently demonstrates the wrapping functionality for gin-gonic, go-chi and aws lambdas.

#### building

```bash
# Compile protobuf definitions.
cd _examples/simple/protobuf
protoc --go_out=. *.proto

# Same process for gin or chi.
cd _examples/simple/cmd/{gin, chi}
go build

# Compile and Zip for lambdas.
cd _examples/simple/cmd/lambda/
GOOS=linux CGO_ENABLED=0 go build main.go && zip function.zip main
```
