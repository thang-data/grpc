gen:
	protoc --go_out=pb --go-grpc_out=pb      proto/*.proto
clean:
	rm -f proto/processor_message.pb.go
run:
	go run main.go
test:
	go test -cover -race ./...