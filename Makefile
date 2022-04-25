gen:
	protoc --go_out=pb --go-grpc_out=pb proto/*.proto
clean:
	rm -f proto/processor_message.pb.go
server:
	go run cmd/server/main.go -port 1323
client:
	go run cmd/client/main.go -address 0.0.0.0:1323
test:
	go test -cover -race ./...