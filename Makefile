all:
	go build -o asn-search-api

.PHONY: clean
clean:
	rm -rf asn-search-api
