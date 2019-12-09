all:
	go build -o asn-search-api main.go

.PHONY: clean
clean:
	rm -rf asn-search-api
