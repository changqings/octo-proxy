build:
	go build -o tcp-proxy cmd/octo/main.go

clean:
	rm -f tcp-proxy

.PHONY: clean
