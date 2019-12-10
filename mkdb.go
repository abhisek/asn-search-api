package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/blevesearch/bleve"
	log "github.com/sirupsen/logrus"
)

func isValidType(pType string) bool {
	if (pType == TYPE_IPv4) || (pType == TYPE_IPv6) {
		return true
	} else {
		return false
	}
}

func getRecordIndexID(r *AsnRecord) string {
	h := sha1.New()
	h.Write([]byte(fmt.Sprint("%s_%s", r.ID, r.Address)))

	return hex.EncodeToString(h.Sum(nil))
}

func removeQuotes(s string) string {
	if len(s) >= 2 {
		if s[0] == '"' && s[len(s)-1] == '"' {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func createIndexedAsnDB(pDB, pType, pFile *string) {
	index, err := bleve.Open(*pDB)
	if err != nil {
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New(*pDB, mapping)

		if err != nil {
			log.Fatalf("Failed to create index at %s error: %+v", *pDB, err)
		}
	}

	file, err := os.Open(*pFile)
	if err != nil {
		log.Fatalf("Failed to open CSV file: %s", *pFile)
	}

	defer file.Close()
	defer index.Close()

	if !isValidType(*pType) {
		log.Fatalf("Invalid type: %s", *pType)
	}

	scanner := bufio.NewScanner(file)
	count := 0

	log.Infof("Indexing records to DB: %s", *pDB)

	var wg sync.WaitGroup
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")
		record := AsnRecord{ID: parts[1],
			Address:      parts[0],
			Organization: removeQuotes(parts[2]),
			Type:         *pType}

		wg.Add(1)
		go func(r *AsnRecord, id int) {
			defer wg.Done()

			riID := getRecordIndexID(&record)
			log.Infof("Indexing record with ASNID: %s ID: %s", record.ID, riID)

			if err := index.Index(riID, r); err != nil {
				log.Fatalf("Failed to index record. Error: %+v", err)
			}
		}(&record, count+1)
		count = count + 1
	}

	wg.Wait()
	log.Infof("Indexed %d entries", count)
}
