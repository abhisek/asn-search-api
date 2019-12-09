# ASN Search API
A Golang API over MaxMind ASN database

## Build

> Ensure go tool chain is [setup correctly](https://golang.org/doc/install#testing)

```
go build -o asn-search-api main.go
```

## Use

```
./asn-search-api
```

```
curl -s http://localhost:8000/domain/example.com
curl -s http://localhost:8000/org/example+technologies
```

## Deploy

### Zeit

> Setup environment for [Zeit CLI](https://zeit.co/download)

```
now
```