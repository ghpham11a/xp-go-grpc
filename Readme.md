# 1. Initialize project

```
export PATH="$PATH:$(go env GOPATH)/bin"
```

```
go mod init xp-go-grpc
go mod tidy
```

```
# Install protoc-gen-go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# Install protoc-gen-go-grpc
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

```
go get google.golang.org/grpc
```

```
make sure Go sees where you installed the binaries
```

Compile protobuf files


Generate protobuf files for Client and server seperatly 

```
protoc \
  -I=. \
  -I=/Users/user/Documents/googleapis \
  --go_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_out=. \
  --go-grpc_opt=paths=source_relative \
  proto/accounts.proto
```

# 2. Run project

```
go run server/main.go
```

In another terminal

```
go run client/main.go --name="Luke"
```