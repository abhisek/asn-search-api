package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const apiName = "asn-search-api"
const apiVersion = "0.1.0"

const listenHost = "0.0.0.0"
const listenPort = 8000

func initLogger() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func initAsnDB() {
	log.Info("Loading ASN database")
}

func queryAsnOrg(w http.ResponseWriter, r *http.Request) {
	org := mux.Vars(r)["org"]
	log.Infof("Querying ASN for organization: %s", org)
}

func queryAsnDomain(w http.ResponseWriter, r *http.Request) {
	domain := mux.Vars(r)["domain"]
	log.Info("Querying ASN for domain: ", domain)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

func main() {
	initLogger()

	r := mux.NewRouter()

	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/domain/{domain}", queryAsnDomain)
	r.HandleFunc("/org/{org}", queryAsnOrg)

	log.Infof("Starting HTTP server on %s:%d", listenHost, listenPort)
	http.ListenAndServe(fmt.Sprintf("%s:%d", listenHost, listenPort), r)
}
