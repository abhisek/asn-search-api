FROM golang:1.12.3
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o asn-search-api

## TODO: Build DB from CSV

FROM alpine:latest
WORKDIR /app

COPY --from=0 /app/asn-search-api .
RUN mkdir data
COPY data/asn.db ./data/asn.db

EXPOSE 8000

CMD ["./asn-search-api"]