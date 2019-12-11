# ASN Search API
A Golang API over MaxMind ASN database

* The API server requires a pre-built database in `data/asn.db`
* The API server listens on `0.0.0.0:8000` by default (Customize with `HOST` and `PORT` environment variable)
* Building database requires [MaxMind ASN CSV](https://dev.maxmind.com/geoip/geoip2/geolite2-asn-csv-database/)

## Getting Started

```bash
docker run -p 8000:8000 abh1sek/asn-search-api
```

## Use

```bash
curl -s http://localhost:8000/domain/example.com
curl -s http://localhost:8000/org/example+technologies
```

> `OrgName` should be [URL Encoded](https://www.w3schools.com/tags/ref_urlencode.asp)

## Build

> Ensure go tool chain is [setup correctly](https://golang.org/doc/install#testing)

```bash
make
```

## Generate ASN Database

1. Download MaxMind ASN Database in CSV Format
2. Use `asn-search-api` tool to create indexed database

```bash
./asn-search-api -mode mkdb \
  -db data/asn.db \
  -type ipv4 \
  -file GeoLite2-ASN-CSV_20190101/GeoLite2-ASN-Blocks-IPv4.csv 

./asn-search-api -mode mkdb \
  -db data/asn.db \
  -type ipv6 \
  -file GeoLite2-ASN-CSV_20190101/GeoLite2-ASN-Blocks-IPv6.csv 
```

## Deploy

### Google Cloud Run

```bash
gcloud run deploy \
  asn-search-api \
  --platform=managed \
  --image=gcr.io/<projectName>/asn-search-api:latest \
  --memory=512Mi \
  --timeout=30s \
  --labels=app=asn-search-api \
  --allow-unauthenticated \
  --region=us-central1
```

## TODO

- [ ] Extract ASN DB ops from `main.go` and create its own package
- [ ] Test cases