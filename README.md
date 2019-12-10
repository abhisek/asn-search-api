# ASN Search API
A Golang API over MaxMind ASN database

## Build

> Ensure go tool chain is [setup correctly](https://golang.org/doc/install#testing)

```
make
```

## Generate ASN Database

1. Download MaxMind ASN Database in CSV Format
2. Use `asn-search-api` tool to create indexed database

```
./asn-search-api -mode mkdb \
  -db data/asn.db \
  -type ipv4 \
  -file GeoLite2-ASN-CSV_20190101/GeoLite2-ASN-Blocks-IPv4.csv 

./asn-search-api -mode mkdb \
  -db data/asn.db \
  -type ipv6 \
  -file GeoLite2-ASN-CSV_20190101/GeoLite2-ASN-Blocks-IPv6.csv 
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

> Ensure database is generated in `data/asn.db`

> Setup environment for [Zeit CLI](https://zeit.co/download)

```
now
```

