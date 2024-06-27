all:
	go run cmd/gochat/main.go

clean:
	rm -rf gochat

.PHONY: all clean