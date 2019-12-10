package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/document"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const apiName = "asn-search-api"
const apiVersion = "0.1.0"

const listenHost = "0.0.0.0"
const listenPort = 8000

const (
	TYPE_IPv4 = "ipv4"
	TYPE_IPv6 = "ipv6"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type AsnRecord struct {
	ID           string `json:"id,omitempty"`
	Address      string `json:"address"`
	Organization string `json:"organization"`
	Type         string `json:"type"`
}

type SuccessResponse struct {
	Count   int         `json:"count"`
	Records []AsnRecord `json:"records"`
}

var pMode = flag.String("mode", "server", "Mode to run (server, mkdb)")
var pDB = flag.String("db", "", "Database path")
var pType = flag.String("type", "", "Type: ipv4/ipv6")
var pFile = flag.String("file", "", "Input CSV file path from MaxMind")

const AsnDBPath = "data/asn.db"

var AsnDB *bleve.Index

func initLogger() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func initAsnDB() {
	log.Info("Loading ASN database")
	db, err := bleve.Open(AsnDBPath)
	if err != nil {
		log.Fatalf("Failed to open ASN DB: %+v", err)
	}

	AsnDB = &db
}

func getFieldValue(fields []document.Field, name string) string {
	var ret string

	for _, f := range fields {
		if f.Name() == name {
			ret = string(f.Value())
			break
		}
	}

	return ret
}

func queryAsnByOrgName(name string) []AsnRecord {
	results := make([]AsnRecord, 0)

	query := bleve.NewMatchQuery(name)
	searchRequest := bleve.NewSearchRequest(query)
	searchResult, err := (*AsnDB).Search(searchRequest)

	if err != nil {
		log.Warnf("Failed to execute search query: %#v", err)
		return results
	}

	for _, hit := range searchResult.Hits {
		doc, _ := (*AsnDB).Document(hit.ID)

		r := AsnRecord{}
		r.ID = getFieldValue(doc.Fields, "ID")
		r.Address = getFieldValue(doc.Fields, "Address")
		r.Organization = getFieldValue(doc.Fields, "Organization")
		r.Type = getFieldValue(doc.Fields, "Type")

		// for _, field := range doc.Fields {
		// 	fmt.Printf("Name: %s Value: %s\n", field.Name(), field.Value())
		// }

		results = append(results, r)
	}

	return results
}

func executeAsnQueryByOrg(org string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if len(org) <= 0 {
		res := ErrorResponse{Error: "error", Message: "Input param invalid"}
		json, _ := json.Marshal(res)

		w.WriteHeader(http.StatusBadRequest)
		w.Write(json)

		return
	}

	records := queryAsnByOrgName(org)
	res := SuccessResponse{
		Count:   len(records),
		Records: records}

	json, _ := json.Marshal(res)

	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func queryAsnOrgHandler(w http.ResponseWriter, r *http.Request) {
	org, _ := url.QueryUnescape(mux.Vars(r)["org"])
	log.Infof("Querying ASN for organization: %s", org)

	executeAsnQueryByOrg(org, w, r)
}

func queryAsnDomainHandler(w http.ResponseWriter, r *http.Request) {
	domain, _ := url.QueryUnescape(mux.Vars(r)["domain"])
	log.Info("Querying ASN for domain: ", domain)

	org := "" // Get org from domain
	executeAsnQueryByOrg(org, w, r)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

func main() {
	initLogger()

	flag.Parse()

	if *pMode == "mkdb" {
		createIndexedAsnDB(pDB, pType, pFile)
		return
	}

	initAsnDB()

	r := mux.NewRouter()

	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/domain/{domain}", queryAsnDomainHandler)
	r.HandleFunc("/org/{org}", queryAsnOrgHandler)

	log.Infof("Starting HTTP server on %s:%d", listenHost, listenPort)
	http.ListenAndServe(fmt.Sprintf("%s:%d", listenHost, listenPort), r)
}
