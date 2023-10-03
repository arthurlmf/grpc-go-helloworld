module github.com/menghanl/grpc-go-helloworld

go 1.13

replace google.golang.org/grpc => ../grpc-go

require (
	golang.org/x/net v0.7.0 // indirect
	google.golang.org/grpc v1.23.1
)
